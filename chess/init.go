package chess

// init is the master init function that calls sub init functions in the correct sequence.
func init() {
	initPieceTable()
	initPromoPieceTypeTable()

	initBBFiles()
	initBBRanks()
	initBBDiagonals()
	initBBAntiDiagonals()

	// Attacks
	initBBWhitePawnCaptures()
	initBBBlackPawnCaptures()
	initBBMagicRookMoves()
	initBBMagicBishopMoves()
}
