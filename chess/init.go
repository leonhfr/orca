package chess

// init is the master init function that calls sub init functions in the correct sequence.
func init() {
	initPieceTable()
	initPromoPieceTypeTable()

	initBBFiles()
	initBBRanks()
	initBBDiagonals()
	initBBAntiDiagonals()
	initBBInBetweens()

	// Attacks
	initBBKingMoves()
	initBBKnightMoves()
	initBBMagicRookMoves()
	initBBMagicBishopMoves()

	// Castle checks
	initCastleChecks()
}
