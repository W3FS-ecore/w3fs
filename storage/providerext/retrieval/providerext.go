package retrievalimpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/storage/event"
	logging "github.com/ipfs/go-log/v2"
	"sync"
	"time"

	"github.com/filecoin-project/go-address"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-statemachine/fsm"
	"github.com/hannahhoward/go-pubsub"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-fil-markets/piecestore"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/askstore"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/dtutils"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/providerstates"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/requestvalidation"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/migrations"
	rmnet "github.com/filecoin-project/go-fil-markets/retrievalmarket/network"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/go-fil-markets/stores"
)

var log = logging.Logger("retrieval-providerext")

// RetrievalProviderOption is a function that configures a retrieval provider
type RetrievalProviderOption func(p *ProviderExt)

// DealDecider is a function that makes a decision about whether to accept a deal
type DealDecider func(ctx context.Context, state retrievalmarket.ProviderDealState) (bool, string, error)

type RetrievalPricingFunc func(ctx context.Context, dealPricingParams retrievalmarket.PricingInput) (retrievalmarket.Ask, error)

var queryTimeout = 5 * time.Second

// Provider is the production implementation of the RetrievalProvider interface
type ProviderExt struct {
	ctx              context.Context
	dataTransfer     datatransfer.Manager
	node             retrievalmarket.RetrievalProviderNode
	sa               retrievalmarket.SectorAccessor
	network          rmnet.RetrievalMarketNetwork
	requestValidator *requestvalidation.ProviderRequestValidator
	revalidator      *requestvalidation.ProviderRevalidator
	minerAddress     address.Address
	pieceStore       piecestore.PieceStore
	readySub         *pubsub.PubSub
	subscribers      *pubsub.PubSub
	//stateMachines        fsm.Group
	receiver             event.DataTransferEventReceiver
	migrateStateMachines func(context.Context) error
	dealDecider          DealDecider
	askStore             retrievalmarket.AskStore
	disableNewDeals      bool
	retrievalPricingFunc RetrievalPricingFunc
	dagStore             stores.DAGStoreWrapper
	stores               *stores.ReadOnlyBlockstores
}

type internalProviderEvent struct {
	evt   retrievalmarket.ProviderEvent
	state retrievalmarket.ProviderDealState
}

func providerDispatcher(evt pubsub.Event, subscriberFn pubsub.SubscriberFn) error {
	ie, ok := evt.(internalProviderEvent)
	if !ok {
		return errors.New("wrong type of event")
	}
	cb, ok := subscriberFn.(retrievalmarket.ProviderSubscriber)
	if !ok {
		return errors.New("wrong type of event")
	}
	log.Debugw("process retrieval provider] listeners", "name", retrievalmarket.ProviderEvents[ie.evt], "proposal cid", ie.state.ID)
	cb(ie.evt, ie.state)
	return nil
}

// DealDeciderOpt sets a custom protocol
func DealDeciderOpt(dd DealDecider) RetrievalProviderOption {
	return func(provider *ProviderExt) {
		provider.dealDecider = dd
	}
}

// DisableNewDeals disables setup for v1 deal protocols
func DisableNewDeals() RetrievalProviderOption {
	return func(provider *ProviderExt) {
		provider.disableNewDeals = true
	}
}

