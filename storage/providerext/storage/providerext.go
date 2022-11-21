package storage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	voucher "github.com/ethereum/go-ethereum/borcontracts/checktxhash"
	file_store "github.com/ethereum/go-ethereum/borcontracts/file-store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/storage/event"
	"github.com/ethereum/go-ethereum/storage/requestvalidationex"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	storageimpl "github.com/filecoin-project/go-fil-markets/storagemarket/impl"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	logging "github.com/ipfs/go-log/v2"
	cbg "github.com/whyrusleeping/cbor-gen"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hannahhoward/go-pubsub"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-commp-utils/ffiwrapper"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	versioning "github.com/filecoin-project/go-ds-versioning/pkg"
	versionedfsm "github.com/filecoin-project/go-ds-versioning/pkg/fsm"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-padreader"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	acrypto "github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/exitcode"
	"github.com/filecoin-project/go-statemachine/fsm"

	"github.com/ethereum/go-ethereum/storage/auth"
	"github.com/ethereum/go-ethereum/storage/clientext"
	"github.com/filecoin-project/go-fil-markets/filestore"
	"github.com/filecoin-project/go-fil-markets/piecestore"
	"github.com/filecoin-project/go-fil-markets/shared"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/connmanager"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/dtutils"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/providerstates"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/providerutils"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/requestvalidation"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-fil-markets/stores"
	ipld "github.com/ipfs/go-ipld-format"
)

var log = logging.Logger("providerext")

// TODO: These are copied from spec-actors master, use spec-actors exports when we update
const DealMaxLabelSize = 256
const HEADORBODYSTATUS = 1

//var _ storagemarket.StorageProvider = &ProviderExt{}
//var _ network.StorageReceiver = &ProviderExt{}

// StoredAsk is an interface which provides access to a StorageAsk
type StoredAsk interface {
	GetAsk() *storagemarket.SignedStorageAsk
	SetAsk(price abi.TokenAmount, verifiedPrice abi.TokenAmount, duration abi.ChainEpoch, options ...storagemarket.StorageAskOption) error
}

// Provider is the production implementation of the StorageProvider interface
type ProviderExt struct {
	ctx                   context.Context
	environment           providerExtDealEnvironment
	net                   network.StorageMarketNetwork
	sa                    retrievalmarket.SectorAccessor
	spn                   storagemarket.StorageProviderNode
	fs                    filestore.FileStore
	pieceStore            piecestore.PieceStore
	conns                 *connmanager.ConnManager
	storedAsk             StoredAsk
	actor                 address.Address
	dataTransfer          datatransfer.Manager
	customDealDeciderFunc DealDeciderFunc
	pubSub                *pubsub.PubSub
	readyMgr              *shared.ReadyManager
	receiver              event.DataTransferEventReceiver
	backend               ethapi.Backend

	migrateDeals func(context.Context) error

	unsubDataTransfer datatransfer.Unsubscribe

	dagStore stores.DAGStoreWrapper
	stores   *stores.ReadWriteBlockstores
}

// StorageProviderOption allows custom configuration of a storage providerext
type StorageProviderOption func(p *ProviderExt)

// DealDeciderFunc is a function which evaluates an incoming deal to decide if
// it its accepted
// It returns:
// - boolean = true if deal accepted, false if rejected
// - string = reason deal was not excepted, if rejected
// - error = if an error occurred trying to decide
type DealDeciderFunc func(context.Context, storagemarket.MinerDeal) (bool, string, error)

// CustomDealDecisionLogic allows a providerext to call custom decision logic when validating incoming
// deal proposals
func CustomDealDecisionLogic(decider DealDeciderFunc) StorageProviderOption {
	return func(p *ProviderExt) {
		p.customDealDeciderFunc = decider
	}
}

// NewProvider returns a new storage providerext
func NewProviderExt(net network.StorageMarketNetwork,
	ds datastore.Batching,
	fs filestore.FileStore,
	sa retrievalmarket.SectorAccessor,
	dagStore stores.DAGStoreWrapper,
	pieceStore piecestore.PieceStore,
	dataTransfer datatransfer.Manager,
	spn storagemarket.StorageProviderNode,
	minerAddress address.Address,
	storedAsk storageimpl.StoredAsk,
	backend ethapi.Backend,
	options ...StorageProviderOption,
) (storagemarket.StorageProvider, error) {
	h := &ProviderExt{
		net:          net,
		spn:          spn,
		fs:           fs,
		pieceStore:   pieceStore,
		conns:        connmanager.NewConnManager(),
		storedAsk:    storedAsk,
		actor:        minerAddress,
		dataTransfer: dataTransfer,
		pubSub:       pubsub.New(providerDispatcher),
		readyMgr:     shared.NewReadyManager(),
		receiver:     event.NewDataTransferEventReceiver(),
		dagStore:     dagStore,
		stores:       stores.NewReadWriteBlockstores(),
		backend:      backend,
		sa:           sa,
	}
	h.Configure(options...)

	// register a data transfer event handler -- this will send events to the state machines based on DT events
	h.unsubDataTransfer = dataTransfer.SubscribeToEvents(event.DataTransferSubscriber(&h.receiver))

	err := dataTransfer.RegisterVoucherType(&requestvalidation.StorageDataTransferVoucher{}, requestvalidationex.NewUnifiedRequestValidatorEx(&providerPushDeals{h}, nil))
	if err != nil {
		return nil, err
	}

	err = dataTransfer.RegisterTransportConfigurer(&requestvalidation.StorageDataTransferVoucher{}, dtutils.TransportConfigurer(&providerStoreGetter{h}))
	if err != nil {
		return nil, err
	}

	return h, nil
}

