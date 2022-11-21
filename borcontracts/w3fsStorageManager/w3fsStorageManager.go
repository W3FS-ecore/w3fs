// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package w3fsStorageManager

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
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

// W3fsStorageManagerSector is an auto generated low-level Go binding around an user-defined struct.
type W3fsStorageManagerSector struct {
	SealProofType *big.Int
	SectorNumber  *big.Int
	TicketEpoch   *big.Int
	SeedEpoch     *big.Int
	SealedCID     []byte
	UnsealedCID   []byte
	Proof         []byte
	Check         bool
	IsReal        bool
}

// W3fsStorageManagerMetaData contains all meta data concerning the W3fsStorageManager contract.
var W3fsStorageManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newPower\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newStorageSize\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"AddNewSealPowerAndSize\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"storageSize\",\"type\":\"uint256\"}],\"name\":\"UpdatePromise\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"baseStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"delegatedStakeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"governance\",\"outputs\":[{\"internalType\":\"contractIGovernance\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"locked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"logger\",\"outputs\":[{\"internalType\":\"contractW3fsStakingInfo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"percentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorPowers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorPromise\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSector\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"SealProofType\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"SectorNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"TicketEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"SeedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"SealedCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"UnsealedCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"Proof\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"Check\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isReal\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorStorageSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_registry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_governance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakingLogger\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"addStakeMount\",\"type\":\"uint256\"}],\"name\":\"checkCanStakeMore\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"addStakeMount\",\"type\":\"uint256\"}],\"name\":\"checkCandelegatorsMore\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"}],\"name\":\"showCanStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"storageSize\",\"type\":\"uint256\"}],\"name\":\"updateStoragePromise\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newStakeLimit\",\"type\":\"uint256\"}],\"name\":\"updateStakeLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newDelegatedStakeLimit\",\"type\":\"uint256\"}],\"name\":\"updateDelegatedStakeLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPercentage\",\"type\":\"uint256\"}],\"name\":\"updatePercentage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sectorNumber\",\"type\":\"uint256\"}],\"name\":\"getSealInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"SealProofType\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"SectorNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"TicketEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"SeedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"SealedCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"UnsealedCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"Proof\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"Check\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isReal\",\"type\":\"bool\"}],\"internalType\":\"structW3fsStorageManager.Sector\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isCheck\",\"type\":\"bool\"}],\"name\":\"getSealInfoAllBySigner\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"SealProofType\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"SectorNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"TicketEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"SeedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"SealedCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"UnsealedCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"Proof\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"Check\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isReal\",\"type\":\"bool\"}],\"internalType\":\"structW3fsStorageManager.Sector[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"isReal\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"votes\",\"type\":\"bytes\"}],\"name\":\"addSealInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256[3][]\",\"name\":\"sigs\",\"type\":\"uint256[3][]\"}],\"name\":\"checkSealSigs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"getValidatorPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"validatorBytes\",\"type\":\"bytes\"}],\"name\":\"getAllValidatorPower\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// W3fsStorageManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use W3fsStorageManagerMetaData.ABI instead.
var W3fsStorageManagerABI = W3fsStorageManagerMetaData.ABI

// W3fsStorageManager is an auto generated Go binding around an Ethereum contract.
type W3fsStorageManager struct {
	W3fsStorageManagerCaller     // Read-only binding to the contract
	W3fsStorageManagerTransactor // Write-only binding to the contract
	W3fsStorageManagerFilterer   // Log filterer for contract events
}

// W3fsStorageManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type W3fsStorageManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// W3fsStorageManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type W3fsStorageManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// W3fsStorageManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type W3fsStorageManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// W3fsStorageManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type W3fsStorageManagerSession struct {
	Contract     *W3fsStorageManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// W3fsStorageManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type W3fsStorageManagerCallerSession struct {
	Contract *W3fsStorageManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// W3fsStorageManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type W3fsStorageManagerTransactorSession struct {
	Contract     *W3fsStorageManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// W3fsStorageManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type W3fsStorageManagerRaw struct {
	Contract *W3fsStorageManager // Generic contract binding to access the raw methods on
}

// W3fsStorageManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type W3fsStorageManagerCallerRaw struct {
	Contract *W3fsStorageManagerCaller // Generic read-only contract binding to access the raw methods on
}

// W3fsStorageManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type W3fsStorageManagerTransactorRaw struct {
	Contract *W3fsStorageManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewW3fsStorageManager creates a new instance of W3fsStorageManager, bound to a specific deployed contract.
func NewW3fsStorageManager(address common.Address, backend bind.ContractBackend) (*W3fsStorageManager, error) {
	contract, err := bindW3fsStorageManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManager{W3fsStorageManagerCaller: W3fsStorageManagerCaller{contract: contract}, W3fsStorageManagerTransactor: W3fsStorageManagerTransactor{contract: contract}, W3fsStorageManagerFilterer: W3fsStorageManagerFilterer{contract: contract}}, nil
}

// NewW3fsStorageManagerCaller creates a new read-only instance of W3fsStorageManager, bound to a specific deployed contract.
func NewW3fsStorageManagerCaller(address common.Address, caller bind.ContractCaller) (*W3fsStorageManagerCaller, error) {
	contract, err := bindW3fsStorageManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManagerCaller{contract: contract}, nil
}

