package chess

// StaticExchange performs a static exchange on a move.
//
// Invokes the callback with the piece type of the least valuable piece,
// alternating colors and starting from the color of p1.
// If the callback returns false, the function exits early.
//
// Intended to be used for static exchange evaluation.
func (pos *Position) StaticExchange(m Move, cb func(pt PieceType) bool) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()

	bbOccupancy := pos.board.bbColors[Black] ^ pos.board.bbColors[White]
	bbMayXray := pos.board.bbPieces[Pawn] ^
		pos.board.bbPieces[Bishop] ^
		pos.board.bbPieces[Rook] ^
		pos.board.bbPieces[Queen]
	bbFrom := s1.bitboard()
	bb := pos.attackedAndDefendedByBitboard(s2, bbOccupancy)

	c := p2.Color()
	pt := p1.Type()

	for ; bbFrom > 0; c = c.Other() {
		if cb(pt) {
			break
		}

		bb ^= bbFrom
		bbOccupancy ^= bbFrom

		if bbFrom&bbMayXray > 0 {
			bb |= pos.xrayAttackedByBitboard(s2, bbOccupancy)
		}

		pt, bbFrom = pos.leastValuablePieceBitboard(bb, c)
	}
}

// leastValuablePieceBitboard returns the piece type and the bitboard of the least valuable
// piece from the passed set from color c.
func (pos *Position) leastValuablePieceBitboard(bb bitboard, c Color) (PieceType, bitboard) {
	for pt := Pawn; pt <= King; pt++ {
		subset := bb & pos.board.bbColors[c] & pos.board.bbPieces[pt]
		if subset > 0 {
			return pt, subset & -subset
		}
	}
	return Pawn, bbEmpty
}

// xrayAttackedByBitboard returns the bitboard of all sliding attacking pieces from the given occupancy.
func (pos *Position) xrayAttackedByBitboard(sq Square, bbOccupancy bitboard) bitboard {
	bbRookMoves := bbMagicRookMoves[rookMagics[sq].index(bbOccupancy)]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[sq].index(bbOccupancy)]
	bb := (pos.board.bbPieces[Queen] | pos.board.bbPieces[Rook]) & bbRookMoves & bbOccupancy
	bb |= (pos.board.bbPieces[Queen] | pos.board.bbPieces[Bishop]) & bbBishopMoves & bbOccupancy
	return bb
}
