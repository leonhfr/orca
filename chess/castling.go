package chess

// side represents a side of the board.
type side uint8

const (
	// aSide represents the A file side.
	aSide side = iota
	// hSide represents the H file side.
	hSide
)

// castlingRights represents the castling right of one combination of side and color.
type castlingRights uint8

// noCastle represents the absence of castling rights.
const noCastle castlingRights = 0

const (
	// castleBlackA represents the black A side castle (queen side in classic chess).
	castleBlackA castlingRights = 1 << iota
	// castleBlackH represents the black H side castle (king side in classic chess).
	castleBlackH
	// castleWhiteA represents the white A side castle (queen side in classic chess).
	castleWhiteA
	// castleWhiteH represents the white H side castle (king side in classic chess).
	castleWhiteH
)

// canCastle returns whether a castle with this combinations of
// color and side is possible.
func (cr castlingRights) canCastle(c Color, s side) bool {
	switch {
	case c == Black && s == aSide:
		return (cr & castleBlackA) > 0
	case c == Black && s == hSide:
		return (cr & castleBlackH) > 0
	case c == White && s == aSide:
		return (cr & castleWhiteA) > 0
	case c == White && s == hSide:
		return (cr & castleWhiteH) > 0
	default:
		return false
	}
}

// castling represents the rook files and the castling rights.
type castling struct {
	files  [2]File
	rights castlingRights
}

// String implements the Stringer interface.
//
// Returns a FEN compatible representation.
func (c castling) String() string {
	if c.rights == noCastle {
		return "-"
	}

	var rights string
	if c.rights.canCastle(White, hSide) {
		rights += "K"
	}
	if c.rights.canCastle(White, aSide) {
		rights += "Q"
	}
	if c.rights.canCastle(Black, hSide) {
		rights += "k"
	}
	if c.rights.canCastle(Black, aSide) {
		rights += "q"
	}
	return rights
}

// castleCheck contains all the data needed to check castling moves.
type castleCheck struct {
	bbKingTravel    bitboard
	bbRookTravel    bitboard
	bbNoEnemyPawn   bitboard
	bbNoEnemyKnight bitboard
	bbNoEnemyKing   bitboard
	bbNoCheck       bitboard // check for rook and bishop check.
	king1           Square
	king2           Square
	rook1           Square
	rook2           Square
}

// newCastleCheck creates a new castle check.
func newCastleCheck(c Color, s side, kings [2]Square, cf [2]File, cr castlingRights) castleCheck {
	if !cr.canCastle(c, s) {
		return castleCheck{}
	}

	king1 := kings[c]
	king2 := newSquare(kingFinalFile[s], castleRank[c])
	rook1 := newSquare(cf[s], castleRank[c])
	rook2 := newSquare(rookFinalFile[s], castleRank[c])

	bbKingTravel := (bbInBetweens[king1][king2] | (king2.bitboard() & ^king1.bitboard())) & ^rook1.bitboard()
	bbRookTravel := (bbInBetweens[rook1][rook2] | (rook2.bitboard() & ^rook1.bitboard())) & ^king1.bitboard()
	bbNoCheck := bbInBetweens[king1][king2] | king1.bitboard() | king2.bitboard()

	var bbNoEnemyPawn, bbNoEnemyKnight, bbNoEnemyKing bitboard
	for bb := bbNoCheck; bb > 0; bb = bb.resetLSB() {
		sq := bb.scanForward()
		bbNoEnemyPawn |= singlePawnCaptureBitboard(sq, c)
		bbNoEnemyKnight |= bbKnightMoves[sq]
		bbNoEnemyKing |= bbKingMoves[sq]
	}

	return castleCheck{
		bbKingTravel:    bbKingTravel,
		bbRookTravel:    bbRookTravel & ^king1.bitboard(),
		bbNoEnemyPawn:   bbNoEnemyPawn,
		bbNoEnemyKnight: bbNoEnemyKnight,
		bbNoEnemyKing:   bbNoEnemyKing,
		bbNoCheck:       bbNoCheck,
		king1:           king1,
		king2:           king2,
		rook1:           rook1,
		rook2:           rook2,
	}
}

var (
	kingFinalFile = [2]File{FileC, FileG} // indexed by side.
	rookFinalFile = [2]File{FileD, FileF} // indexed by side.
	castleRank    = [2]Rank{Rank8, Rank1} // indexed by color.
)
