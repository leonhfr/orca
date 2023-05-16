package uci

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/leonhfr/orca/chess"
)

// Engine is the interface that should be implemented by the search Engine.
type Engine interface {
	// Init initializes the search engine.
	Init() error
	// Close shuts down the resources used by the search engine.
	Close()
	// Options lists the available options.
	Options() []Option
	// SeOption sets an option.
	SetOption(name, value string) error
	// Search runs a search on the given position until the given depth.
	// Cancelling the context stops the search.
	Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan Output
}

// Run runs the program in UCI mode.
//
// Run parses command from the reader, executes them with the provided
// search engine and writes the responses on the writer.
func Run(ctx context.Context, e Engine, r io.Reader, s *State) {
	// graceful shutdown when context canceled
	// sending EOF to the UCI scanner by closing the pipe
	pipeR, pipeW := io.Pipe()
	go func() { _, _ = io.Copy(pipeW, r) }()
	go func() { <-ctx.Done(); pipeW.Close() }()

	for scanner := bufio.NewScanner(pipeR); scanner.Scan(); {
		c := parse(strings.Fields(scanner.Text()))
		if c == nil {
			continue
		}
		c.run(ctx, e, s)
		if _, ok := c.(commandQuit); ok {
			break
		}
	}
}

// Output holds a search output.
type Output struct {
	Depth int          // Search depth in plies.
	Nodes int          // Number of nodes searched.
	Score int          // Score from the engine's point of view in centipawns.
	Mate  int          // Number of moves before mate. Positive for the current player to mate, negative for the current player to be mated.
	PV    []chess.Move // Principal variation, best line found.
}

// OptionType represent an option type.
type OptionType uint8

const OptionInteger OptionType = iota // OptionInteger represents an integer option.

// Option represents an available option.
type Option struct {
	Type    OptionType
	Name    string
	Default string
	Min     string
	Max     string
}
