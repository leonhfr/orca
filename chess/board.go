package chess

// board represents a chess board.
type board struct {
	bbKing   bitboard
	bbQueen  bitboard
	bbRook   bitboard
	bbBishop bitboard
	bbKnight bitboard
	bbPawn   bitboard
	bbWhite  bitboard
	bbBlack  bitboard
}

// newBoard creates a new board.
func newBoard(m map[Square]Piece) board {
	b := board{}
	for sq, p := range m {
		bb := sq.bitboard()
		b.xorBitboard(p.Type(), bb)
		b.xorColor(p.Color(), bb)
	}
	return b
}

// xorBitboard performs a xor operation on one of the piece bitboard.
func (b *board) xorBitboard(pt PieceType, bb bitboard) {
	switch pt {
	case King:
		b.bbKing ^= bb
	case Queen:
		b.bbQueen ^= bb
	case Rook:
		b.bbRook ^= bb
	case Bishop:
		b.bbBishop ^= bb
	case Knight:
		b.bbKnight ^= bb
	case Pawn:
		b.bbPawn ^= bb
	}
}

// xorColor performs a xor operation on one of the color bitboard.
func (b *board) xorColor(c Color, bb bitboard) {
	switch c {
	case White:
		b.bbWhite ^= bb
	case Black:
		b.bbBlack ^= bb
	}
}
