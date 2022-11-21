package ethapi

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	file_store "github.com/ethereum/go-ethereum/borcontracts/file-store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/storage/clientext"
	retrievalimpl2 "github.com/ethereum/go-ethereum/storage/clientext/retrieval"
	"github.com/ethereum/go-ethereum/storage/mock"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/discovery"
	discoveryimpl "github.com/filecoin-project/go-fil-markets/discovery/impl"
	rm "github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/stores"
	"github.com/filecoin-project/go-state-types/abi"
	bigext "github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api"
	lapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/journal"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/filecoin-project/lotus/markets/retrievaladapter"
	"github.com/filecoin-project/lotus/markets/storageadapter"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/node/repo/imports"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-cidutil"
	"github.com/ipfs/go-cidutil/cidenc"
	"github.com/ipfs/go-datastore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	chunker "github.com/ipfs/go-ipfs-chunker"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	files "github.com/ipfs/go-ipfs-files"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	unixfile "github.com/ipfs/go-unixfs/file"
	"github.com/ipfs/go-unixfs/importer/balanced"
	ihelper "github.com/ipfs/go-unixfs/importer/helpers"
	"github.com/ipld/go-car"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var DefaultHashFunction = uint64(mh.BLAKE2B_MIN + 31)

// 8 days ~=  SealDuration + PreCommit + MaxProveCommitDuration + 8 hour buffer
const dealStartBufferHours uint64 = 8 * 24
const DefaultDAGStoreDir = "dagstore"
const DefaultMaxRetrievePrice = "0.01"

type ClientManager struct {
	DealClient                storagemarket.StorageClient
	Retrieval                 rm.RetrievalClient
	Imports                   dtypes.ClientImportMgr
	StorageBlockstoreAccessor storagemarket.BlockstoreAccessor
	RtvlBlockstoreAccessor    rm.BlockstoreAccessor

	DataTransfer dtypes.ClientDataTransfer
	Host         host.Host

	Repo         repo.LockedRepo
	RetDiscovery discovery.PeerResolver

	ClientGraphsyncDataTransfer dtypes.ClientDataTransfer
	Journal                     journal.Journal
	LocalDiscovery              *discoveryimpl.Local
	MetadataDS                  dtypes.MetadataDS
	FullApi                     mock.MockFullNode
}

type retrievalSubscribeEvent struct {
	event rm.ClientEvent
	state rm.ClientDealState
}

const (
	// Init
	DEAL_INIT    int = 0
	DEAL_ING     int = 100
	DEAL_SUCCESS int = 200
	DEAL_TIMEOUT int = 400
	DEAL_ERROR   int = 500
)

type DealStatus struct {
	Status    int       // What type of status it is
	Message   string    // Any clarifying information about the event
	Timestamp time.Time // when the status return
}

type FileInfo struct {
	FileExt  string   // file extensions.such as txt/doc/mp3
	FileSize *big.Int // file size
}

type FileInfoExt struct {
	FileExt   string   // file extensions.such as txt/doc/mp3
	FileSize  *big.Int // file size
	StoreType int      // 1 onefile store 2 head-body store
}

func newDealStatus(status int, message string) DealStatus {
	return DealStatus{status, message, time.Now().UTC()}
}

type RetrieveStatus struct {
	Cid       string
	Status    int       // What type of status it is
	Message   string    // Any clarifying information about the event
	Timestamp time.Time // when the status return
}

func newRetrieveStatus(cid string, status int, message string) RetrieveStatus {
	return RetrieveStatus{cid, status, message, time.Now().UTC()}
}

func NewClientManager() {

}

func (a *ClientManager) importManager() *imports.Manager {
	return a.Imports
}

func (a *ClientManager) ClientImport(ctx context.Context, ref api.FileRef) (res *api.ImportRes, err error) {
	var (
		imgr    = a.importManager()
		id      imports.ID
		root    cid.Cid
		carPath string
	)

	id, err = imgr.CreateImport()
	if err != nil {
		return nil, xerrors.Errorf("failed to create import: %w", err)
	}

	if ref.IsCAR {
		// user gave us a CAR file, use it as-is
		// validate that it's either a carv1 or carv2, and has one root.
		f, err := os.Open(ref.Path)
		if err != nil {
			return nil, xerrors.Errorf("failed to open CAR file: %w", err)
		}
		defer f.Close() //nolint:errcheck

		hd, _, err := car.ReadHeader(bufio.NewReader(f))
		if err != nil {
			return nil, xerrors.Errorf("failed to read CAR header: %w", err)
		}
		if len(hd.Roots) != 1 {
			return nil, xerrors.New("car file can have one and only one header")
		}
		if hd.Version != 1 && hd.Version != 2 {
			return nil, xerrors.Errorf("car version must be 1 or 2, is %d", hd.Version)
		}

		carPath = ref.Path
		root = hd.Roots[0]
	} else {
		carPath, err = imgr.AllocateCAR(id)
		if err != nil {
			return nil, xerrors.Errorf("failed to create car path for import: %w", err)
		}

		// remove the import if something went wrong.
		defer func() {
			if err != nil {
				_ = os.Remove(carPath)
				_ = imgr.Remove(id)
			}
		}()

		// perform the unixfs chunking.
		root, err = a.createUnixFSFilestore(ctx, ref.Path, carPath)
		if err != nil {
			return nil, xerrors.Errorf("failed to import file using unixfs: %w", err)
		}
	}

	if err = imgr.AddLabel(id, imports.LSource, "import"); err != nil {
		return nil, err
	}
	if err = imgr.AddLabel(id, imports.LFileName, ref.Path); err != nil {
		return nil, err
	}
	if err = imgr.AddLabel(id, imports.LCARPath, carPath); err != nil {
		return nil, err
	}
	if err = imgr.AddLabel(id, imports.LRootCid, root.String()); err != nil {
		return nil, err
	}
	return &api.ImportRes{
		Root:     root,
		ImportID: id,
	}, nil
}

