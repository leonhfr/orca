package chess

// Metadata represents a position's metadata.
//
//	32 bits
//	__ square fullMove halfMove ___ CCCC T
//	square      en passant square
//	fullMove    full moves
//	halfMove    half move clock
//	CCCC        castle rights
//	T           turn color
//	_           unused bit
type Metadata uint32

// newMetadata create a new metadata.
func newMetadata(c Color, cr castlingRights, halfMoveClock, fullMoves uint8, enPassant Square) Metadata {
	return Metadata(c) |
		Metadata(cr)<<1 |
		Metadata(halfMoveClock)<<8 |
		Metadata(fullMoves)<<16 |
		Metadata(enPassant)<<24
}

// turn returns the turn color.
func (m Metadata) turn() Color {
	return Color(m & 1)
}

// castleRights returns the castle rights.
func (m Metadata) castleRights() castlingRights {
	return castlingRights((m >> 1) & 15)
}

// halfMoveClock returns the half move clock.
func (m Metadata) halfMoveClock() uint8 {
	return uint8(m >> 8)
}

// fullMoves returns the full moves.
func (m Metadata) fullMoves() uint8 {
	return uint8(m >> 16)
}

// enPassant returns the en passant square.
func (m Metadata) enPassant() Square {
	return Square(m >> 24)
}

// side represents a side of the board.
type side uint8

const (
	// kingSide represents the kings' side.
	kingSide side = iota
	// queenSide represents the queens' side.
	queenSide
)

// castlingRights represents the castling right of one combination of side and color.
type castlingRights uint8

const (
	// castleWhiteKing represents white's king castle.
	castleWhiteKing castlingRights = 1 << iota
	// castleWhiteQueen represents white's queen castle.
	castleWhiteQueen
	// castleBlackKing represents black's king castle.
	castleBlackKing
	// castleBlackQueen represents black's queen castle.
	castleBlackQueen
	// noCastle represents the absence of a castle.
	noCastle castlingRights = 0
)

// canCastle returns whether a castle with this combinations of
// color and side is possible.
func (cr castlingRights) canCastle(c Color, s side) bool {
	switch {
	case c == White && s == kingSide:
		return (cr & castleWhiteKing) > 0
	case c == White && s == queenSide:
		return (cr & castleWhiteQueen) > 0
	case c == Black && s == kingSide:
		return (cr & castleBlackKing) > 0
	case c == Black && s == queenSide:
		return (cr & castleBlackQueen) > 0
	default:
		return false
	}
}

// String implements the Stringer interface.
//
// Returns an UCI-compatible representation.
func (cr castlingRights) String() string {
	if cr == noCastle {
		return "-"
	}

	var rights string
	if cr.canCastle(White, kingSide) {
		rights += "K"
	}
	if cr.canCastle(White, queenSide) {
		rights += "Q"
	}
	if cr.canCastle(Black, kingSide) {
		rights += "k"
	}
	if cr.canCastle(Black, queenSide) {
		rights += "q"
	}
	return rights
}
