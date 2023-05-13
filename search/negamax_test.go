package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

var searchTestPositions = []struct {
	name  string
	fen   string
	depth int
}{
	{
		name:  "checkmate",
		fen:   "8/8/8/5K1k/8/8/8/7R b - - 0 1",
		depth: 1,
	},
	{
		name:  "mate in 1",
		fen:   "8/8/8/5K1k/8/8/8/5R2 w - - 0 1",
		depth: 2,
	},
	{
		name:  "mate in 1",
		fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
		depth: 2,
	},
	{
		name:  "mate in 2",
		fen:   "5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1",
		depth: 4,
	},
	{
		name:  "horizon effect",
		fen:   "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
		depth: 3,
	},
}

func TestNegamax(t *testing.T) {
	results := [5]struct {
		output Output
		moves  []string
	}{
		{
			output: Output{0, 1, -mate, 0, nil},
			moves:  []string{},
		},
		{
			output: Output{1, 39, mate - 1, 1, nil},
			moves:  []string{"f1h1"},
		},
		{
			output: Output{1, 1219, mate - 1, 1, nil},
			moves:  []string{"f6f2"},
		},
		{
			output: Output{3, 4103853, mate - 3, 2, nil},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: Output{3, 9561, 549, 0, nil},
			moves:  []string{"g7b2", "a1b2", "b3b2"},
		},
	}

	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			want := results[i]
			output, err := negamax(context.Background(), unsafeFEN(tt.fen), tt.depth)

			assert.Nil(t, err)
			assert.NotNil(t, output)
			fmt.Println(output)
			assert.Equal(t, want.output.Nodes, output.Nodes)
			assert.Equal(t, want.output.Score, output.Score)
			assert.Equal(t, want.moves, movesString(output.PV))
		})
	}
}

func BenchmarkNegamax(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = negamax(context.Background(), pos, bb.depth)
			}
		})
	}
}

func movesString(moves []chess.Move) []string {
	result := make([]string, len(moves))
	for i, move := range moves {
		result[i] = move.String()
	}
	return result
}
