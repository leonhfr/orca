package search

import (
	"sort"

	"github.com/leonhfr/orca/chess"
)

// oracle orders the moves.
func oracle(moves []chess.Move) {
	sort.Slice(moves, func(i, j int) bool {
		return rank(moves[i]) > rank(moves[j])
	})
}

// rank ranks the move.
//
// Rank is computed according to the following order:
// 1. Queen promotions;
// 2. Knight promotions;
// 3. King side castle;
// 4. Queen side castle;
// 5. Captures according to the MVA-LVA ordering;
// 6. Quiet moves;
// 7. Rook promotions;
// 8. Bishop promotions.
func rank(m chess.Move) int {
	switch {
	case m.HasTag(chess.KingSideCastle):
		return rankKingSideCastle
	case m.HasTag(chess.QueenSideCastle):
		return rankQueenSideCastle
	case m.HasTag(chess.Promotion):
		return promoRank[m.Promo()]
	case m.HasTag(chess.Capture):
		return mvvRank[m.P2()] - lvaRank[m.P1()]
	default:
		return 0
	}
}

const (
	rankQueenSideCastle = 85 // queen side castle rank
	rankKingSideCastle  = 90 // king side castle rank
)

var (
	promoRank = [13]int{0, 0, 95, 95, -10, -10, -5, -5, 100, 100, 0, 0, 0} // promo rank indexed by piece
	mvvRank   = [13]int{10, 10, 20, 20, 30, 30, 40, 40, 60, 60, 70, 70, 0} // mvv rank indexed by piece
	lvaRank   = [13]int{1, 1, 2, 2, 3, 3, 4, 4, 6, 6, 7, 7, 0}             // lva rank indexed by piece
)
