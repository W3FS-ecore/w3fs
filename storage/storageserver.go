package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/node"
	w3fsAuth "github.com/ethereum/go-ethereum/storage/auth"
	"github.com/ethereum/go-ethereum/storage/implext"
	"github.com/ethereum/go-ethereum/storage/lp2pext"
	"github.com/ethereum/go-ethereum/storage/mock"
	"github.com/ethereum/go-ethereum/storage/modulesext"
	"github.com/ethereum/go-ethereum/storage/storageext"
	"github.com/filecoin-project/dagstore"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/chain/gen"
	"github.com/filecoin-project/lotus/extern/sector-storage"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"github.com/filecoin-project/lotus/journal"
	"github.com/filecoin-project/lotus/journal/alerting"
	"github.com/filecoin-project/lotus/lib/lotuslog"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/filecoin-project/lotus/markets/retrievaladapter"
	"github.com/filecoin-project/lotus/markets/sectoraccessor"
	"github.com/filecoin-project/lotus/markets/storageadapter"
	lotusnode "github.com/filecoin-project/lotus/node"
	"github.com/filecoin-project/lotus/node/config"
	"github.com/filecoin-project/lotus/node/impl/common"
	"github.com/filecoin-project/lotus/node/impl/full"
	"github.com/filecoin-project/lotus/node/impl/net"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/storage"
	"github.com/filecoin-project/lotus/storage/sectorblocks"
	storage2 "github.com/filecoin-project/specs-storage/storage"
	"github.com/google/uuid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/libp2p/go-libp2p/p2p/net/conngater"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var lotusLog = logging.Logger("eth")

// fetchKeystore retrieves the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) (*keystore.KeyStore, error) {
	if ks := am.Backends(keystore.KeyStoreType); len(ks) > 0 {
		return ks[0].(*keystore.KeyStore), nil
	}

	return nil, errors.New("local keystore not used")
}

func InitAccountKey(backend ethapi.Backend) {
	lotusLog.Info("Initializing auth")

	addr, err := backend.Coinbase()

	ks, err := fetchKeystore(backend.AccountManager())
	if err != nil {
		lotusLog.Errorf(" fetchKeystore error: %w", err)
		return
	}

	if ks.CheckAccountIsUnlock(accounts.Account{Address: addr}) {
		privateKey, err := ks.GetAccountPrivateKeyWithoutPass(accounts.Account{Address: addr})
		if err != nil {
			lotusLog.Errorf("GetAccountPrivateKeyWithoutPass error: %w", err)
			return
		}

		pribyte := crypto.FromECDSA(privateKey)
		_, res := w3fsAuth.W3FS_Register(string(pribyte))
		if res != 0 {
			lotusLog.Errorf("W3FS_Register error: %w", res)
			return
		}
	}
}

func InitStorageProviderConfig(cfg *node.Config) (*repo.FsRepo, repo.LockedRepo, *config.StorageMiner) {
	lotuslog.SetupLogLevels()
	lotusLog.Info("Initializing storage provider config")
	//ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	//defer cancel()

	//sectorSizeInt, err := units.RAMInBytes("8MiB")
	//if err != nil {
	//	return
	//}
	//ssize := abi.SectorSize(sectorSizeInt)

	lotusLog.Info("Checking if repo exists")
	homedir := cfg.DataDir + "/.w3fsminer"
	r, err := repo.NewFS(homedir)
	if err != nil {
		lotusLog.Errorf("NewFS error %s", err)
		return nil, nil, nil
	}
	lotusLog.Info("Initializing repo config")
	isExist, err := r.Exists()
	if err != nil {
		lotusLog.Errorf("stat repo exist error: %w", err)
		return nil, nil, nil
	}
	if !isExist {
		err = r.Init(repo.StorageMiner)
		if err != nil && err != repo.ErrRepoExists {
			lotusLog.Errorf("repo init error: %w", err)
			return nil, nil, nil
		}
	}
	lr, err := r.Lock(repo.StorageMiner)
	if err != nil {
		lotusLog.Errorf("Lock error %s", err)
		return nil, nil, nil
	}
	c, err := lr.Config()
	if err != nil {
		lotusLog.Errorf("load config for repo, got: %w", err)
		return nil, nil, nil
	}
	//defer lr.Close() //nolint:errcheck
	storageMinerConfig, ok := c.(*config.StorageMiner)
	if !ok {
		lotusLog.Errorf("invalid config for repo, got: %T", c)
		return nil, nil, nil
	}
	//init storage.json sectorstore.json
	lotusLog.Info("Initializing storage.json & sectorstore.json")
	sectorstoreFile := filepath.Join(lr.Path(), "sectorstore.json")
	if _, err := os.Stat(sectorstoreFile); err == nil { // No error. File already exists.
		lotusLog.Infof("already existing sectorstore.json")
	} else if os.IsNotExist(err) {
		// pass
		b, err := json.MarshalIndent(&stores.LocalStorageMeta{
			ID:       stores.ID(uuid.New().String()),
			Weight:   10,
			CanSeal:  true,
			CanStore: true,
		}, "", "  ")
		if err != nil {
			lotusLog.Errorf("marshaling storage config: %w", err)
			return nil, nil, nil
		}

		if err := ioutil.WriteFile(filepath.Join(lr.Path(), "sectorstore.json"), b, 0644); err != nil {
			lotusLog.Errorf("persisting storage metadata (%s): %w", filepath.Join(lr.Path(), "sectorstore.json"), err)
			return nil, nil, nil
		}
	}
	storageFile := filepath.Join(lr.Path(), "storage.json")
	if _, err := os.Stat(storageFile); err == nil { // No error. File already exists.
		lotusLog.Infof("already existing storage.json")
	} else if os.IsNotExist(err) {
		var localPaths []stores.LocalPath
		localPaths = append(localPaths, stores.LocalPath{
			Path: lr.Path(),
		})
		if err := lr.SetStorage(func(sc *stores.StorageConfig) {
			sc.StoragePaths = append(sc.StoragePaths, localPaths...)
		}); err != nil {
			lotusLog.Errorf("set storage config: %w", err)
			return nil, nil, nil
		}
	}
	return r, lr, storageMinerConfig
}

