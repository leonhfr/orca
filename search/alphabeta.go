package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

// alphaBeta performs a search using the Negamax algorithm
// and alpha-beta pruning.
func (si *searchInfo) alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta int32, depth, index uint8) (int32, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		si.nodes++
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()
	originalAlpha := alpha

	entry, inCache := si.table.get(hash)
	if inCache && entry.depth >= depth {
		switch {
		case entry.nodeType == exact:
			return entry.score, nil
		case entry.nodeType == lowerBound && entry.score > alpha:
			alpha = entry.score
		case entry.nodeType == upperBound && entry.score < beta:
			beta = entry.score
		}

		if alpha >= beta {
			return entry.score, nil
		}
	}

	if pos.HasInsufficientMaterial() {
		return draw, nil
	}

	checkData, inCheck := pos.InCheck()
	if inCheck {
		depth++
	}

	if depth == 0 {
		return si.quiesce(ctx, pos, -beta, -alpha)
	}

	var validMoves int
	var score int32 = -mate

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
			return 0, err
		}

		current = -current
		if current > score {
			score = current
			best = move
		}

		if current > alpha {
			alpha = current
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
		si.table.set(hash, best, -mate, exact, depth)
		return -mate, nil
	case validMoves == 0:
		si.table.set(hash, best, draw, exact, depth)
		return draw, nil
	}

	score = incMateDistance(score)
	nodeType := getNodeType(originalAlpha, beta, score)
	si.table.set(hash, best, score, nodeType, depth)

	return score, nil
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
