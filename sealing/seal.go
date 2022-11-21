package sealing

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/borcontracts/w3fsStorageManager"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/lotus/api"

	//"github.com/ethereum/go-ethereum/extern/sector-storage/ffiwrapper"
	//"github.com/ethereum/go-ethereum/extern/sector-storage/ffiwrapper/basicfs"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	saproof2 "github.com/filecoin-project/specs-actors/v2/actors/runtime/proof"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
	"math/big"
)

const (
	DomainSeparationTag_TicketProduction int64 = 1 + iota
	DomainSeparationTag_ElectionProofProduction
	DomainSeparationTag_WinningPoStChallengeSeed
	DomainSeparationTag_WindowedPoStChallengeSeed
	DomainSeparationTag_SealRandomness
	DomainSeparationTag_InteractiveSealChallengeSeed
	DomainSeparationTag_WindowedPoStDeadlineAssignment
	DomainSeparationTag_MarketDealCronSeed
	DomainSeparationTag_PoStChainCommit
)


/*var DefDir = "~/.lotus-bench"


type SealSeed struct {
	Value abi.InteractiveSealRandomness
	Epoch abi.ChainEpoch
}

type PoRepResult struct {
	InitResult     bool
	AddPieceResult bool
	Pc1            bool
	Pc2            bool
	C1             bool
	C2             bool
}

type Sealing struct {
	DataDir      string
	SectorNumber uint64
	Mid          abi.ActorID
	SectorSize   abi.SectorSize
	Spt          abi.RegisteredSealProof
	Sb           *ffiwrapper.Sealer
	Piece        abi.PieceInfo
	Pc1o         storage.PreCommit1Out
	Ticket       abi.SealRandomness
	Seed         SealSeed
	Cids         storage.SectorCids
	C1o          storage.Commit1Out
	Proof        storage.Proof
	RoRepResult  PoRepResult
}*/

/*func NewFfiSealer(dataDir string) (*ffiwrapper.Sealer, error) {
	pathDir, _ := homedir.Expand(dataDir)
	if _, err := os.Stat(pathDir); err != nil {
		os.MkdirAll(pathDir, 0775)
	}
	return ffiwrapper.New(&basicfs.Provider{
		Root: pathDir,
	})
}
*/
/*func New(sectorNumber uint64, mid uint64, sectorSizeKb string, dataDir string, blockNumber uint64, headerHash []byte) *Sealing {
	pathDir, _ := homedir.Expand(dataDir)
	if _, err := os.Stat(pathDir); err != nil {
		os.MkdirAll(pathDir, 0775)
	}
	minerAddr := fmt.Sprintf("%s%s", "t0", strconv.Itoa(int(mid)))

	sectorSizeInt, _ := units.RAMInBytes(sectorSizeKb)
	sectorSize := abi.SectorSize(sectorSizeInt)
	spt, _ := SealProofTypeFromSectorSize(sectorSize, network.Version0)

	ticket, _ := GetTicket(big.NewInt(int64(blockNumber)-900), headerHash, minerAddr, DomainSeparationTag_SealRandomness)
	seedValue, _ := GetTicket(big.NewInt(int64(blockNumber)+150), headerHash, minerAddr, DomainSeparationTag_InteractiveSealChallengeSeed)

	ffiSeal, _ := NewFfiSealer(dataDir)

	return &Sealing{
		SectorNumber: sectorNumber,
		Mid:          abi.ActorID(mid),
		Spt:          spt,
		Sb:           ffiSeal,
		Ticket:       ticket,
		SectorSize:   sectorSize,
		RoRepResult:  PoRepResult{},
		Seed: SealSeed{
			Epoch: abi.ChainEpoch(blockNumber + 150),
			Value: seedValue,
		},
	}
}*/

/*func (s *Sealing) AddPiece() {
	sid := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  s.Mid,
			Number: abi.SectorNumber(s.SectorNumber),
		},
		ProofType: s.Spt,
	}
	r := rand.New(rand.NewSource(100 + int64(s.SectorNumber)))
	pi, err := s.Sb.AddPiece(context.TODO(), sid, nil, abi.PaddedPieceSize(s.SectorSize).Unpadded(), r)
	if err != nil {
		panic("addpiece is error")
	} else {
		s.Piece = pi
		s.RoRepResult.AddPieceResult = true
	}
}*/

/*func (s *Sealing) SealPreCommit1() {
	if s.RoRepResult.AddPieceResult == false {
		panic("SealPreCommit1 not do , AddPiece is error")
	}
	sid := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  s.Mid,
			Number: abi.SectorNumber(s.SectorNumber),
		},
		ProofType: s.Spt,
	}
	pc1o, err := s.Sb.SealPreCommit1(context.TODO(), sid, s.Ticket, []abi.PieceInfo{s.Piece})
	if err == nil {
		s.Pc1o = pc1o
		s.RoRepResult.Pc1 = true
	}
}*/

