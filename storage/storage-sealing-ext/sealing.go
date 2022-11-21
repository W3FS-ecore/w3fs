package sealingext

import (
	"context"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	statemachine "github.com/filecoin-project/go-statemachine"
	"github.com/filecoin-project/specs-storage/storage"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	sectorstorage "github.com/filecoin-project/lotus/extern/sector-storage"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	"github.com/filecoin-project/lotus/extern/storage-sealing/sealiface"
	"github.com/filecoin-project/lotus/node/config"
)

const SectorStorePrefix = "/sectors"

var log = logging.Logger("sectors")

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/api.go -package=mocks . SealingAPI

type Sealing struct {
	Api      sealing.SealingAPI
	DealInfo *sealing.CurrentDealInfoManager

	feeCfg config.MinerFeeConfig
	events sealing.Events

	startupWait sync.WaitGroup

	maddr address.Address

	sealer  sectorstorage.SectorManager
	sectors *statemachine.StateGroup
	sc      sealing.SectorIDCounter
	verif   ffiwrapper.Verifier
	pcp     sealing.PreCommitPolicy

	inputLk        sync.Mutex
	openSectors    map[abi.SectorID]*openSector
	sectorTimers   map[abi.SectorID]*time.Timer
	pendingPieces  map[cid.Cid]*pendingPiece
	assignedPieces map[abi.SectorID][]cid.Cid
	creating       *abi.SectorNumber // used to prevent a race where we could create a new sector more than once

	upgradeLk sync.Mutex
	toUpgrade map[abi.SectorNumber]struct{}

	notifee SectorStateNotifee
	addrSel sealing.AddrSel

	stats SectorStats

	terminator  *sealing.TerminateBatcher
	precommiter *PreCommitBatcher
	commiter    *CommitBatcher

	getConfig sealing.GetSealingConfigFunc

	backend ethapi.Backend
}

type SectorStateNotifee func(before, after SectorInfo)

type openSector struct {
	used abi.UnpaddedPieceSize // change to bitfield/rle when AddPiece gains offset support to better fill sectors

	maybeAccept func(cid.Cid) error // called with inputLk
}

type pendingPiece struct {
	size abi.UnpaddedPieceSize
	deal api.PieceDealInfo

	data storage.Data

	assigned bool // assigned to a sector?
	accepted func(abi.SectorNumber, abi.UnpaddedPieceSize, error)
}

func New(mctx context.Context, api sealing.SealingAPI, fc config.MinerFeeConfig, events sealing.Events, maddr address.Address, ds datastore.Batching, sealer sectorstorage.SectorManager, sc sealing.SectorIDCounter, verif ffiwrapper.Verifier, prov ffiwrapper.Prover, pcp sealing.PreCommitPolicy, gc sealing.GetSealingConfigFunc, notifee SectorStateNotifee, as sealing.AddrSel, backend ethapi.Backend) *Sealing {
	s := &Sealing{
		Api:      api,
		DealInfo: &sealing.CurrentDealInfoManager{api},

		feeCfg: fc,
		events: events,

		maddr:  maddr,
		sealer: sealer,
		sc:     sc,
		verif:  verif,
		pcp:    pcp,

		openSectors:    map[abi.SectorID]*openSector{},
		sectorTimers:   map[abi.SectorID]*time.Timer{},
		pendingPieces:  map[cid.Cid]*pendingPiece{},
		assignedPieces: map[abi.SectorID][]cid.Cid{},
		toUpgrade:      map[abi.SectorNumber]struct{}{},

		notifee: notifee,
		addrSel: as,

		terminator:  sealing.NewTerminationBatcher(mctx, maddr, api, as, fc, gc),
		precommiter: NewPreCommitBatcher(mctx, maddr, api, as, fc, gc),
		commiter:    NewCommitBatcher(mctx, maddr, api, as, fc, gc, prov),

		getConfig: gc,

		stats: SectorStats{
			bySector: map[abi.SectorID]SectorState{},
			byState:  map[SectorState]int64{},
		},

		backend: backend,
	}
	s.startupWait.Add(1)

	s.sectors = statemachine.New(namespace.Wrap(ds, datastore.NewKey(SectorStorePrefix)), s, SectorInfo{})

	return s
}

func (m *Sealing) Run(ctx context.Context) error {
	if err := m.restartSectors(ctx); err != nil {
		log.Errorf("%+v", err)
		return xerrors.Errorf("failed load sector states: %w", err)
	}

	return nil
}

func (m *Sealing) Stop(ctx context.Context) error {
	if err := m.terminator.Stop(ctx); err != nil {
		return err
	}

	if err := m.sectors.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func (m *Sealing) Remove(ctx context.Context, sid abi.SectorNumber) error {
	m.startupWait.Wait()

	return m.sectors.Send(uint64(sid), SectorRemove{})
}

func (m *Sealing) Terminate(ctx context.Context, sid abi.SectorNumber) error {
	m.startupWait.Wait()

	return m.sectors.Send(uint64(sid), SectorTerminate{})
}

func (m *Sealing) TerminateFlush(ctx context.Context) (*cid.Cid, error) {
	return m.terminator.Flush(ctx)
}

func (m *Sealing) TerminatePending(ctx context.Context) ([]abi.SectorID, error) {
	return m.terminator.Pending(ctx)
}

func (m *Sealing) SectorPreCommitFlush(ctx context.Context) ([]sealiface.PreCommitBatchRes, error) {
	return m.precommiter.Flush(ctx)
}

func (m *Sealing) SectorPreCommitPending(ctx context.Context) ([]abi.SectorID, error) {
	return m.precommiter.Pending(ctx)
}

func (m *Sealing) CommitFlush(ctx context.Context) ([]sealiface.CommitBatchRes, error) {
	return m.commiter.Flush(ctx)
}

func (m *Sealing) CommitPending(ctx context.Context) ([]abi.SectorID, error) {
	return m.commiter.Pending(ctx)
}

func (m *Sealing) currentSealProof(ctx context.Context) (abi.RegisteredSealProof, error) {
	mi, err := m.Api.StateMinerInfo(ctx, m.maddr, nil)
	if err != nil {
		return 0, err
	}

	ver, err := m.Api.StateNetworkVersion(ctx, nil)
	if err != nil {
		return 0, err
	}

	return miner.PreferredSealProofTypeFromWindowPoStType(ver, mi.WindowPoStProofType)
}

func (m *Sealing) minerSector(spt abi.RegisteredSealProof, num abi.SectorNumber) storage.SectorRef {
	return storage.SectorRef{
		ID:        m.minerSectorID(num),
		ProofType: spt,
	}
}

func (m *Sealing) minerSectorID(num abi.SectorNumber) abi.SectorID {
	mid, err := address.IDFromAddress(m.maddr)
	if err != nil {
		panic(err)
	}

	return abi.SectorID{
		Number: num,
		Miner:  abi.ActorID(mid),
	}
}

func (m *Sealing) Address() address.Address {
	return m.maddr
}

func getDealPerSectorLimit(size abi.SectorSize) (int, error) {
	if size < 64<<30 {
		return 256, nil
	}
	return 512, nil
}