// NewProvider returns a new retrieval ProviderExt
func NewProviderExt(
	minerAddress address.Address,
	node retrievalmarket.RetrievalProviderNode,
	sa retrievalmarket.SectorAccessor,
	network rmnet.RetrievalMarketNetwork,
	pieceStore piecestore.PieceStore,
	dagStore stores.DAGStoreWrapper,
	dataTransfer datatransfer.Manager,
	ds datastore.Batching,
	retrievalPricingFunc RetrievalPricingFunc,
	opts ...RetrievalProviderOption,
) (retrievalmarket.RetrievalProvider, error) {

	if retrievalPricingFunc == nil {
		//return nil, xerrors.New("retrievalPricingFunc is nil")
		log.Warnf("retrievalPricingFunc is null")
	}

	p := &ProviderExt{
		dataTransfer:         dataTransfer,
		node:                 node,
		sa:                   sa,
		network:              network,
		minerAddress:         minerAddress,
		pieceStore:           pieceStore,
		subscribers:          pubsub.New(providerDispatcher),
		readySub:             pubsub.New(shared.ReadyDispatcher),
		receiver:             event.DataTransferEventReceiver{Sms: sync.Map{}},
		retrievalPricingFunc: retrievalPricingFunc,
		dagStore:             dagStore,
		stores:               stores.NewReadOnlyBlockstores(),
	}

	err := shared.MoveKey(ds, "retrieval-ask", "retrieval-ask/latest")
	if err != nil {
		return nil, err
	}

	askStore, err := askstore.NewAskStore(namespace.Wrap(ds, datastore.NewKey("retrieval-ask")), datastore.NewKey("latest"))
	if err != nil {
		return nil, err
	}
	p.askStore = askStore

	p.Configure(opts...)
	p.requestValidator = requestvalidation.NewProviderRequestValidator(&providerExtValidationEnvironment{p})
	transportConfigurer := dtutils.TransportConfigurer(network.ID(), &providerStoreGetter{p})
	p.revalidator = requestvalidation.NewProviderRevalidator(&providerRevalidatorEnvironment{p})

	if p.disableNewDeals {
		err = p.dataTransfer.RegisterVoucherType(&migrations.DealProposal0{}, p.requestValidator)
		if err != nil {
			return nil, err
		}
		err = p.dataTransfer.RegisterRevalidator(&migrations.DealPayment0{}, p.revalidator)
		if err != nil {
			return nil, err
		}
	} else {
		err = p.dataTransfer.RegisterVoucherType(&retrievalmarket.DealProposal{}, p.requestValidator)
		if err != nil {
			return nil, err
		}
		err = p.dataTransfer.RegisterVoucherType(&migrations.DealProposal0{}, p.requestValidator)
		if err != nil {
			return nil, err
		}

		err = p.dataTransfer.RegisterRevalidator(&retrievalmarket.DealPayment{}, p.revalidator)
		if err != nil {
			return nil, err
		}
		err = p.dataTransfer.RegisterRevalidator(&migrations.DealPayment0{}, requestvalidation.NewLegacyRevalidator(p.revalidator))
		if err != nil {
			return nil, err
		}

		err = p.dataTransfer.RegisterVoucherResultType(&retrievalmarket.DealResponse{})
		if err != nil {
			return nil, err
		}

		err = p.dataTransfer.RegisterTransportConfigurer(&retrievalmarket.DealProposal{}, transportConfigurer)
		if err != nil {
			return nil, err
		}
	}
	err = p.dataTransfer.RegisterVoucherResultType(&migrations.DealResponse0{})
	if err != nil {
		return nil, err
	}
	err = p.dataTransfer.RegisterTransportConfigurer(&migrations.DealProposal0{}, transportConfigurer)
	if err != nil {
		return nil, err
	}
	providerRevalidationEnv.p = p
	providerDealEnv.p = p
	dataTransfer.SubscribeToEvents(ProviderDataTransferSubscriberExt(&p.ctx, &dataTransfer))
	return p, nil
}

func dealProposalFromVoucher(voucher datatransfer.Voucher) (*retrievalmarket.DealProposal, bool) {
	dealProposal, ok := voucher.(*retrievalmarket.DealProposal)
	// if this event is for a transfer not related to storage, ignore
	if ok {
		return dealProposal, true
	}

	legacyProposal, ok := voucher.(*migrations.DealProposal0)
	if !ok {
		return nil, false
	}
	newProposal := migrations.MigrateDealProposal0To1(*legacyProposal)
	return &newProposal, true
}

