package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPseudoMoves(t *testing.T) {
	tests := []struct {
		fen  string
		want []string
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", []string{
			"a2a3", "a2a4", "b1a3", "b1c3", "b2b3", "b2b4", "c2c3", "c2c4", "d2d3", "d2d4",
			"e2e3", "e2e4", "f2f3", "f2f4", "g1f3", "g1h3", "g2g3", "g2g4", "h2h3", "h2h4",
		}},
		{"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23", []string{
			"a4b3", "a6a5", "a6a7", "a6a8", "a6b6", "a6c6", "a6d6", "b7a7", "b7a8", "b7b4",
			"b7b5", "b7b6", "b7b8", "b7c6", "b7c7", "c4b3", "c8a8", "c8b8", "c8c5", "c8c6",
			"c8c7", "c8d8", "c8e8", "c8f8", "d7b6", "d7b8", "d7c5", "d7e5", "d7f6", "d7f8",
			"e7b4", "e7c5", "e7d6", "e7d8", "e7f6", "e7f8", "e7g5", "e7h4", "f7f5", "f7f6",
			"g7g5", "g7g6", "g8f8", "g8h8", "h7h5", "h7h6",
		}},
		{"r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9", []string{
			"a1b1", "a2a3", "a2a4", "b3b4", "c1a3", "c1b2", "c1d2", "c3a4", "c3b1", "c3b5",
			"c3d1", "c3d5", "c3e2", "c3e4", "c4c5", "c4d5", "d3b1", "d3c2", "d3d1", "d3d2",
			"d3e2", "d3e4", "d3f1", "d3f5", "d3g6", "d3h7", "e1d1", "e1d2", "e1e2", "e1f1",
			"e1g1", "e3e4", "f3d2", "f3e5", "f3g1", "f3g5", "f3h4", "g2g3", "g2g4", "h1f1",
			"h1g1", "h2h3", "h2h4",
		}},
		{"r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/R3K2R w KQkq - 3 10", []string{
			"a1b1", "a1c1", "a1d1", "a2a3", "a2a4", "b3b4", "c3a4", "c3b1", "c3b5", "c3d1",
			"c3d5", "c3e2", "c3e4", "c4c5", "c4d5", "d2c1", "d3b1", "d3c2", "d3e2", "d3e4",
			"d3f1", "d3f5", "d3g6", "d3h7", "e1c1", "e1d1", "e1e2", "e1f1", "e1g1", "e3e4",
			"f3e5", "f3g1", "f3g5", "f3h4", "g2g3", "g2g4", "h1f1", "h1g1", "h2h3", "h2h4",
		}},
		{"r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", []string{
			"a2a3", "a2a4", "b1a3", "b1c3", "b2b3", "b2b4", "c2c3", "c2c4", "d1e2", "d2d3",
			"d2d4", "e1e2", "e4d5", "e4e5", "f1a6", "f1b5", "f1c4", "f1d3", "f1e2", "f3d4",
			"f3e5", "f3g1", "f3g5", "f3h4", "g2g3", "g2g4", "h1g1", "h2h3", "h2h4",
		}},
		{"r1bqkbnr/ppp1p1pp/2n5/3pPp2/8/5N2/PPPP1PPP/RNBQKB1R w KQkq f6 0 4", []string{
			"a2a3", "a2a4", "b1a3", "b1c3", "b2b3", "b2b4", "c2c3", "c2c4", "d1e2", "d2d3",
			"d2d4", "e1e2", "e5e6", "e5f6", "f1a6", "f1b5", "f1c4", "f1d3", "f1e2", "f3d4",
			"f3g1", "f3g5", "f3h4", "g2g3", "g2g4", "h1g1", "h2h3", "h2h4",
		}},
		{"r1bqkbnr/ppp1p1pp/2n5/3pPp2/3N4/8/PPPP1PPP/RNBQKB1R b KQkq - 1 4", []string{
			"a7a5", "a7a6", "a8b8", "b7b5", "b7b6", "c6a5", "c6b4", "c6b8", "c6d4", "c6e5",
			"c8d7", "c8e6", "d8d6", "d8d7", "e7e6", "e8d7", "e8f7", "f5f4", "g7g5", "g7g6",
			"g8f6", "g8h6", "h7h5", "h7h6",
		}},
		{"r7/1Pp5/2P3p1/8/6pb/4p1kB/4P1p1/6K1 w - - 0 1", []string{
			"b7a8b", "b7a8n", "b7a8q", "b7a8r", "b7b8b", "b7b8n", "b7b8q", "b7b8r", "h3g2", "h3g4",
			// illegal moves
			"g1f1", "g1f2", "g1g2", "g1h1", "g1h2",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.fen, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var got []string
			for _, move := range pos.PseudoMoves() {
				got = append(got, move.String())
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func BenchmarkPseudoMoves(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.preFEN, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				pos.PseudoMoves()
			}
		})
	}
}
