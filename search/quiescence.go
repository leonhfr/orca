package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

func (si *searchInfo) quiesce(ctx context.Context, pos *chess.Position, alpha, beta int32) (searchResult, error) {
	select {
	case <-ctx.Done():
		return searchResult{}, context.Canceled
	default:
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	if standPat := evaluate(pos); standPat >= beta {
		return searchResult{
			nodes: 1,
			score: beta,
		}, nil
	} else if alpha < standPat {
		alpha = standPat
	}

	moves := pos.LoudMoves()
	scoreLoudMoves(moves)

	var nodes uint32
	for i := 0; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}

		current, err := si.quiesce(ctx, pos, -beta, -alpha)
		if err != nil {
			return searchResult{}, nil
		}

		current.score = -current.score
		nodes += current.nodes
		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if current.score >= beta {
			return searchResult{
				nodes: nodes,
				score: beta,
			}, nil
		} else if current.score > alpha {
			alpha = current.score
		}
	}

	return searchResult{
		nodes: nodes,
		score: alpha,
	}, nil
}
