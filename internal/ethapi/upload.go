package ethapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const (
	// File Upload Status
	AllUpload      = 0 // Full amount to upload
	ContinueUpload = 1 // Continue to upload
	SuccUpload     = 2 // It has been uploaded before
)

var StoreDir = ""

// Query the response body of the file upload status
type CheckState struct {
	State     int      `json:"state"`     // The current uploading status of the file
	ChunkList []string `json:"chunkList"` // Block number that has been uploaded successfully
	File      string   `json:"file"`
}

type FileOK struct {
	File string `json:"file"`
}

// Check whether the file exists
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// Check whether the directory exists
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Calculate the SHA256 of the file
func GetHash(str string) (string, error) {
	file, err := os.Open(str)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sha256h := sha256.New()
	_, err = io.Copy(sha256h, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha256h.Sum(nil)), nil
}

func IsHash(str string) bool {
	pattern := "^0x[0-9a-fA-F]{64}$"
	return IsNumber(pattern, str)
}

func IsNumber(pattern string, str string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}

func CancelUpload(w http.ResponseWriter, r *http.Request) {
	oriHash := r.URL.Query().Get("hash")

	log.Info("CancelUpload", "oriHash", oriHash)

	if !IsHash(oriHash) {
		log.Error("cancelUpload request parameter format error", "oriHash", oriHash)
		http.Error(w, "cancel failed", http.StatusBadRequest)
		return
	}
	tempDir := "." + oriHash
	if IsDir(path.Join(StoreDir, tempDir)) {
		// del block dir
		err := os.RemoveAll(path.Join(StoreDir, tempDir))
		if err != nil {
			log.Error("cancelUpload del block dir error", "oriHash", oriHash, "err", err.Error())
			http.Error(w, "cancel failed", http.StatusInternalServerError)
			return
		}
	}
	return
}

