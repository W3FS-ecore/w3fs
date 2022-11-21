package retrievalimpl

import (
	"context"
	"github.com/ipfs/go-cid"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/dagstore"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"

	"github.com/filecoin-project/go-fil-markets/piecestore"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket/impl/dtutils"
	"github.com/filecoin-project/go-fil-markets/shared"
)

var providerValidationEnv = new(providerExtValidationEnvironment)

type providerExtValidationEnvironment struct {
	p *ProviderExt
}

func (pve *providerExtValidationEnvironment) GetAsk(ctx context.Context, payloadCid cid.Cid, pieceCid *cid.Cid,
	piece piecestore.PieceInfo, isUnsealed bool, client peer.ID) (retrievalmarket.Ask, error) {
	// no need to process it. now do nothing.return nil.
	ask := retrievalmarket.Ask{}
	return ask, nil
}

func (pve *providerExtValidationEnvironment) GetPiece(c cid.Cid, pieceCID *cid.Cid) (piecestore.PieceInfo, bool, error) {
	inPieceCid := cid.Undef
	if pieceCID != nil {
		inPieceCid = *pieceCID
	}

	return pve.p.getPieceInfoFromCid(context.TODO(), c, inPieceCid)
}

// CheckDealParams verifies the given deal params are acceptable
func (pve *providerExtValidationEnvironment) CheckDealParams(ask retrievalmarket.Ask, pricePerByte abi.TokenAmount, paymentInterval uint64, paymentIntervalIncrease uint64, unsealPrice abi.TokenAmount) error {
	// no need to check it.
	return nil
}

// RunDealDecisioningLogic runs custom deal decision logic to decide if a deal is accepted, if present
func (pve *providerExtValidationEnvironment) RunDealDecisioningLogic(ctx context.Context, state retrievalmarket.ProviderDealState) (bool, string, error) {
	if pve.p.dealDecider == nil {
		return true, "", nil
	}
	return pve.p.dealDecider(ctx, state)
}

// StateMachines returns the FSM Group to begin tracking with
func (pve *providerExtValidationEnvironment) BeginTracking(pds retrievalmarket.ProviderDealState) error {
	log.Infof("pds.Identifier():%s", pds.Identifier())
	err := pve.p.receiver.Begin(pds.Identifier(), &pds)
	if err != nil {
		return err
	}

	if pds.UnsealPrice.GreaterThan(big.Zero()) {
		//return pve.p.receiver.Send(pds.Identifier(), retrievalmarket.ProviderEventPaymentRequested, uint64(0))
	}

	//return pve.p.receiver.Send(pds.Identifier(), retrievalmarket.ProviderEventOpen)
	return nil
}

var providerRevalidationEnv = new(providerRevalidatorEnvironment)

type providerRevalidatorEnvironment struct {
	p *ProviderExt
}

func (pre *providerRevalidatorEnvironment) Node() retrievalmarket.RetrievalProviderNode {
	return pre.p.node
}

func (pre *providerRevalidatorEnvironment) SendEvent(dealID retrievalmarket.ProviderDealIdentifier, evt retrievalmarket.ProviderEvent, args ...interface{}) error {
	//return pre.p.receiver.Send(dealID, evt, args...)
	return nil
}

func (pre *providerRevalidatorEnvironment) Get(dealID retrievalmarket.ProviderDealIdentifier) (retrievalmarket.ProviderDealState, error) {
	var deal *retrievalmarket.ProviderDealState
	a, ok := pre.p.receiver.Get(dealID)
	if ok {
		deal = a.(*retrievalmarket.ProviderDealState)
	} else {
		log.Errorf("failed to get dealState")
		return *deal, xerrors.Errorf("failed to get dealState")
	}
	log.Infof("deal:%s", *deal)
	return *deal, nil
}

var providerDealEnv = new(providerDealEnvironment)

type providerDealEnvironment struct {
	p *ProviderExt
}

// Node returns the node interface for this deal
func (pde *providerDealEnvironment) Node() retrievalmarket.RetrievalProviderNode {
	return pde.p.node
}

// UnsealData fetches the piece containing data needed for the retrieval,
// unsealing it if necessary
func (pde *providerDealEnvironment) UnsealData(ctx context.Context, deal retrievalmarket.ProviderDealState) error {
	if err := pde.PrepareBlockstore(ctx, deal.ID, deal.PieceInfo.PieceCID); err != nil {
		return err
	}
	log.Infof("blockstore prepared successfully, firing unseal complete for deal %d", deal.ID)
	return nil
}

// PrepareBlockstore is called when the deal data has been unsealed and we need
// to add all blocks to a blockstore that is used to serve retrieval
func (pde *providerDealEnvironment) PrepareBlockstore(ctx context.Context, dealID retrievalmarket.DealID, pieceCid cid.Cid) error {
	// Load the blockstore that has the deal data
	bs, err := pde.p.dagStore.LoadShard(ctx, pieceCid)
	if err != nil {
		log.Errorf("failed to load blockstore for piece %s: %w", pieceCid, err)
		return xerrors.Errorf("failed to load blockstore for piece %s: %w", pieceCid, err)
	}

	log.Debugf("adding blockstore for deal %d to tracker", dealID)
	_, err = pde.p.stores.Track(dealID.String(), bs)
	log.Debugf("added blockstore for deal %d to tracker", dealID)
	return err
}