// createUnixFSFilestore takes a standard file whose path is src, forms a UnixFS DAG, and
// writes a CARv2 file with positional mapping (backed by the go-filestore library).
func (a *ClientManager) createUnixFSFilestore(ctx context.Context, srcPath string, dstPath string) (cid.Cid, error) {
	// This method uses a two-phase approach with a staging CAR blockstore and
	// a final CAR blockstore.
	//
	// This is necessary because of https://github.com/ipld/go-car/issues/196
	//
	// TODO: do we need to chunk twice? Isn't the first output already in the
	//  right order? Can't we just copy the CAR file and replace the header?

	src, err := os.Open(srcPath)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to open input file: %w", err)
	}
	defer src.Close() //nolint:errcheck

	stat, err := src.Stat()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to stat file :%w", err)
	}

	file, err := files.NewReaderPathFile(srcPath, src, stat)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create reader path file: %w", err)
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create temp file: %w", err)
	}
	_ = f.Close() // close; we only want the path.

	tmp := f.Name()
	defer os.Remove(tmp) //nolint:errcheck

	// Step 1. Compute the UnixFS DAG and write it to a CARv2 file to get
	// the root CID of the DAG.
	fstore, err := stores.ReadWriteFilestore(tmp)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create temporary filestore: %w", err)
	}

	finalRoot1, err := buildUnixFS(ctx, file, fstore, true)
	if err != nil {
		_ = fstore.Close()
		return cid.Undef, xerrors.Errorf("failed to import file to store to compute root: %w", err)
	}

	if err := fstore.Close(); err != nil {
		return cid.Undef, xerrors.Errorf("failed to finalize car filestore: %w", err)
	}

	// Step 2. We now have the root of the UnixFS DAG, and we can write the
	// final CAR for real under `dst`.
	bs, err := stores.ReadWriteFilestore(dstPath, finalRoot1)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create a carv2 read/write filestore: %w", err)
	}

	// rewind file to the beginning.
	if _, err := src.Seek(0, 0); err != nil {
		return cid.Undef, xerrors.Errorf("failed to rewind file: %w", err)
	}

	finalRoot2, err := buildUnixFS(ctx, file, bs, true)
	if err != nil {
		_ = bs.Close()
		return cid.Undef, xerrors.Errorf("failed to create UnixFS DAG with carv2 blockstore: %w", err)
	}

	if err := bs.Close(); err != nil {
		return cid.Undef, xerrors.Errorf("failed to finalize car blockstore: %w", err)
	}

	if finalRoot1 != finalRoot2 {
		return cid.Undef, xerrors.New("roots do not match")
	}

	return finalRoot1, nil
}

// buildUnixFS builds a UnixFS DAG out of the supplied reader,
// and imports the DAG into the supplied service.
func buildUnixFS(ctx context.Context, reader io.Reader, into bstore.Blockstore, filestore bool) (cid.Cid, error) {
	b, err := unixFSCidBuilder()
	if err != nil {
		return cid.Undef, err
	}

	bsvc := blockservice.New(into, offline.Exchange(into))
	dags := merkledag.NewDAGService(bsvc)
	bufdag := ipld.NewBufferedDAG(ctx, dags)

	params := ihelper.DagBuilderParams{
		Maxlinks:   build.UnixfsLinksPerLevel,
		RawLeaves:  true,
		CidBuilder: b,
		Dagserv:    bufdag,
		NoCopy:     filestore,
	}

	db, err := params.New(chunker.NewSizeSplitter(reader, int64(build.UnixfsChunkSize)))
	if err != nil {
		return cid.Undef, err
	}
	nd, err := balanced.Layout(db)
	if err != nil {
		return cid.Undef, err
	}

	if err := bufdag.Commit(); err != nil {
		return cid.Undef, err
	}

	return nd.Cid(), nil
}

func unixFSCidBuilder() (cid.Builder, error) {
	prefix, err := merkledag.PrefixForCidVersion(1)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UnixFS CID Builder: %w", err)
	}
	prefix.MhType = DefaultHashFunction
	b := cidutil.InlineBuilder{
		Builder: prefix,
		Limit:   126,
	}
	return b, nil
}

