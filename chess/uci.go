package chess

// UCI is the UCI move notation.
type UCI struct{}

// Encode encodes a move into a UCI compatible string.
//
// Implements the MoveNotation interface.
func (UCI) Encode(_ *Position, m Move) string {
	if m == NoMove {
		return "null"
	}

	base := m.S1().String() + m.S2().String()
	if promo := m.Promo(); promo != NoPiece {
		base += promo.Type().String()
	}
	return base
}

// Decode decodes a move from a UCI string.
//
// Implements the MoveNotation interface.
func (UCI) Decode(pos *Position, s string) (Move, error) {
	if pos == nil {
		return 0, errMissingPosition
	}

	if len(s) < 4 || len(s) > 5 {
		return 0, errInvalidMove
	}

	s1, err := uciSquare(s[0:2])
	if err != nil {
		return 0, errInvalidMove
	}
	s2, err := uciSquare(s[2:4])
	if err != nil {
		return 0, errInvalidMove
	}

	promo := NoPiece
	if len(s) == 5 {
		r := []byte(s)[4]
		if !('A' <= r && r <= 'z') {
			return 0, errInvalidMove
		}
		promoType := promoPieceTypeTable[r-'A']
		if promoType == NoPieceType {
			return 0, errInvalidMove
		}
		promo = promoType.color(pos.turn)
	}

	p1 := pos.board.pieceAt(s1)
	p2 := pos.board.pieceAt(s2)
	return newMove(p1, p2, s1, s2, pos.enPassant, promo), nil
}
