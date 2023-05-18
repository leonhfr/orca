package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

// alphaBeta performs a search using the Negamax algorithm
// and alpha-beta pruning.
func (e *Engine) alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta, depth int) (uci.Output, error) {
	select {
	case <-ctx.Done():
		return uci.Output{}, context.Canceled
	default:
	}

	alphaOriginal := alpha
	cached, ok := e.table.get(pos.Hash())
	if ok && cached.depth >= depth {
		switch {
		case cached.nodeType == exact:
			return uci.Output{
				Nodes: 1,
				Score: cached.score,
			}, nil
		case cached.nodeType == lowerBound && cached.score > alpha:
			alpha = cached.score
		case cached.nodeType == upperBound && cached.score < beta:
			beta = cached.score
		}

		if alpha >= beta {
			return uci.Output{
				Nodes: 1,
				Score: cached.score,
			}, nil
		}
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
	}

	var validMoves int
	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		current, err := e.alphaBeta(ctx, pos, -beta, -alpha, depth-1)
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

	nodeType := exact
	switch {
	case result.Score <= alphaOriginal:
		nodeType = upperBound
	case result.Score >= beta:
		nodeType = lowerBound
	}
	e.table.set(pos.Hash(), tableEntry{
		score:    result.Score,
		depth:    result.Depth,
		nodeType: nodeType,
	})

	return result, nil
}
