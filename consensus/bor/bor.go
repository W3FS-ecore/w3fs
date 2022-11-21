package bor

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/go-units"
	"github.com/ethereum/go-ethereum/borcontracts/w3fsStorageManager"
	"github.com/ethereum/go-ethereum/sealing"
	fabi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/lotus/api"

	//"github.com/ethereum/go-ethereum/extern/sector-storage/ffiwrapper"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	"github.com/minio/blake2b-simd"
	"io"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
	miningLogAtDepth = 7      // miningLogAtDepth is the number of confirmations before logging successful mining.

)

// Bor protocol constants.
var (
	FrontierBlockReward       = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
	ByzantiumBlockReward      = big.NewInt(3e+18) // Block reward in wei for successfully mining a block upward from Byzantium
	ConstantinopleBlockReward = big.NewInt(2e+18) // Block reward in wei for successfully mining a block upward from Constantinople
	big8                      = big.NewInt(8)
	big32                     = big.NewInt(32)

	defaultSprintLength = uint64(64) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.

	diffInTurn = big.NewInt(2) // Block difficulty for in-turn signatures
	diffNoTurn = big.NewInt(1) // Block difficulty for out-of-turn signatures

	validatorHeaderBytesLength = common.AddressLength + 20 // address + power
	systemAddress              = common.HexToAddress("0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE")
	emptyAddress			   = common.Address{}
	validatorFileCoinContract = "0x0000000000000000000000000000000000001002"
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errInvalidCheckpointBeneficiary is returned if a checkpoint/epoch transition
	// block has a beneficiary set to non-zeroes.
	errInvalidCheckpointBeneficiary = errors.New("beneficiary in checkpoint block non-zero")

	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")

	// errInvalidCheckpointVote is returned if a checkpoint/epoch transition block
	// has a vote nonce set to non-zeroes.
	errInvalidCheckpointVote = errors.New("vote nonce in checkpoint block non-zero")

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	// errExtraValidators is returned if non-sprint-end block contain validator data in
	// their extra-data fields.
	errExtraValidators = errors.New("non-sprint-end block contains extra validator list")

	// errInvalidSpanValidators is returned if a block contains an
	// invalid list of validators (i.e. non divisible by 40 bytes).
	errInvalidSpanValidators = errors.New("invalid validator list on sprint end block")

	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")

	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")

	// errInvalidDifficulty is returned if the difficulty of a block neither 1 or 2.
	errInvalidDifficulty = errors.New("invalid difficulty")

	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// errOutOfRangeChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errOutOfRangeChain = errors.New("out of range or non-contiguous chain")

	errProof = errors.New("VerifyWinningPoSt is error")

	errWinningExtraDecode = errors.New("Decode header.WinningExtraData is error")
)

// SignerFn is a signer callback function to request a header to be signed by a
// backing account.
type SignerFn func(accounts.Account, string, []byte) ([]byte, error)
type SignerTxFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error)

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, signer)
	return signer, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.ExtraPower,
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

// CalcProducerDelay is the block delay algorithm based on block time, period, producerDelay and turn-ness of a signer
func CalcProducerDelay(number uint64, succession int, c *params.BorConfig) uint64 {
	// When the block is the first block of the sprint, it is expected to be delayed by `producerDelay`.
	// That is to allow time for block propagation in the last sprint
	delay := c.Period
	if number%c.Sprint == 0 {
		delay = c.ProducerDelay
	}
	if succession > 0 {
		delay += uint64(succession) * c.BackupMultiplier
	}
	return delay
}

// BorRLP returns the rlp bytes which needs to be signed for the bor
// sealing. The RLP to sign consists of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func BorRLP(header *types.Header) []byte {
	b := new(bytes.Buffer)
	encodeSigHeader(b, header)
	return b.Bytes()
}

// Bor is the matic-bor consensus engine
type Bor struct {
	chainConfig *params.ChainConfig // Chain config
	config      *params.BorConfig   // Consensus engine configuration parameters for bor consensus
	db          ethdb.Database      // Database to store and retrieve snapshot checkpoints

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	signer common.Address // Ethereum address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with
	lock   sync.RWMutex   // Protects the signer fields

	ethAPI                 *ethapi.PublicBlockChainAPI
	GenesisW3fsContractsClient *GenesisW3fsContractsClient
	GenesisContractsClient *GenesisContractsClient
	validatorSetABI        abi.ABI
	stateReceiverABI       abi.ABI
	HeimdallClient         IHeimdallClient
	WithoutHeimdall        bool

	scope event.SubscriptionScope
	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications

	//schedule drand.Schedule
	dataDir string
	storageApi *api.StorageMiner
	txSigner types.Signer
	signTxFn SignerTxFn
}


// New creates a Matic Bor consensus engine.
func New(
	chainConfig *params.ChainConfig,
	db ethdb.Database,
	ethAPI *ethapi.PublicBlockChainAPI,
	heimdallURL string,
	withoutHeimdall bool,
	dataDir string,
) *Bor {
	// get bor config
	borConfig := chainConfig.Bor

	//schedule, _ := drand.RandomMockSchedule()
	//schedule, _ := drand.RandomDrandSchedule(genesisTs)

	// Set any missing consensus parameters to their defaults
	if borConfig != nil && borConfig.Sprint == 0 {
		borConfig.Sprint = defaultSprintLength
	}

	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)
	vABI, _ := abi.JSON(strings.NewReader(validatorsetABI))
	sABI, _ := abi.JSON(strings.NewReader(stateReceiverABI))
	heimdallClient, _ := NewHeimdallClient(heimdallURL)
	genesisContractsClient := NewGenesisContractsClient(chainConfig, borConfig.ValidatorContract, borConfig.StateReceiverContract, ethAPI)
	genesisW3fsContractsClient := NewGenesisW3fsContractsClient(chainConfig, ethAPI)
	c := &Bor{
		chainConfig:            chainConfig,
		config:                 borConfig,
		db:                     db,
		ethAPI:                 ethAPI,
		recents:                recents,
		signatures:             signatures,
		validatorSetABI:        vABI,
		stateReceiverABI:       sABI,
		GenesisContractsClient: genesisContractsClient,
		GenesisW3fsContractsClient: genesisW3fsContractsClient,
		HeimdallClient:         heimdallClient,
		WithoutHeimdall:        withoutHeimdall,
		dataDir: 				dataDir,
		txSigner:               types.NewEIP155Signer(chainConfig.ChainID),
	}
	c.GenesisW3fsContractsClient.Bor = c
	return c
}