// NewW3fsStorageManagerTransactor creates a new write-only instance of W3fsStorageManager, bound to a specific deployed contract.
func NewW3fsStorageManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*W3fsStorageManagerTransactor, error) {
	contract, err := bindW3fsStorageManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManagerTransactor{contract: contract}, nil
}

// NewW3fsStorageManagerFilterer creates a new log filterer instance of W3fsStorageManager, bound to a specific deployed contract.
func NewW3fsStorageManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*W3fsStorageManagerFilterer, error) {
	contract, err := bindW3fsStorageManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManagerFilterer{contract: contract}, nil
}

// bindW3fsStorageManager binds a generic wrapper to an already deployed contract.
func bindW3fsStorageManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(W3fsStorageManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_W3fsStorageManager *W3fsStorageManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _W3fsStorageManager.Contract.W3fsStorageManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_W3fsStorageManager *W3fsStorageManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.W3fsStorageManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_W3fsStorageManager *W3fsStorageManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.W3fsStorageManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_W3fsStorageManager *W3fsStorageManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _W3fsStorageManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_W3fsStorageManager *W3fsStorageManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_W3fsStorageManager *W3fsStorageManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.contract.Transact(opts, method, params...)
}

// BaseStakeAmount is a free data retrieval call binding the contract method 0x71129559.
//
// Solidity: function baseStakeAmount() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) BaseStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "baseStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseStakeAmount is a free data retrieval call binding the contract method 0x71129559.
//
// Solidity: function baseStakeAmount() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) BaseStakeAmount() (*big.Int, error) {
	return _W3fsStorageManager.Contract.BaseStakeAmount(&_W3fsStorageManager.CallOpts)
}

// BaseStakeAmount is a free data retrieval call binding the contract method 0x71129559.
//
// Solidity: function baseStakeAmount() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) BaseStakeAmount() (*big.Int, error) {
	return _W3fsStorageManager.Contract.BaseStakeAmount(&_W3fsStorageManager.CallOpts)
}

// CheckCanStakeMore is a free data retrieval call binding the contract method 0x20edccbc.
//
// Solidity: function checkCanStakeMore(address validatorAddr, uint256 amount, uint256 addStakeMount) view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerCaller) CheckCanStakeMore(opts *bind.CallOpts, validatorAddr common.Address, amount *big.Int, addStakeMount *big.Int) (bool, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "checkCanStakeMore", validatorAddr, amount, addStakeMount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckCanStakeMore is a free data retrieval call binding the contract method 0x20edccbc.
//
// Solidity: function checkCanStakeMore(address validatorAddr, uint256 amount, uint256 addStakeMount) view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerSession) CheckCanStakeMore(validatorAddr common.Address, amount *big.Int, addStakeMount *big.Int) (bool, error) {
	return _W3fsStorageManager.Contract.CheckCanStakeMore(&_W3fsStorageManager.CallOpts, validatorAddr, amount, addStakeMount)
}

// CheckCanStakeMore is a free data retrieval call binding the contract method 0x20edccbc.
//
// Solidity: function checkCanStakeMore(address validatorAddr, uint256 amount, uint256 addStakeMount) view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) CheckCanStakeMore(validatorAddr common.Address, amount *big.Int, addStakeMount *big.Int) (bool, error) {
	return _W3fsStorageManager.Contract.CheckCanStakeMore(&_W3fsStorageManager.CallOpts, validatorAddr, amount, addStakeMount)
}

// CheckCandelegatorsMore is a free data retrieval call binding the contract method 0xe70c9357.
//
// Solidity: function checkCandelegatorsMore(uint256 minerId, uint256 addStakeMount) view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerCaller) CheckCandelegatorsMore(opts *bind.CallOpts, minerId *big.Int, addStakeMount *big.Int) (bool, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "checkCandelegatorsMore", minerId, addStakeMount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckCandelegatorsMore is a free data retrieval call binding the contract method 0xe70c9357.
//
// Solidity: function checkCandelegatorsMore(uint256 minerId, uint256 addStakeMount) view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerSession) CheckCandelegatorsMore(minerId *big.Int, addStakeMount *big.Int) (bool, error) {
	return _W3fsStorageManager.Contract.CheckCandelegatorsMore(&_W3fsStorageManager.CallOpts, minerId, addStakeMount)
}

// CheckCandelegatorsMore is a free data retrieval call binding the contract method 0xe70c9357.
//
// Solidity: function checkCandelegatorsMore(uint256 minerId, uint256 addStakeMount) view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) CheckCandelegatorsMore(minerId *big.Int, addStakeMount *big.Int) (bool, error) {
	return _W3fsStorageManager.Contract.CheckCandelegatorsMore(&_W3fsStorageManager.CallOpts, minerId, addStakeMount)
}

// DelegatedStakeLimit is a free data retrieval call binding the contract method 0xbd6e0e51.
//
// Solidity: function delegatedStakeLimit() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) DelegatedStakeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "delegatedStakeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegatedStakeLimit is a free data retrieval call binding the contract method 0xbd6e0e51.
//
// Solidity: function delegatedStakeLimit() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) DelegatedStakeLimit() (*big.Int, error) {
	return _W3fsStorageManager.Contract.DelegatedStakeLimit(&_W3fsStorageManager.CallOpts)
}

// DelegatedStakeLimit is a free data retrieval call binding the contract method 0xbd6e0e51.
//
// Solidity: function delegatedStakeLimit() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) DelegatedStakeLimit() (*big.Int, error) {
	return _W3fsStorageManager.Contract.DelegatedStakeLimit(&_W3fsStorageManager.CallOpts)
}

