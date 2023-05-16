package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

var searchTestPositions = []struct {
	name  string
	fen   string
	depth int
}{
	{
		name:  "checkmate",
		fen:   "8/8/8/5K1k/8/8/8/7R b - - 0 1",
		depth: 1,
	},
	{
		name:  "mate in 1",
		fen:   "8/8/8/5K1k/8/8/8/5R2 w - - 0 1",
		depth: 2,
	},
	{
		name:  "mate in 1",
		fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
		depth: 2,
	},
	{
		name:  "mate in 2",
		fen:   "5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1",
		depth: 4,
	},
	{
		name:  "horizon effect",
		fen:   "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
		depth: 3,
	},
}

// negamax performs a search using the Negamax algorithm.
//
// Negamax is a variant of minimax that relies on the
// zero-sum property of a two-player game.
func negamax(ctx context.Context, pos *chess.Position, depth int) (uci.Output, error) {
	select {
	case <-ctx.Done():
		return uci.Output{}, context.Canceled
	default:
	}

	moves := pos.PseudoMoves()
	switch inCheck := pos.InCheck(); {
	case len(moves) == 0 && inCheck:
		return uci.Output{
			Nodes: 1,
			Score: -mate,
		}, nil
	case len(moves) == 0:
		return uci.Output{
			Nodes: 1,
			Score: draw,
		}, nil
	case depth == 0:
		return uci.Output{
			Nodes: 1,
			Score: evaluate(pos),
		}, nil
	}

	result := uci.Output{
		Nodes: 1,
		Depth: depth,
		Score: -mate,
		PV:    make([]chess.Move, 0, depth),
	}

	var validMoves int
	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		current, err := negamax(ctx, pos, depth-1)
		if err != nil {
			return uci.Output{}, err
		}

		result.Nodes += current.Nodes
		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append(current.PV, move)
		}

		pos.UnmakeMove(move, metadata)
	}

	if validMoves > 0 {
		result.Nodes--
		result.Score = incMateDistance(result.Score)
	}
	return result, nil
}

func TestNegamax(t *testing.T) {
	results := [5]struct {
		output uci.Output
		moves  []string
	}{
		{
			output: uci.Output{Depth: 0, Nodes: 1, Score: -mate, Mate: 0, PV: nil},
			moves:  []string{},
		},
		{
			output: uci.Output{Depth: 1, Nodes: 39, Score: mate - 1, Mate: 1, PV: nil},
			moves:  []string{"f1h1"},
		},
		{
			output: uci.Output{Depth: 1, Nodes: 1219, Score: mate - 1, Mate: 1, PV: nil},
			moves:  []string{"f6f2"},
		},
		{
			output: uci.Output{Depth: 3, Nodes: 4103853, Score: mate - 3, Mate: 2, PV: nil},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: uci.Output{Depth: 3, Nodes: 9561, Score: 549, Mate: 0, PV: nil},
			moves:  []string{"g7b2", "a1b2", "b3b2"},
		},
	}

	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			want := results[i]
			output, err := negamax(context.Background(), unsafeFEN(tt.fen), tt.depth)

			assert.Nil(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, want.output.Nodes, output.Nodes)
			assert.Equal(t, want.output.Score, output.Score)
			assert.Equal(t, want.moves, movesString(output.PV))
		})
	}
}

func BenchmarkNegamax(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = negamax(context.Background(), pos, bb.depth)
			}
		})
	}
}

func movesString(moves []chess.Move) []string {
	result := make([]string, len(moves))
	for i, move := range moves {
		result[i] = move.String()
	}
	return result
}