func (c *Bor) SetStorageApi(storageApi *api.StorageMiner) {
	c.storageApi = storageApi
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (c *Bor) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, c.signatures)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (c *Bor) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return c.verifyHeader(chain, header, nil)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (c *Bor) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := c.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (c *Bor) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// Don't waste time checking blocks from the future
	if header.Time > uint64(time.Now().Unix()) {
		return consensus.ErrFutureBlock
	}

	if err := validateHeaderExtraField(header.Extra); err != nil {
		return err
	}

	// check extr adata
	isSprintEnd := (number+1)%c.config.Sprint == 0

	// Ensure that the extra-data contains a signer list on checkpoint, but none otherwise
	signersBytes := len(header.Extra) - extraVanity - extraSeal
	signersPowerBytes := len(header.ExtraPower)
	if !isSprintEnd && signersBytes != 0 {
		return errExtraValidators
	}

	if !isSprintEnd && signersPowerBytes != 0 {
		return errExtraValidators
	}

	if isSprintEnd && signersBytes%validatorHeaderBytesLength != 0 {
		return errInvalidSpanValidators
	}

	if isSprintEnd && signersPowerBytes%validatorHeaderBytesLength != 0 {
		return errInvalidSpanValidators
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if header.Difficulty == nil {
			return errInvalidDifficulty
		}
	}
	// If all checks passed, validate any special fields for hard forks
	if err := misc.VerifyForkHashes(chain.Config(), header, false); err != nil {
		return err
	}
	// All basic checks passed, verify cascading fields
	return c.verifyCascadingFields(chain, header, parents)
}


func (c *Bor) verifyProof(chain consensus.ChainHeaderReader, header *types.Header, snap *Snapshot) error {
	number := header.Number.Uint64()
	signer, err := ecrecover(header, c.signatures)
	if err != nil {
		return err
	}
	if number%c.config.Sprint == 0 {
		if err != nil {
			return err
		}
		if snap.ValidatorSealPowerSet.hasTotalPower() && snap.ValidatorSealPowerSet.hasAddressAndPower(signer.Bytes()) {
			if len(header.WinningExtraData) == 0 {
				errNoWinning := &ErrorNoWinningData{Number: number,TotalPower: snap.ValidatorSealPowerSet.getTotalPower(), SignerPower: snap.ValidatorSealPowerSet.getSignerPower(signer.Bytes())}
				log.Error(errNoWinning.Error())
				return errNoWinning
			}
			var sealType uint64
			//cvotes, _ := c.getValidatorVotes(signer, header.ParentHash)
			nelaStorageManagerSectors, _ := c.getSealInfoAllBySigner(signer, header.ParentHash)
			if len(nelaStorageManagerSectors) > 0 {
				sealType = nelaStorageManagerSectors[0].SealProofType.Uint64()
			}
			randData := c.generateRandData(header)
			indexs, _ := c.getWinningPoStSectorChallengeVote(signer, header.Number, randData, uint64(len(nelaStorageManagerSectors)), sealType)
			var docvotes []w3fsStorageManager.Cvote
			for _, index := range indexs {
				docvotes = append(docvotes, w3fsStorageManager.Cvote{
					SectorInx:     nelaStorageManagerSectors[index].SectorNumber.Uint64(),
					SealProofType: nelaStorageManagerSectors[index].SealProofType.Uint64(),
					SealedCID:     nelaStorageManagerSectors[index].SealedCID,
				})
			}
			var winningData w3fsStorageManager.WinningData
			if err = rlp.DecodeBytes(header.WinningExtraData, &winningData); err != nil {
				return errWinningExtraDecode
			}
			mid := c.changeMidByEthAddress(signer)
			mid = 1000
			poStProofRand, _ := sealing.GetTicket(header.Number, randData, fmt.Sprintf("%s%s", "t0", strconv.Itoa(int(mid))), sealing.DomainSeparationTag_WinningPoStChallengeSeed)
			ok, verifyErr := sealing.VerifyWinningPoSt(mid, docvotes, winningData.WinningPostProofs, poStProofRand)
			if verifyErr != nil || !ok {
				errCheckProof := &ErrorCheckProof{Number: number, Signer: signer}
				log.Error(errCheckProof.Error())
				return errCheckProof
			} else {
				log.Info("VerifyWinningPoSt successful !", "signer", signer.String(), "number", number, "indexs", indexs)
			}
		}
	}
	return nil
}



