package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

var searchTestPositions = [6]struct {
	name  string
	fen   string
	depth uint8
}{
	{
		name:  "draw stalemate in 1",
		fen:   "8/2b2k2/2K5/8/8/8/5n2/8 w - - 0 1",
		depth: 3,
	},
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
func negamax(ctx context.Context, pos *chess.Position, depth uint8) (searchResult, error) {
	select {
	case <-ctx.Done():
		return searchResult{}, context.Canceled
	default:
	}

	if pos.HasInsufficientMaterial() {
		return searchResult{
			score: draw,
			nodes: 1,
		}, nil
	}

	hash := pos.Hash()
	inCheck := pos.InCheck()
	moves := pos.PseudoMoves()
	switch {
	case len(moves) == 0 && inCheck:
		return searchResult{
			nodes: 1,
			score: -mate,
		}, nil
	case len(moves) == 0:
		return searchResult{
			nodes: 1,
			score: draw,
		}, nil
	case depth == 0:
		return searchResult{
			nodes: 1,
			score: evaluate(pos),
		}, nil
	}

	result := searchResult{
		nodes: 1,
		score: -mate,
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
			return searchResult{}, err
		}

		result.nodes += current.nodes
		current.score = -current.score
		if current.score > result.score {
			result.score = current.score
			result.pv = append(current.pv, move)
		}

		pos.UnmakeMove(move, metadata, hash)
	}

	if validMoves > 0 {
		result.nodes--
		result.score = incMateDistance(result.score)
	}
	return result, nil
}

func TestNegamax(t *testing.T) {
	results := [6]struct {
		output searchResult
		moves  []string
	}{
		{
			output: searchResult{nodes: 601, score: 0},
			moves:  []string{"c6c7"},
		},
		{
			output: searchResult{nodes: 1, score: -mate},
			moves:  []string{},
		},
		{
			output: searchResult{nodes: 39, score: mate - 1},
			moves:  []string{"f1h1"},
		},
		{
			output: searchResult{nodes: 1219, score: mate - 1},
			moves:  []string{"f6f2"},
		},
		{
			output: searchResult{nodes: 4103853, score: mate - 3},
			moves:  []string{"c1e1", "e2g2", "c6g2"},
		},
		{
			output: searchResult{nodes: 9561, score: 549},
			moves:  []string{"g7b2", "a1b2", "b3b2"},
		},
	}

	for i, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			want := results[i]
			output, err := negamax(context.Background(), unsafeFEN(tt.fen), tt.depth)

			assert.Nil(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, want.output.nodes, output.nodes, fmt.Sprintf("want %d, got %d", want.output.nodes, output.nodes))
			assert.Equal(t, want.output.score, output.score, fmt.Sprintf("want %d, got %d", want.output.score, output.score))
			assert.Equal(t, want.moves, movesString(output.pv))
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
