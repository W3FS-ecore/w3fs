package bor

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth_math "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

type GenesisContractsClient struct {
	validatorSetABI       abi.ABI
	stateReceiverABI      abi.ABI
	ValidatorContract     string
	StateReceiverContract string
	chainConfig           *params.ChainConfig
	ethAPI                *ethapi.PublicBlockChainAPI
}

const validatorsetABI = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"id","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"startBlock","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"endBlock","type":"uint256"}],"name":"NewSpan","type":"event"},{"inputs":[],"name":"BOR_ID","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"CHAIN","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"FIRST_END_BLOCK","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"ROUND_TYPE","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"SPRINT","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"SYSTEM_ADDRESS","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"VOTE_TYPE","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"producers","outputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"power","type":"uint256"},{"internalType":"address","name":"signer","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"spanNumbers","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"spans","outputs":[{"internalType":"uint256","name":"number","type":"uint256"},{"internalType":"uint256","name":"startBlock","type":"uint256"},{"internalType":"uint256","name":"endBlock","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"validators","outputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"power","type":"uint256"},{"internalType":"address","name":"signer","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"currentSprint","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"}],"name":"getSpan","outputs":[{"internalType":"uint256","name":"number","type":"uint256"},{"internalType":"uint256","name":"startBlock","type":"uint256"},{"internalType":"uint256","name":"endBlock","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentSpan","outputs":[{"internalType":"uint256","name":"number","type":"uint256"},{"internalType":"uint256","name":"startBlock","type":"uint256"},{"internalType":"uint256","name":"endBlock","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getNextSpan","outputs":[{"internalType":"uint256","name":"number","type":"uint256"},{"internalType":"uint256","name":"startBlock","type":"uint256"},{"internalType":"uint256","name":"endBlock","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"number","type":"uint256"}],"name":"getSpanByBlock","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"currentSpanNumber","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"}],"name":"getValidatorsTotalStakeBySpan","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"}],"name":"getProducersTotalStakeBySpan","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"},{"internalType":"address","name":"signer","type":"address"}],"name":"getValidatorBySigner","outputs":[{"components":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"power","type":"uint256"},{"internalType":"address","name":"signer","type":"address"}],"internalType":"struct BorValidatorSet.Validator","name":"result","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"},{"internalType":"address","name":"signer","type":"address"}],"name":"isValidator","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"},{"internalType":"address","name":"signer","type":"address"}],"name":"isProducer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"signer","type":"address"}],"name":"isCurrentValidator","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"signer","type":"address"}],"name":"isCurrentProducer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"number","type":"uint256"}],"name":"getBorValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getInitialValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"newSpan","type":"uint256"},{"internalType":"uint256","name":"startBlock","type":"uint256"},{"internalType":"uint256","name":"endBlock","type":"uint256"},{"internalType":"bytes","name":"validatorBytes","type":"bytes"},{"internalType":"bytes","name":"producerBytes","type":"bytes"},{"internalType":"bytes32","name":"stateRootHash","type":"bytes32"}],"name":"commitSpan","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"span","type":"uint256"},{"internalType":"bytes32","name":"dataHash","type":"bytes32"},{"internalType":"bytes","name":"sigs","type":"bytes"}],"name":"getStakePowerBySigs","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getAccountRootHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"rootHash","type":"bytes32"},{"internalType":"bytes32","name":"leaf","type":"bytes32"},{"internalType":"bytes","name":"proof","type":"bytes"}],"name":"checkMembership","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"bytes32","name":"d","type":"bytes32"}],"name":"leafNode","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"bytes32","name":"left","type":"bytes32"},{"internalType":"bytes32","name":"right","type":"bytes32"}],"name":"innerNode","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"}]`
const stateReceiverABI = `[{"constant":true,"inputs":[],"name":"SYSTEM_ADDRESS","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"lastStateId","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"syncTime","type":"uint256"},{"internalType":"bytes","name":"recordBytes","type":"bytes"}],"name":"commitState","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
const W3fsStakeManagerABI = `[{"inputs":[],"stateMutability":"payable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"miner","type":"address"}],"name":"ReceiverRewardEvent","type":"event"},{"inputs":[],"name":"COMMISSION_UPDATE_DELAY","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"INIT_MINERSET_BYTES","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"NFTContract","outputs":[{"internalType":"contract W3fsStakingNFT","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"SYSTEM_ADDRESS","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"UNSTAKE_CLAIM_DELAY","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"delegateShareFactory","outputs":[{"internalType":"contract DelegateShareFactory","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"governance","outputs":[{"internalType":"contract IGovernance","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"lock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"locked","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"logger","outputs":[{"internalType":"contract W3fsStakingInfo","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"registry","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"signerToFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"signerToStorageMiner","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"signers","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"storageManager","outputs":[{"internalType":"contract IW3fsStorageManager","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"storageMiners","outputs":[{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"reward","type":"uint256"},{"internalType":"uint256","name":"activationEpoch","type":"uint256"},{"internalType":"uint256","name":"deactivationEpoch","type":"uint256"},{"internalType":"uint256","name":"jailTime","type":"uint256"},{"internalType":"address","name":"signer","type":"address"},{"internalType":"address","name":"contractAddress","type":"address"},{"internalType":"enum W3fsStakeManagerStorage.Status","name":"status","type":"uint8"},{"internalType":"uint256","name":"commissionRate","type":"uint256"},{"internalType":"uint256","name":"lastCommissionUpdate","type":"uint256"},{"internalType":"uint256","name":"delegatorsReward","type":"uint256"},{"internalType":"uint256","name":"delegatedAmount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"token","outputs":[{"internalType":"contract IERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unlock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"userFeeExit","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"stateMutability":"payable","type":"receive"},{"inputs":[{"internalType":"address","name":"_owner","type":"address"},{"internalType":"address","name":"_registry","type":"address"},{"internalType":"address","name":"_token","type":"address"},{"internalType":"address","name":"_NFTContract","type":"address"},{"internalType":"address","name":"_governance","type":"address"},{"internalType":"address","name":"_stakingLogger","type":"address"},{"internalType":"address","name":"_delegateShareFactory","type":"address"},{"internalType":"address","name":"_storageManager","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"},{"internalType":"uint256","name":"storagePromise","type":"uint256"},{"internalType":"bool","name":"acceptDelegation","type":"bool"},{"internalType":"bytes","name":"signerPubkey","type":"bytes"}],"name":"stakeFor","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"},{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"topupFee","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"uint256","name":"accumFeeAmount","type":"uint256"},{"internalType":"uint256","name":"index","type":"uint256"},{"internalType":"bytes","name":"proof","type":"bytes"}],"name":"claimFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"minerAddr","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"updateRewardsMiner","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"minerAddr","type":"address"},{"internalType":"uint256","name":"slashAmount","type":"uint256"},{"internalType":"bool","name":"doJail","type":"bool"}],"name":"slash","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"miner","type":"address"}],"name":"receiverReward","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"commissionRate","type":"uint256"},{"internalType":"uint256","name":"stakeAmount","type":"uint256"},{"internalType":"uint256","name":"delegatedAmount","type":"uint256"}],"name":"_increaseReward","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"restake","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"}],"name":"unstake","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"}],"name":"unstakeClaim","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"}],"name":"unjail","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"},{"internalType":"int256","name":"amount","type":"int256"}],"name":"updateMinerState","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"decreaseMinerDelegatedAmount","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"}],"name":"withdrawDelegatorsReward","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"user","type":"address"}],"name":"getSorageMinerId","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"address","name":"delegator","type":"address"}],"name":"transferFunds","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newValue","type":"uint256"}],"name":"updateStorageMinerThreshold","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newSpanDuration","type":"uint256"}],"name":"updateSpanDuration","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newReward","type":"uint256"}],"name":"updateMinerReward","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"new_unstake_claim_delay","type":"uint256"},{"internalType":"uint256","name":"new_commission_update_delay","type":"uint256"}],"name":"updateDelay","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"}],"name":"forceUnstake","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"},{"internalType":"uint256","name":"newCommissionRate","type":"uint256"}],"name":"updateCommissionRate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"number","type":"uint256"}],"name":"getCurrentEpoch","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBorMiners","outputs":[{"internalType":"address[]","name":"","type":"address[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getInitialValidators","outputs":[{"internalType":"address[]","name":"","type":"address[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"minerAddr","type":"address"}],"name":"isActiveMiner","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"}],"name":"getMinerBaseInfo","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"minerAddr","type":"address"}],"name":"getMinerId","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"withdrawalDelay","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
const SlashManagerABI = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"miner","type":"address"}],"name":"SlashEvent","type":"event"},{"inputs":[],"name":"PERCENTAGE_SLASH","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"felonyThreshold","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"governance","outputs":[{"internalType":"contract IGovernance","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"indicators","outputs":[{"internalType":"uint256","name":"height","type":"uint256"},{"internalType":"uint256","name":"count","type":"uint256"},{"internalType":"uint256","name":"totalCount","type":"uint256"},{"internalType":"uint256","name":"jailCount","type":"uint256"},{"internalType":"uint256","name":"prevAmount","type":"uint256"},{"internalType":"bool","name":"exist","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"jailThreshold","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"lock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"locked","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"logger","outputs":[{"internalType":"contract W3fsStakingInfo","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"registry","outputs":[{"internalType":"contract Registry","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"slashMap","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unlock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_registry","type":"address"},{"internalType":"address","name":"_governance","type":"address"},{"internalType":"address","name":"_w3fsStakingInfo","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"minerAddr","type":"address"}],"name":"slash","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newPercentage","type":"uint256"}],"name":"updatePercentageSlash","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newFelonyThreshold","type":"uint256"}],"name":"updateFelonyThreshold","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newJailThreshold","type":"uint256"}],"name":"updateJailThreshold","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
const W3fsStorageManagerABI = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"newPower","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"newStorageSize","type":"uint256"},{"indexed":true,"internalType":"address","name":"signer","type":"address"},{"indexed":false,"internalType":"uint256","name":"nonce","type":"uint256"}],"name":"AddNewSealPowerAndSize","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"signer","type":"address"},{"indexed":true,"internalType":"uint256","name":"storageSize","type":"uint256"}],"name":"UpdatePromise","type":"event"},{"inputs":[],"name":"baseStakeAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"delegatedStakeLimit","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"factor","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"governance","outputs":[{"internalType":"contract IGovernance","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"lock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"locked","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"logger","outputs":[{"internalType":"contract W3fsStakingInfo","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"percentage","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"registry","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"stakeLimit","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalPower","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unlock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"validatorNonce","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"validatorPowers","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"validatorPromise","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"validatorSector","outputs":[{"internalType":"uint256","name":"SealProofType","type":"uint256"},{"internalType":"uint256","name":"SectorNumber","type":"uint256"},{"internalType":"uint256","name":"TicketEpoch","type":"uint256"},{"internalType":"uint256","name":"SeedEpoch","type":"uint256"},{"internalType":"bytes","name":"SealedCID","type":"bytes"},{"internalType":"bytes","name":"UnsealedCID","type":"bytes"},{"internalType":"bytes","name":"Proof","type":"bytes"},{"internalType":"bool","name":"Check","type":"bool"},{"internalType":"bool","name":"isReal","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"validatorStorageSize","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_owner","type":"address"},{"internalType":"address","name":"_registry","type":"address"},{"internalType":"address","name":"_governance","type":"address"},{"internalType":"address","name":"_stakingLogger","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"validatorAddr","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"addStakeMount","type":"uint256"}],"name":"checkCanStakeMore","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"minerId","type":"uint256"},{"internalType":"uint256","name":"addStakeMount","type":"uint256"}],"name":"checkCandelegatorsMore","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"validatorAddr","type":"address"}],"name":"showCanStakeAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"signer","type":"address"},{"internalType":"uint256","name":"storageSize","type":"uint256"}],"name":"updateStoragePromise","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newStakeLimit","type":"uint256"}],"name":"updateStakeLimit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newDelegatedStakeLimit","type":"uint256"}],"name":"updateDelegatedStakeLimit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"newPercentage","type":"uint256"}],"name":"updatePercentage","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"signer","type":"address"},{"internalType":"uint256","name":"sectorNumber","type":"uint256"}],"name":"getSealInfo","outputs":[{"components":[{"internalType":"uint256","name":"SealProofType","type":"uint256"},{"internalType":"uint256","name":"SectorNumber","type":"uint256"},{"internalType":"uint256","name":"TicketEpoch","type":"uint256"},{"internalType":"uint256","name":"SeedEpoch","type":"uint256"},{"internalType":"bytes","name":"SealedCID","type":"bytes"},{"internalType":"bytes","name":"UnsealedCID","type":"bytes"},{"internalType":"bytes","name":"Proof","type":"bytes"},{"internalType":"bool","name":"Check","type":"bool"},{"internalType":"bool","name":"isReal","type":"bool"}],"internalType":"struct W3fsStorageManager.Sector","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"signer","type":"address"},{"internalType":"bool","name":"isCheck","type":"bool"}],"name":"getSealInfoAllBySigner","outputs":[{"components":[{"internalType":"uint256","name":"SealProofType","type":"uint256"},{"internalType":"uint256","name":"SectorNumber","type":"uint256"},{"internalType":"uint256","name":"TicketEpoch","type":"uint256"},{"internalType":"uint256","name":"SeedEpoch","type":"uint256"},{"internalType":"bytes","name":"SealedCID","type":"bytes"},{"internalType":"bytes","name":"UnsealedCID","type":"bytes"},{"internalType":"bytes","name":"Proof","type":"bytes"},{"internalType":"bool","name":"Check","type":"bool"},{"internalType":"bool","name":"isReal","type":"bool"}],"internalType":"struct W3fsStorageManager.Sector[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bool","name":"isReal","type":"bool"},{"internalType":"address","name":"signer","type":"address"},{"internalType":"bytes","name":"votes","type":"bytes"}],"name":"addSealInfo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"data","type":"bytes"},{"internalType":"uint256[3][]","name":"sigs","type":"uint256[3][]"}],"name":"checkSealSigs","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"signer","type":"address"}],"name":"getValidatorPower","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"validatorBytes","type":"bytes"}],"name":"getAllValidatorPower","outputs":[{"internalType":"address[]","name":"","type":"address[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"}]`
const W3fsStakeManagerAddress = "0x0000000000000000000000000000000000001003"
const SlashManagerAddress = "0x0000000000000000000000000000000000001005"
const W3fsStorageManagerAddress = "0x0000000000000000000000000000000000001002"

func NewGenesisContractsClient(
	chainConfig *params.ChainConfig,
	validatorContract,
	stateReceiverContract string,
	ethAPI *ethapi.PublicBlockChainAPI,
) *GenesisContractsClient {
	vABI, _ := abi.JSON(strings.NewReader(validatorsetABI))
	sABI, _ := abi.JSON(strings.NewReader(stateReceiverABI))
	return &GenesisContractsClient{
		validatorSetABI:       vABI,
		stateReceiverABI:      sABI,
		ValidatorContract:     validatorContract,
		StateReceiverContract: stateReceiverContract,
		chainConfig:           chainConfig,
		ethAPI:                ethAPI,
	}
}

func (gc *GenesisContractsClient) CommitState(
	event *EventRecordWithTime,
	state *state.StateDB,
	header *types.Header,
	chCtx chainContext,
) error {
	eventRecord := event.BuildEventRecord()
	recordBytes, err := rlp.EncodeToBytes(eventRecord)
	if err != nil {
		return err
	}
	method := "commitState"
	t := event.Time.Unix()
	data, err := gc.stateReceiverABI.Pack(method, big.NewInt(0).SetInt64(t), recordBytes)
	if err != nil {
		log.Error("Unable to pack tx for commitState", "error", err)
		return err
	}
	log.Info("→ committing new state", "eventRecord", event.String())
	msg := getSystemMessage(common.HexToAddress(gc.StateReceiverContract), data)
	if err := applyMessage(msg, state, header, gc.chainConfig, chCtx); err != nil {
		return err
	}
	return nil
}

func (gc *GenesisContractsClient) LastStateId(snapshotNumber uint64) (*big.Int, error) {
	blockNr := rpc.BlockNumber(snapshotNumber)
	method := "lastStateId"
	data, err := gc.stateReceiverABI.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for LastStateId", "error", err)
		return nil, err
	}

	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(gc.StateReceiverContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := gc.ethAPI.Call(context.Background(), ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNr}, nil)
	if err != nil {
		return nil, err
	}

	var ret = new(*big.Int)
	if err := gc.stateReceiverABI.UnpackIntoInterface(ret, method, result); err != nil {
		return nil, err
	}
	return *ret, nil
}

