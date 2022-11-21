package retrievalimpl2

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/storage/event"
	"github.com/ethereum/go-ethereum/storage/mock"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"os"
	"path/filepath"
	"sync"

	"github.com/hannahhoward/go-pubsub"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	logging "github.com/ipfs/go-log/v2"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-statemachine/fsm"

	"crypto/rand"
	"encoding/json"
	storageclientext "github.com/ethereum/go-ethereum/storage/clientext"
	"github.com/filecoin-project/go-fil-markets/discovery"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/clientstates"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/dtutils"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/migrations"
	rmnet "github.com/filecoin-project/go-fil-markets/retrievalmarket/network"
	"github.com/filecoin-project/go-fil-markets/shared"
	bigext "github.com/filecoin-project/go-state-types/big"
	math "math/big"
)

var log = logging.Logger("retrieval")
var storagenetwork *storageclientext.ClientEx

const (
	// Init
	DEAL_INIT    int = 0
	DEAL_ING     int = 100
	DEAL_SUCCESS int = 200
	DEAL_TIMEOUT int = 400
	DEAL_ERROR   int = 500
)

// Client is the production implementation of the RetrievalClient interface
type ClientEx struct {
	network              rmnet.RetrievalMarketNetwork
	dataTransfer         datatransfer.Manager
	node                 retrievalmarket.RetrievalClientNode
	dealIDGen            *shared.TimeCounter
	Environment          clientDealEnvironment
	subscribers          *pubsub.PubSub
	readySub             *pubsub.PubSub
	resolver             discovery.PeerResolver
	stateMachines        fsm.Group
	migrateStateMachines func(context.Context) error
	Ds                   datastore.Batching
	bstores              retrievalmarket.BlockstoreAccessor
	receiver             event.DataTransferEventReceiver
	// Guards concurrent access to Retrieve method
	retrieveLk  sync.Mutex
	mockFullApi mock.MockFullNode
}

type internalEvent struct {
	evt   retrievalmarket.ClientEvent
	state retrievalmarket.ClientDealState
}

type RetrieveType struct {
	Cid      string `json:"Cid"`
	Headflag bool   `json:"Headflag"`
	Auth     bool   `json:"Auth"`
	Txhash   string `json:"Txhash"`
}

type AuthStatus struct {
	Txhash string `json:"Txhash"`
	Status int    `json:"Status"`
}

type GlobalData struct {
	Savepath string
}

func dispatcher(evt pubsub.Event, subscriberFn pubsub.SubscriberFn) error {
	ie, ok := evt.(internalEvent)
	if !ok {
		return errors.New("wrong type of event")
	}
	cb, ok := subscriberFn.(retrievalmarket.ClientSubscriber)
	if !ok {
		return errors.New("wrong type of event")
	}
	log.Debugw("process retrieval ClientEx listeners", "name", retrievalmarket.ClientEvents[ie.evt], "proposal cid", ie.state.ID)
	cb(ie.evt, ie.state)
	return nil
}

var _ retrievalmarket.RetrievalClient = &ClientEx{}

const dealStartBufferHours uint64 = 8 * 24

