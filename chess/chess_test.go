package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasInsufficientMaterial(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{
			name: "king versus king",
			fen:  "8/5k2/8/8/6K1/8/8/8 w - - 0 1",
			want: true,
		},
		{
			name: "king and bishop versus king",
			fen:  "8/2b2k2/8/8/6K1/8/8/8 w - - 0 1",
			want: true,
		},
		{
			name: "king and knight versus king",
			fen:  "8/2n2k2/8/8/6K1/8/8/8 w - - 0 1",
			want: true,
		},
		{
			name: "king and bishop versus king and bishop with the bishops on the same color",
			fen:  "8/2b2k2/8/8/3B2K1/8/8/8 w - - 0 1",
			want: true,
		},
		{
			name: "king and bishop versus king and bishop with the bishops on different colors",
			fen:  "8/2b2k2/8/8/4B1K1/8/8/8 b - - 0 1",
			want: false,
		},
		{
			name: "normal position",
			fen:  startFEN,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			got := pos.HasInsufficientMaterial()
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkHasInsufficientMaterial(b *testing.B) {
	fen := "8/2b2k2/8/8/3B2K1/8/8/8 w - - 0 1"
	pos := unsafeFEN(fen)
	for n := 0; n < b.N; n++ {
		pos.HasInsufficientMaterial()
	}
}

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
			"b7a8q", "b7b8q", "h3g2", "h3g4",
			// illegal moves
			"g1f1", "g1h1",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.fen, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var got []string
			checkData, _ := pos.InCheck()
			moves := pos.PseudoMoves(checkData)
			for _, move := range moves {
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
				checkData, _ := pos.InCheck()
				pos.PseudoMoves(checkData)
			}
		})
	}
}

func TestLoudMoves(t *testing.T) {
	tests := []struct {
		fen  string
		want []string
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", []string{}},
		{"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23", []string{
			"a4b3", "b7b4", "c4b3", "d7e5", "e7b4",
		}},
		{"r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9", []string{
			"c3d5", "c4d5", "d3h7",
		}},
		{"r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/R3K2R w KQkq - 3 10", []string{
			"c3d5", "c4d5", "d3h7",
		}},
		{"r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", []string{
			"e4d5",
		}},
		{"r1bqkbnr/ppp1p1pp/2n5/3pPp2/8/5N2/PPPP1PPP/RNBQKB1R w KQkq f6 0 4", []string{
			"e5f6",
		}},
		{"r1bqkbnr/ppp1p1pp/2n5/3pPp2/3N4/8/PPPP1PPP/RNBQKB1R b KQkq - 1 4", []string{
			"c6d4", "c6e5",
		}},
		{"r7/1Pp5/2P3p1/8/6pb/4p1kB/4P1p1/6K1 w - - 0 1", []string{
			"b7a8q", "h3g2", "h3g4",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.fen, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var got []string
			for _, move := range pos.LoudMoves() {
				got = append(got, move.String())
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func BenchmarkLoudMoves(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.preFEN, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				pos.LoudMoves()
			}
		})
	}
}

type perfTest struct {
	fen   string
	nodes []int
}

var perfResults = []perfTest{
	{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		[]int{
			20, 400, 8902, 197281, 4865609,
			// 119060324, 3195901860, 84998978956, 2439530234167, 69352859712417
		},
	},
	{
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		[]int{
			48, 2039, 97862, 4085603,
			//  193690690
		},
	},
	{
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		[]int{
			14, 191, 2812, 43238, 674624, 11030083,
			// 178633661
		},
	},
	{
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		[]int{
			6, 264, 9467, 422333, 15833292,
			// 706045033
		},
	},
	{
		"r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
		[]int{
			6, 264, 9467, 422333, 15833292,
			// 706045033
		},
	},
	{
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		[]int{
			44, 1486, 62379,
			// 2103487, 89941194,
		},
	},
	{
		"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		[]int{
			46, 2079, 89890, 3894594,
			//  164075551, 6923051137, 287188994746, 11923589843526, 490154852788714
		},
	},
}

func TestPerfResults(t *testing.T) {
	for _, tt := range perfResults {
		for depth := 0; depth < len(tt.nodes); depth++ {
			want := tt.nodes[depth]

			t.Run(fmt.Sprintf("%s depth %d", tt.fen, depth), func(t *testing.T) {
				got := perft(unsafeFEN(tt.fen), depth)
				assert.Equal(t, want, got)
			})
		}
	}
}

func perft(pos *Position, depth int) int {
	if depth == 0 {
		return len(legalMoves(pos))
	}

	var count int
	hash := pos.Hash()
	pseudoMoves := pos.pseudoMoves(bbFull, true, false, false)
	for _, m := range pseudoMoves {
		if meta, ok := pos.MakeMove(m); ok {
			count += perft(pos, depth-1)
			pos.UnmakeMove(m, meta, hash)
		}
	}
	return count
}

func legalMoves(pos *Position) []Move {
	hash := pos.Hash()
	var moves []Move
	pseudoMoves := pos.pseudoMoves(bbFull, true, false, false)
	for _, m := range pseudoMoves {
		if meta, ok := pos.MakeMove(m); ok {
			moves = append(moves, m)
			pos.UnmakeMove(m, meta, hash)
		}
	}
	return moves
}