func (a *ClientManager) StartDeal(ctx context.Context, root cid.Cid, peerId string) (cid.Cid, error) {
	data, err := cid.Parse(root)
	if err != nil {
		return cid.Undef, err
	}

	ref := &storagemarket.DataRef{
		TransferType: storagemarket.TTGraphsync,
		Root:         data,
	}
	addr, _ := address.NewFromString("t01000")
	sdParams := &lapi.StartDealParams{
		Data:               ref,
		Wallet:             address.Undef,
		Miner:              addr,
		EpochPrice:         types.EmptyInt,
		MinBlocksDuration:  uint64(0),
		DealStartEpoch:     abi.ChainEpoch(0),
		FastRetrieval:      true,
		VerifiedDeal:       false,
		ProviderCollateral: bigext.Zero(),
	}
	var proposal cid.Cid
	proposal, err = a.dealStarter(ctx, sdParams, peerId)

	if err != nil {
		return cid.Undef, err
	}

	encoder, err := GetCidEncoder()
	if err != nil {
		return cid.Undef, err
	}

	lotusLog.Info(encoder.Encode(proposal))
	return proposal, nil
}

func (a *ClientManager) dealStarter(ctx context.Context, params *api.StartDealParams, peerId string) (cid.Cid, error) {

	bs, onDone, err := a.dealBlockstore(params.Data.Root)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to find blockstore for root CID: %w", err)
	}
	if has, err := bs.Has(params.Data.Root); err != nil {
		return cid.Undef, xerrors.Errorf("failed to query blockstore for root CID: %w", err)
	} else if !has {
		return cid.Undef, xerrors.Errorf("failed to find root CID in blockstore: %w", err)
	}
	onDone()

	mi, err := a.FullApi.StateMinerInfo(ctx, address.Undef, types.EmptyTSK)
	if err != nil {
		return cid.Undef, xerrors.Errorf("get sector size & type error: %w", err)
	}
	ver, err := a.FullApi.StateNetworkVersion(ctx, types.EmptyTSK)
	if err != nil {
		return cid.Undef, xerrors.Errorf("get network version error: %w", err)
	}

	sp, err := miner.PreferredSealProofTypeFromWindowPoStType(ver, mi.WindowPoStProofType)
	if err != nil {
		return cid.Undef, xerrors.Errorf("get PreferredSealProofTypeFromWindowPoStType error: %w", err)
	}

	//if uint64(params.Data.PieceSize.Padded()) > uint64(mi.SectorSize) {
	//	return nil, xerrors.New("data doesn't fit in a sector")
	//}
	//generate rand data to avoid deal.cid has the same value
	blocksPerHour := 60 * 60 / build.BlockDelaySecs
	dealStartEpoch := abi.ChainEpoch(dealStartBufferHours * blocksPerHour)
	num, _ := rand.Int(rand.Reader, big.NewInt(9999999999))
	dealEndEpoch := dealStartEpoch + abi.ChainEpoch(num.Int64())
	c, _ := cid.Decode(peerId)
	peerid, _ := peer.FromCid(c)
	lotusLog.Infof("peerid from string is %s", peerid)
	result, err := a.DealClient.ProposeStorageDeal(ctx, storagemarket.ProposeStorageDealParams{
		Addr: params.Miner,
		Info: &storagemarket.StorageProviderInfo{
			PeerID:     peerid,
			SectorSize: uint64(mi.SectorSize),
			Address:    params.Miner,
		},
		Data:          params.Data,
		StartEpoch:    dealStartEpoch,
		EndEpoch:      dealEndEpoch,
		Price:         params.EpochPrice,
		Collateral:    params.ProviderCollateral,
		Rt:            sp,
		FastRetrieval: params.FastRetrieval,
		VerifiedDeal:  params.VerifiedDeal,
	})
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to start deal: %w", err)
	}

	return result.ProposalCid, nil
}

func (a *ClientManager) SubscribeToAllTransferEvents(ctx context.Context) {
	unsub := a.DealClient.SubscribeToEvents(func(event storagemarket.ClientEvent, deal storagemarket.ClientDeal) {
		if event == storagemarket.ClientEventDataTransferComplete {
			localCid := peer.ToCid(deal.Miner)
			a.MarkHasStored(deal.DataRef.Root, localCid.String())
			a.FinalizeFile(deal.DataRef.Root, localCid.String())
		}
	})

	go func() {
		defer unsub()
		<-ctx.Done()
	}()
}

// dealBlockstore picks the source blockstore for a storage deal; either the
// IPFS blockstore, or an import CARv2 file. It also returns a function that
// must be called when done.
func (a *ClientManager) dealBlockstore(root cid.Cid) (bstore.Blockstore, func(), error) {
	switch acc := a.StorageBlockstoreAccessor.(type) {
	case *storageadapter.ImportsBlockstoreAccessor:
		bs, err := acc.Get(root)
		if err != nil {
			return nil, nil, xerrors.Errorf("%w", err)
		}

		doneFn := func() {
			_ = acc.Done(root) //nolint:errcheck
		}
		return bs, doneFn, nil

	case *storageadapter.ProxyBlockstoreAccessor:
		return acc.Blockstore, func() {}, nil

	default:
		return nil, nil, xerrors.Errorf("unsupported blockstore accessor type: %T", acc)
	}
}