// Factor is a free data retrieval call binding the contract method 0x54f703f8.
//
// Solidity: function factor() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Factor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "factor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Factor is a free data retrieval call binding the contract method 0x54f703f8.
//
// Solidity: function factor() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) Factor() (*big.Int, error) {
	return _W3fsStorageManager.Contract.Factor(&_W3fsStorageManager.CallOpts)
}

// Factor is a free data retrieval call binding the contract method 0x54f703f8.
//
// Solidity: function factor() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Factor() (*big.Int, error) {
	return _W3fsStorageManager.Contract.Factor(&_W3fsStorageManager.CallOpts)
}

// GetAllValidatorPower is a free data retrieval call binding the contract method 0x3b17b903.
//
// Solidity: function getAllValidatorPower(bytes validatorBytes) view returns(address[], uint256[])
func (_W3fsStorageManager *W3fsStorageManagerCaller) GetAllValidatorPower(opts *bind.CallOpts, validatorBytes []byte) ([]common.Address, []*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "getAllValidatorPower", validatorBytes)

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

// GetAllValidatorPower is a free data retrieval call binding the contract method 0x3b17b903.
//
// Solidity: function getAllValidatorPower(bytes validatorBytes) view returns(address[], uint256[])
func (_W3fsStorageManager *W3fsStorageManagerSession) GetAllValidatorPower(validatorBytes []byte) ([]common.Address, []*big.Int, error) {
	return _W3fsStorageManager.Contract.GetAllValidatorPower(&_W3fsStorageManager.CallOpts, validatorBytes)
}

// GetAllValidatorPower is a free data retrieval call binding the contract method 0x3b17b903.
//
// Solidity: function getAllValidatorPower(bytes validatorBytes) view returns(address[], uint256[])
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) GetAllValidatorPower(validatorBytes []byte) ([]common.Address, []*big.Int, error) {
	return _W3fsStorageManager.Contract.GetAllValidatorPower(&_W3fsStorageManager.CallOpts, validatorBytes)
}

// GetSealInfo is a free data retrieval call binding the contract method 0x578f5b30.
//
// Solidity: function getSealInfo(address signer, uint256 sectorNumber) view returns((uint256,uint256,uint256,uint256,bytes,bytes,bytes,bool,bool))
func (_W3fsStorageManager *W3fsStorageManagerCaller) GetSealInfo(opts *bind.CallOpts, signer common.Address, sectorNumber *big.Int) (W3fsStorageManagerSector, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "getSealInfo", signer, sectorNumber)

	if err != nil {
		return *new(W3fsStorageManagerSector), err
	}

	out0 := *abi.ConvertType(out[0], new(W3fsStorageManagerSector)).(*W3fsStorageManagerSector)

	return out0, err

}

// GetSealInfo is a free data retrieval call binding the contract method 0x578f5b30.
//
// Solidity: function getSealInfo(address signer, uint256 sectorNumber) view returns((uint256,uint256,uint256,uint256,bytes,bytes,bytes,bool,bool))
func (_W3fsStorageManager *W3fsStorageManagerSession) GetSealInfo(signer common.Address, sectorNumber *big.Int) (W3fsStorageManagerSector, error) {
	return _W3fsStorageManager.Contract.GetSealInfo(&_W3fsStorageManager.CallOpts, signer, sectorNumber)
}

// GetSealInfo is a free data retrieval call binding the contract method 0x578f5b30.
//
// Solidity: function getSealInfo(address signer, uint256 sectorNumber) view returns((uint256,uint256,uint256,uint256,bytes,bytes,bytes,bool,bool))
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) GetSealInfo(signer common.Address, sectorNumber *big.Int) (W3fsStorageManagerSector, error) {
	return _W3fsStorageManager.Contract.GetSealInfo(&_W3fsStorageManager.CallOpts, signer, sectorNumber)
}

// GetSealInfoAllBySigner is a free data retrieval call binding the contract method 0x41961249.
//
// Solidity: function getSealInfoAllBySigner(address signer, bool isCheck) view returns((uint256,uint256,uint256,uint256,bytes,bytes,bytes,bool,bool)[])
func (_W3fsStorageManager *W3fsStorageManagerCaller) GetSealInfoAllBySigner(opts *bind.CallOpts, signer common.Address, isCheck bool) ([]W3fsStorageManagerSector, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "getSealInfoAllBySigner", signer, isCheck)

	if err != nil {
		return *new([]W3fsStorageManagerSector), err
	}

	out0 := *abi.ConvertType(out[0], new([]W3fsStorageManagerSector)).(*[]W3fsStorageManagerSector)

	return out0, err

}

// GetSealInfoAllBySigner is a free data retrieval call binding the contract method 0x41961249.
//
// Solidity: function getSealInfoAllBySigner(address signer, bool isCheck) view returns((uint256,uint256,uint256,uint256,bytes,bytes,bytes,bool,bool)[])
func (_W3fsStorageManager *W3fsStorageManagerSession) GetSealInfoAllBySigner(signer common.Address, isCheck bool) ([]W3fsStorageManagerSector, error) {
	return _W3fsStorageManager.Contract.GetSealInfoAllBySigner(&_W3fsStorageManager.CallOpts, signer, isCheck)
}

