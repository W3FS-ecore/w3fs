// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package file_store

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FileStoreStructBaseInfo is an auto generated low-level Go binding around an user-defined struct.
type FileStoreStructBaseInfo struct {
	OriHash           [32]byte
	OwnerAddr         common.Address
	FileSize          *big.Int
	FileExt           string
	Miners            [][32]byte
	DappContractAddrs []common.Address
	CDate             *big.Int
	MDate             *big.Int
}

// FileStoreStructExpireFile is an auto generated low-level Go binding around an user-defined struct.
type FileStoreStructExpireFile struct {
	OriHash     [32]byte
	Index       *big.Int
	StorageType uint8
}

// FileStoreStructFileInfo is an auto generated low-level Go binding around an user-defined struct.
type FileStoreStructFileInfo struct {
	HeadStatus     *big.Int
	HeadHash       [32]byte
	HeadCid        string
	BodyStatus     *big.Int
	BodyHash       [32]byte
	BodyCid        string
	FileCost       *big.Int
	FileCostStatus uint8
	CDate          *big.Int
	MDate          *big.Int
	EDate          *big.Int
}

// FileStoreStructFileMinerInfo is an auto generated low-level Go binding around an user-defined struct.
type FileStoreStructFileMinerInfo struct {
	FileHash [32]byte
	FileCid  string
	MinerIds [][32]byte
}

// FileStoreStructMinerInfo is an auto generated low-level Go binding around an user-defined struct.
type FileStoreStructMinerInfo struct {
	MinerAddr common.Address
	PublicKey string
	PeerId    string
	PeerAddr  string
	ProxyAddr string
}

// FileStoreMetaData contains all meta data concerning the FileStore contract.
var FileStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"headFlag\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"fileInfoChangeEvt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"headHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"bodyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"operTime\",\"type\":\"uint256\"}],\"name\":\"newFileStoreEvt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"ownerAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dappContractAddr\",\"type\":\"address\"}],\"name\":\"registerDappContractAddrEvt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"minerAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"}],\"name\":\"setMinerEvt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldOwnerAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOwnerAddr\",\"type\":\"address\"}],\"name\":\"transferOwnerEvt\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"headFlag\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"checkStorage\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"checkStorage4Entire\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"headHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyHash\",\"type\":\"bytes32\"}],\"name\":\"createFileStoreInfo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"createFileStoreInfo4Entire\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"extendFileDeadline\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"extendFileDeadlineEntire\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"}],\"name\":\"findMiner4EntireFile\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fileHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"fileCid\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"minerIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structFileStoreStruct.FileMinerInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"headFlag\",\"type\":\"bool\"}],\"name\":\"findMiner4File\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fileHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"fileCid\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"minerIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structFileStoreStruct.FileMinerInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"}],\"name\":\"getBaseInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"ownerAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"miners\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"dappContractAddrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"cDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mDate\",\"type\":\"uint256\"}],\"internalType\":\"structFileStoreStruct.BaseInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"}],\"name\":\"getBaseInfo4Entire\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"ownerAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"miners\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"dappContractAddrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"cDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mDate\",\"type\":\"uint256\"}],\"internalType\":\"structFileStoreStruct.BaseInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpireFile\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"enumFileStoreStruct.StorageType\",\"name\":\"storageType\",\"type\":\"uint8\"}],\"internalType\":\"structFileStoreStruct.ExpireFile\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpireFileEntire\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"enumFileStoreStruct.StorageType\",\"name\":\"storageType\",\"type\":\"uint8\"}],\"internalType\":\"structFileStoreStruct.ExpireFile\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"}],\"name\":\"getFileCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"getFileInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"headStatus\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"headHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"headCid\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"bodyStatus\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"bodyHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"bodyCid\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"fileCost\",\"type\":\"uint256\"},{\"internalType\":\"enumFileStoreStruct.Status\",\"name\":\"fileCostStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"cDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"eDate\",\"type\":\"uint256\"}],\"internalType\":\"structFileStoreStruct.FileInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"getFileInfo4Entire\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"headStatus\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"headHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"headCid\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"bodyStatus\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"bodyHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"bodyCid\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"fileCost\",\"type\":\"uint256\"},{\"internalType\":\"enumFileStoreStruct.Status\",\"name\":\"fileCostStatus\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"cDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"eDate\",\"type\":\"uint256\"}],\"internalType\":\"structFileStoreStruct.FileInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"getMinerAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minerAddr\",\"type\":\"address\"}],\"name\":\"getMinerId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"getMinerInfoByMinerId\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"minerAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"peerId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"peerAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"proxyAddr\",\"type\":\"string\"}],\"internalType\":\"structFileStoreStruct.MinerInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"getProxyAddrByMinerId\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"}],\"name\":\"getStoreMiners\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"}],\"name\":\"getStoreMiners4Entire\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_addrs\",\"type\":\"address[]\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"isFileExpire\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"}],\"name\":\"isFileExpireEntire\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLocked\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isEntireFile\",\"type\":\"bool\"}],\"name\":\"lockOrUnlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"dappContractAddr\",\"type\":\"address\"}],\"name\":\"regDappContractAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setFileStoreStorageAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"publicKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"peerId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"peerAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"proxyAddr\",\"type\":\"string\"}],\"name\":\"setMinerInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"ownerAddr\",\"type\":\"address\"}],\"name\":\"transferFileOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"headFlag\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"updateFileStoreInfo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"storeKeyHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"updateFileStoreInfo4Entire\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"updateWithdrawThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"headFlag\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"minerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fileHash\",\"type\":\"bytes32\"}],\"name\":\"validFileInfo\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"enumFileStoreStruct.StorageType\",\"name\":\"storageType\",\"type\":\"uint8\"}],\"name\":\"withdrawRemaining\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// FileStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use FileStoreMetaData.ABI instead.
var FileStoreABI = FileStoreMetaData.ABI

// FileStore is an auto generated Go binding around an Ethereum contract.
type FileStore struct {
	FileStoreCaller     // Read-only binding to the contract
	FileStoreTransactor // Write-only binding to the contract
	FileStoreFilterer   // Log filterer for contract events
}

// FileStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type FileStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FileStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FileStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FileStoreSession struct {
	Contract     *FileStore        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FileStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FileStoreCallerSession struct {
	Contract *FileStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// FileStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FileStoreTransactorSession struct {
	Contract     *FileStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// FileStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type FileStoreRaw struct {
	Contract *FileStore // Generic contract binding to access the raw methods on
}

// FileStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FileStoreCallerRaw struct {
	Contract *FileStoreCaller // Generic read-only contract binding to access the raw methods on
}

// FileStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FileStoreTransactorRaw struct {
	Contract *FileStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFileStore creates a new instance of FileStore, bound to a specific deployed contract.
func NewFileStore(address common.Address, backend bind.ContractBackend) (*FileStore, error) {
	contract, err := bindFileStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FileStore{FileStoreCaller: FileStoreCaller{contract: contract}, FileStoreTransactor: FileStoreTransactor{contract: contract}, FileStoreFilterer: FileStoreFilterer{contract: contract}}, nil
}

// NewFileStoreCaller creates a new read-only instance of FileStore, bound to a specific deployed contract.
func NewFileStoreCaller(address common.Address, caller bind.ContractCaller) (*FileStoreCaller, error) {
	contract, err := bindFileStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FileStoreCaller{contract: contract}, nil
}

// NewFileStoreTransactor creates a new write-only instance of FileStore, bound to a specific deployed contract.
func NewFileStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*FileStoreTransactor, error) {
	contract, err := bindFileStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FileStoreTransactor{contract: contract}, nil
}

// NewFileStoreFilterer creates a new log filterer instance of FileStore, bound to a specific deployed contract.
func NewFileStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*FileStoreFilterer, error) {
	contract, err := bindFileStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FileStoreFilterer{contract: contract}, nil
}

// bindFileStore binds a generic wrapper to an already deployed contract.
func bindFileStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FileStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileStore *FileStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FileStore.Contract.FileStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileStore *FileStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileStore.Contract.FileStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileStore *FileStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileStore.Contract.FileStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileStore *FileStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FileStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileStore *FileStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileStore *FileStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileStore.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FileStore *FileStoreCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FileStore *FileStoreSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _FileStore.Contract.DEFAULTADMINROLE(&_FileStore.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FileStore *FileStoreCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _FileStore.Contract.DEFAULTADMINROLE(&_FileStore.CallOpts)
}

