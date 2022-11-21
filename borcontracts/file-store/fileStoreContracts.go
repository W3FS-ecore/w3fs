package file_store

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"sync"
	"time"
)

var FileStoreContractAddr = "0x0000000000000000000000000000000000003002"
var FABI *abi.ABI

var transactor *FileStoreTransactor
var caller *FileStoreCaller
var FileStoreCli *FileStoreClient

const NOT_FOUND = 0
const VALID_PASS = 1
const VALID_NO_PASS = 2

const CHECK_FINISH = 1
const CHECK_NO_FINISH = 2

func init() {
	FABI, _ = FileStoreMetaData.GetAbi()
}

type FileStoreClient struct {
	transactor *FileStoreTransactor
	caller     *FileStoreCaller
	client     *ethclient.Client
	rmLock     sync.RWMutex
	lastTime   int64
	lastNonce  uint64
}

func NewFileStoreClient(client *ethclient.Client) *FileStoreClient {
	transactor, _ := NewFileStoreTransactor(common.HexToAddress(FileStoreContractAddr), client)
	caller, _ := NewFileStoreCaller(common.HexToAddress(FileStoreContractAddr), client)
	return &FileStoreClient{
		transactor: transactor,
		caller:     caller,
		client:     client,
		rmLock:     sync.RWMutex{},
		lastTime:   time.Now().UnixNano() / 1e6,
	}
}

func (fsc *FileStoreClient) GetCaller() *FileStoreCaller {
	return fsc.caller
}

func (fsc *FileStoreClient) GetTransactor() *FileStoreTransactor {
	return fsc.transactor
}

/**
  do not need invoke me when you just to query.
*/
func (pc *FileStoreClient) GenerateAuthObj(ecdsaPrivateKey *ecdsa.PrivateKey, chainID *big.Int, fromAddress common.Address, data []byte) (auth *bind.TransactOpts, err error) {
	toAddress := common.HexToAddress(FileStoreContractAddr)
	callMsg := ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	}
	gasprice, err := pc.client.SuggestGasPrice(context.Background())
	nonce, err := pc.getPendingNonce(fromAddress)
	callMsg.From = fromAddress
	//gasLimit, err := pc.client.EstimateGas(context.Background(), callMsg)
	gasLimit, err := pc.client.EstimateGasLatest(context.Background(), callMsg)
	//auth = bind.NewKeyedTransactor(ecdsaPrivateKey)
	auth, _ = bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
	auth.GasPrice = gasprice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = gasLimit * 2
	return
}

func (pc *FileStoreClient) getPendingNonce(fromAddress common.Address) (uint64, error) {
	pc.rmLock.Lock()
	defer pc.rmLock.Unlock()
	currentTimestamp := time.Now().UnixNano() / 1e6
	timeDiff := currentTimestamp - pc.lastTime
	var err error
	if pc.lastNonce != 0 && timeDiff < 50 {
		pc.lastNonce = pc.lastNonce + 1
	} else {
		pc.lastNonce, err = pc.client.PendingNonceAt(context.Background(), fromAddress)
	}
	pc.lastTime = currentTimestamp
	return pc.lastNonce, err
}
func (pc *FileStoreClient) Pack4UpdateFileStoreInfo4Entire(storeKeyHash [32]byte, cid string, status uint8) ([]byte, error) {
	data, err := FABI.Pack("updateFileStoreInfo4Entire", storeKeyHash, cid, status)
	return data, err
}

func (fsc *FileStoreClient) UpdateFileStoreInfo4Entire(opts *bind.TransactOpts, storeKeyHash [32]byte, cid string, status uint8) (*types.Transaction, error) {
	return fsc.transactor.UpdateFileStoreInfo4Entire(opts, storeKeyHash, cid, status)
}

func (pc *FileStoreClient) Pack4UpdateFileStoreInfo(oriHash [32]byte, headFlag bool, cid string, status uint8) ([]byte, error) {
	data, err := FABI.Pack("updateFileStoreInfo", oriHash, headFlag, cid, status)
	return data, err
}

// Solidity: function updateFileStoreInfo(uint256 oriHash, bool headFlag, string cid, uint8 status, uint256 operTime) returns()
func (fsc *FileStoreClient) UpdateFileStoreInfo(opts *bind.TransactOpts, oriHash [32]byte, headFlag bool, cid string, status uint8) (*types.Transaction, error) {
	return fsc.transactor.UpdateFileStoreInfo(opts, oriHash, headFlag, cid, status)
}

func (pc *FileStoreClient) Pack4SetMinerInfo(minerId [32]byte, publicKey string, peerId string, peerAddr string, proxyIp string) ([]byte, error) {
	data, err := FABI.Pack("setMinerInfo", minerId, publicKey, peerId, peerAddr, proxyIp)
	return data, err
}

// Solidity: function setMinerInfo(address minerAddr, uint256 minerId, string enode) returns()
func (fsc *FileStoreClient) SetMinerInfo(opts *bind.TransactOpts, minerId [32]byte, publicKey string, peerId string, peerAddr string, proxyIp string) (*types.Transaction, error) {
	return fsc.transactor.SetMinerInfo(opts, minerId, publicKey, peerId, peerAddr, proxyIp)
}

func (pc *FileStoreClient) Pack4WithdrawRemaining(oriHash [32]byte, index *big.Int, storageType uint8) ([]byte, error) {
	data, err := FABI.Pack("withdrawRemaining", oriHash, index, storageType)
	return data, err
}

// Solidity: function WithdrawRemaining() returns()
func (fsc *FileStoreClient) WithdrawRemaining(opts *bind.TransactOpts, oriHash [32]byte, index *big.Int, storageType uint8) (*types.Transaction, error) {
	return fsc.transactor.WithdrawRemaining(opts, oriHash, index, storageType)
}

