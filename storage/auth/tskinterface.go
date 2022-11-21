package w3fsAuth

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include
// #cgo LDFLAGS: -L./lib -lTSKLinux
// #include "TSKInterface_GO.h"
import "C"
import (
	"encoding/hex"
	"fmt"
	"unsafe"
)

//eth privateKey
var lotuspubkey string = ""

func HexDecode(s string) []byte {
	dst := make([]byte, hex.DecodedLen(len(s)))
	n, err := hex.Decode(dst, []byte(s))
	if err != nil {
		return nil
	}
	return dst[:n]
}

func TSK_Init() C.int {
	res := C.TSK_GO_Init()
	return res
}

func TSK_LoginUser(userid []byte) C.int {
	res := C.TSK_GO_LoginUser((*C.uchar)(unsafe.Pointer(&userid)))
	return res
}

func TSK_DigestSha256(priAddresslen int, bytePriAddress []byte) ([]byte, C.int) {
	priSha256 := make([]byte, C.DIGEST_SHA256_LENGTH)
	res := C.TSK_GO_DigestSha256(*(*C.uint)(unsafe.Pointer(&priAddresslen)), (*C.uchar)(unsafe.Pointer(&bytePriAddress[0])), nil, (*C.uchar)(unsafe.Pointer(&priSha256[0])))
	return priSha256, res
}

func TSK_IdentityIssueEx(priSha256 []byte) ([]byte, []byte, []byte, C.int) {

	var pubkeylen C.UINT32
	var prikeylen C.UINT32
	pubkeylen = 200
	prikeylen = 200

	pubkey := make([]byte, 200)
	prikey := make([]byte, 200)
	pkeyid := make([]byte, 20)

	res := C.TSK_GO_IdentityIssueEx((*C.uchar)(unsafe.Pointer(&priSha256[0])), C.DIGEST_SHA256_LENGTH,
		(*C.uchar)(unsafe.Pointer(&pubkey[0])), &pubkeylen, (*C.uchar)(unsafe.Pointer(&prikey[0])), &prikeylen,
		(*C.uchar)(unsafe.Pointer(&pkeyid[0])))
	ipublen := C.int(pubkeylen)
	iprilen := C.int(prikeylen)

	newpubkey := append(pubkey[:ipublen], pubkey[200:]...)
	newprikey := append(prikey[:iprilen], prikey[200:]...)

	return newpubkey, newprikey, pkeyid, res
}

func TSK_SetHoldIdentity(nListVerb C.UINT32, byteuserid []byte, nPermission int16, bytefileprivate []byte, fileprivatelen int) C.int {
	res := C.TSK_GO_SetHoldIdentity(nListVerb, (*C.uchar)(unsafe.Pointer(&byteuserid[0])), *(*C.UINT16)(unsafe.Pointer(&nPermission)), -1, (*C.uchar)(unsafe.Pointer(&bytefileprivate[0])), *(*C.int)(unsafe.Pointer(&fileprivatelen)))
	return res
}

func TSK_IdentityExportByPrivate(priSha256 []byte, action byte) ([]byte, C.UINT32, C.int) {
	var tskid unsafe.Pointer
	var pubkeylen C.UINT32
	pubkeylen = 200
	key := make([]byte, 200)

	res := C.TSK_GO_IdentityIssue(&tskid, (*C.uchar)(unsafe.Pointer(&priSha256[0])), C.DIGEST_SHA256_LENGTH)
	fmt.Println("tskid:", tskid)
	fmt.Println(res)

	if res == 0 {
		res = C.TSK_GO_IdentityExport(tskid, C.IDENTITY_PUBLIC_KEY, 200, (*C.uchar)(unsafe.Pointer(&key[0])), &pubkeylen)
		key = append(key[:pubkeylen], key[200:]...)

		C.TSK_GO_IdentityFree(tskid)
	}

	return key, pubkeylen, res
}

func TSK_FileOp_AdjustByFlow(headSrc []byte, srcLen int, userid []byte, nPermission int16, lasttime int, pubkey []byte, pubkeylen int) ([]byte, C.int, C.int) {
	var dstLen C.int
	dstHead := make([]byte, srcLen+1024)
	clasttime := C.int(lasttime)
	cpubkeylen := C.int(pubkeylen)
	csrclen := C.int(srcLen)
	dstLen = csrclen + 1024
	res := C.TSK_Go_FileOp_AdjustByFlow((*C.uchar)(unsafe.Pointer(&headSrc[0])),
		csrclen, (*C.uchar)(unsafe.Pointer(&dstHead[0])), &dstLen, true,
		(*C.uchar)(unsafe.Pointer(&userid[0])), *(*C.UINT16)(unsafe.Pointer(&nPermission)), clasttime, (*C.uchar)(unsafe.Pointer(&pubkey[0])), cpubkeylen)
	dstHead = dstHead[:dstLen]
	return dstHead, dstLen, res
}

func TSK_IdentityImport(Action C.uchar, BufLen int, buf []byte) (unsafe.Pointer, C.int) {
	var tskid unsafe.Pointer
	res := C.TSK_GO_IdentityImport(Action, *(*C.uint)(unsafe.Pointer(&BufLen)), (*C.uchar)(unsafe.Pointer(&buf[0])), &tskid)
	return tskid, res
}

func TSK_IdentityExport(tskid unsafe.Pointer, Action C.uchar, BufLen C.uint, buf []byte) C.int {
	var dstLen C.uint
	return C.TSK_GO_IdentityExport(tskid, Action, BufLen, (*C.uchar)(unsafe.Pointer(&buf[0])), &dstLen)
}

func TSK_GetKeyIdByPubkey(newpubkey []byte) ([]byte, C.int) {
	newpublen := len(newpubkey)

	if newpublen == 0 {
		return []byte{}, -1
	}

	pkeyid := make([]byte, 20)

	keyid, res := TSK_IdentityImport(C.IDENTITY_PUBLIC_KEY, newpublen, newpubkey)

	if res == 0 {
		TSK_IdentityExport(keyid, C.IDENTITY_KEY_ID, 20, pkeyid)
		C.TSK_GO_IdentityFree(keyid)
	}

	return pkeyid, res
}

func TSK_GetPublicKey() string {
	return lotuspubkey
}

func TSK_SetPublicKey(key string) {
	lotuspubkey = key
}

func TSK_UnInit() {
	C.TSK_GO_UnInit()
}