// UnpauseDeal resumes a deal so we can start sending data after its unsealed
func (pde *providerDealEnvironment) UnpauseDeal(ctx context.Context, deal retrievalmarket.ProviderDealState) error {
	log.Debugf("unpausing data transfer for deal %d", deal.ID)
	err := pde.TrackTransfer(deal)
	if err != nil {
		return err
	}
	if deal.ChannelID != nil {
		log.Debugf("resuming data transfer for deal %d", deal.ID)
		err = pde.ResumeDataTransfer(ctx, *deal.ChannelID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pde *providerDealEnvironment) TrackTransfer(deal retrievalmarket.ProviderDealState) error {
	pde.p.revalidator.TrackChannel(deal)
	return nil
}

func (pde *providerDealEnvironment) UntrackTransfer(deal retrievalmarket.ProviderDealState) error {
	pde.p.revalidator.UntrackChannel(deal)
	return nil
}

func (pde *providerDealEnvironment) ResumeDataTransfer(ctx context.Context, chid datatransfer.ChannelID) error {
	return pde.p.dataTransfer.ResumeDataTransferChannel(ctx, chid)
}

func (pde *providerDealEnvironment) CloseDataTransfer(ctx context.Context, chid datatransfer.ChannelID) error {
	// When we close the data transfer, we also send a cancel message to the peer.
	// Make sure we don't wait too long to send the message.
	ctx, cancel := context.WithTimeout(ctx, shared.CloseDataTransferTimeout)
	defer cancel()

	err := pde.p.dataTransfer.CloseDataTransferChannel(ctx, chid)
	if shared.IsCtxDone(err) {
		log.Warnf("failed to send cancel data transfer channel %s to client within timeout %s",
			chid, shared.CloseDataTransferTimeout)
		return nil
	}
	return err
}

func (pde *providerDealEnvironment) DeleteStore(dealID retrievalmarket.DealID) error {
	// close the read-only blockstore and stop tracking it for the deal
	if err := pde.p.stores.Untrack(dealID.String()); err != nil {
		return xerrors.Errorf("failed to clean read-only blockstore for deal %d: %w", dealID, err)
	}

	return nil
}

func storageDealsForPiece(clientSpecificPiece bool, payloadCID cid.Cid, pieceInfo piecestore.PieceInfo, pieceStore piecestore.PieceStore) ([]abi.DealID, error) {
	var storageDeals []abi.DealID
	var err error
	if clientSpecificPiece {
		//  If the user wants to retrieve the payload from a specific piece,
		//  we only need to inspect storage deals made for that piece to quote a price.
		for _, d := range pieceInfo.Deals {
			storageDeals = append(storageDeals, d.DealID)
		}
	} else {
		// If the user does NOT want to retrieve from a specific piece, we'll have to inspect all storage deals
		// made for that piece to quote a price.
		storageDeals, err = getAllDealsContainingPayload(pieceStore, payloadCID)
		if err != nil {
			return nil, xerrors.Errorf("failed to fetch deals for payload: %w", err)
		}
	}

	if len(storageDeals) == 0 {
		return nil, xerrors.New("no storage deals found")
	}

	return storageDeals, nil
}

func getAllDealsContainingPayload(pieceStore piecestore.PieceStore, payloadCID cid.Cid) ([]abi.DealID, error) {
	cidInfo, err := pieceStore.GetCIDInfo(payloadCID)
	if err != nil {
		log.Errorf("get cid info: %w", err)
		return nil, xerrors.Errorf("get cid info: %w", err)
	}
	var dealsIds []abi.DealID
	var lastErr error

	for _, pieceBlockLocation := range cidInfo.PieceBlockLocations {
		pieceInfo, err := pieceStore.GetPieceInfo(pieceBlockLocation.PieceCID)
		if err != nil {
			lastErr = err
			continue
		}
		for _, d := range pieceInfo.Deals {
			dealsIds = append(dealsIds, d.DealID)
		}
	}

	if lastErr == nil && len(dealsIds) == 0 {
		return nil, xerrors.New("no deals found")
	}

	if lastErr != nil && len(dealsIds) == 0 {
		return nil, xerrors.Errorf("failed to fetch deals containing payload %s: %w", payloadCID, lastErr)
	}

	return dealsIds, nil
}

var _ dtutils.StoreGetter = &providerStoreGetter{}

type providerStoreGetter struct {
	p *ProviderExt
}

func (psg *providerStoreGetter) Get(otherPeer peer.ID, dealID retrievalmarket.DealID) (bstore.Blockstore, error) {
	// var deal retrievalmarket.ProviderDealState
	provDealID := retrievalmarket.ProviderDealIdentifier{Receiver: otherPeer, DealID: dealID}
	_, ok := psg.p.receiver.Get(provDealID)
	if !ok {
		log.Errorf("failed to get deal state: %w", ok)
		return nil, xerrors.Errorf("failed to get deal state: %w", ok)
	}

	//
	// When a request for data is received
	// 1. The data transfer layer calls Get to get the blockstore
	// 2. The data for the deal is unsealed
	// 3. The unsealed data is put into the blockstore (in this case a CAR file)
	// 4. The data is served from the blockstore (using blockstore.Get)
	//
	// So we use a "lazy" blockstore that can be returned in step 1
	// but is only accessed in step 4 after the data has been unsealed.
	//
	return newLazyBlockstore(func() (dagstore.ReadBlockstore, error) {
		return psg.p.stores.Get(dealID.String())
	}), nil
}
