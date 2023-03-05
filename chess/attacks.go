package chess

var (
	// bbWhitePawnCaptures contains a lookup table of white pawn captures bitboard indexed by squares.
	bbWhitePawnCaptures = [64]bitboard{}
	// bbBlackPawnCaptures contains a lookup table of black pawn captures bitboard indexed by squares.
	bbBlackPawnCaptures = [64]bitboard{}
)

// initializes bbWhitePawnCaptures
func initBBWhitePawnCaptures() {
	for sq := A1; sq <= H8; sq++ {
		captureR := (sq.bitboard() & ^bbFileH & ^bbRank8) << 9
		captureL := (sq.bitboard() & ^bbFileA & ^bbRank8) << 7
		bbWhitePawnCaptures[sq] = captureR | captureL
	}
}

// initializes bbBlackPawnCaptures
func initBBBlackPawnCaptures() {
	for sq := A1; sq <= H8; sq++ {
		captureR := (sq.bitboard() & ^bbFileH & ^bbRank1) >> 7
		captureL := (sq.bitboard() & ^bbFileA & ^bbRank1) >> 9
		bbBlackPawnCaptures[sq] = captureR | captureL
	}
}