// Start initializes deal processing on a StorageProvider and restarts in progress deals.
// It also registers the providerext with a StorageMarketNetwork so it can receive incoming
// messages on the storage market's libp2p protocols
func (p *ProviderExt) Start(ctx context.Context) error {
	p.ctx = ctx
	p.environment = providerExtDealEnvironment{p}
	err := p.net.SetDelegate(p)
	if err != nil {
		return err
	}
	go func() {
		err := p.start(ctx)
		if err != nil {
			log.Error(err.Error())
		}
	}()
	return nil
}

// OnReady registers a listener for when the providerext has finished starting up
func (p *ProviderExt) OnReady(ready shared.ReadyFunc) {
	p.readyMgr.OnReady(ready)
}

func (p *ProviderExt) AwaitReady() error {
	return p.readyMgr.AwaitReady()
}

/*
HandleDealStream is called by the network implementation whenever a new message is received on the deal protocol

It initiates the providerext side of the deal flow.

When a providerext receives a DealProposal of the deal protocol, it takes the following steps:

1. Calculates the CID for the received ClientDealProposal.

2. Constructs a MinerDeal to track the state of this deal.

3. Tells its statemachine to begin tracking this deal state by CID of the received ClientDealProposal

4. Tracks the received deal stream by the CID of the ClientDealProposal

4. Triggers a `ProviderEventOpen` event on its statemachine.

From then on, the statemachine controls the deal flow in the client. Other components may listen for events in this flow by calling
`SubscribeToEvents` on the Provider. The Provider handles loading the next block to send to the client.

Documentation of the client state machine can be found at https://godoc.org/github.com/filecoin-project/go-fil-markets/storagemarket/impl/providerstates
*/
func (p *ProviderExt) HandleDealStream(s network.StorageDealStream) {
	log.Info("Handling storage deal proposal!")

	err := p.receiveDeal(s)
	if err != nil {
		log.Errorf("%+v", err)
		s.Close()
		return
	}
}

func (p *ProviderExt) receiveDeal(s network.StorageDealStream) error {
	proposal, err := s.ReadDealProposal()
	if err != nil {
		return xerrors.Errorf("failed to read proposal message: %w", err)
	}

	ishead, err := p.checkauthority(s, proposal)
	if ishead {
		return err
	}

	proposalNd, err := cborutil.AsIpld(proposal.DealProposal)
	if err != nil {
		return err
	}

	// Check if we are already tracking this deal
	var md *storagemarket.MinerDeal
	a, ok := p.receiver.Get(proposalNd.Cid())
	if ok {
		md = a.(*storagemarket.MinerDeal)
		// We are already tracking this deal, for some reason it was re-proposed, perhaps because of a client restart
		// this is ok, just send a response back.
		return p.resendProposalResponse(s, md)
	}

	var path string
	// create an empty CARv2 file at a temp location that Graphysnc will write the incoming blocks to via a CARv2 ReadWrite blockstore wrapper.
	if proposal.Piece.TransferType != storagemarket.TTManual {
		tmp, err := p.fs.CreateTemp()
		if err != nil {
			return xerrors.Errorf("failed to create an empty temp CARv2 file: %w", err)
		}
		if err := tmp.Close(); err != nil {
			_ = os.Remove(string(tmp.OsPath()))
			return xerrors.Errorf("failed to close temp file: %w", err)
		}
		path = string(tmp.OsPath())
	}
	//p.lk.Lock()
	//dealid := abi.DealID(1)
	//dealbuf, err := p.ds.Get(datastore.NewKey("miner-deal-id"))
	//if err != nil {
	//	b, _ := binary.Varint(dealbuf)
	//	dealid = abi.DealID(b)
	//}
	//dealid = dealid + 1
	//
	//dealidbuf := make([]byte, binary.MaxVarintLen64)
	//n := binary.PutVarint(dealidbuf, int64(dealid))
	//if n > 0 {
	//	p.ds.Put(datastore.NewKey("miner-deal-id"), dealidbuf)
	//}
	//p.lk.Unlock()

	c := proposalNd.Cid()
	deal := &storagemarket.MinerDeal{
		Client:             s.RemotePeer(),
		Miner:              p.net.ID(),
		ClientDealProposal: *proposal.DealProposal,
		ProposalCid:        proposalNd.Cid(),
		PublishCid:         &c,
		State:              storagemarket.StorageDealUnknown,
		Ref:                proposal.Piece,
		FastRetrieval:      proposal.FastRetrieval,
		CreationTime:       curTime(),
		InboundCAR:         path,
	}

	err = p.receiver.Begin(proposalNd.Cid(), deal)
	if err != nil {
		return err
	}
	defer p.receiver.End(proposalNd.Cid(), deal)
	err = p.conns.AddStream(proposalNd.Cid(), s)
	if err != nil {
		return err
	}
	//send response
	p.environment.TagPeer(deal.Client, deal.ProposalCid.String())
	if len(proposal.DealProposal.Proposal.Label) > DealMaxLabelSize {
		return xerrors.Errorf("deal label can be at most %d bytes, is %d", DealMaxLabelSize, len(proposal.DealProposal.Proposal.Label))
	}

	if err := proposal.DealProposal.Proposal.PieceSize.Validate(); err != nil {
		return xerrors.Errorf("proposal piece size is invalid: %w", err)
	}

	if !proposal.DealProposal.Proposal.PieceCID.Defined() {
		return xerrors.Errorf("proposal PieceCID undefined")
	}

	if proposal.DealProposal.Proposal.PieceCID.Prefix() != market.PieceCIDPrefix {
		return xerrors.Errorf("proposal PieceCID had wrong prefix")
	}
	// Send intent to accept
	if err := p.environment.SendSignedResponse(nil, &network.Response{
		State:    storagemarket.StorageDealWaitingForData,
		Proposal: deal.ProposalCid,
	}); err != nil {
		return err
	}

	if err := p.environment.Disconnect(deal.ProposalCid); err != nil {
		log.Warnf("closing client connection: %+v", err)
	}

	if err := p.receiver.Wait(p.ctx, deal.ProposalCid); err != nil {
		log.Warnf("event receiver wait: %+v", err)
		return err
	}

	//VerifyData
	if err := p.environment.FinalizeBlockstore(deal.ProposalCid); err != nil {
		return xerrors.Errorf("failed to finalize read/write blockstore: %w", err)
	}

	pieceCid, _, err := p.environment.GeneratePieceCommitment(deal.ProposalCid, deal.InboundCAR, deal.Proposal.PieceSize)
	if err != nil {
		return xerrors.Errorf("error generating CommP: %w", err)
	}

	// Verify CommP matches
	if pieceCid != deal.Proposal.PieceCID {
		return xerrors.Errorf("proposal CommP doesn't match calculated CommP")
	}

	log.Infof("HandoffDeal start %+v", *deal)
	err = p.HandoffDeal(*deal)
	if err != nil {
		log.Errorf("HandoffDeal: %+v", err)
		return err
	}
	log.Info("HandoffDeal End")
	// transfer completed, update the smart contract
	txHash, err := p.updateFileStoreInfoAll(proposal)
	if err != nil {
		log.Errorf("updateFileStoreInfoAll: %+v", err)
		return err
	}
	log.Infof("updateFileStoreInfoAll %s", txHash)
	return nil
}