// CheckStorage is a free data retrieval call binding the contract method 0xcbb3a92a.
//
// Solidity: function checkStorage(bytes32 oriHash, bool headFlag, bytes32 minerId) view returns(uint8)
func (_FileStore *FileStoreCaller) CheckStorage(opts *bind.CallOpts, oriHash [32]byte, headFlag bool, minerId [32]byte) (uint8, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "checkStorage", oriHash, headFlag, minerId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckStorage is a free data retrieval call binding the contract method 0xcbb3a92a.
//
// Solidity: function checkStorage(bytes32 oriHash, bool headFlag, bytes32 minerId) view returns(uint8)
func (_FileStore *FileStoreSession) CheckStorage(oriHash [32]byte, headFlag bool, minerId [32]byte) (uint8, error) {
	return _FileStore.Contract.CheckStorage(&_FileStore.CallOpts, oriHash, headFlag, minerId)
}

// CheckStorage is a free data retrieval call binding the contract method 0xcbb3a92a.
//
// Solidity: function checkStorage(bytes32 oriHash, bool headFlag, bytes32 minerId) view returns(uint8)
func (_FileStore *FileStoreCallerSession) CheckStorage(oriHash [32]byte, headFlag bool, minerId [32]byte) (uint8, error) {
	return _FileStore.Contract.CheckStorage(&_FileStore.CallOpts, oriHash, headFlag, minerId)
}

// CheckStorage4Entire is a free data retrieval call binding the contract method 0x5cdbdbbc.
//
// Solidity: function checkStorage4Entire(bytes32 storeKeyHash, bytes32 minerId) view returns(uint8)
func (_FileStore *FileStoreCaller) CheckStorage4Entire(opts *bind.CallOpts, storeKeyHash [32]byte, minerId [32]byte) (uint8, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "checkStorage4Entire", storeKeyHash, minerId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckStorage4Entire is a free data retrieval call binding the contract method 0x5cdbdbbc.
//
// Solidity: function checkStorage4Entire(bytes32 storeKeyHash, bytes32 minerId) view returns(uint8)
func (_FileStore *FileStoreSession) CheckStorage4Entire(storeKeyHash [32]byte, minerId [32]byte) (uint8, error) {
	return _FileStore.Contract.CheckStorage4Entire(&_FileStore.CallOpts, storeKeyHash, minerId)
}

// CheckStorage4Entire is a free data retrieval call binding the contract method 0x5cdbdbbc.
//
// Solidity: function checkStorage4Entire(bytes32 storeKeyHash, bytes32 minerId) view returns(uint8)
func (_FileStore *FileStoreCallerSession) CheckStorage4Entire(storeKeyHash [32]byte, minerId [32]byte) (uint8, error) {
	return _FileStore.Contract.CheckStorage4Entire(&_FileStore.CallOpts, storeKeyHash, minerId)
}

// FindMiner4EntireFile is a free data retrieval call binding the contract method 0x7521f962.
//
// Solidity: function findMiner4EntireFile(bytes32 storeKeyHash) view returns((bytes32,string,bytes32[]))
func (_FileStore *FileStoreCaller) FindMiner4EntireFile(opts *bind.CallOpts, storeKeyHash [32]byte) (FileStoreStructFileMinerInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "findMiner4EntireFile", storeKeyHash)

	if err != nil {
		return *new(FileStoreStructFileMinerInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructFileMinerInfo)).(*FileStoreStructFileMinerInfo)

	return out0, err

}

// FindMiner4EntireFile is a free data retrieval call binding the contract method 0x7521f962.
//
// Solidity: function findMiner4EntireFile(bytes32 storeKeyHash) view returns((bytes32,string,bytes32[]))
func (_FileStore *FileStoreSession) FindMiner4EntireFile(storeKeyHash [32]byte) (FileStoreStructFileMinerInfo, error) {
	return _FileStore.Contract.FindMiner4EntireFile(&_FileStore.CallOpts, storeKeyHash)
}

// FindMiner4EntireFile is a free data retrieval call binding the contract method 0x7521f962.
//
// Solidity: function findMiner4EntireFile(bytes32 storeKeyHash) view returns((bytes32,string,bytes32[]))
func (_FileStore *FileStoreCallerSession) FindMiner4EntireFile(storeKeyHash [32]byte) (FileStoreStructFileMinerInfo, error) {
	return _FileStore.Contract.FindMiner4EntireFile(&_FileStore.CallOpts, storeKeyHash)
}

// FindMiner4File is a free data retrieval call binding the contract method 0x0867ca08.
//
// Solidity: function findMiner4File(bytes32 oriHash, bool headFlag) view returns((bytes32,string,bytes32[]))
func (_FileStore *FileStoreCaller) FindMiner4File(opts *bind.CallOpts, oriHash [32]byte, headFlag bool) (FileStoreStructFileMinerInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "findMiner4File", oriHash, headFlag)

	if err != nil {
		return *new(FileStoreStructFileMinerInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructFileMinerInfo)).(*FileStoreStructFileMinerInfo)

	return out0, err

}

// FindMiner4File is a free data retrieval call binding the contract method 0x0867ca08.
//
// Solidity: function findMiner4File(bytes32 oriHash, bool headFlag) view returns((bytes32,string,bytes32[]))
func (_FileStore *FileStoreSession) FindMiner4File(oriHash [32]byte, headFlag bool) (FileStoreStructFileMinerInfo, error) {
	return _FileStore.Contract.FindMiner4File(&_FileStore.CallOpts, oriHash, headFlag)
}

// FindMiner4File is a free data retrieval call binding the contract method 0x0867ca08.
//
// Solidity: function findMiner4File(bytes32 oriHash, bool headFlag) view returns((bytes32,string,bytes32[]))
func (_FileStore *FileStoreCallerSession) FindMiner4File(oriHash [32]byte, headFlag bool) (FileStoreStructFileMinerInfo, error) {
	return _FileStore.Contract.FindMiner4File(&_FileStore.CallOpts, oriHash, headFlag)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_FileStore *FileStoreCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_FileStore *FileStoreSession) GetBalance() (*big.Int, error) {
	return _FileStore.Contract.GetBalance(&_FileStore.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_FileStore *FileStoreCallerSession) GetBalance() (*big.Int, error) {
	return _FileStore.Contract.GetBalance(&_FileStore.CallOpts)
}

// GetBaseInfo is a free data retrieval call binding the contract method 0x23847afb.
//
// Solidity: function getBaseInfo(bytes32 oriHash) view returns((bytes32,address,uint256,string,bytes32[],address[],uint256,uint256))
func (_FileStore *FileStoreCaller) GetBaseInfo(opts *bind.CallOpts, oriHash [32]byte) (FileStoreStructBaseInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getBaseInfo", oriHash)

	if err != nil {
		return *new(FileStoreStructBaseInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructBaseInfo)).(*FileStoreStructBaseInfo)

	return out0, err

}

// GetBaseInfo is a free data retrieval call binding the contract method 0x23847afb.
//
// Solidity: function getBaseInfo(bytes32 oriHash) view returns((bytes32,address,uint256,string,bytes32[],address[],uint256,uint256))
func (_FileStore *FileStoreSession) GetBaseInfo(oriHash [32]byte) (FileStoreStructBaseInfo, error) {
	return _FileStore.Contract.GetBaseInfo(&_FileStore.CallOpts, oriHash)
}

// GetBaseInfo is a free data retrieval call binding the contract method 0x23847afb.
//
// Solidity: function getBaseInfo(bytes32 oriHash) view returns((bytes32,address,uint256,string,bytes32[],address[],uint256,uint256))
func (_FileStore *FileStoreCallerSession) GetBaseInfo(oriHash [32]byte) (FileStoreStructBaseInfo, error) {
	return _FileStore.Contract.GetBaseInfo(&_FileStore.CallOpts, oriHash)
}

// GetBaseInfo4Entire is a free data retrieval call binding the contract method 0x98b05cb2.
//
// Solidity: function getBaseInfo4Entire(bytes32 storeKeyHash) view returns((bytes32,address,uint256,string,bytes32[],address[],uint256,uint256))
func (_FileStore *FileStoreCaller) GetBaseInfo4Entire(opts *bind.CallOpts, storeKeyHash [32]byte) (FileStoreStructBaseInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getBaseInfo4Entire", storeKeyHash)

	if err != nil {
		return *new(FileStoreStructBaseInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructBaseInfo)).(*FileStoreStructBaseInfo)

	return out0, err

}

// GetBaseInfo4Entire is a free data retrieval call binding the contract method 0x98b05cb2.
//
// Solidity: function getBaseInfo4Entire(bytes32 storeKeyHash) view returns((bytes32,address,uint256,string,bytes32[],address[],uint256,uint256))
func (_FileStore *FileStoreSession) GetBaseInfo4Entire(storeKeyHash [32]byte) (FileStoreStructBaseInfo, error) {
	return _FileStore.Contract.GetBaseInfo4Entire(&_FileStore.CallOpts, storeKeyHash)
}

// GetBaseInfo4Entire is a free data retrieval call binding the contract method 0x98b05cb2.
//
// Solidity: function getBaseInfo4Entire(bytes32 storeKeyHash) view returns((bytes32,address,uint256,string,bytes32[],address[],uint256,uint256))
func (_FileStore *FileStoreCallerSession) GetBaseInfo4Entire(storeKeyHash [32]byte) (FileStoreStructBaseInfo, error) {
	return _FileStore.Contract.GetBaseInfo4Entire(&_FileStore.CallOpts, storeKeyHash)
}

// GetExpireFile is a free data retrieval call binding the contract method 0xa6e54f46.
//
// Solidity: function getExpireFile() view returns((bytes32,uint256,uint8))
func (_FileStore *FileStoreCaller) GetExpireFile(opts *bind.CallOpts) (FileStoreStructExpireFile, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getExpireFile")

	if err != nil {
		return *new(FileStoreStructExpireFile), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructExpireFile)).(*FileStoreStructExpireFile)

	return out0, err

}

// GetExpireFile is a free data retrieval call binding the contract method 0xa6e54f46.
//
// Solidity: function getExpireFile() view returns((bytes32,uint256,uint8))
func (_FileStore *FileStoreSession) GetExpireFile() (FileStoreStructExpireFile, error) {
	return _FileStore.Contract.GetExpireFile(&_FileStore.CallOpts)
}

