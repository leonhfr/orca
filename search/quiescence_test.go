package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type quiescenceSearchTestResult struct {
	score int32
	nodes uint32
}

func TestQuiescence(t *testing.T) {
	tests := []struct {
		name   string
		fen    string
		depth  uint8
		result quiescenceSearchTestResult
		moves  []string
	}{
		{
			name:   "horizon effect depth 4",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  4,
			result: quiescenceSearchTestResult{nodes: 950, score: 42},
			moves:  []string{"b1a1", "h8h7", "a1b2", "g7f8", "e7f8", "b3b2"},
		},
		{
			name:   "horizon effect depth 5",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  5,
			result: quiescenceSearchTestResult{nodes: 6967, score: 1},
			moves:  []string{"b1c1", "h8g7", "e3d4", "g7f8", "e7f8", "d5d4"},
		},
		{
			name:   "horizon effect depth 6",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  6,
			result: quiescenceSearchTestResult{nodes: 10262, score: 0},
			moves:  []string{"a1b1", "f8a8", "f4f5", "c4d3", "d2d3", "h8g8", "b1a1", "g7a1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})
			pos := unsafeFEN(tt.fen)
			output, err := si.alphaBeta(context.Background(), pos, -mate, mate, tt.depth, 0)

			assert.Equal(t, tt.result.nodes, si.nodes, fmt.Sprintf("want %d, got %d", tt.result.nodes, si.nodes))
			assert.Equal(t, tt.result.score, output.score, fmt.Sprintf("want %d, got %d", tt.result.score, output.score))
			assert.Equal(t, tt.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}
