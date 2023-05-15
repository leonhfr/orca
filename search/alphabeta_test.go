package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/uci"
)

func TestAlphaBeta(t *testing.T) {
	results := [5]struct {
		output uci.Output
		moves  []string
	}{
		{
			output: uci.Output{Depth: 0, Nodes: 1, Score: -mate, Mate: 0, PV: nil},
			moves:  []string{},
		},
		{
			output: uci.Output{Depth: 1, Nodes: 10, Score: mate - 1, Mate: 0, PV: nil},
			moves:  []string{"f1h1"},
		},
		{
			output: uci.Output{Depth: 1, Nodes: 58, Score: mate - 1, Mate: 0, PV: nil},
			moves:  []string{"f6f2"},
		},
		{
			output: uci.Output{Depth: 3, Nodes: 1779, Score: mate - 3, Mate: 0, PV: nil},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: uci.Output{Depth: 3, Nodes: 306, Score: 549, Mate: 0, PV: nil},
			moves:  []string{"g7b2", "a1b2", "b3b2"},
		},
	}

	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			res := results[i]
			pos := unsafeFEN(tt.fen)
			output, err := alphaBeta(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, res.output.Nodes, output.Nodes)
			assert.Equal(t, res.output.Score, output.Score)
			assert.Equal(t, res.moves, movesString(output.PV))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = alphaBeta(context.Background(), pos, -mate, mate, bb.depth)
			}
		})
	}
}