// validateHeaderExtraField validates that the extra-data contains both the vanity and signature.
// header.Extra = header.Vanity + header.ProducerBytes (optional) + header.Seal
func validateHeaderExtraField(extraBytes []byte) error {
	if len(extraBytes) < extraVanity {
		return errMissingVanity
	}
	if len(extraBytes) < extraVanity+extraSeal {
		return errMissingSignature
	}
	return nil
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (c *Bor) verifyCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}

	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	if parent.Time + c.config.Period > header.Time {
		return ErrInvalidTimestamp
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := c.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	//err = drand.ValidateBlockValues(int64(number), c.schedule, []drand.BeaconEntry{{Round: header.Round, Data: header.BeaconData}}, drand.BeaconEntry{
	//			Round: parent.Round,
	//			Data: parent.BeaconData,
	//		})
	//if err != nil {
	//	log.Error("ValidateBlockValues is wrong")
	//}

	// verify the validator list in the last sprint block
	if isSprintStart(number, c.config.Sprint) {
		parentValidatorBytes := parent.Extra[extraVanity : len(parent.Extra)-extraSeal]
		validatorsBytes := make([]byte, len(snap.ValidatorSet.Validators)*validatorHeaderBytesLength)
		currentValidators := snap.ValidatorSet.Copy().Validators
		// sort validator by address
		sort.Sort(ValidatorsByAddress(currentValidators))
		for i, validator := range currentValidators {
			copy(validatorsBytes[i*validatorHeaderBytesLength:], validator.HeaderBytes())
		}
		// len(header.Extra) >= extraVanity+extraSeal has already been validated in validateHeaderExtraField, so this won't result in a panic
		if !bytes.Equal(parentValidatorBytes, validatorsBytes) {
			return &MismatchingValidatorsError{number - 1, validatorsBytes, parentValidatorBytes}
		}

		parentValidatorPowerBytes := parent.ExtraPower[:]
		validatorsPowerBytes := make([]byte, len(snap.ValidatorSealPowerSet.ValidatorSealPowers)*validatorHeaderBytesLength)
		currentValidatorSealPowers := snap.ValidatorSealPowerSet.Copy().ValidatorSealPowers
		sort.Sort(ValidatorSealPowerByAddress(currentValidatorSealPowers))
		for i, validatorPower := range currentValidatorSealPowers {
			copy(validatorsPowerBytes[i*validatorHeaderBytesLength:], validatorPower.HeaderBytes())
		}
		if !bytes.Equal(parentValidatorPowerBytes, validatorsPowerBytes) {
			return &MismatchingValidatorsPowerError{number - 1, validatorsPowerBytes, parentValidatorPowerBytes}
		}
	}

	// All basic checks passed, verify the seal and return
	return c.verifySeal(chain, header, parents)
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (c *Bor) snapshot(chain consensus.ChainHeaderReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)

	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := c.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}

		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(c.config, c.signatures, c.db, hash, c.ethAPI); err == nil {
				log.Trace("Loaded snapshot from disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}

		// If we're at the genesis, snapshot the initial state. Alternatively if we're
		// at a checkpoint block without a parent (light client CHT), or we have piled
		// up more headers than allowed to be reorged (chain reinit from a freezer),
		// consider the checkpoint trusted and snapshot it.
		// TODO fix this
		if number == 0 {
			checkpoint := chain.GetHeaderByNumber(number)
			if checkpoint != nil {
				// get checkpoint data
				hash := checkpoint.Hash()

				// get validators and current span
				//validators, validatorsPowers, err := c.GetCurrentValidators(hash, number+1, true)
				validators, validatorsPowers, err := c.GenesisW3fsContractsClient.GetBorMiners(hash, true)

				if err != nil {
					return nil, err
				}

				// new snap shot
				snap = newSnapshot(c.config, c.signatures, number, hash, validators, validatorsPowers, c.ethAPI)
				if err := snap.store(c.db); err != nil {
					return nil, err
				}
				log.Info("Stored checkpoint snapshot to disk", "number", number, "hash", hash)
				break
			}
		}

		// No snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			header = parents[len(parents)-1]
			if header.Hash() != hash || header.Number.Uint64() != number {
				return nil, consensus.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			header = chain.GetHeader(hash, number)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}
		headers = append(headers, header)
		number, hash = number-1, header.ParentHash
	}

	// check if snapshot is nil
	if snap == nil {
		return nil, fmt.Errorf("Unknown error while retrieving snapshot at block number %v", number)
	}

	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}

	snap, err := snap.apply(headers)
	if err != nil {
		return nil, err
	}
	c.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = snap.store(c.db); err != nil {
			return nil, err
		}
		log.Trace("Stored snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (c *Bor) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (c *Bor) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	return c.verifySeal(chain, header, nil)
}

