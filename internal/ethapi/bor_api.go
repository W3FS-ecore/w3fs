package ethapi

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	file_store "github.com/ethereum/go-ethereum/borcontracts/file-store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	w3fsAuth "github.com/ethereum/go-ethereum/storage/auth"
	"strconv"
	"time"
)

const RPC_JSON_PORT = 8545

// miner Flag
var MinerFlag = false
var isSealing int32

// GetRootHash returns root hash for given start and end block
func (s *PublicBlockChainAPI) GetRootHash(ctx context.Context, starBlockNr uint64, endBlockNr uint64) (string, error) {
	root, err := s.b.GetRootHash(ctx, starBlockNr, endBlockNr)
	if err != nil {
		return "", err
	}
	return root, nil
}

func (s *PublicBlockChainAPI) GetBorBlockReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return s.b.GetBorBlockReceipt(ctx, hash)
}

//
// Bor transaction utils
//

func (s *PublicBlockChainAPI) appendRPCMarshalBorTransaction(ctx context.Context, block *types.Block, fields map[string]interface{}, fullTx bool) map[string]interface{} {
	if block != nil {
		txHash := types.GetDerivedBorTxHash(types.BorReceiptKey(block.Number().Uint64(), block.Hash()))
		borTx, blockHash, blockNumber, txIndex, _ := s.b.GetBorBlockTransactionWithBlockHash(ctx, txHash, block.Hash())
		if borTx != nil {
			formattedTxs := fields["transactions"].([]interface{})
			if fullTx {
				marshalledTx := newRPCTransaction(borTx, blockHash, blockNumber, txIndex, nil)
				// newRPCTransaction calculates hash based on RLP of the transaction data.
				// In case of bor block tx, we need simple derived tx hash (same as function argument) instead of RLP hash
				marshalledTx.Hash = txHash
				fields["transactions"] = append(formattedTxs, marshalledTx)
			} else {
				fields["transactions"] = append(formattedTxs, txHash)
			}
		}
	}
	return fields
}

func changeMidByEthAddress(address common.Address) uint64 {
	encode := hexutil.Encode(address.Bytes())
	mid, _ := strconv.ParseUint(encode[2:6], 16, 32)
	return mid
}

/*func sealAndSend(addr common.Address, curHeader *types.Header, sectorSize string, privateKey *ecdsa.PrivateKey, chainId *big.Int, storageDir string) (common.Hash,error) {
	mid := changeMidByEthAddress(addr)
	currentBlockNumber := uint64(0)
	if curHeader != nil {
		currentBlockNumber = curHeader.Number.Uint64()
	}
	mid = 1000
	sectorNumberBig, _ := borcontracts.PowerCliObject.GetCaller().GetValidatorSectorInx(nil, addr)
	sectorNumber := sectorNumberBig.Uint64()
	seal := sealing.New(sectorNumber, mid, sectorSize, storageDir, currentBlockNumber, curHeader.Hash().Bytes())
	seal.AddPiece()
	seal.SealPreCommit1()
	seal.SealPreCommit2()
	seal.SealCommit1()
	seal.SealCommit2()
	ok, err := seal.VerifySeal()
	if !ok {
		lotusLog.Error(err)
		return common.Hash{}, err
	}
	sectorSizeInt, _ := units.RAMInBytes(sectorSize)
	sealProofType, _ := sealing.SealProofTypeFromSectorSize(fabi.SectorSize(sectorSizeInt), network.Version0)
	var cvotes []borcontracts.Cvote
	cvotes = append(cvotes, borcontracts.Cvote{
		SectorInx:     sectorNumber,
		SealProofType: uint64(sealProofType),
		SealedCID:     seal.Cids.Sealed.Bytes(),
		Proof:         seal.Proof,
	})
	cvotesBytes, _ := rlp.EncodeToBytes(cvotes)
	data, err := borcontracts.PowerCliObject.PackAddValidatorPowerAndProof(false ,addr, new(big.Int).SetUint64(uint64(sealProofType)), cvotesBytes)
	if err != nil {
		return common.Hash{}, errors.New("pack data is err")
	}
	auth, err := borcontracts.PowerCliObject.GenerateAuthObj(privateKey, chainId, addr, data)
	if err != nil {
		return common.Hash{}, err
	}
	txHash, err := borcontracts.PowerCliObject.AddValidatorPowerAndProofByClient(auth, false, addr, new(big.Int).SetUint64(uint64(sealProofType)), cvotesBytes)
	lotusLog.Info("SealAndSendTransaction tx :", txHash)
	if err != nil {
		return common.Hash{}, err
	}
	return txHash, nil
}*/