// DataTransfer Subscriber func
func ProviderDataTransferSubscriberExt(ctx *context.Context, dataTransfer *datatransfer.Manager) datatransfer.Subscriber {
	return func(event datatransfer.Event, channelState datatransfer.ChannelState) {

		dealProposal, ok := dealProposalFromVoucher(channelState.Voucher())
		// if this event is for a transfer not related to storage, ignore
		if !ok {
			return
		}
		log.Infow("processing retrieval provider dt event", "event", datatransfer.Events[event.Code], "peer", channelState.OtherPeer(), "channelState", channelState)

		if event.Code == datatransfer.Accept {
			dealId := retrievalmarket.ProviderDealIdentifier{DealID: dealProposal.ID, Receiver: channelState.Recipient()}
			log.Infof("dealId:%s", dealId)
			dealState, err := providerRevalidationEnv.Get(dealId)
			if err != nil {
				log.Errorf("Get ProviderDealState error! dealId:%s, error:%s", dealId, err)
				return
			}
			// set channelId
			chid := channelState.ChannelID()
			dealState.ChannelID = &chid
			log.Infof("dealId:%s  dealState's Status:%s", dealId, dealState.Status)
			err = providerDealEnv.UnsealData(*ctx, dealState)
			if err != nil {
				log.Errorf("provider UnsealData error! dealId:%s, dealState:%s error:%s", dealId, dealState, err)
				closeChannelWhenError(ctx, dataTransfer, chid)
				return
			}
			err = providerDealEnv.UnpauseDeal(*ctx, dealState)
			if err != nil {
				log.Errorf("provider UnpauseDeal error! dealId:%s, dealState:%s error:%s", dealId, dealState, err)
				closeChannelWhenError(ctx, dataTransfer, chid)
				return
			}
		}
	}
}

// when error occurs,close CloseDataTransferChannel,and then other peer will receives 'cancel' event.
func closeChannelWhenError(ctx *context.Context, dataTransfer *datatransfer.Manager, chid datatransfer.ChannelID) {
	(*dataTransfer).CloseDataTransferChannel(*ctx, chid)
}

func (p *ProviderExt) ListDeals() map[retrievalmarket.ProviderDealIdentifier]retrievalmarket.ProviderDealState {
	panic("implement me")
}

// Stop stops handling incoming requests.
func (p *ProviderExt) Stop() error {
	return p.network.StopHandlingRequests()
}

// Start begins listening for deals on the given host.
// Start must be called in order to accept incoming deals.
func (p *ProviderExt) Start(ctx context.Context) error {
	p.ctx = ctx
	err := p.network.SetDelegate(p)
	if err != nil {
		return err
	}
	return nil
}

// OnReady registers a listener for when the provider has finished starting up
func (p *ProviderExt) OnReady(ready shared.ReadyFunc) {
	p.readySub.Subscribe(ready)
}

func (p *ProviderExt) notifySubscribers(eventName fsm.EventName, state fsm.StateType) {
	evt := eventName.(retrievalmarket.ProviderEvent)
	ds := state.(retrievalmarket.ProviderDealState)
	_ = p.subscribers.Publish(internalProviderEvent{evt, ds})
}

// SubscribeToEvents listens for events that happen related to client retrievals
func (p *ProviderExt) SubscribeToEvents(subscriber retrievalmarket.ProviderSubscriber) retrievalmarket.Unsubscribe {
	return retrievalmarket.Unsubscribe(p.subscribers.Subscribe(subscriber))
}

// GetAsk returns the current deal parameters this provider accepts
func (p *ProviderExt) GetAsk() *retrievalmarket.Ask {
	return p.askStore.GetAsk()
}

// SetAsk sets the deal parameters this provider accepts
func (p *ProviderExt) SetAsk(ask *retrievalmarket.Ask) {

	err := p.askStore.SetAsk(ask)

	if err != nil {
		log.Warnf("Error setting retrieval ask: %w", err)
	}
}

//// ListDeals lists all known retrieval deals
//func (p *ProviderExt) ListDeals() map[retrievalmarket.ProviderDealIdentifier]retrievalmarket.ProviderDealState {
//	var deals []retrievalmarket.ProviderDealState
//	_ = p.stateMachines.List(&deals)
//	dealMap := make(map[retrievalmarket.ProviderDealIdentifier]retrievalmarket.ProviderDealState)
//	for _, deal := range deals {
//		dealMap[retrievalmarket.ProviderDealIdentifier{Receiver: deal.Receiver, DealID: deal.ID}] = deal
//	}
//	return dealMap
//}

