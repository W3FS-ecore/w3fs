package sealingext

import (
	"context"
	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"

	"github.com/filecoin-project/go-state-types/abi"
)

// `curH`-`ts.Height` = `confidence`
type HeightHandler func(ctx context.Context, tok sealing.TipSetToken, curH abi.ChainEpoch) error
type RevertHandler func(ctx context.Context, tok sealing.TipSetToken) error

type Events interface {
	ChainAt(hnd HeightHandler, rev RevertHandler, confidence int, h abi.ChainEpoch) error
}