// verifySeal checks whether the signature contained in the header satisfies the
// consensus protocol requirements. The method accepts an optional list of parent
// headers that aren't yet part of the local blockchain to generate the snapshots
// from.
func (c *Bor) verifySeal(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := c.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// Resolve the authorization key and check against signers
	signer, err := ecrecover(header, c.signatures)
	if err != nil {
		return err
	}

	if !snap.ValidatorSet.HasAddress(signer.Bytes()) {
		// Check the UnauthorizedSignerError.Error() msg to see why we pass number-1
		return &UnauthorizedSignerError{number - 1, signer.Bytes()}
	}

	succession, err := snap.GetSignerSuccessionNumber(signer)
	if err != nil {
		return err
	}

	var parent *types.Header
	if len(parents) > 0 { // if parents is nil, len(parents) is zero
		parent = parents[len(parents)-1]
	} else if number > 0 {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}

	verifyTime := parent.Time + CalcProducerDelay(number, succession, c.config)
	if succession > 0 && isSprintStart(number,c.config.Sprint) { // if no Proposer , must delay porepDelay time
		verifyTime += c.config.PorepDelay
	}
	if parent != nil && header.Time < verifyTime {
		return &BlockTooSoonError{number, succession}
	}

	// check winningData if empty
	if isSprintStart(number, c.config.Sprint) &&
		snap.ValidatorSealPowerSet.hasTotalPower() &&
		snap.ValidatorSealPowerSet.hasAddressAndPower(signer.Bytes()) &&
		header.WinningExtraData == nil {
		return &ErrorNoWinningData{
			Number: number,
			TotalPower: snap.ValidatorSealPowerSet.getTotalPower(),
			SignerPower: snap.ValidatorSealPowerSet.getSignerPower(signer.Bytes()),
		}
	}

	// Ensure that the difficulty corresponds to the turn-ness of the signer
	if !c.fakeDiff {
		difficulty := snap.Difficulty(signer)
		if header.Difficulty.Uint64() != difficulty {
			return &WrongDifficultyError{number, difficulty, header.Difficulty.Uint64(), signer.Bytes()}
		}
	}
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *Bor) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// If the block isn't a checkpoint, cast a random vote (good enough for now)
	header.Coinbase = c.signer
	header.Nonce = types.BlockNonce{}
	number := header.Number.Uint64()

	// Assemble the validator snapshot to check which votes make sense
	snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	// Set the correct difficulty
	header.Difficulty = new(big.Int).SetUint64(snap.Difficulty(c.signer))

	// Ensure the extra data has all it's components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	// get validator set if number
	if (number + 1) % c.config.Sprint == 0 {
		//newValidators, newValidatorsPowers, err := c.GetCurrentValidators(header.ParentHash, number+1, true)
		standard := chain.GetHeaderByNumber(number - miningLogAtDepth)
		newValidators, newValidatorsPowers, err := c.GenesisW3fsContractsClient.GetBorMiners(standard.Hash(), true)
		if err != nil {
			return errors.New("unknown validators")
		}
		// sort validator by address
		sort.Sort(ValidatorsByAddress(newValidators))
		for _, validator := range newValidators {
			header.Extra = append(header.Extra, validator.HeaderBytes()...)
		}

		sort.Sort(ValidatorSealPowerByAddress(newValidatorsPowers))
		for _, validatorPower := range newValidatorsPowers {
			header.ExtraPower = append(header.ExtraPower, validatorPower.HeaderBytes()...)
		}
	}

	// add extra seal space
	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	// Ensure the timestamp has the correct delay
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// ====================================== show delay =========================================
	c.printlnValidatorDelay(snap)
	// ========================================================================================

	var succession int
	if bytes.Compare(c.signer.Bytes(), common.Address{}.Bytes()) != 0 {
		succession, err = snap.GetSignerSuccessionNumber(c.signer)
		if err != nil {
			return err
		}
	}
	header.Time = parent.Time + CalcProducerDelay(number, succession, c.config)
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}
	if number % c.config.Sprint == 0 && succession > 0 { // if no Proposer, must delay porepDelay time
		header.Time += c.config.PorepDelay
	}

	{
		if number > 0 && isSprintStart(number, c.config.Sprint) {
			if snap.ValidatorSealPowerSet.hasTotalPower() && snap.ValidatorSealPowerSet.hasAddressAndPower(c.signer.Bytes()) {
				winningExtraData, err := c.winningPostVote(header)
				if err == nil {
					header.WinningExtraData = winningExtraData
				} else {
					return err
				}
			}
		}
	}
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given.
func (c *Bor) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs *[]*types.Transaction, uncles []*types.Header, receipts *[]*types.Receipt, systemTxs *[]*types.Transaction, usedGas *uint64) error {
	cx := chainContext{Chain: chain, Bor: c}
	stateSyncData := []*types.StateSyncData{}
	var err error
	headerNumber := header.Number.Uint64()
	// The verification can only be done when the state is ready, it can't be done in VerifyHeader.
	snap, err := c.snapshot(chain, headerNumber - 1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	if proofErr := c.verifyProof(chain, header, snap) ; proofErr != nil {
		return proofErr
	}
	if err = c.GenesisW3fsContractsClient.receiverReward(state, header, cx); err != nil {
		log.Error("Error while committing receiverReward", "error", err)
		return err
	} else {
		log.Info("→ Finalize committing new receiverReward", "coinbase", header.Coinbase)
	}

	if doSlash := c.GenesisW3fsContractsClient.checkDoSlash(snap, headerNumber, c.config.Sprint, header.Coinbase); doSlash {
		proposer := snap.ValidatorSet.GetProposer().Address
		err := c.GenesisW3fsContractsClient.slash(proposer, state, header, cx, txs, receipts, systemTxs, usedGas, false)
		if err != nil {
			log.Error("slash wrong" , "error" , err)
		}
	}


	if headerNumber % c.config.Sprint == 0 {
		// check and commit span
		if err := c.checkAndCommitSpan(state, header, cx); err != nil {
			log.Error("Error while committing span", "error", err)
			return nil
		}
		if !c.WithoutHeimdall {
			// commit states
			stateSyncData, err = c.CommitStates(state, header, cx)
			if err != nil {
				log.Error("Error while committing states", "error", err)
				return nil
			}
		}
	}

	// No block rewards in PoA, so the state remains as is and uncles are dropped
	//accumulateRewards(chain.Config(), state, header, uncles)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	// Set state sync data to blockchain
	bc := chain.(*core.BlockChain)
	bc.SetStateSync(stateSyncData)
	return nil
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (c *Bor) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, []*types.Receipt, error) {
	stateSyncData := []*types.StateSyncData{}
	cx := chainContext{Chain: chain, Bor: c}
	headerNumber := header.Number.Uint64()
	log.Info("FinalizeAndAssemble -- headerNumber", "number", headerNumber, "c.config.Sprint", c.config.Sprint, "WithoutHeimdall", c.WithoutHeimdall)

	rewardErr := c.GenesisW3fsContractsClient.receiverReward(state, header, cx)
	if rewardErr != nil {
		log.Error("Error while committing receiverReward", "error", rewardErr)
		return nil, nil, rewardErr
	} else {
		log.Info("→ FinalizeAndAssemble committing new receiverReward", "coinbase", header.Coinbase)
	}

	snap, err := c.snapshot(chain, headerNumber-1, header.ParentHash, nil)
	if err != nil {
		return nil, nil, err
	}
	if doSlash := c.GenesisW3fsContractsClient.checkDoSlash(snap, headerNumber, c.config.Sprint, c.signer); doSlash && c.signTxFn != nil {
		proposer := snap.ValidatorSet.GetProposer().Address
		err := c.GenesisW3fsContractsClient.slash(proposer, state, header, cx, &txs, &receipts, nil, &header.GasUsed, true)
		if err != nil {
			return nil, nil, err
		}
	}

	if headerNumber % c.config.Sprint == 0 {
		// check and commit span
		err := c.checkAndCommitSpan(state, header, cx)
		if err != nil {
			log.Error("Error while committing span", "error", err)
			return nil, nil, err
		}
		if !c.WithoutHeimdall {
			log.Info("start CommitStates")
			stateSyncData, err = c.CommitStates(state, header, cx)
			if err != nil {
				log.Error("Error while committing states", "error", err)
				return nil, nil, err
			}
		}
	}
	// No block rewards in PoA, so the state remains as is and uncles are dropped
	//accumulateRewards(chain.Config(), state, header, uncles)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	// Assemble block
	block := types.NewBlock(header, txs, nil, receipts, new(trie.Trie))
	// set state sync
	bc := chain.(*core.BlockChain)
	bc.SetStateSync(stateSyncData)
	// return the final block for sealing
	return block, receipts, nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (c *Bor) Authorize(signer common.Address, signFn SignerFn, signTxFn SignerTxFn) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.signer = signer
	c.signFn = signFn
	c.signTxFn = signTxFn
}

func (c *Bor) winningPostVote(header *types.Header) ([]byte, error) {
	log.Info("GenerateWinningPoSt start","time", time.Now().Unix())
	//cvotes, err := c.getValidatorVotes(c.signer, header.ParentHash)
	nelaStorageManagerSector, err := c.getSealInfoAllBySigner(c.signer, header.ParentHash)
	if err != nil {
		log.Error("get validatorVotes is wrong", "error" , err)
		return nil,err
	}
	var sealType uint64
	if len(nelaStorageManagerSector) > 0 {
		sealType = nelaStorageManagerSector[0].SealProofType.Uint64()
	}
	randData := c.generateRandData(header)
	indexs, err := c.getWinningPoStSectorChallengeVote(c.signer, header.Number, randData, uint64(len(nelaStorageManagerSector)), sealType)
	log.Info("getWinningPoStSectorChallengeVote success ! " , "indexs" , indexs)
	if err != nil {
		log.Error("getWinningPoStSectorChallengeVote is wrong" , "error", err)
		return nil, err
	}
	var docVotes []w3fsStorageManager.Cvote
	for _, index := range indexs {
		docVotes = append(docVotes, w3fsStorageManager.Cvote{
			SectorInx:     nelaStorageManagerSector[index].SectorNumber.Uint64(),
			SealProofType: nelaStorageManagerSector[index].SealProofType.Uint64(),
			SealedCID:     nelaStorageManagerSector[index].SealedCID,
		})
	}
	mid := c.changeMidByEthAddress(c.signer)
	mid = 1000
	poStProofRand, _ := sealing.GetTicket(header.Number, randData, fmt.Sprintf("%s%s", "t0", strconv.Itoa(int(mid))), sealing.DomainSeparationTag_WinningPoStChallengeSeed)
	poStProofs, err := sealing.GenerateWinningPoSt2(docVotes, poStProofRand, c.storageApi)
	if err != nil {
		return nil ,err
	}
	var winningPostProofs []w3fsStorageManager.WinningPostProof
	for _, proof := range poStProofs {
		winningPostProofs = append(winningPostProofs, w3fsStorageManager.WinningPostProof{
			PoStProof:  uint64(proof.PoStProof),
			ProofBytes: proof.ProofBytes,
		})
	}
	winningExtraData, err := rlp.EncodeToBytes(w3fsStorageManager.WinningData{
		WinningPostProofs: winningPostProofs,
	})
	if err != nil {
		return nil ,err
	}
	log.Info("GenerateWinningPoSt end","time", time.Now().Unix(), "indexs" , indexs)
	return winningExtraData, nil
}


// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (c *Bor) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	header := block.Header()
	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// For 0-period chains, refuse to seal empty blocks (no reward but would spin sealing)
	if c.config.Period == 0 && len(block.Transactions()) == 0 {
		log.Info("Sealing paused, waiting for transactions")
		return nil
	}
	// Don't hold the signer fields for the entire sealing procedure
	c.lock.RLock()
	signer, signFn := c.signer, c.signFn
	c.lock.RUnlock()

	snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	// Bail out if we're unauthorized to sign a block
	if !snap.ValidatorSet.HasAddress(signer.Bytes()) {
		// Check the UnauthorizedSignerError.Error() msg to see why we pass number-1
		return &UnauthorizedSignerError{number - 1, signer.Bytes()}
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(int64(header.Time), 0).Sub(time.Now())

	// Sign all the things!
	sighash, err := signFn(accounts.Account{Address: signer}, accounts.MimetypeBor, BorRLP(header))
	if err != nil {
		return err
	}
	copy(header.Extra[len(header.Extra)-extraSeal:], sighash)

	// Wait until sealing is terminated or delay timeout.
	log.Info("Waiting for slot to sign and propagate", "delay", common.PrettyDuration(delay), "number" , number)
	go func() {
		select {
		case <-stop:
			log.Debug("Discarding sealing operation for block", "number", number)
			return
		case <-time.After(delay):
			log.Info("Sealing successful", "number", number, "delay", delay, "headerDifficulty", header.Difficulty)
		}
		select {
		case results <- block.WithSeal(header):
		default:
			log.Warn("Sealing result was not read by miner", "number", number, "sealhash", SealHash(header))
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (c *Bor) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	snap, err := c.snapshot(chain, parent.Number.Uint64(), parent.Hash(), nil)
	if err != nil {
		return nil
	}
	return new(big.Int).SetUint64(snap.Difficulty(c.signer))
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *Bor) SealHash(header *types.Header) common.Hash {
	return SealHash(header)
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (c *Bor) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "bor",
		Version:   "1.0",
		Service:   &API{chain: chain, bor: c},
		Public:    false,
	}}
}

// Close implements consensus.Engine. It's a noop for bor as there are no background threads.
func (c *Bor) Close() error {
	return nil
}

// GetCurrentSpan get current span from contract
func (c *Bor) GetCurrentSpan(headerHash common.Hash) (*Span, error) {
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)

	// method
	method := "getCurrentSpan"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err := c.validatorSetABI.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for getCurrentSpan", "error", err)
		return nil, err
	}

	msgData := (hexutil.Bytes)(data)
	// c.config.ValidatorContract = 0x0000000000000000000000000000000000001000
	toAddress := common.HexToAddress(c.config.ValidatorContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		return nil, err
	}

	// span result
	ret := new(struct {
		Number     *big.Int
		StartBlock *big.Int
		EndBlock   *big.Int
	})
	if err := c.validatorSetABI.UnpackIntoInterface(ret, method, result); err != nil {
		return nil, err
	}

	// create new span
	span := Span{
		ID:         ret.Number.Uint64(),
		StartBlock: ret.StartBlock.Uint64(),
		EndBlock:   ret.EndBlock.Uint64(),
	}

	return &span, nil
}