// GetSealInfoAllBySigner is a free data retrieval call binding the contract method 0x41961249.
//
// Solidity: function getSealInfoAllBySigner(address signer, bool isCheck) view returns((uint256,uint256,uint256,uint256,bytes,bytes,bytes,bool,bool)[])
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) GetSealInfoAllBySigner(signer common.Address, isCheck bool) ([]W3fsStorageManagerSector, error) {
	return _W3fsStorageManager.Contract.GetSealInfoAllBySigner(&_W3fsStorageManager.CallOpts, signer, isCheck)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address signer) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) GetValidatorPower(opts *bind.CallOpts, signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "getValidatorPower", signer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address signer) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) GetValidatorPower(signer common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.GetValidatorPower(&_W3fsStorageManager.CallOpts, signer)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address signer) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) GetValidatorPower(signer common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.GetValidatorPower(&_W3fsStorageManager.CallOpts, signer)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Governance(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "governance")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerSession) Governance() (common.Address, error) {
	return _W3fsStorageManager.Contract.Governance(&_W3fsStorageManager.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Governance() (common.Address, error) {
	return _W3fsStorageManager.Contract.Governance(&_W3fsStorageManager.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Locked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "locked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerSession) Locked() (bool, error) {
	return _W3fsStorageManager.Contract.Locked(&_W3fsStorageManager.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Locked() (bool, error) {
	return _W3fsStorageManager.Contract.Locked(&_W3fsStorageManager.CallOpts)
}

// Logger is a free data retrieval call binding the contract method 0xf24ccbfe.
//
// Solidity: function logger() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Logger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "logger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Logger is a free data retrieval call binding the contract method 0xf24ccbfe.
//
// Solidity: function logger() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerSession) Logger() (common.Address, error) {
	return _W3fsStorageManager.Contract.Logger(&_W3fsStorageManager.CallOpts)
}

// Logger is a free data retrieval call binding the contract method 0xf24ccbfe.
//
// Solidity: function logger() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Logger() (common.Address, error) {
	return _W3fsStorageManager.Contract.Logger(&_W3fsStorageManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerSession) Owner() (common.Address, error) {
	return _W3fsStorageManager.Contract.Owner(&_W3fsStorageManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Owner() (common.Address, error) {
	return _W3fsStorageManager.Contract.Owner(&_W3fsStorageManager.CallOpts)
}

// Percentage is a free data retrieval call binding the contract method 0xc78ad77f.
//
// Solidity: function percentage() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Percentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "percentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Percentage is a free data retrieval call binding the contract method 0xc78ad77f.
//
// Solidity: function percentage() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) Percentage() (*big.Int, error) {
	return _W3fsStorageManager.Contract.Percentage(&_W3fsStorageManager.CallOpts)
}

// Percentage is a free data retrieval call binding the contract method 0xc78ad77f.
//
// Solidity: function percentage() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Percentage() (*big.Int, error) {
	return _W3fsStorageManager.Contract.Percentage(&_W3fsStorageManager.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerSession) Registry() (common.Address, error) {
	return _W3fsStorageManager.Contract.Registry(&_W3fsStorageManager.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) Registry() (common.Address, error) {
	return _W3fsStorageManager.Contract.Registry(&_W3fsStorageManager.CallOpts)
}

// ShowCanStakeAmount is a free data retrieval call binding the contract method 0xc330e9a7.
//
// Solidity: function showCanStakeAmount(address validatorAddr) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) ShowCanStakeAmount(opts *bind.CallOpts, validatorAddr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "showCanStakeAmount", validatorAddr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ShowCanStakeAmount is a free data retrieval call binding the contract method 0xc330e9a7.
//
// Solidity: function showCanStakeAmount(address validatorAddr) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) ShowCanStakeAmount(validatorAddr common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ShowCanStakeAmount(&_W3fsStorageManager.CallOpts, validatorAddr)
}

// ShowCanStakeAmount is a free data retrieval call binding the contract method 0xc330e9a7.
//
// Solidity: function showCanStakeAmount(address validatorAddr) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) ShowCanStakeAmount(validatorAddr common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ShowCanStakeAmount(&_W3fsStorageManager.CallOpts, validatorAddr)
}

// StakeLimit is a free data retrieval call binding the contract method 0x45ef79af.
//
// Solidity: function stakeLimit() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) StakeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "stakeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeLimit is a free data retrieval call binding the contract method 0x45ef79af.
//
// Solidity: function stakeLimit() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) StakeLimit() (*big.Int, error) {
	return _W3fsStorageManager.Contract.StakeLimit(&_W3fsStorageManager.CallOpts)
}

// StakeLimit is a free data retrieval call binding the contract method 0x45ef79af.
//
// Solidity: function stakeLimit() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) StakeLimit() (*big.Int, error) {
	return _W3fsStorageManager.Contract.StakeLimit(&_W3fsStorageManager.CallOpts)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) TotalPower(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "totalPower")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) TotalPower() (*big.Int, error) {
	return _W3fsStorageManager.Contract.TotalPower(&_W3fsStorageManager.CallOpts)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) TotalPower() (*big.Int, error) {
	return _W3fsStorageManager.Contract.TotalPower(&_W3fsStorageManager.CallOpts)
}

// ValidatorNonce is a free data retrieval call binding the contract method 0x0ce10484.
//
// Solidity: function validatorNonce(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) ValidatorNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "validatorNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorNonce is a free data retrieval call binding the contract method 0x0ce10484.
//
// Solidity: function validatorNonce(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) ValidatorNonce(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorNonce(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorNonce is a free data retrieval call binding the contract method 0x0ce10484.
//
// Solidity: function validatorNonce(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) ValidatorNonce(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorNonce(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorPowers is a free data retrieval call binding the contract method 0xf61aa09d.
//
// Solidity: function validatorPowers(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) ValidatorPowers(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "validatorPowers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorPowers is a free data retrieval call binding the contract method 0xf61aa09d.
//
// Solidity: function validatorPowers(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) ValidatorPowers(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorPowers(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorPowers is a free data retrieval call binding the contract method 0xf61aa09d.
//
// Solidity: function validatorPowers(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) ValidatorPowers(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorPowers(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorPromise is a free data retrieval call binding the contract method 0x1f94d137.
//
// Solidity: function validatorPromise(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) ValidatorPromise(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "validatorPromise", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorPromise is a free data retrieval call binding the contract method 0x1f94d137.
//
// Solidity: function validatorPromise(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) ValidatorPromise(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorPromise(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorPromise is a free data retrieval call binding the contract method 0x1f94d137.
//
// Solidity: function validatorPromise(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) ValidatorPromise(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorPromise(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorSector is a free data retrieval call binding the contract method 0x4cdf3f46.
//
// Solidity: function validatorSector(address , uint256 ) view returns(uint256 SealProofType, uint256 SectorNumber, uint256 TicketEpoch, uint256 SeedEpoch, bytes SealedCID, bytes UnsealedCID, bytes Proof, bool Check, bool isReal)
func (_W3fsStorageManager *W3fsStorageManagerCaller) ValidatorSector(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	SealProofType *big.Int
	SectorNumber  *big.Int
	TicketEpoch   *big.Int
	SeedEpoch     *big.Int
	SealedCID     []byte
	UnsealedCID   []byte
	Proof         []byte
	Check         bool
	IsReal        bool
}, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "validatorSector", arg0, arg1)

	outstruct := new(struct {
		SealProofType *big.Int
		SectorNumber  *big.Int
		TicketEpoch   *big.Int
		SeedEpoch     *big.Int
		SealedCID     []byte
		UnsealedCID   []byte
		Proof         []byte
		Check         bool
		IsReal        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SealProofType = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SectorNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TicketEpoch = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.SeedEpoch = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.SealedCID = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.UnsealedCID = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	outstruct.Proof = *abi.ConvertType(out[6], new([]byte)).(*[]byte)
	outstruct.Check = *abi.ConvertType(out[7], new(bool)).(*bool)
	outstruct.IsReal = *abi.ConvertType(out[8], new(bool)).(*bool)

	return *outstruct, err

}

// ValidatorSector is a free data retrieval call binding the contract method 0x4cdf3f46.
//
// Solidity: function validatorSector(address , uint256 ) view returns(uint256 SealProofType, uint256 SectorNumber, uint256 TicketEpoch, uint256 SeedEpoch, bytes SealedCID, bytes UnsealedCID, bytes Proof, bool Check, bool isReal)
func (_W3fsStorageManager *W3fsStorageManagerSession) ValidatorSector(arg0 common.Address, arg1 *big.Int) (struct {
	SealProofType *big.Int
	SectorNumber  *big.Int
	TicketEpoch   *big.Int
	SeedEpoch     *big.Int
	SealedCID     []byte
	UnsealedCID   []byte
	Proof         []byte
	Check         bool
	IsReal        bool
}, error) {
	return _W3fsStorageManager.Contract.ValidatorSector(&_W3fsStorageManager.CallOpts, arg0, arg1)
}

// ValidatorSector is a free data retrieval call binding the contract method 0x4cdf3f46.
//
// Solidity: function validatorSector(address , uint256 ) view returns(uint256 SealProofType, uint256 SectorNumber, uint256 TicketEpoch, uint256 SeedEpoch, bytes SealedCID, bytes UnsealedCID, bytes Proof, bool Check, bool isReal)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) ValidatorSector(arg0 common.Address, arg1 *big.Int) (struct {
	SealProofType *big.Int
	SectorNumber  *big.Int
	TicketEpoch   *big.Int
	SeedEpoch     *big.Int
	SealedCID     []byte
	UnsealedCID   []byte
	Proof         []byte
	Check         bool
	IsReal        bool
}, error) {
	return _W3fsStorageManager.Contract.ValidatorSector(&_W3fsStorageManager.CallOpts, arg0, arg1)
}

