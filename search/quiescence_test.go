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
			result: quiescenceSearchTestResult{nodes: 8054, score: 49},
			moves:  []string{"b3b2", "a1b2", "c4c3", "d2c3"},
		},
		{
			name:   "horizon effect depth 5",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  5,
			result: quiescenceSearchTestResult{nodes: 44073, score: 0},
			moves:  []string{"c4c3", "a1b2", "c3b2", "e7d6", "g7g5"},
		},
		{
			name:   "horizon effect depth 6",
			fen:    "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
			depth:  6,
			result: quiescenceSearchTestResult{nodes: 173014, score: 0},
			moves:  []string{"c4c3", "a1b2", "c3b2", "e7d6", "g7g5", "f4g5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(newHashMapTable(), noPawnTable{})
			pos := unsafeFEN(tt.fen)
			score, err := si.alphaBeta(context.Background(), pos, -mate, mate, tt.depth, 0)
			pv := si.table.principalVariation(pos)

			assert.Equal(t, tt.result.nodes, si.nodes, fmt.Sprintf("want %d, got %d", tt.result.nodes, si.nodes))
			assert.Equal(t, tt.result.score, score, fmt.Sprintf("want %d, got %d", tt.result.score, score))
			assert.Equal(t, tt.moves, movesString(pv))
			assert.Nil(t, err)
		})
	}
}