// NewClientEx creates a new retrieval ClientEx
func NewClientEx(
	network rmnet.RetrievalMarketNetwork,
	dataTransfer datatransfer.Manager,
	node retrievalmarket.RetrievalClientNode,
	resolver discovery.PeerResolver,
	ds datastore.Batching,
	ba retrievalmarket.BlockstoreAccessor,
	fullApi mock.MockFullNode,
) (retrievalmarket.RetrievalClient, error) {
	c := &ClientEx{
		network:      network,
		dataTransfer: dataTransfer,
		node:         node,
		resolver:     resolver,
		dealIDGen:    shared.NewTimeCounter(),
		subscribers:  pubsub.New(dispatcher),
		readySub:     pubsub.New(shared.ReadyDispatcher),
		receiver:     event.DataTransferEventReceiver{Sms: sync.Map{}},
		bstores:      ba,
		Ds:           ds,
		mockFullApi:  fullApi,
	}

	c.Environment = clientDealEnvironment{c}
	//retrievalMigrations, err := migrations.ClientMigrations.Build()
	//if err != nil {
	//	return nil, err
	//}
	//c.stateMachines, c.migrateStateMachines, err = versionedfsm.NewVersionedFSM(ds, fsm.Parameters{
	//	Environment:     &clientDealEnvironment{c},
	//	StateType:       retrievalmarket.ClientDealState{},
	//	StateKeyField:   "Status",
	//	Events:          clientstates.ClientEvents,
	//	StateEntryFuncs: clientstates.ClientStateEntryFuncs,
	//	FinalityStates:  clientstates.ClientFinalityStates,
	//	Notifier:        c.notifySubscribers,
	//}, retrievalMigrations, "2")
	//if err != nil {
	//	return nil, err
	//}
	err := dataTransfer.RegisterVoucherResultType(&retrievalmarket.DealResponse{})
	if err != nil {
		return nil, err
	}
	err = dataTransfer.RegisterVoucherResultType(&migrations.DealResponse0{})
	if err != nil {
		return nil, err
	}
	err = dataTransfer.RegisterVoucherType(&retrievalmarket.DealProposal{}, nil)
	if err != nil {
		return nil, err
	}
	err = dataTransfer.RegisterVoucherType(&migrations.DealProposal0{}, nil)
	if err != nil {
		return nil, err
	}
	err = dataTransfer.RegisterVoucherType(&retrievalmarket.DealPayment{}, nil)
	if err != nil {
		return nil, err
	}
	err = dataTransfer.RegisterVoucherType(&migrations.DealPayment0{}, nil)
	if err != nil {
		return nil, err
	}
	dataTransfer.SubscribeToEvents(ClientDataTransferSubscriberExt(c))
	transportConfigurer := dtutils.TransportConfigurer(network.ID(), &clientStoreGetter{c})
	err = dataTransfer.RegisterTransportConfigurer(&retrievalmarket.DealProposal{}, transportConfigurer)
	if err != nil {
		return nil, err
	}
	err = dataTransfer.RegisterTransportConfigurer(&migrations.DealProposal0{}, transportConfigurer)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ClientEx) NextID() retrievalmarket.DealID {
	return retrievalmarket.DealID(c.dealIDGen.Next())
}

// Start initialized the ClientEx, performing relevant database migrations
func (c *ClientEx) Start(ctx context.Context) error {
	//go func() {
	//	err := c.migrateStateMachines(ctx)
	//	if err != nil {
	//		log.Errorf("Migrating retrieval ClientEx state machines: %s", err.Error())
	//	}
	//
	//	err = c.readySub.Publish(err)
	//	if err != nil {
	//		log.Warnf("Publish retrieval ClientEx ready event: %s", err.Error())
	//	}
	//}()
	return nil
}

// OnReady registers a listener for when the ClientEx has finished starting up
func (c *ClientEx) OnReady(ready shared.ReadyFunc) {
	c.readySub.Subscribe(ready)
}

// FindProviders uses PeerResolver interface to locate a list of providers who may have a given payload CID.
func (c *ClientEx) FindProviders(payloadCID cid.Cid) []retrievalmarket.RetrievalPeer {
	peers, err := c.resolver.GetPeers(payloadCID)
	if err != nil {
		log.Errorf("failed to get peers: %s", err)
		return []retrievalmarket.RetrievalPeer{}
	}
	return peers
}

/*
#no use#
Query sends a retrieval query to a specific retrieval provider, to determine
if the provider can serve a retrieval request and what its specific parameters for
the request are.

The ClientEx creates a new `RetrievalQueryStream` for the chosen peer ID,
and calls `WriteQuery` on it, which constructs a data-transfer message and writes it to the Query stream.
*/
func (c *ClientEx) Query(ctx context.Context, p retrievalmarket.RetrievalPeer, payloadCID cid.Cid, params retrievalmarket.QueryParams) (retrievalmarket.QueryResponse, error) {
	s, err := c.network.NewQueryStream(p.ID)
	if err != nil {
		log.Warn(err)
		return retrievalmarket.QueryResponseUndefined, err
	}
	defer s.Close()

	err = s.WriteQuery(retrievalmarket.Query{
		PayloadCID:  payloadCID,
		QueryParams: params,
	})
	if err != nil {
		log.Warn(err)
		return retrievalmarket.QueryResponseUndefined, err
	}

	return s.ReadQueryResponse()
}