// ValidatorStorageSize is a free data retrieval call binding the contract method 0xa1ca459a.
//
// Solidity: function validatorStorageSize(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCaller) ValidatorStorageSize(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _W3fsStorageManager.contract.Call(opts, &out, "validatorStorageSize", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorStorageSize is a free data retrieval call binding the contract method 0xa1ca459a.
//
// Solidity: function validatorStorageSize(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerSession) ValidatorStorageSize(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorStorageSize(&_W3fsStorageManager.CallOpts, arg0)
}

// ValidatorStorageSize is a free data retrieval call binding the contract method 0xa1ca459a.
//
// Solidity: function validatorStorageSize(address ) view returns(uint256)
func (_W3fsStorageManager *W3fsStorageManagerCallerSession) ValidatorStorageSize(arg0 common.Address) (*big.Int, error) {
	return _W3fsStorageManager.Contract.ValidatorStorageSize(&_W3fsStorageManager.CallOpts, arg0)
}

// AddSealInfo is a paid mutator transaction binding the contract method 0xd64d3e43.
//
// Solidity: function addSealInfo(bool isReal, address signer, bytes votes) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) AddSealInfo(opts *bind.TransactOpts, isReal bool, signer common.Address, votes []byte) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "addSealInfo", isReal, signer, votes)
}

