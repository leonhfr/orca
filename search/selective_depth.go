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

	return rLateMoveReduction
}

// shouldNullMovePrune determines whether whether the search function should apply null move pruning.
func shouldNullMovePrune(pos *chess.Position, inCheck bool, depth uint8) bool {
	if inCheck || depth <= rNullMovePruning {
		return false
	}

	return pos.CountOwnPieces() != 0
}

const (
	rLateMoveReduction = 1 // Depth reduction in plies for late move reduction.
	rNullMovePruning   = 2 // Depth reduction in plies for null move pruning.
)
