// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package checktxhash

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

// VoucherMetaData contains all meta data concerning the Voucher contract.
var VoucherMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"oriHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"purchaseType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"userFilePubkey\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validityDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"createDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"extraInfo\",\"type\":\"string\"}],\"name\":\"PurchaseVoucher\",\"type\":\"event\"}]",
}

// VoucherABI is the input ABI used to generate the binding from.
// Deprecated: Use VoucherMetaData.ABI instead.
var VoucherABI = VoucherMetaData.ABI

// Voucher is an auto generated Go binding around an Ethereum contract.
type Voucher struct {
	VoucherCaller     // Read-only binding to the contract
	VoucherTransactor // Write-only binding to the contract
	VoucherFilterer   // Log filterer for contract events
}

// VoucherCaller is an auto generated read-only Go binding around an Ethereum contract.
type VoucherCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VoucherTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VoucherTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VoucherFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VoucherFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VoucherSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VoucherSession struct {
	Contract     *Voucher          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VoucherCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VoucherCallerSession struct {
	Contract *VoucherCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// VoucherTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VoucherTransactorSession struct {
	Contract     *VoucherTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// VoucherRaw is an auto generated low-level Go binding around an Ethereum contract.
type VoucherRaw struct {
	Contract *Voucher // Generic contract binding to access the raw methods on
}

// VoucherCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VoucherCallerRaw struct {
	Contract *VoucherCaller // Generic read-only contract binding to access the raw methods on
}

// VoucherTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VoucherTransactorRaw struct {
	Contract *VoucherTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVoucher creates a new instance of Voucher, bound to a specific deployed contract.
func NewVoucher(address common.Address, backend bind.ContractBackend) (*Voucher, error) {
	contract, err := bindVoucher(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Voucher{VoucherCaller: VoucherCaller{contract: contract}, VoucherTransactor: VoucherTransactor{contract: contract}, VoucherFilterer: VoucherFilterer{contract: contract}}, nil
}

// NewVoucherCaller creates a new read-only instance of Voucher, bound to a specific deployed contract.
func NewVoucherCaller(address common.Address, caller bind.ContractCaller) (*VoucherCaller, error) {
	contract, err := bindVoucher(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VoucherCaller{contract: contract}, nil
}

// NewVoucherTransactor creates a new write-only instance of Voucher, bound to a specific deployed contract.
func NewVoucherTransactor(address common.Address, transactor bind.ContractTransactor) (*VoucherTransactor, error) {
	contract, err := bindVoucher(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VoucherTransactor{contract: contract}, nil
}

// NewVoucherFilterer creates a new log filterer instance of Voucher, bound to a specific deployed contract.
func NewVoucherFilterer(address common.Address, filterer bind.ContractFilterer) (*VoucherFilterer, error) {
	contract, err := bindVoucher(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VoucherFilterer{contract: contract}, nil
}

// bindVoucher binds a generic wrapper to an already deployed contract.
func bindVoucher(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VoucherABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Voucher *VoucherRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Voucher.Contract.VoucherCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Voucher *VoucherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voucher.Contract.VoucherTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Voucher *VoucherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Voucher.Contract.VoucherTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Voucher *VoucherCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Voucher.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Voucher *VoucherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voucher.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Voucher *VoucherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Voucher.Contract.contract.Transact(opts, method, params...)
}

// VoucherPurchaseVoucherIterator is returned from FilterPurchaseVoucher and is used to iterate over the raw logs and unpacked data for PurchaseVoucher events raised by the Voucher contract.
type VoucherPurchaseVoucherIterator struct {
	Event *VoucherPurchaseVoucher // Event containing the contract specifics and raw log

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
func (it *VoucherPurchaseVoucherIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VoucherPurchaseVoucher)
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
		it.Event = new(VoucherPurchaseVoucher)
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
func (it *VoucherPurchaseVoucherIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VoucherPurchaseVoucherIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VoucherPurchaseVoucher represents a PurchaseVoucher event raised by the Voucher contract.
type VoucherPurchaseVoucher struct {
	OriHash        [32]byte
	Consumer       common.Address
	PurchaseType   uint8
	UserFilePubkey string
	ValidityDate   *big.Int
	CreateDate     *big.Int
	ExtraInfo      string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPurchaseVoucher is a free log retrieval operation binding the contract event 0x6306ba057289db4ff93c98f6af37b408120f2adb8fb1f13bcb4f0a15a3ffcbe8.
//
// Solidity: event PurchaseVoucher(bytes32 oriHash, address consumer, uint8 purchaseType, string userFilePubkey, uint256 validityDate, uint256 createDate, string extraInfo)
func (_Voucher *VoucherFilterer) FilterPurchaseVoucher(opts *bind.FilterOpts) (*VoucherPurchaseVoucherIterator, error) {

	logs, sub, err := _Voucher.contract.FilterLogs(opts, "PurchaseVoucher")
	if err != nil {
		return nil, err
	}
	return &VoucherPurchaseVoucherIterator{contract: _Voucher.contract, event: "PurchaseVoucher", logs: logs, sub: sub}, nil
}

// WatchPurchaseVoucher is a free log subscription operation binding the contract event 0x6306ba057289db4ff93c98f6af37b408120f2adb8fb1f13bcb4f0a15a3ffcbe8.
//
// Solidity: event PurchaseVoucher(bytes32 oriHash, address consumer, uint8 purchaseType, string userFilePubkey, uint256 validityDate, uint256 createDate, string extraInfo)
func (_Voucher *VoucherFilterer) WatchPurchaseVoucher(opts *bind.WatchOpts, sink chan<- *VoucherPurchaseVoucher) (event.Subscription, error) {

	logs, sub, err := _Voucher.contract.WatchLogs(opts, "PurchaseVoucher")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VoucherPurchaseVoucher)
				if err := _Voucher.contract.UnpackLog(event, "PurchaseVoucher", log); err != nil {
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

// ParsePurchaseVoucher is a log parse operation binding the contract event 0x6306ba057289db4ff93c98f6af37b408120f2adb8fb1f13bcb4f0a15a3ffcbe8.
//
// Solidity: event PurchaseVoucher(bytes32 oriHash, address consumer, uint8 purchaseType, string userFilePubkey, uint256 validityDate, uint256 createDate, string extraInfo)
func (_Voucher *VoucherFilterer) ParsePurchaseVoucher(log types.Log) (*VoucherPurchaseVoucher, error) {
	event := new(VoucherPurchaseVoucher)
	if err := _Voucher.contract.UnpackLog(event, "PurchaseVoucher", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
