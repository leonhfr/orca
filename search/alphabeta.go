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
func (si *searchInfo) alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta int32, depth, index uint8) (searchResult, error) {
	select {
	case <-ctx.Done():
		return searchResult{}, context.Canceled
	default:
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	originalAlpha := alpha

	entry, inCache := si.table.get(hash)
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

	checkData, inCheck := pos.InCheck()
	if inCheck {
		depth++
	}

	if depth == 0 {
		return si.quiesce(ctx, pos, -beta, -alpha)
	}

	var validMoves int
	result := searchResult{
		score: -mate,
	}

	best := chess.NoMove
	if inCache && entry.best != chess.NoMove {
		best = entry.best
	}

	moves := pos.PseudoMoves(checkData)
	scoreMoves(moves, best, si.killers.get(index))

	for i := 0; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		current, err := si.alphaBeta(ctx, pos, -beta, -alpha, depth-1, index+1)
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

		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if alpha >= beta {
			if move.HasTag(chess.Quiet) {
				si.killers.set(move, index)
			}
			break
		}
	}

	switch {
	case validMoves == 0 && inCheck:
		mateResult := searchResult{
			nodes: 1,
			score: -mate,
		}
		si.storeResult(hash, depth, result, exact)
		return mateResult, nil
	case validMoves == 0:
		drawResult := searchResult{
			nodes: 1,
			score: draw,
		}
		si.storeResult(hash, depth, result, exact)
		return drawResult, nil
	}

	result.score = incMateDistance(result.score)
	nodeType := getNodeType(originalAlpha, beta, result.score)
	si.storeResult(hash, depth, result, nodeType)

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
func (si *searchInfo) storeResult(hash chess.Hash, depth uint8, r searchResult, n nodeType) {
	se := searchEntry{
		hash:     hash,
		score:    r.score,
		nodeType: n,
		depth:    depth,
	}
	if len(r.pv) > 0 {
		se.best = r.pv[len(r.pv)-1]
	}
	si.table.set(hash, se)
}
