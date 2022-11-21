package sealing

import (
	"testing"
)

func TestSeal(t *testing.T) {
	/*mid := 1000
	dataDir := fmt.Sprintf("~/.lotus-bench-%s", strconv.Itoa(int(mid)))
	sectorNumber := 0
	size := "2KiB"
	seal := New(uint64(sectorNumber), uint64(mid), size, dataDir, 0, []byte("0000000"))
	seal.AddPiece()
	seal.SealPreCommit1()
	seal.SealPreCommit2()
	seal.SealCommit1()
	seal.SealCommit2()
	ok, err := seal.VerifySeal()
	if !ok || err != nil {
		log.Fatal(err)
	}
	var docVotes []borcontracts.Cvote
	sectorSizeInt, _ := units.RAMInBytes(size)
	sealProofType, _ := SealProofTypeFromSectorSize(fabi.SectorSize(sectorSizeInt), network.Version0)
	docVotes = append(docVotes, borcontracts.Cvote{
		SectorInx:     uint64(sectorNumber),
		SealProofType: uint64(sealProofType),
		SealedCID:     seal.Cids.Sealed.Bytes(),
	})
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 1)
	rval := blake2b.Sum256(buf)
	poStProofRand, _ := GetTicket(big.NewInt(1), rval[:], fmt.Sprintf("%s%d", "t0", mid), DomainSeparationTag_WinningPoStChallengeSeed)
	poStProofs, err := GenerateWinningPoSt(docVotes, poStProofRand, 1000, "")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(poStProofs)
	}
	// check
	var winningPostProofs []borcontracts.WinningPostProof
	for _, proof := range poStProofs {
		winningPostProofs = append(winningPostProofs, borcontracts.WinningPostProof{
			PoStProof:  uint64(proof.PoStProof),
			ProofBytes: proof.ProofBytes,
		})
	}
	poStProofRand2, _ := GetTicket(big.NewInt(1), rval[:], fmt.Sprintf("%s%d", "t0", mid), DomainSeparationTag_WinningPoStChallengeSeed)
	winningData := borcontracts.WinningData{
		WinningPostProofs: winningPostProofs,
	}
	ok, verifyErr := VerifyWinningPoSt(uint64(mid), docVotes, winningData.WinningPostProofs, poStProofRand2)
	if !ok || verifyErr != nil {
		log.Fatal(verifyErr)
	}*/
}
