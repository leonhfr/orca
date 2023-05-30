package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlphaBeta(t *testing.T) {
	results := [6]struct {
		output searchResult
		moves  []string
	}{
		{
			output: searchResult{nodes: 3, score: 0},
			moves:  []string{"c6c7"},
		},
		{
			output: searchResult{nodes: 1, score: -mate},
			moves:  []string{},
		},
		{
			output: searchResult{nodes: 11, score: mate - 1},
			moves:  []string{"f1h1"},
		},
		{
			output: searchResult{nodes: 483, score: mate - 1},
			moves:  []string{"f6f2"},
		},
		{
			output: searchResult{nodes: 16029, score: mate - 3},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: searchResult{nodes: 308, score: 55},
			moves:  []string{"h8h7", "a1b2", "g7f8", "e7f8", "b3b2"},
		},
	}

	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{})
			si.killers.increaseDepth(int(tt.depth))

			res := results[i]
			pos := unsafeFEN(tt.fen)
			output, err := si.alphaBeta(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, res.output.nodes, output.nodes, fmt.Sprintf("want %d, got %d", res.output.nodes, output.nodes))
			assert.Equal(t, res.output.score, output.score, fmt.Sprintf("want %d, got %d", res.output.score, output.score))
			assert.Equal(t, res.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{})
			si.killers.increaseDepth(int(bb.depth))

			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.alphaBeta(context.Background(), pos, -mate, mate, bb.depth)
			}
		})
	}
}