// HandoffDeal hands off a published deal for sealing and commitment in a sector
func (p *ProviderExt) HandoffDeal(deal storagemarket.MinerDeal) error {
	var packingInfo *storagemarket.PackingResult
	var carFilePath string
	if deal.PiecePath != "" {
		// Data for offline deals is stored on disk, so if PiecePath is set,
		// create a Reader from the file path
		file, err := p.environment.FileStore().Open(deal.PiecePath)
		if err != nil {
			return xerrors.Errorf("reading piece at path %s: %w", deal.PiecePath, err)
		}
		carFilePath = string(file.OsPath())

		// Hand the deal off to the process that adds it to a sector
		log.Infow("handing off deal to sealing subsystem", "pieceCid", deal.Proposal.PieceCID, "proposalCid", deal.ProposalCid)
		packingInfo, err = handoffDeal(p.ctx, p.environment, deal, file, uint64(file.Size()))
		if err := file.Close(); err != nil {
			log.Errorw("failed to close imported CAR file", "pieceCid", deal.Proposal.PieceCID, "proposalCid", deal.ProposalCid, "err", err)
		}

		if err != nil {
			err = xerrors.Errorf("packing piece at path %s: %w", deal.PiecePath, err)
			return err
		}
	} else {
		carFilePath = deal.InboundCAR

		v2r, err := p.environment.ReadCAR(deal.InboundCAR)
		if err != nil {
			return xerrors.Errorf("failed to open CARv2 file, proposalCid=%s: %w",
				deal.ProposalCid, err)
		}

		// Hand the deal off to the process that adds it to a sector
		var packingErr error
		log.Infow("handing off deal to sealing subsystem", "pieceCid", deal.Proposal.PieceCID, "proposalCid", deal.ProposalCid)
		packingInfo, packingErr = handoffDeal(p.ctx, p.environment, deal, v2r.DataReader(), v2r.Header.DataSize)
		// Close the reader as we're done reading from it.
		if err := v2r.Close(); err != nil {
			return xerrors.Errorf("failed to close CARv2 reader: %w", err)
		}
		log.Infow("closed car datareader after handing off deal to sealing subsystem", "pieceCid", deal.Proposal.PieceCID, "proposalCid", deal.ProposalCid)
		if packingErr != nil {
			err = xerrors.Errorf("packing piece %s: %w", deal.Ref.PieceCid, packingErr)
			return err
		}
	}

	if err := recordPiece(p.environment, deal, packingInfo.SectorNumber, packingInfo.Offset, packingInfo.Size); err != nil {
		err = xerrors.Errorf("failed to register deal data for piece %s for retrieval: %w", deal.Ref.PieceCid, err)
		log.Error(err.Error())
		_ = err
	}

	// Register the deal data as a "shard" with the DAG store. Later it can be
	// fetched from the DAG store during retrieval.
	if err := p.environment.RegisterShard(p.ctx, deal.Proposal.PieceCID, carFilePath, true); err != nil {
		err = xerrors.Errorf("failed to activate shard: %w", err)
		log.Error(err)
	}

	log.Infow("successfully handed off deal to sealing subsystem", "pieceCid", deal.Proposal.PieceCID, "proposalCid", deal.ProposalCid)
	return nil
}

func handoffDeal(ctx context.Context, environment providerExtDealEnvironment, deal storagemarket.MinerDeal, reader io.Reader, payloadSize uint64) (*storagemarket.PackingResult, error) {
	// because we use the PadReader directly during AP we need to produce the
	// correct amount of zeroes
	// (alternative would be to keep precise track of sector offsets for each
	// piece which is just too much work for a seldom used feature)
	paddedReader, err := padreader.NewInflator(reader, payloadSize, deal.Proposal.PieceSize.Unpadded())
	if err != nil {
		return nil, err
	}

	return environment.Node().OnDealComplete(
		ctx,
		storagemarket.MinerDeal{
			Client:             deal.Client,
			ClientDealProposal: deal.ClientDealProposal,
			ProposalCid:        deal.ProposalCid,
			State:              deal.State,
			Ref:                deal.Ref,
			PublishCid:         deal.PublishCid,
			DealID:             deal.DealID,
			FastRetrieval:      deal.FastRetrieval,
		},
		deal.Proposal.PieceSize.Unpadded(),
		paddedReader,
	)
}

