package chess

var (
	bbWhitePawnCaptures = [64]bitboard{} // bbWhitePawnCaptures contains a lookup table of white pawn captures bitboard indexed by squares.
	bbBlackPawnCaptures = [64]bitboard{} // bbBlackPawnCaptures contains a lookup table of black pawn captures bitboard indexed by squares.
	bbKingMoves         = [64]bitboard{} // bbMagicRookMoves contains a lookup table of king moves indexed by squares.
	bbKnightMoves       = [64]bitboard{} // bbMagicRookMoves contains a lookup table of king moves indexed by squares.
	bbMagicRookMoves    []bitboard       // bbMagicRookMoves contains a lookup table of rook moves indexed by magics.
	bbMagicBishopMoves  []bitboard       // bbMagicBishopMoves contains a lookup table of bishop moves indexed by magics.
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

// initializes bbKingMoves
func initBBKingMoves() {
	for sq := A1; sq <= H8; sq++ {
		var bb bitboard
		for dest, ok := range map[Square]bool{
			sq + 8 - 1: sq.Rank() <= Rank7 && sq.File() >= FileB,
			sq + 8:     sq.Rank() <= Rank7,
			sq + 8 + 1: sq.Rank() <= Rank7 && sq.File() <= FileG,
			sq + 1:     sq.File() <= FileG,
			sq - 8 + 1: sq.Rank() >= Rank2 && sq.File() <= FileG,
			sq - 8:     sq.Rank() >= Rank2,
			sq - 8 - 1: sq.Rank() >= Rank2 && sq.File() >= FileB,
			sq - 1:     sq.File() >= FileB,
		} {
			if ok {
				bb ^= dest.bitboard()
			}
		}
		bbKingMoves[sq] = bb
	}
}

// initializes bbKnightMoves
func initBBKnightMoves() {
	for sq := A1; sq <= H8; sq++ {
		var bb bitboard
		for dest, ok := range map[Square]bool{
			sq + 8 - 2:  sq.Rank() <= Rank7 && sq.File() >= FileC,
			sq + 16 - 1: sq.Rank() <= Rank6 && sq.File() >= FileB,
			sq + 16 + 1: sq.Rank() <= Rank6 && sq.File() <= FileG,
			sq + 8 + 2:  sq.Rank() <= Rank7 && sq.File() <= FileF,
			sq - 8 + 2:  sq.Rank() >= Rank2 && sq.File() <= FileF,
			sq - 16 + 1: sq.Rank() >= Rank3 && sq.File() <= FileG,
			sq - 16 - 1: sq.Rank() >= Rank3 && sq.File() >= FileB,
			sq - 8 - 2:  sq.Rank() >= Rank2 && sq.File() >= FileC,
		} {
			if ok {
				bb ^= dest.bitboard()
			}
		}
		bbKnightMoves[sq] = bb
	}
}

// initializes rookMoves
//
// requires bbFiles, bbRanks, bbDiagonals, bbAntiDiagonals
func initBBMagicRookMoves() {
	for sq := A1; sq <= H8; sq++ {
		moves, _ := slowMoveTable(Rook, sq, rookMagics[sq])
		bbMagicRookMoves = append(bbMagicRookMoves, moves...)
	}
}

// initializes bishopMoves
//
// requires bbFiles, bbRanks, bbDiagonals, bbAntiDiagonals
func initBBMagicBishopMoves() {
	for sq := A1; sq <= H8; sq++ {
		moves, _ := slowMoveTable(Bishop, sq, bishopMagics[sq])
		bbMagicBishopMoves = append(bbMagicBishopMoves, moves...)
	}
}

// moveBitboard returns the move bitboard.
//
// The returned bitboard has to be NOT AND with the bitboard of the color whose turn it is.
func moveBitboard(sq Square, pt PieceType, occupancy bitboard) bitboard {
	switch pt {
	case King:
		return bbKingMoves[sq]
	case Queen:
		rIndex := rookMagics[sq].index(occupancy)
		bIndex := bishopMagics[sq].index(occupancy)
		return bbMagicRookMoves[rIndex] | bbMagicBishopMoves[bIndex]
	case Rook:
		index := rookMagics[sq].index(occupancy)
		return bbMagicRookMoves[index]
	case Bishop:
		index := bishopMagics[sq].index(occupancy)
		return bbMagicBishopMoves[index]
	case Knight:
		return bbKnightMoves[sq]
	default:
		return emptyBitboard
	}
}

// slowMoves computes the move bitboard for each piece type.
//
// This function is intended to be used during initialization of move tables.
func slowMoves(pt PieceType, sq Square, blockers bitboard) bitboard {
	switch pt {
	case Rook:
		return rayBitboard(sq, north, blockers) | rayBitboard(sq, east, blockers) |
			rayBitboard(sq, south, blockers) | rayBitboard(sq, west, blockers)
	case Bishop:
		return rayBitboard(sq, northEast, blockers) | rayBitboard(sq, southEast, blockers) |
			rayBitboard(sq, southWest, blockers) | rayBitboard(sq, northWest, blockers)
	default:
		panic("slow moves not defined for piece type")
	}
}

func rayBitboard(sq Square, d direction, blockers bitboard) bitboard {
	m := map[direction]func(sq Square) bool{
		north:     func(sq Square) bool { return sq.Rank() == Rank8 },
		northEast: func(sq Square) bool { return sq.File() == FileH || sq.Rank() == Rank8 },
		east:      func(sq Square) bool { return sq.File() == FileH },
		southEast: func(sq Square) bool { return sq.File() == FileH || sq.Rank() == Rank1 },
		south:     func(sq Square) bool { return sq.Rank() == Rank1 },
		southWest: func(sq Square) bool { return sq.File() == FileA || sq.Rank() == Rank1 },
		west:      func(sq Square) bool { return sq.File() == FileA },
		northWest: func(sq Square) bool { return sq.File() == FileA || sq.Rank() == Rank8 },
	}

	var bb bitboard
	for !m[d](sq) {
		sq = Square(direction(sq) + d)
		bb |= sq.bitboard()
		if sq.bitboard()&blockers > 0 {
			break
		}
	}
	return bb
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
		panic("mask not defined for piece type")
	}
}