func InitStorageProvider(stack *node.Node, backend ethapi.Backend) <-chan struct{} {
	api.RunningNodeType = api.NodeMiner
	lotuslog.SetupLogLevels()
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	//defer cancel()

	r, lr, cfg := InitStorageProviderConfig(stack.Config())
	if r == nil || lr == nil || cfg == nil {
		lotusLog.Errorf("Initializing storage provider config error")
		return nil
	}

	lotusLog.Info("Initializing libp2p identity")

	k, err := modules.KeyStore(lr)
	if err != nil {
		lotusLog.Errorf("KeyStore error %s", err)
		return nil
	}
	apialg, err := modules.APISecret(k, lr)
	if err != nil {
		lotusLog.Errorf("APISecret error %s", err)
		return nil
	}
	p, err := lp2p.PrivKey(k)
	if err != nil {
		lotusLog.Errorf("PrivKey error %s", err)
		return nil
	}
	peerid, err := peer.IDFromPublicKey(p.GetPublic())
	if err != nil {
		lotusLog.Errorf("IDFromPublicKey error %s", err)
		return nil
	}

	ps := pstoremem.NewPeerstore()
	err = lp2p.PstoreAddSelfKeys(peerid, p, ps)
	if err != nil {
		lotusLog.Errorf("PstoreAddSelfKeys error %s", err)
		return nil
	}

	//addresses := []string{
	//	"/ip4/0.0.0.0/tcp/" + stack.Config().Ip4Port,
	//	"/ip6/::/tcp/0" + stack.Config().Ip6Port,
	//}
	for ind, la := range cfg.Libp2p.ListenAddresses {
		if la == "/ip4/0.0.0.0/tcp/0" {
			cfg.Libp2p.ListenAddresses[ind] = "/ip4/0.0.0.0/tcp/30308"
			break
		}
	}
	startListening := lp2p.StartListening(cfg.Libp2p.ListenAddresses)

	metadataDS, err := modulesext.Datastore(false, ctx, lr) //cfg.Backup.DisableMetadataLog
	if err != nil {
		lotusLog.Errorf("modulesext.Datastore error %s", err)
		return nil
	}
	//defer metadataDS.Close()

	//p2p options
	var Opts [][]libp2p.Option
	opts, err := lp2p.DefaultTransports()
	if err != nil || opts.Opts == nil {
		lotusLog.Errorf("PstoreAddSelfKeys error %s", err)
		return nil
	}
	Opts = append(Opts, opts.Opts)
	addrsFactory := lp2p.AddrsFactory(nil, nil)
	opts, err = addrsFactory()
	Opts = append(Opts, opts.Opts)
	smuxTransport := lp2p.SmuxTransport(true)
	opts, err = smuxTransport()
	Opts = append(Opts, opts.Opts)
	noRelay := lp2p.NoRelay()
	opts, err = noRelay()
	Opts = append(Opts, opts.Opts)
	securityfunc := lp2p.Security(true, false)
	security := securityfunc.(func() (opts lp2p.Libp2pOpts))
	opts = security()
	Opts = append(Opts, opts.Opts)
	connectionManager := lp2p.ConnectionManager(50, 200, 20*time.Second, nil)
	opts, err = connectionManager()
	if err != nil {
		lotusLog.Error("ConnectionManager option %s", err)
		return nil
	}
	Opts = append(Opts, opts.Opts)
	connGater, err := lp2p.ConnGater(metadataDS)
	if err != nil {
		lotusLog.Error("ConnGater option %s", err)
		return nil
	}
	opts, err = lp2p.ConnGaterOption(connGater)
	Opts = append(Opts, opts.Opts)

	opts, reporter := lp2p.BandwidthCounter()
	Opts = append(Opts, opts.Opts)

	opts, err = lp2p.AutoNATService()
	Opts = append(Opts, opts.Opts)

	opts, err = lp2p.NatPortMap()
	Opts = append(Opts, opts.Opts)

	///new Host
	params := lp2pext.P2PHostIn{
		ID:        peerid,
		Peerstore: ps,
		Opts:      Opts,
	}
	rawHost, err := lp2pext.Host(ctx, params)
	//defer rawHost.Close()

	baseIpfsRouting, err := lp2pext.DHTRouting(dht.ModeAuto, ctx, rawHost, metadataDS, modules.RecordValidator(ps), "testnetnet", dtypes.Bootstrapper(false))
	//ipfsDHT := baseIpfsRouting.(*dht.IpfsDHT)
	//defer ipfsDHT.Close()
	routeHost := lp2pext.RoutedHost(rawHost, baseIpfsRouting)
	err = startListening(routeHost)
	if err != nil {
		lotusLog.Errorf("startListening %s", err)
		return nil
	}
	localCid := peer.ToCid(routeHost.ID())
	localAddr := routeHost.Addrs()
	lotusLog.Infof("peerid is %s", localCid)
	lotusLog.Infof("addr is %s", localAddr)
	// set peerId and peerAddr for p2p necessity
	for _, addr := range localAddr {
		if !strings.Contains(addr.String(), "127.0.0.1") && !strings.Contains(addr.String(), "ip6") {
			ethapi.PeerId = localCid.String()
			break
		}
	}
	//save the miner addr
	addrs := peer.AddrInfo{
		ID:    routeHost.ID(),
		Addrs: routeHost.Addrs(),
	}

	var buf bytes.Buffer
	for _, peer := range addrs.Addrs {
		fmt.Fprintf(&buf, "%s/p2p/%s\n", peer, addrs.ID)
	}
	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	minerAddr := backend.GetMinerAddr()
	*minerAddr = refStr

	//borLocalNode.Node()
	stagingBlockstore, _ := modulesext.StagingBlockstore(ctx, lr)
	stagingGraphsync := modulesext.StagingGraphsync(ctx, 20, 20, stagingBlockstore, routeHost)

	dt, _ := modulesext.NewProviderDAGServiceDataTransfer(routeHost, stagingGraphsync, metadataDS, lr)
	dt.SubscribeToEvents(marketevents.DataTransferLogger)
	dt.Start(ctx)
	//defer dt.Stop(ctx)

	///storageminer
	//var fc = config.MinerFeeConfig{
	//	MaxPreCommitGasFee:      types.FIL(types.FromFil(1)),
	//	MaxCommitGasFee:         types.FIL(types.FromFil(1)),
	//	MaxTerminateGasFee:      types.FIL(types.FromFil(1)),
	//	MaxPreCommitBatchGasFee: config.BatchFeeConfig{Base: types.FIL(types.FromFil(3)), PerSector: types.FIL(types.FromFil(1))},
	//	MaxCommitBatchGasFee:    config.BatchFeeConfig{Base: types.FIL(types.FromFil(3)), PerSector: types.FIL(types.FromFil(1))},
	//}
	act := "t01000"
	a, err := address.NewFromString(act)
	if err != nil {
		lotusLog.Errorf("failed parsing actor flag value (%q): %w", act, err)
		return nil
	}
	if err := metadataDS.Put(datastore.NewKey("miner-address"), a.Bytes()); err != nil {
		lotusLog.Errorf("failed put miner-address (%q): %w", act, err)
		return nil
	}
	maddr, err := multiaddr.NewMultiaddr(cfg.API.ListenAddress)
	if err != nil {
		lotusLog.Error("NewMultiaddr %s", err)
		return nil
	}
	lr.SetAPIEndpoint(maddr)
	urls, _ := func(e dtypes.APIEndpoint) (stores.URLs, error) {
		ip := cfg.API.RemoteListenAddress

		var urls stores.URLs
		urls = append(urls, "http://"+ip+"/remote") // TODO: This makes no assumptions, and probably could...
		return urls, nil
	}(maddr)
	envDisabledEvents := journal.EnvDisabledEvents()
	journal, err := modulesext.OpenFilesystemJournal(lr, envDisabledEvents)
	if err != nil {
		lotusLog.Error("OpenFilesystemJournal %s", err)
	}
	//defer journal.Close()
	index := stores.NewIndex()
	localStorage, _ := modulesext.LocalStorage(ctx, lr, index, urls)
	ShutdownChan := make(chan struct{})
	Alerting := alerting.NewAlertingSystem(journal)
	commonApi := &common.CommonAPI{
		Alerting:     Alerting,
		APISecret:    apialg,
		ShutdownChan: ShutdownChan,
	}
	storageAuth, err := modules.StorageAuth(ctx, commonApi)
	if err != nil {
		lotusLog.Error("StorageAuth %s", err)
		return nil
	}
	remoteStorage := modules.RemoteStorage(localStorage, index, storageAuth, cfg.Storage)
	sectorstoragemanager, err := modulesext.SectorStorage(ctx, localStorage, remoteStorage, lr, index, cfg.Storage, metadataDS)
	//defer sectorstoragemanager.Close(ctx)
	getSealingConfigFn, _ := modules.NewGetSealConfigFunc(lr)
	AddressSelectorFunc := modules.AddressSelector(&cfg.Addresses)
	AddressSelector, err := AddressSelectorFunc()
	if err != nil {
		lotusLog.Errorf("AddressSelectorFunc %s", err)
		return nil
	}
	mockfullapi := backend.GetMockFullApi()
	fullNode := &mock.MockFullNode{}
	*mockfullapi = *fullNode
	storageMinerParams := modulesext.StorageMinerParams{
		API:                fullNode,
		MetadataDS:         metadataDS,
		SectorIDCounter:    modules.SectorIDCounter(metadataDS),
		Sealer:             sectorstoragemanager,
		Verifier:           ffiwrapper.ProofVerifier,
		Prover:             ffiwrapper.ProofProver,
		GetSealingConfigFn: getSealingConfigFn,
		Journal:            journal,
		AddrSel:            AddressSelector,
	}
	miner, err := modulesext.StorageMinerEx(ctx, cfg.Fees, storageMinerParams, backend)
	if err != nil {
		lotusLog.Errorf("StorageMinerEx %s", err)
		return nil
	}
	//defer miner.Stop(ctx)

	logger := mock.NewLogger()
	lc := &lifecycleWrapper{
		mock.New(logger),
	}

	//defer lc.Stop(ctx)

	DealPublisherFunc := storageadapter.NewDealPublisher(nil, storageadapter.PublishMsgConfig{})
	dealPublisher := DealPublisherFunc(lc, fullNode, AddressSelector)
	sectorBlock := sectorblocks.NewSectorBlocks(miner, metadataDS)

	NewProviderNodeAdapterFunc := storageadapter.NewProviderNodeAdapter(&cfg.Fees, &cfg.Dealmaking)
	spn, err := NewProviderNodeAdapterFunc(ctx, lc, sectorBlock, fullNode, dealPublisher)
	if err != nil {
		lotusLog.Error("NewProviderNodeAdapter %s %s", err, spn)
		return nil
	}
	//
	//lc.Start(ctx)
	///new Provider
	minerAddress, err := modules.MinerAddress(metadataDS)
	if err != nil {
		lotusLog.Error("MinerAddress %s", err)
		return nil
	}
	providerPieceStore, err := modulesext.NewProviderPieceStore(metadataDS)
	if err != nil {
		lotusLog.Error("NewProviderPieceStore %s", err)
		return nil
	}
	ready := make(chan error, 1)
	providerPieceStore.OnReady(func(err error) {
		ready <- err
	})
	providerPieceStore.Start(ctx)
	if err := <-ready; err != nil {
		lotusLog.Errorf("aborting dagstore start; piecestore failed to start: %s", err)
		return nil
	}
	// instance the SectorAccessor
	pp := sectorstorage.NewPieceProvider(remoteStorage, index, sectorstoragemanager)
	sa := sectoraccessor.NewSectorAccessor(minerAddress, miner, pp, fullNode)
	minerAPI, err := modulesext.NewMinerAPI(ctx, lr, providerPieceStore, sa)
	if err != nil {
		lotusLog.Error("NewMinerAPI %s", err)
		return nil
	}
	err = minerAPI.Start(ctx)
	if err != nil {
		lotusLog.Error("minerAPI.Start %s", err)
		return nil
	}
	dagStore, wrapper, err := modulesext.DAGStore(lr, minerAPI)
	wrapper.Start(ctx)
	if err != nil {
		lotusLog.Error("wrapper.Start %s", err)
		return nil
	}

	storageProvider, err := modulesext.StorageProviderExt(minerAddress, nil, routeHost, metadataDS, lr, providerPieceStore, dt, spn, nil, wrapper, backend, sa)
	if err != nil {
		lotusLog.Error("StorageProviderExt %s", err)
		return nil
	}

	storageProvider.Start(ctx)
	//defer storageProvider.Stop()

	lotusLog.Infof("Start storageProvider success!")

	node := retrievaladapter.NewRetrievalProviderNode(nil)
	network := modules.RetrievalNetwork(routeHost)
	ds := namespace.Wrap(metadataDS, datastore.NewKey("/deals/retrievalProviderext"))

	retrievalProvider, err := modulesext.RetrievalProviderExt(address.Address(minerAddress), node, sa, network, providerPieceStore, wrapper, dt, ds, nil)
	if err != nil {
		lotusLog.Error("RetrievalProviderExt error: %s", err)
		return nil
	}
	retrievalProvider.Start(ctx)
	//defer retrievalProvider.Stop()
	lotusLog.Infof("Start retrievalProvider success!")
	//////enable miner api
	shutdownChan := make(chan struct{})
	//var minerapi api.StorageMiner
	//minerapi = &impl.StorageMinerAPI{}
	minerapi := backend.GetMinerApi()
	stop, err := New(ctx, minerapi,
		Override(new(dtypes.ShutdownChan), shutdownChan),
		Override(new(api.Net), lotusnode.From(new(net.NetAPI))),
		Override(new(api.Common), lotusnode.From(new(common.CommonAPI))),
		Override(new(api.MinerSubsystems), modules.ExtractEnabledMinerSubsystems(cfg.Subsystems)),
		Override(new(v1api.FullNode), fullNode),
		Override(new(stores.SectorIndex), index),
		Override(new(*storage.AddressSelector), AddressSelectorFunc),
		Override(new(dtypes.MetadataDS), metadataDS),
		Override(new(*stores.Local), localStorage),
		Override(new(*stores.Remote), remoteStorage),
		Override(new(*alerting.Alerting), Alerting),
		Override(new(*dtypes.APIAlg), apialg),
		Override(new(lp2p.RawHost), rawHost),
		Override(new(host.Host), routeHost),
		Override(new(lp2p.BaseIpfsRouting), baseIpfsRouting),
		Override(new(*conngater.BasicConnectionGater), connGater),
		Override(new(*dtypes.ScoreKeeper), lp2p.ScoreKeeper),
		Override(new(metrics.Reporter), reporter),
		Override(new(*storageext.MinerExt), miner),
		Override(new(*sectorstorage.Manager), sectorstoragemanager),
		Override(new(sectorstorage.SectorManager), lotusnode.From(new(*sectorstorage.Manager))),
		Override(new(storiface.WorkerReturn), lotusnode.From(new(sectorstorage.SectorManager))),
		Override(new(*sectorblocks.SectorBlocks), sectorBlock),
		Override(new(dtypes.SetSealingConfigFunc), modules.NewSetSealConfigFunc),
		Override(new(dtypes.GetSealingConfigFunc), modules.NewGetSealConfigFunc),
		Override(new(repo.LockedRepo), modules.LockedRepo(lr)), // module handles closing
		Override(new(dtypes.ProviderPieceStore), providerPieceStore),
		Override(new(*dagstore.DAGStore), dagStore),
		Override(new(retrievalmarket.SectorAccessor), sa),
		Override(new(dtypes.ProviderDataTransfer), dt),
		Override(new(dtypes.MinerAddress), minerAddress),
		Override(new(dtypes.MinerID), modules.MinerID),
		Override(new(ffiwrapper.Verifier), ffiwrapper.ProofVerifier),
		Override(new(storage2.Prover), lotusnode.From(new(sectorstorage.SectorManager))),
		Override(new(gen.WinningPoStProver), storageext.NewWinningPoStProver),
	)
	if err != nil {
		lotusLog.Error("creating node: %w", err)
		return nil
	}
	endpoint, err := r.APIEndpoint()
	if err != nil {
		lotusLog.Errorf("getting API endpoint: %w", err)
		return nil
	}

	// Instantiate the miner node handler.
	handler, err := implext.MinerHandler(*minerapi, true)
	if err != nil {
		lotusLog.Errorf("failed to instantiate rpc handler: %w", err)
		return nil
	}

	// Serve the RPC.
	rpcStopper, err := lotusnode.ServeRPC(handler, "bor-miner", endpoint)
	if err != nil {
		lotusLog.Errorf("failed to start json-rpc endpoint: %s", err)
		return nil
	}

	// Monitor for shutdown.
	finishCh := lotusnode.MonitorShutdown(shutdownChan,
		lotusnode.ShutdownHandler{Component: "rpc server", StopFunc: rpcStopper},
		lotusnode.ShutdownHandler{Component: "miner", StopFunc: stop},
	)

	return finishCh
}

