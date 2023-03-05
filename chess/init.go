package chess

// init is the master init function that calls sub init functions in the correct sequence.
func init() {
	initPieceTable()

	// Attacks
	initBBWhitePawnCaptures()
	initBBBlackPawnCaptures()
}
