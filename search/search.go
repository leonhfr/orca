package search

import "math"

const (
	// maxPkgDepth is the maximum depth at which the package will search.
	maxPkgDepth = 64
	// mate is the score of a checkmate.
	mate = math.MaxInt
	// draw is the score of a draw.
	draw = 0
)