func InitStorageClient(stack *node.Node, backend ethapi.Backend) *ethapi.ClientManager {
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()
	homedir := stack.Config().DataDir + "/.w3fs"
	ethapi.RetrievalFileDir = homedir + "/retrieval-file/"
	ethapi.StoreDir = homedir + "/storage-file/"
	r, err := repo.NewFS(homedir)
	if err != nil {
		lotusLog.Errorf("NewFS error %s", err)
		return nil
	}
	// set output path for retrievals file
	dir := filepath.Join(homedir, "/retrieval-file")
	if err := os.MkdirAll(dir, 0755); err != nil {
		lotusLog.Errorf("failed to create retrieval's output directory %s: %w", dir, err)
		return nil
	}
	// set output path for retrievals file
	dir = filepath.Join(homedir, "/storage-file")
	if err := os.MkdirAll(dir, 0755); err != nil {
		lotusLog.Errorf("failed to create storage's input directory %s: %w", dir, err)
		return nil
	}
	err = r.Init(repo.FullNode)
	if err != nil && err != repo.ErrRepoExists {
		lotusLog.Errorf("repo init error: %w", err)
		return nil
	}
	lr, err := r.Lock(repo.FullNode)
	if err != nil {
		lotusLog.Errorf("Lock error %s", err)
		return nil
	}
	k, err := modules.KeyStore(lr)
	if err != nil {
		lotusLog.Errorf("KeyStore error %s", err)
		return nil
	}
	p, err := lp2p.PrivKey(k)
	if err != nil {
		lotusLog.Errorf("PrivKey error %s", err)
		return nil
	}
	peerid, err := peer.IDFromPublicKey(p.GetPublic())
	if err != nil {
		lotusLog.Errorf("IDFromPublicKey error %s", err)
		return nil
	}
	lotusLog.Infof("client peerid is  %s", peerid)
	ps := pstoremem.NewPeerstore()
	err = lp2p.PstoreAddSelfKeys(peerid, p, ps)
	if err != nil {
		lotusLog.Errorf("PstoreAddSelfKeys error %s", err)
		return nil
	}
	addresses := []string{
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
	startListening := lp2p.StartListening(addresses)

	metadataDS, err := modulesext.Datastore(false, ctx, lr) //cfg.Backup.DisableMetadataLog
	//defer metadataDS.Close()

	//p2p options
	var Opts [][]libp2p.Option
	opts, err := lp2p.DefaultTransports()
	if err != nil || opts.Opts == nil {
		lotusLog.Errorf("PstoreAddSelfKeys error %s", err)
	}
	Opts = append(Opts, opts.Opts)
	addrsFactory := lp2p.AddrsFactory(nil, nil)
	opts, err = addrsFactory()
	Opts = append(Opts, opts.Opts)
	smuxTransport := lp2p.SmuxTransport(true)
	opts, err = smuxTransport()
	Opts = append(Opts, opts.Opts)
	noRelay := lp2p.NoRelay()
	opts, err = noRelay()
	Opts = append(Opts, opts.Opts)
	securityfunc := lp2p.Security(true, false)
	security := securityfunc.(func() (opts lp2p.Libp2pOpts))
	opts = security()
	Opts = append(Opts, opts.Opts)
	connectionManager := lp2p.ConnectionManager(50, 200, 20*time.Second, nil)
	opts, err = connectionManager()
	if err != nil {
		lotusLog.Errorf("ConnectionManager option %s", err)
	}
	Opts = append(Opts, opts.Opts)
	connGater, err := lp2p.ConnGater(metadataDS)
	if err != nil {
		lotusLog.Errorf("ConnGater option %s", err)
	}
	opts, err = lp2p.ConnGaterOption(connGater)
	Opts = append(Opts, opts.Opts)

	opts, _ = lp2p.BandwidthCounter()
	Opts = append(Opts, opts.Opts)

	opts, err = lp2p.AutoNATService()
	Opts = append(Opts, opts.Opts)

	opts, err = lp2p.NatPortMap()
	Opts = append(Opts, opts.Opts)

	///new Host
	params := lp2pext.P2PHostIn{
		ID:        peerid,
		Peerstore: ps,
		Opts:      Opts,
	}
	rawHost, err := lp2pext.Host(ctx, params)
	//defer rawHost.Close()

	baseIpfsRouting, err := lp2pext.DHTRouting(dht.ModeAuto, ctx, rawHost, metadataDS, modules.RecordValidator(ps), "testnetnet", dtypes.Bootstrapper(false))
	//ipfsDHT := baseIpfsRouting.(*dht.IpfsDHT)
	//defer ipfsDHT.Close()
	host := lp2pext.RoutedHost(rawHost, baseIpfsRouting)

	err = startListening(host)
	if err != nil {
		lotusLog.Errorf("startListening %s", err)
	}
	addr := host.Addrs()
	lotusLog.Infof("host addr %s", addr)

	stagingBlockstore, _ := modulesext.StagingBlockstore(ctx, lr)
	stagingGraphsync := modulesext.StagingGraphsync(ctx, 20, 20, stagingBlockstore, host)

	dt, _ := modulesext.NewProviderDAGServiceDataTransfer(host, stagingGraphsync, metadataDS, lr)
	dt.SubscribeToEvents(marketevents.DataTransferLogger)
	dt.Start(ctx)
	//defer dt.Stop(ctx)

	////new Client
	clientBlockstore := modules.ClientBlockstore()
	clientDatastore := modules.NewClientDatastore(metadataDS)
	clientImportMgr, err := modulesext.ClientImportMgr(metadataDS, lr)
	if err != nil {
		lotusLog.Errorf("ClientImportMgr %s", err)
	}

	universalBlockstore, err := modulesext.UniversalBlockstore(ctx, lr)
	if err != nil {
		lotusLog.Errorf("UniversalBlockstore %s", err)
	}
	/*if c, ok := universalBlockstore.(io.Closer); ok {
		defer c.Close()
	}*/
	graphsync, err := modulesext.Graphsync(ctx, 20, 20, lr, clientBlockstore, universalBlockstore, host)
	if err != nil {
		lotusLog.Errorf("Graphsync %s", err)
	}
	clientGraphsyncDataTransfer, err := modulesext.NewClientGraphsyncDataTransfer(ctx, host, graphsync, metadataDS, lr)
	if err != nil {
		lotusLog.Errorf("NewClientGraphsyncDataTransfer %s", err)
	}
	//defer clientGraphsyncDataTransfer.Stop(ctx)
	localDiscovery, _ := modulesext.NewLocalDiscovery(ctx, metadataDS)
	storageBlockstoreAccessor := modules.StorageBlockstoreAccessor(clientImportMgr)
	envDisabledEvents := journal.EnvDisabledEvents()
	journal, err := modulesext.OpenFilesystemJournal(lr, envDisabledEvents)
	if err != nil {
		lotusLog.Errorf("OpenFilesystemJournal %s", err)
	}
	//defer journal.Close()
	storageClient, err := modulesext.StorageClient(host, clientGraphsyncDataTransfer, localDiscovery, clientDatastore, nil, storageBlockstoreAccessor, journal)
	if err != nil {
		lotusLog.Errorf("StorageClient %s", err)
	}
	//storageClient.Start(ctx)
	//defer storageClient.Stop()

	retrievalBlockstoreAccessor, err := modules.RetrievalBlockstoreAccessor(lr)
	if err != nil {
		lotusLog.Errorf("RetrievalBlockstoreAccessor err %s", err)
	}

	resolver := modules.RetrievalResolver(localDiscovery)
	chainAPI := full.ChainAPI{}
	stateAPI := full.StateAPI{}
	mockFullApi := backend.GetMockFullApi()

	retrievalClient, err := modulesext.RetrievalClient(host, clientGraphsyncDataTransfer, resolver, metadataDS, chainAPI, stateAPI, *mockFullApi, retrievalBlockstoreAccessor, journal)

	if err != nil {
		lotusLog.Errorf("RetrievalClient %s", err)
	}

	modulesext.RetrievalLinkStorage(storageClient, retrievalClient, ethapi.RetrievalFileDir)

	// new Client
	//retrievalClient.Start(ctx)

	mgr := &ethapi.ClientManager{
		Retrieval:                 retrievalClient,
		DealClient:                storageClient,
		Imports:                   clientImportMgr,
		StorageBlockstoreAccessor: storageBlockstoreAccessor,
		RtvlBlockstoreAccessor:    retrievalBlockstoreAccessor, // retrieval relate
		DataTransfer:              dt,
		Host:                      host,
		Repo:                      lr,
		RetDiscovery:              resolver,
		FullApi:                   *mockFullApi,
	}
	mgr.SubscribeToAllTransferEvents(ctx)
	return mgr
}

type lifecycleWrapper struct{ *mock.Lifecycle }

func (l *lifecycleWrapper) Append(h fx.Hook) {
	l.Lifecycle.Append(mock.Hook{
		OnStart: h.OnStart,
		OnStop:  h.OnStop,
	})
}

// New builds and starts new Filecoin node
func New(ctx context.Context, out *api.StorageMiner, opts ...Option) (lotusnode.StopFunc, error) {
	resAPI := &implext.StorageMinerAPI{}
	*out = resAPI
	invokes := make([]fx.Option, 26)
	invokes[lotusnode.ExtractApiKey] = fx.Populate(resAPI)
	// fill holes in invokes for use in fx.Options
	for i, opt := range invokes {
		if opt == nil {
			invokes[i] = fx.Options()
		}
	}
	options := make([]fx.Option, len(opts))
	for i, opt := range opts {
		options[i] = opt()
	}

	app := fx.New(
		fx.Options(options...),
		fx.Options(invokes...),

		fx.NopLogger,
	)

	// TODO: we probably should have a 'firewall' for Closing signal
	//  on this context, and implement closing logic through lifecycles
	//  correctly
	if err := app.Start(ctx); err != nil {
		// comment fx.NopLogger few lines above for easier debugging
		return nil, xerrors.Errorf("starting node: %w", err)
	}

	return app.Stop, nil
}

func as(in interface{}, as interface{}) interface{} {
	outType := reflect.TypeOf(as)

	if outType.Kind() != reflect.Ptr {
		panic("outType is not a pointer")
	}

	inType := reflect.TypeOf(in)

	if inType.Kind() != reflect.Func || inType.AssignableTo(outType.Elem()) {
		ctype := reflect.FuncOf(nil, []reflect.Type{outType.Elem()}, false)

		return reflect.MakeFunc(ctype, func(args []reflect.Value) (results []reflect.Value) {
			out := reflect.New(outType.Elem())
			out.Elem().Set(reflect.ValueOf(in))

			return []reflect.Value{out.Elem()}
		}).Interface()
	}

	ins := make([]reflect.Type, inType.NumIn())
	outs := make([]reflect.Type, inType.NumOut())

	for i := range ins {
		ins[i] = inType.In(i)
	}
	outs[0] = outType.Elem()
	for i := range outs[1:] {
		outs[i+1] = inType.Out(i + 1)
	}

	ctype := reflect.FuncOf(ins, outs, false)

	return reflect.MakeFunc(ctype, func(args []reflect.Value) (results []reflect.Value) {
		outs := reflect.ValueOf(in).Call(args)

		out := reflect.New(outType.Elem())
		if outs[0].Type().AssignableTo(outType.Elem()) {
			// Out: Iface = In: *Struct; Out: Iface = In: OtherIface
			out.Elem().Set(outs[0])
		} else {
			// Out: Iface = &(In: Struct)
			t := reflect.New(outs[0].Type())
			t.Elem().Set(outs[0])
			out.Elem().Set(t)
		}
		outs[0] = out.Elem()

		return outs
	}).Interface()
}

type special struct{ id int }
type invoke int

type Option func() fx.Option

// Override option changes constructor for a given type
func Override(typ, constructor interface{}) Option {
	return func() fx.Option {
		if _, ok := typ.(invoke); ok {
			return fx.Invoke(constructor)
		}

		if _, ok := typ.(special); ok {
			return fx.Provide(constructor)
		}
		ctor := as(constructor, typ)
		return fx.Provide(ctor)
	}
}
