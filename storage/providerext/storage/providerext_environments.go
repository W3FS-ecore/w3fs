package storage

import (
	"context"
	"io"
	"os"

	"github.com/ipfs/go-cid"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-commp-utils/writer"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-fil-markets/filestore"
	"github.com/filecoin-project/go-fil-markets/piecestore"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/providerstates"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-fil-markets/stores"
)

// -------
// providerDealEnvironment
// -------

type providerExtDealEnvironment struct {
	p *ProviderExt
}

func (p *providerExtDealEnvironment) RegisterShard(ctx context.Context, pieceCid cid.Cid, carPath string, eagerInit bool) error {
	return stores.RegisterShardSync(ctx, p.p.dagStore, pieceCid, carPath, eagerInit)
}

func (p *providerExtDealEnvironment) ReadCAR(path string) (*carv2.Reader, error) {
	return carv2.OpenReader(path)
}

func (p *providerExtDealEnvironment) FinalizeBlockstore(proposalCid cid.Cid) error {
	bs, err := p.p.stores.Get(proposalCid.String())
	if err != nil {
		return xerrors.Errorf("failed to get read/write blockstore: %w", err)
	}

	if err := bs.Finalize(); err != nil {
		return xerrors.Errorf("failed to finalize read/write blockstore: %w", err)
	}

	return nil
}

func (p *providerExtDealEnvironment) TerminateBlockstore(proposalCid cid.Cid, path string) error {
	// stop tracking it.
	if err := p.p.stores.Untrack(proposalCid.String()); err != nil {
		log.Warnf("failed to untrack read write blockstore, proposalCid=%s, car_path=%s: %s", proposalCid, path, err)
	}

	// delete the backing CARv2 file as it was a temporary file we created for
	// this storage deal; the piece has now been handed off, or the deal has failed.
	if err := os.Remove(path); err != nil {
		log.Warnf("failed to delete carv2 file on termination, car_path=%s: %s", path, err)
	}

	return nil
}

func (p *providerExtDealEnvironment) Address() address.Address {
	return p.p.actor
}

func (p *providerExtDealEnvironment) Node() storagemarket.StorageProviderNode {
	return p.p.spn
}

func (p *providerExtDealEnvironment) Ask() storagemarket.StorageAsk {
	sask := p.p.storedAsk.GetAsk()
	if sask == nil {
		return storagemarket.StorageAskUndefined
	}
	return *sask.Ask
}

// GeneratePieceCommitment generates the pieceCid for the CARv1 deal payload in
// the CARv2 file that already exists at the given path.
func (p *providerExtDealEnvironment) GeneratePieceCommitment(proposalCid cid.Cid, carPath string, dealSize abi.PaddedPieceSize) (c cid.Cid, path filestore.Path, finalErr error) {
	rd, err := carv2.OpenReader(carPath)
	if err != nil {
		return cid.Undef, "", xerrors.Errorf("failed to get CARv2 reader, proposalCid=%s, carPath=%s: %w", proposalCid, carPath, err)
	}

	defer func() {
		if err := rd.Close(); err != nil {
			log.Errorf("failed to close CARv2 reader, carPath=%s, err=%s", carPath, err)

			if finalErr == nil {
				c = cid.Undef
				path = ""
				finalErr = xerrors.Errorf("failed to close CARv2 reader, proposalCid=%s, carPath=%s: %w",
					proposalCid, carPath, err)
				return
			}
		}
	}()

	// dump the CARv1 payload of the CARv2 file to the Commp Writer and get back the CommP.
	w := &writer.Writer{}
	written, err := io.Copy(w, rd.DataReader())
	if err != nil {
		return cid.Undef, "", xerrors.Errorf("failed to write to CommP writer: %w", err)
	}
	if written != int64(rd.Header.DataSize) {
		return cid.Undef, "", xerrors.Errorf("number of bytes written to CommP writer %d not equal to the CARv1 payload size %d", written, rd.Header.DataSize)
	}

	cidAndSize, err := w.Sum()
	if err != nil {
		return cid.Undef, "", xerrors.Errorf("failed to get CommP: %w", err)
	}

	if cidAndSize.PieceSize < dealSize {
		// need to pad up!
		rawPaddedCommp, err := commp.PadCommP(
			// we know how long a pieceCid "hash" is, just blindly extract the trailing 32 bytes
			cidAndSize.PieceCID.Hash()[len(cidAndSize.PieceCID.Hash())-32:],
			uint64(cidAndSize.PieceSize),
			uint64(dealSize),
		)
		if err != nil {
			return cid.Undef, "", err
		}
		cidAndSize.PieceCID, _ = commcid.DataCommitmentV1ToCID(rawPaddedCommp)
	}

	return cidAndSize.PieceCID, filestore.Path(""), err
}