// GetCurrentValidators get current validators
/*func (c *Bor) GetCurrentValidators(headerHash common.Hash, blockNumber uint64, isAddPower bool) ([]*Validator, []*ValidatorSealPower, error) {
	// block
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)

	// method
	method := "getBorValidators"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err := c.validatorSetABI.Pack(method, big.NewInt(0).SetUint64(blockNumber))
	if err != nil {
		log.Error("Unable to pack tx for getValidator", "error", err)
		return nil, nil, err
	}

	// call
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(c.config.ValidatorContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		panic(err)
	}

	var (
		ret0 = new([]common.Address)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}

	if err := c.validatorSetABI.UnpackIntoInterface(out, method, result); err != nil {
		return nil, nil, err
	}

	valz := make([]*Validator, len(*ret0))
	for i, a := range *ret0 {
		valz[i] = &Validator{
			Address:     a,
			VotingPower: (*ret1)[i].Int64(),
		}
	}
	var valp []*ValidatorSealPower

	if isAddPower {
		var addrs []common.Address
		for _, v := range valz {
			addrs = append(addrs, v.Address)
		}
		validatorSealPower, err := c.getAllValidatorPower(addrs, headerHash)
		if err == nil {
			valp = validatorSealPower
			for _, vs := range validatorSealPower {
				log.Info(fmt.Sprintf("address =  %s , power = %d", hexutil.Encode(vs.Address.Bytes()), vs.SealPower))
			}
		}
	}
	return valz, valp, nil
}*/