// Retrieve initiates the retrieval deal flow, which involves multiple requests and responses
//
// To start this processes, the ClientEx creates a new `RetrievalDealStream`.  Currently, this connection is
// kept open through the entire deal until completion or failure.  Make deals pauseable as well as surviving
// a restart is a planned future feature.
//
// Retrieve should be called after using FindProviders and Query are used to identify an appropriate provider to
// retrieve the deal from. The parameters identified in Query should be passed to Retrieve to ensure the
// greatest likelihood the provider will accept the deal
//
// When called, the ClientEx takes the following actions:
//
// 1. Creates a deal ID using the next value from its `storedCounter`.
//
// 2. Constructs a `DealProposal` with deal terms
//
// 3. Tells its statemachine to begin tracking this deal state by dealID.
//
// 4. Constructs a `blockio.SelectorVerifier` and adds it to its dealID-keyed map of block verifiers.
//
// 5. Triggers a `ClientEventOpen` event on its statemachine.
//
// From then on, the statemachine controls the deal flow in the ClientEx. Other components may listen for events in this flow by calling
// `SubscribeToEvents` on the ClientEx. The ClientEx handles consuming blocks it receives from the provider, via `ConsumeBlocks` function
//
// Retrieve can use an ID generated through NextID, or can generate an ID if the user passes a zero value.
//
// Use NextID when it's necessary to reserve an ID ahead of time, e.g. to
// associate it with a given blockstore in the BlockstoreAccessor.
//
// Documentation of the ClientEx state machine can be found at https://godoc.org/github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/clientstates
func (c *ClientEx) Retrieve(
	ctx context.Context,
	id retrievalmarket.DealID,
	payloadCID cid.Cid,
	params retrievalmarket.Params,
	totalFunds abi.TokenAmount,
	p retrievalmarket.RetrievalPeer,
	clientWallet address.Address,
	minerWallet address.Address,
) (retrievalmarket.DealID, error) {
	c.retrieveLk.Lock()
	defer c.retrieveLk.Unlock()

	// assign a new ID.
	if id == 0 {
		next := c.dealIDGen.Next()
		id = retrievalmarket.DealID(next)
	}

	dealState := retrievalmarket.ClientDealState{
		DealProposal: retrievalmarket.DealProposal{
			PayloadCID: payloadCID,
			ID:         id,
			Params:     params,
		},
		TotalFunds:       totalFunds,
		ClientWallet:     clientWallet,
		MinerWallet:      minerWallet,
		TotalReceived:    0,
		CurrentInterval:  params.PaymentInterval,
		BytesPaidFor:     0,
		PaymentRequested: abi.NewTokenAmount(0),
		FundsSpent:       abi.NewTokenAmount(0),
		Status:           retrievalmarket.DealStatusNew,
		Sender:           p.ID,
		UnsealFundsPaid:  big.Zero(),
	}
	// start the deal processing
	err := c.receiver.Begin(dealState.ID, &dealState)
	if err != nil {
		return 0, err
	}
	// sends the proposal to the remote peer
	// TODO? //legacy := deal.Status == retrievalmarket.DealStatusRetryLegacy
	legacy := false
	log.Infof("open pull from peerid:%s", p.ID)
	channelID, err := c.Environment.OpenDataTransfer(ctx, p.ID, &dealState.DealProposal, legacy)
	if err != nil {
		return 0, err
	}
	log.Infof("open channelID:%s for other peer:%s", channelID, p.ID)
	return id, nil
}

