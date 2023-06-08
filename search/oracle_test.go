package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestScoreMoves(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		best    chess.Move
		killers [2]chess.Move
	}{
		{
			"tags",
			"rnbq1knr/pPpp2pp/8/Pp6/7Q/8/8/R3K2R w KQ b6 0 1",
			chess.NoMove,
			[2]chess.Move{
				chess.Move(chess.H4) ^
					chess.Move(chess.H5)<<6 ^
					chess.Move(chess.WhiteQueen)<<12 ^
					chess.Move(chess.NoPiece)<<16 ^
					chess.Move(chess.NoPiece)<<20 ^
					chess.Move(chess.Quiet),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			checkData, _ := pos.InCheck()
			original := pos.PseudoMoves(checkData)

			moves := pos.PseudoMoves(checkData)
			scoreMoves(pos, moves, tt.best, tt.killers)

			for i, move := range moves {
				assert.Equal(t, score(pos, original[i], tt.best, tt.killers), move.Score())
			}
		})
	}
}

func TestOrderMoves(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		best    chess.Move
		killers [2]chess.Move
		want    []string
	}{
		{
			"promotions",
			"7k/P7/8/8/8/8/8/K7 w - - 0 1",
			chess.NoMove,
			[2]chess.Move{},
			[]string{
				"a7a8q", "a7a8n", "a1b2", "a1b1", "a1a2",
				"a7a8b", "a7a8r",
			},
		},
		{
			"tags",
			"rnbq1knr/pPpp2pp/8/Pp6/7Q/8/8/R3K2R w KQ b6 0 1",
			chess.NoMove,
			[2]chess.Move{
				newMove(chess.H4, chess.H5, chess.WhiteQueen, chess.NoPiece, chess.NoPiece, chess.Quiet),
			},
			[]string{
				"b7c8q", "b7a8q", "b7c8n", "b7a8n", "e1g1",
				"e1c1", "a5b6", "h4d8", "h4h7", "h4h5",
				"e1f2", "e1d1", "e1f1", "a1b1", "a1c1",
				"a1d1", "a1a2", "a1a3", "a1a4", "h1f1",
				"h1g1", "h1h2", "h1h3", "h4f2", "h4h2",
				"h4g3", "h4h3", "h4a4", "h4b4", "h4c4",
				"h4d4", "h4e4", "h4f4", "h4g4", "h4g5",
				"a5a6", "h4f6", "h4h6", "h4e7", "e1d2",
				"e1e2", "b7c8b", "b7a8r", "b7a8b", "b7c8r",
			},
		},
		{
			"best move",
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			newMove(chess.D2, chess.D4, chess.WhitePawn, chess.NoPiece, chess.NoPiece, chess.Quiet),
			[2]chess.Move{},
			[]string{
				"d2d4", "g1f2", "c4c5", "g1h1", "f3d4",
				"b4c5", "f1f2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			checkData, _ := pos.InCheck()
			moves := pos.PseudoMoves(checkData)
			scoreMoves(pos, moves, tt.best, tt.killers)

			var sorted []chess.Move
			for i := 0; i < len(moves); i++ {
				nextOracle(moves, i)
				sorted = append(sorted, moves[i])
			}

			assert.Equal(t, tt.want, movesString(sorted))
		})
	}
}

func newMove(s1, s2 chess.Square, p1, p2, promo chess.Piece, tags chess.MoveTag) chess.Move {
	return chess.Move(s1) ^ chess.Move(s2)<<6 ^
		chess.Move(p1)<<12 ^ chess.Move(p2)<<16 ^
		chess.Move(promo)<<20 ^
		chess.Move(tags)
}