func recordPiece(environment providerExtDealEnvironment, deal storagemarket.MinerDeal, sectorID abi.SectorNumber, offset, length abi.PaddedPieceSize) error {

	var blockLocations map[cid.Cid]piecestore.BlockLocation
	if deal.MetadataPath != filestore.Path("") {
		var err error
		blockLocations, err = providerutils.LoadBlockLocations(environment.FileStore(), deal.MetadataPath)
		if err != nil {
			return xerrors.Errorf("failed to load block locations: %w", err)
		}
	} else {
		blockLocations = map[cid.Cid]piecestore.BlockLocation{
			deal.Ref.Root: {},
		}
	}

	if err := environment.PieceStore().AddPieceBlockLocations(deal.Proposal.PieceCID, blockLocations); err != nil {
		return xerrors.Errorf("failed to add piece block locations: %s", err)
	}

	err := environment.PieceStore().AddDealForPiece(deal.Proposal.PieceCID, piecestore.DealInfo{
		DealID:   deal.DealID,
		SectorID: sectorID,
		Offset:   offset,
		Length:   length,
	})
	if err != nil {
		return xerrors.Errorf("failed to add deal for piece: %s", err)
	}

	return nil
}

// fetchKeystore retrieves the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) (*keystore.KeyStore, error) {
	if ks := am.Backends(keystore.KeyStoreType); len(ks) > 0 {
		return ks[0].(*keystore.KeyStore), nil
	}
	return nil, errors.New("local keystore not used")
}

// include entire file storage and seperate file storage.
func (p *ProviderExt) updateFileStoreInfoAll(proposal network.Proposal) (common.Hash, error) {
	dataBuf := proposal.DealProposal.ClientSignature.Data[32:]
	strHash := string(dataBuf)
	strHashArr := strings.Split(strHash, ",")
	hashStr := strHashArr[0] //oriHash or storeKey
	storageType := ""
	if len(strHashArr) >= 3 {
		// get storage type
		storageType = strHashArr[2]
		if storageType == ethapi.ENTIRE_FILE {
			// compute sha256
			hashStr = common.BuildNewOriHashWithSha256(hashStr)
		}
	}
	hashByte := common.HexSTrToByte32(hashStr)
	fileFlag, err := strconv.ParseBool(strHashArr[1])
	if err != nil {
		return common.Hash{}, err
	}
	ethBackend := p.backend
	addr, err := ethBackend.Coinbase()
	ks, err := fetchKeystore(ethBackend.AccountManager())
	if err != nil {
		return common.Hash{}, err
	}
	privateKey, err := ks.GetAccountPrivateKeyWithoutPass(accounts.Account{Address: addr})
	if err != nil {
		return common.Hash{}, err
	}
	var data []byte
	if storageType == ethapi.ENTIRE_FILE {
		data, err = file_store.FileStoreCli.Pack4UpdateFileStoreInfo4Entire(hashByte, proposal.Piece.Root.String(), HEADORBODYSTATUS)
	} else {
		data, err = file_store.FileStoreCli.Pack4UpdateFileStoreInfo(hashByte, fileFlag, proposal.Piece.Root.String(), HEADORBODYSTATUS)
	}

	if err != nil {
		return common.Hash{}, errors.New("pack data is err")
	}

	var tx *types.Transaction
	for i, v := range []int{2, 4, 6, 10, 15} {
		// build auth object
		auth, err := file_store.FileStoreCli.GenerateAuthObj(privateKey, ethBackend.ChainConfig().ChainID, addr, data)
		if err != nil {
			continue
		}
		if storageType == ethapi.ENTIRE_FILE {
			log.Infow("UpdateEntireFileStoreInfo params", "storeKey: ", hashStr, "cid: ", proposal.Piece.Root.String(), "nonce:", auth.Nonce, "gasLimit:", auth.GasLimit)
			tx, err = file_store.FileStoreCli.UpdateFileStoreInfo4Entire(auth, hashByte, proposal.Piece.Root.String(), HEADORBODYSTATUS)
			if err == nil {
				return tx.Hash(), nil
			}
			log.Errorw("UpdateEntireFileStoreInfo error:", "storeKey: ", hashStr, " root: ", proposal.Piece.Root.String(), " error: ", err, "retry times: ", i)
		} else {
			log.Infow("updateFileStoreInfo params", "orihash: ", hashStr, "HeadFlag: ", fileFlag, "cid: ", proposal.Piece.Root.String(), "nonce:", auth.Nonce, "gasLimit:", auth.GasLimit)
			tx, err = file_store.FileStoreCli.UpdateFileStoreInfo(auth, hashByte, fileFlag, proposal.Piece.Root.String(), HEADORBODYSTATUS)
			if err == nil {
				return tx.Hash(), nil
			}
			log.Errorw("updateFileStoreInfo error:", "orihash: ", hashStr, "HeadFlag: ", fileFlag, " root: ", proposal.Piece.Root.String(), " error: ", err, "retry times: ", i)
		}

		if strings.Contains(err.Error(), "replacement transaction underpriced") && i < 2 {
			time.Sleep(time.Duration(v) * time.Second)
			continue
		}
		return common.Hash{}, err
	}
	return common.Hash{}, err
}

