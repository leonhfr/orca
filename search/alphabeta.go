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

	hash := pos.Hash()
	alphaOriginal := alpha
	entry, inCache := e.table.get(hash)
	if inCache && entry.depth >= depth {
		switch {
		case entry.nodeType == exact:
			return searchResult{
				nodes: 1,
				score: entry.score,
			}, nil
		case entry.nodeType == lowerBound && entry.score > alpha:
			alpha = entry.score
		case entry.nodeType == upperBound && entry.score < beta:
			beta = entry.score
		}

		if alpha >= beta {
			return searchResult{
				nodes: 1,
				score: entry.score,
			}, nil
		}
	}

	if pos.HasInsufficientMaterial() {
		return searchResult{
			score: draw,
			nodes: 1,
		}, nil
	}

	inCheck := pos.InCheck()
	if inCheck {
		depth++
	}

	if depth == 0 {
		return e.quiesce(ctx, pos, -beta, -alpha)
	}

	var validMoves int
	result := searchResult{
		score: -mate,
	}

	best := chess.NoMove
	if inCache && entry.best != chess.NoMove {
		best = entry.best
	}

	moves := pos.PseudoMoves()
	oracle(moves, best)

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

		pos.UnmakeMove(move, metadata, hash)

		if alpha >= beta {
			break
		}
	}

	switch {
	case validMoves == 0 && inCheck:
		mateResult := searchResult{
			nodes: 1,
			score: -mate,
		}
		e.storeResult(hash, depth, result, exact)
		return mateResult, nil
	case validMoves == 0:
		drawResult := searchResult{
			nodes: 1,
			score: draw,
		}
		e.storeResult(hash, depth, result, exact)
		return drawResult, nil
	}

	result.score = incMateDistance(result.score)
	nodeType := getNodeType(alphaOriginal, beta, result.score)
	e.storeResult(hash, depth, result, nodeType)

	return result, nil
}

// getNodeType returns the node type according to the alpha and beta bounds.
func getNodeType(alpha, beta, score int32) nodeType {
	switch {
	case score <= alpha:
		return upperBound
	case score >= beta:
		return lowerBound
	default:
		return exact
	}
}

// storeResult stores a search result in the transposition table.
func (e *Engine) storeResult(hash chess.Hash, depth uint8, r searchResult, n nodeType) {
	se := searchEntry{
		hash:     hash,
		score:    r.score,
		nodeType: n,
		depth:    depth,
	}
	if len(r.pv) > 0 {
		se.best = r.pv[len(r.pv)-1]
	}
	e.table.set(hash, se)
}
