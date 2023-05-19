package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

// searchResult contains a search result.
type searchResult struct {
	pv    []chess.Move
	score int32
	nodes uint32
}

// alphaBeta performs a search using the Negamax algorithm
// and alpha-beta pruning.
func (e *Engine) alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta int32, depth uint8) (searchResult, error) {
	select {
	case <-ctx.Done():
		return searchResult{}, context.Canceled
	default:
	}

	alphaOriginal := alpha
	cached, ok := e.table.get(pos.Hash())
	if ok && cached.depth >= depth {
		switch {
		case cached.nodeType == exact:
			return searchResult{
				nodes: 1,
				score: cached.score,
			}, nil
		case cached.nodeType == lowerBound && cached.score > alpha:
			alpha = cached.score
		case cached.nodeType == upperBound && cached.score < beta:
			beta = cached.score
		}

		if alpha >= beta {
			return searchResult{
				nodes: 1,
				score: cached.score,
			}, nil
		}
	}

	moves, inCheck := pos.PseudoMoves()
	switch {
	case len(moves) == 0 && inCheck:
		return searchResult{
			nodes: 1,
			score: -mate,
		}, nil
	case len(moves) == 0:
		return searchResult{
			nodes: 1,
			score: draw,
		}, nil
	case depth == 0:
		return searchResult{
			nodes: 1,
			score: evaluate(pos),
		}, nil
	}

	oracle(moves, cached.best)

	result := searchResult{
		nodes: 1,
		score: -mate,
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
			return searchResult{}, err
		}

		result.nodes += current.nodes
		current.score = -current.score
		if current.score > result.score {
			result.score = current.score
			result.pv = append(current.pv, move)
		}

		if current.score > alpha {
			alpha = current.score
		}

		pos.UnmakeMove(move, metadata)

		if alpha >= beta {
			break
		}
	}

	if validMoves > 0 {
		result.nodes--
		result.score = incMateDistance(result.score)
	}

	nodeType := exact
	switch {
	case result.score <= alphaOriginal:
		nodeType = upperBound
	case result.score >= beta:
		nodeType = lowerBound
	}

	se := searchEntry{
		score:    result.score,
		nodeType: nodeType,
	}
	if len(result.pv) > 0 {
		se.best = result.pv[len(result.pv)-1]
	}
	e.table.set(pos.Hash(), se)

	return result, nil
}
