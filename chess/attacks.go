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

// slowMoves computes the move bitboard for each piece type.
//
// This function is intended to be used during initialization of move tables.
func slowMoves(pt PieceType, sq Square, blockers bitboard) bitboard {
	switch pt {
	case Rook:
		return linearBitboard(sq, blockers, bbFiles[sq]) | linearBitboard(sq, blockers, bbRanks[sq])
	case Bishop:
		return linearBitboard(sq, blockers, bbDiagonals[sq]) | linearBitboard(sq, blockers, bbAntiDiagonals[sq])
	default:
		return emptyBitboard
	}
}

// linearBitboard computes a slider attack bitboard.
func linearBitboard(sq Square, occupied, mask bitboard) bitboard {
	inMask := occupied & mask
	return ((inMask - 2*sq.bitboard()) ^ (inMask.reverse() - 2*sq.bitboard().reverse()).reverse()) & mask
}

// slowMasks computes the mask bitboard for each piece type.
//
// This function is intended to be used during initialization of move tables.
func slowMasks(pt PieceType, sq Square) bitboard {
	switch pt {
	case Rook:
		file := bbFiles[sq] & ^(bbRank1 | bbRank8)
		rank := bbRanks[sq] & ^(bbFileA | bbFileH)
		return (file | rank) & ^sq.bitboard()
	case Bishop:
		mask := bbRank1 | bbRank8 | bbFileA | bbFileH
		diagonal := bbDiagonals[sq] & ^mask
		antiDiagonal := bbAntiDiagonals[sq] & ^mask
		return (diagonal | antiDiagonal) & ^sq.bitboard()
	default:
		return emptyBitboard
	}
}