// Check if there's already an active retrieval deal with the same peer
// for the same payload CID
func (c *ClientEx) checkForActiveDeal(payloadCID cid.Cid, pid peer.ID) error {
	var deals []retrievalmarket.ClientDealState
	err := c.stateMachines.List(&deals)
	if err != nil {
		return err
	}

	for _, deal := range deals {
		match := deal.Sender == pid && deal.PayloadCID == payloadCID
		active := !clientstates.IsFinalityState(deal.Status)
		if match && active {
			msg := fmt.Sprintf("there is an active retrieval deal with peer %s ", pid)
			msg += fmt.Sprintf("for payload CID %s ", payloadCID)
			msg += fmt.Sprintf("(retrieval deal ID %d, state %s) - ",
				deal.ID, retrievalmarket.DealStatuses[deal.Status])
			msg += "existing deal must be cancelled before starting a new retrieval deal"
			err := xerrors.Errorf(msg)
			return err
		}
	}
	return nil
}

func (c *ClientEx) notifySubscribers(eventName fsm.EventName, state fsm.StateType) {
	evt := eventName.(retrievalmarket.ClientEvent)
	ds := state.(retrievalmarket.ClientDealState)
	_ = c.subscribers.Publish(internalEvent{evt, ds})
}

func (c *ClientEx) addMultiaddrs(ctx context.Context, p retrievalmarket.RetrievalPeer) error {
	tok, _, err := c.node.GetChainHead(ctx)
	if err != nil {
		return err
	}
	maddrs, err := c.node.GetKnownAddresses(ctx, p, tok)
	if err != nil {
		return err
	}
	if len(maddrs) > 0 {
		c.network.AddAddrs(p.ID, maddrs)
	}
	return nil
}

// SubscribeToEvents allows another component to listen for events on the RetrievalClient
// in order to track deals as they progress through the deal flow
func (c *ClientEx) SubscribeToEvents(subscriber retrievalmarket.ClientSubscriber) retrievalmarket.Unsubscribe {
	return retrievalmarket.Unsubscribe(c.subscribers.Subscribe(subscriber))
}

// V1

// TryRestartInsufficientFunds attempts to restart any deals stuck in the insufficient funds state
// after funds are added to a given payment channel
func (c *ClientEx) TryRestartInsufficientFunds(paymentChannel address.Address) error {
	var deals []retrievalmarket.ClientDealState
	err := c.stateMachines.List(&deals)
	if err != nil {
		return err
	}
	for _, deal := range deals {
		if deal.Status == retrievalmarket.DealStatusInsufficientFunds && deal.PaymentInfo.PayCh == paymentChannel {
			if err := c.stateMachines.Send(deal.ID, retrievalmarket.ClientEventRecheckFunds); err != nil {
				return err
			}
		}
	}
	return nil
}

// CancelDeal attempts to cancel an in progress deal
func (c *ClientEx) CancelDeal(dealID retrievalmarket.DealID) error {
	return c.stateMachines.Send(dealID, retrievalmarket.ClientEventCancel)
}

// GetDeal returns a given deal by deal ID, if it exists
func (c *ClientEx) GetDeal(dealID retrievalmarket.DealID) (retrievalmarket.ClientDealState, error) {
	var deal *retrievalmarket.ClientDealState
	a, ok := c.receiver.Get(dealID)
	if !ok {
		return retrievalmarket.ClientDealState{}, xerrors.Errorf("failed to get client deal state")
	}
	deal = a.(*retrievalmarket.ClientDealState)
	return *deal, nil
}

// ListDeals lists all known retrieval deals
func (c *ClientEx) ListDeals() (map[retrievalmarket.DealID]retrievalmarket.ClientDealState, error) {
	var deals []retrievalmarket.ClientDealState
	err := c.stateMachines.List(&deals)
	if err != nil {
		return nil, err
	}
	dealMap := make(map[retrievalmarket.DealID]retrievalmarket.ClientDealState)
	for _, deal := range deals {
		dealMap[deal.ID] = deal
	}
	return dealMap, nil
}

type clientDealEnvironment struct {
	c *ClientEx
}