/*func (s *Sealing) SealPreCommit2() {
	sid := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  s.Mid,
			Number: abi.SectorNumber(s.SectorNumber),
		},
		ProofType: s.Spt,
	}
	cids, err := s.Sb.SealPreCommit2(context.TODO(), sid, s.Pc1o)
	if err == nil {
		s.Cids = cids
		s.RoRepResult.Pc2 = true
	}
}*/

/*func (s *Sealing) SealCommit1() {
	sid := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  s.Mid,
			Number: abi.SectorNumber(s.SectorNumber),
		},
		ProofType: s.Spt,
	}
	commit1, err := s.Sb.SealCommit1(context.TODO(), sid, s.Ticket, s.Seed.Value, []abi.PieceInfo{s.Piece}, s.Cids)
	if err == nil {
		s.C1o = commit1
		s.RoRepResult.C1 = true
	}
}*/

/*func (s *Sealing) SealCommit2() {
	sid := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  s.Mid,
			Number: abi.SectorNumber(s.SectorNumber),
		},
		ProofType: s.Spt,
	}
	proof, err := s.Sb.SealCommit2(context.TODO(), sid, s.C1o)
	if err == nil {
		s.Proof = proof
		s.RoRepResult.C2 = true
	}
}*/

/*func (s *Sealing) VerifySeal() (bool, error) {
	if !s.RoRepResult.AddPieceResult || !s.RoRepResult.Pc1 || !s.RoRepResult.Pc2 || !s.RoRepResult.C1 || !s.RoRepResult.C2 {
		return false, errors.New("addpiece or pc1 or pc2 or c1 or c2 is wrong")
	}
	sid := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  s.Mid,
			Number: abi.SectorNumber(s.SectorNumber),
		},
		ProofType: s.Spt,
	}
	svi := saproof2.SealVerifyInfo{
		SectorID:              abi.SectorID{Miner: s.Mid, Number: abi.SectorNumber(s.SectorNumber)},
		SealedCID:             s.Cids.Sealed,
		SealProof:             sid.ProofType,
		Proof:                 s.Proof,
		DealIDs:               nil,
		Randomness:            s.Ticket,
		InteractiveRandomness: s.Seed.Value,
		UnsealedCID:           s.Cids.Unsealed,
	}
	ok, err := ffiwrapper.ProofVerifier.VerifySeal(svi)
	return ok, err
}*/

/*func GenerateWinningPoSt(cvotes []borcontracts.Cvote, challenge []byte, mid uint64, dir string) ([]saproof2.PoStProof, error){
	//ffiwrapper, err := NewFfiSealer(fmt.Sprintf("%s-%s",DefDir,strconv.Itoa(int(mid))))
	ffiwrapper, err := NewFfiSealer(dir)
	var sectorInfos []saproof2.SectorInfo
	for _ , v := range cvotes {
		vcid,_ := cid.Cast(v.SealedCID)
		sectorInfos = append(sectorInfos, saproof2.SectorInfo{
			SectorNumber: abi.SectorNumber(v.SectorInx),
			SealProof: abi.RegisteredSealProof(v.SealProofType),
			SealedCID: vcid,
		})
	}
	st, err := ffiwrapper.GenerateWinningPoSt(context.TODO(), abi.ActorID(mid), sectorInfos, challenge)
	return st,err
}*/

func GenerateWinningPoSt2(cvotes []w3fsStorageManager.Cvote, challenge []byte, storageApi *api.StorageMiner) ([]saproof2.PoStProof, error){
	var sectorInfos []saproof2.SectorInfo
	for _ , v := range cvotes {
		vcid,_ := cid.Cast(v.SealedCID)
		sectorInfos = append(sectorInfos, saproof2.SectorInfo{
			SectorNumber: abi.SectorNumber(v.SectorInx),
			SealProof: abi.RegisteredSealProof(v.SealProofType),
			SealedCID: vcid,
		})
	}
	proof, err := (*storageApi).ComputeProof(context.Background(), sectorInfos, challenge)
	return proof, err
}

