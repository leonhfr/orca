package search

import "github.com/leonhfr/orca/chess"

// scoreMoves scores the moves.
func scoreMoves(pos *chess.Position, moves []chess.Move, best chess.Move, killers [2]chess.Move) {
	for i, move := range moves {
		moves[i] = move.WithScore(score(pos, move, best, killers))
	}
}

// quickScoreMoves quickly scores the moves.
func quickScoreMoves(pos *chess.Position, moves []chess.Move) {
	for i, move := range moves {
		moves[i] = move.WithScore(quickScore(pos, move))
	}
}

// scoreLoudMoves scores the loud moves.
func scoreLoudMoves(pos *chess.Position, moves []chess.Move) {
	for i, move := range moves {
		moves[i] = move.WithScore(rankSEE(pos, move))
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

// score ranks the move.
//
// Rank is computed according to the following order:
//
//	score              move
//	 500               best move
//	 490               queen promotion
//	 480               knight promotion
//	 470               king side castle
//	 460               queen side castle
//	 300 + [-100:100]  capture ordered by mvv-lva
//	 150               killer moves
//	 100               quiet moves
//	   0               bishop and rook promotions
func score(pos *chess.Position, m, best chess.Move, killers [2]chess.Move) uint32 {
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
		return uint32(rankCapture + see(pos, m))
	case killers[0] == m || killers[1] == m:
		return rankKiller
	default:
		return rankQuiet
	}
}

// quickScore scores the move without the best and killer moves.
func quickScore(pos *chess.Position, m chess.Move) uint32 {
	switch {
	case m.HasTag(chess.KingSideCastle):
		return rankKingSideCastle
	case m.HasTag(chess.QueenSideCastle):
		return rankQueenSideCastle
	case m.HasTag(chess.Promotion):
		return promoRank[m.Promo()]
	case m.HasTag(chess.Capture):
		return uint32(rankCapture + see(pos, m))
	default:
		return rankQuiet
	}
}

// rankSEE ranks the move by SEE.
func rankSEE(pos *chess.Position, m chess.Move) uint32 {
	return uint32(rankCapture + see(pos, m))
}

const (
	rankBestMove        = 500
	rankKingSideCastle  = 470
	rankQueenSideCastle = 460
	rankCapture         = 300
	rankKiller          = 150
	rankQuiet           = 100
)

var promoRank = [13]uint32{0, 0, 480, 480, 0, 0, 0, 0, 490, 490, 0, 0, 0} // promo rank indexed by piece
