package modulesext

import (
	"context"
	"github.com/ethereum/go-ethereum/storage/clientext"
	retrievalimpl2 "github.com/ethereum/go-ethereum/storage/clientext/retrieval"
	"github.com/ethereum/go-ethereum/storage/mock"
	"github.com/filecoin-project/go-data-transfer/channelmonitor"
	dtimpl "github.com/filecoin-project/go-data-transfer/impl"
	dtnet "github.com/filecoin-project/go-data-transfer/network"
	dtgstransport "github.com/filecoin-project/go-data-transfer/transport/graphsync"
	"github.com/filecoin-project/go-fil-markets/discovery"
	discoveryimpl "github.com/filecoin-project/go-fil-markets/discovery/impl"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	rmnet "github.com/filecoin-project/go-fil-markets/retrievalmarket/network"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	smnet "github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/lotus/journal"
	"github.com/filecoin-project/lotus/markets"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/filecoin-project/lotus/markets/retrievaladapter"
	"github.com/filecoin-project/lotus/node/config"
	"github.com/filecoin-project/lotus/node/impl/full"
	payapi "github.com/filecoin-project/lotus/node/impl/paych"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/node/repo/imports"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/ipfs/go-graphsync"
	graphsyncimpl "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"time"
)

// NewClientGraphsyncDataTransfer returns a data transfer manager that just
// uses the clients's Client DAG service for transfers
func NewClientGraphsyncDataTransfer(ctx context.Context, h host.Host, gs dtypes.Graphsync, ds dtypes.MetadataDS, r repo.LockedRepo) (dtypes.ClientDataTransfer, error) {
	// go-data-transfer protocol retries:
	// 1s, 5s, 25s, 2m5s, 5m x 11 ~= 1 hour
	dtRetryParams := dtnet.RetryParameters(time.Second, 5*time.Minute, 15, 5)
	net := dtnet.NewFromLibp2pHost(h, dtRetryParams)

	dtDs := namespace.Wrap(ds, datastore.NewKey("/datatransfer/client/transfers"))
	transport := dtgstransport.NewTransport(h.ID(), gs, net)
	err := os.MkdirAll(filepath.Join(r.Path(), "data-transfer"), 0755) //nolint: gosec
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	// data-transfer push / pull channel restart configuration:
	dtRestartConfig := dtimpl.ChannelRestartConfig(channelmonitor.Config{
		// Disable Accept and Complete timeouts until this issue is resolved:
		// https://github.com/filecoin-project/lotus/issues/6343#
		// Wait for the other side to respond to an Open channel message
		AcceptTimeout: 0,
		// Wait for the other side to send a Complete message once all
		// data has been sent / received
		CompleteTimeout: 0,

		// When an error occurs, wait a little while until all related errors
		// have fired before sending a restart message
		RestartDebounce: 10 * time.Second,
		// After sending a restart, wait for at least 1 minute before sending another
		RestartBackoff: time.Minute,
		// After trying to restart 3 times, give up and fail the transfer
		MaxConsecutiveRestarts: 3,
	})
	dt, err := dtimpl.NewDataTransfer(dtDs, filepath.Join(r.Path(), "data-transfer"), net, transport, dtRestartConfig)
	if err != nil {
		return nil, err
	}

	dt.OnReady(marketevents.ReadyLogger("client data transfer"))
	dt.SubscribeToEvents(marketevents.DataTransferLogger)
	dt.Start(ctx)
	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//		dt.SubscribeToEvents(marketevents.DataTransferLogger)
	//		return dt.Start(ctx)
	//	},
	//	OnStop: func(ctx context.Context) error {
	//		return dt.Stop(ctx)
	//	},
	//})
	return dt, nil
}

