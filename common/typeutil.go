package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"math/big"
	"net"
	"strconv"
	"strings"
)

func HexSTrToByte32(str string) [32]byte {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	bint, _ := new(big.Int).SetString(str, 16)
	data := bint.Bytes()
	ret := [32]byte{}
	for i := 0; i < len(data); i++ {
		ret[i] = data[i]
	}
	return ret
}

func Byte32ToHexStr(data [32]byte) string {
	bint := new(big.Int).SetBytes(data[:])
	return "0x" + BigIntToStr64(bint)
}

func BigIntToStr64(bint *big.Int) string {
	val := bint.Text(16)
	len := len(val)
	const MAX = 64
	if len < MAX {
		for i := 0; i < MAX-len; i++ {
			val = "0" + val
		}
	}
	return val
}

func BuildNewOriHashWithSha256(oriHashStr string) string {
	data := []byte(oriHashStr)
	sha256h := sha256.New()
	sha256h.Write(data)
	oriHashTmp := sha256h.Sum(nil)
	tmp := "0x" + hex.EncodeToString(oriHashTmp)
	return tmp
}


func FormatFloat2(f float64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return value
}

func IsIP(str string) bool {
	address := net.ParseIP(str)
	if address == nil {
		return false
	}else {
		return true
	}
}

func IsMultiAddr(ipfsMaddr string) bool {
	_, err := multiaddr.NewMultiaddr(ipfsMaddr)
	if err != nil {
		return false
	}
	return true
}

func CheckPort(ipfsMaddr string,port int) bool {
	addr, err := multiaddr.NewMultiaddr(ipfsMaddr)
	if err != nil {
		return false
	}
	netAddr, err1 := manet.ToNetAddr(addr)
	if err1!=nil {
		return false
	}
	_,portStr, err2 := net.SplitHostPort(netAddr.String())
	if err2!=nil {
		return false
	}
	portInt, _ := strconv.Atoi(portStr)
	if (port == portInt) {
		return true
	}
	return false
}