type GenesisW3fsContractsClient struct {
	w3fsStakeManagerABI        abi.ABI
	slashManagerABI            abi.ABI
	w3fsStorageManagerABI      abi.ABI
	W3fsStakeManagerContract   string
	SlashManagerContract       string
	W3fsStorageManagerContract string
	chainConfig                *params.ChainConfig
	ethAPI                     *ethapi.PublicBlockChainAPI
	Bor                        *Bor
	StartSlashNumber			uint64
}

func NewGenesisW3fsContractsClient(chainConfig *params.ChainConfig, ethAPI *ethapi.PublicBlockChainAPI) *GenesisW3fsContractsClient {
	nABI, _ := abi.JSON(strings.NewReader(W3fsStakeManagerABI))
	sABI, _ := abi.JSON(strings.NewReader(SlashManagerABI))
	nsABI, _ := abi.JSON(strings.NewReader(W3fsStorageManagerABI))
	return &GenesisW3fsContractsClient{
		w3fsStakeManagerABI:        nABI,
		slashManagerABI:            sABI,
		w3fsStorageManagerABI:      nsABI,
		W3fsStakeManagerContract:   W3fsStakeManagerAddress,
		SlashManagerContract:       SlashManagerAddress,
		W3fsStorageManagerContract: W3fsStorageManagerAddress,
		chainConfig:                chainConfig,
		ethAPI:                     ethAPI,
		StartSlashNumber:			chainConfig.Bor.Sprint * 10,
	}
}

