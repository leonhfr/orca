package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		depth int
		oo    []uci.Output
	}{
		{
			name:  "mate in 2",
			fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
			depth: 2,
			oo: []uci.Output{
				{Depth: 1, Nodes: 46, Score: 357, Mate: 0, PV: []chess.Move{0x2c322dc}},
				{Depth: 2, Nodes: 58, Score: 9223372036854775806, Mate: 1, PV: []chess.Move{0x2c1836d}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Engine{}.Search(context.Background(), unsafeFEN(tt.fen), tt.depth)
			var outputs []uci.Output
			for o := range output {
				outputs = append(outputs, *o)
			}

			assert.Equal(t, tt.oo, outputs)
		})
	}
}
