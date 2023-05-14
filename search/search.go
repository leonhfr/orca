package search

import (
	"context"
	"math"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

const (
	// maxPkgDepth is the maximum depth at which the package will search.
	maxPkgDepth = 64
	// mate is the score of a checkmate.
	mate = math.MaxInt
	// draw is the score of a draw.
	draw = 0
)

// Engine represents the search engine.
type Engine struct{}

// Search implements the uci.Engine interface.
func (Engine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan *uci.Output {
	output := make(chan *uci.Output)

	if maxDepth == 0 || maxDepth > maxPkgDepth {
		maxDepth = maxPkgDepth
	}

	go func() {
		defer close(output)
		iterativeSearch(ctx, pos, maxDepth, output)
	}()

	return output
}

// iterativeSearch performs an iterative search.
func iterativeSearch(ctx context.Context, pos *chess.Position, maxDepth int, output chan<- *uci.Output) {
	for depth := 1; depth <= maxDepth; depth++ {
		o, err := alphaBeta(ctx, pos, -mate, mate, depth)
		if err != nil {
			return
		}
		o.Mate = mateIn(o.Score)
		output <- o
	}
}

// mateIn returns the number of moves before mate.
func mateIn(score int) int {
	sign := sign(score)
	delta := mate - sign*score
	if delta > maxPkgDepth {
		return 0
	}
	return sign * (delta/2 + delta%2)
}

// sign returns the sign +/-1 of the passed integer.
func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}