// Stop terminates processing of deals on a StorageProvider
func (p *ProviderExt) Stop() error {
	p.readyMgr.Stop()
	p.unsubDataTransfer()
	err := p.receiver.Stop(context.TODO())
	if err != nil {
		return err
	}
	return p.net.StopHandlingRequests()
}

// ImportDataForDeal manually imports data for an offline storage deal
// It will verify that the data in the passed io.Reader matches the expected piece
// cid for the given deal or it will error
func (p *ProviderExt) ImportDataForDeal(ctx context.Context, propCid cid.Cid, data io.Reader) error {
	// TODO: be able to check if we have enough disk space
	var d storagemarket.MinerDeal
	a, ok := p.receiver.Get(propCid)
	d = a.(storagemarket.MinerDeal)
	if !ok {
		return xerrors.Errorf("failed getting deal %s: %w", propCid)
	}

	tempfi, err := p.fs.CreateTemp()
	if err != nil {
		return xerrors.Errorf("failed to create temp file for data import: %w", err)
	}
	defer tempfi.Close()
	cleanup := func() {
		_ = tempfi.Close()
		_ = p.fs.Delete(tempfi.Path())
	}

	log.Debugw("will copy imported file to local file", "propCid", propCid)
	n, err := io.Copy(tempfi, data)
	if err != nil {
		cleanup()
		return xerrors.Errorf("importing deal data failed: %w", err)
	}
	log.Debugw("finished copying imported file to local file", "propCid", propCid)

	_ = n // TODO: verify n?

	carSize := uint64(tempfi.Size())

	_, err = tempfi.Seek(0, io.SeekStart)
	if err != nil {
		cleanup()
		return xerrors.Errorf("failed to seek through temp imported file: %w", err)
	}

	proofType, err := p.spn.GetProofType(ctx, p.actor, nil)
	if err != nil {
		cleanup()
		return xerrors.Errorf("failed to determine proof type: %w", err)
	}
	log.Debugw("fetched proof type", "propCid", propCid)

	pieceCid, err := generatePieceCommitment(proofType, tempfi, carSize)
	if err != nil {
		cleanup()
		return xerrors.Errorf("failed to generate commP: %w", err)
	}
	log.Debugw("generated pieceCid for imported file", "propCid", propCid)

	if carSizePadded := padreader.PaddedSize(carSize).Padded(); carSizePadded < d.Proposal.PieceSize {
		// need to pad up!
		rawPaddedCommp, err := commp.PadCommP(
			// we know how long a pieceCid "hash" is, just blindly extract the trailing 32 bytes
			pieceCid.Hash()[len(pieceCid.Hash())-32:],
			uint64(carSizePadded),
			uint64(d.Proposal.PieceSize),
		)
		if err != nil {
			cleanup()
			return err
		}
		pieceCid, _ = commcid.DataCommitmentV1ToCID(rawPaddedCommp)
	}

	// Verify CommP matches
	if !pieceCid.Equals(d.Proposal.PieceCID) {
		cleanup()
		return xerrors.Errorf("given data does not match expected commP (got: %s, expected %s)", pieceCid, d.Proposal.PieceCID)
	}

	log.Debugw("will fire ProviderEventVerifiedData for imported file", "propCid", propCid)

	return p.receiver.VerifiedData(propCid, tempfi.Path(), filestore.Path(""))
}

func generatePieceCommitment(rt abi.RegisteredSealProof, rd io.Reader, pieceSize uint64) (cid.Cid, error) {
	paddedReader, paddedSize := padreader.New(rd, pieceSize)
	commitment, err := ffiwrapper.GeneratePieceCIDFromFile(rt, paddedReader, paddedSize)
	if err != nil {
		return cid.Undef, err
	}
	return commitment, nil
}

// GetAsk returns the storage miner's ask, or nil if one does not exist.
func (p *ProviderExt) GetAsk() *storagemarket.SignedStorageAsk {
	return p.storedAsk.GetAsk()
}

// AddStorageCollateral adds storage collateral
func (p *ProviderExt) AddStorageCollateral(ctx context.Context, amount abi.TokenAmount) error {
	done := make(chan error, 1)

	mcid, err := p.spn.AddFunds(ctx, p.actor, amount)
	if err != nil {
		return err
	}

	err = p.spn.WaitForMessage(ctx, mcid, func(code exitcode.ExitCode, bytes []byte, finalCid cid.Cid, err error) error {
		if err != nil {
			done <- xerrors.Errorf("AddFunds errored: %w", err)
		} else if code != exitcode.Ok {
			done <- xerrors.Errorf("AddFunds error, exit code: %s", code.String())
		} else {
			done <- nil
		}
		return nil
	})

	if err != nil {
		return err
	}

	return <-done
}

// GetStorageCollateral returns the current collateral balance
func (p *ProviderExt) GetStorageCollateral(ctx context.Context) (storagemarket.Balance, error) {
	tok, _, err := p.spn.GetChainHead(ctx)
	if err != nil {
		return storagemarket.Balance{}, err
	}

	return p.spn.GetBalance(ctx, p.actor, tok)
}

func (p *ProviderExt) RetryDealPublishing(propcid cid.Cid) error {
	return p.receiver.Restart(propcid, storagemarket.ProviderEventRestart)
}

