package uci

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/leonhfr/orca/chess"
)

// Output holds a search output.
type Output struct {
	Depth int          // Search depth in plies.
	Nodes int          // Number of nodes searched.
	Score int          // Score from the engine's point of view in centipawns.
	Mate  int          // Number of moves before mate. Positive for the current player to mate, negative for the current player to be mated.
	PV    []chess.Move // Principal variation, best line found.
}

// engine is the interface implemented by the search engine.
type engine interface {
	// Search runs a search on the given position until the given depth.
	//
	// Cancelling the context stops the search.
	Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan *Output
}

// Run runs the program in UCI mode.
//
// Run parses command from the reader, executes them with the provided
// search engine and writes the responses on the writer.
func Run(ctx context.Context, e engine, r io.Reader, s *State) {
	for scanner := bufio.NewScanner(r); scanner.Scan(); {
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