// AddSealInfo is a paid mutator transaction binding the contract method 0xd64d3e43.
//
// Solidity: function addSealInfo(bool isReal, address signer, bytes votes) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) AddSealInfo(isReal bool, signer common.Address, votes []byte) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.AddSealInfo(&_W3fsStorageManager.TransactOpts, isReal, signer, votes)
}

// AddSealInfo is a paid mutator transaction binding the contract method 0xd64d3e43.
//
// Solidity: function addSealInfo(bool isReal, address signer, bytes votes) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) AddSealInfo(isReal bool, signer common.Address, votes []byte) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.AddSealInfo(&_W3fsStorageManager.TransactOpts, isReal, signer, votes)
}

// CheckSealSigs is a paid mutator transaction binding the contract method 0xaf7c33a3.
//
// Solidity: function checkSealSigs(bytes data, uint256[3][] sigs) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) CheckSealSigs(opts *bind.TransactOpts, data []byte, sigs [][3]*big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "checkSealSigs", data, sigs)
}

// CheckSealSigs is a paid mutator transaction binding the contract method 0xaf7c33a3.
//
// Solidity: function checkSealSigs(bytes data, uint256[3][] sigs) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) CheckSealSigs(data []byte, sigs [][3]*big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.CheckSealSigs(&_W3fsStorageManager.TransactOpts, data, sigs)
}

// CheckSealSigs is a paid mutator transaction binding the contract method 0xaf7c33a3.
//
// Solidity: function checkSealSigs(bytes data, uint256[3][] sigs) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) CheckSealSigs(data []byte, sigs [][3]*big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.CheckSealSigs(&_W3fsStorageManager.TransactOpts, data, sigs)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address _owner, address _registry, address _governance, address _stakingLogger) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _registry common.Address, _governance common.Address, _stakingLogger common.Address) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "initialize", _owner, _registry, _governance, _stakingLogger)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address _owner, address _registry, address _governance, address _stakingLogger) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) Initialize(_owner common.Address, _registry common.Address, _governance common.Address, _stakingLogger common.Address) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.Initialize(&_W3fsStorageManager.TransactOpts, _owner, _registry, _governance, _stakingLogger)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address _owner, address _registry, address _governance, address _stakingLogger) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) Initialize(_owner common.Address, _registry common.Address, _governance common.Address, _stakingLogger common.Address) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.Initialize(&_W3fsStorageManager.TransactOpts, _owner, _registry, _governance, _stakingLogger)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) Lock() (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.Lock(&_W3fsStorageManager.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) Lock() (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.Lock(&_W3fsStorageManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.RenounceOwnership(&_W3fsStorageManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.RenounceOwnership(&_W3fsStorageManager.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.TransferOwnership(&_W3fsStorageManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.TransferOwnership(&_W3fsStorageManager.TransactOpts, newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) Unlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "unlock")
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) Unlock() (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.Unlock(&_W3fsStorageManager.TransactOpts)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) Unlock() (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.Unlock(&_W3fsStorageManager.TransactOpts)
}

// UpdateDelegatedStakeLimit is a paid mutator transaction binding the contract method 0xfd548467.
//
// Solidity: function updateDelegatedStakeLimit(uint256 newDelegatedStakeLimit) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) UpdateDelegatedStakeLimit(opts *bind.TransactOpts, newDelegatedStakeLimit *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "updateDelegatedStakeLimit", newDelegatedStakeLimit)
}

// UpdateDelegatedStakeLimit is a paid mutator transaction binding the contract method 0xfd548467.
//
// Solidity: function updateDelegatedStakeLimit(uint256 newDelegatedStakeLimit) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) UpdateDelegatedStakeLimit(newDelegatedStakeLimit *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdateDelegatedStakeLimit(&_W3fsStorageManager.TransactOpts, newDelegatedStakeLimit)
}

// UpdateDelegatedStakeLimit is a paid mutator transaction binding the contract method 0xfd548467.
//
// Solidity: function updateDelegatedStakeLimit(uint256 newDelegatedStakeLimit) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) UpdateDelegatedStakeLimit(newDelegatedStakeLimit *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdateDelegatedStakeLimit(&_W3fsStorageManager.TransactOpts, newDelegatedStakeLimit)
}

// UpdatePercentage is a paid mutator transaction binding the contract method 0x5c2930b6.
//
// Solidity: function updatePercentage(uint256 newPercentage) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) UpdatePercentage(opts *bind.TransactOpts, newPercentage *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "updatePercentage", newPercentage)
}

// UpdatePercentage is a paid mutator transaction binding the contract method 0x5c2930b6.
//
// Solidity: function updatePercentage(uint256 newPercentage) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) UpdatePercentage(newPercentage *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdatePercentage(&_W3fsStorageManager.TransactOpts, newPercentage)
}

// UpdatePercentage is a paid mutator transaction binding the contract method 0x5c2930b6.
//
// Solidity: function updatePercentage(uint256 newPercentage) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) UpdatePercentage(newPercentage *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdatePercentage(&_W3fsStorageManager.TransactOpts, newPercentage)
}

