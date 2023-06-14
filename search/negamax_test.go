package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

type searchTestResult struct {
	score int32
	nodes uint32
	moves []string
}

var searchTestPositions = []struct {
	name               string
	fen                string
	depth              uint8
	negamax            searchTestResult
	alphaBeta          searchTestResult
	principalVariation searchTestResult
	zeroWindow         searchTestResult
}{
	{
		name:  "draw stalemate in 1",
		fen:   "8/2b2k2/2K5/8/8/8/5n2/8 w - - 0 1",
		depth: 3,
		negamax: searchTestResult{
			score: 0,
			nodes: 718,
		},
		alphaBeta: searchTestResult{
			score: 0,
			nodes: 60,
			moves: []string{"c6c7"},
		},
		principalVariation: searchTestResult{
			score: 0,
			nodes: 60,
			moves: []string{"c6c7"},
		},
		zeroWindow: searchTestResult{
			score: mate - 1,
			nodes: 60,
		},
	},
	{
		name:  "checkmate",
		fen:   "8/8/8/5K1k/8/8/8/7R b - - 0 1",
		depth: 1,
		negamax: searchTestResult{
			score: -mate,
			nodes: 1,
		},
		alphaBeta: searchTestResult{
			score: -mate,
			nodes: 1,
			moves: []string{},
		},
		principalVariation: searchTestResult{
			score: -mate,
			nodes: 1,
			moves: []string{},
		},
		zeroWindow: searchTestResult{
			score: mate - 1,
			nodes: 1,
		},
	},
	{
		name:  "mate in 1",
		fen:   "8/8/8/5K1k/8/8/8/5R2 w - - 0 1",
		depth: 2,
		negamax: searchTestResult{
			score: mate - 1,
			nodes: 54,
		},
		alphaBeta: searchTestResult{
			score: mate - 1,
			nodes: 43,
			moves: []string{"f1h1"},
		},
		principalVariation: searchTestResult{
			score: mate - 1,
			nodes: 88,
			moves: []string{"f1h1"},
		},
		zeroWindow: searchTestResult{
			score: mate,
			nodes: 35,
		},
	},
	{
		name:  "mate in 1",
		fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
		depth: 2,
		negamax: searchTestResult{
			score: mate - 1,
			nodes: 1265,
		},
		alphaBeta: searchTestResult{
			score: mate - 1,
			nodes: 1187,
			moves: []string{"f6f2"},
		},
		principalVariation: searchTestResult{
			score: mate - 1,
			nodes: 247,
			moves: []string{"f6f2"},
		},
		zeroWindow: searchTestResult{
			score: mate,
			nodes: 55,
		},
	},
	{
		name:  "mate in 2",
		fen:   "5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1",
		depth: 4,
		negamax: searchTestResult{
			score: mate - 3,
			nodes: 4195950,
		},
		alphaBeta: searchTestResult{
			score: mate - 3,
			nodes: 28981,
			moves: []string{"c6g2", "e2g2", "c1e1"},
		},
		principalVariation: searchTestResult{
			score: mate - 3,
			nodes: 21558,
			moves: []string{"c6g2", "e2g2", "c1e1"},
		},
		zeroWindow: searchTestResult{
			score: mate,
			nodes: 2542,
		},
	},
	{
		name:  "horizon effect",
		fen:   "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
		depth: 3,
		negamax: searchTestResult{
			score: 575,
			nodes: 10065,
		},
		alphaBeta: searchTestResult{
			score: 49,
			nodes: 2259,
			moves: []string{"b3b2", "a1b2", "c4c3"},
		},
		principalVariation: searchTestResult{
			score: 49,
			nodes: 2918,
			moves: []string{"b3b2", "a1b2", "c4c3"},
		},
		zeroWindow: searchTestResult{
			score: mate - 1,
			nodes: 812,
		},
	},
}

// negamax performs a search using the Negamax algorithm.
//
// Negamax is a variant of minimax that relies on the
// zero-sum property of a two-player game.
func (si *searchInfo) negamax(ctx context.Context, pos *chess.Position, depth uint8) (int32, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		si.nodes++
	}

	if pos.HasInsufficientMaterial() {
		return draw, nil
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	checkData, inCheck := pos.InCheck()
	moves := pos.PseudoMoves(checkData)
	switch {
	case len(moves) == 0 && inCheck:
		return -mate, nil
	case len(moves) == 0:
		return draw, nil
	case depth == 0:
		return si.evaluate(pos), nil
	}

	var score int32 = -mate

	var validMoves int
	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		current, err := si.negamax(ctx, pos, depth-1)

		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if err != nil {
			return 0, err
		}

		current = -current
		if current > score {
			score = current
		}
	}

	if validMoves > 0 {
		score = incMateDistance(score)
	}

	return score, nil
}

func TestNegamax(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})
			score, err := si.negamax(context.Background(), unsafeFEN(tt.fen), tt.depth)

			want := tt.negamax
			assert.Nil(t, err)
			assert.NotNil(t, score)
			assert.Equal(t, want.nodes, si.nodes, fmt.Sprintf("nodes: want %d, got %d", want.nodes, si.nodes))
			assert.Equal(t, want.score, score, fmt.Sprintf("score: want %d, got %d", want.score, score))
		})
	}
}

func BenchmarkNegamax(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{}, noPawnTable{})
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.negamax(context.Background(), pos, bb.depth)
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
