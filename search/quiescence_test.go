package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type quiescenceSearchTestResult struct {
	score int32
	nodes uint32
}

func TestQuiescence(t *testing.T) {
	t.Parallel()
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
			result: quiescenceSearchTestResult{nodes: 4064, score: 4},
			moves:  []string{"b3b2", "a1b2", "c4c3", "d2c3"},
		},
		{
			name:   "horizon effect depth 5",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  5,
			result: quiescenceSearchTestResult{nodes: 17115, score: 3},
			moves:  []string{"c4c3", "d2c3", "b3b2", "b1b2", "b5b4"},
		},
		{
			name:   "horizon effect depth 6",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  6,
			result: quiescenceSearchTestResult{nodes: 16838, score: 0},
			moves:  []string{"g7a1", "b1a1", "h8g8", "d2d3", "c4d3", "g5g6", "f8a8", "a1b1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			si := newSearchInfo(newHashMapTable(), noPawnTable{})
			pos := unsafeFEN(tt.fen)
			score, err := si.alphaBeta(context.Background(), pos, -mate, mate, tt.depth, 0)
			pv := si.table.principalVariation(pos)

			assert.Equal(t, tt.result.nodes, si.nodes, "want %d, got %d", tt.result.nodes, si.nodes)
			assert.Equal(t, tt.result.score, score, "want %d, got %d", tt.result.score, score)
			assert.Equal(t, tt.moves, movesString(pv))
			assert.NoError(t, err)
		})
	}
}