/*func (c *Bor) getValidatorSectorInx(addr common.Address, headerHash common.Hash) (uint64, error){
	if bytes.Compare(addr.Bytes(), common.Address{}.Bytes()) == 0 {
		// TODO ERROR
	}
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)
	// method
	method := "getValidatorSectorInx"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := borcontracts.FABI.Pack(method, addr)
	if err != nil {
		return 0, err
	}
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(validatorFileCoinContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		 return 0, err
	}
	var (
		ret0 = new(big.Int)
	)
	out := &[]interface{}{
		ret0,
	}
	if err := borcontracts.FABI.UnpackIntoInterface(out, method, result); err != nil {
		return 0, err
	}
	return ret0.Uint64(), nil
}*/


func (c *Bor) checkAndCommitSpan(
	state *state.StateDB,
	header *types.Header,
	chain core.ChainContext,
) error {
	headerNumber := header.Number.Uint64()
	span, err := c.GetCurrentSpan(header.ParentHash)
	if err != nil {
		return err
	}
	if c.needToCommitSpan(span, headerNumber) {
		err := c.fetchAndCommitSpan(span.ID+1, state, header, chain)
		return err
	}
	return nil
}



func (c *Bor) needToCommitSpan(span *Span, headerNumber uint64) bool {
	// if span is nil
	if span == nil {
		return false
	}

	// check span is not set initially
	if span.EndBlock == 0 {
		return true
	}

	// if current block is first block of last sprint in current span
	if span.EndBlock > c.config.Sprint && span.EndBlock-c.config.Sprint+1 == headerNumber {
		return true
	}

	return false
}


func (c *Bor) fetchAndCommitSpan(
	newSpanID uint64,
	state *state.StateDB,
	header *types.Header,
	chain core.ChainContext,
) error {
	var heimdallSpan HeimdallSpan

	if c.WithoutHeimdall {

		s, err := c.getNextHeimdallSpanForTest(newSpanID, state, header, chain)
		if err != nil {
			return err
		}
		heimdallSpan = *s
	} else {
		response, err := c.HeimdallClient.FetchWithRetry(fmt.Sprintf("bor/span/%d", newSpanID), "")
		if err != nil {
			return err
		}

		if err := json.Unmarshal(response.Result, &heimdallSpan); err != nil {
			return err
		}
	}

	// check if chain id matches with heimdall span
	if heimdallSpan.ChainID != c.chainConfig.ChainID.String() {
		return fmt.Errorf(
			"Chain id proposed span, %s, and bor chain id, %s, doesn't match",
			heimdallSpan.ChainID,
			c.chainConfig.ChainID,
		)
	}
	// get validators bytes
	var validators []MinimalVal
	for _, val := range heimdallSpan.ValidatorSet.Validators {
		validators = append(validators, val.MinimalVal())
	}
	validatorBytes, err := rlp.EncodeToBytes(validators)
	if err != nil {
		return err
	}

	// get producers bytes
	var producers []MinimalVal
	for _, val := range heimdallSpan.SelectedProducers {
		producers = append(producers, val.MinimalVal())
	}
	producerBytes, err := rlp.EncodeToBytes(producers)
	if err != nil {
		return err
	}
	method := "commitSpan"
	log.Info("✅ Committing new span", "id", heimdallSpan.ID, "startBlock", heimdallSpan.StartBlock, "endBlock", heimdallSpan.EndBlock, "validatorBytes", hex.EncodeToString(validatorBytes), "producerBytes", hex.EncodeToString(producerBytes), "accountRootHashBytes" , common.HexToHash(heimdallSpan.AccountRootHash))
	// get packed data
	data, err := c.validatorSetABI.Pack(
		method,
		big.NewInt(0).SetUint64(heimdallSpan.ID),
		big.NewInt(0).SetUint64(heimdallSpan.StartBlock),
		big.NewInt(0).SetUint64(heimdallSpan.EndBlock),
		validatorBytes,
		producerBytes,
		common.HexToHash(heimdallSpan.AccountRootHash),
	)
	if err != nil {
		log.Error("Unable to pack tx for commitSpan", "error", err)
		return err
	}
	// get system message
	msg := getSystemMessage(common.HexToAddress(c.config.ValidatorContract), data)
	err = applyMessage(msg, state, header, c.chainConfig, chain)
	if err != nil {
		return err
	}

	// =============== update storage promise
	/*
		var validatorPromise []MiniPromiseVal
		for _, val := range heimdallSpan.ValidatorSet.Validators {
			validatorPromise = append(validatorPromise, val.MiniPromiseVal())
		}
		validatorPromiseBytes, err := rlp.EncodeToBytes(validatorPromise)
		if err != nil {
			return err
		}
		methodP := "updateValidatorPromiseBySystem"
		log.Info("✅ Committing updateValidatorPromiseBySystem", "validatorPromiseBytes", hex.EncodeToString(validatorPromiseBytes))
		dataP, err :=  borcontracts.FABI.Pack(methodP, validatorPromiseBytes)
		if err != nil {
			log.Error("Unable to pack tx for commitSpan", "error", err)
			return err
		}
		// get system message
		msgP := getSystemMessage(common.HexToAddress(validatorFileCoinContract), dataP)
		err = applyMessage(msgP, state, header, c.chainConfig, chain)
		if err != nil {
			return err
		}
	*/
	return nil
}