func (gnc *GenesisW3fsContractsClient) GetBorMiners(headerHash common.Hash, isAddPower bool) ([]*Validator, []*ValidatorSealPower, error) {
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)
	// method
	method := "getBorMiners"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := gnc.w3fsStakeManagerABI.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for getValidator", "error", err)
		return nil, nil, err
	}
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(gnc.W3fsStakeManagerContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := gnc.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Error("call getBorMiners error", "error", err)
	}
	var (
		ret0 = new([]common.Address)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	if err := gnc.w3fsStakeManagerABI.UnpackIntoInterface(out, method, result); err != nil {
		return nil, nil, err
	}
	validatorSet := make([]*Validator, len(*ret0))
	for i, a := range *ret0 {
		validatorSet[i] = &Validator{
			Address:     a,
			VotingPower: gnc.weiToEth((*ret1)[i]),
		}
	}
	if isAddPower {
		var addrs []common.Address
		for _, v := range validatorSet {
			addrs = append(addrs, v.Address)
		}
		validatorSealPowers, err := gnc.getAllValidatorPower(headerHash, addrs)
		if err == nil {
			for _, vs := range validatorSealPowers {
				log.Info(fmt.Sprintf("address =  %s , power = %d", hexutil.Encode(vs.Address.Bytes()), vs.SealPower))
			}
		}
		return validatorSet, validatorSealPowers, nil
	} else {
		var validatorSealPowers []*ValidatorSealPower
		for _, v := range validatorSet {
			validatorSealPowers = append(validatorSealPowers, &ValidatorSealPower{
				Address:   v.Address,
				SealPower: 0,
			})
		}
		return validatorSet, validatorSealPowers, nil
	}
}

