package search

import (
	"context"
	"math"
	"sync"

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
type Engine struct {
	once      sync.Once
	tableSize int
	table     *transpositionTable
}

// New creates a new search engine.
func New() *Engine {
	return &Engine{
		tableSize: 64,
	}
}

// Init initializes the search engine.
//
// Implements the uci.Engine interface.
func (e *Engine) Init() error {
	var err error
	e.once.Do(func() {
		e.table, err = newTable(e.tableSize)
	})
	return err
}

// Close shuts down the resources used by the search engine.
//
// Implements the uci.Engine interface.
func (e *Engine) Close() {
	_ = e.Init()
	e.table.close()
}

// Options lists the available options.
//
// Implements the uci.Engine interface.
func (e *Engine) Options() []uci.Option {
	return nil
}

// SetOption sets an option.
//
// Implements the uci.Engine interface.
func (e *Engine) SetOption(_, _ string) error {
	return nil
}

// Search runs a search on the given position until the given depth.
// Cancelling the context stops the search.
//
// Implements the uci.Engine interface.
func (e *Engine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan *uci.Output {
	_ = e.Init()
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
