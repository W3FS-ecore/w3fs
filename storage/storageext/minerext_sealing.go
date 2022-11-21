package storageext

import (
	"context"

	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/specs-storage/storage"

	sealingext "github.com/ethereum/go-ethereum/storage/storage-sealing-ext"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/extern/storage-sealing/sealiface"
	"github.com/filecoin-project/lotus/storage/sectorblocks"
)

// TODO: refactor this to be direct somehow

func (m *MinerExt) Address() address.Address {
	return m.sealing.Address()
}

func (m *MinerExt) StartPackingSector(sectorNum abi.SectorNumber) error {
	return m.sealing.StartPacking(sectorNum)
}

func (m *MinerExt) ListSectors() ([]sealingext.SectorInfo, error) {
	return m.sealing.ListSectors()
}

func (m *MinerExt) PledgeSector(ctx context.Context) (storage.SectorRef, error) {
	return m.sealing.PledgeSector(ctx)
}

func (m *MinerExt) ForceSectorState(ctx context.Context, id abi.SectorNumber, state sealingext.SectorState) error {
	return m.sealing.ForceSectorState(ctx, id, state)
}

func (m *MinerExt) RemoveSector(ctx context.Context, id abi.SectorNumber) error {
	return m.sealing.Remove(ctx, id)
}

func (m *MinerExt) TerminateSector(ctx context.Context, id abi.SectorNumber) error {
	return m.sealing.Terminate(ctx, id)
}

func (m *MinerExt) TerminateFlush(ctx context.Context) (*cid.Cid, error) {
	return m.sealing.TerminateFlush(ctx)
}

func (m *MinerExt) TerminatePending(ctx context.Context) ([]abi.SectorID, error) {
	return m.sealing.TerminatePending(ctx)
}

func (m *MinerExt) SectorPreCommitFlush(ctx context.Context) ([]sealiface.PreCommitBatchRes, error) {
	return m.sealing.SectorPreCommitFlush(ctx)
}

func (m *MinerExt) SectorPreCommitPending(ctx context.Context) ([]abi.SectorID, error) {
	return m.sealing.SectorPreCommitPending(ctx)
}

func (m *MinerExt) CommitFlush(ctx context.Context) ([]sealiface.CommitBatchRes, error) {
	return m.sealing.CommitFlush(ctx)
}

func (m *MinerExt) CommitPending(ctx context.Context) ([]abi.SectorID, error) {
	return m.sealing.CommitPending(ctx)
}

func (m *MinerExt) MarkForUpgrade(id abi.SectorNumber) error {
	return m.sealing.MarkForUpgrade(id)
}

func (m *MinerExt) IsMarkedForUpgrade(id abi.SectorNumber) bool {
	return m.sealing.IsMarkedForUpgrade(id)
}

func (m *MinerExt) SectorAddPieceToAny(ctx context.Context, size abi.UnpaddedPieceSize, r storage.Data, d api.PieceDealInfo) (api.SectorOffset, error) {
	return m.sealing.SectorAddPieceToAny(ctx, size, r, d)
}

func (m *MinerExt) SectorsStatus(ctx context.Context, sid abi.SectorNumber, showOnChainInfo bool) (api.SectorInfo, error) {
	if showOnChainInfo {
		return api.SectorInfo{}, xerrors.Errorf("on-chain info not supported")
	}

	info, err := m.sealing.GetSectorInfo(sid)
	if err != nil {
		return api.SectorInfo{}, err
	}

	deals := make([]abi.DealID, len(info.Pieces))
	pieces := make([]api.SectorPiece, len(info.Pieces))
	for i, piece := range info.Pieces {
		pieces[i].Piece = piece.Piece
		if piece.DealInfo == nil {
			continue
		}

		pdi := *piece.DealInfo // copy
		pieces[i].DealInfo = &pdi

		deals[i] = piece.DealInfo.DealID
	}

	log := make([]api.SectorLog, len(info.Log))
	for i, l := range info.Log {
		log[i] = api.SectorLog{
			Kind:      l.Kind,
			Timestamp: l.Timestamp,
			Trace:     l.Trace,
			Message:   l.Message,
		}
	}

	sInfo := api.SectorInfo{
		SectorID: sid,
		State:    api.SectorState(info.State),
		CommD:    info.CommD,
		CommR:    info.CommR,
		Proof:    info.Proof,
		Deals:    deals,
		Pieces:   pieces,
		Ticket: api.SealTicket{
			Value: info.TicketValue,
			Epoch: info.TicketEpoch,
		},
		Seed: api.SealSeed{
			Value: info.SeedValue,
			Epoch: info.SeedEpoch,
		},
		PreCommitMsg: info.PreCommitMessage,
		CommitMsg:    info.CommitMessage,
		Retries:      info.InvalidProofs,
		ToUpgrade:    m.IsMarkedForUpgrade(sid),

		LastErr: info.LastErr,
		Log:     log,
		// on chain info
		SealProof:          info.SectorType,
		Activation:         0,
		Expiration:         0,
		DealWeight:         big.Zero(),
		VerifiedDealWeight: big.Zero(),
		InitialPledge:      big.Zero(),
		OnTime:             0,
		Early:              0,
	}

	return sInfo, nil
}

var _ sectorblocks.SectorBuilder = &MinerExt{}
