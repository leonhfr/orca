package chess

import "strings"

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

// pieceAt returns the piece, if any, present at the given square.
func (b board) pieceAt(sq Square) Piece {
	p := BlackPawn
	if b.bbWhite&sq.bitboard() > 0 {
		p = WhitePawn
	}
	for bb := sq.bitboard(); p <= WhiteKing; p += 2 {
		if (b.getBitboard(p) & bb) > 0 {
			return p
		}
	}
	return NoPiece
}

// getBitboard returns the bitboard of the given piece.
func (b board) getBitboard(p Piece) bitboard {
	switch p {
	case WhiteKing:
		return b.bbWhite & b.bbKing
	case WhiteQueen:
		return b.bbWhite & b.bbQueen
	case WhiteRook:
		return b.bbWhite & b.bbRook
	case WhiteBishop:
		return b.bbWhite & b.bbBishop
	case WhiteKnight:
		return b.bbWhite & b.bbKnight
	case WhitePawn:
		return b.bbWhite & b.bbPawn
	case BlackKing:
		return b.bbBlack & b.bbKing
	case BlackQueen:
		return b.bbBlack & b.bbQueen
	case BlackRook:
		return b.bbBlack & b.bbRook
	case BlackBishop:
		return b.bbBlack & b.bbBishop
	case BlackKnight:
		return b.bbBlack & b.bbKnight
	case BlackPawn:
		return b.bbBlack & b.bbPawn
	default:
		panic("unknown piece")
	}
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

// String implements the Stringer interface.
//
// Returns an UCI-compatible representation.
func (b board) String() string {
	var fields []string
	for rank := 7; rank >= 0; rank-- {
		var field []byte
		for file := FileA; file <= FileH; file++ {
			switch p := b.pieceAt(newSquare(file, Rank(rank))); {
			case p != NoPiece:
				field = append(field, []byte(p.String())...)
			case len(field) == 0:
				field = append(field, '1')
			case '8' < field[len(field)-1]:
				field = append(field, '1')
			default:
				field[len(field)-1]++
			}
		}
		fields = append(fields, string(field))
	}
	return strings.Join(fields, "/")
}