func (gnc *GenesisW3fsContractsClient) getAllValidatorPower(headerHash common.Hash, minerAddr []common.Address) ([]*ValidatorSealPower, error) {
	blockNr := rpc.BlockNumberOrHashWithHash(headerHash, false)
	method := "getAllValidatorPower"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var validators []MinimalAddrVal
	for _, value := range minerAddr {
		validators = append(validators, MinimalAddrVal{
			Signer: value,
		})
	}
	validatorBytes, _ := rlp.EncodeToBytes(validators)
	data, err := gnc.w3fsStorageManagerABI.Pack(method, validatorBytes)
	if err != nil {
		log.Error("Unable to pack tx for getValidator", "error", err)
		return nil, err
	}
	msgData := (hexutil.Bytes)(data)
	toAddress := common.HexToAddress(gnc.W3fsStorageManagerContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))
	result, err := gnc.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, blockNr, nil)
	if err != nil {
		log.Error("call getBorMiners error", "error", err)
	}
	var (
		ret0 = new([]common.Address)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}

	if err := gnc.w3fsStorageManagerABI.UnpackIntoInterface(out, method, result); err != nil {
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
}

func (gnc *GenesisW3fsContractsClient) weiToEth(wei *big.Int) int64 {
	weiToEth := eth_math.BigPow(10, 18)
	div := new(big.Int).Div(wei, weiToEth)
	return div.Int64()
}

