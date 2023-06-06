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
	name      string
	fen       string
	depth     uint8
	negamax   searchTestResult
	alphaBeta searchTestResult
}{
	{
		name:  "draw stalemate in 1",
		fen:   "8/2b2k2/2K5/8/8/8/5n2/8 w - - 0 1",
		depth: 3,
		negamax: searchTestResult{
			score: 0,
			nodes: 601,
			moves: []string{"c6c7"},
		},
		alphaBeta: searchTestResult{
			score: 0,
			nodes: 3,
			moves: []string{"c6c7"},
		},
	},
	{
		name:  "checkmate",
		fen:   "8/8/8/5K1k/8/8/8/7R b - - 0 1",
		depth: 1,
		negamax: searchTestResult{
			score: -mate,
			nodes: 1,
			moves: []string{},
		},
		alphaBeta: searchTestResult{
			score: -mate,
			nodes: 1,
			moves: []string{},
		},
	},
	{
		name:  "mate in 1",
		fen:   "8/8/8/5K1k/8/8/8/5R2 w - - 0 1",
		depth: 2,
		negamax: searchTestResult{
			score: mate - 1,
			nodes: 39,
			moves: []string{"f1h1"},
		},
		alphaBeta: searchTestResult{
			score: mate - 1,
			nodes: 11,
			moves: []string{"f1h1"},
		},
	},
	{
		name:  "mate in 1",
		fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
		depth: 2,
		negamax: searchTestResult{
			score: mate - 1,
			nodes: 1219,
			moves: []string{"f6f2"},
		},
		alphaBeta: searchTestResult{
			score: mate - 1,
			nodes: 482,
			moves: []string{"f6f2"},
		},
	},
	{
		name:  "mate in 2",
		fen:   "5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1",
		depth: 4,
		negamax: searchTestResult{
			score: mate - 3,
			nodes: 4103853,
			moves: []string{"c1e1", "e2g2", "c6g2"},
		},
		alphaBeta: searchTestResult{
			score: mate - 3,
			nodes: 16112,
			moves: []string{"c1e1", "e2g2", "c6g2"},
		},
	},
	{
		name:  "horizon effect",
		fen:   "5r1k/4Qpq1/4p3/1p1p2P1/2p2P2/1p2P3/3P4/BK6 b - - 0 1",
		depth: 3,
		negamax: searchTestResult{
			score: 555,
			nodes: 9561,
			moves: []string{"g7b2", "a1b2", "b3b2"},
		},
		alphaBeta: searchTestResult{
			score: 34,
			nodes: 307,
			moves: []string{"h8h7", "a1b2", "g7f8", "e7f8", "b3b2"},
		},
	},
}

// negamax performs a search using the Negamax algorithm.
//
// Negamax is a variant of minimax that relies on the
// zero-sum property of a two-player game.
func (si *searchInfo) negamax(ctx context.Context, pos *chess.Position, depth uint8) (searchResult, error) {
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
	pawnHash := pos.PawnHash()
	checkData, inCheck := pos.InCheck()
	moves := pos.PseudoMoves(checkData)
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
			score: si.evaluate(pos),
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

		current, err := si.negamax(ctx, pos, depth-1)
		if err != nil {
			return searchResult{}, err
		}

		result.nodes += current.nodes
		current.score = -current.score
		if current.score > result.score {
			result.score = current.score
			result.pv = append(current.pv, move)
		}

		pos.UnmakeMove(move, metadata, hash, pawnHash)
	}

	if validMoves > 0 {
		result.nodes--
		result.score = incMateDistance(result.score)
	}
	return result, nil
}

func TestNegamax(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})
			output, err := si.negamax(context.Background(), unsafeFEN(tt.fen), tt.depth)

			want := tt.negamax
			assert.Nil(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, want.nodes, output.nodes, fmt.Sprintf("nodes: want %d, got %d", want.nodes, output.nodes))
			assert.Equal(t, want.score, output.score, fmt.Sprintf("score: want %d, got %d", want.score, output.score))
			assert.Equal(t, want.moves, movesString(output.pv))
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
