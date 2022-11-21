package clientext

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipld/go-ipld-prime"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/xerrors"

	datatransfer "github.com/filecoin-project/go-data-transfer"

	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
)

// -------
// clientDealEnvironment
// -------

type ClientDealEnvironment struct {
	c *ClientEx
}

func (c *ClientDealEnvironment) NewDealStream(ctx context.Context, p peer.ID) (network.StorageDealStream, error) {
	return c.c.net.NewDealStream(ctx, p)
}

func (c *ClientDealEnvironment) Node() storagemarket.StorageClientNode {
	return c.c.node
}

func (c *ClientDealEnvironment) CleanBlockstore(payloadCid cid.Cid) error {
	return c.c.bstores.Done(payloadCid)
}

func (c *ClientDealEnvironment) StartDataTransfer(ctx context.Context, to peer.ID, voucher datatransfer.Voucher, baseCid cid.Cid, selector ipld.Node) (datatransfer.ChannelID,
	error) {
	chid, err := c.c.dataTransfer.OpenPushDataChannel(ctx, to, voucher, baseCid, selector)
	return chid, err
}

func (c *ClientDealEnvironment) RestartDataTransfer(ctx context.Context, channelId datatransfer.ChannelID) error {
	return c.c.dataTransfer.RestartDataTransferChannel(ctx, channelId)
}

func (c *ClientDealEnvironment) GetProviderDealState(ctx context.Context, proposalCid cid.Cid) (*storagemarket.ProviderDealState, error) {
	return c.c.GetProviderDealState(ctx, proposalCid)
}

func (c *ClientDealEnvironment) PollingInterval() time.Duration {
	return c.c.pollingInterval
}

type clientStoreGetter struct {
	c *ClientEx
}

func (csg *clientStoreGetter) Get(proposalCid cid.Cid) (bstore.Blockstore, error) {
	var deal *storagemarket.ClientDeal
	a, ok := csg.c.receiver.Get(proposalCid)
	if !ok {
		return nil, xerrors.Errorf("failed to get client deal state")
	}
	deal = a.(*storagemarket.ClientDeal)
	bs, err := csg.c.bstores.Get(deal.DataRef.Root)
	if err != nil {
		return nil, xerrors.Errorf("failed to get blockstore for %s: %w", proposalCid, err)
	}

	return bs, nil
}

func (c *ClientDealEnvironment) TagPeer(peer peer.ID, tag string) {
	c.c.net.TagPeer(peer, tag)
}

func (c *ClientDealEnvironment) UntagPeer(peer peer.ID, tag string) {
	c.c.net.UntagPeer(peer, tag)
}

type clientPullDeals struct {
	c *ClientEx
}

func (cpd *clientPullDeals) Get(proposalCid cid.Cid) (storagemarket.ClientDeal, error) {
	var deal storagemarket.ClientDeal
	err := cpd.c.statemachines.GetSync(context.TODO(), proposalCid, &deal)
	return deal, err
}