// receiver reward
func (gnc *GenesisW3fsContractsClient) receiverReward(
	state *state.StateDB,
	header *types.Header,
	chCtx chainContext,
) error {
	method := "receiverReward"
	data, err := gnc.w3fsStakeManagerABI.Pack(method, header.Coinbase)
	if err != nil {
		log.Error("Unable to pack tx for receiverReward", "error", err)
		return err
	}
	msg := getSystemMessage(common.HexToAddress(gnc.W3fsStakeManagerContract), data)
	if err := applyMessage(msg, state, header, gnc.chainConfig, chCtx); err != nil {
		return err
	}
	return nil
}

func (gnc *GenesisW3fsContractsClient) checkDoSlash(snap *Snapshot, headerNumber uint64, sprint uint64, signer common.Address) bool {
	if gnc.StartSlashNumber < headerNumber && bytes.Compare(snap.ValidatorSet.GetProposer().Address.Bytes(), signer.Bytes()) != 0 {
		/*currentValidator := snap.ValidatorSet.GetProposer().Address
		if headerNumber > 0 && !isSprintStart(headerNumber, sprint) {
			if bytes.Compare(currentValidator.Bytes(), snap.Recents[headerNumber-1].Bytes()) != 0 {
				return true
			}
		}*/
	}
	return false
}