// Graphsync creates a graphsync instance from the given loader and storer
func Graphsync(ctx context.Context, parallelTransfersForStorage uint64, parallelTransfersForRetrieval uint64, r repo.LockedRepo, clientBs dtypes.ClientBlockstore, chainBs dtypes.ExposedBlockstore, h host.Host) (dtypes.Graphsync, error) {
	graphsyncNetwork := gsnet.NewFromLibp2pHost(h)
	lsys := storeutil.LinkSystemForBlockstore(clientBs)

	gs := graphsyncimpl.New(ctx,
		graphsyncNetwork,
		lsys,
		graphsyncimpl.RejectAllRequestsByDefault(),
		graphsyncimpl.MaxInProgressIncomingRequests(parallelTransfersForStorage),
		graphsyncimpl.MaxInProgressOutgoingRequests(parallelTransfersForRetrieval),
		graphsyncimpl.MaxLinksPerIncomingRequests(config.MaxTraversalLinks),
		graphsyncimpl.MaxLinksPerOutgoingRequests(config.MaxTraversalLinks))
	chainLinkSystem := storeutil.LinkSystemForBlockstore(chainBs)
	err := gs.RegisterPersistenceOption("chainstore", chainLinkSystem)
	if err != nil {
		return nil, err
	}
	gs.RegisterIncomingRequestHook(func(p peer.ID, requestData graphsync.RequestData, hookActions graphsync.IncomingRequestHookActions) {
		_, has := requestData.Extension("chainsync")
		if has {
			// TODO: we should confirm the selector is a reasonable one before we validate
			// TODO: this code will get more complicated and should probably not live here eventually
			hookActions.ValidateRequest()
			hookActions.UsePersistenceOption("chainstore")
		}
	})
	gs.RegisterOutgoingRequestHook(func(p peer.ID, requestData graphsync.RequestData, hookActions graphsync.OutgoingRequestHookActions) {
		_, has := requestData.Extension("chainsync")
		if has {
			hookActions.UsePersistenceOption("chainstore")
		}
	})

	//graphsyncStats(mctx, lc, gs)

	return gs, nil
}

func ClientImportMgr(ds dtypes.MetadataDS, r repo.LockedRepo) (dtypes.ClientImportMgr, error) {
	// store the imports under the repo's `imports` subdirectory.
	dir := filepath.Join(r.Path(), "imports")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, xerrors.Errorf("failed to create directory %s: %w", dir, err)
	}

	ns := namespace.Wrap(ds, datastore.NewKey("/client"))
	return imports.NewManager(ns, dir), nil
}

func StorageClient(h host.Host, dataTransfer dtypes.ClientDataTransfer, discovery *discoveryimpl.Local,
	deals dtypes.ClientDatastore, scn storagemarket.StorageClientNode, accessor storagemarket.BlockstoreAccessor, j journal.Journal) (storagemarket.StorageClient, error) {
	// go-fil-markets protocol retries:
	// 1s, 5s, 25s, 2m5s, 5m x 11 ~= 1 hour
	marketsRetryParams := smnet.RetryParameters(time.Second, 5*time.Minute, 15, 5)
	net := smnet.NewFromLibp2pHost(h, marketsRetryParams)

	c, err := clientext.NewClientEx(net, dataTransfer, discovery, deals, scn, accessor, clientext.DealPollingInterval(time.Second), clientext.MaxTraversalLinks(config.MaxTraversalLinks))
	if err != nil {
		return nil, err
	}
	c.OnReady(marketevents.ReadyLogger("storage client"))
	//c.SubscribeToEvents(marketevents.StorageClientLogger)
	//evtType := j.RegisterEventType("markets/storage/client", "state_change")
	//c.SubscribeToEvents(markets.StorageClientJournaler(j, evtType))
	//c.Start(ctx)
	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//		c.SubscribeToEvents(marketevents.StorageClientLogger)
	//
	//		evtType := j.RegisterEventType("markets/storage/client", "state_change")
	//		c.SubscribeToEvents(markets.StorageClientJournaler(j, evtType))
	//
	//		return c.Start(ctx)
	//	},
	//	OnStop: func(context.Context) error {
	//		return c.Stop()
	//	},
	//})
	return c, nil
}

// RetrievalClient creates a new retrieval client attached to the client blockstore
func RetrievalClient(h host.Host, dt dtypes.ClientDataTransfer, resolver discovery.PeerResolver,
	ds dtypes.MetadataDS, chainAPI full.ChainAPI, stateAPI full.StateAPI, fullApi mock.MockFullNode, accessor retrievalmarket.BlockstoreAccessor, j journal.Journal) (retrievalmarket.RetrievalClient, error) {
	payAPI := payapi.PaychAPI{}
	adapter := retrievaladapter.NewRetrievalClientNode(payAPI, chainAPI, stateAPI)
	network := rmnet.NewFromLibp2pHost(h)
	ds = namespace.Wrap(ds, datastore.NewKey("/retrievals/client"))
	client, err := retrievalimpl2.NewClientEx(network, dt, adapter, resolver, ds, accessor, fullApi)
	if err != nil {
		return nil, err
	}
	client.OnReady(marketevents.ReadyLogger("retrieval client"))

	client.SubscribeToEvents(marketevents.RetrievalClientLogger)
	evtType := j.RegisterEventType("markets/retrieval/client", "state_change")
	client.SubscribeToEvents(markets.RetrievalClientJournaler(j, evtType))
	return client, nil
}

func RetrievalLinkStorage(stroage interface{}, retrieval interface{}, retrievalpath string) {
	retrievalnetwork := retrieval.(*retrievalimpl2.ClientEx)

	retrievalnetwork.MarkRetrievePath(retrievalpath)

	retrievalnetwork.SetStroageNetwork(stroage)
}