/*
HandleQueryStream is called by the network implementation whenever a new message is received on the query protocol

A Provider handling a retrieval `Query` does the following:

1. Get the node's chain head in order to get its miner worker address.

2. Look in its piece store to determine if it can serve the given payload CID.

3. Combine these results with its existing parameters for retrieval deals to construct a `retrievalmarket.QueryResponse` struct.

4. Writes this response to the `Query` stream.

The connection is kept open only as long as the query-response exchange.
*/
func (p *ProviderExt) HandleQueryStream(stream rmnet.RetrievalQueryStream) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryTimeout)
	defer cancel()

	defer stream.Close()
	query, err := stream.ReadQuery()
	if err != nil {
		return
	}

	sendResp := func(resp retrievalmarket.QueryResponse) {
		if err := stream.WriteQueryResponse(resp); err != nil {
			log.Errorf("Retrieval query: writing query response: %s", err)
		}
	}

	answer := retrievalmarket.QueryResponse{
		Status:          retrievalmarket.QueryResponseUnavailable,
		PieceCIDFound:   retrievalmarket.QueryItemUnavailable,
		MinPricePerByte: big.Zero(),
		UnsealPrice:     big.Zero(),
	}

	// fetch the piece from which the payload will be retrieved.
	// if user has specified the Piece in the request, we use that.
	// Otherwise, we prefer a Piece which can retrieved from an unsealed sector.
	pieceCID := cid.Undef
	if query.PieceCID != nil {
		pieceCID = *query.PieceCID
	}
	pieceInfo, _, err := p.getPieceInfoFromCid(ctx, query.PayloadCID, pieceCID)
	if err != nil {
		log.Errorf("Retrieval query: getPieceInfoFromCid: %s", err)
		if !xerrors.Is(err, retrievalmarket.ErrNotFound) {
			answer.Status = retrievalmarket.QueryResponseError
			answer.Message = fmt.Sprintf("failed to fetch piece to retrieve from: %s", err)
		}

		sendResp(answer)
		return
	}

	answer.Status = retrievalmarket.QueryResponseAvailable
	answer.Size = uint64(pieceInfo.Deals[0].Length.Unpadded()) // TODO: verify on intermediate
	answer.PieceCIDFound = retrievalmarket.QueryItemAvailable

	/**
	storageDeals, err := storageDealsForPiece(query.PieceCID != nil, query.PayloadCID, pieceInfo, p.pieceStore)
	if err != nil {
		log.Errorf("Retrieval query: storageDealsForPiece: %s", err)
		answer.Status = retrievalmarket.QueryResponseError
		answer.Message = fmt.Sprintf("failed to fetch storage deals containing payload: %s", err)
		sendResp(answer)
		return
	}

	input := retrievalmarket.PricingInput{
		// piece from which the payload will be retrieved
		// If user hasn't given a PieceCID, we try to choose an unsealed piece in the call to `getPieceInfoFromCid` above.
		PieceCID: pieceInfo.PieceCID,

		PayloadCID: query.PayloadCID,
		Unsealed:   isUnsealed,
		Client:     stream.RemotePeer(),
	}*/

	answer.MinPricePerByte = abi.NewTokenAmount(0)
	answer.MaxPaymentInterval = 1 << 30
	answer.MaxPaymentIntervalIncrease = 0
	answer.UnsealPrice = abi.NewTokenAmount(0)
	addr, err := address.NewFromString("t01000")
	answer.PaymentAddress = addr
	sendResp(answer)
}

func (p *ProviderExt) getPieceInfoFromCid(ctx context.Context, payloadCID, pieceCID cid.Cid) (piecestore.PieceInfo, bool, error) {
	cidInfo, err := p.pieceStore.GetCIDInfo(payloadCID)
	if err != nil {
		log.Error("get cid info error : ", err)
		return piecestore.PieceInfoUndefined, false, xerrors.Errorf("get cid info: %w", err)
	}
	var lastErr error
	var sealedPieceInfo *piecestore.PieceInfo

	for _, pieceBlockLocation := range cidInfo.PieceBlockLocations {
		pieceInfo, err := p.pieceStore.GetPieceInfo(pieceBlockLocation.PieceCID)
		if err != nil {
			lastErr = err
			continue
		}

		// if client wants to retrieve the payload from a specific piece, just return that piece.
		if pieceCID.Defined() && pieceInfo.PieceCID.Equals(pieceCID) {
			return pieceInfo, p.pieceInUnsealedSector(ctx, pieceInfo), nil
		}

		// if client dosen't have a preference for a particular piece, prefer a piece
		// for which an unsealed sector exists.
		if pieceCID.Equals(cid.Undef) {
			if p.pieceInUnsealedSector(ctx, pieceInfo) {
				return pieceInfo, true, nil
			}

			if sealedPieceInfo == nil {
				sealedPieceInfo = &pieceInfo
			}
		}

	}

	if sealedPieceInfo != nil {
		return *sealedPieceInfo, false, nil
	}

	if lastErr == nil {
		log.Error("unknown pieceCID : ", pieceCID.String())
		lastErr = xerrors.Errorf("unknown pieceCID %s", pieceCID.String())
	}

	return piecestore.PieceInfoUndefined, false, xerrors.Errorf("could not locate piece: %w", lastErr)
}

