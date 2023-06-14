package search

import "github.com/leonhfr/orca/chess"

// lateMoveReduction determines the ply reduction for late moves.
func lateMoveReduction(validMoves int, inCheck bool, depth uint8, move chess.Move) uint8 {
	if validMoves <= 4 || inCheck || depth <= 3 {
		return 0
	}

	if !move.HasTag(chess.Quiet) {
		return 0
	}

	return 1
}