// GetExpireFile is a free data retrieval call binding the contract method 0xa6e54f46.
//
// Solidity: function getExpireFile() view returns((bytes32,uint256,uint8))
func (_FileStore *FileStoreCallerSession) GetExpireFile() (FileStoreStructExpireFile, error) {
	return _FileStore.Contract.GetExpireFile(&_FileStore.CallOpts)
}

// GetExpireFileEntire is a free data retrieval call binding the contract method 0x9de81bf2.
//
// Solidity: function getExpireFileEntire() view returns((bytes32,uint256,uint8))
func (_FileStore *FileStoreCaller) GetExpireFileEntire(opts *bind.CallOpts) (FileStoreStructExpireFile, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getExpireFileEntire")

	if err != nil {
		return *new(FileStoreStructExpireFile), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructExpireFile)).(*FileStoreStructExpireFile)

	return out0, err

}

// GetExpireFileEntire is a free data retrieval call binding the contract method 0x9de81bf2.
//
// Solidity: function getExpireFileEntire() view returns((bytes32,uint256,uint8))
func (_FileStore *FileStoreSession) GetExpireFileEntire() (FileStoreStructExpireFile, error) {
	return _FileStore.Contract.GetExpireFileEntire(&_FileStore.CallOpts)
}

// GetExpireFileEntire is a free data retrieval call binding the contract method 0x9de81bf2.
//
// Solidity: function getExpireFileEntire() view returns((bytes32,uint256,uint8))
func (_FileStore *FileStoreCallerSession) GetExpireFileEntire() (FileStoreStructExpireFile, error) {
	return _FileStore.Contract.GetExpireFileEntire(&_FileStore.CallOpts)
}

