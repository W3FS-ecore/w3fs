package w3fsAuth

import "C"
import (
	"encoding/base64"
	"encoding/hex"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("w3fsauth")

func W3FS_Register(key string) ([]byte, C.int) {
	b64key := base64.StdEncoding.EncodeToString([]byte(key))
	byteKey := []byte(b64key)
	keylen := len(byteKey)

	keysha256, res := TSK_DigestSha256(keylen, byteKey)

	if res != 0 {
		return []byte{}, res
	}

	newpubkey, newprikey, pkeyid, res := TSK_IdentityIssueEx(keysha256)

	if res != 0 {
		return []byte{}, res
	}

	TSK_SetPublicKey(hex.EncodeToString(newpubkey))

	res = TSK_LoginUser(pkeyid)
	if res != 0 {
		return []byte{}, res
	}

	var nPermission int16 = -1
	prikeylen := len(newprikey)

	res = TSK_SetHoldIdentity(3, pkeyid, nPermission, newprikey, prikeylen)
	if res != 0 {
		return []byte{}, res
	}

	res = TSK_SetHoldIdentity(1, pkeyid, nPermission, newprikey, prikeylen)
	if res != 0 {
		return []byte{}, res
	}

	return keysha256, res
}

func W3FS_Auth(userFilePubkey string, oldflow []byte, lasttime int, nPermissionModify int16) ([]byte, C.int, C.int) {
	flowlen := len(oldflow)
	pubbyte, _ := hex.DecodeString(userFilePubkey)

	log.Info("W3FS_Auth pubbyte:\n", pubbyte)
	log.Info("W3FS_Auth userFilePubkey:\n", userFilePubkey)
	log.Info("W3FS_Auth oldflow:\n", oldflow)
	log.Info("W3FS_Auth oldflow:\n", base64.StdEncoding.EncodeToString(oldflow))
	publen := len(pubbyte)

	pkeyid, res := TSK_GetKeyIdByPubkey(pubbyte)

	if res == 0 {
		log.Error("TSK_GetKeyIdByPubkey :\n", pkeyid)
		return TSK_FileOp_AdjustByFlow(oldflow, flowlen, pkeyid, nPermissionModify, lasttime, pubbyte, publen)
	}

	return []byte{}, 0, res
}
