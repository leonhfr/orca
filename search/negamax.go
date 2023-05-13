package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

// negamax performs a search using the Negamax algorithm.
//
// Negamax is a variant of minimax that relies on the
// zero-sum property of a two-player game.
func negamax(ctx context.Context, pos *chess.Position, depth int) (*uci.Output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	moves := pos.PseudoMoves()
	switch inCheck := pos.InCheck(); {
	case len(moves) == 0 && inCheck:
		return &uci.Output{
			Nodes: 1,
			Score: -mate,
		}, nil
	case len(moves) == 0:
		return &uci.Output{
			Nodes: 1,
			Score: draw,
		}, nil
	case depth == 0:
		return &uci.Output{
			Nodes: 1,
			Score: evaluate(pos),
		}, nil
	}

	result := &uci.Output{
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

		current, err := negamax(ctx, pos, depth-1)
		if err != nil {
			return nil, err
		}

		result.Nodes += current.Nodes
		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append(current.PV, move)
		}

		pos.UnmakeMove(move, metadata)
	}

	if validMoves > 0 {
		result.Nodes--
		result.Score = incMateDistance(result.Score)
	}
	return result, nil
}