func (p *providerExtDealEnvironment) FileStore() filestore.FileStore {
	return p.p.fs
}

func (p *providerExtDealEnvironment) PieceStore() piecestore.PieceStore {
	return p.p.pieceStore
}

func (p *providerExtDealEnvironment) SendSignedResponse(ctx context.Context, resp *network.Response) error {
	s, err := p.p.conns.DealStream(resp.Proposal)
	if err != nil {
		return xerrors.Errorf("couldn't send response: %w", err)
	}

	sig, err := p.p.sign(ctx, resp)
	if err != nil {
		return xerrors.Errorf("failed to sign response message: %w", err)
	}

	signedResponse := network.SignedResponse{
		Response:  *resp,
		Signature: sig,
	}

	err = s.WriteDealResponse(signedResponse, p.p.sign)
	if err != nil {
		// Assume client disconnected
		_ = p.p.conns.Disconnect(resp.Proposal)
	}
	return err
}

func (p *providerExtDealEnvironment) Disconnect(proposalCid cid.Cid) error {
	return p.p.conns.Disconnect(proposalCid)
}

func (p *providerExtDealEnvironment) RunCustomDecisionLogic(ctx context.Context, deal storagemarket.MinerDeal) (bool, string, error) {
	if p.p.customDealDeciderFunc == nil {
		return true, "", nil
	}
	return p.p.customDealDeciderFunc(ctx, deal)
}

func (p *providerExtDealEnvironment) TagPeer(id peer.ID, s string) {
	p.p.net.TagPeer(id, s)
}

func (p *providerExtDealEnvironment) UntagPeer(id peer.ID, s string) {
	p.p.net.UntagPeer(id, s)
}

var _ providerstates.ProviderDealEnvironment = &providerExtDealEnvironment{}

type providerStoreGetter struct {
	p *ProviderExt
}

func (psg *providerStoreGetter) Get(proposalCid cid.Cid) (bstore.Blockstore, error) {
	// Wait for the providerext to be ready
	err := awaitProviderReady(psg.p)
	if err != nil {
		return nil, err
	}

	var deal *storagemarket.MinerDeal
	d, ok := psg.p.receiver.Get(proposalCid)
	if !ok {
		return nil, xerrors.Errorf("failed to get deal state: %w", err)
	}
	deal = d.(*storagemarket.MinerDeal)

	return psg.p.stores.GetOrOpen(proposalCid.String(), deal.InboundCAR, deal.Ref.Root)
}

type providerPushDeals struct {
	p *ProviderExt
}

func (ppd *providerPushDeals) Get(proposalCid cid.Cid) (storagemarket.MinerDeal, error) {
	// Wait for the providerext to be ready
	var deal storagemarket.MinerDeal
	err := awaitProviderReady(ppd.p)
	if err != nil {
		return deal, err
	}

	err = ppd.p.receiver.GetSync(context.TODO(), proposalCid, &deal)
	return deal, err
}

// awaitProviderReady waits for the providerext to startup
func awaitProviderReady(p *ProviderExt) error {
	err := p.AwaitReady()
	if err != nil {
		return xerrors.Errorf("could not get deal with proposal CID %s: error waiting for providerext startup: %w")
	}

	return nil
}
