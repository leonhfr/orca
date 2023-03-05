package chess

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