// CommitStates commit states
func (c *Bor) CommitStates(
	state *state.StateDB,
	header *types.Header,
	chain chainContext,
) ([]*types.StateSyncData, error) {
	stateSyncs := make([]*types.StateSyncData, 0)
	number := header.Number.Uint64()
	_lastStateID, err := c.GenesisContractsClient.LastStateId(number - 1)
	log.Info("print _lastStateID", "_lastStateID", _lastStateID)
	if err != nil {
		return nil, err
	}

	to := time.Unix(int64(chain.Chain.GetHeaderByNumber(number-c.config.Sprint).Time), 0)
	lastStateID := _lastStateID.Uint64()
	log.Info(
		"Fetching state updates from Heimdall",
		"fromID", lastStateID+1,
		"to", to.Format(time.RFC3339))
	eventRecords, err := c.HeimdallClient.FetchStateSyncEvents(lastStateID+1, to.Unix())
	if c.config.OverrideStateSyncRecords != nil {
		if val, ok := c.config.OverrideStateSyncRecords[strconv.FormatUint(number, 10)]; ok {
			eventRecords = eventRecords[0:val]
		}
	}

	chainID := c.chainConfig.ChainID.String()
	for _, eventRecord := range eventRecords {
		if eventRecord.ID <= lastStateID {
			continue
		}
		if err := validateEventRecord(eventRecord, number, to, lastStateID, chainID); err != nil {
			log.Error(err.Error())
			break
		}

		stateData := types.StateSyncData{
			ID:       eventRecord.ID,
			Contract: eventRecord.Contract,
			Data:     hex.EncodeToString(eventRecord.Data),
			TxHash:   eventRecord.TxHash,
		}
		stateSyncs = append(stateSyncs, &stateData)

		if err := c.GenesisContractsClient.CommitState(eventRecord, state, header, chain); err != nil {
			return nil, err
		}
		lastStateID++
	}
	return stateSyncs, nil
}

func validateEventRecord(eventRecord *EventRecordWithTime, number uint64, to time.Time, lastStateID uint64, chainID string) error {
	// event id should be sequential and event.Time should lie in the range [from, to)
	if lastStateID+1 != eventRecord.ID || eventRecord.ChainID != chainID || !eventRecord.Time.Before(to) {
		return &InvalidStateReceivedError{number, lastStateID, &to, eventRecord}
	}
	return nil
}

func (c *Bor) SetHeimdallClient(h IHeimdallClient) {
	c.HeimdallClient = h
}

//
// Private methods
//

func (c *Bor) getNextHeimdallSpanForTest(
	newSpanID uint64,
	state *state.StateDB,
	header *types.Header,
	chain core.ChainContext,
) (*HeimdallSpan, error) {
	headerNumber := header.Number.Uint64()
	span, err := c.GetCurrentSpan(header.ParentHash)
	if err != nil {
		return nil, err
	}

	// get local chain context object
	localContext := chain.(chainContext)
	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := c.snapshot(localContext.Chain, headerNumber-1, header.ParentHash, nil)
	if err != nil {
		return nil, err
	}

	// new span
	span.ID = newSpanID
	if span.EndBlock == 0 {
		span.StartBlock = 256
	} else {
		span.StartBlock = span.EndBlock + 1
	}
	span.EndBlock = span.StartBlock + (100 * c.config.Sprint) - 1

	selectedProducers := make([]Validator, len(snap.ValidatorSet.Validators))
	for i, v := range snap.ValidatorSet.Validators {
		selectedProducers[i] = *v
	}
	heimdallSpan := &HeimdallSpan{
		Span:              *span,
		ValidatorSet:      *snap.ValidatorSet,
		SelectedProducers: selectedProducers,
		ChainID:           c.chainConfig.ChainID.String(),
		AccountRootHash:  "0xb4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30",
	}

	return heimdallSpan, nil
}

//
// Chain context
//

// chain context
type chainContext struct {
	Chain consensus.ChainHeaderReader
	Bor   consensus.Engine
}

func (c chainContext) Engine() consensus.Engine {
	return c.Bor
}

func (c chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	return c.Chain.GetHeader(hash, number)
}

// callmsg implements core.Message to allow passing it as a transaction simulator.
type callmsg struct {
	ethereum.CallMsg
}

func (m callmsg) From() common.Address { return m.CallMsg.From }
func (m callmsg) Nonce() uint64        { return 0 }
func (m callmsg) CheckNonce() bool     { return false }
func (m callmsg) To() *common.Address  { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int   { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64          { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int      { return m.CallMsg.Value }
func (m callmsg) Data() []byte         { return m.CallMsg.Data }

// get system message
func getSystemMessage(toAddress common.Address, data []byte) callmsg {
	return callmsg{
		ethereum.CallMsg{
			From:     systemAddress,
			Gas:      math.MaxUint64 / 2,
			GasPrice: big.NewInt(0),
			Value:    big.NewInt(0),
			To:       &toAddress,
			Data:     data,
		},
	}
}

// apply message
func applyMessage(
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainConfig *params.ChainConfig,
	chainContext core.ChainContext,
) error {
	// Create a new context to be used in the EVM environment
	blockContext := core.NewEVMBlockContext(header, chainContext, &header.Coinbase)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, state, chainConfig, vm.Config{})
	// Apply the transaction to the current state (included in the env)
	_, _, err := vmenv.Call(
		vm.AccountRef(msg.From()),
		*msg.To(),
		msg.Data(),
		msg.Gas(),
		msg.Value(),
	)
	// Update the state with pending changes
	if err != nil {
		state.Finalise(true)
	}

	return nil
}

func validatorContains(a []*Validator, x *Validator) (*Validator, bool) {
	for _, n := range a {
		if bytes.Compare(n.Address.Bytes(), x.Address.Bytes()) == 0 {
			return n, true
		}
	}
	return nil, false
}

func getUpdatedValidatorSet(oldValidatorSet *ValidatorSet, newVals []*Validator) *ValidatorSet {
	v := oldValidatorSet
	oldVals := v.Validators

	var changes []*Validator
	for _, ov := range oldVals {
		if f, ok := validatorContains(newVals, ov); ok {
			ov.VotingPower = f.VotingPower
		} else {
			ov.VotingPower = 0
		}

		changes = append(changes, ov)
	}

	for _, nv := range newVals {
		if _, ok := validatorContains(changes, nv); !ok {
			changes = append(changes, nv)
		}
	}

	v.UpdateWithChangeSet(changes)
	return v
}

func isSprintStart(number, sprint uint64) bool {
	return number%sprint == 0
}

func accumulateRewards(config *params.ChainConfig, state *state.StateDB, header *types.Header, uncles []*types.Header) {
	// Skip block reward in catalyst mode
	if config.IsCatalyst(header.Number) {
		return
	}
	blockReward := FrontierBlockReward
	if config.IsByzantium(header.Number) {
		blockReward = ByzantiumBlockReward
	}
	if config.IsConstantinople(header.Number) {
		blockReward = ConstantinopleBlockReward
	}
	reward := new(big.Int).Set(blockReward)
	state.AddBalance(header.Coinbase, reward)
}

/*func (c *Bor) getValidatorPower(signer common.Address, headerHash common.Hash) uint64 {
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)
	method := "getValidatorPower"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := borcontracts.FABI.Pack(method, signer)
	if err != nil {
		fmt.Println(err)
	}
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(validatorFileCoinContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	var ret = new(*big.Int)
	if err := borcontracts.FABI.UnpackIntoInterface(ret, method, result); err != nil {
		return 0
	}
	u := (*ret).Uint64()
	return u
}*/

func (c *Bor) getSealInfoAllBySigner(signer common.Address, headerHash common.Hash) ([]w3fsStorageManager.W3fsStorageManagerSector, error) {
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)
	method := "getSealInfoAllBySigner"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nelaAbi, _ := w3fsStorageManager.W3fsStorageManagerMetaData.GetAbi()
	data, err := (*nelaAbi).Pack(method, signer , false)
	if err != nil {
		return nil, err
	}
	msgData := (hexutil.Bytes)(data)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	toAddress := common.HexToAddress(W3fsStorageManagerAddress)
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		panic(err)
	}
	res, err := (*nelaAbi).Unpack(method, result)
	if err != nil {
		return *new([]w3fsStorageManager.W3fsStorageManagerSector) , err
	}
	out := *abi.ConvertType(res[0], new([]w3fsStorageManager.W3fsStorageManagerSector)).(*[]w3fsStorageManager.W3fsStorageManagerSector)
	return out, nil
}


