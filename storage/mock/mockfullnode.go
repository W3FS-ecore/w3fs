package mock

import (
	context "context"
	json "encoding/json"
	address "github.com/filecoin-project/go-address"
	bitfield "github.com/filecoin-project/go-bitfield"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	retrievalmarket "github.com/filecoin-project/go-fil-markets/retrievalmarket"
	storagemarket "github.com/filecoin-project/go-fil-markets/storagemarket"
	auth "github.com/filecoin-project/go-jsonrpc/auth"
	abi "github.com/filecoin-project/go-state-types/abi"
	big "github.com/filecoin-project/go-state-types/big"
	acrypto "github.com/filecoin-project/go-state-types/crypto"
	crypto "github.com/filecoin-project/go-state-types/crypto"
	dline "github.com/filecoin-project/go-state-types/dline"
	network "github.com/filecoin-project/go-state-types/network"
	api "github.com/filecoin-project/lotus/api"
	apitypes "github.com/filecoin-project/lotus/api/types"
	"github.com/filecoin-project/lotus/chain/actors/builtin/market"
	miner "github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/store"
	types "github.com/filecoin-project/lotus/chain/types"
	alerting "github.com/filecoin-project/lotus/journal/alerting"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	dtypes "github.com/filecoin-project/lotus/node/modules/dtypes"
	imports "github.com/filecoin-project/lotus/node/repo/imports"
	miner0 "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	paych "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	builtin6 "github.com/filecoin-project/specs-actors/v6/actors/builtin"
	uuid "github.com/google/uuid"
	cid "github.com/ipfs/go-cid"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	network0 "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	"math/rand"
)

// MockFullNode is a mock of FullNode interface.
type MockFullNode struct {
}

// NewMockFullNode creates a new mock instance.
func NewMockFullNode() *MockFullNode {
	mock := &MockFullNode{}
	return mock
}

// AuthNew mocks base method.
func (m *MockFullNode) AuthNew(arg0 context.Context, arg1 []auth.Permission) ([]byte, error) {
	return nil, nil
}

// AuthVerify mocks base method.
func (m *MockFullNode) AuthVerify(arg0 context.Context, arg1 string) ([]auth.Permission, error) {
	return nil, nil
}

// BeaconGetEntry mocks base method.
func (m *MockFullNode) BeaconGetEntry(arg0 context.Context, arg1 abi.ChainEpoch) (*types.BeaconEntry, error) {
	return nil, nil
}

// ChainBlockstoreInfo mocks base method.
func (m *MockFullNode) ChainBlockstoreInfo(arg0 context.Context) (map[string]interface{}, error) {
	return nil, nil
}

// ChainCheckBlockstore mocks base method.
func (m *MockFullNode) ChainCheckBlockstore(arg0 context.Context) error {
	return nil
}

// ChainDeleteObj mocks base method.
func (m *MockFullNode) ChainDeleteObj(arg0 context.Context, arg1 cid.Cid) error {
	return nil
}

// ChainExport mocks base method.
func (m *MockFullNode) ChainExport(arg0 context.Context, arg1 abi.ChainEpoch, arg2 bool, arg3 types.TipSetKey) (<-chan []byte, error) {
	return nil, nil
}

// ChainGetBlock mocks base method.
func (m *MockFullNode) ChainGetBlock(arg0 context.Context, arg1 cid.Cid) (*types.BlockHeader, error) {
	return nil, nil
}

// ChainGetBlockMessages mocks base method.
func (m *MockFullNode) ChainGetBlockMessages(arg0 context.Context, arg1 cid.Cid) (*api.BlockMessages, error) {
	return nil, nil
}

// ChainGetGenesis mocks base method.
func (m *MockFullNode) ChainGetGenesis(arg0 context.Context) (*types.TipSet, error) {
	return nil, nil
}

// ChainGetMessage mocks base method.
func (m *MockFullNode) ChainGetMessage(arg0 context.Context, arg1 cid.Cid) (*types.Message, error) {
	return nil, nil
}

// ChainGetMessagesInTipset mocks base method.
func (m *MockFullNode) ChainGetMessagesInTipset(arg0 context.Context, arg1 types.TipSetKey) ([]api.Message, error) {
	return nil, nil
}

