// Package search provides a search engine that implements the uci.Engine interface.
//
// Handles alpha-beta minimax search, quiescence search, position evaluation, move ordering,
// and caching of search results using transposition tables.
package search

import (
	"bytes"
	"context"
	"math"
	"math/rand"
	"sync"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/data/books"
	"github.com/leonhfr/orca/uci"
)

const (
	// maxSearchDepth is the maximum depth at which the package will search.
	maxSearchDepth = 64
	// mate is the score of a checkmate.
	mate = math.MaxInt32
	// draw is the score of a draw.
	draw = 0
)

// Engine represents the search engine.
//
// Implements the uci.Engine interface.
//
//nolint:govet
type Engine struct {
	book      *chess.Book
	killers   *killerList
	once      sync.Once
	ownBook   bool
	tableSize int
	table     transpositionTable
	pawnTable transpositionPawnTable
}

// NewEngine creates a new search engine.
func NewEngine(options ...func(*Engine)) *Engine {
	e := &Engine{
		book:      chess.NewBook(),
		killers:   newKillerList(),
		table:     noTable{},
		pawnTable: noPawnTable{},
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
		err = e.book.Init(performance)
		e.killers = newKillerList()
		e.table = newArrayTable(e.tableSize)
		e.pawnTable = newArrayPawnTable(8)
	})
	return err
}

// Close shuts down the resources used by the search engine.
//
// Implements the uci.Engine interface.
func (e *Engine) Close() {
	_ = e.Init()
	e.table.close()
	e.pawnTable.close()
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
func (e *Engine) Search(ctx context.Context, pos *chess.Position, maxDepth, maxNodes int) <-chan uci.Output {
	_ = e.Init()
	output := make(chan uci.Output)

	go func() {
		defer close(output)

		if e.ownBook {
			moves := e.book.Lookup(pos)
			if move := weightedRandomMove(moves); move != chess.NoMove {
				output <- uci.Output{
					PV:    []chess.Move{move},
					Depth: 1,
					Nodes: 1,
				}

				return
			}
		}

		e.table.inc()
		e.iterativeSearch(ctx, pos, maxDepth, maxNodes, output)
	}()

	return output
}

// searchInfo contains info on the running search.
type searchInfo struct {
	killers   *killerList
	table     transpositionTable
	pawnTable transpositionPawnTable
	nodes     uint32
}

// newSearchInfo returns a new searchInfo.
func newSearchInfo(table transpositionTable, pawnTable transpositionPawnTable) *searchInfo {
	return &searchInfo{
		killers:   newKillerList(),
		table:     table,
		pawnTable: pawnTable,
	}
}

// iterativeSearch performs an iterative search.
func (e *Engine) iterativeSearch(ctx context.Context, pos *chess.Position, maxDepth, maxNodes int, output chan<- uci.Output) {
	si := newSearchInfo(e.table, e.pawnTable)

	if maxDepth <= 0 || maxDepth > maxSearchDepth {
		maxDepth = maxSearchDepth
	}

	if maxNodes <= 0 {
		maxNodes = math.MaxInt
	}

	for depth := 1; depth <= maxDepth; depth++ {
		o, err := si.alphaBeta(ctx, pos, -mate, mate, uint8(depth), 0)
		if err != nil {
			return
		}

		pv := e.table.principalVariation(pos)

		maxDepth := depth
		if len(pv) > maxDepth {
			maxDepth = len(pv)
		}

		nodes := int(si.nodes)

		output <- uci.Output{
			Depth: maxDepth,
			Score: int(o.score),
			Nodes: nodes,
			Mate:  int(mateIn(o.score)),
			PV:    pv,
		}

		if nodes >= maxNodes {
			break
		}
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
func mateIn(score int32) int32 {
	sign := sign(score)
	delta := mate - sign*score
	if delta > maxSearchDepth {
		return 0
	}
	return sign * (delta/2 + delta%2)
}

// sign returns the sign +/-1 of the passed integer.
func sign(n int32) int32 {
	if n < 0 {
		return -1
	}
	return 1
}
