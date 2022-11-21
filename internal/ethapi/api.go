// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethapi

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/docker/go-units"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	file_store "github.com/ethereum/go-ethereum/borcontracts/file-store"
	"github.com/ethereum/go-ethereum/borcontracts/w3fsStorageManager"
	"github.com/ethereum/go-ethereum/extern/gcache"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/fatih/color"
	"github.com/filecoin-project/go-address"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-jsonrpc/auth"
	bigext "github.com/filecoin-project/go-state-types/big"
	lapi "github.com/filecoin-project/lotus/api"
	chaintypes "github.com/filecoin-project/lotus/chain/types"
	lcli "github.com/filecoin-project/lotus/cli"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	"github.com/filecoin-project/lotus/extern/sector-storage/fsutil"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	lsealing "github.com/filecoin-project/lotus/extern/storage-sealing"
	"github.com/filecoin-project/lotus/lib/lotuslog"
	"github.com/filecoin-project/lotus/lib/tablewriter"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/accounts/scwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	fabi "github.com/filecoin-project/go-state-types/abi"
	logging "github.com/ipfs/go-log/v2"
	"github.com/tyler-smith/go-bip39"
)

var lotusLog = logging.Logger("ethapi")

// the are for lotus p2p
var PeerId = ""

const ENTIRE_FILE = "ENTIRE"
const STORED = "stored"
const UNSTORED = "unstored"
const BOTH = "both"

// PublicEthereumAPI provides an API to access Ethereum related information.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicEthereumAPI struct {
	b Backend
}

// NewPublicEthereumAPI creates a new Ethereum protocol API.
func NewPublicEthereumAPI(b Backend) *PublicEthereumAPI {
	return &PublicEthereumAPI{b}
}

// GasPrice returns a suggestion for a gas price for legacy transactions.
func (s *PublicEthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	tipcap, err := s.b.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	if head := s.b.CurrentHeader(); head.BaseFee != nil {
		tipcap.Add(tipcap, head.BaseFee)
	}
	return (*hexutil.Big)(tipcap), err
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (s *PublicEthereumAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	tipcap, err := s.b.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(tipcap), err
}

type feeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

func (s *PublicEthereumAPI) FeeHistory(ctx context.Context, blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*feeHistoryResult, error) {
	oldest, reward, baseFee, gasUsed, err := s.b.FeeHistory(ctx, int(blockCount), lastBlock, rewardPercentiles)
	if err != nil {
		return nil, err
	}
	results := &feeHistoryResult{
		OldestBlock:  (*hexutil.Big)(oldest),
		GasUsedRatio: gasUsed,
	}
	if reward != nil {
		results.Reward = make([][]*hexutil.Big, len(reward))
		for i, w := range reward {
			results.Reward[i] = make([]*hexutil.Big, len(w))
			for j, v := range w {
				results.Reward[i][j] = (*hexutil.Big)(v)
			}
		}
	}
	if baseFee != nil {
		results.BaseFee = make([]*hexutil.Big, len(baseFee))
		for i, v := range baseFee {
			results.BaseFee[i] = (*hexutil.Big)(v)
		}
	}
	return results, nil
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *PublicEthereumAPI) Syncing() (interface{}, error) {
	progress := s.b.Downloader().Progress()

	// Return not syncing if the synchronisation already completed
	if progress.CurrentBlock >= progress.HighestBlock {
		return false, nil
	}
	// Otherwise gather the block sync stats
	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
	}, nil
}

// PublicTxPoolAPI offers and API for the transaction pool. It only operates on data that is non confidential.
type PublicTxPoolAPI struct {
	b Backend
}

// NewPublicTxPoolAPI creates a new tx pool service that gives information about the transaction pool.
func NewPublicTxPoolAPI(b Backend) *PublicTxPoolAPI {
	return &PublicTxPoolAPI{b}
}

// Content returns the transactions contained within the transaction pool.
func (s *PublicTxPoolAPI) Content() map[string]map[string]map[string]*RPCTransaction {
	content := map[string]map[string]map[string]*RPCTransaction{
		"pending": make(map[string]map[string]*RPCTransaction),
		"queued":  make(map[string]map[string]*RPCTransaction),
	}
	pending, queue := s.b.TxPoolContent()
	curHeader := s.b.CurrentHeader()
	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
		}
		content["pending"][account.Hex()] = dump
	}
	// Flatten the queued transactions
	for account, txs := range queue {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}

// ContentFrom returns the transactions contained within the transaction pool.
func (s *PublicTxPoolAPI) ContentFrom(addr common.Address) map[string]map[string]*RPCTransaction {
	content := make(map[string]map[string]*RPCTransaction, 2)
	pending, queue := s.b.TxPoolContentFrom(addr)
	curHeader := s.b.CurrentHeader()

	// Build the pending transactions
	dump := make(map[string]*RPCTransaction, len(pending))
	for _, tx := range pending {
		dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
	}
	content["pending"] = dump

	// Build the queued transactions
	dump = make(map[string]*RPCTransaction, len(queue))
	for _, tx := range queue {
		dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
	}
	content["queued"] = dump

	return content
}

// Status returns the number of pending and queued transaction in the pool.
func (s *PublicTxPoolAPI) Status() map[string]hexutil.Uint {
	pending, queue := s.b.Stats()
	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(pending),
		"queued":  hexutil.Uint(queue),
	}
}

// Inspect retrieves the content of the transaction pool and flattens it into an
// easily inspectable list.
func (s *PublicTxPoolAPI) Inspect() map[string]map[string]map[string]string {
	content := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}
	pending, queue := s.b.TxPoolContent()

	// Define a formatter to flatten a transaction into a string
	var format = func(tx *types.Transaction) string {
		if to := tx.To(); to != nil {
			return fmt.Sprintf("%s: %v wei + %v gas × %v wei", tx.To().Hex(), tx.Value(), tx.Gas(), tx.GasPrice())
		}
		return fmt.Sprintf("contract creation: %v wei + %v gas × %v wei", tx.Value(), tx.Gas(), tx.GasPrice())
	}
	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	// Flatten the queued transactions
	for account, txs := range queue {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}

// PublicAccountAPI provides an API to access accounts managed by this node.
// It offers only methods that can retrieve accounts.
type PublicAccountAPI struct {
	am *accounts.Manager
}

// NewPublicAccountAPI creates a new PublicAccountAPI.
func NewPublicAccountAPI(am *accounts.Manager) *PublicAccountAPI {
	return &PublicAccountAPI{am: am}
}

// Accounts returns the collection of accounts this node manages
func (s *PublicAccountAPI) Accounts() []common.Address {
	return s.am.Accounts()
}

// PrivateAccountAPI provides an API to access accounts managed by this node.
// It offers methods to create, (un)lock en list accounts. Some methods accept
// passwords and are therefore considered private by default.
type PrivateAccountAPI struct {
	am        *accounts.Manager
	nonceLock *AddrLocker
	b         Backend
}

// NewPrivateAccountAPI create a new PrivateAccountAPI.
func NewPrivateAccountAPI(b Backend, nonceLock *AddrLocker) *PrivateAccountAPI {
	return &PrivateAccountAPI{
		am:        b.AccountManager(),
		nonceLock: nonceLock,
		b:         b,
	}
}

// listAccounts will return a list of addresses for accounts this node manages.
func (s *PrivateAccountAPI) ListAccounts() []common.Address {
	return s.am.Accounts()
}

// rawWallet is a JSON representation of an accounts.Wallet interface, with its
// data contents extracted into plain fields.
type rawWallet struct {
	URL      string             `json:"url"`
	Status   string             `json:"status"`
	Failure  string             `json:"failure,omitempty"`
	Accounts []accounts.Account `json:"accounts,omitempty"`
}

// ListWallets will return a list of wallets this node manages.
func (s *PrivateAccountAPI) ListWallets() []rawWallet {
	wallets := make([]rawWallet, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		status, failure := wallet.Status()

		raw := rawWallet{
			URL:      wallet.URL().String(),
			Status:   status,
			Accounts: wallet.Accounts(),
		}
		if failure != nil {
			raw.Failure = failure.Error()
		}
		wallets = append(wallets, raw)
	}
	return wallets
}

// OpenWallet initiates a hardware wallet opening procedure, establishing a USB
// connection and attempting to authenticate via the provided passphrase. Note,
// the method may return an extra challenge requiring a second open (e.g. the
// Trezor PIN matrix challenge).
func (s *PrivateAccountAPI) OpenWallet(url string, passphrase *string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}
	pass := ""
	if passphrase != nil {
		pass = *passphrase
	}
	return wallet.Open(pass)
}

// DeriveAccount requests a HD wallet to derive a new account, optionally pinning
// it for later reuse.
func (s *PrivateAccountAPI) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return accounts.Account{}, err
	}
	derivPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.Account{}, err
	}
	if pin == nil {
		pin = new(bool)
	}
	return wallet.Derive(derivPath, *pin)
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI) NewAccount(password string) (common.Address, error) {
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := ks.NewAccount(password)
	if err == nil {
		log.Info("Your new key was generated", "address", acc.Address)
		log.Warn("Please backup your key file!", "path", acc.URL.Path)
		log.Warn("Please remember your password!")
		return acc.Address, nil
	}
	return common.Address{}, err
}

// fetchKeystore retrieves the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) (*keystore.KeyStore, error) {
	if ks := am.Backends(keystore.KeyStoreType); len(ks) > 0 {
		return ks[0].(*keystore.KeyStore), nil
	}
	return nil, errors.New("local keystore not used")
}

// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (s *PrivateAccountAPI) ImportRawKey(privkey string, password string) (common.Address, error) {
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := ks.ImportECDSA(key, password)
	return acc.Address, err
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (s *PrivateAccountAPI) UnlockAccount(ctx context.Context, addr common.Address, password string, duration *uint64) (bool, error) {
	// When the API is exposed by external RPC(http, ws etc), unless the user
	// explicitly specifies to allow the insecure account unlocking, otherwise
	// it is disabled.
	if s.b.ExtRPCEnabled() && !s.b.AccountManager().Config().InsecureUnlockAllowed {
		return false, errors.New("account unlock with HTTP access is forbidden")
	}

	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return false, err
	}
	err = ks.TimedUnlock(accounts.Account{Address: addr}, password, d)
	if err != nil {
		log.Warn("Failed account unlock attempt", "address", addr, "err", err)
	}
	return err == nil, err
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PrivateAccountAPI) LockAccount(addr common.Address) bool {
	if ks, err := fetchKeystore(s.am); err == nil {
		return ks.Lock(addr) == nil
	}
	return false
}

// signTransaction sets defaults and signs the given transaction
// NOTE: the caller needs to ensure that the nonceLock is held, if applicable,
// and release it after the transaction has been submitted to the tx pool
func (s *PrivateAccountAPI) signTransaction(ctx context.Context, args *TransactionArgs, passwd string) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: args.from()}
	wallet, err := s.am.Find(account)
	if err != nil {
		return nil, err
	}
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and sign with the wallet
	tx := args.toTransaction()

	return wallet.SignTxWithPassphrase(account, passwd, tx, s.b.ChainConfig().ChainID)
}

// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given
// passwd isn't able to decrypt the key it fails.
func (s *PrivateAccountAPI) SendTransaction(ctx context.Context, args TransactionArgs, passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.from())
		defer s.nonceLock.UnlockAddr(args.from())
	}
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.from(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, signed)
}

// SignTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given passwd isn't
// able to decrypt the key it fails. The transaction is returned in RLP-form, not broadcast
// to other nodes
func (s *PrivateAccountAPI) SignTransaction(ctx context.Context, args TransactionArgs, passwd string) (*SignTransactionResult, error) {
	// No need to obtain the noncelock mutex, since we won't be sending this
	// tx into the transaction pool, but right back to the user
	if args.From == nil {
		return nil, fmt.Errorf("sender not specified")
	}
	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil && (args.MaxFeePerGas == nil || args.MaxPriorityFeePerGas == nil) {
		return nil, fmt.Errorf("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	// Before actually signing the transaction, ensure the transaction fee is reasonable.
	tx := args.toTransaction()
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction sign attempt", "from", args.from(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}

// Sign calculates an Ethereum ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func (s *PrivateAccountAPI) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Assemble sign the data with the wallet
	signature, err := wallet.SignTextWithPassphrase(account, passwd, data)
	if err != nil {
		log.Warn("Failed data sign attempt", "address", addr, "err", err)
		return nil, err
	}
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be 27 or 28 for legacy reasons.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecover
func (s *PrivateAccountAPI) EcRecover(ctx context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.SigToPub(accounts.TextHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

// SignAndSendTransaction was renamed to SendTransaction. This method is deprecated
// and will be removed in the future. It primary goal is to give clients time to update.
func (s *PrivateAccountAPI) SignAndSendTransaction(ctx context.Context, args TransactionArgs, passwd string) (common.Hash, error) {
	return s.SendTransaction(ctx, args, passwd)
}

// InitializeWallet initializes a new wallet at the provided URL, by generating and returning a new private key.
func (s *PrivateAccountAPI) InitializeWallet(ctx context.Context, url string) (string, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return "", err
	}

	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	seed := bip39.NewSeed(mnemonic, "")

	switch wallet := wallet.(type) {
	case *scwallet.Wallet:
		return mnemonic, wallet.Initialize(seed)
	default:
		return "", fmt.Errorf("specified wallet does not support initialization")
	}
}

// Unpair deletes a pairing between wallet and geth.
func (s *PrivateAccountAPI) Unpair(ctx context.Context, url string, pin string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}

	switch wallet := wallet.(type) {
	case *scwallet.Wallet:
		return wallet.Unpair([]byte(pin))
	default:
		return fmt.Errorf("specified wallet does not support pairing")
	}
}

// PublicBlockChainAPI provides an API to access the Ethereum blockchain.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicBlockChainAPI struct {
	b Backend
}

// NewPublicBlockChainAPI creates a new Ethereum blockchain API.
func NewPublicBlockChainAPI(b Backend) *PublicBlockChainAPI {
	return &PublicBlockChainAPI{b}
}

// GetTransactionReceiptsByBlock returns the transaction receipts for the given block number or hash.
func (s *PublicBlockChainAPI) GetTransactionReceiptsByBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	block, err := s.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}

	receipts, err := s.b.GetReceipts(ctx, block.Hash())
	if err != nil {
		return nil, err
	}

	txs := block.Transactions()

	var txHash common.Hash

	borReceipt := rawdb.ReadBorReceipt(s.b.ChainDb(), block.Hash(), block.NumberU64())
	if borReceipt != nil {
		receipts = append(receipts, borReceipt)
		txHash = types.GetDerivedBorTxHash(types.BorReceiptKey(block.Number().Uint64(), block.Hash()))
		if txHash != (common.Hash{}) {
			borTx, _, _, _, _ := s.b.GetBorBlockTransactionWithBlockHash(ctx, txHash, block.Hash())
			txs = append(txs, borTx)
		}
	}

	if len(txs) != len(receipts) {
		return nil, fmt.Errorf("txs length doesn't equal to receipts' length", len(txs), len(receipts))
	}

	txReceipts := make([]map[string]interface{}, 0, len(txs))
	for idx, receipt := range receipts {
		tx := txs[idx]
		var signer types.Signer = types.FrontierSigner{}
		if tx.Protected() {
			signer = types.NewEIP155Signer(tx.ChainId())
		}
		from, _ := types.Sender(signer, tx)

		fields := map[string]interface{}{
			"blockHash":         block.Hash(),
			"blockNumber":       hexutil.Uint64(block.NumberU64()),
			"transactionHash":   tx.Hash(),
			"transactionIndex":  hexutil.Uint64(idx),
			"from":              from,
			"to":                tx.To(),
			"gasUsed":           hexutil.Uint64(receipt.GasUsed),
			"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
			"contractAddress":   nil,
			"logs":              receipt.Logs,
			"logsBloom":         receipt.Bloom,
		}

		// Assign receipt status or post state.
		if len(receipt.PostState) > 0 {
			fields["root"] = hexutil.Bytes(receipt.PostState)
		} else {
			fields["status"] = hexutil.Uint(receipt.Status)
		}
		if receipt.Logs == nil {
			fields["logs"] = [][]*types.Log{}
		}
		if borReceipt != nil {
			fields["transactionHash"] = txHash
		}
		// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
		if receipt.ContractAddress != (common.Address{}) {
			fields["contractAddress"] = receipt.ContractAddress
		}

		txReceipts = append(txReceipts, fields)
	}

	return txReceipts, nil
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (api *PublicBlockChainAPI) ChainId() (*hexutil.Big, error) {
	// if current block is at or past the EIP-155 replay-protection fork block, return chainID from config
	if config := api.b.ChainConfig(); config.IsEIP155(api.b.CurrentBlock().Number()) {
		return (*hexutil.Big)(config.ChainID), nil
	}
	return nil, fmt.Errorf("chain not synced beyond EIP-155 replay-protection fork block")
}

// BlockNumber returns the block number of the chain head.
func (s *PublicBlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	return hexutil.Uint64(header.Number.Uint64())
}

// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (s *PublicBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}

// Result structs for GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// GetProof returns the Merkle-proof for a given account and optionally some storage keys.
func (s *PublicBlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}

	storageTrie := state.StorageTrie(address)
	storageHash := types.EmptyRootHash
	codeHash := state.GetCodeHash(address)
	storageProof := make([]StorageResult, len(storageKeys))

	// if we have a storageTrie, (which means the account exists), we can update the storagehash
	if storageTrie != nil {
		storageHash = storageTrie.Hash()
	} else {
		// no storageTrie means the account does not exist, so the codeHash is the hash of an empty bytearray.
		codeHash = crypto.Keccak256Hash(nil)
	}

	// create the proof for the storageKeys
	for i, key := range storageKeys {
		if storageTrie != nil {
			proof, storageError := state.GetStorageProof(address, common.HexToHash(key))
			if storageError != nil {
				return nil, storageError
			}
			storageProof[i] = StorageResult{key, (*hexutil.Big)(state.GetState(address, common.HexToHash(key)).Big()), toHexSlice(proof)}
		} else {
			storageProof[i] = StorageResult{key, &hexutil.Big{}, []string{}}
		}
	}

	// create the accountProof
	accountProof, proofErr := state.GetProof(address)
	if proofErr != nil {
		return nil, proofErr
	}

	return &AccountResult{
		Address:      address,
		AccountProof: toHexSlice(accountProof),
		Balance:      (*hexutil.Big)(state.GetBalance(address)),
		CodeHash:     codeHash,
		Nonce:        hexutil.Uint64(state.GetNonce(address)),
		StorageHash:  storageHash,
		StorageProof: storageProof,
	}, state.Error()
}