func (a *ClientManager) makeRetrievalQuery(ctx context.Context, rp rm.RetrievalPeer, payload cid.Cid, piece *cid.Cid, qp rm.QueryParams) api.QueryOffer {
	queryResponse, err := a.Retrieval.Query(ctx, rp, payload, qp)
	if err != nil {
		return api.QueryOffer{Err: err.Error(), Miner: rp.Address, MinerPeer: rp}
	}
	var errStr string
	switch queryResponse.Status {
	case rm.QueryResponseAvailable:
		errStr = ""
	case rm.QueryResponseUnavailable:
		errStr = fmt.Sprintf("retrieval query offer was unavailable: %s", queryResponse.Message)
	case rm.QueryResponseError:
		errStr = fmt.Sprintf("retrieval query offer errored: %s", queryResponse.Message)
	}

	return api.QueryOffer{
		Root:        payload,
		Piece:       piece,
		Size:        queryResponse.Size,
		MinPrice:    queryResponse.PieceRetrievalPrice(),
		UnsealPrice: queryResponse.UnsealPrice,
		//PaymentInterval:         queryResponse.MaxPaymentInterval,
		//PaymentIntervalIncrease: queryResponse.MaxPaymentIntervalIncrease,
		PaymentInterval:         queryResponse.Size,
		PaymentIntervalIncrease: 0,
		Miner:                   queryResponse.PaymentAddress, // TODO: check
		MinerPeer:               rp,
		Err:                     errStr,
	}
}

func (a *ClientManager) ClientHasLocal(_ context.Context, root cid.Cid) (bool, error) {
	_, onDone, err := a.dealBlockstore(root)
	if err != nil {
		return false, err
	}
	onDone()
	return true, nil
}

//func (a *ClientManager) ClientFindDataNew(ctx context.Context, dataCid cid.Cid, piece *cid.Cid) error {
//
//	// Check if we already have this data locally
//	has, err := a.ClientHasLocal(ctx, dataCid)
//	if has {
//		lotusLog.Info("LOCAL")
//	}
//
//	offers, err := a.ClientFindData(ctx, dataCid, nil)
//	if err != nil {
//		return err
//	}
//
//	for _, offer := range offers {
//		if offer.Err != "" {
//			lotusLog.Errorf("ERR %s@%s: %s\n", offer.Miner, offer.MinerPeer.ID, offer.Err)
//			continue
//		}
//		lotusLog.Infof("RETRIEVAL %s@%s-%s-%s\n", offer.Miner, offer.MinerPeer.ID, "Free", types.SizeStr(types.NewInt(offer.Size)))
//	}
//
//	return nil
//}
//
//func (a *ClientManager) ClientFindData(ctx context.Context, root cid.Cid, piece *cid.Cid) ([]api.QueryOffer, error) {
//	peers, err := a.RetDiscovery.GetPeers(root)
//	if err != nil {
//		return nil, err
//	}
//
//	out := make([]api.QueryOffer, 0, len(peers))
//	for _, p := range peers {
//		if piece != nil && !piece.Equals(*p.PieceCID) {
//			continue
//		}
//
//		lotusLog.Infof("p.address:%s  p.id:%s", p.Address, p.ID)
//		pp := rm.RetrievalPeer{
//			Address: p.Address,
//			ID:      p.ID,
//		}
//		out = append(out, a.makeRetrievalQuery(ctx, pp, root, piece, rm.QueryParams{}))
//	}
//
//	return out, nil
//}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func removeIfNecessary(path string, dir string) {
	if filepath.IsAbs(path) {
		return
	}
	path = dir + "/" + path
	if common.FileExist(path) {
		os.Remove(path)
	}
}

// get Storage result.return boolean value.
func (a *ClientManager) GetStorageStatus(cid cid.Cid, peerId string) bool {
	return a.HasStored(cid, peerId)
}

func (a *ClientManager) MarkFileHashAndFlagToCid(oriHash string, headFlag bool, peerId string, cid cid.Cid) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(oriHash + strconv.FormatBool(headFlag) + peerId + "_FileHashAndFlagToCid")
	err := clientEx.Ds.Put(key, []byte(cid.String()))
	if err != nil {
		lotusLog.Errorf("MarkFileHashAndFlagToCid error:%w", err)
		return false
	}
	return true
}

func (a *ClientManager) GetFileHashAndFlagToCid(oriHash string, headFlag bool, peerId string) string {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(oriHash + strconv.FormatBool(headFlag) + peerId + "_FileHashAndFlagToCid")
	valbuf, err := clientEx.Ds.Get(key)
	if err == nil {
		strHash := string(valbuf)
		return strHash
	}
	return ""
}

func (a *ClientManager) MarkFileHashAndFlag(cid cid.Cid, peerId string, strHash string) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId + "_FileHashAndFlag")
	err := clientEx.Ds.Put(key, []byte(strHash))
	if err != nil {
		lotusLog.Errorf("MarkFileHashAndFlag error:%w", err)
		return false
	}
	return true
}

