package modulesext

import (
	"context"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	"golang.org/x/xerrors"

	"github.com/multiformats/go-multiaddr"

	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/blockstore"
	"github.com/filecoin-project/lotus/markets/retrievaladapter"
	"github.com/filecoin-project/lotus/markets/storageadapter"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
)

func IpfsStorageBlockstoreAccessor(ipfsBlockstore dtypes.ClientBlockstore) storagemarket.BlockstoreAccessor {
	return storageadapter.NewFixedBlockstoreAccessor(bstore.Blockstore(ipfsBlockstore))
}

func IpfsRetrievalBlockstoreAccessor(ipfsBlockstore dtypes.ClientBlockstore) retrievalmarket.BlockstoreAccessor {
	return retrievaladapter.NewFixedBlockstoreAccessor(bstore.Blockstore(ipfsBlockstore))
}

// IpfsClientBlockstore returns a ClientBlockstore implementation backed by an IPFS node.
// If ipfsMaddr is empty, a local IPFS node is assumed considering IPFS_PATH configuration.
// If ipfsMaddr is not empty, it will connect to the remote IPFS node with the provided multiaddress.
// The flag useForRetrieval indicates if the IPFS node will also be used for storing retrieving deals.
func IpfsClientBlockstore(ctx context.Context, ipfsMaddr string, onlineMode bool, localStore dtypes.ClientImportMgr) (dtypes.ClientBlockstore, error) {
	var err error
	var ipfsbs blockstore.BasicBlockstore
	if ipfsMaddr != "" {
		var ma multiaddr.Multiaddr
		ma, err = multiaddr.NewMultiaddr(ipfsMaddr)
		if err != nil {
			return nil, xerrors.Errorf("parsing ipfs multiaddr: %w", err)
		}
		ipfsbs, err = blockstore.NewRemoteIPFSBlockstore(ctx, ma, onlineMode)
	} else {
		ipfsbs, err = blockstore.NewLocalIPFSBlockstore(ctx, onlineMode)
	}
	if err != nil {
		return nil, xerrors.Errorf("constructing ipfs blockstore: %w", err)
	}
	return blockstore.WrapIDStore(ipfsbs), nil
}
