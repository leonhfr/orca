package search

import "github.com/leonhfr/orca/chess"

// evaluate returns the score of a position.
func (si *searchInfo) evaluate(pos *chess.Position) int32 {
	player := pos.Turn()
	knights, bishops, rooks, queens := pos.CountPieces()
	phase := int32(knights + bishops + 2*rooks + 4*queens)
	phase = min(phase, 24) // in case of early promotion

	var mgMaterial [2]int32
	var egMaterial [2]int32

	pawnCount := pos.PawnCount()
	mgMaterial[chess.Black] = int32(pawnCount[chess.Black]) * pestoMGPieceValues[chess.Pawn]
	mgMaterial[chess.White] = int32(pawnCount[chess.White]) * pestoMGPieceValues[chess.Pawn]
	egMaterial[chess.Black] = int32(pawnCount[chess.Black]) * pestoEGPieceValues[chess.Pawn]
	egMaterial[chess.White] = int32(pawnCount[chess.White]) * pestoEGPieceValues[chess.Pawn]

	mg, eg := si.evaluatePawns(pos)
	fd := pos.FileData()

	pos.PieceMap(func(p chess.Piece, sq chess.Square, mobility int, properties chess.PieceProperty) {
		c := p.Color()
		pt := p.Type()

		mgMaterial[c] += pestoMGPieceValues[pt]
		egMaterial[c] += pestoEGPieceValues[pt]

		mgValue := pestoMGPieceSquareTable[p][sq]
		egValue := pestoEGPieceSquareTable[p][sq]

		mgValue += mobilityTermsMG[pt][mobility]
		egValue += mobilityTermsEG[pt][mobility]

		if properties.HasProperty(chess.Trapped) {
			mgValue += trappedPiecePenalty[pt]
			egValue += trappedPiecePenalty[pt]
		}

		if properties.HasProperty(chess.Lost) {
			mgValue += lostPiecePenalty[pt]
			egValue += lostPiecePenalty[pt]
		}

		if pt == chess.Rook {
			switch {
			case sq.Rank() == rookPenultimateRank[c]:
				mgValue += rookPenultimateRankBonus
				egValue += rookPenultimateRankBonus
			case fd.OnOpenFile(sq):
				mgValue += rookOpenFileBonus
				egValue += rookOpenFileBonus
			case fd.OnHalfOpenFile(sq, c.Other()):
				mgValue += rookHalfOpenFileBonus
				egValue += rookHalfOpenFileBonus
			}
		}

		if c == chess.White {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}
	})

	materialValue := [2]int32{
		taperedEval(mgMaterial[chess.Black], egMaterial[chess.Black], phase),
		taperedEval(mgMaterial[chess.White], egMaterial[chess.White], phase),
	}

	pos.KingMap(func(p chess.Piece, sq chess.Square, shieldDefects, openFiles, halfOpenFiles int) {
		c := p.Color()

		factor := materialValue[c.Other()] / initialMaterialValue

		mgValue := pestoMGPieceSquareTable[p][sq]
		egValue := pestoEGPieceSquareTable[p][sq]

		mgValue += shieldDefectsPenaltyMG[shieldDefects] * factor
		egValue += shieldDefectsPenaltyEG[shieldDefects] * factor

		mgValue += (int32(openFiles)*openFilePenaltyMG + int32(halfOpenFiles)*halfOpenFilePenaltyMG) * factor
		egValue += (int32(openFiles)*openFilePenaltyEG + int32(halfOpenFiles)*halfOpenFilePenaltyEG) * factor

		if c == chess.White {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}
	})

	mg += mgMaterial[chess.White] - mgMaterial[chess.Black]
	eg += egMaterial[chess.White] - egMaterial[chess.Black]

	if player == chess.Black {
		mg, eg = -mg, -eg
	}
	mg += tempo

	return taperedEval(mg, eg, phase)
}

// evaluatePawns evaluate the pawn structure.
//
// Always returns the evaluation from White's point of view.
func (si *searchInfo) evaluatePawns(pos *chess.Position) (int32, int32) {
	pawnHash := pos.PawnHash()

	if entry, inCache := si.pawnTable.get(pawnHash); inCache {
		return entry.mg(), entry.eg()
	}

	var mg, eg int32
	pos.PawnMap(func(p chess.Piece, sq chess.Square, properties chess.PawnProperty) {
		color, file := p.Color(), sq.File()

		mgValue := pestoMGPieceSquareTable[p][sq]
		egValue := pestoEGPieceSquareTable[p][sq]

		if properties.HasProperty(chess.Doubled) {
			mgValue += doubledPenaltyMG[file]
			egValue += doubledPenaltyEG[file]
		}

		if properties.HasProperty(chess.Isolani) {
			mgValue += isolaniPenaltyMG[file]
			egValue += isolaniPenaltyEG[file]
		}

		if properties.HasProperty(chess.Passed) {
			mgValue += passedBonusMG[color][sq]
			egValue += passedBonusEG[color][sq]
		}

		if color == chess.White {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}
	})

	si.pawnTable.set(pawnHash, mg, eg)

	return mg, eg
}

