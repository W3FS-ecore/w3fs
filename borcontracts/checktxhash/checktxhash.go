package checktxhash

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	logging "github.com/ipfs/go-log/v2"
	"math/big"
	"strconv"
	"strings"
)

var log = logging.Logger("checktxhash")

type VoucherClient struct {
	client *ethclient.Client
}

var big8 = big.NewInt(8)

var VoucherCli *VoucherClient

func SetVoucherCli(cli *ethclient.Client) {
	VoucherCli = &VoucherClient{cli}
}

func GetPubKey(tx *eth.Transaction) (common.Address, []byte) {

	var signer eth.Signer = eth.FrontierSigner{}
	if tx.Protected() {
		signer = eth.NewEIP155Signer(tx.ChainId())
	}

	sighash := signer.Hash(tx)
	Vb, R, S := tx.RawSignatureValues()

	// EIP155 support
	var V byte
	if Vb.Int64() > 28 {
		v := new(big.Int).Sub(Vb, tx.ChainId())
		v = new(big.Int).Sub(v, tx.ChainId())
		v = new(big.Int).Sub(v, big8)
		V = byte(v.Uint64() - 27)
	} else {
		V = byte(Vb.Uint64() - 27)
	}

	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// recover the public key from the signature
	var addr common.Address
	pubkey, _ := crypto.Ecrecover(sighash[:], sig)
	copy(addr[:], crypto.Keccak256(pubkey[1:])[12:])
	return addr, pubkey
}

func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

func Ethinit(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func Ethcreatewallets() (string, string, string) {

	ethprivateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("GenerateKey:", err)
		return "", "", ""
	}
	privateKeyBytes := crypto.FromECDSA(ethprivateKey)
	publicKey := ethprivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", "", ""
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	ethpublicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	return address, hexutil.Encode(ethpublicKeyBytes), hexutil.Encode(privateKeyBytes)
}

func Ethgetaddressbyprivate(privatekey string) string {
	ethprivateKey, err := crypto.HexToECDSA(privatekey)

	if err != nil {
		return ""
	}
	publicKey := ethprivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return ""
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address
}

func InitVoucherFilter(conn *ethclient.Client, contractAddress string) (*VoucherFilterer, error) {
	contractCall, err := NewVoucherFilterer(common.HexToAddress(contractAddress), conn)
	return contractCall, err
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

func (v *VoucherClient) GetOwnerbytxhash(txhash string) common.Address {
	var owneraddr common.Address
	txHash := common.HexToHash(txhash)
	tx, _, err := v.client.TransactionByHash(context.Background(), txHash)

	if err == nil {
		owneraddr, _ = GetPubKey(tx)
	}

	return owneraddr
}

func (v *VoucherClient) TransactionCheck(txhash string) (*VoucherPurchaseVoucher, common.Address, error) {
	var contractAddr common.Address
	voucherInfo := new(VoucherPurchaseVoucher)
	txHash := common.HexToHash(txhash)

	receipt, err := v.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Error("TransactionReceipt:", err)
		return voucherInfo, contractAddr, err
	}

	//contractAddress := common.HexToAddress(NelaShareAddress)
	query := ethereum.FilterQuery{
		FromBlock: receipt.BlockNumber, //big.NewInt(11922)
		ToBlock:   receipt.BlockNumber, //big.NewInt(11922)
		//Addresses: []common.Address{contractAddress,},
	}

	logs, err := v.client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error("FilterLogs:", err)
		return voucherInfo, contractAddr, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(VoucherMetaData.ABI)))
	if err != nil {
		return voucherInfo, contractAddr, err
	}

	eventID := contractAbi.Events["PurchaseVoucher"]

	for _, vLog := range logs {

		if eventID.ID.String() == vLog.Topics[0].String() {

			event, err := contractAbi.EventByID(vLog.Topics[0])

			if err != nil {
				log.Error("Failed to look up ABI method: %v", err)
				return voucherInfo, contractAddr, err
			}

			if event == nil {
				log.Error("We should find a event for topic")
				return voucherInfo, contractAddr, err
			}

			voucherInterface, _ := InitVoucherFilter(v.client, vLog.Address.String())
			contractAddr = vLog.Address
			unpackerr := voucherInterface.contract.UnpackLog(voucherInfo, "PurchaseVoucher", vLog)
			if unpackerr != nil {
				log.Error("UnpackLog:", unpackerr)
				return voucherInfo, contractAddr, unpackerr
			}

		}

	}

	return voucherInfo, contractAddr, err
}
