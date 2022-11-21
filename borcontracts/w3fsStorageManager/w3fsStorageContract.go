package w3fsStorageManager

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

var (
	w3fsStorageManagerAddress  = common.HexToAddress("0x0000000000000000000000000000000000001002")
	GlobalStorageManagerClient *StorageManagerClient
)

type StorageSealInfo struct {
	SealProofType uint64
	SectorNumber  uint64
	TicketEpoch   uint64
	SeedEpoch     uint64
	SealedCID     []byte
	UnsealedCID   []byte
	Proof         []byte
}


type Cvote struct {
	SectorInx     uint64
	SealProofType uint64
	SealedCID     []byte
	Proof 		  []byte
}

type WinningPostProof struct {
	PoStProof  uint64
	ProofBytes []byte
}

type WinningData struct {
	Cvotes            []Cvote
	WinningPostProofs []WinningPostProof
}




type StorageManagerClient struct {
	caller     *W3fsStorageManagerCaller
	transactor *W3fsStorageManagerTransactor
	client     *ethclient.Client
}

func NewW3fsStorageManagerClient(client *ethclient.Client) *StorageManagerClient {
	transactor, _ := NewW3fsStorageManagerTransactor(w3fsStorageManagerAddress, client)
	caller, _ := NewW3fsStorageManagerCaller(w3fsStorageManagerAddress, client)
	return &StorageManagerClient{
		transactor: transactor,
		caller:     caller,
		client:     client,
	}
}

func (sm *StorageManagerClient) PackAddSealInfo(isReal bool, signer common.Address, votes []byte) ([]byte, error) {
	abi, err := W3fsStorageManagerMetaData.GetAbi()
	data, err := (*abi).Pack("addSealInfo", isReal, signer, votes)
	return data, err
}

func (sm *StorageManagerClient) Test(signer common.Address) {
	bySigner, err := sm.caller.GetSealInfoAllBySigner(nil, signer, false)
	if err  == nil {
		fmt.Println(bySigner)
	}
}

func (sm *StorageManagerClient) AddSealInfo(auth *bind.TransactOpts, isReal bool, signer common.Address, votes []byte) (common.Hash, error) {
	tx, err := sm.transactor.AddSealInfo(auth, isReal, signer, votes)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (sm *StorageManagerClient) GetValidatorPower(signer common.Address) uint64 {
	power, err := sm.caller.GetValidatorPower(nil, signer)
	if err != nil {
		return 0
	}
	return power.Uint64()
}

func (sm *StorageManagerClient) GetValidatorPromise(signer common.Address) uint64 {
	promise, err := sm.caller.ValidatorPromise(nil, signer)
	if err != nil {
		log.Error("get validatorPromise error", "error", err)
		return 0
	}
	return promise.Uint64()
}

func (sm *StorageManagerClient) GetValidatorStorageSize(signer common.Address) uint64 {
	storageSize, err := sm.caller.ValidatorStorageSize(nil, signer)
	if err != nil {
		log.Error("get validatorPromise error", "error", err)
	}
	return storageSize.Uint64()
}

func (sm *StorageManagerClient) GenerateAuthObj(ecdsaPrivateKey *ecdsa.PrivateKey, chainID *big.Int, fromAddress common.Address, data []byte) (auth *bind.TransactOpts, err error) {
	callMsg := ethereum.CallMsg{
		To:   &w3fsStorageManagerAddress,
		Data: data,
	}
	gasprice, err := sm.client.SuggestGasPrice(context.Background())
	nonce, err := sm.client.PendingNonceAt(context.Background(), fromAddress)
	callMsg.From = fromAddress
	//gasLimit, err := pc.client.EstimateGas(context.Background(), callMsg)
	gasLimit, err := sm.client.EstimateGasLatest(context.Background(), callMsg)
	//auth = bind.NewKeyedTransactor(ecdsaPrivateKey)
	auth, _ = bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
	auth.GasPrice = gasprice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = gasLimit * 2
	return
}

