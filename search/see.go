package search

import (
	"math"

	"github.com/leonhfr/orca/chess"
)

// see performs a static exchange evaluation on a move and returns the value.
func see(pos *chess.Position, m chess.Move) int32 {
	var depth int
	var gains [32]int32
	gains[0] = values[m.P2().Type()]

	pos.StaticExchange(m, func(pt chess.PieceType) bool {
		depth++
		gains[depth] = values[pt] - gains[depth-1]
		return max32(-gains[depth-1], gains[depth]) < 0
	})

	for depth--; depth > 0; depth-- {
		gains[depth-1] = -max32(-gains[depth-1], gains[depth])
	}

	return gains[0]
}

// temp values for piece types.
var values = [6]int32{100, 325, 325, 500, 1000, math.MaxInt32}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