// ListLocalDeals lists deals processed by this storage providerext
func (p *ProviderExt) ListLocalDeals() ([]storagemarket.MinerDeal, error) {
	var out []storagemarket.MinerDeal
	if err := p.receiver.List(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetAsk configures the storage miner's ask with the provided price,
// duration, and options. Any previously-existing ask is replaced.
func (p *ProviderExt) SetAsk(price abi.TokenAmount, verifiedPrice abi.TokenAmount, duration abi.ChainEpoch, options ...storagemarket.StorageAskOption) error {
	return p.storedAsk.SetAsk(price, verifiedPrice, duration, options...)
}

/*
HandleAskStream is called by the network implementation whenever a new message is received on the ask protocol

A Provider handling a `AskRequest` does the following:

1. Reads the current signed storage ask from storage

2. Wraps the signed ask in an AskResponse and writes it on the StorageAskStream

The connection is kept open only as long as the request-response exchange.
*/
func (p *ProviderExt) HandleAskStream(s network.StorageAskStream) {
	defer s.Close()
	ar, err := s.ReadAskRequest()
	if err != nil {
		log.Errorf("failed to read AskRequest from incoming stream: %s", err)
		return
	}

	var ask *storagemarket.SignedStorageAsk
	if p.actor != ar.Miner {
		log.Warnf("storage providerext for address %s receive ask for miner with address %s", p.actor, ar.Miner)
	} else {
		ask = p.storedAsk.GetAsk()
	}

	resp := network.AskResponse{
		Ask: ask,
	}

	if err := s.WriteAskResponse(resp, p.sign); err != nil {
		log.Errorf("failed to write ask response: %s", err)
		return
	}
}

/*
HandleDealStatusStream is called by the network implementation whenever a new message is received on the deal status protocol

A Provider handling a `DealStatuRequest` does the following:

1. Lots the deal state from the Provider FSM

2. Verifies the signature on the DealStatusRequest matches the Client for this deal

3. Constructs a ProviderDealState from the deal state

4. Signs the ProviderDealState with its private key

5. Writes a DealStatusResponse with the ProviderDealState and signature onto the DealStatusStream

The connection is kept open only as long as the request-response exchange.
*/
func (p *ProviderExt) HandleDealStatusStream(s network.DealStatusStream) {
	ctx := context.TODO()
	defer s.Close()
	request, err := s.ReadDealStatusRequest()
	if err != nil {
		log.Errorf("failed to read DealStatusRequest from incoming stream: %s", err)
		return
	}

	dealState, err := p.processDealStatusRequest(ctx, &request)
	if err != nil {
		log.Errorf("failed to process deal status request: %s", err)
		dealState = &storagemarket.ProviderDealState{
			State:   storagemarket.StorageDealError,
			Message: err.Error(),
		}
	}

	signature, err := p.sign(ctx, dealState)
	if err != nil {
		log.Errorf("failed to sign deal status response: %s", err)
		return
	}

	response := network.DealStatusResponse{
		DealState: *dealState,
		Signature: *signature,
	}

	if err := s.WriteDealStatusResponse(response, p.sign); err != nil {
		log.Warnf("failed to write deal status response: %s", err)
		return
	}
}

func (p *ProviderExt) processDealStatusRequest(ctx context.Context, request *network.DealStatusRequest) (*storagemarket.ProviderDealState, error) {
	// fetch deal state
	var md = storagemarket.MinerDeal{}
	m, ok := p.receiver.Get(request.Proposal)
	md = m.(storagemarket.MinerDeal)
	if !ok {
		log.Errorf("proposal doesn't exist in state store: %s", ok)
		return nil, xerrors.Errorf("no such proposal")
	}

	// verify query signature
	buf, err := cborutil.Dump(&request.Proposal)
	if err != nil {
		log.Errorf("failed to serialize status request: %s", err)
		return nil, xerrors.Errorf("internal error")
	}

	tok, _, err := p.spn.GetChainHead(ctx)
	if err != nil {
		log.Errorf("failed to get chain head: %s", err)
		return nil, xerrors.Errorf("internal error")
	}

	err = providerutils.VerifySignature(ctx, request.Signature, md.ClientDealProposal.Proposal.Client, buf, tok, p.spn.VerifySignature)
	if err != nil {
		log.Errorf("invalid deal status request signature: %s", err)
		return nil, xerrors.Errorf("internal error")
	}

	return &storagemarket.ProviderDealState{
		State:         md.State,
		Message:       md.Message,
		Proposal:      &md.Proposal,
		ProposalCid:   &md.ProposalCid,
		AddFundsCid:   md.AddFundsCid,
		PublishCid:    md.PublishCid,
		DealID:        md.DealID,
		FastRetrieval: md.FastRetrieval,
	}, nil
}

// Configure applies the given list of StorageProviderOptions after a StorageProvider
// is initialized
func (p *ProviderExt) Configure(options ...StorageProviderOption) {
	for _, option := range options {
		option(p)
	}
}

// SubscribeToEvents allows another component to listen for events on the StorageProvider
// in order to track deals as they progress through the deal flow
func (p *ProviderExt) SubscribeToEvents(subscriber storagemarket.ProviderSubscriber) shared.Unsubscribe {
	return shared.Unsubscribe(p.pubSub.Subscribe(subscriber))
}

// dispatch puts the fsm event into a form that pubSub can consume,
// then publishes the event
func (p *ProviderExt) dispatch(eventName fsm.EventName, deal fsm.StateType) {
	evt, ok := eventName.(storagemarket.ProviderEvent)
	if !ok {
		log.Errorf("dropped bad event %s", eventName)
	}
	realDeal, ok := deal.(storagemarket.MinerDeal)
	if !ok {
		log.Errorf("not a MinerDeal %v", deal)
	}
	pubSubEvt := internalProviderEvent{evt, realDeal}

	log.Debugw("process storage providerext listeners", "name", storagemarket.ProviderEvents[evt], "proposal cid", realDeal.ProposalCid)
	if err := p.pubSub.Publish(pubSubEvt); err != nil {
		log.Errorf("failed to publish event %d", evt)
	}
}

func (p *ProviderExt) start(ctx context.Context) error {
	// Run datastore and DAG store migrations
	deals, err := p.runMigrations(ctx)
	publishErr := p.readyMgr.FireReady(err)
	if publishErr != nil {
		log.Warnf("publish storage providerext ready event: %s", err.Error())
	}
	if err != nil {
		return err
	}

	// Fire restart event on all active deals
	if err := p.restartDeals(deals); err != nil {
		return fmt.Errorf("failed to restart deals: %w", err)
	}
	return nil
}

func (p *ProviderExt) runMigrations(ctx context.Context) ([]storagemarket.MinerDeal, error) {
	// Perform datastore migration
	//err := p.migrateDeals(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("migrating storage providerext state machines: %w", err)
	//}

	var deals []storagemarket.MinerDeal
	err := p.receiver.List(&deals)
	if err != nil {
		return nil, xerrors.Errorf("failed to fetch deals during startup: %w", err)
	}

	// re-track all deals for whom we still have a local blockstore.
	for _, d := range deals {
		if _, err := os.Stat(d.InboundCAR); err == nil && d.Ref != nil {
			_, _ = p.stores.GetOrOpen(d.ProposalCid.String(), d.InboundCAR, d.Ref.Root)
		}
	}

	// migrate deals to the dagstore if still not migrated.
	if ok, err := p.dagStore.MigrateDeals(ctx, deals); err != nil {
		return nil, fmt.Errorf("failed to migrate deals to DAG store: %w", err)
	} else if ok {
		log.Info("dagstore migration completed successfully")
	} else {
		log.Info("no dagstore migration necessary")
	}

	return deals, nil
}

func (p *ProviderExt) restartDeals(deals []storagemarket.MinerDeal) error {
	for _, deal := range deals {
		if p.receiver.IsTerminated(deal) {
			continue
		}

		err := p.receiver.Restart(deal.ProposalCid, storagemarket.ProviderEventRestart)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProviderExt) sign(ctx context.Context, data interface{}) (*crypto.Signature, error) {
	//tok, _, err := p.spn.GetChainHead(ctx)
	//if err != nil {
	//	return nil, xerrors.Errorf("couldn't get chain head: %w", err)
	//}
	//
	//return providerutils.SignMinerData(ctx, data, p.actor, tok, p.spn.GetMinerWorkerAddress, p.spn.SignBytes)
	return &acrypto.Signature{
		Type: crypto.SigTypeBLS,
		Data: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}, nil
}

func (p *ProviderExt) resendProposalResponse(s network.StorageDealStream, md *storagemarket.MinerDeal) error {
	resp := &network.Response{State: md.State, Message: md.Message, Proposal: md.ProposalCid}
	sig, err := p.sign(context.TODO(), resp)
	if err != nil {
		return xerrors.Errorf("failed to sign response message: %w", err)
	}

	err = s.WriteDealResponse(network.SignedResponse{Response: *resp, Signature: sig}, p.sign)

	if closeErr := s.Close(); closeErr != nil {
		log.Warnf("closing connection: %v", err)
	}

	return err
}

func newProviderStateMachine(ds datastore.Batching, env fsm.Environment, notifier fsm.Notifier, storageMigrations versioning.VersionedMigrationList, target versioning.VersionKey) (fsm.Group, func(context.Context) error, error) {
	return versionedfsm.NewVersionedFSM(ds, fsm.Parameters{
		Environment:     env,
		StateType:       storagemarket.MinerDeal{},
		StateKeyField:   "State",
		Events:          providerstates.ProviderEvents,
		StateEntryFuncs: providerstates.ProviderStateEntryFuncs,
		FinalityStates:  providerstates.ProviderFinalityStates,
		Notifier:        notifier,
	}, storageMigrations, target)
}

type internalProviderEvent struct {
	evt  storagemarket.ProviderEvent
	deal storagemarket.MinerDeal
}

func providerDispatcher(evt pubsub.Event, fn pubsub.SubscriberFn) error {
	ie, ok := evt.(internalProviderEvent)
	if !ok {
		return xerrors.New("wrong type of event")
	}
	cb, ok := fn.(storagemarket.ProviderSubscriber)
	if !ok {
		return xerrors.New("wrong type of callback")
	}
	cb(ie.evt, ie.deal)
	return nil
}

// ProviderFSMParameterSpec is a valid set of parameters for a providerext FSM - used in doc generation
var ProviderFSMParameterSpec = fsm.Parameters{
	Environment:     &providerExtDealEnvironment{},
	StateType:       storagemarket.MinerDeal{},
	StateKeyField:   "State",
	Events:          providerstates.ProviderEvents,
	StateEntryFuncs: providerstates.ProviderStateEntryFuncs,
	FinalityStates:  providerstates.ProviderFinalityStates,
}

func curTime() cbg.CborTime {
	now := time.Now()
	return cbg.CborTime(time.Unix(0, now.UnixNano()).UTC())
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

func (p *ProviderExt) messageResponse(md *storagemarket.MinerDeal, s network.StorageDealStream, proposalNd ipld.Node) {
	//Message return new cid
	resp := &network.Response{State: md.State, Message: md.Message, Proposal: md.ProposalCid}
	sig, err := p.sign(context.TODO(), resp)
	if err != nil {
		log.Warnf("sign: %v", err)
	}

	err = s.WriteDealResponse(network.SignedResponse{Response: *resp, Signature: sig}, p.sign)

	if closeErr := s.Close(); closeErr != nil {
		log.Warnf("closing connection: %v", err)
	}

	if err := p.environment.Disconnect(proposalNd.Cid()); err != nil {
		log.Warnf("closing client connection: %+v", err)
	}
}

func (p *ProviderExt) checkauthority(s network.StorageDealStream, proposal network.Proposal) (bool, error) {

	//dealinfo := clientext.BytesToDealInfo(proposal.DealProposal.ClientSignature.Data)
	var dealinfo clientext.DealInfo
	log.Info("Recv byte len:", len(proposal.DealProposal.ClientSignature.Data))
	log.Info("Recv byte :", string(proposal.DealProposal.ClientSignature.Data))
	//var buffer bytes.Buffer
	//b, err := proposal.DealProposal.ClientSignature.MarshalBinary()
	//if err == nil {
	//	buffer.Write(b)
	//}
	//b := bytes.NewReader()
	//dec := gob.NewDecoder(&buffer)
	//err = dec.Decode(&dealinfo)
	dataBuf := proposal.DealProposal.ClientSignature.Data
	err := json.Unmarshal(dataBuf, &dealinfo)

	if err == nil {

		//var carpath string
		pieceCID := cid.Undef

		proposalNd, err := cborutil.AsIpld(proposal.DealProposal)
		if err != nil {
			return dealinfo.Ishead, err
		}

		c := proposalNd.Cid()
		deal := &storagemarket.MinerDeal{
			Client:             s.RemotePeer(),
			Miner:              p.net.ID(),
			ClientDealProposal: *proposal.DealProposal,
			ProposalCid:        proposalNd.Cid(),
			PublishCid:         &c,
			State:              storagemarket.StorageDealUnknown,
			Ref:                proposal.Piece,
			FastRetrieval:      proposal.FastRetrieval,
			CreationTime:       curTime(),
			InboundCAR:         "",
		}

		err = p.receiver.Begin(proposalNd.Cid(), deal)
		if err != nil {
			return dealinfo.Ishead, err
		}

		defer p.receiver.End(proposalNd.Cid(), deal)
		err = p.conns.AddStream(proposalNd.Cid(), s)
		if err != nil {
			return dealinfo.Ishead, err
		}

		var md *storagemarket.MinerDeal
		a, ok := p.receiver.Get(proposalNd.Cid())
		if ok {
			md = a.(*storagemarket.MinerDeal)
		}
		md.State = storagemarket.StorageDealWaitingForData

		log.Info("PayloadCID  :", dealinfo)
		curCid, err := cid.Decode(dealinfo.FileCid)
		if err != nil {
			p.messageResponse(md, s, proposalNd)
			log.Errorf("Decode: %s", err)
			return dealinfo.Ishead, err
		}

		pieceInfo, _, err := p.getPieceInfoFromCid(p.ctx, curCid, pieceCID)
		if err != nil {
			p.messageResponse(md, s, proposalNd)
			log.Errorf("storage query: getPieceInfoFromCid: %s", err)
			return dealinfo.Ishead, err
		}

		bs, err := p.dagStore.LoadShard(p.ctx, pieceInfo.PieceCID)
		if err != nil {
			log.Errorf("failed to load blockstore for piece %s: %w", dealinfo.FileCid, err)
			p.messageResponse(md, s, proposalNd)
			return dealinfo.Ishead, err
		}
		defer bs.Close()

		blk, err := bs.Get(curCid)

		if err != nil {

		} else {

			voucherInfo, contractAddr, err := voucher.VoucherCli.TransactionCheck(dealinfo.Txhash)
			if err != nil {
				log.Error("failed to TransactionCheck: %w", err)
				//return dealinfo.RetrievalFile.Ishead, xerrors.Errorf("failed to TransactionCheck: %w", err)
			} else {
				//check validity
				//file_store.FileStoreCli
				storeval, err := file_store.FileStoreCli.GetBaseInfo(nil, voucherInfo.OriHash)

				log.Info("storeval : ", storeval)

				if err != nil {
					p.messageResponse(md, s, proposalNd)
					return dealinfo.Ishead, xerrors.Errorf("failed getfilestorebaseinfo: %w", err)
				}

				log.Info("voucherInfo : ", voucherInfo)
				//contract address check
				//buy or rent
				if voucherInfo.PurchaseType == 0 {

					//ori owner check by orihash
					txowner := voucher.VoucherCli.GetOwnerbytxhash(dealinfo.Txhash)
					if txowner == storeval.OwnerAddr {
						newhead, _, res := w3fsAuth.W3FS_Auth(voucherInfo.UserFilePubkey, blk.RawData(), -1, -1)
						md.Message = base64.StdEncoding.EncodeToString(newhead)
						if res != 0 {
							log.Error("message 0 : ", res)
							md.Message = ""
						}
					}

				} else {
					for _, temp := range storeval.DappContractAddrs {
						if contractAddr == temp {
							expirydate, _ := strconv.Atoi(voucherInfo.ValidityDate.String())
							newhead, _, res := w3fsAuth.W3FS_Auth(voucherInfo.UserFilePubkey, blk.RawData(), expirydate, 0)
							md.Message = base64.StdEncoding.EncodeToString(newhead)
							if res != 0 {
								log.Error("message  1: ", res)
								md.Message = ""
							}
						}
					}

				}

			}

		}

		//Message return new cid
		p.messageResponse(md, s, proposalNd)
		return dealinfo.Ishead, err
	} else {
		log.Info("Decode deal : ", err)

		if strings.Contains(string(proposal.DealProposal.ClientSignature.Data), ",") {
			return false, nil
		}
	}

	return false, nil
}
