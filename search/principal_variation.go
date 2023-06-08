package search

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

// principalVariation performs a search using the principal variation algorithm.
func (si *searchInfo) principalVariation(ctx context.Context, pos *chess.Position, alpha, beta int32, depth, index uint8) (int32, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		si.nodes++
	}

	hash := pos.Hash()
	pawnHash := pos.PawnHash()

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
		si.table.set(hash, chess.NoMove, draw, exact, depth)
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
	var best chess.Move
	nt := upperBound

	moves := pos.PseudoMoves(checkData)
	scoreMoves(pos, moves, entry.best, si.killers.get(index))

	for i, searchPv := 0, true; i < len(moves); i++ {
		nextOracle(moves, i)
		move := moves[i]

		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}
		validMoves++

		var score int32
		var err error

		if searchPv {
			score, err = si.principalVariation(ctx, pos, -beta, -alpha, depth-1, index+1)
			if err != nil {
				return 0, err
			}

			score = -score
		} else {
			score, err = si.zeroWindow(ctx, pos, -alpha, depth-1)
			if err != nil {
				return 0, err
			}

			score = -score

			if score > alpha {
				score, err = si.principalVariation(ctx, pos, -beta, -alpha, depth-1, index+1)
				if err != nil {
					return 0, err
				}

				score = -score
			}
		}

		pos.UnmakeMove(move, metadata, hash, pawnHash)

		if score >= beta {
			if move.HasTag(chess.Quiet) {
				si.killers.set(move, index)
			}

			beta = incMateDistance(beta)
			si.table.set(hash, move, beta, lowerBound, depth)
			return beta, nil
		}

		if score > alpha {
			alpha = score
			best = move
			nt = exact
			searchPv = false
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

	alpha = incMateDistance(alpha)
	si.table.set(hash, best, alpha, nt, depth)

	return alpha, nil
}