func (gnc *GenesisW3fsContractsClient) slash(
	miner common.Address,
	state *state.StateDB,
	header *types.Header,
	chain core.ChainContext,
	txs *[]*types.Transaction,
	receipts *[]*types.Receipt,
	receivedTxs *[]*types.Transaction,
	usedGas *uint64,
	mining bool,
) error {
	method := "slash"
	data, err := gnc.slashManagerABI.Pack(method, miner)
	if err != nil {
		log.Error("Unable to pack tx for slash", "error", err)
		return err
	}
	to := common.HexToAddress(gnc.SlashManagerContract)
	msg := callmsg{
		ethereum.CallMsg{
			From:     header.Coinbase,
			Gas:      math.MaxUint64 / 2,
			GasPrice: big.NewInt(0),
			Value:    big.NewInt(0),
			To:       &to,
			Data:     data,
		},
	}
	return gnc.Bor.applyW3fsTransaction(msg, state, header, chain, txs, receipts, receivedTxs, usedGas, mining)
}

func (c *Bor) applyW3fsTransaction(
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainContext core.ChainContext,
	txs *[]*types.Transaction, receipts *[]*types.Receipt,
	receivedTxs *[]*types.Transaction, usedGas *uint64, mining bool,
) (err error) {
	nonce := state.GetNonce(msg.From())
	expectedTx := types.NewTransaction(nonce, *msg.To(), msg.Value(), msg.Gas(), msg.GasPrice(), msg.Data())
	expectedHash := c.txSigner.Hash(expectedTx)
	if msg.From() == c.signer && mining {
		expectedTx, err = c.signTxFn(accounts.Account{Address: msg.From()}, expectedTx, c.chainConfig.ChainID)
		if err != nil {
			return err
		}
	} else {
		if receivedTxs == nil || len(*receivedTxs) == 0 || (*receivedTxs)[0] == nil {
			return errors.New("supposed to get a actual transaction, but get none")
		}
		actualTx := (*receivedTxs)[0]
		// actualTx应该是收到的交易，判断收到的交易和自己封装的交易的hash是否相同
		if !bytes.Equal(c.txSigner.Hash(actualTx).Bytes(), expectedHash.Bytes()) {
			return fmt.Errorf("expected tx hash %v, get %v, nonce %d, to %s, value %s, gas %d, gasPrice %s, data %s", expectedHash.String(), actualTx.Hash().String(),
				expectedTx.Nonce(),
				expectedTx.To().String(),
				expectedTx.Value().String(),
				expectedTx.Gas(),
				expectedTx.GasPrice().String(),
				hex.EncodeToString(expectedTx.Data()),
			)
		}
		expectedTx = actualTx
		// move to next
		*receivedTxs = (*receivedTxs)[1:]
	}
	state.Prepare(expectedTx.Hash(), len(*txs))
	gasUsed, err := applyW3fsMessage(msg, state, header, c.chainConfig, chainContext)
	if err != nil {
		return err
	}
	*txs = append(*txs, expectedTx)
	var root []byte
	if c.chainConfig.IsByzantium(header.Number) {
		state.Finalise(true)
	} else {
		root = state.IntermediateRoot(c.chainConfig.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gasUsed
	receipt := types.NewReceipt(root, false, *usedGas)
	receipt.TxHash = expectedTx.Hash()
	log.Info("system applyW3fsTransaction", "hash", expectedTx.Hash())
	receipt.GasUsed = gasUsed

	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = state.GetLogs(expectedTx.Hash(), header.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = header.Hash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(state.TxIndex())
	*receipts = append(*receipts, receipt)
	state.SetNonce(msg.From(), nonce+1)
	return nil
}

func applyW3fsMessage(
	msg callmsg,
	state *state.StateDB,
	header *types.Header,
	chainConfig *params.ChainConfig,
	chainContext core.ChainContext,
) (uint64, error) {
	context := core.NewEVMBlockContext(header, chainContext, &header.Coinbase)
	vmenv := vm.NewEVM(context, vm.TxContext{Origin: msg.From(), GasPrice: big.NewInt(0)}, state, chainConfig, vm.Config{})
	ret, returnGas, err := vmenv.Call(
		vm.AccountRef(msg.From()),
		*msg.To(),
		msg.Data(),
		msg.Gas(),
		msg.Value(),
	)
	if err != nil {
		log.Error("apply message failed", "msg", string(ret), "err", err)
	}
	return msg.Gas() - returnGas, err
}

func (c *Bor) IsSystemTransaction(tx *types.Transaction, header *types.Header) (bool, error) {
	if tx.To() == nil {
		return false, nil
	}
	sender, err := types.Sender(c.txSigner, tx)
	if err != nil {
		return false, errors.New("UnAuthorized transaction")
	}
	if sender == header.Coinbase && isToSystemContract(*tx.To()) && tx.GasPrice().Cmp(big.NewInt(0)) == 0 {
		return true, nil
	}
	return false, nil
}

func isToSystemContract(to common.Address) bool {
	systemContracts := map[common.Address]bool{
		common.HexToAddress("0x0000000000000000000000000000000000001003"): true,
		common.HexToAddress("0x0000000000000000000000000000000000001004"): true,
		common.HexToAddress("0x0000000000000000000000000000000000001005"): true,
	}
	return systemContracts[to]
}
