package search

import "github.com/leonhfr/orca/chess"

// evaluate returns the score of a position.
func (si *searchInfo) evaluate(pos *chess.Position) int32 {
	player := pos.Turn()
	knights, bishops, rooks, queens := pos.CountPieces()
	phase := int32(knights + bishops + 2*rooks + 4*queens)
	if phase > 24 {
		phase = 24 // in case of early promotion
	}

	mg, eg := si.evaluatePawns(pos)

	pos.PieceMap(func(p chess.Piece, sq chess.Square, mobility int) {
		mgValue := pestoMGPieceSquareTable[p][sq]
		egValue := pestoEGPieceSquareTable[p][sq]

		pt := p.Type()
		mgValue += mobilityTermsMG[pt][mobility]
		egValue += mobilityTermsEG[pt][mobility]

		if p.Color() == chess.White {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}
	})

	pos.KingMap(func(p chess.Piece, sq chess.Square, shieldDefects int) {
		mgValue := pestoMGPieceSquareTable[p][sq]
		egValue := pestoEGPieceSquareTable[p][sq]

		mgValue += shieldDefectsPenaltyMG[shieldDefects]
		egValue += shieldDefectsPenaltyEG[shieldDefects]

		if p.Color() == chess.White {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}
	})

	if player == chess.Black {
		mg, eg = -mg, -eg
	}
	mg += tempo

	return (phase*mg + (24-phase)*eg) / 24
}

// evaluatePawns evaluate the pawn structure.
//
// Always returns the evaluation from White's point of view.
func (si *searchInfo) evaluatePawns(pos *chess.Position) (int32, int32) {
	pawnHash := pos.PawnHash()

	if entry, inCache := si.pawnTable.get(pawnHash); inCache {
		return entry.mg, entry.eg
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

	si.pawnTable.set(pawnHash, pawnEntry{
		hash: pawnHash,
		mg:   mg,
		eg:   eg,
	})

	return mg, eg
}

// tempo is a bonus for the player having the right to move.
// Only applied for middle game.
const tempo = 6

var (
	pestoMGPieceSquareTable = [12][64]int32{}                            // Middle game piece square table. Includes material advantage. Indexed by square and piece.
	pestoEGPieceSquareTable = [12][64]int32{}                            // End game piece square table. Includes material advantage. Indexed by square and piece.
	doubledPenaltyMG        = [8]int32{-10, -6, -6, -6, -6, -6, -6, -10} // Middle game penalty for double pawns. Indexed by file.
	doubledPenaltyEG        = [8]int32{-5, -3, -3, -3, -3, -3, -3, -5}   // End game penalty for double pawns. Indexed by file.
	isolaniPenaltyMG        = [8]int32{-10, -6, -6, -6, -6, -6, -6, -10} // Middle game penalty for isolated pawns. Indexed by file.
	isolaniPenaltyEG        = [8]int32{-5, -3, -3, -3, -3, -3, -3, -5}   // End game penalty for isolated pawns. Indexed by file.
	passedBonusMG           = [2][64]int32{}                             // Middle game bonus for passed pawns. Indexed by square and color.
	passedBonusEG           = [2][64]int32{}                             // End game bonus for passed pawns. Indexed by square and color.
	shieldDefectsPenaltyMG  = [4]int32{0, -10, -20, -30}                 // Middle game penalty for shield defects.
	shieldDefectsPenaltyEG  = [4]int32{0, 0, 0, 0}                       // End game penalty for shield defects.
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
