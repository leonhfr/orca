package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroWindow(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			res := tt.zeroWindow
			pos := unsafeFEN(tt.fen)
			score, err := si.zeroWindow(context.Background(), pos, mate, tt.depth)

			assert.Equal(t, res.nodes, si.nodes, fmt.Sprintf("want %d, got %d", res.nodes, si.nodes))
			assert.Equal(t, res.score, score, fmt.Sprintf("want %d, got %d", res.score, score))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkZeroWindow(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.zeroWindow(context.Background(), pos, mate, bb.depth)
			}
		})
	}
}
