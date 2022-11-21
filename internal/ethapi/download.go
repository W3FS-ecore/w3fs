package ethapi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/extern/gcache"
	"github.com/ethereum/go-ethereum/log"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var RetrievalFileDir = ""
var Lcm ClientManager
var GCacheForFileStoreEnable = false
var GCacheForFileStore gcache.Cache
var GCacheForMemoryEnable = false
var GCacheForMemory gcache.Cache

func GetFileHash(file os.File) (string, error) {
	sha256h := sha256.New()
	_, err := io.Copy(sha256h, &file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha256h.Sum(nil)), nil
}

func CheckTxhash(txhash string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return re.MatchString(txhash)
}

func download(w http.ResponseWriter, r *http.Request, gKey string, filename string, filePath string) {

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		log.Error("downloadFile open file error", "fileName", filename, "err", err.Error())
		if GCacheForFileStoreEnable && GCacheForFileStore.Has(gKey) {
			GCacheForFileStore.Remove(gKey)
		}
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error("downloadFile get file info error", "fileName", filename, "err", err.Error())
		if GCacheForFileStoreEnable && GCacheForFileStore.Has(gKey) {
			GCacheForFileStore.Remove(gKey)
		}
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	etag, err := GetFileHash(*file)
	if err != nil {
		log.Error("downloadFile get file SHA256 error", "fileName", filename, "err", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// set response head
	w.Header().Add("Accept-Ranges", "bytes")
	w.Header().Add("Content-Disposition", "attachment; filename="+filename)
	// format:Mon, 13 Dec 2021 08:24:29 GMT
	w.Header().Add("Last-Modified", fileInfo.ModTime().Local().UTC().Format(http.TimeFormat))
	w.Header().Add("Etag", "0x"+etag)
	//w.Header().Add("Cache-Control", "no-store") // disable browser cache

	var start, end int64
	// check whether Range is included   format:bytes=10-
	if r := r.Header.Get("Range"); r != "" {
		if strings.Contains(r, "bytes=") && strings.Contains(r, "-") {

			fmt.Sscanf(r, "bytes=%d-%d", &start, &end)
			// return all
			if end == 0 {
				end = fileInfo.Size() - 1
			}
			// already downloaded
			if start == fileInfo.Size() {
				w.WriteHeader(http.StatusPartialContent)
				return
			}
			if start > end || start < 0 || end < 0 || end >= fileInfo.Size() {
				log.Error("downloadFile range error", "fileName", filename, "start", start, "end", end, "fileSize", fileInfo.Size())
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}
			w.Header().Add("Content-Length", strconv.FormatInt(end-start+1, 10))
			w.Header().Add("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, fileInfo.Size()))
			w.WriteHeader(http.StatusPartialContent)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// first request returns all
		w.Header().Add("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		start = 0
		end = fileInfo.Size() - 1
	}
	// set the cursor position
	_, err = file.Seek(start, 0)
	if err != nil {
		log.Error("downloadFile set the cursor position error", "fileName", filename, "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	n := 10 * 1024 * 1024
	buf := make([]byte, n)
	for {
		if end-start+1 < int64(n) {
			n = int(end - start + 1)
		}
		_, err := file.Read(buf[:n])
		if err != nil {
			log.Error("downloadFile read bytes error", "fileName", filename, "err", err.Error())
			if err != io.EOF {
				log.Error("downloadFile", "err", err.Error())
			}
			return
		}
		err = nil
		_, err = w.Write(buf[:n])
		if err != nil {
			log.Warn("downloadFile write bytes fail", "fileName", filename, "start", start, "end", end, "fileSize", fileInfo.Size(), "bufSize", n, "err", err.Error())
			return
		}
		start += int64(n)
		if start >= end+1 {
			return
		}
	}
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {

	var (
		gKey        string
		filePath    string
		dstFileName string
	)

	oriHash := r.URL.Query().Get("oriHash")
	storageType := r.URL.Query().Get("storageType")
	if storageType == ENTIRE_FILE {
		storeKey := r.URL.Query().Get("storeKey")
		oriHash = common.BuildNewOriHashWithSha256(storeKey)
	}
	headFlag := r.URL.Query().Get("headFlag")
	txHash := r.URL.Query().Get("txHash")
	fileName := r.URL.Query().Get("fileName")

	flag, err := strconv.ParseBool(headFlag)

	if !IsHash(oriHash) || err != nil {
		log.Error("downloadFile request parameter format error", "oriHash", oriHash, "headFlag", headFlag)
		http.Error(w, "download failed", http.StatusBadRequest)
		return
	}

	ishash := CheckTxhash(txHash)
	gKey = Createretrievalkey(oriHash, flag, txHash)

	rs := Lcm.GetRetrievalStatusExt(gKey, ishash)
	dstFileName = rs.Cid

	log.Info("DownloadFile", "oriHash", oriHash, "headFlag", headFlag, "storeKey", r.URL.Query().Get("storeKey"), "storageType", storageType, "fileName", fileName, "range", r.Header.Get("Range"), "cid", dstFileName)

	if dstFileName == "" {
		log.Error("downloadFile get cid error", "oriHash", oriHash, "headFlag", headFlag)
		http.Error(w, "download failed", http.StatusBadRequest)
		return
	}

	if fileName == "" {
		fileName = dstFileName
	}

	filePath = path.Join(RetrievalFileDir, dstFileName)
	download(w, r, gKey, fileName, filePath)
}