// Node returns the node interface for this deal
func (c *clientDealEnvironment) Node() retrievalmarket.RetrievalClientNode {
	return c.c.node
}

func (c *clientDealEnvironment) OpenDataTransfer(ctx context.Context, to peer.ID, proposal *retrievalmarket.DealProposal, legacy bool) (datatransfer.ChannelID, error) {
	sel := selectorparse.CommonSelector_ExploreAllRecursively
	if proposal.SelectorSpecified() {
		var err error
		sel, err = retrievalmarket.DecodeNode(proposal.Selector)
		if err != nil {
			return datatransfer.ChannelID{}, xerrors.Errorf("selector is invalid: %w", err)
		}
	}

	var vouch datatransfer.Voucher = proposal
	if legacy {
		vouch = &migrations.DealProposal0{
			PayloadCID: proposal.PayloadCID,
			ID:         proposal.ID,
			Params0: migrations.Params0{
				Selector:                proposal.Selector,
				PieceCID:                proposal.PieceCID,
				PricePerByte:            proposal.PricePerByte,
				PaymentInterval:         proposal.PaymentInterval,
				PaymentIntervalIncrease: proposal.PaymentIntervalIncrease,
				UnsealPrice:             proposal.UnsealPrice,
			},
		}
	}

	blocksPerHour := 60 * 60 / build.BlockDelaySecs
	dealStartEpoch := abi.ChainEpoch(dealStartBufferHours * blocksPerHour)
	num, _ := rand.Int(rand.Reader, math.NewInt(9999999999))
	dealEndEpoch := dealStartEpoch + abi.ChainEpoch(num.Int64())

	data, err := cid.Parse(proposal.PayloadCID)
	if err != nil {
		return datatransfer.ChannelID{}, xerrors.Errorf("selector is invalid: %w", err)
	}

	refdata := &storagemarket.DataRef{
		TransferType: storagemarket.TTGraphsync,
		Root:         data,
	}

	addr, _ := address.NewFromString("t01000")

	mockfullapi := c.c.mockFullApi
	mi, err := mockfullapi.StateMinerInfo(ctx, addr, types.EmptyTSK)
	if err != nil {
		return datatransfer.ChannelID{}, xerrors.Errorf("get sector size & type error: %w", err)
	}
	ver, err := mockfullapi.StateNetworkVersion(ctx, types.EmptyTSK)
	if err != nil {
		return datatransfer.ChannelID{}, xerrors.Errorf("get network version error: %w", err)
	}

	sp, err := miner.PreferredSealProofTypeFromWindowPoStType(ver, mi.WindowPoStProofType)
	if err != nil {
		return datatransfer.ChannelID{}, xerrors.Errorf("get PreferredSealProofTypeFromWindowPoStType error: %w", err)
	}

	deal := storagemarket.ProposeStorageDealParams{
		Addr: addr,
		Info: &storagemarket.StorageProviderInfo{
			PeerID:     to,
			SectorSize: uint64(mi.SectorSize),
			Address:    addr,
		},
		Data:          refdata,
		StartEpoch:    dealStartEpoch,
		EndEpoch:      dealEndEpoch,
		Price:         types.EmptyInt,
		Collateral:    bigext.Zero(),
		Rt:            sp,
		FastRetrieval: true,
		VerifiedDeal:  false,
	}

	res := c.c.GetRetrievalType(proposal.PayloadCID)
	if res.Headflag == true && res.Auth == true { //res.Headflag == true
		//receive provider data
		newhead, err := storagenetwork.CreateNewHead(ctx, proposal, deal, res.Txhash)
		if err == nil {

			if err != nil {
				return datatransfer.ChannelID{}, xerrors.Errorf("get new head %w", err)
			}

			if len(newhead) == 0 {
				return datatransfer.ChannelID{}, xerrors.Errorf("new head length 0")
			}

			savepath := c.c.GetRetrievalPath()

			newheadpath := filepath.Join(savepath, fmt.Sprintf("%s", res.Txhash))
			file, err := os.Create(newheadpath)
			if err != nil {
				return datatransfer.ChannelID{}, xerrors.Errorf("failed to create car file for import: %w", err)
			}

			bytehead, err := base64.StdEncoding.DecodeString(newhead)
			if err != nil {
				file.Close()
				return datatransfer.ChannelID{}, xerrors.Errorf("failed to base64 decode: %w", err)
			}

			file.Write(bytehead)

			// close the file before returning the path.
			if err := file.Close(); err != nil {
				return datatransfer.ChannelID{}, xerrors.Errorf("failed to close  file: %w", err)
			}

			c.c.MarkeAutStatus(res.Txhash, DEAL_SUCCESS)
		}

	}

	return c.c.dataTransfer.OpenPullDataChannel(ctx, to, vouch, proposal.PayloadCID, sel)
}