func (a *ClientManager) GetFileHashAndFlag(cid cid.Cid, peerId string) string {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId + "_FileHashAndFlag")
	valbuf, err := clientEx.Ds.Get(key)
	if err == nil {
		strHash := string(valbuf)
		clientEx.Ds.Delete(key)
		return strHash
	}
	return ""
}

func (a *ClientManager) MarkImportId(cid cid.Cid, peerId string, id imports.ID) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId + "_ImportId")
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	err := clientEx.Ds.Put(key, b)
	if err != nil {
		lotusLog.Errorf("MarkImportId error:%w", err)
		return false
	}
	return true
}

func (a *ClientManager) MarkFile(cid cid.Cid, peerId string, file string) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId + "_SFile")
	err := clientEx.Ds.Put(key, []byte(file))
	if err != nil {
		lotusLog.Errorf("MarkFile error:%w", err)
		return false
	}
	return true
}

//func (a *ClientManager) MarkStorageType(cid cid.Cid, peerId string, storageType string) bool {
//	clientEx := a.DealClient.(*clientext.ClientEx)
//	key := datastore.NewKey(cid.String() + peerId + "_StorageType")
//	err := clientEx.Ds.Put(key, []byte(storageType))
//	if err != nil {
//		lotusLog.Errorf("MarkFile error:%w", err)
//		return false
//	}
//	return true
//}

func (a *ClientManager) FinalizeFile(cid cid.Cid, peerId string) {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId + "_SFile")
	valbuf, err := clientEx.Ds.Get(key)
	if err != nil {
		lotusLog.Errorf("get key: %s, error:%w", key, err)
	} else {
		file := string(valbuf)
		removeIfNecessary(file, a.Repo.Path()+"/storage-file")
		clientEx.Ds.Delete(key)
	}

	imgr := a.importManager()
	carPath, err := imgr.CARPathFor(cid)
	if err == nil {
		os.Remove(carPath)
	}

	key = datastore.NewKey(cid.String() + peerId + "_ImportId")
	valbuf, err = clientEx.Ds.Get(key)
	if err != nil {
		lotusLog.Errorf("get ImportId by key:%s error:%w", key, err)
	} else {
		if len(valbuf) >= 8 {
			importId := imports.ID(binary.BigEndian.Uint64(valbuf))
			imgr.Remove(importId)
		}
		clientEx.Ds.Delete(key)
	}
}

func (a *ClientManager) HasStoring(cid cid.Cid, peerId string) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId)
	valbuf, err := clientEx.Ds.Get(key)
	if err != nil {
		return false
	}
	flg := string(valbuf)
	if flg == "storing" {
		return true
	}
	return false
}

func (a *ClientManager) MarkHasStoring(cid cid.Cid, peerId string) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId)
	err := clientEx.Ds.Put(key, []byte("storing"))
	if err != nil {
		lotusLog.Errorf("MarkHasStoring error:%w", err)
		return false
	}
	return true
}

func (a *ClientManager) HasStored(cid cid.Cid, peerId string) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId)
	valbuf, err := clientEx.Ds.Get(key)
	if err != nil {
		return false
	}
	flg := string(valbuf)
	if flg == "stored" {
		return true
	}
	return false
}

func (a *ClientManager) MarkHasStored(cid cid.Cid, peerId string) bool {
	clientEx := a.DealClient.(*clientext.ClientEx)
	key := datastore.NewKey(cid.String() + peerId)
	err := clientEx.Ds.Put(key, []byte("stored"))
	if err != nil {
		lotusLog.Errorf("MarkHasStored error:%w", err)
		return false
	}
	return true
}

// get Retrieval status.return boolean value.
func (a *ClientManager) GetRetrievalStatus(oriHash string, headFlag bool) RetrieveStatus {
	clientEx := a.DealClient.(*clientext.ClientEx)

	var tmpKey = ""
	if headFlag {
		tmpKey = oriHash + "_head"
	} else {
		tmpKey = oriHash + "_body"
	}

	key := datastore.NewKey(tmpKey)
	valbuf, err := clientEx.Ds.Get(key)
	if err != nil {
		return newRetrieveStatus("", DEAL_INIT, err.Error())
	}
	var rs RetrieveStatus
	err2 := json.Unmarshal(valbuf, &rs)
	if err2 != nil {
		return newRetrieveStatus("", DEAL_INIT, err2.Error())
	}
	// if status in db is sucess,then need to check the file if exists in disk!
	if rs.Status == DEAL_SUCCESS {
		cid := rs.Cid
		outPath := a.Repo.Path() + "/retrieval-file/" + cid
		fileExist := PathExists(outPath)
		if fileExist {
			return newRetrieveStatus(cid, DEAL_SUCCESS, "")
		}
		return newRetrieveStatus(cid, DEAL_INIT, "")
	}
	return rs
}

