package search

import (
	"context"
	"fmt"
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
			output: searchResult{nodes: 11, score: mate - 1},
			moves:  []string{"f1h1"},
		},
		{
			output: searchResult{nodes: 485, score: mate - 1},
			moves:  []string{"f6f2"},
		},
		{
			output: searchResult{nodes: 15616, score: mate - 3},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: searchResult{nodes: 444, score: 49},
			moves:  []string{"f8d8", "e7a3", "d5d4"},
		},
	}

	e := New()
	e.table = noTable{}
	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			res := results[i]
			pos := unsafeFEN(tt.fen)
			output, err := e.alphaBeta(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, res.output.nodes, output.nodes, fmt.Sprintf("want %d, got %d", res.output.nodes, output.nodes))
			assert.Equal(t, res.output.score, output.score, fmt.Sprintf("want %d, got %d", res.output.score, output.score))
			assert.Equal(t, res.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	e := New()
	e.table = noTable{}
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = e.alphaBeta(context.Background(), pos, -mate, mate, bb.depth)
			}
		})
	}
}
