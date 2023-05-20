package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlphaBeta(t *testing.T) {
	results := [5]struct {
		output searchResult
		moves  []string
	}{
		{
			output: searchResult{nodes: 1, score: -mate},
			moves:  []string{},
		},
		{
			output: searchResult{nodes: 23, score: mate - 1},
			moves:  []string{"f1h1"},
		},
		{
			output: searchResult{nodes: 58, score: mate - 1},
			moves:  []string{"f6f2"},
		},
		{
			output: searchResult{nodes: 1914, score: mate - 3},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: searchResult{nodes: 303, score: 549},
			moves:  []string{"g7b2", "a1b2", "b3b2"},
		},
	}

	e := New()
	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			res := results[i]
			pos := unsafeFEN(tt.fen)
			output, err := e.alphaBeta(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, res.output.nodes, output.nodes)
			assert.Equal(t, res.output.score, output.score)
			assert.Equal(t, res.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	e := New()
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = e.alphaBeta(context.Background(), pos, -mate, mate, bb.depth)
			}
		})
	}
}