// get Retrieval status.return boolean value.
func (a *ClientManager) GetRetrievalStatusExt(filekey string, istxhash bool) RetrieveStatus {

	var rs RetrieveStatus
	lotusLog.Info("GetRetrievalStatusExt:\n", filekey, istxhash)
	if istxhash {
		clientEx := a.Retrieval.(*retrievalimpl2.ClientEx)
		res := clientEx.GetAuthStatus(filekey)
		rs.Status = res

		// if status in db is sucess,then need to check the file if exists in disk!
		if rs.Status == DEAL_SUCCESS || rs.Status == DEAL_INIT {
			outPath := a.Repo.Path() + "/retrieval-file/" + filekey
			fileExist := PathExists(outPath)
			if fileExist {
				rs.Cid = filekey
				rs.Status = DEAL_SUCCESS
				return rs
			}
			rs.Status = DEAL_ERROR //faild
		}

	} else {
		clientEx := a.DealClient.(*clientext.ClientEx)
		key := datastore.NewKey(filekey)
		valbuf, err := clientEx.Ds.Get(key)
		if err != nil {
			return newRetrieveStatus("", DEAL_INIT, err.Error())
		}

		err2 := json.Unmarshal(valbuf, &rs)
		if err2 != nil {
			return newRetrieveStatus("", DEAL_INIT, err2.Error())
		}
		// if status in db is sucess,then need to check the file if exists in disk!
		if rs.Status == DEAL_SUCCESS {
			cid := rs.Cid
			outPath := a.Repo.Path() + "/retrieval-file/" + cid
			fileExist := PathExists(outPath)
			if fileExist {
				return newRetrieveStatus(cid, DEAL_SUCCESS, "")
			}
			return newRetrieveStatus(cid, DEAL_INIT, "")
		}
	}

	return rs
}

// delete deal status in memory
func (a *ClientManager) MarkRetrieveResult(oriHash string, headFlag bool, cid cid.Cid, status int, err string) bool {
	rs := newRetrieveStatus(cid.String(), status, err)
	clientEx := a.DealClient.(*clientext.ClientEx)
	var tmpKey = ""
	if headFlag {
		tmpKey = oriHash + "_head"
	} else {
		tmpKey = oriHash + "_body"
	}
	key := datastore.NewKey(tmpKey)
	data, _ := json.Marshal(rs)
	err2 := clientEx.Ds.Put(key, data)
	if err2 != nil {
		lotusLog.Errorf("MarkHasRetrieved error:%w", err2)
		return false
	}
	return true
}

func (a *ClientManager) GetRetrievalOrder(cid cid.Cid, rp rm.RetrievalPeer) (*lapi.RetrievalOrder, error) {
	var order *lapi.RetrievalOrder
	var offer api.QueryOffer
	// addr, _ := address.NewFromString("f01000")
	rp.Address = address.Undef
	offer = api.QueryOffer{
		Root:                    cid,
		Piece:                   nil,
		Size:                    999999, // no need to set real size.
		MinPrice:                types.NewInt(0),
		UnsealPrice:             types.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
		Miner:                   address.Undef,
		MinerPeer:               rp,
		Err:                     "",
	}

	if offer.Err != "" {
		return nil, fmt.Errorf("The received offer errored: %s", offer.Err)
	}
	o := offer.Order(address.Undef)
	order = &o
	return order, nil
}

// retrieval file with oriHash
func (s *PublicStorageAPI) ClientRetrieveNew(oriHash string, headFlag bool, cid cid.Cid, out file_store.FileStoreStructFileMinerInfo) error {
	a := s.b.ClientManager()
	// mark as retrieving.
	a.MarkRetrieveResult(oriHash, headFlag, cid, DEAL_ING, "")
	// set out path
	outPath := a.Repo.Path() + "/retrieval-file/" + cid.String()
	lotusLog.Infof("Start Retrieval file, file will be stored to %s", outPath)
	ref := &lapi.FileRef{
		Path:  outPath,
		IsCAR: false,
	}

	// 1 hours as timeout
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Minute)

	updates, err := s.ClientRetrieveWithEvents(ctx, oriHash, out, ref)
	if err != nil {
		return xerrors.Errorf("error setting up retrieval: %w", err)
	}

	var logList []string
	go func() {
		var prevStatus rm.DealStatus
		for {
			select {
			case evt, ok := <-updates:
				if ok {
					if evt.Event == rm.ClientEventComplete {
						var str string
						if evt.Status == rm.DealStatusErrored {
							str = fmt.Sprintf("> Error : provider process error!")
						} else {
							str = fmt.Sprintf("> Recv Finish, %s (%s)",
								rm.ClientEvents[evt.Event],
								rm.DealStatuses[evt.Status])
						}
						logList = append(logList, str)
					} else if evt.Event == rm.ClientEventDataTransferError {
						lotusLog.Errorf("DataTransfer Error")
						a.MarkRetrieveResult(oriHash, headFlag, cid, DEAL_ERROR, "DataTransfer Error")
						return
					} else {
						str := fmt.Sprintf("> Recv: %s, %s (%s)",
							types.SizeStr(types.NewInt(evt.BytesReceived)),
							rm.ClientEvents[evt.Event],
							rm.DealStatuses[evt.Status])
						lotusLog.Infof("%s", str)
						logList = append(logList, str)
					}
					prevStatus = evt.Status
				}

				if evt.Err != "" {
					lotusLog.Errorf("> Error : %s", evt.Err)
					a.MarkRetrieveResult(oriHash, headFlag, cid, DEAL_ERROR, evt.Err)
					return
				}

				if !ok {
					if prevStatus == rm.DealStatusCompleted {
						for _, str := range logList {
							lotusLog.Infof(str)
						}
						lotusLog.Infof("> Retrieval Success")
						// mark retrieve success : delete the key.
						a.MarkRetrieveResult(oriHash, headFlag, cid, DEAL_SUCCESS, "")
						cancel()
					} else {
						lotusLog.Warnf("saw final deal state %s instead of expected success state DealStatusCompleted\n",
							rm.DealStatuses[prevStatus])
					}
					return
				}

			case <-ctx.Done():
				lotusLog.Errorf("retrieval timed out")
				a.MarkRetrieveResult(oriHash, headFlag, cid, DEAL_TIMEOUT, "retrieval fail:timed out")
				return
			}
		}
	}()
	return nil
}

