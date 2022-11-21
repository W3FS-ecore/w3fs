package myTest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/multiformats/go-multiaddr"
	"testing"
)

func TestMultiAddress(t *testing.T) {
	var ma multiaddr.Multiaddr
	ipfsMaddr := "/ip4/125.124.162.231/tcp/8545"
	ma, err := multiaddr.NewMultiaddr(ipfsMaddr)
	if err != nil {
		fmt.Errorf("parsing ipfs multiaddr: %w", err)
	}
	fmt.Println(ma.MarshalJSON())
	fmt.Println(ma.String())
	b2 := common.CheckPort(ipfsMaddr, 2092)
	fmt.Println(b2)
}

func Test1(t *testing.T) {
	b := common.IsIP("12.102.2.255")
	fmt.Println("b:", b)
	freeSpace :=1000
	totalFreeSpace :=2000
	radio := float64(freeSpace) / float64(totalFreeSpace)
	FreeSpaceRatio := common.FormatFloat2(radio)
	fmt.Println(FreeSpaceRatio)
	oriHashStr := "0x28e2b887c80e16130bc7aac822267a20ded156e141097de3160332508919a7910x211fA1DDbE3d000e1a42921eC56bBE7A923A6BeD"
	fmt.Println(oriHashStr)
	fmt.Println("0x06ea9248fac2f1dae63786724170ad3aadaa3bac3a4b072bb2f7dc04259cb606")
	data := []byte(oriHashStr)
	//data, _ := hex.DecodeString(oriHashStr)
	// new sha256
	sha256h := sha256.New()
	sha256h.Write(data)
	oriHashTmp := sha256h.Sum(nil)
	fmt.Println("---------------------------")
	//fmt.Printf("%x\n", sha256.Sum256([]byte(oriHashStr)))
	tmp := hex.EncodeToString(oriHashTmp)
	fmt.Println("0x" + tmp)
}
func TestAbc(t *testing.T) {
	str1 := "05872aac4f85c37e1d5fbf663d95cbdb6cb7da9cbcaac805f1197738297f8826"

	data1 := common.HexSTrToByte32(str1)
	fmt.Println(len(data1))

	str2 := "0x05872aac4f85c37e1d5fbf663d95cbdb6cb7da9cbcaac805f1197738297f8826"

	data2 := common.HexSTrToByte32(str2)
	fmt.Println(len(data2))

	str3 := common.Byte32ToHexStr(data2)
	fmt.Println("str3=", str3)
	//
	//fmt.Println("----------------------------------------------------------------------------------------------")
	////    f0ef386506354b62c891f08ad4bf90d9d58f1e32579923cf85baaa81efdcc1f5
	//x := "7b5631180bcd14111972639d07e1837317129fd94fe791d93362bfa7e2b10c9c"
	////bc1b7d6a03a0d3346ed057ca91fe6213742b1135c133f931558e25e028a20b30
	//fmt.Println(x)
	//y := ethapi.Str64ToMinerId(x)
	//fmt.Printf("OK: y=%s \n", y)
	//
	//
	//fmt.Println(ethapi.MinerIdToStr64(y))
	//
	//fmt.Println("      ")
	//fmt.Println("------------------------------------------------------------------------------------------------")
	//z := ethapi.MinerIdToNodeId("5332767211998331076075451763451634259909720778516484101592828555541560887761")
	//fmt.Println(z)
}

func TestGenerateKey(t *testing.T) {
	// GenerateKey
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed GenerateKey with %s.", err)
	}
	fmt.Println(hexutil.Encode(crypto.FromECDSA(key)))
	fmt.Println(hex.EncodeToString(crypto.FromECDSA(key)))
	fmt.Println(hexutil.Encode(crypto.FromECDSAPub(&key.PublicKey)))
	fmt.Println(hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey)))
	fmt.Println(crypto.PubkeyToAddress(key.PublicKey))
}

func TestGenerateKeyAndSaveFile(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed GenerateKey with %s.", err)
	}
	if err := crypto.SaveECDSA("privatekey", key); err != nil {
		t.Fatalf(fmt.Sprintf("Failed to persist node key: %v", err))
	}
}

func TestDemo(t *testing.T) {
	address := common.HexToAddress("1fa9E8F2Bd8F5D35AFbC20cB7aEF6f9e4318A97a")
	privKey, _ := crypto.HexToECDSA("3fea189abf64e5ecab825f22d47fe0660f883ca191b86c10938d6184e7091b84")
	msg := crypto.Keccak256([]byte("bor"))
	sign, _ := crypto.Sign(msg, privKey)
	recoveredPub, _ := crypto.Ecrecover(msg, sign)
	fmt.Println("address：", address, "pubKey：", hexutil.Encode(recoveredPub))

	// get public key
	recoveredPub_2, _ := crypto.SigToPub(msg, sign)
	fmt.Println("address：", crypto.PubkeyToAddress(*recoveredPub_2), "pubKey：", hexutil.Encode(crypto.FromECDSAPub(recoveredPub_2)))

}
