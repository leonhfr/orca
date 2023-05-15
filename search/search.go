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
func New(options ...func(*Engine)) *Engine {
	e := &Engine{}
	for _, o := range availableOptions {
		o.defaultFunc()(e)
	}
	for _, fn := range options {
		fn(e)
	}
	return e
}

// WithTableSize sets the size of the transposition table in MB.
func WithTableSize(size int) func(*Engine) {
	return func(e *Engine) {
		e.tableSize = size
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
	options := make([]uci.Option, len(availableOptions))
	for i, option := range availableOptions {
		options[i] = option.uci()
	}
	return options
}

// SetOption sets an option.
//
// Implements the uci.Engine interface.
func (e *Engine) SetOption(name, value string) error {
	for _, option := range availableOptions {
		if option.String() == name {
			fn, err := option.optionFunc(value)
			if err != nil {
				return err
			}
			fn(e)
			return nil
		}
	}

	return errOptionName
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