func (s *PublicStorageAPI) ClientRetrieveWithEvents(ctx context.Context, oriHash string, out file_store.FileStoreStructFileMinerInfo, ref *lapi.FileRef) (<-chan marketevents.RetrievalEvent, error) {
	events := make(chan marketevents.RetrievalEvent)
	go s.clientRetrieve(ctx, events, oriHash, out, ref)
	return events, nil
}

func (s *PublicStorageAPI) clientRetrieve(ctx context.Context, events chan marketevents.RetrievalEvent, oriHash string, out file_store.FileStoreStructFileMinerInfo, ref *lapi.FileRef) {
	defer close(events)

	finish := func(e error) {
		if e != nil {
			events <- marketevents.RetrievalEvent{Err: e.Error(), FundsSpent: bigext.Zero()}
		}
	}

	var logPrefix = fmt.Sprintf("[oriHash=%s]", oriHash)
	minerIDs := out.MinerIds
	for i := 0; i < len(minerIDs); i++ {
		lotusLog.Infof("%s storage succeed miner id %s: %s", logPrefix, strconv.Itoa(i), common.Byte32ToHexStr(minerIDs[i]))
	}

	minerNum := len(minerIDs) - 1
	for i := 0; i < len(minerIDs); i++ {
		minerID := minerIDs[i]
		err := s.retrieve(ctx, events, out, minerID, *ref)
		if err != nil {
			errMsg := fmt.Sprintf("logPrefix:%s, minerID:%s fileHash:%s err:%s", logPrefix, common.Byte32ToHexStr(minerID), common.Byte32ToHexStr(out.FileHash), err.Error())
			if i == minerNum {
				lotusLog.Error(errMsg)
				finish(err)
				return
			}
			lotusLog.Warnf(errMsg)
			continue
		} else {
			lotusLog.Infof("retrieve miner info, logPrefix:%s minerID:%s minerSerialNumber:%s minerNum:%s", logPrefix, common.Byte32ToHexStr(minerID), i+1, len(minerIDs))
			break
		}
	}
}

