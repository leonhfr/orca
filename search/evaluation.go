package search

import "github.com/leonhfr/orca/chess"

// evaluate returns the score of a position.
//
// The function implements the PeSTO (Piece-Square Tables Only)
// evaluation function by Ronald Friedrich.
//
// It performs a tapered evaluation to interpolate by current game stage
// between piece-square tables values for opening and endgame.
//
// Source: https://www.chessprogramming.org/PeSTO%27s_Evaluation_Function
func (si *searchInfo) evaluate(pos *chess.Position) int32 {
	player := pos.Turn()
	knights, bishops, rooks, queens := pos.CountPieces()
	phase := int32(knights + bishops + 2*rooks + 4*queens)
	if phase > 24 {
		phase = 24 // in case of early promotion
	}

	mg, eg := si.evaluatePawns(pos)
	if player == chess.Black {
		mg, eg = -mg, -eg
	}

	if phase <= 6 || pos.FullMoves() > 16 {
		pos.PieceMap(func(p chess.Piece, sq chess.Square) {
			mgValue := pestoMGPieceTables[p][sq]
			egValue := pestoEGPieceTables[p][sq]
			if p.Color() == player {
				mg += mgValue
				eg += egValue
			} else {
				mg -= mgValue
				eg -= egValue
			}
		})

		return (phase*mg + (24-phase)*eg) / 24
	}

	pos.UniquePieceMap(func(p chess.Piece, sq chess.Square) {
		mgValue := pestoMGPieceTables[p][sq]
		egValue := pestoEGPieceTables[p][sq]
		if p.Color() == player {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}
	})

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
	// TODO: evaluation
	pos.PawnMap(func(p chess.Piece, sq chess.Square, properties chess.PawnProperty) {
		mgValue := pestoMGPieceTables[p][sq]
		egValue := pestoEGPieceTables[p][sq]
		if p.Color() == chess.White {
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
