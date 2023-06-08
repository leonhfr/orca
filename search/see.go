package search

import "github.com/leonhfr/orca/chess"

// see performs a static exchange evaluation on a move and returns the value.
func see(pos *chess.Position, m chess.Move) int {
	var depth int
	var gains [32]int
	gains[0] = values[m.P2().Type()]

	pos.StaticExchange(m, func(pt chess.PieceType) bool {
		depth++
		gains[depth] = values[pt] - gains[depth-1]
		return maxInt(-gains[depth-1], gains[depth]) < 0
	})

	for depth--; depth > 0; depth-- {
		gains[depth-1] = -maxInt(-gains[depth-1], gains[depth])
	}

	return gains[0]
}

// values for piece types.
var values = [7]int{10, 30, 35, 50, 100, 0, 0}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
