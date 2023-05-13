package uci

import (
	"context"

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
	Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan Output
}
