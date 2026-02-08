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

	meta := pos.Metadata()
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

		if ok := pos.MakeMove(move); !ok {
			continue
		}

		score, err := si.quiesce(ctx, pos, -beta, -alpha)

		score = -score
		pos.UnmakeMove(move, meta, hash, pawnHash)

		if err != nil {
			return 0, err
		}

		if score >= beta {
			return beta, nil
		} else if score > alpha {
			alpha = score
		}
	}

	return alpha, nil
}