func (c *clientDealEnvironment) SendDataTransferVoucher(ctx context.Context, channelID datatransfer.ChannelID, payment *retrievalmarket.DealPayment, legacy bool) error {
	var vouch datatransfer.Voucher = payment
	if legacy {
		vouch = &migrations.DealPayment0{
			ID:             payment.ID,
			PaymentChannel: payment.PaymentChannel,
			PaymentVoucher: payment.PaymentVoucher,
		}
	}
	return c.c.dataTransfer.SendVoucher(ctx, channelID, vouch)
}

func (c *clientDealEnvironment) CloseDataTransfer(ctx context.Context, channelID datatransfer.ChannelID) error {
	// When we close the data transfer, we also send a cancel message to the peer.
	// Make sure we don't wait too long to send the message.
	ctx, cancel := context.WithTimeout(ctx, shared.CloseDataTransferTimeout)
	defer cancel()

	err := c.c.dataTransfer.CloseDataTransferChannel(ctx, channelID)
	if shared.IsCtxDone(err) {
		log.Warnf("failed to send cancel data transfer channel %s to provider within timeout %s",
			channelID, shared.CloseDataTransferTimeout)
		return nil
	}
	return err
}

// FinalizeBlockstore is called when all blocks have been received
func (c *clientDealEnvironment) FinalizeBlockstore(dealID retrievalmarket.DealID) error {
	return c.c.bstores.Done(dealID)
}

type clientStoreGetter struct {
	c *ClientEx
}

