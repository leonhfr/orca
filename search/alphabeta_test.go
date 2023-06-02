package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlphaBeta(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			res := tt.alphaBeta
			pos := unsafeFEN(tt.fen)
			output, err := si.alphaBeta(context.Background(), pos, -mate, mate, tt.depth, 0)

			assert.Equal(t, res.nodes, output.nodes, fmt.Sprintf("want %d, got %d", res.nodes, output.nodes))
			assert.Equal(t, res.score, output.score, fmt.Sprintf("want %d, got %d", res.score, output.score))
			assert.Equal(t, res.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.alphaBeta(context.Background(), pos, -mate, mate, bb.depth, 0)
			}
		})
	}
}
