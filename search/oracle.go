package search

import (
	"sort"

	"github.com/leonhfr/orca/chess"
)

// oracle orders the moves.
func oracle(moves []chess.Move, best chess.Move) {
	sort.Slice(moves, func(i, j int) bool {
		return rank(moves[i], best) > rank(moves[j], best)
	})
}

// loudOracle orders the moves only by MVV-LVA.
func loudOracle(moves []chess.Move) {
	sort.Slice(moves, func(i, j int) bool {
		return rankMvvLva(moves[i]) > rankMvvLva(moves[j])
	})
}

// rank ranks the move.
//
// Rank is computed according to the following order:
// 0. Best move, if any;
// 1. Queen promotions;
// 2. Knight promotions;
// 3. King side castle;
// 4. Queen side castle;
// 5. Captures according to the MVA-LVA ordering;
// 6. Quiet moves;
// 7. Rook promotions;
// 8. Bishop promotions.
func rank(m, best chess.Move) int {
	switch {
	case m == best:
		return rankBestMove
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

// rankMvvLva ranks the move by MVV-LVA.
func rankMvvLva(m chess.Move) int {
	return mvvRank[m.P2()] - lvaRank[m.P1()]
}

const (
	rankBestMove        = 100 // best move
	rankQueenSideCastle = 75  // queen side castle rank
	rankKingSideCastle  = 80  // king side castle rank
)

var (
	promoRank = [13]int{0, 0, 85, 85, -10, -10, -5, -5, 90, 90, 0, 0, 0}   // promo rank indexed by piece
	mvvRank   = [13]int{10, 10, 20, 20, 30, 30, 40, 40, 60, 60, 70, 70, 0} // mvv rank indexed by piece
	lvaRank   = [13]int{1, 1, 2, 2, 3, 3, 4, 4, 6, 6, 7, 7, 0}             // lva rank indexed by piece
)
