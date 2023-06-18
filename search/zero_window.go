package search

import (
	"context"

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

	meta := pos.Metadata()
	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	moves := pos.PseudoMoves(checkData)
	quickScoreMoves(pos, moves)

	for i := 0; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		if ok := pos.MakeMove(move); !ok {
			continue
		}

		score, err := si.zeroWindow(ctx, pos, 1-beta, depth-1)

		pos.UnmakeMove(move, meta, hash, pawnHash)

		if err != nil {
			return 0, err
		}
		score = -score

		if score >= beta {
			return beta, nil
		}
	}

	return beta - 1, nil
}
