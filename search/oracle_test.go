package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderMoves(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			"promotions",
			"7k/P7/8/8/8/8/8/K7 w - - 0 1",
			[]string{
				"a7a8q", "a7a8n", "a1b1", "a1a2", "a1b2",
				"a7a8r", "a7a8b",
			},
		},
		{
			"tags",
			"rnbq1knr/pPpp2pp/8/Pp6/7Q/8/8/R3K2R w KQ b6 0 1",
			[]string{
				"b7a8q", "b7c8q", "b7a8n", "b7c8n", "e1g1",
				"e1c1", "h4d8", "h4h7", "h4f2", "h4b4",
				"e1f2", "e1e2", "a1b1", "a1c1", "a1d1",
				"a1a2", "a1a3", "a1a4", "h1f1", "h1g1",
				"h1h2", "h1h3", "e1d2", "h4h2", "h4g3",
				"h4h3", "h4a4", "e1f1", "h4c4", "h4d4",
				"h4e4", "h4f4", "h4g4", "h4g5", "h4h5",
				"h4f6", "h4h6", "h4e7", "e1d1", "a5a6",
				"a5b6", "b7a8r", "b7c8r", "b7c8b", "b7a8b",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := unsafeFEN(tt.args).PseudoMoves()
			oracle(moves)

			assert.Equal(t, tt.want, movesString(moves))
		})
	}
}
