package eth

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/ethereum/go-ethereum/extern/gcache"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

var persistDir = ""
var cacheFile = "cf"
var evictFile = "ef"
var evictItems = make(map[interface{}]interface{})

func NewFileCache(stack *node.Node, size int, ci time.Duration) gcache.Cache {
	g := gcache.New(size).
		ARC().
		EvictedFunc(func(key, value interface{}) {
			// Called when the cache is expelled
			evictItems[key] = value
		}).
		Build()

	persistDir = stack.Config().DataDir + "/.gcache"
	if !fileIsExist(persistDir) {
		err := os.Mkdir(persistDir, 0755)
		if err != nil {
			log.Error("create cache dir fail", "err", err.Error())
		}
	} else {
		// Cache recovery
		err := cacheRecover(g)
		if err != nil {
			log.Error("cache recover fail", "err", err.Error())
			// del local cache file
			err := os.Remove(path.Join(persistDir, cacheFile))
			if err != nil {
				log.Error("del local cache file fail", "err", err.Error())
			}
		}
	}

	// Periodic persistence to disks
	go run(g, ci)
	return g
}

func cacheRecover(g gcache.Cache) error {
	cf := path.Join(persistDir, cacheFile)
	if !fileIsExist(cf) {
		return nil
	}

	rf, err := ioutil.ReadFile(cf)
	if err != nil {
		return err
	}

	items := Decoder(rf)
	for key, value := range items {
		err := g.Set(key, value)
		if err != nil {
			g.Purge()
			return err
		}
	}

	log.Info("cache recover succ", "items", len(g.GetALL(false)))
	return nil
}

func run(g gcache.Cache, ci time.Duration) {
	ticker := time.NewTicker(ci)
	for {
		select {
		case <-ticker.C:
			cacheSync(g)
		}
	}
}

func cacheSync(g gcache.Cache) {
	items := g.GetALL(false)
	log.Info("cache sync", "CacheItems", len(items), "EvictItems", len(evictItems),
		"HitCount", g.HitCount(), "MissCount", g.MissCount(), "LookupCount", g.LookupCount(), "HitRate", g.HitRate())

	if len(items) != 0 {
		cf, _ := os.OpenFile(path.Join(persistDir, cacheFile), os.O_WRONLY|os.O_CREATE, 0666)
		cb := Encoder(items)
		cf.Write(cb)
		cf.Close()

	}
	if len(evictItems) != 0 {
		year := time.Now().Format("2006")
		month := time.Now().Format("01")
		day := time.Now().Format("02")
		fn := evictFile + "_" + year + "_" + month + "_" + day
		ef, _ := os.OpenFile(path.Join(persistDir, fn), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		w := bufio.NewWriter(ef)
		for key, value := range evictItems {
			w.WriteString(str2val(value) + "\n")
			delete(evictItems, key)
		}
		w.Flush()
		ef.Close()
	}
}

// Check whether the file exists
func fileIsExist(fn string) bool {
	var exist = true
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// serialization
func Encoder(items map[interface{}]interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(items)
	return buffer.Bytes()
}

// deserialization
func Decoder(rf []byte) map[interface{}]interface{} {
	var items map[interface{}]interface{}
	decoder := gob.NewDecoder(bytes.NewReader(rf))
	err := decoder.Decode(&items)
	if err != nil {
		log.Error("cache file load fail", "err", err.Error())
	}
	return items
}

// struct to string
func str2val(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func NewMemoryCache(size int) gcache.Cache {
	return gcache.New(size).LRU().Build()
}
