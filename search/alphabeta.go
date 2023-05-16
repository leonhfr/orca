package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

// alphaBeta performs a search using the Negamax algorithm
// and alpha-beta pruning.
func alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta, depth int) (uci.Output, error) {
	select {
	case <-ctx.Done():
		return uci.Output{}, context.Canceled
	default:
	}

	moves := pos.PseudoMoves()
	switch inCheck := pos.InCheck(); {
	case len(moves) == 0 && inCheck:
		return uci.Output{
			Nodes: 1,
			Score: -mate,
		}, nil
	case len(moves) == 0:
		return uci.Output{
			Nodes: 1,
			Score: draw,
		}, nil
	case depth == 0:
		return uci.Output{
			Nodes: 1,
			Score: evaluate(pos),
		}, nil
	}

	oracle(moves)

	result := uci.Output{
		Nodes: 1,
		Depth: depth,
		Score: -mate,
		PV:    make([]chess.Move, 0, depth),
	}

	var validMoves int
	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		current, err := alphaBeta(ctx, pos, -beta, -alpha, depth-1)
		if err != nil {
			return uci.Output{}, err
		}

		result.Nodes += current.Nodes
		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append(current.PV, move)
		}

		if current.Score > alpha {
			alpha = current.Score
		}

		pos.UnmakeMove(move, metadata)

		if alpha >= beta {
			break
		}
	}

	if validMoves > 0 {
		result.Nodes--
		result.Score = incMateDistance(result.Score)
	}
	return result, nil
}