// GetFileCost is a free data retrieval call binding the contract method 0x642164f3.
//
// Solidity: function getFileCost(uint256 fileSize) pure returns(uint256)
func (_FileStore *FileStoreCaller) GetFileCost(opts *bind.CallOpts, fileSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getFileCost", fileSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFileCost is a free data retrieval call binding the contract method 0x642164f3.
//
// Solidity: function getFileCost(uint256 fileSize) pure returns(uint256)
func (_FileStore *FileStoreSession) GetFileCost(fileSize *big.Int) (*big.Int, error) {
	return _FileStore.Contract.GetFileCost(&_FileStore.CallOpts, fileSize)
}

// GetFileCost is a free data retrieval call binding the contract method 0x642164f3.
//
// Solidity: function getFileCost(uint256 fileSize) pure returns(uint256)
func (_FileStore *FileStoreCallerSession) GetFileCost(fileSize *big.Int) (*big.Int, error) {
	return _FileStore.Contract.GetFileCost(&_FileStore.CallOpts, fileSize)
}

// GetFileInfo is a free data retrieval call binding the contract method 0x7733d418.
//
// Solidity: function getFileInfo(bytes32 oriHash, bytes32 minerId) view returns((uint256,bytes32,string,uint256,bytes32,string,uint256,uint8,uint256,uint256,uint256))
func (_FileStore *FileStoreCaller) GetFileInfo(opts *bind.CallOpts, oriHash [32]byte, minerId [32]byte) (FileStoreStructFileInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getFileInfo", oriHash, minerId)

	if err != nil {
		return *new(FileStoreStructFileInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructFileInfo)).(*FileStoreStructFileInfo)

	return out0, err

}

// GetFileInfo is a free data retrieval call binding the contract method 0x7733d418.
//
// Solidity: function getFileInfo(bytes32 oriHash, bytes32 minerId) view returns((uint256,bytes32,string,uint256,bytes32,string,uint256,uint8,uint256,uint256,uint256))
func (_FileStore *FileStoreSession) GetFileInfo(oriHash [32]byte, minerId [32]byte) (FileStoreStructFileInfo, error) {
	return _FileStore.Contract.GetFileInfo(&_FileStore.CallOpts, oriHash, minerId)
}

// GetFileInfo is a free data retrieval call binding the contract method 0x7733d418.
//
// Solidity: function getFileInfo(bytes32 oriHash, bytes32 minerId) view returns((uint256,bytes32,string,uint256,bytes32,string,uint256,uint8,uint256,uint256,uint256))
func (_FileStore *FileStoreCallerSession) GetFileInfo(oriHash [32]byte, minerId [32]byte) (FileStoreStructFileInfo, error) {
	return _FileStore.Contract.GetFileInfo(&_FileStore.CallOpts, oriHash, minerId)
}

// GetFileInfo4Entire is a free data retrieval call binding the contract method 0x7d7a51fb.
//
// Solidity: function getFileInfo4Entire(bytes32 storeKeyHash, bytes32 minerId) view returns((uint256,bytes32,string,uint256,bytes32,string,uint256,uint8,uint256,uint256,uint256))
func (_FileStore *FileStoreCaller) GetFileInfo4Entire(opts *bind.CallOpts, storeKeyHash [32]byte, minerId [32]byte) (FileStoreStructFileInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getFileInfo4Entire", storeKeyHash, minerId)

	if err != nil {
		return *new(FileStoreStructFileInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructFileInfo)).(*FileStoreStructFileInfo)

	return out0, err

}

// GetFileInfo4Entire is a free data retrieval call binding the contract method 0x7d7a51fb.
//
// Solidity: function getFileInfo4Entire(bytes32 storeKeyHash, bytes32 minerId) view returns((uint256,bytes32,string,uint256,bytes32,string,uint256,uint8,uint256,uint256,uint256))
func (_FileStore *FileStoreSession) GetFileInfo4Entire(storeKeyHash [32]byte, minerId [32]byte) (FileStoreStructFileInfo, error) {
	return _FileStore.Contract.GetFileInfo4Entire(&_FileStore.CallOpts, storeKeyHash, minerId)
}

// GetFileInfo4Entire is a free data retrieval call binding the contract method 0x7d7a51fb.
//
// Solidity: function getFileInfo4Entire(bytes32 storeKeyHash, bytes32 minerId) view returns((uint256,bytes32,string,uint256,bytes32,string,uint256,uint8,uint256,uint256,uint256))
func (_FileStore *FileStoreCallerSession) GetFileInfo4Entire(storeKeyHash [32]byte, minerId [32]byte) (FileStoreStructFileInfo, error) {
	return _FileStore.Contract.GetFileInfo4Entire(&_FileStore.CallOpts, storeKeyHash, minerId)
}

// GetMinerAddr is a free data retrieval call binding the contract method 0xb9a8fbf7.
//
// Solidity: function getMinerAddr(bytes32 minerId) view returns(address)
func (_FileStore *FileStoreCaller) GetMinerAddr(opts *bind.CallOpts, minerId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getMinerAddr", minerId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMinerAddr is a free data retrieval call binding the contract method 0xb9a8fbf7.
//
// Solidity: function getMinerAddr(bytes32 minerId) view returns(address)
func (_FileStore *FileStoreSession) GetMinerAddr(minerId [32]byte) (common.Address, error) {
	return _FileStore.Contract.GetMinerAddr(&_FileStore.CallOpts, minerId)
}

// GetMinerAddr is a free data retrieval call binding the contract method 0xb9a8fbf7.
//
// Solidity: function getMinerAddr(bytes32 minerId) view returns(address)
func (_FileStore *FileStoreCallerSession) GetMinerAddr(minerId [32]byte) (common.Address, error) {
	return _FileStore.Contract.GetMinerAddr(&_FileStore.CallOpts, minerId)
}

// GetMinerId is a free data retrieval call binding the contract method 0xe2dea715.
//
// Solidity: function getMinerId(address minerAddr) view returns(bytes32)
func (_FileStore *FileStoreCaller) GetMinerId(opts *bind.CallOpts, minerAddr common.Address) ([32]byte, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getMinerId", minerAddr)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetMinerId is a free data retrieval call binding the contract method 0xe2dea715.
//
// Solidity: function getMinerId(address minerAddr) view returns(bytes32)
func (_FileStore *FileStoreSession) GetMinerId(minerAddr common.Address) ([32]byte, error) {
	return _FileStore.Contract.GetMinerId(&_FileStore.CallOpts, minerAddr)
}

// GetMinerId is a free data retrieval call binding the contract method 0xe2dea715.
//
// Solidity: function getMinerId(address minerAddr) view returns(bytes32)
func (_FileStore *FileStoreCallerSession) GetMinerId(minerAddr common.Address) ([32]byte, error) {
	return _FileStore.Contract.GetMinerId(&_FileStore.CallOpts, minerAddr)
}

// GetMinerInfoByMinerId is a free data retrieval call binding the contract method 0x028092cb.
//
// Solidity: function getMinerInfoByMinerId(bytes32 minerId) view returns((address,string,string,string,string))
func (_FileStore *FileStoreCaller) GetMinerInfoByMinerId(opts *bind.CallOpts, minerId [32]byte) (FileStoreStructMinerInfo, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getMinerInfoByMinerId", minerId)

	if err != nil {
		return *new(FileStoreStructMinerInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(FileStoreStructMinerInfo)).(*FileStoreStructMinerInfo)

	return out0, err

}

// GetMinerInfoByMinerId is a free data retrieval call binding the contract method 0x028092cb.
//
// Solidity: function getMinerInfoByMinerId(bytes32 minerId) view returns((address,string,string,string,string))
func (_FileStore *FileStoreSession) GetMinerInfoByMinerId(minerId [32]byte) (FileStoreStructMinerInfo, error) {
	return _FileStore.Contract.GetMinerInfoByMinerId(&_FileStore.CallOpts, minerId)
}

// GetMinerInfoByMinerId is a free data retrieval call binding the contract method 0x028092cb.
//
// Solidity: function getMinerInfoByMinerId(bytes32 minerId) view returns((address,string,string,string,string))
func (_FileStore *FileStoreCallerSession) GetMinerInfoByMinerId(minerId [32]byte) (FileStoreStructMinerInfo, error) {
	return _FileStore.Contract.GetMinerInfoByMinerId(&_FileStore.CallOpts, minerId)
}

// GetProxyAddrByMinerId is a free data retrieval call binding the contract method 0x928441da.
//
// Solidity: function getProxyAddrByMinerId(bytes32 minerId) view returns(string)
func (_FileStore *FileStoreCaller) GetProxyAddrByMinerId(opts *bind.CallOpts, minerId [32]byte) (string, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getProxyAddrByMinerId", minerId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetProxyAddrByMinerId is a free data retrieval call binding the contract method 0x928441da.
//
// Solidity: function getProxyAddrByMinerId(bytes32 minerId) view returns(string)
func (_FileStore *FileStoreSession) GetProxyAddrByMinerId(minerId [32]byte) (string, error) {
	return _FileStore.Contract.GetProxyAddrByMinerId(&_FileStore.CallOpts, minerId)
}

// GetProxyAddrByMinerId is a free data retrieval call binding the contract method 0x928441da.
//
// Solidity: function getProxyAddrByMinerId(bytes32 minerId) view returns(string)
func (_FileStore *FileStoreCallerSession) GetProxyAddrByMinerId(minerId [32]byte) (string, error) {
	return _FileStore.Contract.GetProxyAddrByMinerId(&_FileStore.CallOpts, minerId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FileStore *FileStoreCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FileStore *FileStoreSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _FileStore.Contract.GetRoleAdmin(&_FileStore.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FileStore *FileStoreCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _FileStore.Contract.GetRoleAdmin(&_FileStore.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_FileStore *FileStoreCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_FileStore *FileStoreSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _FileStore.Contract.GetRoleMember(&_FileStore.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_FileStore *FileStoreCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _FileStore.Contract.GetRoleMember(&_FileStore.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_FileStore *FileStoreCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_FileStore *FileStoreSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _FileStore.Contract.GetRoleMemberCount(&_FileStore.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_FileStore *FileStoreCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _FileStore.Contract.GetRoleMemberCount(&_FileStore.CallOpts, role)
}

// GetStoreMiners is a free data retrieval call binding the contract method 0x802af5f0.
//
// Solidity: function getStoreMiners(bytes32 oriHash) view returns(bytes32[])
func (_FileStore *FileStoreCaller) GetStoreMiners(opts *bind.CallOpts, oriHash [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getStoreMiners", oriHash)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetStoreMiners is a free data retrieval call binding the contract method 0x802af5f0.
//
// Solidity: function getStoreMiners(bytes32 oriHash) view returns(bytes32[])
func (_FileStore *FileStoreSession) GetStoreMiners(oriHash [32]byte) ([][32]byte, error) {
	return _FileStore.Contract.GetStoreMiners(&_FileStore.CallOpts, oriHash)
}

// GetStoreMiners is a free data retrieval call binding the contract method 0x802af5f0.
//
// Solidity: function getStoreMiners(bytes32 oriHash) view returns(bytes32[])
func (_FileStore *FileStoreCallerSession) GetStoreMiners(oriHash [32]byte) ([][32]byte, error) {
	return _FileStore.Contract.GetStoreMiners(&_FileStore.CallOpts, oriHash)
}

// GetStoreMiners4Entire is a free data retrieval call binding the contract method 0xbd816ace.
//
// Solidity: function getStoreMiners4Entire(bytes32 storeKeyHash) view returns(bytes32[])
func (_FileStore *FileStoreCaller) GetStoreMiners4Entire(opts *bind.CallOpts, storeKeyHash [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "getStoreMiners4Entire", storeKeyHash)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetStoreMiners4Entire is a free data retrieval call binding the contract method 0xbd816ace.
//
// Solidity: function getStoreMiners4Entire(bytes32 storeKeyHash) view returns(bytes32[])
func (_FileStore *FileStoreSession) GetStoreMiners4Entire(storeKeyHash [32]byte) ([][32]byte, error) {
	return _FileStore.Contract.GetStoreMiners4Entire(&_FileStore.CallOpts, storeKeyHash)
}

// GetStoreMiners4Entire is a free data retrieval call binding the contract method 0xbd816ace.
//
// Solidity: function getStoreMiners4Entire(bytes32 storeKeyHash) view returns(bytes32[])
func (_FileStore *FileStoreCallerSession) GetStoreMiners4Entire(storeKeyHash [32]byte) ([][32]byte, error) {
	return _FileStore.Contract.GetStoreMiners4Entire(&_FileStore.CallOpts, storeKeyHash)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FileStore *FileStoreCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FileStore *FileStoreSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _FileStore.Contract.HasRole(&_FileStore.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FileStore *FileStoreCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _FileStore.Contract.HasRole(&_FileStore.CallOpts, role, account)
}

// IsFileExpire is a free data retrieval call binding the contract method 0x4b7cc556.
//
// Solidity: function isFileExpire(bytes32 oriHash, bytes32 minerId) view returns(bool)
func (_FileStore *FileStoreCaller) IsFileExpire(opts *bind.CallOpts, oriHash [32]byte, minerId [32]byte) (bool, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "isFileExpire", oriHash, minerId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFileExpire is a free data retrieval call binding the contract method 0x4b7cc556.
//
// Solidity: function isFileExpire(bytes32 oriHash, bytes32 minerId) view returns(bool)
func (_FileStore *FileStoreSession) IsFileExpire(oriHash [32]byte, minerId [32]byte) (bool, error) {
	return _FileStore.Contract.IsFileExpire(&_FileStore.CallOpts, oriHash, minerId)
}

// IsFileExpire is a free data retrieval call binding the contract method 0x4b7cc556.
//
// Solidity: function isFileExpire(bytes32 oriHash, bytes32 minerId) view returns(bool)
func (_FileStore *FileStoreCallerSession) IsFileExpire(oriHash [32]byte, minerId [32]byte) (bool, error) {
	return _FileStore.Contract.IsFileExpire(&_FileStore.CallOpts, oriHash, minerId)
}

// IsFileExpireEntire is a free data retrieval call binding the contract method 0x3c0bf4ae.
//
// Solidity: function isFileExpireEntire(bytes32 storeKeyHash, bytes32 minerId) view returns(bool)
func (_FileStore *FileStoreCaller) IsFileExpireEntire(opts *bind.CallOpts, storeKeyHash [32]byte, minerId [32]byte) (bool, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "isFileExpireEntire", storeKeyHash, minerId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFileExpireEntire is a free data retrieval call binding the contract method 0x3c0bf4ae.
//
// Solidity: function isFileExpireEntire(bytes32 storeKeyHash, bytes32 minerId) view returns(bool)
func (_FileStore *FileStoreSession) IsFileExpireEntire(storeKeyHash [32]byte, minerId [32]byte) (bool, error) {
	return _FileStore.Contract.IsFileExpireEntire(&_FileStore.CallOpts, storeKeyHash, minerId)
}

// IsFileExpireEntire is a free data retrieval call binding the contract method 0x3c0bf4ae.
//
// Solidity: function isFileExpireEntire(bytes32 storeKeyHash, bytes32 minerId) view returns(bool)
func (_FileStore *FileStoreCallerSession) IsFileExpireEntire(storeKeyHash [32]byte, minerId [32]byte) (bool, error) {
	return _FileStore.Contract.IsFileExpireEntire(&_FileStore.CallOpts, storeKeyHash, minerId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_FileStore *FileStoreCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_FileStore *FileStoreSession) Paused() (bool, error) {
	return _FileStore.Contract.Paused(&_FileStore.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_FileStore *FileStoreCallerSession) Paused() (bool, error) {
	return _FileStore.Contract.Paused(&_FileStore.CallOpts)
}

// ValidFileInfo is a free data retrieval call binding the contract method 0x768fe148.
//
// Solidity: function validFileInfo(bytes32 oriHash, bool headFlag, bytes32 minerId, bytes32 fileHash) view returns(uint8)
func (_FileStore *FileStoreCaller) ValidFileInfo(opts *bind.CallOpts, oriHash [32]byte, headFlag bool, minerId [32]byte, fileHash [32]byte) (uint8, error) {
	var out []interface{}
	err := _FileStore.contract.Call(opts, &out, "validFileInfo", oriHash, headFlag, minerId, fileHash)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// ValidFileInfo is a free data retrieval call binding the contract method 0x768fe148.
//
// Solidity: function validFileInfo(bytes32 oriHash, bool headFlag, bytes32 minerId, bytes32 fileHash) view returns(uint8)
func (_FileStore *FileStoreSession) ValidFileInfo(oriHash [32]byte, headFlag bool, minerId [32]byte, fileHash [32]byte) (uint8, error) {
	return _FileStore.Contract.ValidFileInfo(&_FileStore.CallOpts, oriHash, headFlag, minerId, fileHash)
}

// ValidFileInfo is a free data retrieval call binding the contract method 0x768fe148.
//
// Solidity: function validFileInfo(bytes32 oriHash, bool headFlag, bytes32 minerId, bytes32 fileHash) view returns(uint8)
func (_FileStore *FileStoreCallerSession) ValidFileInfo(oriHash [32]byte, headFlag bool, minerId [32]byte, fileHash [32]byte) (uint8, error) {
	return _FileStore.Contract.ValidFileInfo(&_FileStore.CallOpts, oriHash, headFlag, minerId, fileHash)
}

// CreateFileStoreInfo is a paid mutator transaction binding the contract method 0x4bacf4a4.
//
// Solidity: function createFileStoreInfo(bytes32 oriHash, uint256 fileSize, string fileExt, bytes32 minerId, bytes32 headHash, bytes32 bodyHash) payable returns()
func (_FileStore *FileStoreTransactor) CreateFileStoreInfo(opts *bind.TransactOpts, oriHash [32]byte, fileSize *big.Int, fileExt string, minerId [32]byte, headHash [32]byte, bodyHash [32]byte) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "createFileStoreInfo", oriHash, fileSize, fileExt, minerId, headHash, bodyHash)
}

// CreateFileStoreInfo is a paid mutator transaction binding the contract method 0x4bacf4a4.
//
// Solidity: function createFileStoreInfo(bytes32 oriHash, uint256 fileSize, string fileExt, bytes32 minerId, bytes32 headHash, bytes32 bodyHash) payable returns()
func (_FileStore *FileStoreSession) CreateFileStoreInfo(oriHash [32]byte, fileSize *big.Int, fileExt string, minerId [32]byte, headHash [32]byte, bodyHash [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.CreateFileStoreInfo(&_FileStore.TransactOpts, oriHash, fileSize, fileExt, minerId, headHash, bodyHash)
}

// CreateFileStoreInfo is a paid mutator transaction binding the contract method 0x4bacf4a4.
//
// Solidity: function createFileStoreInfo(bytes32 oriHash, uint256 fileSize, string fileExt, bytes32 minerId, bytes32 headHash, bytes32 bodyHash) payable returns()
func (_FileStore *FileStoreTransactorSession) CreateFileStoreInfo(oriHash [32]byte, fileSize *big.Int, fileExt string, minerId [32]byte, headHash [32]byte, bodyHash [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.CreateFileStoreInfo(&_FileStore.TransactOpts, oriHash, fileSize, fileExt, minerId, headHash, bodyHash)
}

// CreateFileStoreInfo4Entire is a paid mutator transaction binding the contract method 0xecbc1191.
//
// Solidity: function createFileStoreInfo4Entire(bytes32 storeKeyHash, uint256 fileSize, string fileExt, bytes32 minerId) payable returns()
func (_FileStore *FileStoreTransactor) CreateFileStoreInfo4Entire(opts *bind.TransactOpts, storeKeyHash [32]byte, fileSize *big.Int, fileExt string, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "createFileStoreInfo4Entire", storeKeyHash, fileSize, fileExt, minerId)
}

// CreateFileStoreInfo4Entire is a paid mutator transaction binding the contract method 0xecbc1191.
//
// Solidity: function createFileStoreInfo4Entire(bytes32 storeKeyHash, uint256 fileSize, string fileExt, bytes32 minerId) payable returns()
func (_FileStore *FileStoreSession) CreateFileStoreInfo4Entire(storeKeyHash [32]byte, fileSize *big.Int, fileExt string, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.CreateFileStoreInfo4Entire(&_FileStore.TransactOpts, storeKeyHash, fileSize, fileExt, minerId)
}

// CreateFileStoreInfo4Entire is a paid mutator transaction binding the contract method 0xecbc1191.
//
// Solidity: function createFileStoreInfo4Entire(bytes32 storeKeyHash, uint256 fileSize, string fileExt, bytes32 minerId) payable returns()
func (_FileStore *FileStoreTransactorSession) CreateFileStoreInfo4Entire(storeKeyHash [32]byte, fileSize *big.Int, fileExt string, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.CreateFileStoreInfo4Entire(&_FileStore.TransactOpts, storeKeyHash, fileSize, fileExt, minerId)
}

// ExtendFileDeadline is a paid mutator transaction binding the contract method 0xdb9234f8.
//
// Solidity: function extendFileDeadline(bytes32 oriHash, uint256 fileSize, bytes32 minerId) payable returns()
func (_FileStore *FileStoreTransactor) ExtendFileDeadline(opts *bind.TransactOpts, oriHash [32]byte, fileSize *big.Int, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "extendFileDeadline", oriHash, fileSize, minerId)
}

// ExtendFileDeadline is a paid mutator transaction binding the contract method 0xdb9234f8.
//
// Solidity: function extendFileDeadline(bytes32 oriHash, uint256 fileSize, bytes32 minerId) payable returns()
func (_FileStore *FileStoreSession) ExtendFileDeadline(oriHash [32]byte, fileSize *big.Int, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.ExtendFileDeadline(&_FileStore.TransactOpts, oriHash, fileSize, minerId)
}

// ExtendFileDeadline is a paid mutator transaction binding the contract method 0xdb9234f8.
//
// Solidity: function extendFileDeadline(bytes32 oriHash, uint256 fileSize, bytes32 minerId) payable returns()
func (_FileStore *FileStoreTransactorSession) ExtendFileDeadline(oriHash [32]byte, fileSize *big.Int, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.ExtendFileDeadline(&_FileStore.TransactOpts, oriHash, fileSize, minerId)
}

// ExtendFileDeadlineEntire is a paid mutator transaction binding the contract method 0x6af7fa9a.
//
// Solidity: function extendFileDeadlineEntire(bytes32 storeKeyHash, uint256 fileSize, bytes32 minerId) payable returns()
func (_FileStore *FileStoreTransactor) ExtendFileDeadlineEntire(opts *bind.TransactOpts, storeKeyHash [32]byte, fileSize *big.Int, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "extendFileDeadlineEntire", storeKeyHash, fileSize, minerId)
}

// ExtendFileDeadlineEntire is a paid mutator transaction binding the contract method 0x6af7fa9a.
//
// Solidity: function extendFileDeadlineEntire(bytes32 storeKeyHash, uint256 fileSize, bytes32 minerId) payable returns()
func (_FileStore *FileStoreSession) ExtendFileDeadlineEntire(storeKeyHash [32]byte, fileSize *big.Int, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.ExtendFileDeadlineEntire(&_FileStore.TransactOpts, storeKeyHash, fileSize, minerId)
}

// ExtendFileDeadlineEntire is a paid mutator transaction binding the contract method 0x6af7fa9a.
//
// Solidity: function extendFileDeadlineEntire(bytes32 storeKeyHash, uint256 fileSize, bytes32 minerId) payable returns()
func (_FileStore *FileStoreTransactorSession) ExtendFileDeadlineEntire(storeKeyHash [32]byte, fileSize *big.Int, minerId [32]byte) (*types.Transaction, error) {
	return _FileStore.Contract.ExtendFileDeadlineEntire(&_FileStore.TransactOpts, storeKeyHash, fileSize, minerId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.GrantRole(&_FileStore.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.GrantRole(&_FileStore.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] _addrs) returns()
func (_FileStore *FileStoreTransactor) Initialize(opts *bind.TransactOpts, _addrs []common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "initialize", _addrs)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] _addrs) returns()
func (_FileStore *FileStoreSession) Initialize(_addrs []common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.Initialize(&_FileStore.TransactOpts, _addrs)
}

// Initialize is a paid mutator transaction binding the contract method 0xa224cee7.
//
// Solidity: function initialize(address[] _addrs) returns()
func (_FileStore *FileStoreTransactorSession) Initialize(_addrs []common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.Initialize(&_FileStore.TransactOpts, _addrs)
}

// LockOrUnlock is a paid mutator transaction binding the contract method 0xb603c7d2.
//
// Solidity: function lockOrUnlock(bytes32 oriHash, bytes32 minerId, bool isLocked, bool isEntireFile) returns()
func (_FileStore *FileStoreTransactor) LockOrUnlock(opts *bind.TransactOpts, oriHash [32]byte, minerId [32]byte, isLocked bool, isEntireFile bool) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "lockOrUnlock", oriHash, minerId, isLocked, isEntireFile)
}

// LockOrUnlock is a paid mutator transaction binding the contract method 0xb603c7d2.
//
// Solidity: function lockOrUnlock(bytes32 oriHash, bytes32 minerId, bool isLocked, bool isEntireFile) returns()
func (_FileStore *FileStoreSession) LockOrUnlock(oriHash [32]byte, minerId [32]byte, isLocked bool, isEntireFile bool) (*types.Transaction, error) {
	return _FileStore.Contract.LockOrUnlock(&_FileStore.TransactOpts, oriHash, minerId, isLocked, isEntireFile)
}

// LockOrUnlock is a paid mutator transaction binding the contract method 0xb603c7d2.
//
// Solidity: function lockOrUnlock(bytes32 oriHash, bytes32 minerId, bool isLocked, bool isEntireFile) returns()
func (_FileStore *FileStoreTransactorSession) LockOrUnlock(oriHash [32]byte, minerId [32]byte, isLocked bool, isEntireFile bool) (*types.Transaction, error) {
	return _FileStore.Contract.LockOrUnlock(&_FileStore.TransactOpts, oriHash, minerId, isLocked, isEntireFile)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_FileStore *FileStoreTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_FileStore *FileStoreSession) Pause() (*types.Transaction, error) {
	return _FileStore.Contract.Pause(&_FileStore.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_FileStore *FileStoreTransactorSession) Pause() (*types.Transaction, error) {
	return _FileStore.Contract.Pause(&_FileStore.TransactOpts)
}

// RegDappContractAddr is a paid mutator transaction binding the contract method 0x7078264c.
//
// Solidity: function regDappContractAddr(bytes32 oriHash, address dappContractAddr) returns()
func (_FileStore *FileStoreTransactor) RegDappContractAddr(opts *bind.TransactOpts, oriHash [32]byte, dappContractAddr common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "regDappContractAddr", oriHash, dappContractAddr)
}

// RegDappContractAddr is a paid mutator transaction binding the contract method 0x7078264c.
//
// Solidity: function regDappContractAddr(bytes32 oriHash, address dappContractAddr) returns()
func (_FileStore *FileStoreSession) RegDappContractAddr(oriHash [32]byte, dappContractAddr common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.RegDappContractAddr(&_FileStore.TransactOpts, oriHash, dappContractAddr)
}

// RegDappContractAddr is a paid mutator transaction binding the contract method 0x7078264c.
//
// Solidity: function regDappContractAddr(bytes32 oriHash, address dappContractAddr) returns()
func (_FileStore *FileStoreTransactorSession) RegDappContractAddr(oriHash [32]byte, dappContractAddr common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.RegDappContractAddr(&_FileStore.TransactOpts, oriHash, dappContractAddr)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.RenounceRole(&_FileStore.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.RenounceRole(&_FileStore.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.RevokeRole(&_FileStore.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FileStore *FileStoreTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.RevokeRole(&_FileStore.TransactOpts, role, account)
}

// SetFileStoreStorageAddress is a paid mutator transaction binding the contract method 0x669f1427.
//
// Solidity: function setFileStoreStorageAddress(address _addr) returns()
func (_FileStore *FileStoreTransactor) SetFileStoreStorageAddress(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "setFileStoreStorageAddress", _addr)
}

// SetFileStoreStorageAddress is a paid mutator transaction binding the contract method 0x669f1427.
//
// Solidity: function setFileStoreStorageAddress(address _addr) returns()
func (_FileStore *FileStoreSession) SetFileStoreStorageAddress(_addr common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.SetFileStoreStorageAddress(&_FileStore.TransactOpts, _addr)
}

// SetFileStoreStorageAddress is a paid mutator transaction binding the contract method 0x669f1427.
//
// Solidity: function setFileStoreStorageAddress(address _addr) returns()
func (_FileStore *FileStoreTransactorSession) SetFileStoreStorageAddress(_addr common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.SetFileStoreStorageAddress(&_FileStore.TransactOpts, _addr)
}

// SetMinerInfo is a paid mutator transaction binding the contract method 0x331e4884.
//
// Solidity: function setMinerInfo(bytes32 minerId, string publicKey, string peerId, string peerAddr, string proxyAddr) returns()
func (_FileStore *FileStoreTransactor) SetMinerInfo(opts *bind.TransactOpts, minerId [32]byte, publicKey string, peerId string, peerAddr string, proxyAddr string) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "setMinerInfo", minerId, publicKey, peerId, peerAddr, proxyAddr)
}

// SetMinerInfo is a paid mutator transaction binding the contract method 0x331e4884.
//
// Solidity: function setMinerInfo(bytes32 minerId, string publicKey, string peerId, string peerAddr, string proxyAddr) returns()
func (_FileStore *FileStoreSession) SetMinerInfo(minerId [32]byte, publicKey string, peerId string, peerAddr string, proxyAddr string) (*types.Transaction, error) {
	return _FileStore.Contract.SetMinerInfo(&_FileStore.TransactOpts, minerId, publicKey, peerId, peerAddr, proxyAddr)
}

// SetMinerInfo is a paid mutator transaction binding the contract method 0x331e4884.
//
// Solidity: function setMinerInfo(bytes32 minerId, string publicKey, string peerId, string peerAddr, string proxyAddr) returns()
func (_FileStore *FileStoreTransactorSession) SetMinerInfo(minerId [32]byte, publicKey string, peerId string, peerAddr string, proxyAddr string) (*types.Transaction, error) {
	return _FileStore.Contract.SetMinerInfo(&_FileStore.TransactOpts, minerId, publicKey, peerId, peerAddr, proxyAddr)
}

// TransferFileOwner is a paid mutator transaction binding the contract method 0x717cc80c.
//
// Solidity: function transferFileOwner(bytes32 oriHash, address ownerAddr) returns()
func (_FileStore *FileStoreTransactor) TransferFileOwner(opts *bind.TransactOpts, oriHash [32]byte, ownerAddr common.Address) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "transferFileOwner", oriHash, ownerAddr)
}

// TransferFileOwner is a paid mutator transaction binding the contract method 0x717cc80c.
//
// Solidity: function transferFileOwner(bytes32 oriHash, address ownerAddr) returns()
func (_FileStore *FileStoreSession) TransferFileOwner(oriHash [32]byte, ownerAddr common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.TransferFileOwner(&_FileStore.TransactOpts, oriHash, ownerAddr)
}

// TransferFileOwner is a paid mutator transaction binding the contract method 0x717cc80c.
//
// Solidity: function transferFileOwner(bytes32 oriHash, address ownerAddr) returns()
func (_FileStore *FileStoreTransactorSession) TransferFileOwner(oriHash [32]byte, ownerAddr common.Address) (*types.Transaction, error) {
	return _FileStore.Contract.TransferFileOwner(&_FileStore.TransactOpts, oriHash, ownerAddr)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_FileStore *FileStoreTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_FileStore *FileStoreSession) Unpause() (*types.Transaction, error) {
	return _FileStore.Contract.Unpause(&_FileStore.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_FileStore *FileStoreTransactorSession) Unpause() (*types.Transaction, error) {
	return _FileStore.Contract.Unpause(&_FileStore.TransactOpts)
}

// UpdateFileStoreInfo is a paid mutator transaction binding the contract method 0x992a534d.
//
// Solidity: function updateFileStoreInfo(bytes32 oriHash, bool headFlag, string cid, uint8 status) payable returns()
func (_FileStore *FileStoreTransactor) UpdateFileStoreInfo(opts *bind.TransactOpts, oriHash [32]byte, headFlag bool, cid string, status uint8) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "updateFileStoreInfo", oriHash, headFlag, cid, status)
}

// UpdateFileStoreInfo is a paid mutator transaction binding the contract method 0x992a534d.
//
// Solidity: function updateFileStoreInfo(bytes32 oriHash, bool headFlag, string cid, uint8 status) payable returns()
func (_FileStore *FileStoreSession) UpdateFileStoreInfo(oriHash [32]byte, headFlag bool, cid string, status uint8) (*types.Transaction, error) {
	return _FileStore.Contract.UpdateFileStoreInfo(&_FileStore.TransactOpts, oriHash, headFlag, cid, status)
}

// UpdateFileStoreInfo is a paid mutator transaction binding the contract method 0x992a534d.
//
// Solidity: function updateFileStoreInfo(bytes32 oriHash, bool headFlag, string cid, uint8 status) payable returns()
func (_FileStore *FileStoreTransactorSession) UpdateFileStoreInfo(oriHash [32]byte, headFlag bool, cid string, status uint8) (*types.Transaction, error) {
	return _FileStore.Contract.UpdateFileStoreInfo(&_FileStore.TransactOpts, oriHash, headFlag, cid, status)
}

// UpdateFileStoreInfo4Entire is a paid mutator transaction binding the contract method 0xf3b2a74c.
//
// Solidity: function updateFileStoreInfo4Entire(bytes32 storeKeyHash, string cid, uint8 status) payable returns()
func (_FileStore *FileStoreTransactor) UpdateFileStoreInfo4Entire(opts *bind.TransactOpts, storeKeyHash [32]byte, cid string, status uint8) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "updateFileStoreInfo4Entire", storeKeyHash, cid, status)
}

// UpdateFileStoreInfo4Entire is a paid mutator transaction binding the contract method 0xf3b2a74c.
//
// Solidity: function updateFileStoreInfo4Entire(bytes32 storeKeyHash, string cid, uint8 status) payable returns()
func (_FileStore *FileStoreSession) UpdateFileStoreInfo4Entire(storeKeyHash [32]byte, cid string, status uint8) (*types.Transaction, error) {
	return _FileStore.Contract.UpdateFileStoreInfo4Entire(&_FileStore.TransactOpts, storeKeyHash, cid, status)
}

// UpdateFileStoreInfo4Entire is a paid mutator transaction binding the contract method 0xf3b2a74c.
//
// Solidity: function updateFileStoreInfo4Entire(bytes32 storeKeyHash, string cid, uint8 status) payable returns()
func (_FileStore *FileStoreTransactorSession) UpdateFileStoreInfo4Entire(storeKeyHash [32]byte, cid string, status uint8) (*types.Transaction, error) {
	return _FileStore.Contract.UpdateFileStoreInfo4Entire(&_FileStore.TransactOpts, storeKeyHash, cid, status)
}

// UpdateWithdrawThreshold is a paid mutator transaction binding the contract method 0x8608301d.
//
// Solidity: function updateWithdrawThreshold(uint256 threshold) returns()
func (_FileStore *FileStoreTransactor) UpdateWithdrawThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "updateWithdrawThreshold", threshold)
}

// UpdateWithdrawThreshold is a paid mutator transaction binding the contract method 0x8608301d.
//
// Solidity: function updateWithdrawThreshold(uint256 threshold) returns()
func (_FileStore *FileStoreSession) UpdateWithdrawThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _FileStore.Contract.UpdateWithdrawThreshold(&_FileStore.TransactOpts, threshold)
}

// UpdateWithdrawThreshold is a paid mutator transaction binding the contract method 0x8608301d.
//
// Solidity: function updateWithdrawThreshold(uint256 threshold) returns()
func (_FileStore *FileStoreTransactorSession) UpdateWithdrawThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _FileStore.Contract.UpdateWithdrawThreshold(&_FileStore.TransactOpts, threshold)
}

// WithdrawRemaining is a paid mutator transaction binding the contract method 0xcbe40abc.
//
// Solidity: function withdrawRemaining(bytes32 oriHash, uint256 index, uint8 storageType) payable returns()
func (_FileStore *FileStoreTransactor) WithdrawRemaining(opts *bind.TransactOpts, oriHash [32]byte, index *big.Int, storageType uint8) (*types.Transaction, error) {
	return _FileStore.contract.Transact(opts, "withdrawRemaining", oriHash, index, storageType)
}

// WithdrawRemaining is a paid mutator transaction binding the contract method 0xcbe40abc.
//
// Solidity: function withdrawRemaining(bytes32 oriHash, uint256 index, uint8 storageType) payable returns()
func (_FileStore *FileStoreSession) WithdrawRemaining(oriHash [32]byte, index *big.Int, storageType uint8) (*types.Transaction, error) {
	return _FileStore.Contract.WithdrawRemaining(&_FileStore.TransactOpts, oriHash, index, storageType)
}

// WithdrawRemaining is a paid mutator transaction binding the contract method 0xcbe40abc.
//
// Solidity: function withdrawRemaining(bytes32 oriHash, uint256 index, uint8 storageType) payable returns()
func (_FileStore *FileStoreTransactorSession) WithdrawRemaining(oriHash [32]byte, index *big.Int, storageType uint8) (*types.Transaction, error) {
	return _FileStore.Contract.WithdrawRemaining(&_FileStore.TransactOpts, oriHash, index, storageType)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_FileStore *FileStoreTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _FileStore.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_FileStore *FileStoreSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _FileStore.Contract.Fallback(&_FileStore.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_FileStore *FileStoreTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _FileStore.Contract.Fallback(&_FileStore.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FileStore *FileStoreTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileStore.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FileStore *FileStoreSession) Receive() (*types.Transaction, error) {
	return _FileStore.Contract.Receive(&_FileStore.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FileStore *FileStoreTransactorSession) Receive() (*types.Transaction, error) {
	return _FileStore.Contract.Receive(&_FileStore.TransactOpts)
}

// FileStorePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the FileStore contract.
type FileStorePausedIterator struct {
	Event *FileStorePaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStorePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStorePaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStorePaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStorePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStorePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStorePaused represents a Paused event raised by the FileStore contract.
type FileStorePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_FileStore *FileStoreFilterer) FilterPaused(opts *bind.FilterOpts) (*FileStorePausedIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &FileStorePausedIterator{contract: _FileStore.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_FileStore *FileStoreFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *FileStorePaused) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStorePaused)
				if err := _FileStore.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_FileStore *FileStoreFilterer) ParsePaused(log types.Log) (*FileStorePaused, error) {
	event := new(FileStorePaused)
	if err := _FileStore.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the FileStore contract.
type FileStoreRoleAdminChangedIterator struct {
	Event *FileStoreRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreRoleAdminChanged represents a RoleAdminChanged event raised by the FileStore contract.
type FileStoreRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FileStore *FileStoreFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*FileStoreRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &FileStoreRoleAdminChangedIterator{contract: _FileStore.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FileStore *FileStoreFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *FileStoreRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreRoleAdminChanged)
				if err := _FileStore.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FileStore *FileStoreFilterer) ParseRoleAdminChanged(log types.Log) (*FileStoreRoleAdminChanged, error) {
	event := new(FileStoreRoleAdminChanged)
	if err := _FileStore.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the FileStore contract.
type FileStoreRoleGrantedIterator struct {
	Event *FileStoreRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreRoleGranted represents a RoleGranted event raised by the FileStore contract.
type FileStoreRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FileStore *FileStoreFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FileStoreRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FileStoreRoleGrantedIterator{contract: _FileStore.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FileStore *FileStoreFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *FileStoreRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreRoleGranted)
				if err := _FileStore.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FileStore *FileStoreFilterer) ParseRoleGranted(log types.Log) (*FileStoreRoleGranted, error) {
	event := new(FileStoreRoleGranted)
	if err := _FileStore.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the FileStore contract.
type FileStoreRoleRevokedIterator struct {
	Event *FileStoreRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreRoleRevoked represents a RoleRevoked event raised by the FileStore contract.
type FileStoreRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FileStore *FileStoreFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FileStoreRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FileStoreRoleRevokedIterator{contract: _FileStore.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FileStore *FileStoreFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *FileStoreRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreRoleRevoked)
				if err := _FileStore.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FileStore *FileStoreFilterer) ParseRoleRevoked(log types.Log) (*FileStoreRoleRevoked, error) {
	event := new(FileStoreRoleRevoked)
	if err := _FileStore.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the FileStore contract.
type FileStoreUnpausedIterator struct {
	Event *FileStoreUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreUnpaused represents a Unpaused event raised by the FileStore contract.
type FileStoreUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_FileStore *FileStoreFilterer) FilterUnpaused(opts *bind.FilterOpts) (*FileStoreUnpausedIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &FileStoreUnpausedIterator{contract: _FileStore.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_FileStore *FileStoreFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FileStoreUnpaused) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreUnpaused)
				if err := _FileStore.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_FileStore *FileStoreFilterer) ParseUnpaused(log types.Log) (*FileStoreUnpaused, error) {
	event := new(FileStoreUnpaused)
	if err := _FileStore.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreFileInfoChangeEvtIterator is returned from FilterFileInfoChangeEvt and is used to iterate over the raw logs and unpacked data for FileInfoChangeEvt events raised by the FileStore contract.
type FileStoreFileInfoChangeEvtIterator struct {
	Event *FileStoreFileInfoChangeEvt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreFileInfoChangeEvtIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreFileInfoChangeEvt)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreFileInfoChangeEvt)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreFileInfoChangeEvtIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreFileInfoChangeEvtIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreFileInfoChangeEvt represents a FileInfoChangeEvt event raised by the FileStore contract.
type FileStoreFileInfoChangeEvt struct {
	OriHash  [32]byte
	MinerId  [32]byte
	HeadFlag bool
	Status   uint8
	Cid      string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFileInfoChangeEvt is a free log retrieval operation binding the contract event 0x669d194c406eba24a400019dc8a87ed1c62929fd3f9ee44d6c85a2ed10c550c0.
//
// Solidity: event fileInfoChangeEvt(bytes32 oriHash, bytes32 minerId, bool headFlag, uint8 status, string cid)
func (_FileStore *FileStoreFilterer) FilterFileInfoChangeEvt(opts *bind.FilterOpts) (*FileStoreFileInfoChangeEvtIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "fileInfoChangeEvt")
	if err != nil {
		return nil, err
	}
	return &FileStoreFileInfoChangeEvtIterator{contract: _FileStore.contract, event: "fileInfoChangeEvt", logs: logs, sub: sub}, nil
}

// WatchFileInfoChangeEvt is a free log subscription operation binding the contract event 0x669d194c406eba24a400019dc8a87ed1c62929fd3f9ee44d6c85a2ed10c550c0.
//
// Solidity: event fileInfoChangeEvt(bytes32 oriHash, bytes32 minerId, bool headFlag, uint8 status, string cid)
func (_FileStore *FileStoreFilterer) WatchFileInfoChangeEvt(opts *bind.WatchOpts, sink chan<- *FileStoreFileInfoChangeEvt) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "fileInfoChangeEvt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreFileInfoChangeEvt)
				if err := _FileStore.contract.UnpackLog(event, "fileInfoChangeEvt", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFileInfoChangeEvt is a log parse operation binding the contract event 0x669d194c406eba24a400019dc8a87ed1c62929fd3f9ee44d6c85a2ed10c550c0.
//
// Solidity: event fileInfoChangeEvt(bytes32 oriHash, bytes32 minerId, bool headFlag, uint8 status, string cid)
func (_FileStore *FileStoreFilterer) ParseFileInfoChangeEvt(log types.Log) (*FileStoreFileInfoChangeEvt, error) {
	event := new(FileStoreFileInfoChangeEvt)
	if err := _FileStore.contract.UnpackLog(event, "fileInfoChangeEvt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreNewFileStoreEvtIterator is returned from FilterNewFileStoreEvt and is used to iterate over the raw logs and unpacked data for NewFileStoreEvt events raised by the FileStore contract.
type FileStoreNewFileStoreEvtIterator struct {
	Event *FileStoreNewFileStoreEvt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreNewFileStoreEvtIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreNewFileStoreEvt)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreNewFileStoreEvt)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreNewFileStoreEvtIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreNewFileStoreEvtIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreNewFileStoreEvt represents a NewFileStoreEvt event raised by the FileStore contract.
type FileStoreNewFileStoreEvt struct {
	OriHash  [32]byte
	UserAddr common.Address
	FileSize *big.Int
	FileExt  string
	MinerId  [32]byte
	HeadHash [32]byte
	BodyHash [32]byte
	OperTime *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewFileStoreEvt is a free log retrieval operation binding the contract event 0x22a6a0618cb6b5ca64a2b7a8ccde9aa6340c73ff189ca6ce132d116fb9d48928.
//
// Solidity: event newFileStoreEvt(bytes32 oriHash, address userAddr, uint256 fileSize, string fileExt, bytes32 minerId, bytes32 headHash, bytes32 bodyHash, uint256 operTime)
func (_FileStore *FileStoreFilterer) FilterNewFileStoreEvt(opts *bind.FilterOpts) (*FileStoreNewFileStoreEvtIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "newFileStoreEvt")
	if err != nil {
		return nil, err
	}
	return &FileStoreNewFileStoreEvtIterator{contract: _FileStore.contract, event: "newFileStoreEvt", logs: logs, sub: sub}, nil
}

// WatchNewFileStoreEvt is a free log subscription operation binding the contract event 0x22a6a0618cb6b5ca64a2b7a8ccde9aa6340c73ff189ca6ce132d116fb9d48928.
//
// Solidity: event newFileStoreEvt(bytes32 oriHash, address userAddr, uint256 fileSize, string fileExt, bytes32 minerId, bytes32 headHash, bytes32 bodyHash, uint256 operTime)
func (_FileStore *FileStoreFilterer) WatchNewFileStoreEvt(opts *bind.WatchOpts, sink chan<- *FileStoreNewFileStoreEvt) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "newFileStoreEvt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreNewFileStoreEvt)
				if err := _FileStore.contract.UnpackLog(event, "newFileStoreEvt", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewFileStoreEvt is a log parse operation binding the contract event 0x22a6a0618cb6b5ca64a2b7a8ccde9aa6340c73ff189ca6ce132d116fb9d48928.
//
// Solidity: event newFileStoreEvt(bytes32 oriHash, address userAddr, uint256 fileSize, string fileExt, bytes32 minerId, bytes32 headHash, bytes32 bodyHash, uint256 operTime)
func (_FileStore *FileStoreFilterer) ParseNewFileStoreEvt(log types.Log) (*FileStoreNewFileStoreEvt, error) {
	event := new(FileStoreNewFileStoreEvt)
	if err := _FileStore.contract.UnpackLog(event, "newFileStoreEvt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreRegisterDappContractAddrEvtIterator is returned from FilterRegisterDappContractAddrEvt and is used to iterate over the raw logs and unpacked data for RegisterDappContractAddrEvt events raised by the FileStore contract.
type FileStoreRegisterDappContractAddrEvtIterator struct {
	Event *FileStoreRegisterDappContractAddrEvt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreRegisterDappContractAddrEvtIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreRegisterDappContractAddrEvt)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreRegisterDappContractAddrEvt)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreRegisterDappContractAddrEvtIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreRegisterDappContractAddrEvtIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreRegisterDappContractAddrEvt represents a RegisterDappContractAddrEvt event raised by the FileStore contract.
type FileStoreRegisterDappContractAddrEvt struct {
	OriHash          [32]byte
	OwnerAddr        common.Address
	DappContractAddr common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRegisterDappContractAddrEvt is a free log retrieval operation binding the contract event 0xb57d6197d3a659be7be2ea7fc4f17d13ab3f2ca1c5d0b8d937c97cae1beefc76.
//
// Solidity: event registerDappContractAddrEvt(bytes32 oriHash, address ownerAddr, address dappContractAddr)
func (_FileStore *FileStoreFilterer) FilterRegisterDappContractAddrEvt(opts *bind.FilterOpts) (*FileStoreRegisterDappContractAddrEvtIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "registerDappContractAddrEvt")
	if err != nil {
		return nil, err
	}
	return &FileStoreRegisterDappContractAddrEvtIterator{contract: _FileStore.contract, event: "registerDappContractAddrEvt", logs: logs, sub: sub}, nil
}

// WatchRegisterDappContractAddrEvt is a free log subscription operation binding the contract event 0xb57d6197d3a659be7be2ea7fc4f17d13ab3f2ca1c5d0b8d937c97cae1beefc76.
//
// Solidity: event registerDappContractAddrEvt(bytes32 oriHash, address ownerAddr, address dappContractAddr)
func (_FileStore *FileStoreFilterer) WatchRegisterDappContractAddrEvt(opts *bind.WatchOpts, sink chan<- *FileStoreRegisterDappContractAddrEvt) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "registerDappContractAddrEvt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreRegisterDappContractAddrEvt)
				if err := _FileStore.contract.UnpackLog(event, "registerDappContractAddrEvt", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegisterDappContractAddrEvt is a log parse operation binding the contract event 0xb57d6197d3a659be7be2ea7fc4f17d13ab3f2ca1c5d0b8d937c97cae1beefc76.
//
// Solidity: event registerDappContractAddrEvt(bytes32 oriHash, address ownerAddr, address dappContractAddr)
func (_FileStore *FileStoreFilterer) ParseRegisterDappContractAddrEvt(log types.Log) (*FileStoreRegisterDappContractAddrEvt, error) {
	event := new(FileStoreRegisterDappContractAddrEvt)
	if err := _FileStore.contract.UnpackLog(event, "registerDappContractAddrEvt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreSetMinerEvtIterator is returned from FilterSetMinerEvt and is used to iterate over the raw logs and unpacked data for SetMinerEvt events raised by the FileStore contract.
type FileStoreSetMinerEvtIterator struct {
	Event *FileStoreSetMinerEvt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreSetMinerEvtIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreSetMinerEvt)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreSetMinerEvt)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreSetMinerEvtIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreSetMinerEvtIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreSetMinerEvt represents a SetMinerEvt event raised by the FileStore contract.
type FileStoreSetMinerEvt struct {
	MinerId   [32]byte
	MinerAddr common.Address
	PublicKey string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetMinerEvt is a free log retrieval operation binding the contract event 0x6395949b438cd6c015a95f3a15dd3f0bf64a249ff80455c6684d67a33212451d.
//
// Solidity: event setMinerEvt(bytes32 minerId, address minerAddr, string publicKey)
func (_FileStore *FileStoreFilterer) FilterSetMinerEvt(opts *bind.FilterOpts) (*FileStoreSetMinerEvtIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "setMinerEvt")
	if err != nil {
		return nil, err
	}
	return &FileStoreSetMinerEvtIterator{contract: _FileStore.contract, event: "setMinerEvt", logs: logs, sub: sub}, nil
}

// WatchSetMinerEvt is a free log subscription operation binding the contract event 0x6395949b438cd6c015a95f3a15dd3f0bf64a249ff80455c6684d67a33212451d.
//
// Solidity: event setMinerEvt(bytes32 minerId, address minerAddr, string publicKey)
func (_FileStore *FileStoreFilterer) WatchSetMinerEvt(opts *bind.WatchOpts, sink chan<- *FileStoreSetMinerEvt) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "setMinerEvt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreSetMinerEvt)
				if err := _FileStore.contract.UnpackLog(event, "setMinerEvt", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetMinerEvt is a log parse operation binding the contract event 0x6395949b438cd6c015a95f3a15dd3f0bf64a249ff80455c6684d67a33212451d.
//
// Solidity: event setMinerEvt(bytes32 minerId, address minerAddr, string publicKey)
func (_FileStore *FileStoreFilterer) ParseSetMinerEvt(log types.Log) (*FileStoreSetMinerEvt, error) {
	event := new(FileStoreSetMinerEvt)
	if err := _FileStore.contract.UnpackLog(event, "setMinerEvt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileStoreTransferOwnerEvtIterator is returned from FilterTransferOwnerEvt and is used to iterate over the raw logs and unpacked data for TransferOwnerEvt events raised by the FileStore contract.
type FileStoreTransferOwnerEvtIterator struct {
	Event *FileStoreTransferOwnerEvt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileStoreTransferOwnerEvtIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileStoreTransferOwnerEvt)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FileStoreTransferOwnerEvt)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FileStoreTransferOwnerEvtIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileStoreTransferOwnerEvtIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileStoreTransferOwnerEvt represents a TransferOwnerEvt event raised by the FileStore contract.
type FileStoreTransferOwnerEvt struct {
	OriHash      [32]byte
	OldOwnerAddr common.Address
	NewOwnerAddr common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferOwnerEvt is a free log retrieval operation binding the contract event 0xbe683ff9a6f425521c79acd1ca4b7174994034e74a01ec06c90e7de6862f42d6.
//
// Solidity: event transferOwnerEvt(bytes32 oriHash, address oldOwnerAddr, address newOwnerAddr)
func (_FileStore *FileStoreFilterer) FilterTransferOwnerEvt(opts *bind.FilterOpts) (*FileStoreTransferOwnerEvtIterator, error) {

	logs, sub, err := _FileStore.contract.FilterLogs(opts, "transferOwnerEvt")
	if err != nil {
		return nil, err
	}
	return &FileStoreTransferOwnerEvtIterator{contract: _FileStore.contract, event: "transferOwnerEvt", logs: logs, sub: sub}, nil
}

// WatchTransferOwnerEvt is a free log subscription operation binding the contract event 0xbe683ff9a6f425521c79acd1ca4b7174994034e74a01ec06c90e7de6862f42d6.
//
// Solidity: event transferOwnerEvt(bytes32 oriHash, address oldOwnerAddr, address newOwnerAddr)
func (_FileStore *FileStoreFilterer) WatchTransferOwnerEvt(opts *bind.WatchOpts, sink chan<- *FileStoreTransferOwnerEvt) (event.Subscription, error) {

	logs, sub, err := _FileStore.contract.WatchLogs(opts, "transferOwnerEvt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileStoreTransferOwnerEvt)
				if err := _FileStore.contract.UnpackLog(event, "transferOwnerEvt", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferOwnerEvt is a log parse operation binding the contract event 0xbe683ff9a6f425521c79acd1ca4b7174994034e74a01ec06c90e7de6862f42d6.
//
// Solidity: event transferOwnerEvt(bytes32 oriHash, address oldOwnerAddr, address newOwnerAddr)
func (_FileStore *FileStoreFilterer) ParseTransferOwnerEvt(log types.Log) (*FileStoreTransferOwnerEvt, error) {
	event := new(FileStoreTransferOwnerEvt)
	if err := _FileStore.contract.UnpackLog(event, "transferOwnerEvt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