func (csg *clientStoreGetter) Get(_ peer.ID, id retrievalmarket.DealID) (bstore.Blockstore, error) {
	var deal *retrievalmarket.ClientDealState
	a, ok := csg.c.receiver.Get(id)
	if !ok {
		return nil, xerrors.Errorf("failed to get client deal state")
	}
	deal = a.(*retrievalmarket.ClientDealState)
	payloadCID := deal.DealProposal.PayloadCID
	return csg.c.bstores.Get(id, payloadCID)
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

// DataTransfer Subscriber Func
func ClientDataTransferSubscriberExt(c *ClientEx) datatransfer.Subscriber {
	return func(event datatransfer.Event, channelState datatransfer.ChannelState) {

		dealProposal, ok := dealProposalFromVoucher(channelState.Voucher())
		// if this event is for a transfer not related to storage, ignore
		if !ok {
			return
		}

		// get clientDealState
		clientDealState, err := c.GetDeal(dealProposal.ID)
		if err != nil {
			log.Errorf("get clientDealState error:%s", err)
			c.notifySubscribers(retrievalmarket.ClientEventDealNotFound, clientDealState)
			return
		}

		log.Infow("processing retrieval client dt event", "event", datatransfer.Events[event.Code], "peer", channelState.OtherPeer(), "channelState", channelState)

		if event.Code == datatransfer.DataReceivedProgress {
			// log.Infof("### data Received:%d/%d  %d,%d",channelState.Received(),channelState.TotalSize(),channelState.Sent(),channelState.ReceivedCidsTotal())
			clientDealState.TotalReceived = channelState.Received()
			clientDealState.Status = retrievalmarket.DealStatusOngoing
			c.notifySubscribers(retrievalmarket.ClientEventBlocksReceived, clientDealState)
		} else if event.Code == datatransfer.FinishTransfer {
			log.Infof("dealProposal.ID:%s", dealProposal.ID)
			err = c.Environment.FinalizeBlockstore(dealProposal.ID)
			if err != nil {
				log.Errorf("FinalizeBlockstore error:%s", err)
				c.notifySubscribers(retrievalmarket.ClientEventFinalizeBlockstoreErrored, clientDealState)
				return
			}
			clientDealState.Status = retrievalmarket.DealStatusCompleted
			c.notifySubscribers(retrievalmarket.ClientEventComplete, clientDealState)
		} else if event.Code == datatransfer.Cancel {
			// log.Infof("### data Received:%d/%d  %d,%d",channelState.Received(),channelState.TotalSize(),channelState.Sent(),channelState.ReceivedCidsTotal())
			clientDealState.TotalReceived = channelState.Received()
			clientDealState.Status = retrievalmarket.DealStatusErrored
			c.notifySubscribers(retrievalmarket.ClientEventComplete, clientDealState)
		} else if event.Code == datatransfer.Error {
			// if error,send errorEvent
			clientDealState.TotalReceived = channelState.Received()
			clientDealState.Status = retrievalmarket.DealStatusErrored
			c.notifySubscribers(retrievalmarket.ClientEventDataTransferError, clientDealState)
		}
	}
}

func (c *ClientEx) GetRetrievalType(cid cid.Cid) RetrieveType {
	key := datastore.NewKey(cid.String())
	valbuf, err := c.Ds.Get(key)
	if err != nil {
		return RetrieveType{"", false, false, ""}
	}
	var rs RetrieveType
	err = json.Unmarshal(valbuf, &rs)
	if err != nil {
		return RetrieveType{"", false, false, ""}
	}
	c.Ds.Delete(key)
	return rs
}

func NewAuthStatus(txhash string, status int) AuthStatus {
	return AuthStatus{txhash, status}
}

func NewRetrieveType(cid string, headflag bool, auth bool, txhash string) RetrieveType {
	return RetrieveType{cid, headflag, auth, txhash}
}

func (c *ClientEx) GetAuthStatus(filekey string) int {
	key := datastore.NewKey(filekey)
	valbuf, err := c.Ds.Get(key)
	if err != nil {
		return -1
	}
	var rs AuthStatus
	err = json.Unmarshal(valbuf, &rs)
	if err != nil {
		return -1
	}
	return rs.Status
}

func (c *ClientEx) MarkeAutStatus(txhash string, status int) bool {
	rs := NewAuthStatus(txhash, status)
	key := datastore.NewKey(txhash)
	data, _ := json.Marshal(rs)
	err2 := c.Ds.Put(key, data)
	if err2 != nil {
		log.Error("MarkeAutStatus", err2)
		return false
	}
	return true
}

func (c *ClientEx) MarkRetrievePath(retrievepath string) bool {
	rs := GlobalData{retrievepath}
	key := datastore.NewKey("retrievepath")
	data, _ := json.Marshal(rs)
	err2 := c.Ds.Put(key, data)
	if err2 != nil {
		return false
	}
	return true
}

func (c *ClientEx) GetRetrievalPath() string {
	key := datastore.NewKey("retrievepath")
	valbuf, err := c.Ds.Get(key)
	if err != nil {
		return ""
	}
	var rs GlobalData
	err = json.Unmarshal(valbuf, &rs)
	if err != nil {
		return ""
	}
	return rs.Savepath
}

func (c *ClientEx) SetStroageNetwork(stroage interface{}) {
	storagenetwork = stroage.(*storageclientext.ClientEx)
}

func (c *ClientEx) MarkRetrieveType(cid cid.Cid, headflag bool, auth bool, txhash string) bool {
	rs := NewRetrieveType(cid.String(), headflag, auth, txhash)
	key := datastore.NewKey(cid.String())
	data, _ := json.Marshal(rs)
	err2 := c.Ds.Put(key, data)
	if err2 != nil {
		log.Errorf("MarkRetrieveType error:%w", err2)
		return false
	}
	return true
}
