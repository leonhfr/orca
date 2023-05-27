package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuiescence(t *testing.T) {
	tests := []struct {
		name   string
		fen    string
		depth  uint8
		result searchResult
		moves  []string
	}{
		{
			name:   "horizon effect depth 4",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  4,
			result: searchResult{nodes: 2087, score: 2},
			moves:  []string{"c4c3", "a1b2", "g7f8", "e7f8", "d5d4"},
		},
		{
			name:   "horizon effect depth 5",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  5,
			result: searchResult{nodes: 5267, score: 2},
			moves:  []string{"b1c1", "h8g7", "e3d4", "g7f8", "e7f8", "d5d4"},
		},
		{
			name:   "horizon effect depth 6",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  6,
			result: searchResult{nodes: 19904, score: 2},
			moves:  []string{"g7g6", "b1c1", "h8g7", "e3d4", "g7f8", "e7f8", "d5d4"},
		},
	}

	e := NewEngine()
	e.table = noTable{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			output, err := e.alphaBeta(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, tt.result.nodes, output.nodes, fmt.Sprintf("want %d, got %d", tt.result.nodes, output.nodes))
			assert.Equal(t, tt.result.score, output.score, fmt.Sprintf("want %d, got %d", tt.result.score, output.score))
			assert.Equal(t, tt.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}