func MergeChunk(w http.ResponseWriter, r *http.Request) {
	oriHash := r.URL.Query().Get("hash")
	fileName := r.URL.Query().Get("name")

	log.Info("MergeChunk", "oriHash", oriHash, "fileName", fileName)

	if !IsHash(oriHash) {
		log.Error("mergeChunk request parameter format error", "oriHash", oriHash)
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}

	tempDir := "." + oriHash
	pathDir := path.Join(StoreDir, tempDir, "*")

	// Gets the names of all block files under this folder
	chunkList, err := filepath.Glob(pathDir)
	if err != nil {
		log.Error("mergeChunk get dir all files error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	// Create the original file
	file, err := os.OpenFile(path.Join(StoreDir, oriHash), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("mergeChunk open file error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	for i := 1; i < len(chunkList)+1; i++ {
		chunkFile, err := os.OpenFile(path.Join(StoreDir, tempDir, strconv.Itoa(i)), os.O_RDONLY, 0666)
		if err != nil {
			log.Error("mergeChunk open chunk file error", "oriHash", oriHash, "chunkId", i, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadAll(chunkFile)
		if err != nil {
			log.Error("mergeChunk read chunk file error", "oriHash", oriHash, "chunkId", i, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
		_, err = file.Write(bytes)
		if err != nil {
			log.Error("mergeChunk write file error", "oriHash", oriHash, "chunkId", i, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
		chunkFile.Close()
	}
	defer file.Close()

	// Delete block dir
	err = os.RemoveAll(path.Join(StoreDir, tempDir))
	if err != nil {
		log.Error("mergeChunk del block dir error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	sha256Sum, err := GetHash(path.Join(StoreDir, oriHash))
	if err != nil {
		log.Error("mergeChunk get file SHA256 error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	hash := "0x" + sha256Sum
	if hash != oriHash {
		err = os.Remove(path.Join(StoreDir, oriHash))
		if err != nil {
			log.Error("mergeChunk remove file error", "oriHash", oriHash, "err", err.Error())
		}
		log.Error("mergeChunk check hash failure", "oriHash", oriHash, "computational", hash)
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	fileOK := FileOK{File: hash}
	bytes, err := json.Marshal(&fileOK)
	if err != nil {
		log.Error("mergeChunk struct conversion json error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	// upload successful
	_, err = w.Write(bytes)
	if err != nil {
		log.Error("mergeChunk write bytes to writer error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	return
}

func UploadChunk(w http.ResponseWriter, r *http.Request) {
	chunkId := r.URL.Query().Get("chunkid")
	oriHash := r.URL.Query().Get("hash")

	log.Info("UploadChunk", "oriHash", oriHash, "chunkId", chunkId)

	if !IsHash(oriHash) || !IsNumber("^+?[1-9][0-9]*$", chunkId) {
		log.Error("uploadChunk request parameter format error", "oriHash", oriHash, "chunkId", chunkId)
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		log.Error("uploadChunk get form file error", "oriHash", oriHash, "chunkId", chunkId, "err", err.Error())
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	// Create a temporary directory
	tempDir := path.Join(StoreDir, "."+oriHash)
	if !IsDir(tempDir) {
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			log.Error("uploadChunk mkdir error", "oriHash", oriHash, "chunkId", chunkId, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
	}

	chunkFile, err := os.OpenFile(path.Join(tempDir, chunkId), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error("uploadChunk open file error", "oriHash", oriHash, "chunkId", chunkId, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	defer chunkFile.Close()

	_, err = io.Copy(chunkFile, formFile)
	if err != nil {
		log.Error("uploadChunk copy file error", "oriHash", oriHash, "chunkId", chunkId, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	return
}

func CheckChunk(w http.ResponseWriter, r *http.Request) {
	oriHash := r.URL.Query().Get("hash")

	log.Info("CheckChunk", "oriHash", oriHash)

	if !IsHash(oriHash) {
		log.Error("checkChunk request parameter format error", "oriHash", oriHash)
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}

	var chunkList []string

	// a pass
	if IsFile(path.Join(StoreDir, oriHash)) {
		checkState := CheckState{State: SuccUpload, File: oriHash}
		bytes, err := json.Marshal(&checkState)
		if err != nil {
			log.Error("checkChunk SuccUpload struct conversion json error", "oriHash", oriHash, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			log.Error("checkChunk SuccUpload write bytes to writer error", "oriHash", oriHash, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
		log.Info("CheckChunk succ upload", "oriHash", oriHash)
		return
	}

	tempDir := path.Join(StoreDir, "."+oriHash)
	if IsDir(tempDir) {
		// Gets the block that has been uploaded successfully
		files, err := ioutil.ReadDir(tempDir)
		if err != nil {
			log.Error("checkChunk read temp dir error", "oriHash", oriHash, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			fileName := file.Name()
			chunkList = append(chunkList, fileName)
		}

		checkState := CheckState{State: ContinueUpload, ChunkList: chunkList}
		bytes, err := json.Marshal(&checkState)
		if err != nil {
			log.Error("checkChunk ContinueUpload struct conversion json error", "oriHash", oriHash, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			log.Error("checkChunk ContinueUpload write bytes to writer error", "oriHash", oriHash, "err", err.Error())
			http.Error(w, "upload failed", http.StatusInternalServerError)
			return
		}
		log.Info("CheckChunk continue upload", "oriHash", oriHash, "chunkList", chunkList)
		return
	}

	// This is a new file
	checkState := CheckState{State: AllUpload, ChunkList: chunkList}
	bytes, err := json.Marshal(&checkState)
	if err != nil {
		log.Error("checkChunk AllUpload struct conversion json error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Error("checkChunk AllUpload write bytes to writer error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	log.Info("CheckChunk all upload", "oriHash", oriHash)
	return
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	oriHash := r.URL.Query().Get("hash")
	fileName := r.URL.Query().Get("name")

	log.Info("UploadFile", "oriHash", oriHash, "fileName", fileName)

	if !IsHash(oriHash) {
		log.Error("uploadFile request parameter format error", "oriHash", oriHash)
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}

	formFile, _, err := r.FormFile("file")
	if err != nil {
		log.Error("uploadFile get form file error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	file, err := os.OpenFile(path.Join(StoreDir, oriHash), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error("uploadFile open file error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, formFile)
	if err != nil {
		log.Error("uploadFile copy file error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	sha256Sum, err := GetHash(path.Join(StoreDir, oriHash))
	if err != nil {
		log.Error("uploadFile get file SHA256 error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	hash := "0x" + sha256Sum
	if hash != oriHash {
		file.Close()
		err := os.Remove(path.Join(StoreDir, oriHash))
		if err != nil {
			log.Error("uploadFile remove file error", "oriHash", oriHash, "err", err.Error())
		}
		log.Error("uploadFile check hash failure", "oriHash", oriHash, "computational", hash)
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	fileOK := FileOK{File: hash}
	bytes, err := json.Marshal(&fileOK)
	if err != nil {
		log.Error("uploadFile struct conversion json error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Error("uploadFile write bytes to writer error", "oriHash", oriHash, "err", err.Error())
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}
	return
}

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func GetCurrentTime(w http.ResponseWriter, r *http.Request) {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	w.Write([]byte(currentTime))
	return
}
