package chess

// UCI is the UCI move notation.
type UCI struct{}

// Encode encodes a move into a UCI compatible string.
//
// Implements the MoveNotation interface.
func (UCI) Encode(_ *Position, m Move) string {
	switch {
	case m == NoMove:
		return "null"
	case m.HasTag(Promotion):
		return m.S1().String() + m.S2().String() + m.Promo().Type().String()
	default:
		return m.S1().String() + m.S2().String()
	}
}

// Decode decodes a move from a UCI string.
//
// Implements the MoveNotation interface.
func (UCI) Decode(pos *Position, s string) (Move, error) {
	if pos == nil {
		return 0, errMissingPosition
	}

	s1, s2, promo, err := uciMove(s, pos.turn)
	if err != nil {
		return NoMove, err
	}

	p1 := pos.board.pieceAt(s1)
	p2 := pos.board.pieceAt(s2)
	return newMove(p1, p2, s1, s2, pos.enPassant, promo), nil
}

// UCIChess960 is the UCI move notation compatible with Chess960.
//
// Instead of encoding castling moves as 2 spaces king jumps, it encodes them as
// the king taking hiw own rook.
type UCIChess960 struct{}

// Encode encodes a move into a UCI Chess960 compatible string.
//
// Implements the MoveNotation interface.
func (UCIChess960) Encode(pos *Position, m Move) string {
	switch {
	case m == NoMove:
		return "null"
	case m.HasTag(ASideCastle):
		return m.S1().String() + newSquare(pos.castling.files[aSide], m.S1().Rank()).String()
	case m.HasTag(HSideCastle):
		return m.S1().String() + newSquare(pos.castling.files[hSide], m.S1().Rank()).String()
	case m.HasTag(Promotion):
		return m.S1().String() + m.S2().String() + m.Promo().Type().String()
	default:
		return m.S1().String() + m.S2().String()
	}
}

// Decode decodes a move from a UCI Chess960 string.
//
// Implements the MoveNotation interface.
func (UCIChess960) Decode(pos *Position, s string) (Move, error) {
	if pos == nil {
		return 0, errMissingPosition
	}

	s1, s2, promo, err := uciMove(s, pos.turn)
	if err != nil {
		return NoMove, err
	}

	p1 := pos.board.pieceAt(s1)
	p2 := pos.board.pieceAt(s2)
	return newChess960Move(p1, p2, s1, s2, pos.enPassant, promo, pos.castling.files), nil
}

// uciMove parses a uci move.
func uciMove(s string, turn Color) (Square, Square, Piece, error) {
	if len(s) < 4 || len(s) > 5 {
		return NoSquare, NoSquare, NoPiece, errInvalidMove
	}

	s1, err := uciSquare(s[0:2])
	if err != nil {
		return NoSquare, NoSquare, NoPiece, errInvalidMove
	}
	s2, err := uciSquare(s[2:4])
	if err != nil {
		return NoSquare, NoSquare, NoPiece, errInvalidMove
	}

	promo := NoPiece
	if len(s) == 5 {
		r := []byte(s)[4]
		if r < 'A' || r > 'z' {
			return NoSquare, NoSquare, NoPiece, errInvalidMove
		}
		promoType := promoPieceTypeTable[r-'A']
		if promoType == NoPieceType {
			return NoSquare, NoSquare, NoPiece, errInvalidMove
		}
		promo = promoType.color(turn)
	}

	return s1, s2, promo, nil
}