func (p *ProviderExt) pieceInUnsealedSector(ctx context.Context, pieceInfo piecestore.PieceInfo) bool {
	for _, di := range pieceInfo.Deals {
		isUnsealed, err := p.sa.IsUnsealed(ctx, di.SectorID, di.Offset.Unpadded(), di.Length.Unpadded())
		if err != nil {
			log.Errorf("failed to find out if sector %d is unsealed, err=%s", di.SectorID, err)
			continue
		}
		if isUnsealed {
			return true
		}
	}

	return false
}

// GetDynamicAsk quotes a dynamic price for the retrieval deal by calling the user configured
// dynamic pricing function. It passes the static price parameters set in the Ask Store to the pricing function.
func (p *ProviderExt) GetDynamicAsk(ctx context.Context, input retrievalmarket.PricingInput, storageDeals []abi.DealID) (retrievalmarket.Ask, error) {
	dp, err := p.node.GetRetrievalPricingInput(ctx, input.PieceCID, storageDeals)
	if err != nil {
		return retrievalmarket.Ask{}, xerrors.Errorf("GetRetrievalPricingInput: %s", err)
	}
	// currAsk cannot be nil as we initialize the ask store with a default ask.
	// Users can then change the values in the ask store using SetAsk but not remove it.
	currAsk := p.GetAsk()
	if currAsk == nil {
		return retrievalmarket.Ask{}, xerrors.New("no ask configured in ask-store")
	}

	dp.PayloadCID = input.PayloadCID
	dp.PieceCID = input.PieceCID
	dp.Unsealed = input.Unsealed
	dp.Client = input.Client
	dp.CurrentAsk = *currAsk

	ask, err := p.retrievalPricingFunc(ctx, dp)
	if err != nil {
		return retrievalmarket.Ask{}, xerrors.Errorf("retrievalPricingFunc: %w", err)
	}
	return ask, nil
}

// Configure reconfigures a provider after initialization
func (p *ProviderExt) Configure(opts ...RetrievalProviderOption) {
	for _, opt := range opts {
		opt(p)
	}
}

// ProviderFSMParameterSpec is a valid set of parameters for a provider FSM - used in doc generation
var ProviderFSMParameterSpec = fsm.Parameters{
	Environment:     &providerDealEnvironment{},
	StateType:       retrievalmarket.ProviderDealState{},
	StateKeyField:   "Status",
	Events:          providerstates.ProviderEvents,
	StateEntryFuncs: providerstates.ProviderStateEntryFuncs,
}

// DefaultPricingFunc is the default pricing policy that will be used to price retrieval deals.
var DefaultPricingFunc = func(VerifiedDealsFreeTransfer bool) func(ctx context.Context, pricingInput retrievalmarket.PricingInput) (retrievalmarket.Ask, error) {
	return func(ctx context.Context, pricingInput retrievalmarket.PricingInput) (retrievalmarket.Ask, error) {
		ask := pricingInput.CurrentAsk

		// don't charge for Unsealing if we have an Unsealed copy.
		if pricingInput.Unsealed {
			ask.UnsealPrice = big.Zero()
		}

		// don't charge for data transfer for verified deals if it's been configured to do so.
		if pricingInput.VerifiedDeal && VerifiedDealsFreeTransfer {
			ask.PricePerByte = big.Zero()
		}

		return ask, nil
	}
}
