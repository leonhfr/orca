package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

// alphaBeta performs a search using the Negamax algorithm
// and alpha-beta pruning.
func (si *searchInfo) alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta int32, depth, index uint8) (int32, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		si.nodes++
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	originalAlpha := alpha

	entry, inCache := si.table.get(hash)
	if inCache && entry.depth >= depth {
		switch {
		case entry.nodeType == exact:
			return entry.score, nil
		case entry.nodeType == lowerBound && entry.score > alpha:
			alpha = entry.score
		case entry.nodeType == upperBound && entry.score < beta:
			beta = entry.score
		}

		if alpha >= beta {
			return entry.score, nil
		}
	}

	if pos.HasInsufficientMaterial() {
		return draw, nil
	}

	checkData, inCheck := pos.InCheck()
	if inCheck {
		depth++
	}

	if depth == 0 {
		return si.quiesce(ctx, pos, -beta, -alpha)
	}

	var validMoves int
	var score int32 = -mate

	best := chess.NoMove
	if inCache && entry.best != chess.NoMove {
		best = entry.best
	}

	moves := pos.PseudoMoves(checkData)
	scoreMoves(pos, moves, best, si.killers.get(index))

	for i := 0; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		lmr := lateMoveReduction(validMoves, inCheck, depth, move)
		current, err := si.alphaBeta(ctx, pos, -beta, -alpha, depth-lmr-1, index+1)

		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if err != nil {
			return 0, err
		}

		current = -current
		if current > score {
			score = current
			best = move
		}

		if current > alpha {
			alpha = current
		}

		if alpha >= beta {
			if move.HasTag(chess.Quiet) {
				si.killers.set(move, index)
			}
			break
		}
	}

	switch {
	case validMoves == 0 && inCheck:
		si.table.set(hash, best, -mate, exact, depth)
		return -mate, nil
	case validMoves == 0:
		si.table.set(hash, best, draw, exact, depth)
		return draw, nil
	default:
		score = incMateDistance(score)
		nodeType := getNodeType(originalAlpha, beta, score)
		si.table.set(hash, best, score, nodeType, depth)
		return score, nil
	}
}

// getNodeType returns the node type according to the alpha and beta bounds.
func getNodeType(alpha, beta, score int32) nodeType {
	switch {
	case score <= alpha:
		return upperBound
	case score >= beta:
		return lowerBound
	default:
		return exact
	}
}

func TestAlphaBeta(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(newHashMapTable(), noPawnTable{})

			res := tt.alphaBeta
			pos := unsafeFEN(tt.fen)
			score, err := si.alphaBeta(context.Background(), pos, -mate, mate, tt.depth, 0)
			pv := si.table.principalVariation(pos)

			assert.Equal(t, res.nodes, si.nodes, fmt.Sprintf("want %d, got %d", res.nodes, si.nodes))
			assert.Equal(t, res.score, score, fmt.Sprintf("want %d, got %d", res.score, score))
			assert.Equal(t, res.moves, movesString(pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.alphaBeta(context.Background(), pos, -mate, mate, bb.depth, 0)
			}
		})
	}
}
