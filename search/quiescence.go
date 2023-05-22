package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

func (e *Engine) quiesce(ctx context.Context, pos *chess.Position, alpha, beta int32) (searchResult, error) {
	select {
	case <-ctx.Done():
		return searchResult{}, context.Canceled
	default:
	}

	hash := pos.Hash()
	if standPat := evaluate(pos); standPat >= beta {
		return searchResult{
			nodes: 1,
			score: beta,
		}, nil
	} else if alpha < standPat {
		alpha = standPat
	}

	moves := pos.LoudMoves()
	loudOracle(moves)

	var nodes uint32
	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}

		current, err := e.quiesce(ctx, pos, -beta, -alpha)
		if err != nil {
			return searchResult{}, nil
		}

		current.score = -current.score
		nodes += current.nodes
		pos.UnmakeMove(move, metadata, hash)

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