func (fsc *FileStoreClient) GetExpireFileEntire(opts *bind.CallOpts) (FileStoreStructExpireFile, error) {
	return fsc.caller.GetExpireFileEntire(opts)
}

// Solidity: function findMiner4File(uint256 oriHash, bool headFlag) view returns((uint256,string,uint256))
func (fsc *FileStoreClient) GetExpireFile(opts *bind.CallOpts) (FileStoreStructExpireFile, error) {
	return fsc.caller.GetExpireFile(opts)
}

func (fsc *FileStoreClient) FindMiner4EntireFile(opts *bind.CallOpts, storeKeyHash [32]byte) (FileStoreStructFileMinerInfo, error) {
	return fsc.caller.FindMiner4EntireFile(opts, storeKeyHash)
}

// Solidity: function findMiner4File(uint256 oriHash, bool headFlag) view returns((uint256,string,uint256))
func (fsc *FileStoreClient) FindMiner4File(opts *bind.CallOpts, oriHash [32]byte, headFlag bool) (FileStoreStructFileMinerInfo, error) {
	return fsc.caller.FindMiner4File(opts, oriHash, headFlag)
}

func (fsc *FileStoreClient) GetBaseInfo4Entire(opts *bind.CallOpts, storeKeyHash [32]byte) (FileStoreStructBaseInfo, error) {
	return fsc.caller.GetBaseInfo4Entire(opts, storeKeyHash)
}

// Solidity: function getFileStoreBaseInfo(uint256 oriHash) view returns((uint256,address,uint32,string,uint256[],address[],uint256,uint256))
func (fsc *FileStoreClient) GetBaseInfo(opts *bind.CallOpts, oriHash [32]byte) (FileStoreStructBaseInfo, error) {
	return fsc.caller.GetBaseInfo(opts, oriHash)
}

// Solidity: function getMinerInfoByMinerId(uint256 minerId) view returns((string,string))
func (fsc *FileStoreClient) GetMinerInfoByMinerId(opts *bind.CallOpts, minerId [32]byte) FileStoreStructMinerInfo {
	minerInfo, err := fsc.caller.GetMinerInfoByMinerId(opts, minerId)
	if err != nil {
		return FileStoreStructMinerInfo{}
	}
	return minerInfo
}

func (fsc *FileStoreClient) GetMinerId(opts *bind.CallOpts, minerAddr common.Address) string {
	data, err := fsc.caller.GetMinerId(opts, minerAddr)
	if err != nil {
		return ""
	}
	return common.Byte32ToHexStr(data)
}

func (fsc *FileStoreClient) GetMinerAddr(opts *bind.CallOpts, minerId [32]byte) common.Address {
	data, err := fsc.caller.GetMinerAddr(opts, minerId)
	if err != nil {
		return data
	}
	return data
}

/**
  return int.
   0-not found.1- valid pass  2- no pass
*/
func (fsc *FileStoreClient) ValidFileInfo4Entire(storeKeyHash [32]byte, storeKey string, fileHashHex string) uint8 {
	baseInfo, err2 := fsc.caller.GetBaseInfo4Entire(nil, storeKeyHash)
	if err2 != nil || baseInfo.CDate.Cmp(big.NewInt(0)) == 0 {
		return NOT_FOUND
	}
	// caculte storeKey: user address (with 0x) + sha256 value (without 0x)
	hashAndAddr := baseInfo.OwnerAddr.String() + fileHashHex
	// check 2022-07-21 update:Change the comparison to all lower case because the address is sometimes all lower case and sometimes mixed case
	if strings.ToLower(hashAndAddr) == strings.ToLower(storeKey) {
		return VALID_PASS
	}
	return VALID_NO_PASS
}

/**
  return int.
   0-not found.1- valid pass  2- no pass
*/
func (fsc *FileStoreClient) ValidFileInfo(opts *bind.CallOpts, oriHash [32]byte, headFlag bool, minerId [32]byte, fileHash [32]byte) uint8 {
	result, err := fsc.caller.ValidFileInfo(opts, oriHash, headFlag, minerId, fileHash)
	if err != nil {
		return VALID_NO_PASS
	}
	return result
}

func (fsc *FileStoreClient) CheckStorage4Entire(opts *bind.CallOpts, storeKeyHash [32]byte, minerId [32]byte) uint8 {
	result, err := fsc.caller.CheckStorage4Entire(opts, storeKeyHash, minerId)
	if err != nil {
		return CHECK_NO_FINISH
	}
	return result
}

// Solidity: function checkStorage(uint256 oriHash, bool headFlag, uint256 minerId) view returns(uint8)
// @return int  0-not found  1-finish  2- unFinish
func (fsc *FileStoreClient) CheckStorage(opts *bind.CallOpts, oriHash [32]byte, headFlag bool, minerId [32]byte) uint8 {
	result, err := fsc.caller.CheckStorage(opts, oriHash, headFlag, minerId)
	if err != nil {
		return CHECK_NO_FINISH
	}
	return result
}

// GetEntireStoreMiners only for entire file storeage
func (fsc *FileStoreClient) GetStoreMiners4Entire(opts *bind.CallOpts, storeKeyHash [32]byte) ([][32]byte, error) {
	result, err := fsc.caller.GetStoreMiners4Entire(opts, storeKeyHash)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetStoreMiners is a free data retrieval call binding the contract method 0x82c5a1ca.
//
// Solidity: function getStoreMiners(uint256 oriHash) view returns(uint256[])
func (fsc *FileStoreClient) GetStoreMiners(opts *bind.CallOpts, oriHash [32]byte) ([][32]byte, error) {
	result, err := fsc.caller.GetStoreMiners(opts, oriHash)
	if err != nil {
		return nil, err
	}
	return result, nil
}
