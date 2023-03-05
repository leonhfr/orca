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
// Returns an UCI-compatible representation.
func (c Color) String() string {
	return colorName[c : c+1]
}

// PieceType is the type of a piece.
type PieceType uint8

const (
	// Pawn represents a pawn.
	Pawn PieceType = iota << 1
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
//
// Returns an UCI-compatible representation.
func (pt PieceType) String() string {
	return pieceTypeName[pt/2 : pt/2+1]
}

// Piece is a piece type with a color.
type Piece uint8

const (
	// BlackPawn represents a black pawn.
	BlackPawn Piece = Piece(Black) | Piece(Pawn)
	// WhitePawn represents a white pawn.
	WhitePawn Piece = Piece(White) | Piece(Pawn)
	// BlackKnight represents a black knight.
	BlackKnight Piece = Piece(Black) | Piece(Knight)
	// WhiteKnight represents a white knight.
	WhiteKnight Piece = Piece(White) | Piece(Knight)
	// BlackBishop represents a black bishop.
	BlackBishop Piece = Piece(Black) | Piece(Bishop)
	// WhiteBishop represents a white bishop.
	WhiteBishop Piece = Piece(White) | Piece(Bishop)
	// BlackRook represents a black rook.
	BlackRook Piece = Piece(Black) | Piece(Rook)
	// WhiteRook represents a white rook.
	WhiteRook Piece = Piece(White) | Piece(Rook)
	// BlackQueen represents a black queen.
	BlackQueen Piece = Piece(Black) | Piece(Queen)
	// WhiteQueen represents a white queen.
	WhiteQueen Piece = Piece(White) | Piece(Queen)
	// BlackKing represents a black king.
	BlackKing Piece = Piece(Black) | Piece(King)
	// WhiteKing represents a white king.
	WhiteKing Piece = Piece(White) | Piece(King)
	// NoPiece represents an absence of Piece.
	NoPiece Piece = 12
)

const pieceName = "pPnNbBrRqQkK-"

// String implements the Stringer interface.
//
// Returns an UCI-compatible representation.
func (p Piece) String() string {
	return pieceName[p : p+1]
}

// Color returns the color of the piece.
func (p Piece) Color() Color {
	return Color(p & 1)
}

// Type returns the type of the piece.
func (p Piece) Type() PieceType {
	return PieceType(p & ^Piece(1))
}

// pieceTable is a lookup type of pieces indexed by their byte representation minus 'A'.
var pieceTable = [58]Piece{}

// initializes pieceTable
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
