package bor

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	sealing2 "github.com/ethereum/go-ethereum/sealing"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	proof5 "github.com/filecoin-project/specs-actors/v5/actors/runtime/proof"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"
	"math"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
	"github.com/xsleonard/go-merkle"
	"golang.org/x/crypto/sha3"
)

var (
	// MaxCheckpointLength is the maximum number of blocks that can be requested for constructing a checkpoint root hash
	MaxCheckpointLength = uint64(math.Pow(2, 15))
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain         consensus.ChainHeaderReader
	bor           *Bor
	rootHashCache *lru.ARCCache
}

type SendSealArgs struct {
	Signer        common.Address  `json:"signer"`
	SealProofType *hexutil.Uint64 `json:"seal_proof_type"`
	SectorNumber  *hexutil.Uint64 `json:"sector_number"`
	TicketEpoch   *hexutil.Uint64 `json:"ticket_epoch"`
	SeedEpoch     *hexutil.Uint64 `json:"seed_epoch"`
	SealedCID     *hexutil.Bytes  `json:"sealed_cid"`
	UnsealedCID   *hexutil.Bytes  `json:"unsealed_cid"`
	Proof         *hexutil.Bytes  `json:"proof"`
}

func (api *API) VerifySeal(sealData SendSealArgs) (bool, error) {
	var randomness, interactiveRandomness []byte
	fromString, _ := address.NewFromString("t01000")
	buf := new(bytes.Buffer)
	if err := fromString.MarshalCBOR(buf); err != nil {
		return false, err
	}
	ticketBlockHeader := api.chain.GetHeaderByNumber(uint64(*sealData.TicketEpoch))
	if ticketBlockHeader != nil {
		randData := blake2b.Sum256(ticketBlockHeader.Hash().Bytes())
		randomness, _ = sealing2.DrawRandomness(randData[:], ticketBlockHeader.Number.Int64(), buf.Bytes(), sealing2.DomainSeparationTag_SealRandomness)
	}
	seedBLockHeader := api.chain.GetHeaderByNumber(uint64(*sealData.SeedEpoch))
	if seedBLockHeader != nil {
		randData := blake2b.Sum256(seedBLockHeader.Hash().Bytes())
		interactiveRandomness, _ = sealing2.DrawRandomness(randData[:], seedBLockHeader.Number.Int64(), buf.Bytes(), sealing2.DomainSeparationTag_InteractiveSealChallengeSeed)
	}
	CommR, _ := cid.Cast(*sealData.SealedCID)
	CommD, _ := cid.Cast(*sealData.UnsealedCID)
	svi := proof5.SealVerifyInfo{
		SectorID:              abi.SectorID{Miner: abi.ActorID(1000), Number: abi.SectorNumber(*sealData.SectorNumber)},
		SealedCID:             CommR,
		SealProof:             abi.RegisteredSealProof(*sealData.SealProofType),
		Proof:                 *sealData.Proof,
		Randomness:            randomness,
		InteractiveRandomness: interactiveRandomness,
		UnsealedCID:           CommD,
	}
	ok, err := ffiwrapper.ProofVerifier.VerifySeal(svi)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.bor.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetAuthor retrieves the author a block.
func (api *API) GetAuthor(number *rpc.BlockNumber) (*common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	author, err := api.bor.Author(header)
	return &author, err
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.bor.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSigners retrieves the list of authorized signers at the specified block.
func (api *API) GetSigners(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.bor.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.signers(), nil
}

// GetSignersAtHash retrieves the list of authorized signers at the specified block.
func (api *API) GetSignersAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.bor.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.signers(), nil
}

// GetCurrentProposer gets the current proposer
func (api *API) GetCurrentProposer() (common.Address, error) {
	snap, err := api.GetSnapshot(nil)
	if err != nil {
		return common.Address{}, err
	}
	return snap.ValidatorSet.GetProposer().Address, nil
}

// GetCurrentValidators gets the current validators
func (api *API) GetCurrentValidators() ([]*Validator, error) {
	snap, err := api.GetSnapshot(nil)
	if err != nil {
		return make([]*Validator, 0), err
	}
	return snap.ValidatorSet.Validators, nil
}

// GetRootHash returns the merkle root of the start to end block headers
func (api *API) GetRootHash(start uint64, end uint64) (string, error) {
	if err := api.initializeRootHashCache(); err != nil {
		return "", err
	}
	key := getRootHashKey(start, end)
	if root, known := api.rootHashCache.Get(key); known {
		return root.(string), nil
	}
	length := uint64(end - start + 1)
	if length > MaxCheckpointLength {
		return "", &MaxCheckpointLengthExceededError{start, end}
	}
	currentHeaderNumber := api.chain.CurrentHeader().Number.Uint64()
	if start > end || end > currentHeaderNumber {
		return "", &InvalidStartEndBlockError{start, end, currentHeaderNumber}
	}
	blockHeaders := make([]*types.Header, end-start+1)
	wg := new(sync.WaitGroup)
	concurrent := make(chan bool, 20)
	for i := start; i <= end; i++ {
		wg.Add(1)
		concurrent <- true
		go func(number uint64) {
			blockHeaders[number-start] = api.chain.GetHeaderByNumber(uint64(number))
			<-concurrent
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(concurrent)

	headers := make([][32]byte, nextPowerOfTwo(length))
	for i := 0; i < len(blockHeaders); i++ {
		blockHeader := blockHeaders[i]
		header := crypto.Keccak256(appendBytes32(
			blockHeader.Number.Bytes(),
			new(big.Int).SetUint64(blockHeader.Time).Bytes(),
			blockHeader.TxHash.Bytes(),
			blockHeader.ReceiptHash.Bytes(),
		))

		var arr [32]byte
		copy(arr[:], header)
		headers[i] = arr
	}

	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeaves: true})
	if err := tree.Generate(convert(headers), sha3.NewLegacyKeccak256()); err != nil {
		return "", err
	}
	root := hex.EncodeToString(tree.Root().Hash)
	api.rootHashCache.Add(key, root)
	return root, nil
}

func (api *API) initializeRootHashCache() error {
	var err error
	if api.rootHashCache == nil {
		api.rootHashCache, err = lru.NewARC(10)
	}
	return err
}

func getRootHashKey(start uint64, end uint64) string {
	return strconv.FormatUint(start, 10) + "-" + strconv.FormatUint(end, 10)
}
