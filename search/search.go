package search

import (
	"bytes"
	"context"
	"math"
	"math/rand"
	"sync"

	"github.com/leonhfr/orca/chess"
	books "github.com/leonhfr/orca/data/books"
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
//
//nolint:govet
type Engine struct {
	book      *chess.Book
	once      sync.Once
	ownBook   bool
	tableSize int
	table     transpositionTable
}

// New creates a new search engine.
func New(options ...func(*Engine)) *Engine {
	e := &Engine{
		book:  chess.NewBook(),
		table: noTable{},
	}
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

// WithOwnBook determines the use of the search engine's own opening book.
func WithOwnBook(on bool) func(*Engine) {
	return func(e *Engine) {
		e.ownBook = on
	}
}

// Init initializes the search engine.
//
// Implements the uci.Engine interface.
func (e *Engine) Init() error {
	var err error
	e.once.Do(func() {
		performance := bytes.NewReader(books.Performance)
		if err = e.book.Init(performance); err != nil {
			return
		}
		e.table, err = newRistrettoTable(e.tableSize)
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
func (e *Engine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan uci.Output {
	_ = e.Init()
	output := make(chan uci.Output)

	go func() {
		defer close(output)

		if e.ownBook {
			if moves := e.book.Lookup(pos); len(moves) > 0 {
				if move := weightedRandomMove(moves); move != chess.NoMove {
					output <- uci.Output{
						PV:    []chess.Move{move},
						Depth: 1,
						Nodes: 1,
						Score: 1,
					}

					return
				}
			}
		}

		e.iterativeSearch(ctx, pos, maxDepth, output)
	}()

	return output
}

// iterativeSearch performs an iterative search.
func (e *Engine) iterativeSearch(ctx context.Context, pos *chess.Position, maxDepth int, output chan<- uci.Output) {
	if maxDepth <= 0 || maxDepth > maxPkgDepth {
		maxDepth = maxPkgDepth
	}

	for depth := 1; depth <= maxDepth; depth++ {
		o, err := e.alphaBeta(ctx, pos, -mate, mate, depth)
		if err != nil {
			return
		}
		pv := make([]chess.Move, len(o.PV))
		for i, m := range o.PV {
			pv[len(o.PV)-i-1] = m
		}
		o.Mate = mateIn(o.Score)
		o.PV = pv
		output <- o
	}
}

// weightedRandomMove randomly selects a move with weighted probabilities.
func weightedRandomMove(moves []chess.WeightedMove) chess.Move {
	var sum int
	for _, move := range moves {
		sum += move.Weight
	}
	if sum <= 0 {
		return chess.NoMove
	}
	index := rand.Intn(sum) //nolint:gosec
	for _, move := range moves {
		if index < move.Weight {
			return move.Move
		}
		index -= move.Weight
	}
	return chess.NoMove
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
