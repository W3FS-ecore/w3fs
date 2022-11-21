package modulesext

import (
	"context"
	discoveryimpl "github.com/filecoin-project/go-fil-markets/discovery/impl"
	"github.com/filecoin-project/lotus/journal"
	"github.com/filecoin-project/lotus/journal/fsjournal"
	marketevents "github.com/filecoin-project/lotus/markets/loggers"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
)

func NewLocalDiscovery(ctx context.Context, ds dtypes.MetadataDS) (*discoveryimpl.Local, error) {
	local, err := discoveryimpl.NewLocal(namespace.Wrap(ds, datastore.NewKey("/deals/local")))
	if err != nil {
		return nil, err
	}
	local.OnReady(marketevents.ReadyLogger("discovery"))
	local.Start(ctx)
	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//		return local.Start(ctx)
	//	},
	//})
	return local, nil
}

func OpenFilesystemJournal(lr repo.LockedRepo, disabled journal.DisabledEvents) (journal.Journal, error) {
	jrnl, err := fsjournal.OpenFSJournal(lr, disabled)
	if err != nil {
		return nil, err
	}

	//lc.Append(fx.Hook{
	//	OnStop: func(_ context.Context) error { return jrnl.Close() },
	//})

	return jrnl, err
}
