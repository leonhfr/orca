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
func evaluate(pos *chess.Position) int {
	var mg, eg, phase int

	pos.PieceMap(func(p chess.Piece, sq chess.Square) {
		mgValue := pestoMGPieceTables[p][sq]
		egValue := pestoEGPieceTables[p][sq]

		if p.Color() == pos.Turn() {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}

		phase += pestoGamePhaseInc[p.Type()]
	})

	if phase > 24 {
		phase = 24 // in case of early promotion
	}

	return (phase*mg + (24-phase)*eg) / 24
}

// incMateDistance increases the distance to the mate by a count of one.
//
// In case of a positive score, it is decreased by 1.
// In case of a negative score, it is increased by 1.
func incMateDistance(score int) int {
	sign := 1
	if score < 0 {
		sign = -1
	}
	delta := mate - sign*score
	if delta <= maxPkgDepth {
		return score - sign
	}
	return score
}