func VerifyWinningPoSt(mid uint64, cvotes []w3fsStorageManager.Cvote, winningPostProofs []w3fsStorageManager.WinningPostProof, randomness []byte) (bool, error) {
	var sectorInfos []saproof2.SectorInfo
	var postProofs []saproof2.PoStProof
	for _, vote := range cvotes {
		cid, _ := cid.Cast(vote.SealedCID)
		sectorInfos = append(sectorInfos, saproof2.SectorInfo{
			SectorNumber: abi.SectorNumber(vote.SectorInx),
			SealProof:    abi.RegisteredSealProof(vote.SealProofType),
			SealedCID:    cid,
		})
	}
	for _, proof := range winningPostProofs {
		postProofs = append(postProofs, saproof2.PoStProof{
			PoStProof:  abi.RegisteredPoStProof(proof.PoStProof),
			ProofBytes: proof.ProofBytes,
		})
	}
	pvi := saproof2.WinningPoStVerifyInfo{
		Randomness:        abi.PoStRandomness(randomness[:]),
		ChallengedSectors: sectorInfos,
		Prover:            abi.ActorID(mid),
		Proofs:            postProofs,
	}
	return ffiwrapper.ProofVerifier.VerifyWinningPoSt(context.TODO(), pvi)
}




func SectorSizeFromSealProofType(sealType abi.RegisteredSealProof) uint64 {
	switch sealType {
	case abi.RegisteredSealProof_StackedDrg2KiBV1 :
		return 2 << 10
	case abi.RegisteredSealProof_StackedDrg8MiBV1 :
		return 8 << 20
	case abi.RegisteredSealProof_StackedDrg512MiBV1 :
		return 512 << 20
	case abi.RegisteredSealProof_StackedDrg32GiBV1 :
		return 32 << 30
	case abi.RegisteredSealProof_StackedDrg64GiBV1 :
		return  64 << 30
	}
	return 0
}

func GetSealProofType(sealType uint64) string {
	if sealType == 0 || sealType == 5 {
		return "2KiB"
	}
	if sealType == 1 || sealType == 6 {
		return "8MiB"
	}
	if sealType == 2 || sealType == 7 {
		return "512MiB"
	}
	if sealType == 3 || sealType == 8 {
		return "32GiB"
	}
	if sealType == 4 || sealType == 9 {
		return "64GiB"
	}
	return ""
}


func SealProofTypeFromSectorSize(ssize abi.SectorSize, nv network.Version) (abi.RegisteredSealProof, error) {
	switch {
	case nv < network.Version7:
		switch ssize {
		case 2 << 10:
			return abi.RegisteredSealProof_StackedDrg2KiBV1, nil
		case 8 << 20:
			return abi.RegisteredSealProof_StackedDrg8MiBV1, nil
		case 512 << 20:
			return abi.RegisteredSealProof_StackedDrg512MiBV1, nil
		case 32 << 30:
			return abi.RegisteredSealProof_StackedDrg32GiBV1, nil
		case 64 << 30:
			return abi.RegisteredSealProof_StackedDrg64GiBV1, nil
		default:
			return 0, xerrors.Errorf("unsupported sector size for miner: %v", ssize)
		}
	case nv >= network.Version7:
		switch ssize {
		case 2 << 10:
			return abi.RegisteredSealProof_StackedDrg2KiBV1_1, nil
		case 8 << 20:
			return abi.RegisteredSealProof_StackedDrg8MiBV1_1, nil
		case 512 << 20:
			return abi.RegisteredSealProof_StackedDrg512MiBV1_1, nil
		case 32 << 30:
			return abi.RegisteredSealProof_StackedDrg32GiBV1_1, nil
		case 64 << 30:
			return abi.RegisteredSealProof_StackedDrg64GiBV1_1, nil
		default:
			return 0, xerrors.Errorf("unsupported sector size for miner: %v", ssize)
		}
	}

	return 0, xerrors.Errorf("unsupported network version")
}

func GetTicket(blockNumber *big.Int, rbase []byte, addressStr string, randtype int64) ([]byte, error) {
	fromString, _ := address.NewFromString(addressStr)
	buf := new(bytes.Buffer)
	err := fromString.MarshalCBOR(buf)
	if err != nil {
		panic("marchalCBOR is error")
	}
	//tipset := new(big.Int).Add(blockNumber)
	challenge, err := DrawRandomness(rbase, blockNumber.Int64(), buf.Bytes(), randtype)
	return challenge, err
}

func DrawRandomness(rbase []byte, round int64, entropy []byte, randtype int64) ([]byte, error) {
	h := blake2b.New256()
	if err := binary.Write(h, binary.BigEndian, randtype); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	VRFDigest := blake2b.Sum256(rbase)
	_, err := h.Write(VRFDigest[:])
	if err != nil {
		return nil, xerrors.Errorf("hashing VRFDigest: %w", err)
	}
	if err := binary.Write(h, binary.BigEndian, round); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	_, err = h.Write(entropy)
	if err != nil {
		return nil, xerrors.Errorf("hashing entropy: %w", err)
	}
	return h.Sum(nil), nil
}



