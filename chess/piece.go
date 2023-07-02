package chess

// Color represents the color of a chess piece.
type Color uint8

const (
	// Black represents the black color.
	Black Color = iota
	// White represents the white color.
	White
)

const colorName = "bw"

// String implements the Stringer interface.
//
// Returns a FEN compatible representation.
func (c Color) String() string {
	return colorName[c : c+1]
}

// Other returns the Other color.
func (c Color) Other() Color {
	return Color((c + 1) & 1)
}

// PieceType is the type of a piece.
type PieceType uint8

const (
	// Pawn represents a pawn.
	Pawn PieceType = iota
	// Knight represents a knight.
	Knight
	// Bishop represents a bishop.
	Bishop
	// Rook represents a rook.
	Rook
	// Queen represents a queen.
	Queen
	// King represents a king.
	King
	// NoPieceType represents an absence of PieceType.
	NoPieceType
)

const pieceTypeName = "pnbrqk-"

// String implements the Stringer interface.
func (pt PieceType) String() string {
	return pieceTypeName[pt : pt+1]
}

// color returns a piece of the passed piece type and color.
func (pt PieceType) color(c Color) Piece {
	return Piece(pt<<1) ^ Piece(c)
}

// Piece is a piece type with a color.
type Piece uint8

const (
	// BlackPawn represents a black pawn.
	BlackPawn Piece = iota
	// WhitePawn represents a white pawn.
	WhitePawn
	// BlackKnight represents a black knight.
	BlackKnight
	// WhiteKnight represents a white knight.
	WhiteKnight
	// BlackBishop represents a black bishop.
	BlackBishop
	// WhiteBishop represents a white bishop.
	WhiteBishop
	// BlackRook represents a black rook.
	BlackRook
	// WhiteRook represents a white rook.
	WhiteRook
	// BlackQueen represents a black queen.
	BlackQueen
	// WhiteQueen represents a white queen.
	WhiteQueen
	// BlackKing represents a black king.
	BlackKing
	// WhiteKing represents a white king.
	WhiteKing
	// NoPiece represents an absence of Piece.
	NoPiece
)

const pieceName = "pPnNbBrRqQkK-"

// String implements the Stringer interface.
//
// Returns a FEN compatible representation.
func (p Piece) String() string {
	return pieceName[p : p+1]
}

// Color returns the color of the piece.
func (p Piece) Color() Color {
	return Color(p & 1)
}

// Type returns the type of the piece.
func (p Piece) Type() PieceType {
	return PieceType(p >> 1)
}

var (
	// pieceTable is a lookup type of pieces indexed by their byte representation minus 'A'.
	pieceTable = [58]Piece{}
	// promoPieceTypeTable is a lookup type of piece types indexed by their byte representation minus 'A'.
	promoPieceTypeTable = [58]PieceType{}
)

// initializes pieceTable.
func initPieceTable() {
	m := map[rune]Piece{
		'K': WhiteKing, 'k': BlackKing,
		'Q': WhiteQueen, 'q': BlackQueen,
		'R': WhiteRook, 'r': BlackRook,
		'B': WhiteBishop, 'b': BlackBishop,
		'N': WhiteKnight, 'n': BlackKnight,
		'P': WhitePawn, 'p': BlackPawn,
	}

	for r := 'A'; r <= 'z'; r++ {
		if p, ok := m[r]; ok {
			pieceTable[r-'A'] = p
		} else {
			pieceTable[r-'A'] = NoPiece
		}
	}
}

// initializes promoPieceTypeTable.
func initPromoPieceTypeTable() {
	m := map[rune]PieceType{
		'q': Queen,
		'r': Rook,
		'b': Bishop,
		'n': Knight,
	}

	for r := 'A'; r <= 'z'; r++ {
		if pt, ok := m[r]; ok {
			promoPieceTypeTable[r-'A'] = pt
		} else {
			promoPieceTypeTable[r-'A'] = NoPieceType
		}
	}
}
