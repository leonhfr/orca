package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

// principalVariation performs a search using the principal variation algorithm.
func (si *searchInfo) principalVariation(ctx context.Context, pos *chess.Position, alpha, beta int32, depth, index uint8) (int32, error) {
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

	best := chess.NoMove
	if inCache && entry.best != chess.NoMove {
		best = entry.best
	}

	moves := pos.PseudoMoves(checkData)
	scoreMoves(moves, best, si.killers.get(index))

	for i, searchPv := 0, true; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		var score int32
		var err error

		if searchPv {
			score, err = si.principalVariation(ctx, pos, -beta, -alpha, depth-1, index+1)
			if err != nil {
				return 0, err
			}

			score = -score
		} else {
			score, err = si.zeroWindow(ctx, pos, -alpha, depth-1)
			if err != nil {
				return 0, err
			}

			score = -score

			if score > alpha {
				score, err = si.principalVariation(ctx, pos, -beta, -alpha, depth-1, index+1)
				if err != nil {
					return 0, err
				}

				score = -score
			}
		}

		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if score >= beta {
			if move.HasTag(chess.Quiet) {
				si.killers.set(move, index)
			}

			beta = incMateDistance(beta)
			si.table.set(hash, move, beta, lowerBound, depth)
			return beta, nil
		}

		if score > alpha {
			alpha = score
			best = move
			searchPv = false
		}
	}

	switch {
	case validMoves == 0 && inCheck:
		si.table.set(hash, best, -mate, exact, depth)
		return -mate, nil
	case validMoves == 0:
		si.table.set(hash, best, draw, exact, depth)
		return draw, nil
	}

	alpha = incMateDistance(alpha)
	nodeType := getNodeType(originalAlpha, beta, alpha)
	si.table.set(hash, best, alpha, nodeType, depth)

	return alpha, nil
}

func TestPrincipalVariation(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(newHashMapTable(), noPawnTable{})

			res := tt.principalVariation
			pos := unsafeFEN(tt.fen)
			score, err := si.principalVariation(context.Background(), pos, -mate, mate, tt.depth, 0)
			pv := si.table.principalVariation(pos)

			assert.Equal(t, res.nodes, si.nodes, fmt.Sprintf("want %d, got %d", res.nodes, si.nodes))
			assert.Equal(t, res.score, score, fmt.Sprintf("want %d, got %d", res.score, score))
			assert.Equal(t, res.moves, movesString(pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkPrincipalVariation(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.principalVariation(context.Background(), pos, -mate, mate, bb.depth, 0)
			}
		})
	}
}