// ChainGetNode mocks base method.
func (m *MockFullNode) ChainGetNode(arg0 context.Context, arg1 string) (*api.IpldObject, error) {
	var ret []interface{}
	ret0, _ := ret[0].(*api.IpldObject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChainGetParentMessages mocks base method.
func (m *MockFullNode) ChainGetParentMessages(arg0 context.Context, arg1 cid.Cid) ([]api.Message, error) {
	return nil, nil
}

// ChainGetParentReceipts mocks base method.
func (m *MockFullNode) ChainGetParentReceipts(arg0 context.Context, arg1 cid.Cid) ([]*types.MessageReceipt, error) {
	return nil, nil
}

// ChainGetPath mocks base method.
func (m *MockFullNode) ChainGetPath(arg0 context.Context, arg1, arg2 types.TipSetKey) ([]*api.HeadChange, error) {
	return nil, nil
}

// ChainGetTipSet mocks base method.
func (m *MockFullNode) ChainGetTipSet(arg0 context.Context, arg1 types.TipSetKey) (*types.TipSet, error) {
	return nil, nil
}

// ChainGetTipSetAfterHeight mocks base method.
func (m *MockFullNode) ChainGetTipSetAfterHeight(arg0 context.Context, arg1 abi.ChainEpoch, arg2 types.TipSetKey) (*types.TipSet, error) {
	return nil, nil
}

// ChainGetTipSetByHeight mocks base method.
func (m *MockFullNode) ChainGetTipSetByHeight(arg0 context.Context, arg1 abi.ChainEpoch, arg2 types.TipSetKey) (*types.TipSet, error) {
	return nil, nil
}

// ChainHasObj mocks base method.
func (m *MockFullNode) ChainHasObj(arg0 context.Context, arg1 cid.Cid) (bool, error) {
	return false, nil
}

// ChainHead mocks base method.
func (m *MockFullNode) ChainHead(arg0 context.Context) (*types.TipSet, error) {
	a, _ := address.NewFromString("t0000")
	//b, _ := address.NewFromString("t02")
	dummyCid, _ := cid.Parse("bafkqaaa")
	h := abi.ChainEpoch(1000)
	var ts, err = types.NewTipSet([]*types.BlockHeader{
		{
			Height: h,
			Miner:  a,

			Parents: nil,

			Ticket: &types.Ticket{VRFProof: []byte{byte(h % 2)}},

			ParentStateRoot:       dummyCid,
			Messages:              dummyCid,
			ParentMessageReceipts: dummyCid,

			BlockSig:     &crypto.Signature{Type: crypto.SigTypeBLS},
			BLSAggregate: &crypto.Signature{Type: crypto.SigTypeBLS},
		},
		{
			Height: h,
			Miner:  a,

			Parents: nil,

			Ticket: &types.Ticket{VRFProof: []byte{byte((h + 1) % 2)}},

			ParentStateRoot:       dummyCid,
			Messages:              dummyCid,
			ParentMessageReceipts: dummyCid,

			BlockSig:     &crypto.Signature{Type: crypto.SigTypeBLS},
			BLSAggregate: &crypto.Signature{Type: crypto.SigTypeBLS},
		},
	})

	if err != nil {
		return nil, nil
	}

	return ts, nil
}

// ChainNotify mocks base method.
func (m *MockFullNode) ChainNotify(arg0 context.Context) (<-chan []*api.HeadChange, error) {
	out := make(chan []*api.HeadChange, 16)
	out <- []*api.HeadChange{{
		Type: store.HCCurrent,
		Val:  &types.TipSet{},
	}}
	return out, nil
}

// ChainReadObj mocks base method.
func (m *MockFullNode) ChainReadObj(arg0 context.Context, arg1 cid.Cid) ([]byte, error) {
	return nil, nil
}

// ChainSetHead mocks base method.
func (m *MockFullNode) ChainSetHead(arg0 context.Context, arg1 types.TipSetKey) error {
	var ret []interface{}
	ret0, _ := ret[0].(error)
	return ret0
}

// ChainStatObj mocks base method.
func (m *MockFullNode) ChainStatObj(arg0 context.Context, arg1, arg2 cid.Cid) (api.ObjStat, error) {
	return api.ObjStat{}, nil
}

// ChainTipSetWeight mocks base method.
func (m *MockFullNode) ChainTipSetWeight(arg0 context.Context, arg1 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// ClientCalcCommP mocks base method.
func (m *MockFullNode) ClientCalcCommP(arg0 context.Context, arg1 string) (*api.CommPRet, error) {
	return nil, nil
}

// ClientCancelDataTransfer mocks base method.
func (m *MockFullNode) ClientCancelDataTransfer(arg0 context.Context, arg1 datatransfer.TransferID, arg2 peer.ID, arg3 bool) error {
	return nil
}

// ClientCancelRetrievalDeal mocks base method.
func (m *MockFullNode) ClientCancelRetrievalDeal(arg0 context.Context, arg1 retrievalmarket.DealID) error {
	return nil
}

// ClientDataTransferUpdates mocks base method.
func (m *MockFullNode) ClientDataTransferUpdates(arg0 context.Context) (<-chan api.DataTransferChannel, error) {
	return nil, nil
}

// ClientDealPieceCID mocks base method.
func (m *MockFullNode) ClientDealPieceCID(arg0 context.Context, arg1 cid.Cid) (api.DataCIDSize, error) {
	return api.DataCIDSize{}, nil
}

// ClientDealSize mocks base method.
func (m *MockFullNode) ClientDealSize(arg0 context.Context, arg1 cid.Cid) (api.DataSize, error) {
	return api.DataSize{}, nil
}

// ClientFindData mocks base method.
func (m *MockFullNode) ClientFindData(arg0 context.Context, arg1 cid.Cid, arg2 *cid.Cid) ([]api.QueryOffer, error) {
	return nil, nil
}

// ClientGenCar mocks base method.
func (m *MockFullNode) ClientGenCar(arg0 context.Context, arg1 api.FileRef, arg2 string) error {
	return nil
}

// ClientGetDealInfo mocks base method.
func (m *MockFullNode) ClientGetDealInfo(arg0 context.Context, arg1 cid.Cid) (*api.DealInfo, error) {
	return nil, nil
}

// ClientGetDealStatus mocks base method.
func (m *MockFullNode) ClientGetDealStatus(arg0 context.Context, arg1 uint64) (string, error) {
	return "", nil
}

// ClientGetDealUpdates mocks base method.
func (m *MockFullNode) ClientGetDealUpdates(arg0 context.Context) (<-chan api.DealInfo, error) {
	return nil, nil
}

// ClientGetRetrievalUpdates mocks base method.
func (m *MockFullNode) ClientGetRetrievalUpdates(arg0 context.Context) (<-chan api.RetrievalInfo, error) {
	return nil, nil
}

// ClientHasLocal mocks base method.
func (m *MockFullNode) ClientHasLocal(arg0 context.Context, arg1 cid.Cid) (bool, error) {
	return true, nil
}

// ClientImport mocks base method.
func (m *MockFullNode) ClientImport(arg0 context.Context, arg1 api.FileRef) (*api.ImportRes, error) {
	return nil, nil
}

// ClientListDataTransfers mocks base method.
func (m *MockFullNode) ClientListDataTransfers(arg0 context.Context) ([]api.DataTransferChannel, error) {
	return nil, nil
}

// ClientListDeals mocks base method.
func (m *MockFullNode) ClientListDeals(arg0 context.Context) ([]api.DealInfo, error) {
	return nil, nil
}

// ClientListImports mocks base method.
func (m *MockFullNode) ClientListImports(arg0 context.Context) ([]api.Import, error) {
	return nil, nil
}

// ClientListRetrievals mocks base method.
func (m *MockFullNode) ClientListRetrievals(arg0 context.Context) ([]api.RetrievalInfo, error) {
	return nil, nil
}

// ClientMinerQueryOffer mocks base method.
func (m *MockFullNode) ClientMinerQueryOffer(arg0 context.Context, arg1 address.Address, arg2 cid.Cid, arg3 *cid.Cid) (api.QueryOffer, error) {
	return api.QueryOffer{}, nil
}

// ClientQueryAsk mocks base method.
func (m *MockFullNode) ClientQueryAsk(arg0 context.Context, arg1 peer.ID, arg2 address.Address) (*storagemarket.StorageAsk, error) {
	return nil, nil
}

// ClientRemoveImport mocks base method.
func (m *MockFullNode) ClientRemoveImport(arg0 context.Context, arg1 imports.ID) error {
	return nil
}

// ClientRestartDataTransfer mocks base method.
func (m *MockFullNode) ClientRestartDataTransfer(arg0 context.Context, arg1 datatransfer.TransferID, arg2 peer.ID, arg3 bool) error {
	return nil
}

// ClientRetrieve mocks base method.
func (m *MockFullNode) ClientRetrieve(arg0 context.Context, arg1 api.RetrievalOrder, arg2 *api.FileRef) error {
	return nil
}

// ClientRetrieveTryRestartInsufficientFunds mocks base method.
func (m *MockFullNode) ClientRetrieveTryRestartInsufficientFunds(arg0 context.Context, arg1 address.Address) error {
	return nil
}

// ClientRetrieveWithEvents mocks base method.
func (m *MockFullNode) ClientRetrieveWithEvents(arg0 context.Context, arg1 api.RetrievalOrder, arg2 *api.FileRef) (<-chan marketevents.RetrievalEvent, error) {
	return nil, nil
}

// ClientStartDeal mocks base method.
func (m *MockFullNode) ClientStartDeal(arg0 context.Context, arg1 *api.StartDealParams) (*cid.Cid, error) {
	return nil, nil
}

// ClientStatelessDeal mocks base method.
func (m *MockFullNode) ClientStatelessDeal(arg0 context.Context, arg1 *api.StartDealParams) (*cid.Cid, error) {
	return nil, nil
}

// Closing mocks base method.
func (m *MockFullNode) Closing(arg0 context.Context) (<-chan struct{}, error) {
	return nil, nil
}

// CreateBackup mocks base method.
func (m *MockFullNode) CreateBackup(arg0 context.Context, arg1 string) error {
	return nil
}

// Discover mocks base method.
func (m *MockFullNode) Discover(arg0 context.Context) (apitypes.OpenRPCDocument, error) {
	return nil, nil
}

// GasEstimateFeeCap mocks base method.
func (m *MockFullNode) GasEstimateFeeCap(arg0 context.Context, arg1 *types.Message, arg2 int64, arg3 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// GasEstimateGasLimit mocks base method.
func (m *MockFullNode) GasEstimateGasLimit(arg0 context.Context, arg1 *types.Message, arg2 types.TipSetKey) (int64, error) {
	return 0, nil
}

// GasEstimateGasPremium mocks base method.
func (m *MockFullNode) GasEstimateGasPremium(arg0 context.Context, arg1 uint64, arg2 address.Address, arg3 int64, arg4 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// GasEstimateMessageGas mocks base method.
func (m *MockFullNode) GasEstimateMessageGas(arg0 context.Context, arg1 *types.Message, arg2 *api.MessageSendSpec, arg3 types.TipSetKey) (*types.Message, error) {
	return nil, nil
}

// ID mocks base method.
func (m *MockFullNode) ID(arg0 context.Context) (peer.ID, error) {
	return "", nil
}

// LogAlerts mocks base method.
func (m *MockFullNode) LogAlerts(arg0 context.Context) ([]alerting.Alert, error) {
	return nil, nil
}

// LogList mocks base method.
func (m *MockFullNode) LogList(arg0 context.Context) ([]string, error) {
	return nil, nil
}

// LogSetLevel mocks base method.
func (m *MockFullNode) LogSetLevel(arg0 context.Context, arg1, arg2 string) error {
	return nil
}

// MarketAddBalance mocks base method.
func (m *MockFullNode) MarketAddBalance(arg0 context.Context, arg1, arg2 address.Address, arg3 big.Int) (cid.Cid, error) {
	return cid.Undef, nil
}

// MarketGetReserved mocks base method.
func (m *MockFullNode) MarketGetReserved(arg0 context.Context, arg1 address.Address) (big.Int, error) {
	return big.Zero(), nil
}

// MarketReleaseFunds mocks base method.
func (m *MockFullNode) MarketReleaseFunds(arg0 context.Context, arg1 address.Address, arg2 big.Int) error {
	return nil
}

// MarketReserveFunds mocks base method.
func (m *MockFullNode) MarketReserveFunds(arg0 context.Context, arg1, arg2 address.Address, arg3 big.Int) (cid.Cid, error) {
	return cid.Undef, nil
}

// MarketWithdraw mocks base method.
func (m *MockFullNode) MarketWithdraw(arg0 context.Context, arg1, arg2 address.Address, arg3 big.Int) (cid.Cid, error) {
	return cid.Undef, nil
}

// MinerCreateBlock mocks base method.
func (m *MockFullNode) MinerCreateBlock(arg0 context.Context, arg1 *api.BlockTemplate) (*types.BlockMsg, error) {
	return nil, nil
}

// MinerGetBaseInfo mocks base method.
func (m *MockFullNode) MinerGetBaseInfo(arg0 context.Context, arg1 address.Address, arg2 abi.ChainEpoch, arg3 types.TipSetKey) (*api.MiningBaseInfo, error) {
	return nil, nil
}

// MpoolBatchPush mocks base method.
func (m *MockFullNode) MpoolBatchPush(arg0 context.Context, arg1 []*types.SignedMessage) ([]cid.Cid, error) {
	return nil, nil
}

// MpoolBatchPushMessage mocks base method.
func (m *MockFullNode) MpoolBatchPushMessage(arg0 context.Context, arg1 []*types.Message, arg2 *api.MessageSendSpec) ([]*types.SignedMessage, error) {
	return nil, nil
}

// MpoolBatchPushUntrusted mocks base method.
func (m *MockFullNode) MpoolBatchPushUntrusted(arg0 context.Context, arg1 []*types.SignedMessage) ([]cid.Cid, error) {
	return nil, nil
}

// MpoolCheckMessages mocks base method.
func (m *MockFullNode) MpoolCheckMessages(arg0 context.Context, arg1 []*api.MessagePrototype) ([][]api.MessageCheckStatus, error) {
	return nil, nil
}

// MpoolCheckPendingMessages mocks base method.
func (m *MockFullNode) MpoolCheckPendingMessages(arg0 context.Context, arg1 address.Address) ([][]api.MessageCheckStatus, error) {
	return nil, nil
}

// MpoolCheckReplaceMessages mocks base method.
func (m *MockFullNode) MpoolCheckReplaceMessages(arg0 context.Context, arg1 []*types.Message) ([][]api.MessageCheckStatus, error) {
	return nil, nil
}

// MpoolClear mocks base method.
func (m *MockFullNode) MpoolClear(arg0 context.Context, arg1 bool) error {
	return nil
}

// MpoolGetConfig mocks base method.
func (m *MockFullNode) MpoolGetConfig(arg0 context.Context) (*types.MpoolConfig, error) {
	return nil, nil
}

// MpoolGetNonce mocks base method.
func (m *MockFullNode) MpoolGetNonce(arg0 context.Context, arg1 address.Address) (uint64, error) {
	return 0, nil
}

// MpoolPending mocks base method.
func (m *MockFullNode) MpoolPending(arg0 context.Context, arg1 types.TipSetKey) ([]*types.SignedMessage, error) {
	return nil, nil
}

// MpoolPush mocks base method.
func (m *MockFullNode) MpoolPush(arg0 context.Context, arg1 *types.SignedMessage) (cid.Cid, error) {
	return cid.Undef, nil
}

// MpoolPushMessage mocks base method.
func (m *MockFullNode) MpoolPushMessage(arg0 context.Context, arg1 *types.Message, arg2 *api.MessageSendSpec) (*types.SignedMessage, error) {
	smsg := &types.SignedMessage{Message: *arg1,
		Signature: acrypto.Signature{
			Type: crypto.SigTypeBLS,
			Data: arg1.Params[:64],
		}}
	return smsg, nil
}

// MpoolPushUntrusted mocks base method.
func (m *MockFullNode) MpoolPushUntrusted(arg0 context.Context, arg1 *types.SignedMessage) (cid.Cid, error) {
	return cid.Undef, nil
}

// MpoolSelect mocks base method.
func (m *MockFullNode) MpoolSelect(arg0 context.Context, arg1 types.TipSetKey, arg2 float64) ([]*types.SignedMessage, error) {
	return nil, nil
}

// MpoolSetConfig mocks base method.
func (m *MockFullNode) MpoolSetConfig(arg0 context.Context, arg1 *types.MpoolConfig) error {
	return nil
}

// MpoolSub mocks base method.
func (m *MockFullNode) MpoolSub(arg0 context.Context) (<-chan api.MpoolUpdate, error) {
	return nil, nil
}

// MsigAddApprove mocks base method.
func (m *MockFullNode) MsigAddApprove(arg0 context.Context, arg1, arg2 address.Address, arg3 uint64, arg4, arg5 address.Address, arg6 bool) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigAddCancel mocks base method.
func (m *MockFullNode) MsigAddCancel(arg0 context.Context, arg1, arg2 address.Address, arg3 uint64, arg4 address.Address, arg5 bool) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigAddPropose mocks base method.
func (m *MockFullNode) MsigAddPropose(arg0 context.Context, arg1, arg2, arg3 address.Address, arg4 bool) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigApprove mocks base method.
func (m *MockFullNode) MsigApprove(arg0 context.Context, arg1 address.Address, arg2 uint64, arg3 address.Address) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigApproveTxnHash mocks base method.
func (m *MockFullNode) MsigApproveTxnHash(arg0 context.Context, arg1 address.Address, arg2 uint64, arg3, arg4 address.Address, arg5 big.Int, arg6 address.Address, arg7 uint64, arg8 []byte) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigCancel mocks base method.
func (m *MockFullNode) MsigCancel(arg0 context.Context, arg1 address.Address, arg2 uint64, arg3 address.Address, arg4 big.Int, arg5 address.Address, arg6 uint64, arg7 []byte) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigCreate mocks base method.
func (m *MockFullNode) MsigCreate(arg0 context.Context, arg1 uint64, arg2 []address.Address, arg3 abi.ChainEpoch, arg4 big.Int, arg5 address.Address, arg6 big.Int) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigGetAvailableBalance mocks base method.
func (m *MockFullNode) MsigGetAvailableBalance(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// MsigGetPending mocks base method.
func (m *MockFullNode) MsigGetPending(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) ([]*api.MsigTransaction, error) {
	return nil, nil
}

// MsigGetVested mocks base method.
func (m *MockFullNode) MsigGetVested(arg0 context.Context, arg1 address.Address, arg2, arg3 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// MsigGetVestingSchedule mocks base method.
func (m *MockFullNode) MsigGetVestingSchedule(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (api.MsigVesting, error) {
	return api.MsigVesting{}, nil
}

// MsigPropose mocks base method.
func (m *MockFullNode) MsigPropose(arg0 context.Context, arg1, arg2 address.Address, arg3 big.Int, arg4 address.Address, arg5 uint64, arg6 []byte) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigRemoveSigner mocks base method.
func (m *MockFullNode) MsigRemoveSigner(arg0 context.Context, arg1, arg2, arg3 address.Address, arg4 bool) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigSwapApprove mocks base method.
func (m *MockFullNode) MsigSwapApprove(arg0 context.Context, arg1, arg2 address.Address, arg3 uint64, arg4, arg5, arg6 address.Address) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigSwapCancel mocks base method.
func (m *MockFullNode) MsigSwapCancel(arg0 context.Context, arg1, arg2 address.Address, arg3 uint64, arg4, arg5 address.Address) (*api.MessagePrototype, error) {
	return nil, nil
}

// MsigSwapPropose mocks base method.
func (m *MockFullNode) MsigSwapPropose(arg0 context.Context, arg1, arg2, arg3, arg4 address.Address) (*api.MessagePrototype, error) {
	return nil, nil
}

// NetAddrsListen mocks base method.
func (m *MockFullNode) NetAddrsListen(arg0 context.Context) (peer.AddrInfo, error) {
	return peer.AddrInfo{}, nil
}

// NetAgentVersion mocks base method.
func (m *MockFullNode) NetAgentVersion(arg0 context.Context, arg1 peer.ID) (string, error) {
	return "", nil
}

// NetAutoNatStatus mocks base method.
func (m *MockFullNode) NetAutoNatStatus(arg0 context.Context) (api.NatInfo, error) {
	return api.NatInfo{}, nil
}

// NetBandwidthStats mocks base method.
func (m *MockFullNode) NetBandwidthStats(arg0 context.Context) (metrics.Stats, error) {
	return metrics.Stats{}, nil
}

// NetBandwidthStatsByPeer mocks base method.
func (m *MockFullNode) NetBandwidthStatsByPeer(arg0 context.Context) (map[string]metrics.Stats, error) {
	return nil, nil
}

// NetBandwidthStatsByProtocol mocks base method.
func (m *MockFullNode) NetBandwidthStatsByProtocol(arg0 context.Context) (map[protocol.ID]metrics.Stats, error) {
	return nil, nil
}

// NetBlockAdd mocks base method.
func (m *MockFullNode) NetBlockAdd(arg0 context.Context, arg1 api.NetBlockList) error {
	return nil
}

// NetBlockList mocks base method.
func (m *MockFullNode) NetBlockList(arg0 context.Context) (api.NetBlockList, error) {
	return api.NetBlockList{}, nil
}

// NetBlockRemove mocks base method.
func (m *MockFullNode) NetBlockRemove(arg0 context.Context, arg1 api.NetBlockList) error {
	var ret []interface{}
	ret0, _ := ret[0].(error)
	return ret0
}

// NetConnect mocks base method.
func (m *MockFullNode) NetConnect(arg0 context.Context, arg1 peer.AddrInfo) error {
	return nil
}

// NetConnectedness mocks base method.
func (m *MockFullNode) NetConnectedness(arg0 context.Context, arg1 peer.ID) (network0.Connectedness, error) {
	return 0, nil
}

// NetDisconnect mocks base method.
func (m *MockFullNode) NetDisconnect(arg0 context.Context, arg1 peer.ID) error {
	return nil
}

// NetFindPeer mocks base method.
func (m *MockFullNode) NetFindPeer(arg0 context.Context, arg1 peer.ID) (peer.AddrInfo, error) {
	return peer.AddrInfo{}, nil
}

// NetPeerInfo mocks base method.
func (m *MockFullNode) NetPeerInfo(arg0 context.Context, arg1 peer.ID) (*api.ExtendedPeerInfo, error) {
	return nil, nil
}

// NetPeers mocks base method.
func (m *MockFullNode) NetPeers(arg0 context.Context) ([]peer.AddrInfo, error) {
	return nil, nil
}

// NetPubsubScores mocks base method.
func (m *MockFullNode) NetPubsubScores(arg0 context.Context) ([]api.PubsubScore, error) {
	return nil, nil
}

// NodeStatus mocks base method.
func (m *MockFullNode) NodeStatus(arg0 context.Context, arg1 bool) (api.NodeStatus, error) {
	return api.NodeStatus{}, nil
}

// PaychAllocateLane mocks base method.
func (m *MockFullNode) PaychAllocateLane(arg0 context.Context, arg1 address.Address) (uint64, error) {
	return 0, nil
}

// PaychAvailableFunds mocks base method.
func (m *MockFullNode) PaychAvailableFunds(arg0 context.Context, arg1 address.Address) (*api.ChannelAvailableFunds, error) {
	return nil, nil
}

// PaychAvailableFundsByFromTo mocks base method.
func (m *MockFullNode) PaychAvailableFundsByFromTo(arg0 context.Context, arg1, arg2 address.Address) (*api.ChannelAvailableFunds, error) {
	return nil, nil
}

// PaychCollect mocks base method.
func (m *MockFullNode) PaychCollect(arg0 context.Context, arg1 address.Address) (cid.Cid, error) {
	return cid.Undef, nil
}

// PaychGet mocks base method.
func (m *MockFullNode) PaychGet(arg0 context.Context, arg1, arg2 address.Address, arg3 big.Int) (*api.ChannelInfo, error) {
	return nil, nil
}

// PaychGetWaitReady mocks base method.
func (m *MockFullNode) PaychGetWaitReady(arg0 context.Context, arg1 cid.Cid) (address.Address, error) {
	return address.Undef, nil
}

// PaychList mocks base method.
func (m *MockFullNode) PaychList(arg0 context.Context) ([]address.Address, error) {
	return []address.Address{address.Undef}, nil
}

// PaychNewPayment mocks base method.
func (m *MockFullNode) PaychNewPayment(arg0 context.Context, arg1, arg2 address.Address, arg3 []api.VoucherSpec) (*api.PaymentInfo, error) {
	return nil, nil
}

// PaychSettle mocks base method.
func (m *MockFullNode) PaychSettle(arg0 context.Context, arg1 address.Address) (cid.Cid, error) {
	return cid.Undef, nil
}

// PaychStatus mocks base method.
func (m *MockFullNode) PaychStatus(arg0 context.Context, arg1 address.Address) (*api.PaychStatus, error) {
	return nil, nil
}

// PaychVoucherAdd mocks base method.
func (m *MockFullNode) PaychVoucherAdd(arg0 context.Context, arg1 address.Address, arg2 *paych.SignedVoucher, arg3 []byte, arg4 big.Int) (big.Int, error) {
	return big.Zero(), nil
}

// PaychVoucherCheckSpendable mocks base method.
func (m *MockFullNode) PaychVoucherCheckSpendable(arg0 context.Context, arg1 address.Address, arg2 *paych.SignedVoucher, arg3, arg4 []byte) (bool, error) {
	return false, nil
}

// PaychVoucherCheckValid mocks base method.
func (m *MockFullNode) PaychVoucherCheckValid(arg0 context.Context, arg1 address.Address, arg2 *paych.SignedVoucher) error {
	return nil
}

// PaychVoucherCreate mocks base method.
func (m *MockFullNode) PaychVoucherCreate(arg0 context.Context, arg1 address.Address, arg2 big.Int, arg3 uint64) (*api.VoucherCreateResult, error) {
	return nil, nil
}

// PaychVoucherList mocks base method.
func (m *MockFullNode) PaychVoucherList(arg0 context.Context, arg1 address.Address) ([]*paych.SignedVoucher, error) {
	return nil, nil
}

// PaychVoucherSubmit mocks base method.
func (m *MockFullNode) PaychVoucherSubmit(arg0 context.Context, arg1 address.Address, arg2 *paych.SignedVoucher, arg3, arg4 []byte) (cid.Cid, error) {
	return cid.Undef, nil
}

// Session mocks base method.
func (m *MockFullNode) Session(arg0 context.Context) (uuid.UUID, error) {
	return uuid.Nil, nil
}

// Shutdown mocks base method.
func (m *MockFullNode) Shutdown(arg0 context.Context) error {
	return nil
}

// StateAccountKey mocks base method.
func (m *MockFullNode) StateAccountKey(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (address.Address, error) {
	return address.Undef, nil
}

// StateAllMinerFaults mocks base method.
func (m *MockFullNode) StateAllMinerFaults(arg0 context.Context, arg1 abi.ChainEpoch, arg2 types.TipSetKey) ([]*api.Fault, error) {
	return nil, nil
}

// StateCall mocks base method.
func (m *MockFullNode) StateCall(arg0 context.Context, arg1 *types.Message, arg2 types.TipSetKey) (*api.InvocResult, error) {
	return nil, nil
}

// StateChangedActors mocks base method.
func (m *MockFullNode) StateChangedActors(arg0 context.Context, arg1, arg2 cid.Cid) (map[string]types.Actor, error) {
	return nil, nil
}

// StateCirculatingSupply mocks base method.
func (m *MockFullNode) StateCirculatingSupply(arg0 context.Context, arg1 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// StateCompute mocks base method.
func (m *MockFullNode) StateCompute(arg0 context.Context, arg1 abi.ChainEpoch, arg2 []*types.Message, arg3 types.TipSetKey) (*api.ComputeStateOutput, error) {
	return nil, nil
}

// StateDealProviderCollateralBounds mocks base method.
func (m *MockFullNode) StateDealProviderCollateralBounds(arg0 context.Context, arg1 abi.PaddedPieceSize, arg2 bool, arg3 types.TipSetKey) (api.DealCollateralBounds, error) {
	return api.DealCollateralBounds{}, nil
}

// StateDecodeParams mocks base method.
func (m *MockFullNode) StateDecodeParams(arg0 context.Context, arg1 address.Address, arg2 abi.MethodNum, arg3 []byte, arg4 types.TipSetKey) (interface{}, error) {
	return nil, nil
}

// StateEncodeParams mocks base method.
func (m *MockFullNode) StateEncodeParams(arg0 context.Context, arg1 cid.Cid, arg2 abi.MethodNum, arg3 json.RawMessage) ([]byte, error) {
	return nil, nil
}

// StateGetActor mocks base method.
func (m *MockFullNode) StateGetActor(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (*types.Actor, error) {
	return &types.Actor{
		Code: builtin6.StorageMinerActorCodeID,
		Head: cid.Undef,
	}, nil
}

// StateGetRandomnessFromBeacon mocks base method.
func (m *MockFullNode) StateGetRandomnessFromBeacon(arg0 context.Context, arg1 crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, arg3 []byte, arg4 types.TipSetKey) (abi.Randomness, error) {
	out := make([]byte, 32)
	_, _ = rand.New(rand.NewSource(int64(randEpoch * 1000))).Read(out) //nolint
	return out, nil
}

// StateGetRandomnessFromTickets mocks base method.
func (m *MockFullNode) StateGetRandomnessFromTickets(arg0 context.Context, arg1 crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, arg3 []byte, arg4 types.TipSetKey) (abi.Randomness, error) {
	out := make([]byte, 32)
	_, _ = rand.New(rand.NewSource(int64(randEpoch * 1000))).Read(out) //nolint
	return out, nil
}

// StateListActors mocks base method.
func (m *MockFullNode) StateListActors(arg0 context.Context, arg1 types.TipSetKey) ([]address.Address, error) {
	return []address.Address{address.Undef}, nil
}

// StateListMessages mocks base method.
func (m *MockFullNode) StateListMessages(arg0 context.Context, arg1 *api.MessageMatch, arg2 types.TipSetKey, arg3 abi.ChainEpoch) ([]cid.Cid, error) {
	return []cid.Cid{cid.Undef}, nil
}

// StateListMiners mocks base method.
func (m *MockFullNode) StateListMiners(arg0 context.Context, arg1 types.TipSetKey) ([]address.Address, error) {
	return []address.Address{address.Undef}, nil
}

// StateLookupID mocks base method.
func (m *MockFullNode) StateLookupID(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (address.Address, error) {
	return address.Undef, nil
}

// StateMarketBalance mocks base method.
func (m *MockFullNode) StateMarketBalance(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (api.MarketBalance, error) {
	return api.MarketBalance{}, nil
}

// StateMarketDeals mocks base method.
func (m *MockFullNode) StateMarketDeals(arg0 context.Context, arg1 types.TipSetKey) (map[string]api.MarketDeal, error) {
	return nil, nil
}

// StateMarketParticipants mocks base method.
func (m *MockFullNode) StateMarketParticipants(arg0 context.Context, arg1 types.TipSetKey) (map[string]api.MarketBalance, error) {
	return nil, nil
}

// StateMarketStorageDeal mocks base method.
func (m *MockFullNode) StateMarketStorageDeal(arg0 context.Context, arg1 abi.DealID, arg2 types.TipSetKey) (*api.MarketDeal, error) {
	return &api.MarketDeal{
		Proposal: market.DealProposal{},
	}, nil
}

// StateMinerActiveSectors mocks base method.
func (m *MockFullNode) StateMinerActiveSectors(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) ([]*miner.SectorOnChainInfo, error) {
	return nil, nil
}

// StateMinerAvailableBalance mocks base method.
func (m *MockFullNode) StateMinerAvailableBalance(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// StateMinerDeadlines mocks base method.
func (m *MockFullNode) StateMinerDeadlines(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) ([]api.Deadline, error) {
	return []api.Deadline{}, nil
}

// StateMinerFaults mocks base method.
func (m *MockFullNode) StateMinerFaults(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (bitfield.BitField, error) {
	return bitfield.New(), nil
}

// StateMinerInfo mocks base method.
func (m *MockFullNode) StateMinerInfo(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (miner.MinerInfo, error) {

	return miner.MinerInfo{
		WindowPoStProofType: abi.RegisteredPoStProof_StackedDrgWindow512MiBV1,
		SectorSize:          512 * 1 << 20,
	}, nil
	/*return miner.MinerInfo{
		WindowPoStProofType: abi.RegisteredPoStProof_StackedDrgWindow2KiBV1,
		SectorSize:          2 << 10,
	}, nil*/
}

// StateMinerInitialPledgeCollateral mocks base method.
func (m *MockFullNode) StateMinerInitialPledgeCollateral(arg0 context.Context, arg1 address.Address, arg2 miner0.SectorPreCommitInfo, arg3 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// StateMinerPartitions mocks base method.
func (m *MockFullNode) StateMinerPartitions(arg0 context.Context, arg1 address.Address, arg2 uint64, arg3 types.TipSetKey) ([]api.Partition, error) {
	return nil, nil
}

// StateMinerPower mocks base method.
func (m *MockFullNode) StateMinerPower(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (*api.MinerPower, error) {
	return nil, nil
}

// StateMinerPreCommitDepositForPower mocks base method.
func (m *MockFullNode) StateMinerPreCommitDepositForPower(arg0 context.Context, arg1 address.Address, arg2 miner0.SectorPreCommitInfo, arg3 types.TipSetKey) (big.Int, error) {
	return big.Zero(), nil
}

// StateMinerProvingDeadline mocks base method.
func (m *MockFullNode) StateMinerProvingDeadline(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (*dline.Info, error) {
	return &dline.Info{}, nil
}

// StateMinerRecoveries mocks base method.
func (m *MockFullNode) StateMinerRecoveries(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (bitfield.BitField, error) {
	return bitfield.New(), nil
}

// StateMinerSectorAllocated mocks base method.
func (m *MockFullNode) StateMinerSectorAllocated(arg0 context.Context, arg1 address.Address, arg2 abi.SectorNumber, arg3 types.TipSetKey) (bool, error) {
	return false, nil
}

// StateMinerSectorCount mocks base method.
func (m *MockFullNode) StateMinerSectorCount(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (api.MinerSectors, error) {
	return api.MinerSectors{}, nil
}

// StateMinerSectors mocks base method.
func (m *MockFullNode) StateMinerSectors(arg0 context.Context, arg1 address.Address, arg2 *bitfield.BitField, arg3 types.TipSetKey) ([]*miner.SectorOnChainInfo, error) {
	return nil, nil
}

// StateNetworkName mocks base method.
func (m *MockFullNode) StateNetworkName(arg0 context.Context) (dtypes.NetworkName, error) {
	return "", nil
}

// StateNetworkVersion mocks base method.
func (m *MockFullNode) StateNetworkVersion(arg0 context.Context, arg1 types.TipSetKey) (network.Version, error) {
	return network.Version10, nil
}

// StateReadState mocks base method.
func (m *MockFullNode) StateReadState(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (*api.ActorState, error) {
	return nil, nil
}

// StateReplay mocks base method.
func (m *MockFullNode) StateReplay(arg0 context.Context, arg1 types.TipSetKey, arg2 cid.Cid) (*api.InvocResult, error) {
	return nil, nil
}

// StateSearchMsg mocks base method.
func (m *MockFullNode) StateSearchMsg(arg0 context.Context, arg1 types.TipSetKey, arg2 cid.Cid, arg3 abi.ChainEpoch, arg4 bool) (*api.MsgLookup, error) {
	return nil, nil
}

// StateSectorExpiration mocks base method.
func (m *MockFullNode) StateSectorExpiration(arg0 context.Context, arg1 address.Address, arg2 abi.SectorNumber, arg3 types.TipSetKey) (*miner.SectorExpiration, error) {
	return nil, nil
}

// StateSectorGetInfo mocks base method.
func (m *MockFullNode) StateSectorGetInfo(arg0 context.Context, arg1 address.Address, arg2 abi.SectorNumber, arg3 types.TipSetKey) (*miner.SectorOnChainInfo, error) {
	return nil, nil
}

// StateSectorPartition mocks base method.
func (m *MockFullNode) StateSectorPartition(arg0 context.Context, arg1 address.Address, arg2 abi.SectorNumber, arg3 types.TipSetKey) (*miner.SectorLocation, error) {
	return nil, nil
}

// StateSectorPreCommitInfo mocks base method.
func (m *MockFullNode) StateSectorPreCommitInfo(arg0 context.Context, arg1 address.Address, arg2 abi.SectorNumber, arg3 types.TipSetKey) (miner.SectorPreCommitOnChainInfo, error) {
	return miner.SectorPreCommitOnChainInfo{
		Info:             miner.SectorPreCommitInfo{},
		PreCommitDeposit: big.Zero()}, nil
}

// StateVMCirculatingSupplyInternal mocks base method.
func (m *MockFullNode) StateVMCirculatingSupplyInternal(arg0 context.Context, arg1 types.TipSetKey) (api.CirculatingSupply, error) {
	return api.CirculatingSupply{}, nil
}

// StateVerifiedClientStatus mocks base method.
func (m *MockFullNode) StateVerifiedClientStatus(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (*big.Int, error) {
	return nil, nil
}

// StateVerifiedRegistryRootKey mocks base method.
func (m *MockFullNode) StateVerifiedRegistryRootKey(arg0 context.Context, arg1 types.TipSetKey) (address.Address, error) {
	return address.Undef, nil
}

// StateVerifierStatus mocks base method.
func (m *MockFullNode) StateVerifierStatus(arg0 context.Context, arg1 address.Address, arg2 types.TipSetKey) (*big.Int, error) {
	return nil, nil
}

// StateWaitMsg mocks base method.
func (m *MockFullNode) StateWaitMsg(arg0 context.Context, arg1 cid.Cid, arg2 uint64, arg3 abi.ChainEpoch, arg4 bool) (*api.MsgLookup, error) {
	return &api.MsgLookup{
		Receipt: types.MessageReceipt{
			ExitCode: 0,
		},
	}, nil
}

// SyncCheckBad mocks base method.
func (m *MockFullNode) SyncCheckBad(arg0 context.Context, arg1 cid.Cid) (string, error) {
	return "", nil
}

// SyncCheckpoint mocks base method.
func (m *MockFullNode) SyncCheckpoint(arg0 context.Context, arg1 types.TipSetKey) error {
	return nil
}

// SyncIncomingBlocks mocks base method.
func (m *MockFullNode) SyncIncomingBlocks(arg0 context.Context) (<-chan *types.BlockHeader, error) {
	return nil, nil
}

// SyncMarkBad mocks base method.
func (m *MockFullNode) SyncMarkBad(arg0 context.Context, arg1 cid.Cid) error {
	return nil
}

// SyncState mocks base method.
func (m *MockFullNode) SyncState(arg0 context.Context) (*api.SyncState, error) {
	return nil, nil
}

// SyncSubmitBlock mocks base method.
func (m *MockFullNode) SyncSubmitBlock(arg0 context.Context, arg1 *types.BlockMsg) error {
	return nil
}

// SyncUnmarkAllBad mocks base method.
func (m *MockFullNode) SyncUnmarkAllBad(arg0 context.Context) error {
	return nil
}

// SyncUnmarkBad mocks base method.
func (m *MockFullNode) SyncUnmarkBad(arg0 context.Context, arg1 cid.Cid) error {
	return nil
}

// SyncValidateTipset mocks base method.
func (m *MockFullNode) SyncValidateTipset(arg0 context.Context, arg1 types.TipSetKey) (bool, error) {
	return false, nil
}

// Version mocks base method.
func (m *MockFullNode) Version(arg0 context.Context) (api.APIVersion, error) {
	return api.APIVersion{}, nil
}

// WalletBalance mocks base method.
func (m *MockFullNode) WalletBalance(arg0 context.Context, arg1 address.Address) (big.Int, error) {
	return big.Zero(), nil
}

// WalletDefaultAddress mocks base method.
func (m *MockFullNode) WalletDefaultAddress(arg0 context.Context) (address.Address, error) {
	return address.Undef, nil
}

// WalletDelete mocks base method.
func (m *MockFullNode) WalletDelete(arg0 context.Context, arg1 address.Address) error {
	return nil
}

// WalletExport mocks base method.
func (m *MockFullNode) WalletExport(arg0 context.Context, arg1 address.Address) (*types.KeyInfo, error) {
	return nil, nil
}

// WalletHas mocks base method.
func (m *MockFullNode) WalletHas(arg0 context.Context, arg1 address.Address) (bool, error) {
	return true, nil
}

// WalletImport mocks base method.
func (m *MockFullNode) WalletImport(arg0 context.Context, arg1 *types.KeyInfo) (address.Address, error) {
	return address.Undef, nil
}

// WalletList mocks base method.
func (m *MockFullNode) WalletList(arg0 context.Context) ([]address.Address, error) {
	return []address.Address{address.Undef}, nil
}

// WalletNew mocks base method.
func (m *MockFullNode) WalletNew(arg0 context.Context, arg1 types.KeyType) (address.Address, error) {
	return address.Undef, nil
}

// WalletSetDefault mocks base method.
func (m *MockFullNode) WalletSetDefault(arg0 context.Context, arg1 address.Address) error {
	return nil
}

// WalletSign mocks base method.
func (m *MockFullNode) WalletSign(arg0 context.Context, arg1 address.Address, arg2 []byte) (*crypto.Signature, error) {
	return nil, nil
}

// WalletSignMessage mocks base method.
func (m *MockFullNode) WalletSignMessage(arg0 context.Context, arg1 address.Address, arg2 *types.Message) (*types.SignedMessage, error) {
	return nil, nil
}

// WalletValidateAddress mocks base method.
func (m *MockFullNode) WalletValidateAddress(arg0 context.Context, arg1 string) (address.Address, error) {
	return address.Undef, nil
}

// WalletVerify mocks base method.
func (m *MockFullNode) WalletVerify(arg0 context.Context, arg1 address.Address, arg2 []byte, arg3 *crypto.Signature) (bool, error) {
	return false, nil
}
