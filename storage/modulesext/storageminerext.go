package modulesext

import (
	"context"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	retrievalimpl "github.com/ethereum/go-ethereum/storage/providerext/retrieval"
	storageImpl "github.com/ethereum/go-ethereum/storage/providerext/storage"
	"github.com/ethereum/go-ethereum/storage/storageext"
	"github.com/filecoin-project/go-address"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	dtimpl "github.com/filecoin-project/go-data-transfer/impl"
	dtnet "github.com/filecoin-project/go-data-transfer/network"
	dtgstransport "github.com/filecoin-project/go-data-transfer/transport/graphsync"
	piecefilestore "github.com/filecoin-project/go-fil-markets/filestore"
	"github.com/filecoin-project/go-fil-markets/piecestore"
	piecestoreimpl "github.com/filecoin-project/go-fil-markets/piecestore/impl"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	rmnet "github.com/filecoin-project/go-fil-markets/retrievalmarket/network"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/storedask"
	smnet "github.com/filecoin-project/go-fil-markets/storagemarket/network"
	filstores "github.com/filecoin-project/go-fil-markets/stores"
	"github.com/filecoin-project/go-statestore"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/blockstore"
	sectorstorage "github.com/filecoin-project/lotus/extern/sector-storage"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"
	"github.com/filecoin-project/lotus/journal"
	"github.com/filecoin-project/lotus/markets/dagstore"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/filecoin-project/lotus/node/config"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/storage"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	graphsync "github.com/ipfs/go-graphsync/impl"
	graphsyncimpl "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/libp2p/go-libp2p-core/host"
	"golang.org/x/xerrors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	StorageCounterDSPrefix = "/storage/nextid"
	StagingAreaDirName     = "deal-staging"
)

