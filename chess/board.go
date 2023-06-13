package chess

import "strings"

// board represents a chess board.
type board struct {
	bbPieces [6]bitboard
	bbColors [2]bitboard
}

// newBoard creates a new board.
func newBoard(m map[Square]Piece) board {
	b := board{}
	for sq, p := range m {
		bb := sq.bitboard()
		b.bbPieces[p.Type()] ^= bb
		b.bbColors[p.Color()] ^= bb
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
	switch bb := sq.bitboard(); {
	case b.bbPieces[Pawn]&bb > 0:
		return Pawn.color(c)
	case b.bbPieces[Knight]&bb > 0:
		return Knight.color(c)
	case b.bbPieces[Bishop]&bb > 0:
		return Bishop.color(c)
	case b.bbPieces[Rook]&bb > 0:
		return Rook.color(c)
	case b.bbPieces[Queen]&bb > 0:
		return Queen.color(c)
	case b.bbPieces[King]&bb > 0:
		return King.color(c)
	default:
		return NoPiece
	}
}

// kingSquare returns the king's square.
func (b *board) kingSquare(c Color) Square {
	if c == White {
		return (b.bbColors[White] & b.bbPieces[King]).scanForward()
	}
	return (b.bbColors[Black] & b.bbPieces[King]).scanForward()
}

// makeMove makes and unmakes a move on the board.
func (b *board) makeMove(m Move) {
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
	case c == White && m.HasTag(KingSideCastle): // white king side castle
		b.bbPieces[Rook] ^= bbWhiteKingCastle
		b.bbColors[White] ^= bbWhiteKingCastle
	case c == White && m.HasTag(QueenSideCastle): // white queen side castle
		b.bbPieces[Rook] ^= bbWhiteQueenCastle
		b.bbColors[White] ^= bbWhiteQueenCastle
	case c == Black && m.HasTag(KingSideCastle): // black king side castle
		b.bbPieces[Rook] ^= bbBlackKingCastle
		b.bbColors[Black] ^= bbBlackKingCastle
	case c == Black && m.HasTag(QueenSideCastle): // black queen side castle
		b.bbPieces[Rook] ^= bbBlackQueenCastle
		b.bbColors[Black] ^= bbBlackQueenCastle
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
