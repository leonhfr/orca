package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerft(t *testing.T) {
	tests := []struct {
		fen   string
		nodes []int
	}{
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			[]int{
				20, 400, 8902, 197281, 4865609, 119060324,
				// 3195901860, 84998978956, 2439530234167, 69352859712417
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
				44, 1486, 62379, 2103487, 89941194,
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

	for _, tt := range tests {
		for depth := 0; depth < len(tt.nodes); depth++ {
			t.Run(fmt.Sprintf("%s depth %d", tt.fen, depth+1), func(t *testing.T) {
				pos := unsafeFEN(tt.fen)
				got := pos.Perft(depth + 1)
				assert.Equal(t, tt.nodes[depth], got.nodes)
			})
		}
	}
}
