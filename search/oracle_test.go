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
		done chess.Move
		want []string
	}{
		{
			"promotions",
			"7k/P7/8/8/8/8/8/K7 w - - 0 1",
			chess.NoMove,
			[]string{"a7a8q", "a1b1", "a1a2", "a1b2"},
		},
		{
			"tags",
			"rnbq1knr/pPpp2pp/8/Pp6/7Q/8/8/R3K2R w KQ b6 0 1",
			chess.NoMove,
			[]string{
				"b7a8q", "b7c8q", "e1g1", "e1c1", "h4d8",
				"h4h7", "h1g1", "h4g3", "e1f1", "e1f2",
				"e1e2", "a1b1", "a1c1", "a1d1", "a1a2",
				"a1a3", "a1a4", "h1f1", "e1d1", "h1h2",
				"h1h3", "h4f2", "h4h2", "a5a6", "h4h3",
				"h4a4", "h4b4", "h4c4", "h4d4", "h4e4",
				"h4f4", "h4g4", "h4g5", "h4h5", "h4f6",
				"h4h6", "h4e7", "e1d2", "a5b6",
			},
		},
		{
			"done move",
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			chess.Move(chess.D2) ^
				chess.Move(chess.D4)<<6 ^
				chess.Move(chess.WhitePawn)<<12 ^
				chess.Move(chess.NoPiece)<<16 ^
				chess.Move(chess.NoPiece)<<20 ^
				chess.Move(chess.Quiet),
			[]string{
				"g1h1", "g1f2", "c4c5", "f3d4", "b4c5",
				"f1f2", "d2d4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves, _ := unsafeFEN(tt.fen).PseudoMoves()
			oracle(moves, tt.done)

			assert.Equal(t, tt.want, movesString(moves))
		})
	}
}
