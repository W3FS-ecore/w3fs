module github.com/ethereum/go-ethereum

go 1.15

require (
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-storage-blob-go v0.7.0
	github.com/Azure/go-autorest/autorest/adal v0.8.0 // indirect
	github.com/VictoriaMetrics/fastcache v1.6.0
	github.com/aws/aws-sdk-go-v2 v1.11.2
	github.com/aws/aws-sdk-go-v2/config v1.1.1
	github.com/aws/aws-sdk-go-v2/credentials v1.1.1
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.1.1
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/buger/goterm v1.0.3
	github.com/cespare/cp v0.1.0
	github.com/cloudflare/cloudflare-go v0.14.0
	github.com/consensys/gnark-crypto v0.4.1-0.20210426202927-39ac3d4b3f1f
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/deepmap/oapi-codegen v1.8.2 // indirect
	github.com/detailyang/go-fallocate v0.0.0-20180908115635-432fa640bd2e
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf
	github.com/docker/go-units v0.4.0
	github.com/dop251/goja v0.0.0-20200721192441-a695b0cdd498
	github.com/drand/drand v1.2.1
	github.com/drand/kyber v1.1.4
	github.com/edsrzf/mmap-go v1.0.0
	github.com/fatih/color v1.13.0
	github.com/filecoin-project/dagstore v0.4.3
	github.com/filecoin-project/filecoin-ffi v0.30.4-0.20200910194244-f640612a1a1f
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-bitfield v0.2.4
	github.com/filecoin-project/go-cbor-util v0.0.1
	github.com/filecoin-project/go-commp-utils v0.1.2
	github.com/filecoin-project/go-data-transfer v1.11.6
	github.com/filecoin-project/go-ds-versioning v0.1.0
	github.com/filecoin-project/go-fil-commcid v0.1.0
	github.com/filecoin-project/go-fil-commp-hashhash v0.1.0
	github.com/filecoin-project/go-fil-markets v1.13.3
	github.com/filecoin-project/go-jsonrpc v0.1.5
	github.com/filecoin-project/go-padreader v0.0.1
	github.com/filecoin-project/go-paramfetch v0.0.2
	github.com/filecoin-project/go-state-types v0.1.1
	github.com/filecoin-project/go-statemachine v1.0.1
	github.com/filecoin-project/go-statestore v0.1.1
	github.com/filecoin-project/lotus v1.13.0
	github.com/filecoin-project/specs-actors v0.9.14
	github.com/filecoin-project/specs-actors/v2 v2.3.5
	github.com/filecoin-project/specs-actors/v5 v5.0.4
	github.com/filecoin-project/specs-actors/v6 v6.0.1
	github.com/filecoin-project/specs-storage v0.1.1-0.20201105051918-5188d9774506
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff
	github.com/go-sourcemap/sourcemap v2.1.2+incompatible // indirect
	github.com/go-stack/stack v1.8.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/graph-gophers/graphql-go v0.0.0-20201113091052-beb923fada29
	github.com/hannahhoward/go-pubsub v0.0.0-20200423002714-8d62886cc36e
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/bloomfilter/v2 v2.0.3
	github.com/holiman/uint256 v1.2.0
	github.com/huin/goupnp v1.0.2
	github.com/influxdata/influxdb v1.8.3
	github.com/influxdata/influxdb-client-go/v2 v2.4.0
	github.com/influxdata/line-protocol v0.0.0-20210311194329-9aa0e372d097 // indirect
	github.com/ipfs/go-block-format v0.0.3
	github.com/ipfs/go-blockservice v0.1.7
	github.com/ipfs/go-cid v0.1.0
	github.com/ipfs/go-cidutil v0.0.2
	github.com/ipfs/go-datastore v0.4.6
	github.com/ipfs/go-graphsync v0.10.4
	github.com/ipfs/go-ipfs-blockstore v1.0.4
	github.com/ipfs/go-ipfs-chunker v0.0.5
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-files v0.0.9
	github.com/ipfs/go-ipfs-routing v0.1.0
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-log/v2 v2.3.0
	github.com/ipfs/go-merkledag v0.4.1
	github.com/ipfs/go-unixfs v0.2.6
	github.com/ipld/go-car v0.3.2-0.20211001225732-32d0d9933823
	github.com/ipld/go-car/v2 v2.0.3-0.20210811121346-c514a30114d7
	github.com/ipld/go-ipld-prime v0.12.3
	github.com/ipsn/go-secp256k1 v0.0.0-20180726113642-9d62b9f0bc52
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/jedisct1/go-minisign v0.0.0-20190909160543-45766022959e
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/libp2p/go-libp2p v0.15.0
	github.com/libp2p/go-libp2p-core v0.11.0
	github.com/libp2p/go-libp2p-kad-dht v0.13.0
	github.com/libp2p/go-libp2p-peerstore v0.3.0
	github.com/libp2p/go-libp2p-pubsub v0.5.6
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-libp2p-routing-helpers v0.2.3
	github.com/mattn/go-colorable v0.1.9
	github.com/mattn/go-isatty v0.0.14
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-multiaddr v0.4.1
	github.com/multiformats/go-multibase v0.0.3
	github.com/multiformats/go-multihash v0.0.16
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/olekukonko/tablewriter v0.0.5
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/prometheus/tsdb v0.7.1
	github.com/rjeczalik/notify v0.9.1
	github.com/rs/cors v1.7.0
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible
	github.com/status-im/keycard-go v0.0.0-20190316090335-8537d3370df4
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/whyrusleeping/cbor-gen v0.0.0-20210713220151-be142a5ae1a8
	github.com/xlab/c-for-go v0.0.0-20201223145653-3ba5db515dcb // indirect
	github.com/xsleonard/go-merkle v1.1.0
	github.com/zondax/hid v0.9.0
	go.opencensus.io v0.23.0
	go.uber.org/dig v1.13.0 // indirect
	go.uber.org/fx v1.9.0
	go.uber.org/multierr v1.7.0
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210917161153-d61c044b1678
	golang.org/x/text v0.3.7
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gopkg.in/urfave/cli.v1 v1.20.0
	gotest.tools v2.2.0+incompatible

)

replace github.com/libp2p/go-libp2p-core v0.11.0 => github.com/libp2p/go-libp2p-core v0.9.0

replace github.com/filecoin-project/lotus => ./extern/storage-lib

replace github.com/filecoin-project/filecoin-ffi => ./extern/storage-lib/extern/filecoin-ffi

//replace github.com/filecoin-project/test-vectors => ./extern/storage-lib/extern/test-vectors

//replace github.com/filecoin-project/sector-storage => ./extern/storage-lib/extern/sector-storage

//replace github.com/filecoin-project/storage-sealing => ./extern/storage-lib/extern/storage-sealing