// UpdateStakeLimit is a paid mutator transaction binding the contract method 0xdb70e8e8.
//
// Solidity: function updateStakeLimit(uint256 newStakeLimit) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) UpdateStakeLimit(opts *bind.TransactOpts, newStakeLimit *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "updateStakeLimit", newStakeLimit)
}

// UpdateStakeLimit is a paid mutator transaction binding the contract method 0xdb70e8e8.
//
// Solidity: function updateStakeLimit(uint256 newStakeLimit) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) UpdateStakeLimit(newStakeLimit *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdateStakeLimit(&_W3fsStorageManager.TransactOpts, newStakeLimit)
}

// UpdateStakeLimit is a paid mutator transaction binding the contract method 0xdb70e8e8.
//
// Solidity: function updateStakeLimit(uint256 newStakeLimit) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) UpdateStakeLimit(newStakeLimit *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdateStakeLimit(&_W3fsStorageManager.TransactOpts, newStakeLimit)
}

// UpdateStoragePromise is a paid mutator transaction binding the contract method 0xcfacc8dc.
//
// Solidity: function updateStoragePromise(address signer, uint256 storageSize) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactor) UpdateStoragePromise(opts *bind.TransactOpts, signer common.Address, storageSize *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.contract.Transact(opts, "updateStoragePromise", signer, storageSize)
}

// UpdateStoragePromise is a paid mutator transaction binding the contract method 0xcfacc8dc.
//
// Solidity: function updateStoragePromise(address signer, uint256 storageSize) returns()
func (_W3fsStorageManager *W3fsStorageManagerSession) UpdateStoragePromise(signer common.Address, storageSize *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdateStoragePromise(&_W3fsStorageManager.TransactOpts, signer, storageSize)
}

// UpdateStoragePromise is a paid mutator transaction binding the contract method 0xcfacc8dc.
//
// Solidity: function updateStoragePromise(address signer, uint256 storageSize) returns()
func (_W3fsStorageManager *W3fsStorageManagerTransactorSession) UpdateStoragePromise(signer common.Address, storageSize *big.Int) (*types.Transaction, error) {
	return _W3fsStorageManager.Contract.UpdateStoragePromise(&_W3fsStorageManager.TransactOpts, signer, storageSize)
}