// GetHeaderByNumber returns the requested canonical block header.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
func (s *PublicBlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, err := s.b.HeaderByNumber(ctx, number)
	if header != nil && err == nil {
		response := s.rpcMarshalHeader(ctx, header)
		if number == rpc.PendingBlockNumber {
			// Pending header need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetHeaderByHash returns the requested header by hash.
func (s *PublicBlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	header, _ := s.b.HeaderByHash(ctx, hash)
	if header != nil {
		return s.rpcMarshalHeader(ctx, header)
	}
	return nil
}

// GetBlockByNumber returns the requested canonical block.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
// * When fullTx is true all transactions in the block are returned, otherwise
//   only the transaction hash is returned.
func (s *PublicBlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, number)
	if block != nil && err == nil {
		response, err := s.rpcMarshalBlock(ctx, block, true, fullTx)
		if err == nil && number == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}

		// append marshalled bor transaction
		if err == nil && response != nil {
			response = s.appendRPCMarshalBorTransaction(ctx, block, response, fullTx)
		}

		return response, err
	}
	return nil, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, hash)
	if block != nil {
		response, err := s.rpcMarshalBlock(ctx, block, true, fullTx)
		// append marshalled bor transaction
		if err == nil && response != nil {
			return s.appendRPCMarshalBorTransaction(ctx, block, response, fullTx), err
		}
		return response, err
	}
	return nil, err
}

// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", blockNr, "hash", block.Hash(), "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}

// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", block.Number(), "hash", blockHash, "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}

// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number
func (s *PublicBlockChainAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash
func (s *PublicBlockChainAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetCode returns the code stored at the given address in the state for the given block number.
func (s *PublicBlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (s *PublicBlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	res := state.GetState(address, common.HexToHash(key))
	return res[:], state.Error()
}

// OverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[common.Address]OverrideAccount

// Apply overrides the fields of specified accounts into the given state.
func (diff *StateOverride) Apply(state *state.StateDB) error {
	if diff == nil {
		return nil
	}
	for addr, account := range *diff {
		// Override account nonce.
		if account.Nonce != nil {
			state.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			state.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			state.SetBalance(addr, (*big.Int)(*account.Balance))
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			state.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				state.SetState(addr, key, value)
			}
		}
	}
	return nil
}

func DoCall(ctx context.Context, b Backend, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, timeout time.Duration, globalGasCap uint64) (*core.ExecutionResult, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	if err := overrides.Apply(state); err != nil {
		return nil, err
	}
	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Get a new instance of the EVM.
	msg, err := args.ToMessage(globalGasCap, header.BaseFee)
	if err != nil {
		return nil, err
	}
	evm, vmError, err := b.GetEVM(ctx, msg, state, header, &vm.Config{NoBaseFee: true})
	if err != nil {
		return nil, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()

	// Execute the message.
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	result, err := core.ApplyMessage(evm, msg, gp)
	if err := vmError(); err != nil {
		return nil, err
	}

	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		return nil, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}
	if err != nil {
		return result, fmt.Errorf("err: %w (supplied gas %d)", err, msg.Gas())
	}
	return result, nil
}

func newRevertError(result *core.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompassas an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

// Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.
func (s *PublicBlockChainAPI) Call(ctx context.Context, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Bytes, error) {
	result, err := DoCall(ctx, s.b, args, blockNrOrHash, overrides, 10*time.Second, s.b.RPCGasCap())
	if err != nil {
		return nil, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return nil, newRevertError(result)
	}
	return result.Return(), result.Err
}

func DoEstimateGas(ctx context.Context, b Backend, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, gasCap uint64) (hexutil.Uint64, error) {
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	// Use zero address if sender unspecified.
	if args.From == nil {
		args.From = new(common.Address)
	}
	// Determine the highest gas limit can be used during the estimation.
	if args.Gas != nil && uint64(*args.Gas) >= params.TxGas {
		hi = uint64(*args.Gas)
	} else {
		// Retrieve the block to act as the gas ceiling
		block, err := b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return 0, err
		}
		if block == nil {
			return 0, errors.New("block not found")
		}
		hi = block.GasLimit()
	}
	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	} else {
		feeCap = common.Big0
	}
	// Recap the highest gas limit with account's available balance.
	if feeCap.BitLen() != 0 {
		state, _, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return 0, err
		}
		balance := state.GetBalance(*args.From) // from can't be nil
		available := new(big.Int).Set(balance)
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, args.Value.ToInt())
		}
		allowance := new(big.Int).Div(available, feeCap)

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := args.Value
			if transfer == nil {
				transfer = new(hexutil.Big)
			}
			log.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer.ToInt(), "maxFeePerGas", feeCap, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	// Recap the highest gas allowance with specified gascap.
	if gasCap != 0 && hi > gasCap {
		log.Debug("Caller gas above allowance, capping", "requested", hi, "cap", gasCap)
		hi = gasCap
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *core.ExecutionResult, error) {
		args.Gas = (*hexutil.Uint64)(&gas)

		result, err := DoCall(ctx, b, args, blockNrOrHash, nil, 0, gasCap)
		if err != nil {
			if errors.Is(err, core.ErrIntrinsicGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return result.Failed(), result, nil
	}
	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)

		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		failed, result, err := executable(hi)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.Err != vm.ErrOutOfGas {
				if len(result.Revert()) > 0 {
					return 0, newRevertError(result)
				}
				return 0, result.Err
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hexutil.Uint64(hi), nil
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (s *PublicBlockChainAPI) EstimateGas(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	return DoEstimateGas(ctx, s.b, args, bNrOrHash, s.b.RPCGasCap())
}

// ExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type ExecutionResult struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc      uint64             `json:"pc"`
	Op      string             `json:"op"`
	Gas     uint64             `json:"gas"`
	GasCost uint64             `json:"gasCost"`
	Depth   int                `json:"depth"`
	Error   string             `json:"error,omitempty"`
	Stack   *[]string          `json:"stack,omitempty"`
	Memory  *[]string          `json:"memory,omitempty"`
	Storage *map[string]string `json:"storage,omitempty"`
}

// FormatLogs formats EVM returned structured logs for json output
func FormatLogs(logs []vm.StructLog) []StructLogRes {
	formatted := make([]StructLogRes, len(logs))
	for index, trace := range logs {
		formatted[index] = StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   trace.Depth,
			Error:   trace.ErrorString(),
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = stackValue.Hex()
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			memory := make([]string, 0, (len(trace.Memory)+31)/32)
			for i := 0; i+32 <= len(trace.Memory); i += 32 {
				memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
			}
			formatted[index].Memory = &memory
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted
}

// RPCMarshalHeader converts the given header to the RPC output .
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	result := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(head.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}

	if head.BaseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
	}

	return result
}

// RPCMarshalBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func RPCMarshalBlock(block *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields := RPCMarshalHeader(block.Header())
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}
		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return newRPCTransactionFromBlockHash(block, tx.Hash()), nil
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range txs {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}
	uncles := block.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes

	return fields, nil
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicBlockchainAPI`.
func (s *PublicBlockChainAPI) rpcMarshalHeader(ctx context.Context, header *types.Header) map[string]interface{} {
	fields := RPCMarshalHeader(header)
	fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, header.Hash()))
	return fields
}

// rpcMarshalBlock uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicBlockchainAPI`.
func (s *PublicBlockChainAPI) rpcMarshalBlock(ctx context.Context, b *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields, err := RPCMarshalBlock(b, inclTx, fullTx)
	if err != nil {
		return nil, err
	}
	if inclTx {
		fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, b.Hash()))
	}
	return fields, err
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash      `json:"blockHash"`
	BlockNumber      *hexutil.Big      `json:"blockNumber"`
	From             common.Address    `json:"from"`
	Gas              hexutil.Uint64    `json:"gas"`
	GasPrice         *hexutil.Big      `json:"gasPrice"`
	GasFeeCap        *hexutil.Big      `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash       `json:"hash"`
	Input            hexutil.Bytes     `json:"input"`
	Nonce            hexutil.Uint64    `json:"nonce"`
	To               *common.Address   `json:"to"`
	TransactionIndex *hexutil.Uint64   `json:"transactionIndex"`
	Value            *hexutil.Big      `json:"value"`
	Type             hexutil.Uint64    `json:"type"`
	Accesses         *types.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big      `json:"chainId,omitempty"`
	V                *hexutil.Big      `json:"v"`
	R                *hexutil.Big      `json:"r"`
	S                *hexutil.Big      `json:"s"`
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, baseFee *big.Int) *RPCTransaction {
	// Determine the signer. For replay-protected transactions, use the most permissive
	// signer, because we assume that signers are backwards-compatible with old
	// transactions. For non-protected transactions, the homestead signer signer is used
	// because the return value of ChainId is zero for those transactions.
	var signer types.Signer
	if tx.Protected() {
		signer = types.LatestSignerForChainID(tx.ChainId())
	} else {
		signer = types.HomesteadSigner{}
	}
	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case types.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			// price = min(tip, gasFeeCap - baseFee) + baseFee
			price := math.BigMin(new(big.Int).Add(tx.GasTipCap(), baseFee), tx.GasFeeCap())
			result.GasPrice = (*hexutil.Big)(price)
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
	}
	return result
}

// newRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func newRPCPendingTransaction(tx *types.Transaction, current *types.Header, config *params.ChainConfig) *RPCTransaction {
	var baseFee *big.Int
	if current != nil {
		baseFee = misc.CalcBaseFee(config, current)
	}
	return newRPCTransaction(tx, common.Hash{}, 0, 0, baseFee)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *types.Block, index uint64) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return newRPCTransaction(txs[index], b.Hash(), b.NumberU64(), index, b.BaseFee())
}

// newRPCRawTransactionFromBlockIndex returns the bytes of a transaction given a block and a transaction index.
func newRPCRawTransactionFromBlockIndex(b *types.Block, index uint64) hexutil.Bytes {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := txs[index].MarshalBinary()
	return blob
}

// newRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockHash(b *types.Block, hash common.Hash) *RPCTransaction {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newRPCTransactionFromBlockIndex(b, uint64(idx))
		}
	}
	return nil
}

// accessListResult returns an optional accesslist
// Its the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

// CreateAccessList creates a EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
func (s *PublicBlockChainAPI) CreateAccessList(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	acl, gasUsed, vmerr, err := AccessList(ctx, s.b, bNrOrHash, args)
	if err != nil {
		return nil, err
	}
	result := &accessListResult{Accesslist: &acl, GasUsed: hexutil.Uint64(gasUsed)}
	if vmerr != nil {
		result.Error = vmerr.Error()
	}
	return result, nil
}

// AccessList creates an access list for the given transaction.
// If the accesslist creation fails an error is returned.
// If the transaction itself fails, an vmErr is returned.
func AccessList(ctx context.Context, b Backend, blockNrOrHash rpc.BlockNumberOrHash, args TransactionArgs) (acl types.AccessList, gasUsed uint64, vmErr error, err error) {
	// Retrieve the execution context
	db, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if db == nil || err != nil {
		return nil, 0, nil, err
	}
	// If the gas amount is not set, extract this as it will depend on access
	// lists and we'll need to reestimate every time
	nogas := args.Gas == nil

	// Ensure any missing fields are filled, extract the recipient and input data
	if err := args.setDefaults(ctx, b); err != nil {
		return nil, 0, nil, err
	}
	var to common.Address
	if args.To != nil {
		to = *args.To
	} else {
		to = crypto.CreateAddress(args.from(), uint64(*args.Nonce))
	}
	// Retrieve the precompiles since they don't need to be added to the access list
	precompiles := vm.ActivePrecompiles(b.ChainConfig().Rules(header.Number))

	// Create an initial tracer
	prevTracer := vm.NewAccessListTracer(nil, args.from(), to, precompiles)
	if args.AccessList != nil {
		prevTracer = vm.NewAccessListTracer(*args.AccessList, args.from(), to, precompiles)
	}
	for {
		// Retrieve the current access list to expand
		accessList := prevTracer.AccessList()
		log.Trace("Creating access list", "input", accessList)

		// If no gas amount was specified, each unique access list needs it's own
		// gas calculation. This is quite expensive, but we need to be accurate
		// and it's convered by the sender only anyway.
		if nogas {
			args.Gas = nil
			if err := args.setDefaults(ctx, b); err != nil {
				return nil, 0, nil, err // shouldn't happen, just in case
			}
		}
		// Copy the original db so we don't modify it
		statedb := db.Copy()
		// Set the accesslist to the last al
		args.AccessList = &accessList
		msg, err := args.ToMessage(b.RPCGasCap(), header.BaseFee)
		if err != nil {
			return nil, 0, nil, err
		}

		// Apply the transaction with the access list tracer
		tracer := vm.NewAccessListTracer(accessList, args.from(), to, precompiles)
		config := vm.Config{Tracer: tracer, Debug: true, NoBaseFee: true}
		vmenv, _, err := b.GetEVM(ctx, msg, statedb, header, &config)
		if err != nil {
			return nil, 0, nil, err
		}
		res, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
		if err != nil {
			return nil, 0, nil, fmt.Errorf("failed to apply transaction: %v err: %v", args.toTransaction().Hash(), err)
		}
		if tracer.Equal(prevTracer) {
			return accessList, res.UsedGas, res.Err, nil
		}
		prevTracer = tracer
	}
}

// PublicTransactionPoolAPI exposes methods for the RPC interface
type PublicTransactionPoolAPI struct {
	b         Backend
	nonceLock *AddrLocker
	signer    types.Signer
	quit      chan int
	mtx       sync.Mutex
}

// NewPublicTransactionPoolAPI creates a new RPC service with methods specific for the transaction pool.
func NewPublicTransactionPoolAPI(b Backend, nonceLock *AddrLocker) *PublicTransactionPoolAPI {
	// The signer used by the API should always be the 'latest' known one because we expect
	// signers to be backwards-compatible with old transactions.
	signer := types.LatestSigner(b.ChainConfig())
	return &PublicTransactionPoolAPI{b: b, nonceLock: nonceLock, signer: signer, quit: nil}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (s *PublicTransactionPoolAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	// Ask transaction pool for the nonce which includes pending transactions
	if blockNr, ok := blockNrOrHash.Number(); ok && blockNr == rpc.PendingBlockNumber {
		nonce, err := s.b.GetPoolNonce(ctx, address)
		if err != nil {
			return nil, err
		}
		return (*hexutil.Uint64)(&nonce), nil
	}
	// Resolve block number and use its state to ask for the nonce
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	nonce := state.GetNonce(address)
	return (*hexutil.Uint64)(&nonce), state.Error()
}

// GetTransactionByHash returns the transaction for the given hash
func (s *PublicTransactionPoolAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	borTx := false

	// Try to return an already finalized transaction
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}

	// fetch bor block tx if necessary
	if tx == nil {
		if tx, blockHash, blockNumber, index, err = s.b.GetBorBlockTransaction(ctx, hash); err != nil {
			return nil, err
		}

		borTx = true
	}

	if tx != nil {
		resultTx := newRPCTransaction(tx, blockHash, blockNumber, index, nil)
		if borTx {
			// newRPCTransaction calculates hash based on RLP of the transaction data.
			// In case of bor block tx, we need simple derived tx hash (same as function argument) instead of RLP hash
			resultTx.Hash = hash
		}
		return resultTx, nil
	}

	// No finalized transaction, try to retrieve it from the pool
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return newRPCPendingTransaction(tx, s.b.CurrentHeader(), s.b.ChainConfig()), nil
	}

	// Transaction unknown, return as such
	return nil, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (s *PublicTransactionPoolAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	// Retrieve a finalized transaction, or a pooled otherwise
	tx, _, _, _, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			// Transaction not found anywhere, abort
			return nil, nil
		}
	}
	// Serialize to RLP and return
	return tx.MarshalBinary()
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	borTx := false

	tx, blockHash, blockNumber, index := rawdb.ReadTransaction(s.b.ChainDb(), hash)
	if tx == nil {
		tx, blockHash, blockNumber, index = rawdb.ReadBorTransaction(s.b.ChainDb(), hash)
		borTx = true
	}

	if tx == nil {
		return nil, nil
	}

	var receipt *types.Receipt

	if borTx {
		// Fetch bor block receipt
		receipt = rawdb.ReadBorReceipt(s.b.ChainDb(), blockHash, blockNumber)
	} else {
		receipts, err := s.b.GetReceipts(ctx, blockHash)
		if err != nil {
			return nil, err
		}
		if len(receipts) <= int(index) {
			return nil, nil
		}
		receipt = receipts[index]
	}

	// Derive the sender.
	bigblock := new(big.Int).SetUint64(blockNumber)
	signer := types.MakeSigner(s.b.ChainConfig(), bigblock)
	from, _ := types.Sender(signer, tx)

	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   hash,
		"transactionIndex":  hexutil.Uint64(index),
		"from":              from,
		"to":                tx.To(),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(tx.Type()),
	}
	// Assign the effective gas price paid
	if !s.b.ChainConfig().IsLondon(bigblock) {
		fields["effectiveGasPrice"] = hexutil.Uint64(tx.GasPrice().Uint64())
	} else {
		header, err := s.b.HeaderByHash(ctx, blockHash)
		if err != nil {
			return nil, err
		}
		gasPrice := new(big.Int).Add(header.BaseFee, tx.EffectiveGasTipValue(header.BaseFee))
		fields["effectiveGasPrice"] = hexutil.Uint64(gasPrice.Uint64())
	}
	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields, nil
}

// sign is a helper function that signs a transaction with the private key of the given address.
func (s *PublicTransactionPoolAPI) sign(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Request the wallet to sign the transaction
	return wallet.SignTx(account, tx, s.b.ChainConfig().ChainID)
}

// SubmitTransaction is a helper function that submits tx to txPool and logs a message.
func SubmitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
	// If the transaction fee cap is already specified, ensure the
	// fee of the given transaction is _reasonable_.
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}
	if !b.UnprotectedAllowed() && !tx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}
	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}
	// Print a log with full tx details for manual investigations and interventions
	signer := types.MakeSigner(b.ChainConfig(), b.CurrentBlock().Number())
	from, err := types.Sender(signer, tx)
	if err != nil {
		return common.Hash{}, err
	}

	if tx.To() == nil {
		addr := crypto.CreateAddress(from, tx.Nonce())
		log.Info("Submitted contract creation", "hash", tx.Hash().Hex(), "from", from, "nonce", tx.Nonce(), "contract", addr.Hex(), "value", tx.Value())
	} else {
		log.Info("Submitted transaction", "hash", tx.Hash().Hex(), "from", from, "nonce", tx.Nonce(), "recipient", tx.To(), "value", tx.Value())
	}
	return tx.Hash(), nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (s *PublicTransactionPoolAPI) SendTransaction(ctx context.Context, args TransactionArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: args.from()}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.from())
		defer s.nonceLock.UnlockAddr(args.from())
	}

	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	// Assemble the transaction and sign with the wallet
	tx := args.toTransaction()

	signed, err := wallet.SignTx(account, tx, s.b.ChainConfig().ChainID)
	if err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, signed)
}

/*func (s *PublicTransactionPoolAPI) getValidatorSectorInx(signer common.Address) (uint64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	method := "getValidatorSectorInx"
	data, _ := borcontracts.FABI.Pack(method, signer)
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(borcontracts.ValidatorFileCoinContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	curHeader, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber)
	result, err := DoCall(ctx, s.b, TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, rpc.BlockNumberOrHashWithHash(curHeader.Hash(), false), nil, 5*time.Second, s.b.RPCGasCap())
	if err != nil {
		return 0, errors.New(err.Error())
	}
	var ret = new(*big.Int)
	if err := borcontracts.FABI.UnpackIntoInterface(ret, method, result.Return()); err != nil {
		return 0, errors.New(err.Error())
	}
	u := (*ret).Uint64()
	return u, nil
}*/

type ClientDealArgs struct {
	File           string `json:"file"`
	PeerId         string `json:"peerId"`
	Addr1          string `json:"addr1"`
	Addr2          string `json:"addr2"`
	Addr3          string `json:"addr3"`
	Cid            string `json:"cid"`
	OutPath        string `json:"outPath"`
	MinerNum       int    `json:"minerNum"`
	Type           string `json:"type"`
	OriHash        string `json:"oriHash"`
	StoreKey       string `json:"storeKey"`
	FileHash       string `json:"fileHash"`
	HeadFlag       bool   `json:"headFlag"`
	MinerId        string `json:"minerId"`
	TxHash         string `json:"txHash"`
	Perm           string `json:"perm"`
	StorageType    string `json:"storageType"`
	WorkId         string `json:"workId"`
	SectorNum      string `json:"sectorNum"`
	WaitDealsDelay string `json:"waitDealsDelay"`
	ReallyDoIt     bool   `json:"reallyDoIt"`
	NewState       string `json:"newState"`
	Verbose        bool   `json:"verbose"`
	Completed      bool   `json:"completed"`
	Watch          bool   `json:"watch"`
	ShowFailed     bool   `json:"showFailed"`
	TransferID     string `json:"transferID"`
	Initiator      bool   `json:"initiator"`
	CancelTimeout  string `json:"cancelTimeout"`
	Concurrency    uint   `json:"concurrency"`
	Sealed         bool   `json:"sealed"`
}

type MinerStatus struct {
	MinerId        string  `json:"minerId"`
	Status         bool    `json:"status"`
	PublicKey      string  `json:"publicKey"`
	FreeSpaceRatio float64 `json:"freeSpaceRatio"`
}

func checkIfFileExist(path string, dir string) (bool, string) {
	if !filepath.IsAbs(path) {
		path = dir + "/" + path
	}
	if common.FileExist(path) {
		return true, path
	}
	return false, path
}

// Calculate the SHA256 of the file
func getSHA256ByFile(str string) (string, error) {
	file, err := os.Open(str)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sha256h := sha256.New()
	_, err = io.Copy(sha256h, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sha256h.Sum(nil)), nil
}

// PublicStorageAPI exposes methods for the RPC interface
type PublicStorageAPI struct {
	b         Backend
	nonceLock *AddrLocker
	signer    types.Signer
}

// NewPublicStorageAPI creates a new RPC service with methods specific for the storage.
func NewPublicStorageAPI(b Backend, nonceLock *AddrLocker) *PublicStorageAPI {
	// The signer used by the API should always be the 'latest' known one because we expect
	// signers to be backwards-compatible with old transactions.
	signer := types.LatestSigner(b.ChainConfig())
	return &PublicStorageAPI{b, nonceLock, signer}
}

func (s *PublicStorageAPI) ValidFileInfoAll(args ClientDealArgs, path string) error {
	storageType := args.StorageType
	var hashStr string
	if storageType == ENTIRE_FILE {
		hashStr = common.BuildNewOriHashWithSha256(args.StoreKey)
	} else {
		hashStr = args.OriHash
	}
	hashByte := common.HexSTrToByte32(hashStr)
	minerIdByte := common.HexSTrToByte32(args.MinerId)
	// Calculates the SHA256 of the merged file
	sha256Sum, err := getSHA256ByFile(path)
	if err != nil {
		return err
	}
	var retStatus uint8
	if storageType == ENTIRE_FILE {
		retStatus = file_store.FileStoreCli.ValidFileInfo4Entire(hashByte, args.StoreKey, sha256Sum)
		if retStatus == file_store.VALID_NO_PASS {
			lotusLog.Errorw("ValidFileInfo failed. Params", "storeKey", hashStr, "headFlag", args.HeadFlag, "minerId", args.MinerId, "fileHash", sha256Sum)
			return errors.New("file hash is not equal to the upload target file")
		} else if retStatus == file_store.NOT_FOUND {
			lotusLog.Errorw("ValidFileInfo failed. Params", "storeKey", hashStr, "headFlag", args.HeadFlag, "minerId", args.MinerId, "fileHash", sha256Sum)
			return errors.New("The current node does not find the contract data, the possible reasons are: 1. The contract data has not been registered; 2. The contract data is being synchronized")
		}
	} else {
		fileHash := common.HexSTrToByte32(sha256Sum)
		retStatus = file_store.FileStoreCli.ValidFileInfo(nil, hashByte, args.HeadFlag, minerIdByte, fileHash)
		if retStatus == file_store.VALID_NO_PASS {
			lotusLog.Errorw("ValidFileInfo failed. Params", "oriHashStr", hashStr, "headFlag", args.HeadFlag, "minerId", args.MinerId, "fileHash", sha256Sum)
			return errors.New("file hash is not equal to the upload target file")
		} else if retStatus == file_store.NOT_FOUND {
			lotusLog.Errorw("ValidFileInfo failed. Params", "oriHashStr", hashStr, "headFlag", args.HeadFlag, "minerId", args.MinerId, "fileHash", sha256Sum)
			return errors.New("The current node does not find the contract data, the possible reasons are: 1. The contract data has not been registered; 2. The contract data is being synchronized")
		}
	}
	return nil
}

// for entire file storage,not seperate to head and body
func (s *PublicStorageAPI) ClientSend4Entire(ctx context.Context, args ClientDealArgs) DealStatus {
	args.HeadFlag = false
	args.StorageType = ENTIRE_FILE
	return s.ClientSend(ctx, args)
}

func (s *PublicStorageAPI) ClientSend(ctx context.Context, args ClientDealArgs) DealStatus {
	lotuslog.SetupLogLevels()
	var dealStatus DealStatus
	storageType := args.StorageType
	var hashStr string
	if storageType == ENTIRE_FILE {
		dealStatus = s.GetStorageStatus4Entire(ctx, args)
		hashStr = args.StoreKey
	} else {
		dealStatus = s.GetStorageStatus(ctx, args)
		hashStr = args.OriHash
	}

	if dealStatus.Status == DEAL_SUCCESS {
		return newDealStatus(DEAL_SUCCESS, "success")
	}
	minerId := common.HexSTrToByte32(args.MinerId)
	addrInfo, peerId, err := s.getAddrInfoByMinerId(minerId)
	if err != nil {
		lotusLog.Errorw("cliendSend error", "oriHash", hashStr, "headFlag", args.HeadFlag, "minerId", args.MinerId, "peerId", peerId, "getAddrInfoByMinerId error", err)
		return newDealStatus(DEAL_ERROR, err.Error())
	}
	//addrInfo, peerId := s.getAddrInfo(args)
	mgr := s.b.ClientManager()
	//mgr.DealClient.Start(ctx)
	exist, path := checkIfFileExist(args.File, mgr.Repo.Path()+"/storage-file")
	if !exist {
		lotusLog.Errorw("cliendSend error", "oriHash/storeKey", hashStr, "headFlag", args.HeadFlag, "errorMsg", "file not exist, please upload file first or check the path")
		return newDealStatus(DEAL_ERROR, "file not exist, please upload file first or check the path")
	}

	err = s.ValidFileInfoAll(args, path)
	if err != nil {
		lotusLog.Errorw("cliendSend error", "oriHash/storeKey", hashStr, "headFlag", args.HeadFlag, "errorMsg", "validate the upload file error: %s", err)
		return newDealStatus(DEAL_ERROR, "validate the upload file error: "+err.Error())
	}

	res, err := clientImport(ctx, mgr, path)
	if err != nil {
		lotusLog.Errorw("cliendSend error", "oriHash/storeKey", hashStr, "headFlag", args.HeadFlag, "errorMsg", "client import %s", err)
		return newDealStatus(DEAL_ERROR, "client import error: "+err.Error())
	}
	imgr := mgr.importManager()

	if !mgr.HasStored(res.Root, peerId) && !mgr.HasStoring(res.Root, peerId) {
		mgr.MarkHasStoring(res.Root, peerId)
		mgr.MarkFileHashAndFlagToCid(hashStr, args.HeadFlag, peerId, res.Root)
		mgr.MarkImportId(res.Root, peerId, res.ImportID)
		mgr.MarkFile(res.Root, peerId, args.File)
		mgr.MarkFileHashAndFlag(res.Root, peerId, hashStr+","+strconv.FormatBool(args.HeadFlag)+","+args.StorageType)
		//mgr.MarkStorageType(res.Root, peerId, args.StorageType)
		err = mgr.Host.Connect(ctx, addrInfo)
		if err != nil {
			lotusLog.Errorw("cliendSend error", "oriHash/storeKey", hashStr, "headFlag", args.HeadFlag, "errorMsg", "client connect  %s", err)
			return newDealStatus(DEAL_ERROR, "client connect "+err.Error())
		}

		// in gorouting to avoid blocking the client
		go func() {
			_, errRet := mgr.StartDeal(context.Background(), res.Root, peerId)
			if errRet != nil {
				lotusLog.Errorw("cliendSend error", "oriHash/storeKey", hashStr, "headFlag", args.HeadFlag, "errorMsg", "mgr.StartDeal error: %s", errRet)
			}
		}()
	} else {
		if mgr.HasStored(res.Root, peerId) {
			if strings.Contains(path, "/storage-file") {
				os.Remove(path)
			}
		}
		carPath, err := imgr.CARPathFor(res.Root)
		if err == nil {
			os.Remove(carPath)
		}
		imgr.Remove(res.ImportID)
	}

	return newDealStatus(DEAL_SUCCESS, res.Root.String())
}

func getMinerStatus(minerId [32]byte, storeMiners [][32]byte) bool {
	if storeMiners != nil {
		for _, miner := range storeMiners {
			if miner == minerId {
				return true
			}
		}
	}
	return false
}

func (s *PublicTransactionPoolAPI) AuthNew(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	perm := args.Perm
	idx := 0
	for i, p := range lapi.AllPermissions {
		if auth.Permission(perm) == p {
			idx = i + 1
		}
	}

	if idx == 0 {
		err := fmt.Errorf("--perm flag has to be one of: %s", lapi.AllPermissions)
		lotusLog.Error(err)
		return "", err
	}

	// slice on [:idx] so for example: 'sign' gives you [read, write, sign]
	token, err := (*minerApi).AuthNew(ctx, lapi.AllPermissions[:idx])
	if err != nil {
		lotusLog.Error(err)
		return "", err
	}

	currentEnv, _, _ := cliutil.EnvsForAPIInfos(repo.StorageMiner)
	rdir := s.b.GetDataDir()
	ma, err := APIEndpoint(rdir + "/.w3fsminer")
	if err != nil {
		lotusLog.Error(err)
		return "", err
	}
	return fmt.Sprintf("%s=%s:%s", currentEnv, string(token), ma), nil
}

func (s *PublicTransactionPoolAPI) WorkerJobs(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	jobs, err := (*minerApi).WorkerJobs(ctx)
	if err != nil {
		return "", xerrors.Errorf("getting worker jobs: %w", err)
	}

	type line struct {
		storiface.WorkerJob
		wid uuid.UUID
	}

	lines := make([]line, 0)

	for wid, jobs := range jobs {
		for _, job := range jobs {
			lines = append(lines, line{
				WorkerJob: job,
				wid:       wid,
			})
		}
	}

	// oldest first
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].RunWait != lines[j].RunWait {
			return lines[i].RunWait < lines[j].RunWait
		}
		if lines[i].Start.Equal(lines[j].Start) {
			return lines[i].ID.ID.String() < lines[j].ID.ID.String()
		}
		return lines[i].Start.Before(lines[j].Start)
	})

	workerHostnames := map[uuid.UUID]string{}

	wst, err := (*minerApi).WorkerStats(ctx)
	if err != nil {
		return "", xerrors.Errorf("getting worker stats: %w", err)
	}

	for wid, st := range wst {
		workerHostnames[wid] = st.Info.Hostname
	}
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 2, 4, 2, ' ', 0)
	_, _ = fmt.Fprintf(tw, "ID\tSector\tWorker\tHostname\tTask\tState\tTime\n")

	for _, l := range lines {
		state := "running"
		switch {
		case l.RunWait > 1:
			state = fmt.Sprintf("assigned(%d)", l.RunWait-1)
		case l.RunWait == storiface.RWPrepared:
			state = "prepared"
		case l.RunWait == storiface.RWRetDone:
			state = "ret-done"
		case l.RunWait == storiface.RWReturned:
			state = "returned"
		case l.RunWait == storiface.RWRetWait:
			state = "ret-wait"
		}
		dur := "n/a"
		if !l.Start.IsZero() {
			dur = time.Now().Sub(l.Start).Truncate(time.Millisecond * 100).String()
		}

		hostname, ok := workerHostnames[l.wid]
		if !ok {
			hostname = l.Hostname
		}

		_, _ = fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n",
			hex.EncodeToString(l.ID.ID[:4]),
			l.Sector.Number,
			hex.EncodeToString(l.wid[:4]),
			hostname,
			l.Task.Short(),
			state,
			dur)
	}
	err = tw.Flush()

	return buf.String(), err
}

func (s *PublicTransactionPoolAPI) ListWorker(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	stats, err := (*minerApi).WorkerStats(ctx)
	if err != nil {
		return "", err
	}

	type sortableStat struct {
		id uuid.UUID
		storiface.WorkerStats
	}

	st := make([]sortableStat, 0, len(stats))
	for id, stat := range stats {
		st = append(st, sortableStat{id, stat})
	}

	sort.Slice(st, func(i, j int) bool {
		return st[i].id.String() < st[j].id.String()
	})

	var buf bytes.Buffer

	for _, stat := range st {
		gpuUse := "not "
		gpuCol := color.FgBlue
		if stat.GpuUsed {
			gpuCol = color.FgGreen
			gpuUse = ""
		}

		var disabled string
		if !stat.Enabled {
			disabled = color.RedString(" (disabled)")
		}

		tmpStr := fmt.Sprintf("Worker %s, host %s%s\n", stat.id, color.MagentaString(stat.Info.Hostname), disabled)
		buf.WriteString(tmpStr)

		var barCols = uint64(64)
		cpuBars := int(stat.CpuUse * barCols / stat.Info.Resources.CPUs)
		cpuBar := strings.Repeat("|", cpuBars)
		if int(barCols)-cpuBars >= 0 {
			cpuBar += strings.Repeat(" ", int(barCols)-cpuBars)
		}

		tmpStr = fmt.Sprintf("\tCPU:  [%s] %d/%d core(s) in use\n",
			color.GreenString(cpuBar), stat.CpuUse, stat.Info.Resources.CPUs)
		buf.WriteString(tmpStr)

		ramBarsRes := int(stat.Info.Resources.MemReserved * barCols / stat.Info.Resources.MemPhysical)
		ramBarsUsed := int(stat.MemUsedMin * barCols / stat.Info.Resources.MemPhysical)
		ramRepeatSpace := int(barCols) - (ramBarsUsed + ramBarsRes)

		colorFunc := color.YellowString
		if ramRepeatSpace < 0 {
			ramRepeatSpace = 0
			colorFunc = color.RedString
		}

		ramBar := colorFunc(strings.Repeat("|", ramBarsRes)) +
			color.GreenString(strings.Repeat("|", ramBarsUsed)) +
			strings.Repeat(" ", ramRepeatSpace)

		vmem := stat.Info.Resources.MemPhysical + stat.Info.Resources.MemSwap

		vmemBarsRes := int(stat.Info.Resources.MemReserved * barCols / vmem)
		vmemBarsUsed := int(stat.MemUsedMax * barCols / vmem)
		vmemRepeatSpace := int(barCols) - (vmemBarsUsed + vmemBarsRes)

		colorFunc = color.YellowString
		if vmemRepeatSpace < 0 {
			vmemRepeatSpace = 0
			colorFunc = color.RedString
		}

		vmemBar := colorFunc(strings.Repeat("|", vmemBarsRes)) +
			color.GreenString(strings.Repeat("|", vmemBarsUsed)) +
			strings.Repeat(" ", vmemRepeatSpace)

		tmpStr = fmt.Sprintf("\tRAM:  [%s] %d%% %s/%s\n", ramBar,
			(stat.Info.Resources.MemReserved+stat.MemUsedMin)*100/stat.Info.Resources.MemPhysical,
			chaintypes.SizeStr(chaintypes.NewInt(stat.Info.Resources.MemReserved+stat.MemUsedMin)),
			chaintypes.SizeStr(chaintypes.NewInt(stat.Info.Resources.MemPhysical)))
		buf.WriteString(tmpStr)

		tmpStr = fmt.Sprintf("\tVMEM: [%s] %d%% %s/%s\n", vmemBar,
			(stat.Info.Resources.MemReserved+stat.MemUsedMax)*100/vmem,
			chaintypes.SizeStr(chaintypes.NewInt(stat.Info.Resources.MemReserved+stat.MemUsedMax)),
			chaintypes.SizeStr(chaintypes.NewInt(vmem)))
		buf.WriteString(tmpStr)

		for _, gpu := range stat.Info.Resources.GPUs {
			tmpStr = fmt.Sprintf("\tGPU: %s\n", color.New(gpuCol).Sprintf("%s, %sused", gpu, gpuUse))
			buf.WriteString(tmpStr)
		}
	}

	return buf.String(), nil
}

func (s *PublicTransactionPoolAPI) SchedDiag(ctx context.Context) (interface{}, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	st, err := (*minerApi).SealingSchedDiag(ctx, false)
	if err != nil {
		return "", err
	}

	return st, nil

	//j, err := json.MarshalIndent(&st, "", "  ")
	//if err != nil {
	//	return "", err
	//}

	//return string(j), nil
}

func (s *PublicTransactionPoolAPI) AbortRunningJob(ctx context.Context, args ClientDealArgs) error {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	jobs, err := (*minerApi).WorkerJobs(ctx)
	if err != nil {
		return xerrors.Errorf("getting worker jobs: %w", err)
	}

	var job *storiface.WorkerJob
outer:
	for _, workerJobs := range jobs {
		for _, j := range workerJobs {
			if strings.HasPrefix(j.ID.ID.String(), args.WorkId) {
				j := j
				job = &j
				break outer
			}
		}
	}

	if job == nil {
		return xerrors.Errorf("job with specified id prefix not found")
	}

	lotusLog.Info("aborting job %s, task %s, sector %d, running on host %s\n", job.ID.String(), job.Task.Short(), job.Sector.Number, job.Hostname)

	return (*minerApi).SealingAbort(ctx, job.ID)
}

func (s *PublicTransactionPoolAPI) StorageList(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	st, err := (*minerApi).StorageList(ctx)
	if err != nil {
		return "", err
	}

	local, err := (*minerApi).StorageLocal(ctx)
	if err != nil {
		return "", err
	}

	type fsInfo struct {
		stores.ID
		sectors []stores.Decl
		stat    fsutil.FsStat
	}

	sorted := make([]fsInfo, 0, len(st))
	for id, decls := range st {
		st, err := (*minerApi).StorageStat(ctx, id)
		if err != nil {
			sorted = append(sorted, fsInfo{ID: id, sectors: decls})
			continue
		}

		sorted = append(sorted, fsInfo{id, decls, st})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].stat.Capacity != sorted[j].stat.Capacity {
			return sorted[i].stat.Capacity > sorted[j].stat.Capacity
		}
		return sorted[i].ID < sorted[j].ID
	})

	var buf bytes.Buffer

	for _, s := range sorted {

		var cnt [3]int
		for _, decl := range s.sectors {
			for i := range cnt {
				if decl.SectorFileType&(1<<i) != 0 {
					cnt[i]++
				}
			}
		}

		buf.WriteString(fmt.Sprintf("%s:\n", s.ID))

		pingStart := time.Now()
		st, err := (*minerApi).StorageStat(ctx, s.ID)
		if err != nil {
			buf.WriteString(fmt.Sprintf("\t%s: %s:\n", color.RedString("Error"), err))
			continue
		}
		ping := time.Now().Sub(pingStart)

		safeRepeat := func(s string, count int) string {
			if count < 0 {
				return ""
			}
			return strings.Repeat(s, count)
		}

		var barCols = int64(50)

		// filesystem use bar
		{
			usedPercent := (st.Capacity - st.FSAvailable) * 100 / st.Capacity

			percCol := color.FgGreen
			switch {
			case usedPercent > 98:
				percCol = color.FgRed
			case usedPercent > 90:
				percCol = color.FgYellow
			}

			set := (st.Capacity - st.FSAvailable) * barCols / st.Capacity
			used := (st.Capacity - (st.FSAvailable + st.Reserved)) * barCols / st.Capacity
			reserved := set - used
			bar := safeRepeat("#", int(used)) + safeRepeat("*", int(reserved)) + safeRepeat(" ", int(barCols-set))

			desc := ""
			if st.Max > 0 {
				desc = " (filesystem)"
			}

			buf.WriteString(fmt.Sprintf("\t[%s] %s/%s %s%s\n", color.New(percCol).Sprint(bar),
				chaintypes.SizeStr(chaintypes.NewInt(uint64(st.Capacity-st.FSAvailable))),
				chaintypes.SizeStr(chaintypes.NewInt(uint64(st.Capacity))),
				color.New(percCol).Sprintf("%d%%", usedPercent), desc))
		}

		// optional configured limit bar
		if st.Max > 0 {
			usedPercent := st.Used * 100 / st.Max

			percCol := color.FgGreen
			switch {
			case usedPercent > 98:
				percCol = color.FgRed
			case usedPercent > 90:
				percCol = color.FgYellow
			}

			set := st.Used * barCols / st.Max
			used := (st.Used + st.Reserved) * barCols / st.Max
			reserved := set - used
			bar := safeRepeat("#", int(used)) + safeRepeat("*", int(reserved)) + safeRepeat(" ", int(barCols-set))

			buf.WriteString(fmt.Sprintf("\t[%s] %s/%s %s (limit)\n", color.New(percCol).Sprint(bar),
				chaintypes.SizeStr(chaintypes.NewInt(uint64(st.Used))),
				chaintypes.SizeStr(chaintypes.NewInt(uint64(st.Max))),
				color.New(percCol).Sprintf("%d%%", usedPercent)))
		}

		buf.WriteString(fmt.Sprintf("\t%s; %s; %s; Reserved: %s\n",
			color.YellowString("Unsealed: %d", cnt[0]),
			color.GreenString("Sealed: %d", cnt[1]),
			color.BlueString("Caches: %d", cnt[2]),
			chaintypes.SizeStr(chaintypes.NewInt(uint64(st.Reserved)))))

		si, err := (*minerApi).StorageInfo(ctx, s.ID)
		if err != nil {
			return "", err
		}

		buf.WriteString(fmt.Sprintf("\t"))
		if si.CanSeal || si.CanStore {
			buf.WriteString(fmt.Sprintf("Weight: %d; Use: ", si.Weight))
			if si.CanSeal {
				buf.WriteString(fmt.Sprintf(color.MagentaString("Seal ")))
			}
			if si.CanStore {
				buf.WriteString(fmt.Sprintf(color.CyanString("Store")))
			}
			buf.WriteString(fmt.Sprintf("\n"))
		} else {
			buf.WriteString(fmt.Sprintf(color.HiYellowString("Use: ReadOnly")))
		}

		if localPath, ok := local[s.ID]; ok {
			buf.WriteString(fmt.Sprintf("\tLocal: %s\n", color.GreenString(localPath)))
		}
		for i, l := range si.URLs {
			var rtt string
			if _, ok := local[s.ID]; !ok && i == 0 {
				rtt = " (latency: " + ping.Truncate(time.Microsecond*100).String() + ")"
			}

			buf.WriteString(fmt.Sprintf("\tURL: %s%s\n", l, rtt)) // TODO; try pinging maybe?? print latency?
		}
		buf.WriteString(fmt.Sprintf("\n"))
	}

	return buf.String(), nil
}

type storedSector struct {
	id    stores.ID
	store stores.SectorStorageInfo

	unsealed, sealed, cache bool
}

func (s *PublicTransactionPoolAPI) StorageFind(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	ma, err := (*minerApi).ActorAddress(ctx)
	if err != nil {
		return "", err
	}

	mid, err := address.IDFromAddress(ma)
	if err != nil {
		return "", err
	}

	snum, err := strconv.ParseUint(args.SectorNum, 10, 64)
	if err != nil {
		return "", err
	}

	sid := fabi.SectorID{
		Miner:  fabi.ActorID(mid),
		Number: fabi.SectorNumber(snum),
	}

	u, err := (*minerApi).StorageFindSector(ctx, sid, storiface.FTUnsealed, 0, false)
	if err != nil {
		return "", xerrors.Errorf("finding unsealed: %w", err)
	}

	ssi, err := (*minerApi).StorageFindSector(ctx, sid, storiface.FTSealed, 0, false)
	if err != nil {
		return "", xerrors.Errorf("finding sealed: %w", err)
	}

	c, err := (*minerApi).StorageFindSector(ctx, sid, storiface.FTCache, 0, false)
	if err != nil {
		return "", xerrors.Errorf("finding cache: %w", err)
	}

	byId := map[stores.ID]*storedSector{}
	for _, info := range u {
		sts, ok := byId[info.ID]
		if !ok {
			sts = &storedSector{
				id:    info.ID,
				store: info,
			}
			byId[info.ID] = sts
		}
		sts.unsealed = true
	}
	for _, info := range ssi {
		sts, ok := byId[info.ID]
		if !ok {
			sts = &storedSector{
				id:    info.ID,
				store: info,
			}
			byId[info.ID] = sts
		}
		sts.sealed = true
	}
	for _, info := range c {
		sts, ok := byId[info.ID]
		if !ok {
			sts = &storedSector{
				id:    info.ID,
				store: info,
			}
			byId[info.ID] = sts
		}
		sts.cache = true
	}

	local, err := (*minerApi).StorageLocal(ctx)
	if err != nil {
		return "", err
	}

	var out []*storedSector
	for _, sector := range byId {
		out = append(out, sector)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].id < out[j].id
	})

	var buf bytes.Buffer

	for _, info := range out {
		var types string
		if info.unsealed {
			types += "Unsealed, "
		}
		if info.sealed {
			types += "Sealed, "
		}
		if info.cache {
			types += "Cache, "
		}

		buf.WriteString(fmt.Sprintf("In %s (%s)\n", info.id, types[:len(types)-2]))
		buf.WriteString(fmt.Sprintf("\tSealing: %t; Storage: %t\n", info.store.CanSeal, info.store.CanStore))
		if localPath, ok := local[info.id]; ok {
			buf.WriteString(fmt.Sprintf("\tLocal (%s)\n", localPath))
		} else {
			buf.WriteString(fmt.Sprintf("\tRemote\n"))
		}
		for _, l := range info.store.URLs {
			buf.WriteString(fmt.Sprintf("\tURL: %s\n", l))
		}
	}
	str := buf.String()
	if len(str) == 0 {
		str = "\n"
	}
	return str, nil
}

func (s *PublicTransactionPoolAPI) StorageCleanup(ctx context.Context) error {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	sectors, err := (*minerApi).SectorsList(ctx)
	if err != nil {
		return err
	}

	maddr, err := (*minerApi).ActorAddress(ctx)
	if err != nil {
		return err
	}

	ssize, err := (*minerApi).ActorSectorSize(ctx, maddr)
	if err != nil {
		return err
	}

	aid, err := address.IDFromAddress(maddr)
	if err != nil {
		return err
	}

	sid := func(sn fabi.SectorNumber) fabi.SectorID {
		return fabi.SectorID{
			Miner:  fabi.ActorID(aid),
			Number: sn,
		}
	}

	toRemove := map[fabi.SectorNumber]struct{}{}

	for _, sector := range sectors {
		st, err := (*minerApi).SectorsStatus(ctx, sector, false)
		if err != nil {
			return xerrors.Errorf("getting sector status for sector %d: %w", sector, err)
		}

		if lsealing.SectorState(st.State) != lsealing.Removed {
			continue
		}

		for _, ft := range storiface.PathTypes {
			si, err := (*minerApi).StorageFindSector(ctx, sid(sector), ft, ssize, false)
			if err != nil {
				return xerrors.Errorf("find sector %d: %w", sector, err)
			}

			if len(si) > 0 {
				toRemove[sector] = struct{}{}
			}
		}
	}

	for sn := range toRemove {
		lotusLog.Infof("cleaning up data for sector %d\n", sn)
		err := (*minerApi).SectorRemove(ctx, sn)
		if err != nil {
			lotusLog.Errorf("cleaning up data for sector %d, err %s", sn, err)
		}
	}

	return nil
}

func (s *PublicTransactionPoolAPI) StorageListSectors(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	sectors, err := (*minerApi).SectorsList(ctx)
	if err != nil {
		return "", xerrors.Errorf("listing sectors: %w", err)
	}

	maddr, err := (*minerApi).ActorAddress(ctx)
	if err != nil {
		return "", err
	}

	aid, err := address.IDFromAddress(maddr)
	if err != nil {
		return "", err
	}

	ssize, err := (*minerApi).ActorSectorSize(ctx, maddr)
	if err != nil {
		return "", err
	}

	sid := func(sn fabi.SectorNumber) fabi.SectorID {
		return fabi.SectorID{
			Miner:  fabi.ActorID(aid),
			Number: sn,
		}
	}

	type entry struct {
		id      fabi.SectorNumber
		storage stores.ID
		ft      storiface.SectorFileType
		urls    string

		primary, seal, store bool

		state lapi.SectorState
	}

	var list []entry

	for _, sector := range sectors {
		st, err := (*minerApi).SectorsStatus(ctx, sector, false)
		if err != nil {
			return "", xerrors.Errorf("getting sector status for sector %d: %w", sector, err)
		}

		for _, ft := range storiface.PathTypes {
			si, err := (*minerApi).StorageFindSector(ctx, sid(sector), ft, ssize, false)
			if err != nil {
				return "", xerrors.Errorf("find sector %d: %w", sector, err)
			}

			for _, info := range si {

				list = append(list, entry{
					id:      sector,
					storage: info.ID,
					ft:      ft,
					urls:    strings.Join(info.URLs, ";"),

					primary: info.Primary,
					seal:    info.CanSeal,
					store:   info.CanStore,

					state: st.State,
				})
			}
		}

	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].store != list[j].store {
			return list[i].store
		}

		if list[i].storage != list[j].storage {
			return list[i].storage < list[j].storage
		}

		if list[i].id != list[j].id {
			return list[i].id < list[j].id
		}

		return list[i].ft < list[j].ft
	})

	tw := tablewriter.New(
		tablewriter.Col("Storage"),
		tablewriter.Col("Sector"),
		tablewriter.Col("Type"),
		tablewriter.Col("State"),
		tablewriter.Col("Primary"),
		tablewriter.Col("Path use"),
		tablewriter.Col("URLs"),
	)

	if len(list) == 0 {
		return "\n", nil
	}

	lastS := list[0].storage
	sc1, sc2 := color.FgBlue, color.FgCyan

	for _, e := range list {
		if e.storage != lastS {
			lastS = e.storage
			sc1, sc2 = sc2, sc1
		}

		m := map[string]interface{}{
			"Storage":  color.New(sc1).Sprint(e.storage),
			"Sector":   e.id,
			"Type":     e.ft.String(),
			"State":    color.New(stateOrder[lsealing.SectorState(e.state)].col).Sprint(e.state),
			"Primary":  maybeStr(e.seal, color.FgGreen, "primary"),
			"Path use": maybeStr(e.seal, color.FgMagenta, "seal ") + maybeStr(e.store, color.FgCyan, "store"),
			"URLs":     e.urls,
		}
		tw.Write(m)
	}
	var buf bytes.Buffer
	err = tw.Flush(&buf)
	str := buf.String()
	if len(str) == 0 {
		str = "\n"
	}

	return str, err
}

func (s *PublicTransactionPoolAPI) SectorsStatus(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	id, err := strconv.ParseUint(args.SectorNum, 10, 64)
	if err != nil {
		return "", err
	}

	status, err := (*minerApi).SectorsStatus(ctx, fabi.SectorNumber(id), true)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	tmpStr := fmt.Sprintf("SectorID:\t%d\n", status.SectorID)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Status:\t\t%s\n", status.State)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("CIDcommD:\t%s\n", status.CommD)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("CIDcommR:\t%s\n", status.CommR)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Ticket:\t\t%x\n", status.Ticket.Value)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("TicketH:\t%d\n", status.Ticket.Epoch)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Seed:\t\t%x\n", status.Seed.Value)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("SeedH:\t\t%d\n", status.Seed.Epoch)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Precommit:\t%s\n", status.PreCommitMsg)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Commit:\t\t%s\n", status.CommitMsg)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Proof:\t\t%x\n", status.Proof)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Deals:\t\t%v\n", status.Deals)
	buf.WriteString(tmpStr)
	tmpStr = fmt.Sprintf("Retries:\t%d\n", status.Retries)
	buf.WriteString(tmpStr)
	if status.LastErr != "" {
		tmpStr = fmt.Sprintf("Last Error:\t\t%s\n", status.LastErr)
		buf.WriteString(tmpStr)
	}

	if true {
		tmpStr = fmt.Sprintf("\nSector On Chain Info\n")
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("SealProof:\t\t%x\n", status.SealProof)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("Activation:\t\t%v\n", status.Activation)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("Expiration:\t\t%v\n", status.Expiration)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("DealWeight:\t\t%v\n", status.DealWeight)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("VerifiedDealWeight:\t\t%v\n", status.VerifiedDealWeight)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("InitialPledge:\t\t%v\n", status.InitialPledge)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("\nExpiration Info\n")
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("OnTime:\t\t%v\n", status.OnTime)
		buf.WriteString(tmpStr)
		tmpStr = fmt.Sprintf("Early:\t\t%v\n", status.Early)
		buf.WriteString(tmpStr)
	}

	//if true {
	//	fullApi, nCloser, err := lcli.GetFullNodeAPI(cctx)
	//	if err != nil {
	//		return "", err
	//	}
	//	defer nCloser()
	//
	//	maddr, err := getActorAddress(ctx, cctx)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	mact, err := fullApi.StateGetActor(ctx, maddr, chaintypes.EmptyTSK)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	tbs := blockstore.NewTieredBstore(blockstore.NewAPIBlockstore(fullApi), blockstore.NewMemory())
	//	mas, err := miner.Load(adt.WrapStore(ctx, cbor.NewCborStore(tbs)), mact)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	errFound := errors.New("found")
	//	if err := mas.ForEachDeadline(func(dlIdx uint64, dl miner.Deadline) error {
	//		return dl.ForEachPartition(func(partIdx uint64, part miner.Partition) error {
	//			pas, err := part.AllSectors()
	//			if err != nil {
	//				return err
	//			}
	//
	//			set, err := pas.IsSet(id)
	//			if err != nil {
	//				return err
	//			}
	//			if set {
	//				tmpStr = fmt.Sprintf("\nDeadline:\t%d\n", dlIdx)
	//				buf.WriteString(tmpStr)
	//				tmpStr = fmt.Sprintf("Partition:\t%d\n", partIdx)
	//				buf.WriteString(tmpStr)
	//
	//				checkIn := func(name string, bg func() (bitfield.BitField, error)) error {
	//					bf, err := bg()
	//					if err != nil {
	//						return err
	//					}
	//
	//					set, err := bf.IsSet(id)
	//					if err != nil {
	//						return err
	//					}
	//					setstr := "no"
	//					if set {
	//						setstr = "yes"
	//					}
	//					tmpStr = fmt.Sprintf("%s:   \t%s\n", name, setstr)
	//					buf.WriteString(tmpStr)
	//					return nil
	//				}
	//
	//				if err := checkIn("Unproven", part.UnprovenSectors); err != nil {
	//					return err
	//				}
	//				if err := checkIn("Live", part.LiveSectors); err != nil {
	//					return err
	//				}
	//				if err := checkIn("Active", part.ActiveSectors); err != nil {
	//					return err
	//				}
	//				if err := checkIn("Faulty", part.FaultySectors); err != nil {
	//					return err
	//				}
	//				if err := checkIn("Recovering", part.RecoveringSectors); err != nil {
	//					return err
	//				}
	//
	//				return errFound
	//			}
	//
	//			return nil
	//		})
	//	}); err != errFound {
	//		if err != nil {
	//			return "", err
	//		}
	//
	//		tmpStr = fmt.Sprintf("\nNot found in any partition")
	//		buf.WriteString(tmpStr)
	//	}
	//}

	if true {
		tmpStr = fmt.Sprintf("--------\nEvent Log:\n")
		buf.WriteString(tmpStr)

		for i, l := range status.Log {
			tmpStr = fmt.Sprintf("%d.\t%s:\t[%s]\t%s\n", i, time.Unix(int64(l.Timestamp), 0), l.Kind, l.Message)
			buf.WriteString(tmpStr)
			if l.Trace != "" {
				tmpStr = fmt.Sprintf("\t%s\n", l.Trace)
				buf.WriteString(tmpStr)
			}
		}
	}
	str := buf.String()
	if len(str) == 0 {
		str = "\n"
	}

	return str, nil
}

func (s *PublicTransactionPoolAPI) SectorsList(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	var list []fabi.SectorNumber

	showRemoved := true
	//var states []lapi.SectorState

	//if cctx.IsSet("states") {
	//	showRemoved = true
	//	sList := strings.Split(cctx.String("states"), ",")
	//	states = make([]lapi.SectorState, len(sList))
	//	for i := range sList {
	//		states[i] = lapi.SectorState(sList[i])
	//	}
	//}
	//
	//if cctx.Bool("unproven") {
	//	for state := range lsealing.ExistSectorStateList {
	//		if state == lsealing.Proving {
	//			continue
	//		}
	//		states = append(states, lapi.SectorState(state))
	//	}
	//}

	//if len(states) == 0 {
	list, err := (*minerApi).SectorsList(ctx)
	//} else {
	//	list, err = (*minerApi).SectorsListInStates(ctx, states)
	//}

	if err != nil {
		return "", err
	}

	//maddr, err := (*minerApi).ActorAddress(ctx)
	//if err != nil {
	//	return "", err
	//}

	//activeSet, err := fullApi.StateMinerActiveSectors(ctx, maddr, head.Key())
	//if err != nil {
	//	return err
	//}
	//activeIDs := make(map[fabi.SectorNumber]struct{}, len(activeSet))
	//for _, info := range activeSet {
	//	activeIDs[info.SectorNumber] = struct{}{}
	//}
	//
	//sset, err := fullApi.StateMinerSectors(ctx, maddr, nil, head.Key())
	//if err != nil {
	//	return err
	//}
	//commitedIDs := make(map[fabi.SectorNumber]struct{}, len(sset))
	//for _, info := range sset {
	//	commitedIDs[info.SectorNumber] = struct{}{}
	//}

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	tw := tablewriter.New(
		tablewriter.Col("ID"),
		tablewriter.Col("State"),
		tablewriter.Col("OnChain"),
		tablewriter.Col("Active"),
		tablewriter.Col("Expiration"),
		tablewriter.Col("SealTime"),
		tablewriter.Col("Events"),
		tablewriter.Col("Deals"),
		tablewriter.Col("DealWeight"),
		tablewriter.Col("VerifiedPower"),
		tablewriter.NewLineCol("Error"),
		tablewriter.NewLineCol("RecoveryTimeout"))

	for _, s := range list {
		st, err := (*minerApi).SectorsStatus(ctx, s, true)
		if err != nil {
			tw.Write(map[string]interface{}{
				"ID":    s,
				"Error": err,
			})
			continue
		}

		if !showRemoved && st.State == lapi.SectorState(lsealing.Removed) {
			continue
		}

		//_, inSSet := commitedIDs[s]
		//_, inASet := activeIDs[s]

		const verifiedPowerGainMul = 9

		dw, vp := .0, .0
		estimate := st.Expiration-st.Activation <= 0
		if !estimate {
			rdw := bigext.Add(st.DealWeight, st.VerifiedDealWeight)
			dw = float64(bigext.Div(rdw, bigext.NewInt(int64(st.Expiration-st.Activation))).Uint64())
			vp = float64(bigext.Div(bigext.Mul(st.VerifiedDealWeight, bigext.NewInt(verifiedPowerGainMul)), bigext.NewInt(int64(st.Expiration-st.Activation))).Uint64())
		} else {
			for _, piece := range st.Pieces {
				if piece.DealInfo != nil {
					dw += float64(piece.Piece.Size)
					if piece.DealInfo.DealProposal != nil && piece.DealInfo.DealProposal.VerifiedDeal {
						vp += float64(piece.Piece.Size) * verifiedPowerGainMul
					}
				}
			}
		}

		var deals int
		for _, deal := range st.Deals {
			if deal != 0 {
				deals++
			}
		}

		exp := st.Expiration
		if st.OnTime > 0 && st.OnTime < exp {
			exp = st.OnTime // Can be different when the sector was CC upgraded
		}

		m := map[string]interface{}{
			"ID":    s,
			"State": color.New(stateOrder[lsealing.SectorState(st.State)].col).Sprint(st.State),
			//"OnChain": yesno(inSSet),
			//"Active":  yesno(inASet),
		}

		if deals > 0 {
			m["Deals"] = color.GreenString("%d", deals)
		} else {
			m["Deals"] = color.BlueString("CC")
			if st.ToUpgrade {
				m["Deals"] = color.CyanString("CC(upgrade)")
			}
		}

		//if true {
		//	if !inSSet {
		//		m["Expiration"] = "n/a"
		//	} else {
		//		m["Expiration"] = lcli.EpochTime(head.Height(), exp)
		//		if st.Early > 0 {
		//			m["RecoveryTimeout"] = color.YellowString(lcli.EpochTime(head.Height(), st.Early))
		//		}
		//	}
		//}

		if deals > 0 {
			estWrap := func(s string) string {
				if !estimate {
					return s
				}
				return fmt.Sprintf("[%s]", s)
			}

			m["DealWeight"] = estWrap(units.BytesSize(dw))
			if vp > 0 {
				m["VerifiedPower"] = estWrap(color.GreenString(units.BytesSize(vp)))
			}
		}

		if true {
			var events int
			for _, sectorLog := range st.Log {
				if !strings.HasPrefix(sectorLog.Kind, "event") {
					continue
				}
				if sectorLog.Kind == "event;sealingext.SectorRestart" {
					continue
				}
				events++
			}

			pieces := len(st.Deals)

			switch {
			case events < 12+pieces:
				m["Events"] = color.GreenString("%d", events)
			case events < 20+pieces:
				m["Events"] = color.YellowString("%d", events)
			default:
				m["Events"] = color.RedString("%d", events)
			}
		}

		if len(st.Log) > 1 {
			start := time.Unix(int64(st.Log[0].Timestamp), 0)

			for _, sectorLog := range st.Log {
				if sectorLog.Kind == "event;sealingext.SectorProving" { // todo: figure out a good way to not hardcode
					end := time.Unix(int64(sectorLog.Timestamp), 0)
					dur := end.Sub(start)

					switch {
					case dur < 12*time.Hour:
						m["SealTime"] = color.GreenString("%s", dur)
					case dur < 24*time.Hour:
						m["SealTime"] = color.YellowString("%s", dur)
					default:
						m["SealTime"] = color.RedString("%s", dur)
					}

					break
				}
			}
		}

		tw.Write(m)
	}
	var buf bytes.Buffer
	tw.Flush(&buf)

	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) SectorsRefs(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	refs, err := (*minerApi).SectorsRefs(ctx)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	for name, refs := range refs {
		buf.WriteString(fmt.Sprintf("Block %s:\n", name))
		for _, ref := range refs {
			buf.WriteString(fmt.Sprintf("\t%d+%d %d bytes\n", ref.SectorID, ref.Offset, ref.Size))
		}
	}
	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) SectorsRemove(ctx context.Context, args ClientDealArgs) (string, error) {
	if !args.ReallyDoIt {
		return "", xerrors.Errorf("this is a command for advanced users, only use it if you are sure of what you are doing")
	}
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	id, err := strconv.ParseUint(args.SectorNum, 10, 64)
	if err != nil {
		return "", xerrors.Errorf("could not parse sector number: %w", err)
	}

	err = (*minerApi).SectorRemove(ctx, fabi.SectorNumber(id))
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *PublicTransactionPoolAPI) SectorMarkForUpgrade(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	id, err := strconv.ParseUint(args.SectorNum, 10, 64)
	if err != nil {
		return "", xerrors.Errorf("could not parse sector number: %w", err)
	}

	err = (*minerApi).SectorMarkForUpgrade(ctx, fabi.SectorNumber(id))
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *PublicTransactionPoolAPI) SectorsStartSeal(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	id, err := strconv.ParseUint(args.SectorNum, 10, 64)
	if err != nil {
		return "", xerrors.Errorf("could not parse sector number: %w", err)
	}

	err = (*minerApi).SectorStartSealing(ctx, fabi.SectorNumber(id))
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *PublicTransactionPoolAPI) SectorsSealDelay(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	hs, err := strconv.ParseUint(args.WaitDealsDelay, 10, 64)
	if err != nil {
		return "", xerrors.Errorf("could not parse WaitDealsDelay number: %w", err)
	}

	delay := hs * uint64(time.Minute)

	err = (*minerApi).SectorSetSealDelay(ctx, time.Duration(delay))
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *PublicTransactionPoolAPI) SectorsUpdate(ctx context.Context, args ClientDealArgs) (string, error) {
	if !args.ReallyDoIt {
		return "", xerrors.Errorf("this is a command for advanced users, only use it if you are sure of what you are doing")
	}
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	id, err := strconv.ParseUint(args.SectorNum, 10, 64)
	if err != nil {
		return "", xerrors.Errorf("could not parse sector number: %w", err)
	}

	var buf bytes.Buffer

	if _, ok := lsealing.ExistSectorStateList[lsealing.SectorState(args.NewState)]; !ok {
		buf.WriteString(fmt.Sprintf(" \"%s\" is not a valid state. Possible states for sectors are: \n", args.NewState))
		for state := range lsealing.ExistSectorStateList {
			buf.WriteString(fmt.Sprintf("%s\n", string(state)))
		}
		return buf.String(), nil
	}
	err = (*minerApi).SectorsUpdate(ctx, fabi.SectorNumber(id), lapi.SectorState(args.NewState))
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *PublicTransactionPoolAPI) ListPieces(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	pieceCids, err := (*minerApi).PiecesListPieces(ctx)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, pc := range pieceCids {
		buf.WriteString(pc.String() + "\n")
	}
	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) ListCidInfos(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	cids, err := (*minerApi).PiecesListCidInfos(ctx)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, c := range cids {
		buf.WriteString(c.String() + "\n")
	}
	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) GetPieceInfo(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	c, err := cid.Decode(args.Cid)
	if err != nil {
		return "", err
	}

	pi, err := (*minerApi).PiecesGetPieceInfo(ctx, c)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	tmpStr := fmt.Sprintln("Piece: ", pi.PieceCID)
	buf.WriteString(tmpStr)
	w := tabwriter.NewWriter(&buf, 4, 4, 2, ' ', 0)
	fmt.Fprintln(w, "Deals:\nDealID\tSectorID\tLength\tOffset")
	for _, d := range pi.Deals {
		fmt.Fprintf(w, "%d\t%d\t%d\t%d\n", d.DealID, d.SectorID, d.Length, d.Offset)
	}
	w.Flush()

	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) GetCIDInfo(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	c, err := cid.Decode(args.Cid)
	if err != nil {
		return "", err
	}

	ci, err := (*minerApi).PiecesGetCIDInfo(ctx, c)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	tmpStr := fmt.Sprintln("Info for: ", ci.CID)
	buf.WriteString(tmpStr)

	w := tabwriter.NewWriter(&buf, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, "PieceCid\tOffset\tSize\n")
	for _, loc := range ci.PieceBlockLocations {
		fmt.Fprintf(w, "%s\t%d\t%d\n", loc.PieceCID, loc.RelOffset, loc.BlockSize)
	}
	w.Flush()

	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) TransfersList(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	channels, err := (*minerApi).MarketListDataTransfers(ctx)
	if err != nil {
		return "", err
	}

	verbose := args.Verbose
	completed := args.Completed
	watch := args.Watch
	showFailed := args.ShowFailed
	if watch {
		channelUpdates, err := (*minerApi).MarketDataTransferUpdates(ctx)
		if err != nil {
			return "", err
		}

		for {
			tm.Clear() // Clear current screen

			tm.MoveCursor(1, 1)

			lcli.OutputDataTransferChannels(tm.Screen, channels, verbose, completed, showFailed)

			tm.Flush()

			select {
			case <-ctx.Done():
				return "", nil
			case channelUpdate := <-channelUpdates:
				var found bool
				for i, existing := range channels {
					if existing.TransferID == channelUpdate.TransferID &&
						existing.OtherPeer == channelUpdate.OtherPeer &&
						existing.IsSender == channelUpdate.IsSender &&
						existing.IsInitiator == channelUpdate.IsInitiator {
						channels[i] = channelUpdate
						found = true
						break
					}
				}
				if !found {
					channels = append(channels, channelUpdate)
				}
			}
		}
	}
	var buf bytes.Buffer
	lcli.OutputDataTransferChannels(&buf, channels, verbose, completed, showFailed)
	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) RestartTransfer(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	transferUint, err := strconv.ParseUint(args.TransferID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("Error reading transfer ID: %w", err)
	}
	transferID := datatransfer.TransferID(transferUint)
	initiator := args.Initiator
	var other peer.ID
	if pidstr := args.PeerId; pidstr != "" {
		p, err := peer.Decode(pidstr)
		if err != nil {
			return "", err
		}
		other = p
	} else {
		channels, err := (*minerApi).MarketListDataTransfers(ctx)
		if err != nil {
			return "", err
		}
		found := false
		for _, channel := range channels {
			if channel.IsInitiator == initiator && channel.TransferID == transferID {
				other = channel.OtherPeer
				found = true
				break
			}
		}
		if !found {
			return "", errors.New("unable to find matching data transfer")
		}
	}

	err = (*minerApi).MarketRestartDataTransfer(ctx, transferID, other, initiator)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (s *PublicTransactionPoolAPI) CancelTransfer(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	cancelTimeout, err := time.ParseDuration(args.CancelTimeout)
	if err != nil {
		cancelTimeout = 5 * time.Second
	}

	transferUint, err := strconv.ParseUint(args.TransferID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("Error reading transfer ID: %w", err)
	}
	transferID := datatransfer.TransferID(transferUint)
	initiator := args.Initiator
	var other peer.ID
	if pidstr := args.PeerId; pidstr != "" {
		p, err := peer.Decode(pidstr)
		if err != nil {
			return "", err
		}
		other = p
	} else {
		channels, err := (*minerApi).MarketListDataTransfers(ctx)
		if err != nil {
			return "", err
		}
		found := false
		for _, channel := range channels {
			if channel.IsInitiator == initiator && channel.TransferID == transferID {
				other = channel.OtherPeer
				found = true
				break
			}
		}
		if !found {
			return "", errors.New("unable to find matching data transfer")
		}
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, cancelTimeout)
	defer cancel()
	err = (*minerApi).MarketCancelDataTransfer(timeoutCtx, transferID, other, initiator)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (s *PublicTransactionPoolAPI) DagstoreListShards(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	shards, err := (*minerApi).DagstoreListShards(ctx)
	if err != nil {
		return "", err
	}

	if len(shards) == 0 {
		return "\n", nil
	}

	tw := tablewriter.New(
		tablewriter.Col("Key"),
		tablewriter.Col("State"),
		tablewriter.Col("Error"),
	)

	colors := map[string]color.Attribute{
		"ShardStateAvailable": color.FgGreen,
		"ShardStateServing":   color.FgBlue,
		"ShardStateErrored":   color.FgRed,
		"ShardStateNew":       color.FgYellow,
	}

	for _, s := range shards {
		m := map[string]interface{}{
			"Key": s.Key,
			"State": func() string {
				if c, ok := colors[s.State]; ok {
					return color.New(c).Sprint(s.State)
				}
				return s.State
			}(),
			"Error": s.Error,
		}
		tw.Write(m)
	}

	var buf bytes.Buffer
	tw.Flush(&buf)

	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) DagstoreInitializeShard(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	err := (*minerApi).DagstoreInitializeShard(ctx, args.StoreKey)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (s *PublicTransactionPoolAPI) DagstoreRecoverShard(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	err := (*minerApi).DagstoreRecoverShard(ctx, args.StoreKey)
	if err != nil {
		return "", err
	}

	return "success", nil
}

func (s *PublicTransactionPoolAPI) DagstoreInitializeAll(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	concurrency := args.Concurrency
	sealed := args.Sealed

	params := lapi.DagstoreInitializeAllParams{
		MaxConcurrency: int(concurrency),
		IncludeSealed:  sealed,
	}

	ch, err := (*minerApi).DagstoreInitializeAll(ctx, params)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	for {
		select {
		case evt, ok := <-ch:
			if !ok {
				refStr := buf.String()
				if len(refStr) == 0 {
					refStr = "\n"
				}
				return refStr, nil
			}
			_, _ = fmt.Fprint(&buf, color.New(color.BgHiBlack).Sprintf("(%d/%d)", evt.Current, evt.Total))
			_, _ = fmt.Fprint(&buf, " ")
			if evt.Event == "start" {
				_, _ = fmt.Fprintln(&buf, evt.Key, color.New(color.Reset).Sprint("STARTING"))
			} else {
				if evt.Success {
					_, _ = fmt.Fprintln(&buf, evt.Key, color.New(color.FgGreen).Sprint("SUCCESS"))
				} else {
					_, _ = fmt.Fprintln(&buf, evt.Key, color.New(color.FgRed).Sprint("ERROR"), evt.Error)
				}
			}

		case <-ctx.Done():
			return "", fmt.Errorf("aborted")
		}
	}
}

func (s *PublicTransactionPoolAPI) DagstoreGc(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	collected, err := (*minerApi).DagstoreGC(ctx)
	if err != nil {
		return "", err
	}

	if len(collected) == 0 {
		return fmt.Sprintln("no shards collected"), nil
	}

	var buf bytes.Buffer
	for _, e := range collected {
		if e.Error == "" {
			buf.WriteString(fmt.Sprintln(e.Key, color.New(color.FgGreen).Sprint("SUCCESS")))
		} else {
			buf.WriteString(fmt.Sprintln(e.Key, color.New(color.FgRed).Sprint("ERROR"), e.Error))
		}
	}

	refStr := buf.String()
	if len(refStr) == 0 {
		refStr = "\n"
	}
	return refStr, nil
}

func (s *PublicTransactionPoolAPI) PledgeSector(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerApi := s.b.GetMinerApi()

	id, err := (*minerApi).PledgeSector(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprint("Created CC sector: ", id.Number), nil
}

func (s *PublicTransactionPoolAPI) NetListen(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	minerAddr := s.b.GetMinerAddr()

	return *minerAddr, nil
}

func (s *PublicTransactionPoolAPI) WithdrawRemaining(ctx context.Context) (string, error) {
	lotuslog.SetupLogLevels()
	ethBackend := s.b
	addr, err := ethBackend.Coinbase()
	if s.b.ChainConfig().Bor == nil {
		return "", xerrors.Errorf("the operation not allow")
	}
	ks, err := fetchKeystore(ethBackend.AccountManager())
	if err != nil {
		return "", err
	}
	privateKey, err := ks.GetAccountPrivateKeyWithoutPass(accounts.Account{Address: addr})
	if err != nil {
		return "", err
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.quit == nil {
		s.quit = make(chan int)
	} else {
		s.quit <- 1
		s.quit = nil
		return color.New(color.FgRed).Sprintln("withdraw stop"), nil
	}

	go func() {
		s.mtx.Lock()
		quit := s.quit
		s.mtx.Unlock()
		period := time.Duration(s.b.ChainConfig().Bor.Period)
		timer := time.NewTicker(time.Second * (period*2 + 2))
		for {
			select {
			case <-timer.C:
				if s.withdrawRemaining(context.Background(), addr, privateKey) {
					timer.Stop()
					close(quit)
					s.mtx.Lock()
					s.quit = nil
					s.mtx.Unlock()
					return
				}
			case <-quit:
				timer.Stop()
				close(quit)
				return
			}
		}
	}()

	return color.New(color.FgRed).Sprintln("withdraw start"), nil
}

func (s *PublicTransactionPoolAPI) withdrawRemaining(ctx context.Context, addr common.Address, privateKey *ecdsa.PrivateKey) bool {
	expireFile, err := file_store.FileStoreCli.GetExpireFile(&bind.CallOpts{From: addr, Context: ctx})
	if err != nil {
		lotusLog.Errorf("GetExpireFile error: %s", err)
		return false
	}
	entireExpireFile, err := file_store.FileStoreCli.GetExpireFileEntire(&bind.CallOpts{From: addr, Context: ctx})
	if err != nil {
		lotusLog.Errorf("GetExpireFile error: %s", err)
		return false
	}
	if expireFile.StorageType == 0 && entireExpireFile.StorageType == 0 {
		return true
	}
	if expireFile.StorageType != 0 {
		data, err := file_store.FileStoreCli.Pack4WithdrawRemaining(expireFile.OriHash, expireFile.Index, expireFile.StorageType)
		if err != nil {
			lotusLog.Errorf("Pack4WithdrawRemaining error: %s", err)
			return false
		}
		auth, err := file_store.FileStoreCli.GenerateAuthObj(privateKey, s.b.ChainConfig().ChainID, addr, data)
		if err != nil {
			lotusLog.Errorf("GenerateAuthObj error: %s", err)
			return false
		}
		_, err = file_store.FileStoreCli.WithdrawRemaining(auth, expireFile.OriHash, expireFile.Index, expireFile.StorageType)
		if err != nil {
			lotusLog.Errorf("WithdrawRemaining error: %s", err)
			return false
		}
	}
	if entireExpireFile.StorageType != 0 {
		data, err := file_store.FileStoreCli.Pack4WithdrawRemaining(entireExpireFile.OriHash, entireExpireFile.Index, entireExpireFile.StorageType)
		if err != nil {
			lotusLog.Errorf("Pack4WithdrawRemaining error: %s", err)
			return false
		}
		auth, err := file_store.FileStoreCli.GenerateAuthObj(privateKey, s.b.ChainConfig().ChainID, addr, data)
		if err != nil {
			lotusLog.Errorf("GenerateAuthObj error: %s", err)
			return false
		}
		_, err = file_store.FileStoreCli.WithdrawRemaining(auth, entireExpireFile.OriHash, entireExpireFile.Index, entireExpireFile.StorageType)
		if err != nil {
			lotusLog.Errorf("WithdrawRemaining error: %s", err)
			return false
		}
	}
	return false
}

func maybeStr(c bool, col color.Attribute, s string) string {
	if !c {
		return ""
	}

	return color.New(col).Sprint(s)
}

func yesno(b bool) string {
	if b {
		return color.GreenString("YES")
	}
	return color.RedString("NO")
}

func APIEndpoint(dir string) (string, error) {
	p := filepath.Join(dir, "api")

	f, err := os.Open(p)
	if os.IsNotExist(err) {
		return "", repo.ErrNoAPIEndpoint
	} else if err != nil {
		return "", err
	}
	defer f.Close() //nolint: errcheck // Read only op

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", xerrors.Errorf("failed to read %q: %w", p, err)
	}
	strma := string(data)
	strma = strings.TrimSpace(strma)
	return strma, nil
}
func (s *PublicStorageAPI) GetMiners4Entire(ctx context.Context, args ClientDealArgs) []MinerStatus {
	p2pServer := s.b.GetP2pServer()
	peers := p2pServer.Peers()

	if args.MinerNum > len(peers) {
		args.MinerNum = len(peers)
	}
	storeKey := args.StoreKey
	storeKeyHashStr := common.BuildNewOriHashWithSha256(storeKey)
	storeHashByte := common.HexSTrToByte32(storeKeyHashStr)
	// get miners who has store the same oriHash
	storeMiners, err := file_store.FileStoreCli.GetStoreMiners4Entire(nil, storeHashByte)
	if err != nil {
		storeMiners = nil
	}
	if args.MinerNum > 0 {
		return s.CommonGetMiners(args.Type, storeMiners, args.MinerNum)
	}
	return s.commonGetLocalMiner(args.Type, storeMiners)
}

func (s *PublicStorageAPI) GetMiners(ctx context.Context, args ClientDealArgs) []MinerStatus {
	p2pServer := s.b.GetP2pServer()
	peers := p2pServer.Peers()
	if args.MinerNum > len(peers) {
		args.MinerNum = len(peers)
	}
	oriHashByte := common.HexSTrToByte32(args.OriHash)
	// get file info from contract.
	storeMiners, err := file_store.FileStoreCli.GetStoreMiners(nil, oriHashByte)
	if err != nil {
		storeMiners = nil
	}
	if args.MinerNum > 0 {
		return s.CommonGetMiners(args.Type, storeMiners, args.MinerNum)
	}
	return s.commonGetLocalMiner(args.Type, storeMiners)
}

func (s *PublicStorageAPI) commonGetLocalMiner(queryType string, storeMiners [][32]byte) []MinerStatus {
	p2pServer := s.b.GetP2pServer()
	localNode := p2pServer.LocalNode()
	minerId := common.HexSTrToByte32(localNode.ID().String())
	minerIdStr := "0x" + localNode.ID().String()
	minerInfo := file_store.FileStoreCli.GetMinerInfoByMinerId(nil, minerId)
	var miners []MinerStatus
	isUsed := getMinerStatus(minerId, storeMiners)
	if isUsed && queryType == UNSTORED {
		return miners
	}
	if !isUsed && queryType == STORED {
		return miners
	}
	minerStatus := MinerStatus{
		MinerId:   minerIdStr,
		Status:    isUsed,
		PublicKey: minerInfo.PublicKey,
	}
	miners = append(miners, minerStatus)
	return miners
}

func (s *PublicStorageAPI) CommonGetMiners(queryType string, storeMiners [][32]byte, minerNum int) []MinerStatus {
	if queryType == "" {
		queryType = UNSTORED
	}
	p2pServer := s.b.GetP2pServer()
	peers := p2pServer.Peers()
	if minerNum > 0 {
		count := 0
		var miners []MinerStatus
		freeSpaceMap := make(map[string]uint64)
		for _, peer := range peers {
			if count > minerNum {
				break
			}
			minerId := common.HexSTrToByte32(peer.ID().String())
			minerIdStr := "0x" + peer.ID().String()
			minerInfo := file_store.FileStoreCli.GetMinerInfoByMinerId(nil, minerId)
			if minerInfo.PublicKey != "" {
				isUsed := getMinerStatus(minerId, storeMiners)
				if isUsed && queryType == UNSTORED {
					continue
				}
				if !isUsed && queryType == STORED {
					continue
				}
				minerStatus := MinerStatus{
					MinerId:   minerIdStr,
					Status:    isUsed,
					PublicKey: minerInfo.PublicKey,
				}
				var freeSpace uint64
				fromChain := true
				minerFreeStorageKey := fmt.Sprintf("minerFreeStorageKey_%s", minerIdStr)
				if GCacheForMemoryEnable {
					val, err := GCacheForMemory.Get(minerFreeStorageKey)
					if err == nil && val != nil {
						freeSpace = val.(uint64)
						fromChain = false
						lotusLog.Infof("get miner(id:%s) freeSpace from cache,value is %s", minerIdStr, freeSpace)
					}
				}
				if fromChain {
					freeSpace = getFreeSpaceFromChain(minerIdStr)
					lotusLog.Infof("get miner(id:%s) freeSpace from chain,value is %s", minerIdStr, freeSpace)
					// save to cache.
					GCacheForMemory.SetWithExpire(minerFreeStorageKey, freeSpace, time.Minute*30)
				}
				freeSpaceMap[minerIdStr] = freeSpace
				miners = append(miners, minerStatus)
				count++
			}
		}
		if GCacheForMemoryEnable && len(miners) > 0 {
			// calculate ratio fro miners.
			calculateFreeStorageRatioAndSort(miners, freeSpaceMap)
		}
		return miners
	}
	return nil
}

// get free space from chain.
func getFreeSpaceFromChain(minerIdStr string) uint64 {
	minerId := common.HexSTrToByte32(minerIdStr)
	minerAddr := file_store.FileStoreCli.GetMinerAddr(nil, minerId)
	lotusLog.Info("minerAddr:", minerAddr)
	promiseSize := w3fsStorageManager.GlobalStorageManagerClient.GetValidatorPromise(minerAddr)
	usedSize := w3fsStorageManager.GlobalStorageManagerClient.GetValidatorStorageSize(minerAddr)
	leftSize := promiseSize - usedSize
	if leftSize < 0 {
		leftSize = 0
	}
	return leftSize
}

// calculate every miner's ratio in the miner list
func calculateFreeStorageRatioAndSort(miners []MinerStatus, minerMap map[string]uint64) {
	var totalFreeSpace uint64
	for _, miner := range miners {
		totalFreeSpace = totalFreeSpace + minerMap[miner.MinerId]
	}
	// calculate
	for i, miner := range miners {
		if totalFreeSpace <= 0 {
			miners[i].FreeSpaceRatio = 0
		} else {
			freeSpace := minerMap[miner.MinerId]
			radio := float64(freeSpace) / float64(totalFreeSpace)
			miners[i].FreeSpaceRatio = common.FormatFloat2(radio)
		}
	}
	// sort by freeSpaceRatio.
	sort.Slice(miners, func(i, j int) bool { return miners[i].FreeSpaceRatio > miners[j].FreeSpaceRatio })
}

func (s *PublicTransactionPoolAPI) LotusImport(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	mgr := s.b.ClientManager()

	exist, path := checkIfFileExist(args.File, mgr.Repo.Path()+"/storage-file")
	if !exist {
		return "", errors.New("file not exist, please upload file first or check the path")
	}

	cid, err := ClientImport(ctx, mgr, path)
	if err != nil {
		lotusLog.Errorf("client import %s", err)
		return "", err
	}
	return cid.String(), nil
}

func (s *PublicStorageAPI) GetStorageStatus4Entire(ctx context.Context, args ClientDealArgs) DealStatus {
	args.HeadFlag = false
	args.StorageType = ENTIRE_FILE
	return s.GetStorageStatus(ctx, args)
}

func (s *PublicStorageAPI) GetStorageStatus(ctx context.Context, args ClientDealArgs) DealStatus {
	if len(args.MinerId) == 0 {
		return newDealStatus(DEAL_ERROR, "param minerId cannot be empty!")
	}
	if len(args.OriHash) == 0 && len(args.StoreKey) == 0 {
		return newDealStatus(DEAL_ERROR, "param oriHash/storeKey cannot be empty!")
	}
	var hashStr string
	storageType := args.StorageType
	if storageType == ENTIRE_FILE {
		hashStr = common.BuildNewOriHashWithSha256(args.StoreKey)
	} else {
		hashStr = args.OriHash
	}
	hashByte := common.HexSTrToByte32(hashStr)
	// get miner's address info.
	minerIdByte := common.HexSTrToByte32(args.MinerId)
	_, peerId, err := s.getAddrInfoByMinerId(minerIdByte)
	if err != nil {
		return newDealStatus(DEAL_ERROR, "getAddrInfoByMinerId error: MinerId: "+args.MinerId)
	}
	mgr := s.b.ClientManager()
	// decode cid
	var strCid string
	if storageType == ENTIRE_FILE {
		strCid = mgr.GetFileHashAndFlagToCid(args.StoreKey, args.HeadFlag, peerId)
	} else {
		strCid = mgr.GetFileHashAndFlagToCid(args.OriHash, args.HeadFlag, peerId)
	}
	dataCid, _ := cid.Decode(strCid)
	var msg = "deal init"
	status := DEAL_INIT
	hasStored := mgr.HasStored(dataCid, peerId)
	if mgr.HasStoring(dataCid, peerId) || hasStored {
		status = DEAL_ING
		msg = "deal being processed"
	}
	var retStatus uint8
	if storageType == ENTIRE_FILE {
		retStatus = file_store.FileStoreCli.CheckStorage4Entire(nil, hashByte, minerIdByte)
	} else {
		retStatus = file_store.FileStoreCli.CheckStorage(nil, hashByte, args.HeadFlag, minerIdByte)
	}

	if hasStored && retStatus == file_store.CHECK_FINISH {
		status = DEAL_SUCCESS
		msg = "success"
	}
	return newDealStatus(status, msg)
}

func (s *PublicTransactionPoolAPI) ClientStorage(ctx context.Context, args ClientDealArgs) (string, error) {
	lotuslog.SetupLogLevels()
	peerId := args.PeerId
	cid, err := cid.Decode(args.Cid)
	if err != nil {
		lotusLog.Errorf("cid parse  %s", err)
		return "", err
	}

	addrInfo, peerId := s.getAddrInfo(args)

	mgr := s.b.ClientManager()
	mgr.DealClient.Start(ctx)

	err = mgr.Host.Connect(ctx, addrInfo)
	if err != nil {
		lotusLog.Errorf("client connect  %s", err)
		return "", err
	}

	mgr.StartDeal(ctx, cid, peerId)
	//	mgr.DealClient.Stop()

	return args.Cid, nil
}

func (s *PublicStorageAPI) findMiner4EntireFile(storeKey string) (file_store.FileStoreStructFileMinerInfo, error) {
	storeKeyHash := common.BuildNewOriHashWithSha256(storeKey)
	storeKeyHashByte := common.HexSTrToByte32(storeKeyHash)
	// get file info from contract.
	return file_store.FileStoreCli.FindMiner4EntireFile(nil, storeKeyHashByte)
}

func (s *PublicStorageAPI) findMiner4File(oriHashStr string, headFlag bool) (file_store.FileStoreStructFileMinerInfo, error) {
	oriHash := common.HexSTrToByte32(oriHashStr)
	// get file info from contract.
	return file_store.FileStoreCli.FindMiner4File(nil, oriHash, headFlag)
}

func Createretrievalkey(orihash string, headflag bool, txhash string) string {
	ishash := CheckTxhash(txhash)
	var key string
	if headflag {
		if ishash {
			key = txhash
		} else {
			key = orihash + "_head"
		}
	} else {
		key = orihash + "_body"
	}

	return key
}

func (s *PublicStorageAPI) GetRetrievalStatus4Entire(ctx context.Context, args ClientDealArgs) DealStatus {
	args.HeadFlag = false
	args.StorageType = ENTIRE_FILE
	return s.GetRetrievalStatus(ctx, args)
}

// @see DealStatus's define
func (s *PublicStorageAPI) GetRetrievalStatus(ctx context.Context, args ClientDealArgs) DealStatus {
	oriHashStr := args.OriHash
	if args.StorageType == ENTIRE_FILE {
		oriHashStr = common.BuildNewOriHashWithSha256(args.StoreKey)
	}

	log.Info("Request GetRetrievalStatus", "oriHash", oriHashStr, "headFlag", args.HeadFlag, "storeKey", args.StoreKey, "storageType", args.StorageType)

	var gKey string
	gKey = Createretrievalkey(oriHashStr, args.HeadFlag, args.TxHash)
	if GCacheForFileStoreEnable && GCacheForFileStore.Has(gKey) {
		// is cache
		log.Info("GetRetrievalStatus have the cache", "oriHash", oriHashStr, "headFlag", args.HeadFlag, "storeKey", args.StoreKey, "storageType", args.StorageType)
		ds := DealStatus{DEAL_SUCCESS, "cache", time.Now().UTC()}
		return ds
	}

	lotuslog.SetupLogLevels()
	mgr := s.b.ClientManager()
	// get params from ctx
	// headFlag
	headFlag := args.HeadFlag
	ishash := CheckTxhash(args.TxHash)
	rs := mgr.GetRetrievalStatusExt(gKey, ishash)
	if rs.Status == DEAL_SUCCESS && GCacheForFileStoreEnable && !GCacheForFileStore.Has(gKey) {
		// add cache
		err := GCacheForFileStore.Set(gKey, rs.Cid)
		if err != nil {
			log.Error("add cache fail", "oriHashStr", oriHashStr, "headFlag", headFlag, "cid", rs.Cid, "err", err.Error())
		}
	}

	ds := DealStatus{rs.Status, rs.Message, rs.Timestamp}
	return ds
}

func (s *PublicStorageAPI) retrieval(ctx context.Context, args ClientDealArgs, auth bool) DealStatus {
	// get param oriHash from ctx
	oriHashStr := args.OriHash
	storageType := args.StorageType
	if storageType == ENTIRE_FILE {
		// compute sha256 as new OriHash
		oriHashStr = common.BuildNewOriHashWithSha256(args.StoreKey)
	}
	// get param headFlag from ctx
	headFlag := args.HeadFlag

	log.Info("Request retrieval", "oriHash", oriHashStr, "headFlag", headFlag, "storeKey", args.StoreKey, "storageType", storageType)

	hashByte := common.HexSTrToByte32(oriHashStr)
	// get fileExt/fileSize from contract.
	fileAttr, err := getFileAttr(hashByte, storageType)
	if err != nil {
		return newDealStatus(DEAL_ERROR, err.Error())
	}

	gKey := Createretrievalkey(oriHashStr, headFlag, args.TxHash)
	if GCacheForFileStoreEnable {
		_, err := GCacheForFileStore.Get(gKey)
		if err != gcache.KeyNotFoundError {
			return newDealStatus(DEAL_SUCCESS, fileAttr)
		}
	}

	lotuslog.SetupLogLevels()
	mgr := s.b.ClientManager()
	isHash := CheckTxhash(args.TxHash)
	retrieveStatus := mgr.GetRetrievalStatusExt(gKey, isHash)
	if retrieveStatus.Status == DEAL_ING {
		errMsg := "During file retrieval, please do not submit again!"
		lotusLog.Infof(errMsg)
		return newDealStatus(DEAL_ERROR, errMsg)
	} else if retrieveStatus.Status == DEAL_SUCCESS {
		if GCacheForFileStoreEnable && !GCacheForFileStore.Has(gKey) {
			GCacheForFileStore.Set(gKey, retrieveStatus.Cid)
		}
		return newDealStatus(DEAL_SUCCESS, fileAttr)
	}

	mgr.Retrieval.Start(ctx)
	// logPre
	var logPrefix = fmt.Sprintf("[oriHash=%s]", oriHashStr)
	var out file_store.FileStoreStructFileMinerInfo
	if storageType == ENTIRE_FILE {
		out, err = s.findMiner4EntireFile(args.StoreKey)
		if err != nil {
			lotusLog.Errorf(logPrefix + err.Error())
			return newDealStatus(DEAL_ERROR, err.Error())
		}
	} else {
		out, err = s.findMiner4File(oriHashStr, headFlag)
		if err != nil {
			lotusLog.Errorf(logPrefix + err.Error())
			return newDealStatus(DEAL_ERROR, err.Error())
		}
	}

	if len(out.MinerIds) == 0 {
		errMsg := "The OriHash could not be found in the contract, 1 file expire, 2 store failed."
		lotusLog.Errorf(logPrefix + errMsg)
		return newDealStatus(DEAL_ERROR, errMsg)
	}

	dataCid, _ := cid.Decode(out.FileCid)
	// Start Retrieval file
	if auth == true {
		mgr.MarkeAutStatus(args.TxHash, DEAL_INIT)
		mgr.MarkRetrieveType(dataCid, headFlag, auth, args.TxHash)
	}

	err = s.ClientRetrieveNew(oriHashStr, headFlag, dataCid, out)
	if err != nil {
		return newDealStatus(DEAL_ERROR, err.Error())
	}

	return newDealStatus(DEAL_SUCCESS, fileAttr)
}

func getFileAttr(hashByte [32]byte, storageType string) (string, error) {
	var fileBaseInfo file_store.FileStoreStructBaseInfo
	var err2 error
	if storageType == ENTIRE_FILE {
		fileBaseInfo, err2 = file_store.FileStoreCli.GetBaseInfo4Entire(nil, hashByte)
		if err2 != nil {
			return "", errors.New("cannot get entire file's extension from the store contract")
		}
	} else {
		fileBaseInfo, err2 = file_store.FileStoreCli.GetBaseInfo(nil, hashByte)
		if err2 != nil {
			return "", errors.New("cannot get file's extension from the store contract")
		}
	}
	// important: set msg value as fileExt,then return as 'Message'
	fileInfo := FileInfo{}
	fileInfo.FileExt = fileBaseInfo.FileExt
	fileInfo.FileSize = fileBaseInfo.FileSize
	fileInfoJsonBytes, _ := json.Marshal(fileInfo)
	return string(fileInfoJsonBytes), nil
}

func (s *PublicStorageAPI) GetAuthorizedHead(ctx context.Context, args ClientDealArgs) DealStatus {
	return s.retrieval(ctx, args, true)
}

func (s *PublicStorageAPI) ClientRetrieval4Entire(ctx context.Context, args ClientDealArgs) DealStatus {
	args.HeadFlag = false
	args.StorageType = ENTIRE_FILE
	return s.retrieval(ctx, args, false)
}

func (s *PublicStorageAPI) ClientRetrieval(ctx context.Context, args ClientDealArgs) DealStatus {
	return s.retrieval(ctx, args, false)
}

// parseNode parses a node record and verifies its signature.
func parseNode(source string) (*enode.Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return enode.ParseV4(source)
	}
	r, err := parseRecord(source)
	if err != nil {
		return nil, err
	}
	return enode.New(enode.ValidSchemes, r)
}

// parseRecord parses a node record from hex, base64, or raw binary input.
func parseRecord(source string) (*enr.Record, error) {
	bin := []byte(source)
	if d, ok := decodeRecordHex(bytes.TrimSpace(bin)); ok {
		bin = d
	} else if d, ok := decodeRecordBase64(bytes.TrimSpace(bin)); ok {
		bin = d
	}
	var r enr.Record
	err := rlp.DecodeBytes(bin, &r)
	return &r, err
}

func decodeRecordHex(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("0x")) {
		b = b[2:]
	}
	dec := make([]byte, hex.DecodedLen(len(b)))
	_, err := hex.Decode(dec, b)
	return dec, err == nil
}

func decodeRecordBase64(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("enr:")) {
		b = b[4:]
	}
	dec := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dec, b)
	return dec[:n], err == nil
}

func (s *PublicStorageAPI) getAddrInfoByMinerId(minerId [32]byte) (peer.AddrInfo, string, error) {
	var minerIdStr = common.Byte32ToHexStr(minerId)
	// get from contract
	minerInfo := file_store.FileStoreCli.GetMinerInfoByMinerId(nil, minerId)
	peerid := minerInfo.PeerId
	addr := minerInfo.PeerAddr
	if addr == "" || peerid == "" {
		return peer.AddrInfo{}, "", fmt.Errorf("can't get peerId or peerAddr from chain data,minerId:%s", minerIdStr)
	}
	m1, err := ma.NewMultiaddr(addr)
	if err != nil {
		return peer.AddrInfo{}, "", fmt.Errorf("NewMultiaddr error: %s, minerId:%s", err, minerIdStr)
	}
	ccid, err := cid.Decode(peerid)
	if err != nil {
		return peer.AddrInfo{}, "", fmt.Errorf("cid decode error: %s, minerId:%s", err, minerIdStr)
	}

	remotepeerid, _ := peer.FromCid(ccid)

	addrInfo := peer.AddrInfo{
		Addrs: []ma.Multiaddr{m1},
		ID:    remotepeerid,
	}
	return addrInfo, peerid, nil
}

func (s *PublicTransactionPoolAPI) getAddrInfo(args ClientDealArgs) (peer.AddrInfo, string) {
	var addr string
	var peerid string
	p2pServer := s.b.GetP2pServer()
	peers := p2pServer.Peers()
	//shuffle the array
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(peers), func(i int, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})

	for _, peer := range peers {
		var remoteListenPort uint64
		entryListenPort := enr.WithEntry("Bor_Remote_Port", &remoteListenPort)
		peer.Node().Load(entryListenPort)
		if remoteListenPort != 0 {
			Enode := peer.Node().URLv4()
			tmpIndex := strings.LastIndex(Enode, ":")
			Enode = Enode[:tmpIndex+1] + strconv.FormatUint(remoteListenPort, 10)
			enode, err := parseNode(Enode)
			if err != nil {
				continue
			}
			localNode := p2pServer.DiscV4().Resolve(enode)
			entryAddr := enr.WithEntry("Lotus_IP", &addr)
			localNode.Load(entryAddr)
			entryPeerId := enr.WithEntry("Lotus_PeerId", &peerid)
			localNode.Load(entryPeerId)
			if addr != "" && peerid != "" {
				break
			}
		}
	}
	if addr == "" || peerid == "" {
		localNode := p2pServer.Self()
		entryAddr := enr.WithEntry("Lotus_IP", &addr)
		localNode.Load(entryAddr)
		entryPeerId := enr.WithEntry("Lotus_PeerId", &peerid)
		localNode.Load(entryPeerId)
	}

	m1, _ := ma.NewMultiaddr(addr)
	ccid, _ := cid.Decode(peerid)

	remotepeerid, _ := peer.FromCid(ccid)

	addrInfo := peer.AddrInfo{
		Addrs: []ma.Multiaddr{m1},
		ID:    remotepeerid,
	}
	return addrInfo, peerid
}

func getMinerConfig() []string {
	content, err := ioutil.ReadFile("miner.txt")
	if err != nil && !os.IsNotExist(err) {
		lotusLog.Errorf("error reading miner:", err)
	}
	var miners []string
	for _, a := range bytes.Split(content, []byte("\n")) {
		if len(a) > 0 && a[0] != '#' {
			miners = append(miners, string(a))
		}
	}
	return miners
}

func ClientImport(ctx context.Context, api ClientManager, file string) (cid cid.Cid, err error) {
	c, err := clientImport(ctx, api, file)
	if err != nil {
		return cid, err
	}

	return c.Root, nil
}

func clientImport(ctx context.Context, api ClientManager, file string) (res *lapi.ImportRes, err error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	ref := lapi.FileRef{
		Path:  absPath,
		IsCAR: false,
	}
	c, err := api.ClientImport(ctx, ref)
	if err != nil {
		return nil, err
	}
	encoder, err := GetCidEncoder()
	if err != nil {
		return nil, err
	}

	lotusLog.Infof("Import %d, Root ", c.ImportID)

	lotusLog.Infof(encoder.Encode(c.Root))
	return c, nil
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (s *PublicTransactionPoolAPI) FillTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and obtain rlp
	tx := args.toTransaction()
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(input); err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, tx)
}

// Sign calculates an ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (s *PublicTransactionPoolAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Sign the requested hash with the wallet
	signature, err := wallet.SignText(account, data)
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, err
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *PublicTransactionPoolAPI) SignTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil && (args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil) {
		return nil, fmt.Errorf("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Before actually sign the transaction, ensure the transaction fee is reasonable.
	tx := args.toTransaction()
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := s.sign(args.from(), tx)
	if err != nil {
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (s *PublicTransactionPoolAPI) PendingTransactions() ([]*RPCTransaction, error) {
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := make(map[common.Address]struct{})
	for _, wallet := range s.b.AccountManager().Wallets() {
		for _, account := range wallet.Accounts() {
			accounts[account.Address] = struct{}{}
		}
	}
	curHeader := s.b.CurrentHeader()
	transactions := make([]*RPCTransaction, 0, len(pending))
	for _, tx := range pending {
		from, _ := types.Sender(s.signer, tx)
		if _, exists := accounts[from]; exists {
			transactions = append(transactions, newRPCPendingTransaction(tx, curHeader, s.b.ChainConfig()))
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *PublicTransactionPoolAPI) Resend(ctx context.Context, sendArgs TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	if sendArgs.Nonce == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
	}
	if err := sendArgs.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	matchTx := sendArgs.toTransaction()

	// Before replacing the old transaction, ensure the _new_ transaction fee is reasonable.
	var price = matchTx.GasPrice()
	if gasPrice != nil {
		price = gasPrice.ToInt()
	}
	var gas = matchTx.Gas()
	if gasLimit != nil {
		gas = uint64(*gasLimit)
	}
	if err := checkTxFee(price, gas, s.b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}
	// Iterate the pending list for replacement
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return common.Hash{}, err
	}
	for _, p := range pending {
		wantSigHash := s.signer.Hash(matchTx)
		pFrom, err := types.Sender(s.signer, p)
		if err == nil && pFrom == sendArgs.from() && s.signer.Hash(p) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				sendArgs.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				sendArgs.Gas = gasLimit
			}
			signedTx, err := s.sign(sendArgs.from(), sendArgs.toTransaction())
			if err != nil {
				return common.Hash{}, err
			}
			if err = s.b.SendTx(ctx, signedTx); err != nil {
				return common.Hash{}, err
			}
			return signedTx.Hash(), nil
		}
	}
	return common.Hash{}, fmt.Errorf("transaction %#x not found", matchTx.Hash())
}

// PublicDebugAPI is the collection of Ethereum APIs exposed over the public
// debugging endpoint.
type PublicDebugAPI struct {
	b Backend
}

// NewPublicDebugAPI creates a new API definition for the public debug methods
// of the Ethereum service.
func NewPublicDebugAPI(b Backend) *PublicDebugAPI {
	return &PublicDebugAPI{b: b}
}

// GetBlockRlp retrieves the RLP encoded for of a single block.
func (api *PublicDebugAPI) GetBlockRlp(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	encoded, err := rlp.EncodeToBytes(block)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", encoded), nil
}

// TestSignCliqueBlock fetches the given block number, and attempts to sign it as a clique header with the
// given address, returning the address of the recovered signature
//
// This is a temporary method to debug the externalsigner integration,
// TODO: Remove this method when the integration is mature
func (api *PublicDebugAPI) TestSignCliqueBlock(ctx context.Context, address common.Address, number uint64) (common.Address, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return common.Address{}, fmt.Errorf("block #%d not found", number)
	}
	header := block.Header()
	header.Extra = make([]byte, 32+65)
	encoded := clique.CliqueRLP(header)

	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: address}
	wallet, err := api.b.AccountManager().Find(account)
	if err != nil {
		return common.Address{}, err
	}

	signature, err := wallet.SignData(account, accounts.MimetypeClique, encoded)
	if err != nil {
		return common.Address{}, err
	}
	sealHash := clique.SealHash(header).Bytes()
	log.Info("test signing of clique block",
		"Sealhash", fmt.Sprintf("%x", sealHash),
		"signature", fmt.Sprintf("%x", signature))
	pubkey, err := crypto.Ecrecover(sealHash, signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PublicDebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return spew.Sdump(block), nil
}

// SeedHash retrieves the seed hash of a block.
func (api *PublicDebugAPI) SeedHash(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return fmt.Sprintf("0x%x", ethash.SeedHash(number)), nil
}

// PrivateDebugAPI is the collection of Ethereum APIs exposed over the private
// debugging endpoint.
type PrivateDebugAPI struct {
	b Backend
}

// NewPrivateDebugAPI creates a new API definition for the private debug methods
// of the Ethereum service.
func NewPrivateDebugAPI(b Backend) *PrivateDebugAPI {
	return &PrivateDebugAPI{b: b}
}

// ChaindbProperty returns leveldb properties of the key-value database.
func (api *PrivateDebugAPI) ChaindbProperty(property string) (string, error) {
	if property == "" {
		property = "leveldb.stats"
	} else if !strings.HasPrefix(property, "leveldb.") {
		property = "leveldb." + property
	}
	return api.b.ChainDb().Stat(property)
}

// ChaindbCompact flattens the entire key-value database into a single level,
// removing all unused slots and merging all keys.
func (api *PrivateDebugAPI) ChaindbCompact() error {
	for b := byte(0); b < 255; b++ {
		log.Info("Compacting chain database", "range", fmt.Sprintf("0x%0.2X-0x%0.2X", b, b+1))
		if err := api.b.ChainDb().Compact([]byte{b}, []byte{b + 1}); err != nil {
			log.Error("Database compaction failed", "err", err)
			return err
		}
	}
	return nil
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *PrivateDebugAPI) SetHead(number hexutil.Uint64) {
	api.b.SetHead(uint64(number))
}

// PublicNetAPI offers network related RPC methods
type PublicNetAPI struct {
	net            *p2p.Server
	networkVersion uint64
}

// NewPublicNetAPI creates a new net API instance.
func NewPublicNetAPI(net *p2p.Server, networkVersion uint64) *PublicNetAPI {
	return &PublicNetAPI{net, networkVersion}
}

// Listening returns an indication if the node is listening for network connections.
func (s *PublicNetAPI) Listening() bool {
	return true // always listening
}

// PeerCount returns the number of connected peers
func (s *PublicNetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(s.net.PeerCount())
}

// Version returns the current ethereum protocol version.
func (s *PublicNetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.Ether)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f ether) exceeds the configured cap (%.2f ether)", feeFloat, cap)
	}
	return nil
}

// toHexSlice creates a slice of hex-strings based on []byte.
func toHexSlice(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = hexutil.Encode(b[i])
	}
	return r
}

type stateMeta struct {
	i     int
	col   color.Attribute
	state lsealing.SectorState
}

var stateOrder = map[lsealing.SectorState]stateMeta{}
var stateList = []stateMeta{
	{col: 39, state: "Total"},
	{col: color.FgGreen, state: lsealing.Proving},

	{col: color.FgBlue, state: lsealing.Empty},
	{col: color.FgBlue, state: lsealing.WaitDeals},
	{col: color.FgBlue, state: lsealing.AddPiece},

	{col: color.FgRed, state: lsealing.UndefinedSectorState},
	{col: color.FgYellow, state: lsealing.Packing},
	{col: color.FgYellow, state: lsealing.GetTicket},
	{col: color.FgYellow, state: lsealing.PreCommit1},
	{col: color.FgYellow, state: lsealing.PreCommit2},
	{col: color.FgYellow, state: lsealing.PreCommitting},
	{col: color.FgYellow, state: lsealing.PreCommitWait},
	{col: color.FgYellow, state: lsealing.SubmitPreCommitBatch},
	{col: color.FgYellow, state: lsealing.PreCommitBatchWait},
	{col: color.FgYellow, state: lsealing.WaitSeed},
	{col: color.FgYellow, state: lsealing.Committing},
	{col: color.FgYellow, state: lsealing.CommitFinalize},
	{col: color.FgYellow, state: lsealing.SubmitCommit},
	{col: color.FgYellow, state: lsealing.CommitWait},
	{col: color.FgYellow, state: lsealing.SubmitCommitAggregate},
	{col: color.FgYellow, state: lsealing.CommitAggregateWait},
	{col: color.FgYellow, state: lsealing.FinalizeSector},

	{col: color.FgCyan, state: lsealing.Terminating},
	{col: color.FgCyan, state: lsealing.TerminateWait},
	{col: color.FgCyan, state: lsealing.TerminateFinality},
	{col: color.FgCyan, state: lsealing.TerminateFailed},
	{col: color.FgCyan, state: lsealing.Removing},
	{col: color.FgCyan, state: lsealing.Removed},

	{col: color.FgRed, state: lsealing.FailedUnrecoverable},
	{col: color.FgRed, state: lsealing.AddPieceFailed},
	{col: color.FgRed, state: lsealing.SealPreCommit1Failed},
	{col: color.FgRed, state: lsealing.SealPreCommit2Failed},
	{col: color.FgRed, state: lsealing.PreCommitFailed},
	{col: color.FgRed, state: lsealing.ComputeProofFailed},
	{col: color.FgRed, state: lsealing.CommitFailed},
	{col: color.FgRed, state: lsealing.CommitFinalizeFailed},
	{col: color.FgRed, state: lsealing.PackingFailed},
	{col: color.FgRed, state: lsealing.FinalizeFailed},
	{col: color.FgRed, state: lsealing.Faulty},
	{col: color.FgRed, state: lsealing.FaultReported},
	{col: color.FgRed, state: lsealing.FaultedFinal},
	{col: color.FgRed, state: lsealing.RemoveFailed},
	{col: color.FgRed, state: lsealing.DealsExpired},
	{col: color.FgRed, state: lsealing.RecoverDealIDs},
}