func taperedEval(mg, eg, phase int32) int32 {
	return (phase*mg + (24-phase)*eg) / 24
}

// tempo is a bonus for the player having the right to move.
// Only applied for middle game.
const tempo = 6

var (
	initialMaterialValue    int32
	pestoMGPieceValues      = [6]int32{82, 337, 365, 477, 1025, 0}       // Middle game piece material values. Indexed by piece type.
	pestoEGPieceValues      = [6]int32{94, 281, 297, 512, 936, 0}        // End game piece material values. Indexed by piece type.
	pestoMGPieceSquareTable = [12][64]int32{}                            // Middle game piece square table. Indexed by square and piece.
	pestoEGPieceSquareTable = [12][64]int32{}                            // End game piece square table. Indexed by square and piece.
	doubledPenaltyMG        = [8]int32{-10, -6, -6, -6, -6, -6, -6, -10} // Middle game penalty for double pawns. Indexed by file.
	doubledPenaltyEG        = [8]int32{-5, -3, -3, -3, -3, -3, -3, -5}   // End game penalty for double pawns. Indexed by file.
	isolaniPenaltyMG        = [8]int32{-10, -6, -6, -6, -6, -6, -6, -10} // Middle game penalty for isolated pawns. Indexed by file.
	isolaniPenaltyEG        = [8]int32{-5, -3, -3, -3, -3, -3, -3, -5}   // End game penalty for isolated pawns. Indexed by file.
	passedBonusMG           = [2][64]int32{}                             // Middle game bonus for passed pawns. Indexed by square and color.
	passedBonusEG           = [2][64]int32{}                             // End game bonus for passed pawns. Indexed by square and color.
	shieldDefectsPenaltyMG  = [4]int32{0, -20, -40, -60}                 // Middle game penalty for shield defects.
	shieldDefectsPenaltyEG  = [4]int32{0, 0, 0, 0}                       // End game penalty for shield defects.
	trappedPiecePenalty     = [6]int32{0, -20, -50, -40, 0, 0}           // Penalty for trapped pieces. Indexed by piece type.
	lostPiecePenalty        = [6]int32{0, 0, -150, 0, 0, 0}              // Penalty for lost pieces. Indexed by piece type.
)

const (
	openFilePenaltyMG     = -20
	openFilePenaltyEG     = 0
	halfOpenFilePenaltyMG = -10
	halfOpenFilePenaltyEG = 0
)

var rookPenultimateRank = [2]chess.Rank{chess.Rank2, chess.Rank7} // rookPenultimateRank indicates the rook penultimate rank. Indexed by color.

const (
	rookPenultimateRankBonus = 80 // Bonus for rooks on penultimate rank.
	rookOpenFileBonus        = 60 // Bonus for rooks on open files.
	rookHalfOpenFileBonus    = 40 // Bonus for rooks on half open files.
)

var (
	mobilityTermsMG = [6][]int32{
		{}, // pawns
		{
			-104,
			-45, -22, -8, 6, 11, 19, 20, 45,
		},
		{
			-99,
			-46, -16, -4, 6, 14, 17, 19, 19, 27, 26,
			52, 55, 83,
		},
		{
			-127,
			-56, -25, -12, -10, -12, -11, -4, 4, 9, 11,
			19, 19, 37, 97,
		},
		{
			-111,
			-253, -127, -46, -20, -9, -1, 2, 8, 10, 15,
			17, 20, 23, 22, 21, 24, 16, 13, 18, 25,
			38, 34, 28, 10, 7, -42, -23,
		},
	}
	mobilityTermsEG = [6][]int32{
		{}, // pawns
		{
			-139,
			-114, -37, 3, 15, 34, 38, 37, 17,
		},
		{
			-186,
			-124, -54, -14, 1, 20, 35, 39, 49, 48, 48,
			32, 47, 2,
		},
		{
			-148,
			-127, -85, -28, 2, 27, 42, 46, 52, 55, 64,
			68, 73, 60, 15,
		},
		{
			-273,
			-401, -228, -236, -173, -86, -35, -1, 8, 31, 37,
			55, 46, 57, 58, 64, 62, 65, 63, 48, 30,
			8, -12, -29, -44, -79, -30, -50,
		},
	}
)

// incMateDistance increases the distance to the mate by a count of one.
//
// In case of a positive score, it is decreased by 1.
// In case of a negative score, it is increased by 1.
func incMateDistance(score int32) int32 {
	var sign int32 = 1
	if score < 0 {
		sign = -1
	}
	delta := mate - sign*score
	if delta <= maxSearchDepth {
		return score - sign
	}
	return score
}
