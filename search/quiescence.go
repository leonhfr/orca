package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

func (si *searchInfo) quiesce(ctx context.Context, pos *chess.Position, alpha, beta int32) (int32, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		si.nodes++
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	if standPat := si.evaluate(pos); standPat >= beta {
		return beta, nil
	} else if alpha < standPat {
		alpha = standPat
	}

	moves := pos.LoudMoves()
	scoreLoudMoves(pos, moves)

	for i := 0; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}

		score, err := si.quiesce(ctx, pos, -beta, -alpha)
		if err != nil {
			return 0, nil
		}

		score = -score
		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if score >= beta {
			return beta, nil
		} else if score > alpha {
			alpha = score
		}
	}

	return alpha, nil
}
