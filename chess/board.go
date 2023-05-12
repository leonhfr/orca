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
	switch bb := sq.bitboard(); {
	case b.bbWhite&bb > 0:
		return b.pieceByColor(sq, White)
	case b.bbBlack&bb > 0:
		return b.pieceByColor(sq, Black)
	default:
		return NoPiece
	}
}

// pieceByColor returns the piece present at the given square.
func (b board) pieceByColor(sq Square, c Color) Piece {
	switch bb := sq.bitboard(); {
	case b.bbPawn&bb > 0:
		return Pawn.color(c)
	case b.bbKnight&bb > 0:
		return Knight.color(c)
	case b.bbBishop&bb > 0:
		return Bishop.color(c)
	case b.bbRook&bb > 0:
		return Rook.color(c)
	case b.bbQueen&bb > 0:
		return Queen.color(c)
	case b.bbKing&bb > 0:
		return King.color(c)
	default:
		return NoPiece
	}
}

// makeMove makes and unmakes a move on the board.
func (b *board) makeMove(m Move) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()
	c := p1.Color()

	s1bb, s2bb := s1.bitboard(), s2.bitboard()
	mbb := s1bb ^ s2bb

	if promo := m.Promo(); promo == NoPiece {
		b.xorBitboard(p1.Type(), mbb)
	} else {
		// promotion
		b.xorBitboard(p1.Type(), s1bb)
		b.xorBitboard(promo.Type(), s2bb)
	}

	b.xorColor(c, mbb)

	if m.HasTag(Quiet) {
		return
	}

	switch enPassant := m.HasTag(EnPassant); {
	case p2 != NoPiece && !enPassant: // capture
		b.xorBitboard(p2.Type(), s2bb)
		b.xorColor(p2.Color(), s2bb)
	case c == White && enPassant: // white en passant
		bb := s2.bitboard() >> 8
		b.bbPawn ^= bb
		b.bbBlack ^= bb
	case c == Black && enPassant: // black en passant
		bb := s2.bitboard() << 8
		b.bbPawn ^= bb
		b.bbWhite ^= bb
	case c == White && m.HasTag(KingSideCastle): // white king side castle
		b.bbRook ^= bbWhiteKingCastle
		b.bbWhite ^= bbWhiteKingCastle
	case c == White && m.HasTag(QueenSideCastle): // white queen side castle
		b.bbRook ^= bbWhiteQueenCastle
		b.bbWhite ^= bbWhiteQueenCastle
	case c == Black && m.HasTag(KingSideCastle): // black king side castle
		b.bbRook ^= bbBlackKingCastle
		b.bbBlack ^= bbBlackKingCastle
	case c == Black && m.HasTag(QueenSideCastle): // black queen side castle
		b.bbRook ^= bbBlackQueenCastle
		b.bbBlack ^= bbBlackQueenCastle
	}
}

// getBitboard returns the bitboard of the given piece type and color.
func (b board) getBitboard(pt PieceType, c Color) bitboard {
	bbColor := b.bbWhite
	if c == Black {
		bbColor = b.bbBlack
	}
	switch pt {
	case King:
		return bbColor & b.bbKing
	case Queen:
		return bbColor & b.bbQueen
	case Rook:
		return bbColor & b.bbRook
	case Bishop:
		return bbColor & b.bbBishop
	case Knight:
		return bbColor & b.bbKnight
	case Pawn:
		return bbColor & b.bbPawn
	default:
		panic("unknown piece")
	}
}

// getColor returns the bitboard of the given color.
func (b board) getColor(c Color) bitboard {
	if c == White {
		return b.bbWhite
	}
	return b.bbBlack
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
