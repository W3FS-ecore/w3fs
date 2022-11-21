package modulesext

import (
	"context"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo"
)

// UniversalBlockstore returns a single universal blockstore that stores both
// chain data and state data. It can be backed by a blockstore directly
// (e.g. Badger), or by a Splitstore.
func UniversalBlockstore(ctx context.Context, r repo.LockedRepo) (dtypes.UniversalBlockstore, error) {
	bs, err := r.Blockstore(ctx, repo.UniversalBlockstore)
	if err != nil {
		return nil, err
	}
	//if c, ok := bs.(io.Closer); ok {
	//	lc.Append(fx.Hook{
	//		OnStop: func(_ context.Context) error {
	//			return c.Close()
	//		},
	//	})
	//}
	return bs, err
}
