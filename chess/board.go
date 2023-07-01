package chess

import "strings"

// board represents a chess board.
type board struct {
	bbPieces [6]bitboard
	bbColors [2]bitboard
	sqKings  [2]Square
}

// newBoard creates a new board.
func newBoard(m map[Square]Piece) board {
	b := board{}
	for sq, p := range m {
		bb := sq.bitboard()
		b.bbPieces[p.Type()] ^= bb
		b.bbColors[p.Color()] ^= bb

		if p.Type() == King {
			b.sqKings[p.Color()] = sq
		}
	}
	return b
}

// pieceAt returns the piece, if any, present at the given square.
func (b board) pieceAt(sq Square) Piece {
	switch bb := sq.bitboard(); {
	case b.bbColors[White]&bb > 0:
		return b.pieceByColor(sq, White)
	case b.bbColors[Black]&bb > 0:
		return b.pieceByColor(sq, Black)
	default:
		return NoPiece
	}
}

// pieceByColor returns the piece present at the given square.
func (b board) pieceByColor(sq Square, c Color) Piece {
	bb := sq.bitboard()
	for pt := Pawn; pt <= King; pt++ {
		if b.bbPieces[pt]&bb > 0 {
			return pt.color(c)
		}
	}
	return NoPiece
}

// makeMove makes and unmakes a move on the board.
func (b *board) makeMove(m Move, cf [2]File) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()
	c := p1.Color()

	s1bb, s2bb := s1.bitboard(), s2.bitboard()
	mbb := s1bb ^ s2bb

	if promo := m.Promo(); promo == NoPiece {
		b.bbPieces[p1.Type()] ^= mbb
	} else {
		// promotion
		b.bbPieces[p1.Type()] ^= s1bb
		b.bbPieces[promo.Type()] ^= s2bb
	}

	b.bbColors[c] ^= mbb

	if p1.Type() == King {
		b.sqKings[c] = s2
	}

	if m.HasTag(Quiet) {
		return
	}

	switch enPassant := m.HasTag(EnPassant); {
	case p2 != NoPiece && !enPassant: // capture
		b.bbPieces[p2.Type()] ^= s2bb
		b.bbColors[p2.Color()] ^= s2bb
	case c == White && enPassant: // white en passant
		bb := s2.bitboard().southOne()
		b.bbPieces[Pawn] ^= bb
		b.bbColors[Black] ^= bb
	case c == Black && enPassant: // black en passant
		bb := s2.bitboard().northOne()
		b.bbPieces[Pawn] ^= bb
		b.bbColors[White] ^= bb
	case c == Black && m.HasTag(ASideCastle): // black A side castle
		bb := newSquare(cf[aSide], Rank8).bitboard() | D8.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[Black] ^= bb
	case c == Black && m.HasTag(HSideCastle): // black H side castle
		bb := newSquare(cf[hSide], Rank8).bitboard() | F8.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[Black] ^= bb
	case c == White && m.HasTag(ASideCastle): // white A side castle
		bb := newSquare(cf[aSide], Rank1).bitboard() | D1.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[White] ^= bb
	case c == White && m.HasTag(HSideCastle): // white H side castle
		bb := newSquare(cf[hSide], Rank1).bitboard() | F1.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[White] ^= bb
	}
}

// unmakeMove unmakes a move on the board.
func (b *board) unmakeMove(m Move, cf [2]File) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()
	c := p1.Color()

	s1bb, s2bb := s1.bitboard(), s2.bitboard()
	mbb := s1bb ^ s2bb

	if promo := m.Promo(); promo == NoPiece {
		b.bbPieces[p1.Type()] ^= mbb
	} else {
		// promotion
		b.bbPieces[p1.Type()] ^= s1bb
		b.bbPieces[promo.Type()] ^= s2bb
	}

	b.bbColors[c] ^= mbb

	if p1.Type() == King {
		b.sqKings[c] = s1
	}

	if m.HasTag(Quiet) {
		return
	}

	switch enPassant := m.HasTag(EnPassant); {
	case p2 != NoPiece && !enPassant: // capture
		b.bbPieces[p2.Type()] ^= s2bb
		b.bbColors[p2.Color()] ^= s2bb
	case c == White && enPassant: // white en passant
		bb := s2.bitboard().southOne()
		b.bbPieces[Pawn] ^= bb
		b.bbColors[Black] ^= bb
	case c == Black && enPassant: // black en passant
		bb := s2.bitboard().northOne()
		b.bbPieces[Pawn] ^= bb
		b.bbColors[White] ^= bb
	case c == Black && m.HasTag(ASideCastle): // black A side castle
		bb := newSquare(cf[aSide], Rank8).bitboard() | D8.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[Black] ^= bb
	case c == Black && m.HasTag(HSideCastle): // black H side castle
		bb := newSquare(cf[hSide], Rank8).bitboard() | F8.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[Black] ^= bb
	case c == White && m.HasTag(ASideCastle): // white A side castle
		bb := newSquare(cf[aSide], Rank1).bitboard() | D1.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[White] ^= bb
	case c == White && m.HasTag(HSideCastle): // white H side castle
		bb := newSquare(cf[hSide], Rank1).bitboard() | F1.bitboard()
		b.bbPieces[Rook] ^= bb
		b.bbColors[White] ^= bb
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