// W3fsStorageManagerAddNewSealPowerAndSizeIterator is returned from FilterAddNewSealPowerAndSize and is used to iterate over the raw logs and unpacked data for AddNewSealPowerAndSize events raised by the W3fsStorageManager contract.
type W3fsStorageManagerAddNewSealPowerAndSizeIterator struct {
	Event *W3fsStorageManagerAddNewSealPowerAndSize // Event containing the contract specifics and raw log

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
func (it *W3fsStorageManagerAddNewSealPowerAndSizeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(W3fsStorageManagerAddNewSealPowerAndSize)
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
		it.Event = new(W3fsStorageManagerAddNewSealPowerAndSize)
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
func (it *W3fsStorageManagerAddNewSealPowerAndSizeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *W3fsStorageManagerAddNewSealPowerAndSizeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// W3fsStorageManagerAddNewSealPowerAndSize represents a AddNewSealPowerAndSize event raised by the W3fsStorageManager contract.
type W3fsStorageManagerAddNewSealPowerAndSize struct {
	NewPower       *big.Int
	NewStorageSize *big.Int
	Signer         common.Address
	Nonce          *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAddNewSealPowerAndSize is a free log retrieval operation binding the contract event 0x9720a8bb36ec795699396855eeed0f102e16c5934b62691884b8e1e898772d0d.
//
// Solidity: event AddNewSealPowerAndSize(uint256 indexed newPower, uint256 indexed newStorageSize, address indexed signer, uint256 nonce)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) FilterAddNewSealPowerAndSize(opts *bind.FilterOpts, newPower []*big.Int, newStorageSize []*big.Int, signer []common.Address) (*W3fsStorageManagerAddNewSealPowerAndSizeIterator, error) {

	var newPowerRule []interface{}
	for _, newPowerItem := range newPower {
		newPowerRule = append(newPowerRule, newPowerItem)
	}
	var newStorageSizeRule []interface{}
	for _, newStorageSizeItem := range newStorageSize {
		newStorageSizeRule = append(newStorageSizeRule, newStorageSizeItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _W3fsStorageManager.contract.FilterLogs(opts, "AddNewSealPowerAndSize", newPowerRule, newStorageSizeRule, signerRule)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManagerAddNewSealPowerAndSizeIterator{contract: _W3fsStorageManager.contract, event: "AddNewSealPowerAndSize", logs: logs, sub: sub}, nil
}

// WatchAddNewSealPowerAndSize is a free log subscription operation binding the contract event 0x9720a8bb36ec795699396855eeed0f102e16c5934b62691884b8e1e898772d0d.
//
// Solidity: event AddNewSealPowerAndSize(uint256 indexed newPower, uint256 indexed newStorageSize, address indexed signer, uint256 nonce)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) WatchAddNewSealPowerAndSize(opts *bind.WatchOpts, sink chan<- *W3fsStorageManagerAddNewSealPowerAndSize, newPower []*big.Int, newStorageSize []*big.Int, signer []common.Address) (event.Subscription, error) {

	var newPowerRule []interface{}
	for _, newPowerItem := range newPower {
		newPowerRule = append(newPowerRule, newPowerItem)
	}
	var newStorageSizeRule []interface{}
	for _, newStorageSizeItem := range newStorageSize {
		newStorageSizeRule = append(newStorageSizeRule, newStorageSizeItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _W3fsStorageManager.contract.WatchLogs(opts, "AddNewSealPowerAndSize", newPowerRule, newStorageSizeRule, signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(W3fsStorageManagerAddNewSealPowerAndSize)
				if err := _W3fsStorageManager.contract.UnpackLog(event, "AddNewSealPowerAndSize", log); err != nil {
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

// ParseAddNewSealPowerAndSize is a log parse operation binding the contract event 0x9720a8bb36ec795699396855eeed0f102e16c5934b62691884b8e1e898772d0d.
//
// Solidity: event AddNewSealPowerAndSize(uint256 indexed newPower, uint256 indexed newStorageSize, address indexed signer, uint256 nonce)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) ParseAddNewSealPowerAndSize(log types.Log) (*W3fsStorageManagerAddNewSealPowerAndSize, error) {
	event := new(W3fsStorageManagerAddNewSealPowerAndSize)
	if err := _W3fsStorageManager.contract.UnpackLog(event, "AddNewSealPowerAndSize", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// W3fsStorageManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the W3fsStorageManager contract.
type W3fsStorageManagerOwnershipTransferredIterator struct {
	Event *W3fsStorageManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *W3fsStorageManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(W3fsStorageManagerOwnershipTransferred)
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
		it.Event = new(W3fsStorageManagerOwnershipTransferred)
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
func (it *W3fsStorageManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *W3fsStorageManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// W3fsStorageManagerOwnershipTransferred represents a OwnershipTransferred event raised by the W3fsStorageManager contract.
type W3fsStorageManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*W3fsStorageManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _W3fsStorageManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManagerOwnershipTransferredIterator{contract: _W3fsStorageManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *W3fsStorageManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _W3fsStorageManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(W3fsStorageManagerOwnershipTransferred)
				if err := _W3fsStorageManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) ParseOwnershipTransferred(log types.Log) (*W3fsStorageManagerOwnershipTransferred, error) {
	event := new(W3fsStorageManagerOwnershipTransferred)
	if err := _W3fsStorageManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// W3fsStorageManagerUpdatePromiseIterator is returned from FilterUpdatePromise and is used to iterate over the raw logs and unpacked data for UpdatePromise events raised by the W3fsStorageManager contract.
type W3fsStorageManagerUpdatePromiseIterator struct {
	Event *W3fsStorageManagerUpdatePromise // Event containing the contract specifics and raw log

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
func (it *W3fsStorageManagerUpdatePromiseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(W3fsStorageManagerUpdatePromise)
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
		it.Event = new(W3fsStorageManagerUpdatePromise)
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
func (it *W3fsStorageManagerUpdatePromiseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *W3fsStorageManagerUpdatePromiseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// W3fsStorageManagerUpdatePromise represents a UpdatePromise event raised by the W3fsStorageManager contract.
type W3fsStorageManagerUpdatePromise struct {
	Signer      common.Address
	StorageSize *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdatePromise is a free log retrieval operation binding the contract event 0xc1c594ec6d2acc15c3f7a06152f646bcf01e56ee0193482de75c51e606c7c1ca.
//
// Solidity: event UpdatePromise(address indexed signer, uint256 indexed storageSize)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) FilterUpdatePromise(opts *bind.FilterOpts, signer []common.Address, storageSize []*big.Int) (*W3fsStorageManagerUpdatePromiseIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var storageSizeRule []interface{}
	for _, storageSizeItem := range storageSize {
		storageSizeRule = append(storageSizeRule, storageSizeItem)
	}

	logs, sub, err := _W3fsStorageManager.contract.FilterLogs(opts, "UpdatePromise", signerRule, storageSizeRule)
	if err != nil {
		return nil, err
	}
	return &W3fsStorageManagerUpdatePromiseIterator{contract: _W3fsStorageManager.contract, event: "UpdatePromise", logs: logs, sub: sub}, nil
}

// WatchUpdatePromise is a free log subscription operation binding the contract event 0xc1c594ec6d2acc15c3f7a06152f646bcf01e56ee0193482de75c51e606c7c1ca.
//
// Solidity: event UpdatePromise(address indexed signer, uint256 indexed storageSize)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) WatchUpdatePromise(opts *bind.WatchOpts, sink chan<- *W3fsStorageManagerUpdatePromise, signer []common.Address, storageSize []*big.Int) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var storageSizeRule []interface{}
	for _, storageSizeItem := range storageSize {
		storageSizeRule = append(storageSizeRule, storageSizeItem)
	}

	logs, sub, err := _W3fsStorageManager.contract.WatchLogs(opts, "UpdatePromise", signerRule, storageSizeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(W3fsStorageManagerUpdatePromise)
				if err := _W3fsStorageManager.contract.UnpackLog(event, "UpdatePromise", log); err != nil {
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

// ParseUpdatePromise is a log parse operation binding the contract event 0xc1c594ec6d2acc15c3f7a06152f646bcf01e56ee0193482de75c51e606c7c1ca.
//
// Solidity: event UpdatePromise(address indexed signer, uint256 indexed storageSize)
func (_W3fsStorageManager *W3fsStorageManagerFilterer) ParseUpdatePromise(log types.Log) (*W3fsStorageManagerUpdatePromise, error) {
	event := new(W3fsStorageManagerUpdatePromise)
	if err := _W3fsStorageManager.contract.UnpackLog(event, "UpdatePromise", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}