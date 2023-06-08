package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func (si *searchInfo) zeroWindow(ctx context.Context, pos *chess.Position, beta int32, depth uint8) (int32, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		si.nodes++
	}

	if pos.HasInsufficientMaterial() {
		return draw, nil
	}

	checkData, inCheck := pos.InCheck()
	if inCheck {
		depth++
	}

	if depth == 0 {
		return si.quiesce(ctx, pos, beta-1, beta)
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	moves := pos.PseudoMoves(checkData)
	quickScoreMoves(moves)

	for i := 0; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}

		score, err := si.zeroWindow(ctx, pos, 1-beta, depth-1)
		if err != nil {
			return 0, err
		}
		score = -score

		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if score >= beta {
			return beta, nil
		}
	}

	return beta - 1, nil
}

func TestZeroWindow(t *testing.T) {
	for _, tt := range searchTestPositions {
		t.Run(tt.name, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			res := tt.zeroWindow
			pos := unsafeFEN(tt.fen)
			score, err := si.zeroWindow(context.Background(), pos, mate, tt.depth)

			assert.Equal(t, res.nodes, si.nodes, fmt.Sprintf("want %d, got %d", res.nodes, si.nodes))
			assert.Equal(t, res.score, score, fmt.Sprintf("want %d, got %d", res.score, score))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkZeroWindow(b *testing.B) {
	for _, bb := range searchTestPositions {
		b.Run(bb.name, func(b *testing.B) {
			si := newSearchInfo(noTable{}, noPawnTable{})

			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = si.zeroWindow(context.Background(), pos, mate, bb.depth)
			}
		})
	}
}
