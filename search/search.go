package search

import (
	"math"

	"github.com/leonhfr/orca/chess"
)

const (
	// maxPkgDepth is the maximum depth at which the package will search.
	maxPkgDepth = 64
	// mate is the score of a checkmate.
	mate = math.MaxInt
	// draw is the score of a draw.
	draw = 0
)

// Output holds a search output.
type Output struct {
	Depth int          // Search depth in plies.
	Nodes int          // Number of nodes searched.
	Score int          // Score from the engine's point of view in centipawns.
	Mate  int          // Number of moves before mate. Positive for the current player to mate, negative for the current player to be mated.
	PV    []chess.Move // Principal variation, best line found.
}
