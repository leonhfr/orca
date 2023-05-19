package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestOrderMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		best chess.Move
		want []string
	}{
		{
			"promotions",
			"7k/P7/8/8/8/8/8/K7 w - - 0 1",
			chess.NoMove,
			[]string{
				"a7a8q", "a7a8n", "a1b1", "a1a2", "a1b2",
				"a7a8r", "a7a8b",
			},
		},
		{
			"tags",
			"rnbq1knr/pPpp2pp/8/Pp6/7Q/8/8/R3K2R w KQ b6 0 1",
			chess.NoMove,
			[]string{
				"b7a8q", "b7c8q", "b7a8n", "b7c8n", "e1g1",
				"e1c1", "h4d8", "h4h7", "h4f2", "h4b4",
				"e1f2", "e1e2", "a1b1", "a1c1", "a1d1",
				"a1a2", "a1a3", "a1a4", "h1f1", "h1g1",
				"h1h2", "h1h3", "e1d2", "h4h2", "h4g3",
				"h4h3", "h4a4", "e1f1", "h4c4", "h4d4",
				"h4e4", "h4f4", "h4g4", "h4g5", "h4h5",
				"h4f6", "h4h6", "h4e7", "e1d1", "a5a6",
				"a5b6", "b7a8r", "b7c8r", "b7c8b", "b7a8b",
			},
		},
		{
			"best move",
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			chess.Move(chess.D2) ^
				chess.Move(chess.D4)<<6 ^
				chess.Move(chess.WhitePawn)<<12 ^
				chess.Move(chess.NoPiece)<<16 ^
				chess.Move(chess.NoPiece)<<20 ^
				chess.Move(chess.Quiet),
			[]string{
				"d2d4", "b4a3", "b4a5", "h6f7", "a4b3",
				"g2g3", "g2g4", "h2h4", "f3e1", "f3d4",
				"f3h4", "f3e5", "f3g5", "h6g4", "h6f5",
				"c4c5", "h6g8", "a4c2", "d2d3", "h2h3",
				"b4c3", "e4e5", "b4c5", "b4d6", "b4e7",
				"b4f8", "a1b1", "a1c1", "f1e1", "f1f2",
				"d1b1", "d1c1", "d1e1", "d1c2", "d1e2",
				"d1b3", "g1h1", "g1f2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := unsafeFEN(tt.fen).PseudoMoves()
			oracle(moves, tt.best)

			assert.Equal(t, tt.want, movesString(moves))
		})
	}
}