func NewProviderDAGServiceDataTransfer(h host.Host, gs dtypes.StagingGraphsync, ds dtypes.MetadataDS, r repo.LockedRepo) (dtypes.ProviderDataTransfer, error) {
	net := dtnet.NewFromLibp2pHost(h)

	dtDs := namespace.Wrap(ds, datastore.NewKey("/datatransfer/providerext/transfers"))
	transport := dtgstransport.NewTransport(h.ID(), gs, net)
	err := os.MkdirAll(filepath.Join(r.Path(), "data-transfer"), 0755) //nolint: gosec
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	dt, err := dtimpl.NewDataTransfer(dtDs, filepath.Join(r.Path(), "data-transfer"), net, transport)
	if err != nil {
		return nil, err
	}

	dt.OnReady(marketevents.ReadyLogger("providerext data transfer"))

	//dt.OnReady(marketevents.ReadyLogger("providerext data transfer"))
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

func StagingGraphsyncNotBS(mctx context.Context, h host.Host) dtypes.StagingGraphsync {
	graphsyncNetwork := gsnet.NewFromLibp2pHost(h)
	lsys := cidlink.DefaultLinkSystem()
	lsys.StorageReadOpener = func(ipld.LinkContext, ipld.Link) (io.Reader, error) {
		return nil, nil
	}
	lsys.StorageWriteOpener = func(ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
		return nil, nil, nil
	}
	gs := graphsync.New(mctx, graphsyncNetwork, lsys,
		graphsync.RejectAllRequestsByDefault(),
		graphsync.MaxInProgressIncomingRequests(20),
		graphsync.MaxInProgressOutgoingRequests(20),
		graphsyncimpl.MaxLinksPerIncomingRequests(config.MaxTraversalLinks),
		graphsyncimpl.MaxLinksPerOutgoingRequests(config.MaxTraversalLinks))

	return gs
}

// StagingGraphsync creates a graphsync instance which reads and writes blocks
// to the StagingBlockstore
func StagingGraphsync(ctx context.Context, parallelTransfersForStorage uint64, parallelTransfersForRetrieval uint64, ibs dtypes.StagingBlockstore, h host.Host) dtypes.StagingGraphsync {
	graphsyncNetwork := gsnet.NewFromLibp2pHost(h)
	lsys := storeutil.LinkSystemForBlockstore(ibs)
	gs := graphsync.New(ctx,
		graphsyncNetwork,
		lsys,
		graphsync.RejectAllRequestsByDefault(),
		graphsync.MaxInProgressIncomingRequests(parallelTransfersForRetrieval),
		graphsync.MaxInProgressOutgoingRequests(parallelTransfersForStorage),
		graphsyncimpl.MaxLinksPerIncomingRequests(config.MaxTraversalLinks),
		graphsyncimpl.MaxLinksPerOutgoingRequests(config.MaxTraversalLinks))

	//graphsyncStats(mctx, lc, gs)

	return gs
}

// StagingBlockstore creates a blockstore for staging blocks for a miner
// in a storage deal, prior to sealing
func StagingBlockstore(ctx context.Context, r repo.LockedRepo) (dtypes.StagingBlockstore, error) {
	stagingds, err := r.Datastore(ctx, "/staging")
	if err != nil {
		return nil, err
	}

	return blockstore.FromDatastore(stagingds), nil
}

func StorageProviderExt(minerAddress dtypes.MinerAddress,
	storedAsk *storedask.StoredAsk,
	h host.Host, ds dtypes.MetadataDS,
	r repo.LockedRepo,
	pieceStore dtypes.ProviderPieceStore,
	dataTransfer dtypes.ProviderDataTransfer,
	spn storagemarket.StorageProviderNode,
	df dtypes.StorageDealFilter,
	dsw *dagstore.Wrapper,
	backend ethapi.Backend,
	sa retrievalmarket.SectorAccessor,
) (storagemarket.StorageProvider, error) {
	net := smnet.NewFromLibp2pHost(h)

	dir := filepath.Join(r.Path(), StagingAreaDirName)

	// migrate temporary files that were created directly under the repo, by
	// moving them to the new directory and symlinking them.
	oldDir := r.Path()
	if err := migrateDealStaging(oldDir, dir); err != nil {
		return nil, xerrors.Errorf("failed to make deal staging directory %w", err)
	}

	store, err := piecefilestore.NewLocalFileStore(piecefilestore.OsPath(dir))
	if err != nil {
		return nil, err
	}

	opt := storageImpl.CustomDealDecisionLogic(storageImpl.DealDeciderFunc(df))

	return storageImpl.NewProviderExt(
		net,
		namespace.Wrap(ds, datastore.NewKey("/deals/providerext")),
		store,
		sa,
		dsw,
		pieceStore,
		dataTransfer,
		spn,
		address.Address(minerAddress),
		storedAsk,
		backend,
		opt,
	)
}

func RetrievalProviderExt(
	minerAddress address.Address,
	node retrievalmarket.RetrievalProviderNode,
	sa retrievalmarket.SectorAccessor,
	network rmnet.RetrievalMarketNetwork,
	pieceStore piecestore.PieceStore,
	dagStore filstores.DAGStoreWrapper,
	dataTransfer datatransfer.Manager,
	ds datastore.Batching,
	retrievalPricingFunc retrievalimpl.RetrievalPricingFunc,
) (retrievalmarket.RetrievalProvider, error) {
	opt := retrievalimpl.DealDeciderOpt(nil)
	return retrievalimpl.NewProviderExt(minerAddress, node, sa, network, pieceStore, dagStore, dataTransfer, ds, retrievalPricingFunc, opt)
}

func migrateDealStaging(oldPath, newPath string) error {
	dirInfo, err := os.Stat(newPath)
	if err == nil {
		if !dirInfo.IsDir() {
			return xerrors.Errorf("%s is not a directory", newPath)
		}
		// The newPath exists already, below migration has already occurred.
		return nil
	}

	// if the directory doesn't exist, create it
	if os.IsNotExist(err) {
		if err := os.MkdirAll(newPath, 0755); err != nil {
			return xerrors.Errorf("failed to mk directory %s for deal staging: %w", newPath, err)
		}
	} else { // if we failed for other reasons, abort.
		return err
	}

	// if this is the first time we created the directory, symlink all staged deals into it. "Migration"
	// get a list of files in the miner repo
	dirEntries, err := os.ReadDir(oldPath)
	if err != nil {
		return xerrors.Errorf("failed to list directory %s for deal staging: %w", oldPath, err)
	}

	for _, entry := range dirEntries {
		// ignore directories, they are not the deals.
		if entry.IsDir() {
			continue
		}
		// the FileStore from fil-storage-market creates temporary staged deal files with the pattern "fstmp"
		// https://github.com/filecoin-project/go-fil-markets/blob/00ff81e477d846ac0cb58a0c7d1c2e9afb5ee1db/filestore/filestore.go#L69
		name := entry.Name()
		if strings.Contains(name, "fstmp") {
			// from the miner repo
			oldPath := filepath.Join(oldPath, name)
			// to its subdir "deal-staging"
			newPath := filepath.Join(newPath, name)
			// create a symbolic link in the new deal staging directory to preserve existing staged deals.
			// all future staged deals will be created here.
			if err := os.Rename(oldPath, newPath); err != nil {
				return xerrors.Errorf("failed to move %s to %s: %w", oldPath, newPath, err)
			}
			if err := os.Symlink(newPath, oldPath); err != nil {
				return xerrors.Errorf("failed to symlink %s to %s: %w", oldPath, newPath, err)
			}
			log.Infow("symlinked staged deal", "from", oldPath, "to", newPath)
		}
	}
	return nil
}

// NewProviderPieceStore creates a statestore for storing metadata about pieces
// shared by the storage and retrieval providers
func NewProviderPieceStore(ds dtypes.MetadataDS) (dtypes.ProviderPieceStore, error) {
	ps, err := piecestoreimpl.NewPieceStore(namespace.Wrap(ds, datastore.NewKey("/storagemarket")))
	if err != nil {
		return nil, err
	}
	ps.OnReady(marketevents.ReadyLogger("piecestore"))
	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//		return ps.Start(ctx)
	//	},
	//})
	return ps, nil
}

type StorageMinerParams struct {
	API                v1api.FullNode
	MetadataDS         dtypes.MetadataDS
	Sealer             sectorstorage.SectorManager
	SectorIDCounter    sealing.SectorIDCounter
	Verifier           ffiwrapper.Verifier
	Prover             ffiwrapper.Prover
	GetSealingConfigFn dtypes.GetSealingConfigFunc
	Journal            journal.Journal
	AddrSel            *storage.AddressSelector
}

func StorageMinerEx(ctx context.Context, fc config.MinerFeeConfig, params StorageMinerParams, backend ethapi.Backend) (*storageext.MinerExt, error) {
	var (
		ds     = params.MetadataDS
		api    = params.API
		sealer = params.Sealer
		sc     = params.SectorIDCounter
		verif  = params.Verifier
		prover = params.Prover
		gsd    = params.GetSealingConfigFn
		j      = params.Journal
		as     = params.AddrSel
	)

	maddr, err := minerAddrFromDS(ds)
	if err != nil {
		return nil, err
	}

	//fps, err := storage.NewWindowedPoStScheduler(api, fc, as, sealer, verif, sealer, j, maddr)
	if err != nil {
		return nil, err
	}

	sm, err := storageext.NewMinerExt(api, maddr, ds, sealer, sc, verif, prover, gsd, fc, j, as, backend)
	if err != nil {
		return nil, err
	}

	//go fps.Run(ctx)
	sm.Run(ctx)

	//lc.Append(fx.Hook{
	//	OnStart: func(context.Context) error {
	//		go fps.Run(ctx)
	//		return sm.Run(ctx)
	//	},
	//	OnStop: sm.Stop,
	//})

	return sm, nil
}

func minerAddrFromDS(ds dtypes.MetadataDS) (address.Address, error) {
	maddrb, err := ds.Get(datastore.NewKey("miner-address"))
	if err != nil {
		return address.Undef, err
	}

	return address.NewFromBytes(maddrb)
}

func LocalStorage(ctx context.Context, ls stores.LocalStorage, si stores.SectorIndex, urls stores.URLs) (*stores.Local, error) {
	return stores.NewLocal(ctx, ls, si, urls)
}

func SectorStorage(ctx context.Context, lstor *stores.Local, stor *stores.Remote, ls stores.LocalStorage, si stores.SectorIndex, sc sectorstorage.SealerConfig, ds dtypes.MetadataDS) (*sectorstorage.Manager, error) {
	wsts := statestore.New(namespace.Wrap(ds, modules.WorkerCallsPrefix))
	smsts := statestore.New(namespace.Wrap(ds, modules.ManagerWorkPrefix))

	sst, err := sectorstorage.New(ctx, lstor, stor, ls, si, sc, wsts, smsts)
	if err != nil {
		return nil, err
	}

	//lc.Append(fx.Hook{
	//	OnStop: sst.Close,
	//})

	return sst, nil
}