func (s *PublicStorageAPI) retrieve(ctx context.Context, events chan marketevents.RetrievalEvent, out file_store.FileStoreStructFileMinerInfo, minerID [32]byte, ref lapi.FileRef) error {
	a := s.b.ClientManager()
	// get addr and peerId from contract (on chain).
	addrInfo, remotePeerId, err := s.getAddrInfoByMinerId(minerID)
	if err != nil {
		return xerrors.Errorf("cannot find miner's info in store contract")
	}
	peerCid, _ := cid.Decode(remotePeerId)
	ID, _ := peer.FromCid(peerCid)
	err = a.Host.Connect(ctx, addrInfo)
	if err != nil {
		return xerrors.Errorf("client connect to remote peer error: %s", err.Error())
	}
	dataCid, err := cid.Decode(out.FileCid)
	var order *lapi.RetrievalOrder
	rp := rm.RetrievalPeer{
		Address: address.Undef,
		ID:      ID,
	}
	lotusLog.Infow(ID.String())
	order, _ = a.GetRetrievalOrder(dataCid, rp)
	if err != nil {
		return xerrors.Errorf("cid not found, please check the cid value: %s", out.FileCid)
	}

	var id rm.DealID
	sel := shared.AllSelector()
	carBss, retrieveIntoCAR := a.RtvlBlockstoreAccessor.(*retrievaladapter.CARBlockstoreAccessor)
	carPath := order.FromLocalCAR

	if !retrieveIntoCAR {
		// we don't recognize the block store accessor.
		return xerrors.Errorf("unsupported retrieval block store accessor")
	}

	if order.MinerPeer == nil || order.MinerPeer.ID == "" {
		order.MinerPeer = &rm.RetrievalPeer{
			ID:      "todo",
			Address: order.Miner,
		}
	}

	if order.Total.Int == nil {
		return xerrors.Errorf("cannot make retrieval deal for null total")
	}

	if order.Size == 0 {
		return xerrors.Errorf("cannot make retrieval deal for zero bytes")
	}

	ppb := types.BigDiv(order.Total, types.NewInt(order.Size))
	order.PaymentInterval = 1 << 40
	order.PaymentIntervalIncrease = 0
	params, err := rm.NewParamsV1(ppb, order.PaymentInterval, order.PaymentIntervalIncrease, sel, order.Piece, order.UnsealPrice)
	if err != nil {
		return xerrors.Errorf("Error in retrieval params: %s", err)
	}

	// Subscribe to events before retrieving to avoid losing events.
	subscribeEvents := make(chan retrievalSubscribeEvent, 1)
	subscribeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	unsubscribe := a.Retrieval.SubscribeToEvents(func(event rm.ClientEvent, state rm.ClientDealState) {
		// We'll check the deal IDs inside consumeAllEvents.
		if state.PayloadCID.Equals(order.Root) {
			select {
			case <-subscribeCtx.Done():
			case subscribeEvents <- retrievalSubscribeEvent{event, state}:
			}
		}
	})

	id = a.Retrieval.NextID()
	lotusLog.Infof("address:%s cid:%s, ID:%s", order.Miner.String(), order.Root.String(), (*order.MinerPeer).ID.String())

	id, err = a.Retrieval.Retrieve(
		ctx,
		id,
		order.Root,
		params,
		order.Total,
		*order.MinerPeer,
		order.Client,
		order.Miner,
	)

	if err != nil {
		unsubscribe()
		return xerrors.Errorf("Retrieve failed: %w", err)
	}

	lotusLog.Infof("submit retrieval (id:%s) success!", id)
	lotusLog.Infof("consumeAllEvents start")
	err = consumeAllEvents(ctx, id, subscribeEvents, events)
	lotusLog.Infof("consumeAllEvents end")
	unsubscribe()
	if err != nil {
		return xerrors.Errorf("Retrieve: %w", err)
	}

	if retrieveIntoCAR {
		carPath = carBss.PathFor(id)
	}

	// determine where did the retrieval go
	var retrievalBs bstore.Blockstore
	cbs, err := stores.ReadOnlyFilestore(carPath)
	if err != nil {
		return err
	}
	defer cbs.Close() //nolint:err check
	retrievalBs = cbs

	// we are extracting a UnixFS file.
	ds := merkledag.NewDAGService(blockservice.New(retrievalBs, offline.Exchange(retrievalBs)))
	root := order.Root

	nd, err := ds.Get(ctx, root)
	if err != nil {
		return xerrors.Errorf("ClientRetrieve: %w", err)
	}
	file, err := unixfile.NewUnixfsFile(ctx, ds, nd)
	if err != nil {
		return xerrors.Errorf("ClientRetrieve: %w", err)
	}

	err = files.WriteTo(file, ref.Path)
	carFile := a.Repo.Path() + "/retrievals/" + id.String() + ".car"
	err = os.Remove(carFile)
	if err != nil {
		// delete fail
		lotusLog.Errorf("delete file %s error", carFile)
	}
	return nil
}

func consumeAllEvents(ctx context.Context, dealID rm.DealID, subscribeEvents chan retrievalSubscribeEvent, events chan marketevents.RetrievalEvent) error {
	for {
		var subscribeEvent retrievalSubscribeEvent
		select {
		case <-ctx.Done():
			return xerrors.New("Retrieval Timed Out")
		case subscribeEvent = <-subscribeEvents:
			if subscribeEvent.state.ID != dealID {
				// we can't check the deal ID ahead of time because:
				// 1. We need to subscribe before retrieving.
				// 2. We won't know the deal ID until after retrieving.
				continue
			}
		}

		select {
		case <-ctx.Done():
			return xerrors.New("Retrieval Timed Out")
		case events <- marketevents.RetrievalEvent{
			Event:         subscribeEvent.event,
			Status:        subscribeEvent.state.Status,
			BytesReceived: subscribeEvent.state.TotalReceived,
			FundsSpent:    subscribeEvent.state.FundsSpent,
		}:
		}

		state := subscribeEvent.state
		lotusLog.Infof("state's Status:%s", state.Status)
		switch state.Status {
		case rm.DealStatusCompleted:
			return nil
		case rm.DealStatusRejected:
			return xerrors.Errorf("Retrieval Proposal Rejected: %s", state.Message)
		case rm.DealStatusCancelled:
			return xerrors.Errorf("Retrieval was cancelled externally: %s", state.Message)
		case
			rm.DealStatusDealNotFound,
			rm.DealStatusErrored:
			return xerrors.Errorf("Retrieval Error: %s", state.Message)
		default:
			continue // must be continue,if not,will receive only one event
		}
	}
}

// GetCidEncoder returns an encoder using the `cid-base` flag if provided, or
// the default (Base32) encoder if not.
func GetCidEncoder() (cidenc.Encoder, error) {
	return cidenc.Encoder{Base: multibase.MustNewEncoder(multibase.Base32)}, nil
}

func (a *ClientManager) MarkeAutStatus(txhash string, status int) bool {
	clientEx := a.Retrieval.(*retrievalimpl2.ClientEx)
	return clientEx.MarkeAutStatus(txhash, status)
}

func (a *ClientManager) MarkRetrieveType(cid cid.Cid, headflag bool, auth bool, txhash string) bool {
	clientEx := a.Retrieval.(*retrievalimpl2.ClientEx)
	return clientEx.MarkRetrieveType(cid, headflag, auth, txhash)
}