//Deprecated
func (s *PrivateAccountAPI) SealAndSendTransaction(ctx context.Context, password string) (string, error) {
	/*storageDir := s.b.GetDataDir() + "/.w3fsminer"
	closeTask := func() {
		atomic.StoreInt32(&isSealing, 0)
	}
	if atomic.LoadInt32(&isSealing) == 1 {
		return "" , errors.New("there is already a sealing mission in progress")
	}
	atomic.StoreInt32(&isSealing, 1)
	chainId := s.b.ChainConfig().ChainID
	ks, err := fetchKeystore(s.am)
	if err != nil {
		closeTask()
		return "", err
	}
	sectorSize := s.b.ChainConfig().Bor.ProofType
	coinbase, _ := s.b.Coinbase()
	privateKey, err := ks.GetAccountPrivateKey(accounts.Account{Address: coinbase}, password)
	if err != nil {
		closeTask()
		return "", err
	}
	curHeader, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber)
	// get all validators
	validators := s.b.Engine().GetCurrentAllValidators(curHeader.Hash(), curHeader.Number.Uint64())
	isValidator := false
	if isValidator {
		closeTask()
		return "", errors.New("validators is empty")
	}
	for _, value := range validators {
		if bytes.Compare(coinbase[:], value[:]) == 0 {
			isValidator = true
			break
		}
	}
	if isValidator == false {
		closeTask()
		return "", errors.New("it is not a validator")
	}
	go func() {
		defer closeTask()
		sealAndSend(coinbase, curHeader, sectorSize, privateKey, chainId, storageDir)
	}()
	return "the pledge start task execution .." , nil

	*/
	return "Deprecated function...", nil
}

func (s *PrivateAccountAPI) SetMinerInfo(ctx context.Context, password string, peerAddr string, proxyAddr string) (common.Hash, error) {
	// check peerId and peerAddr
	if PeerId == "" {
		return common.Hash{}, errors.New("error: miner's peerId is null.")
	}
	if proxyAddr == "" {
		return common.Hash{}, errors.New("error: proxyAddr is null.")
	}
	b := common.IsMultiAddr(proxyAddr)
	if !b {
		return common.Hash{}, errors.New("error: proxyAddr's format is wrong.")
	}
	// check is port = RPC_JSON_PORT
	b2 := common.CheckPort(proxyAddr, RPC_JSON_PORT)
	if !b2 {
		return common.Hash{}, errors.New("error: proxyAddr's port must be " + strconv.Itoa(RPC_JSON_PORT) + "! Please ensure value of param --http.port is " + strconv.Itoa(RPC_JSON_PORT))
	}
	if peerAddr == "" {
		return common.Hash{}, errors.New("error: miner's peerAddr is null.")
	}
	if !common.IsMultiAddr(peerAddr) {
		return common.Hash{}, errors.New("error: miner's peerAddr is not a Multiaddr.")
	}
	// check current is miner?
	if !MinerFlag {
		return common.Hash{}, errors.New("current node is not a miner,no need to setMinerInfo.")
	}
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return common.Hash{}, err
	}
	coinbase, _ := s.b.Coinbase()
	privateKey, err := ks.GetAccountPrivateKey(accounts.Account{Address: coinbase}, password)
	if err != nil {
		return common.Hash{}, err
	}
	localNode := s.b.GetP2pServer().LocalNode()
	minerId := common.HexSTrToByte32(localNode.ID().String())
	lotusLog.Info("localNode id:", localNode.ID().String(), "minerId(big.int):", minerId, "peerAddr:", peerAddr)
	publicKeyStr := w3fsAuth.TSK_GetPublicKey()
	lotusLog.Info("proxyAddr:", proxyAddr, " publicKey:", publicKeyStr)
	data, err := file_store.FileStoreCli.Pack4SetMinerInfo(minerId, publicKeyStr, PeerId, peerAddr, proxyAddr)
	if err != nil {
		return common.Hash{}, errors.New("pack data is err")
	}
	auth, err := file_store.FileStoreCli.GenerateAuthObj(privateKey, s.b.ChainConfig().ChainID, coinbase, data)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := file_store.FileStoreCli.SetMinerInfo(auth, minerId, publicKeyStr, PeerId, peerAddr, proxyAddr)
	if err != nil {
		return common.Hash{}, err
	}
	const SLEEP_TIME = 8 * time.Second
	time.Sleep(SLEEP_TIME)
	// after n seconds, query from chian,check whether set successful.
	retMinerIdStr := file_store.FileStoreCli.GetMinerId(nil, coinbase)
	lotusLog.Info("[1]retMinerIdStr:", retMinerIdStr, "minerAddr:", coinbase)
	if retMinerIdStr == "" || retMinerIdStr == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		time.Sleep(SLEEP_TIME)
		retMinerIdStr = file_store.FileStoreCli.GetMinerId(nil, coinbase)
		lotusLog.Info("[2]retMinerIdStr:", retMinerIdStr, "minerAddr:", coinbase)
		if retMinerIdStr == "" || retMinerIdStr == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			return common.Hash{}, errors.New("setMinerInfo error,please try again later!")
		}
	}
	return tx.Hash(), nil
}
