package search

import "github.com/leonhfr/orca/chess"

// scoreMoves scores the moves.
func scoreMoves(moves []chess.Move, best chess.Move, killers [2]chess.Move) {
	for i, move := range moves {
		moves[i] = move.WithScore(rank(move, best, killers))
	}
}

// scoreLoudMoves scores the loud moves.
func scoreLoudMoves(moves []chess.Move) {
	for i, move := range moves {
		moves[i] = move.WithScore(rankMvvLva(move))
	}
}

// nextOracle predicts the next best move.
func nextOracle(moves []chess.Move, start int) {
	for i := start + 1; i < len(moves); i++ {
		if moves[i].Score() > moves[start].Score() {
			moves[start], moves[i] = moves[i], moves[start]
		}
	}
}

// rank ranks the move.
//
// Rank is computed according to the following order:
//
//	rank           move
//	 500           best move
//	 490           queen promotion
//	 480           knight promotion
//	 470           king side castle
//	 460           queen side castle
//	 300 + [0:70]  capture ordered by mvv-lva
//	 200           killer moves
//	 100           quiet moves
//	   0           bishop and rook promotions
func rank(m, best chess.Move, killers [2]chess.Move) uint32 {
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
		return rankCapture + mvvRank[m.P2()] - lvaRank[m.P1()]
	case killers[0] == m || killers[1] == m:
		return rankKiller
	default:
		return rankQuiet
	}
}

// rankMvvLva ranks the move by MVV-LVA.
func rankMvvLva(m chess.Move) uint32 {
	return mvvRank[m.P2()] - lvaRank[m.P1()]
}

const (
	rankBestMove        = 500
	rankKingSideCastle  = 470
	rankQueenSideCastle = 460
	rankCapture         = 300
	rankKiller          = 200
	rankQuiet           = 100
)

var (
	promoRank = [13]uint32{0, 0, 480, 480, 0, 0, 0, 0, 490, 490, 0, 0, 0}      // promo rank indexed by piece
	mvvRank   = [13]uint32{10, 10, 20, 20, 30, 30, 40, 40, 60, 60, 70, 70, 10} // mvv rank indexed by piece
	lvaRank   = [13]uint32{1, 1, 2, 2, 3, 3, 4, 4, 6, 6, 7, 7, 0}              // lva rank indexed by piece
)