/*func (c *Bor) getValidatorVotes(signer common.Address, headerHash common.Hash) ([]*borcontracts.Cvote, error) {
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)
	method := "getValidatorVotes"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := borcontracts.FABI.Pack(method, signer)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(validatorFileCoinContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		panic(err)
		// return nil, err
	}

	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
		ret2 = new([][]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}

	if err := borcontracts.FABI.UnpackIntoInterface(out, method, result); err != nil {
		return nil, err
	}
	valz := make([]*borcontracts.Cvote, len(*ret0))
	for i, _ := range *ret0 {
		valz[i] = &borcontracts.Cvote{
			SectorInx:     (*ret0)[i].Uint64(),
			SealProofType: (*ret1)[i].Uint64(),
			SealedCID:     (*ret2)[i],
		}
	}
	return valz, nil
}*/

/*func (c *Bor) getTotalPower(snapshotNumber uint64) (uint64, error) {
	blockNr := rpc.BlockNumber(snapshotNumber)
	method := "totalPower"
	data, err := borcontracts.FABI.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for LastStateId", "error", err)
		return 0, err
	}

	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(borcontracts.ValidatorFileCoinContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(context.Background(), ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNr}, nil)
	if err != nil {
		return 0, err
	}

	var ret = new(*big.Int)
	if err := borcontracts.FABI.UnpackIntoInterface(ret, method, result); err != nil {
		return 0, err
	}
	return (*ret).Uint64(), nil
}*/

/*func (c *Bor) getAllValidatorPower(addrs []common.Address, headerHash common.Hash) ([]*ValidatorSealPower, error) {
	if len(addrs) == 0 {
		return nil, errors.New("addrs is empty")
	}
	var validators []MinimalAddrVal
	for _, value := range addrs {
		validators = append(validators, MinimalAddrVal{
			Signer: value,
		})
	}
	validatorBytes, _ := rlp.EncodeToBytes(validators)

	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)

	// method
	method := "getAllValidatorPower"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err := borcontracts.FABI.Pack(method, validatorBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(validatorFileCoinContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := c.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		panic(err)
		// return nil, err
	}
	var (
		ret0 = new([]common.Address)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}

	if err := borcontracts.FABI.UnpackIntoInterface(out, method, result); err != nil {
		return nil, err
	}
	valz := make([]*ValidatorSealPower, len(*ret0))
	for i, a := range *ret0 {
		valz[i] = &ValidatorSealPower{
			Address:   a,
			SealPower: (*ret1)[i].Uint64(),
		}
	}
	return valz, nil
}*/

func (c *Bor) printlnValidatorDelay(snap *Snapshot) {
	validators := snap.ValidatorSet.Validators
	for _, value := range validators {
		succession, _ := snap.GetSignerSuccessionNumber(value.Address)
		log.Info(fmt.Sprintf("address =  %s , ProposerPriority = %d , VotingPower = %d, delay = %d", value.Address, value.ProposerPriority, value.VotingPower, succession))
	}
}

func (c *Bor) getWinningPoStSectorChallengeVote(signer common.Address, blockNumber *big.Int, rbaseData []byte, sectorLen uint64, sealType uint64) ([]uint64, error) {
	mid := c.changeMidByEthAddress(signer)
	prand, _ := sealing.GetTicket(blockNumber, rbaseData, fmt.Sprintf("%s%s", "t0", strconv.Itoa(int(mid))), sealing.DomainSeparationTag_WinningPoStChallengeSeed)
	sectorSizeInt, _ := units.RAMInBytes(sealing.GetSealProofType(sealType))
	sealProofType, _ := sealing.SealProofTypeFromSectorSize(fabi.SectorSize(sectorSizeInt), network.Version0)
	sectorChallengeIndex, err := ffiwrapper.ProofVerifier.GenerateWinningPoStSectorChallenge(context.TODO(), fabi.RegisteredPoStProof(uint64(sealProofType)), fabi.ActorID(mid), prand[:], sectorLen)
	if err != nil {
		return nil, err
	} else {
		return sectorChallengeIndex, nil
	}
}





// generate challenge rand by parentHash
func (c *Bor) generateRandData(header *types.Header) []byte {
	randData := blake2b.Sum256(header.ParentHash.Bytes())
	return randData[:]
}


func (c *Bor) changeMidByEthAddress(address common.Address) uint64 {
	encode := hexutil.Encode(address.Bytes())
	mid, _ := strconv.ParseUint(encode[2:6], 16, 32)
	return mid
}
