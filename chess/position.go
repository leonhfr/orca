package chess

import (
	"fmt"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board          board
	attackMap      attackMap
	turn           Color
	castlingRights castlingRights
	enPassant      Square
	halfMoveClock  uint8
	fullMoves      uint8
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// NewPosition creates a position from a FEN string.
func NewPosition(fen string) (*Position, error) {
	fields := strings.Fields(strings.TrimSpace(fen))
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid fen (%s), must have 6 fields", fen)
	}

	var err error
	pos := &Position{}

	pos.board, err = fenBoard(fields[0])
	if err != nil {
		return nil, err
	}

	pos.turn, err = fenTurn(fields[1])
	if err != nil {
		return nil, err
	}

	pos.castlingRights, err = fenCastlingRights(fields[2])
	if err != nil {
		return nil, err
	}

	pos.enPassant, err = fenEnPassantSquare(fields[3])
	if err != nil {
		return nil, err
	}

	pos.halfMoveClock, err = fenHalfMoveClock(fields[4])
	if err != nil {
		return nil, err
	}

	pos.fullMoves, err = fenFullMoves(fields[5])
	if err != nil {
		return nil, err
	}

	pos.attackMap = newAttackMap(pos)

	return pos, nil
}

// StartingPosition returns the starting position.
func StartingPosition() *Position {
	pos, _ := NewPosition(startFEN)
	return pos
}

// MakeMove makes a move.
//
// Does not check the legality of the resulting position.
// The returned metadata can be used to return the position
// to the same anterior state.
func (pos *Position) MakeMove(m Move) Metadata {
	metadata := newMetadata(pos.turn, pos.castlingRights,
		pos.halfMoveClock, pos.fullMoves, pos.enPassant)

	pos.board.makeMove(m)

	if m.HasTag(EnPassant) || m.HasTag(KingSideCastle) || m.HasTag(QueenSideCastle) {
		pos.attackMap = newAttackMap(pos)
	} else {
		p := m.P1()
		if promo := m.Promo(); promo != NoPiece {
			p = promo
		}
		pos.attackMap.updateAttackMap(m.S1(), m.S2(), NoPiece, p, pos.board.bbWhite|pos.board.bbBlack,
			pos.board.bbRook|pos.board.bbQueen, pos.board.bbBishop|pos.board.bbQueen)
	}

	pos.turn = pos.turn.other()
	pos.castlingRights = moveCastlingRights(pos.castlingRights, m)
	pos.enPassant = moveEnPassant(m)

	if m.P1().Type() == Pawn || m.HasTag(Capture) {
		pos.halfMoveClock = 0
	} else {
		pos.halfMoveClock++
	}

	if pos.turn == White {
		pos.fullMoves++
	}

	return metadata
}

// UnmakeMove unmakes a move and restores the previous position.
func (pos *Position) UnmakeMove(m Move, meta Metadata) {
	pos.board.makeMove(m)
	if m.HasTag(EnPassant) || m.HasTag(KingSideCastle) || m.HasTag(QueenSideCastle) {
		pos.attackMap = newAttackMap(pos)
	} else {
		pos.attackMap.updateAttackMap(m.S1(), m.S2(), m.P1(), m.P2(), pos.board.bbWhite|pos.board.bbBlack,
			pos.board.bbRook|pos.board.bbQueen, pos.board.bbBishop|pos.board.bbQueen)
	}
	pos.turn = meta.turn()
	pos.castlingRights = meta.castleRights()
	pos.enPassant = meta.enPassant()
	pos.halfMoveClock = meta.halfMoveClock()
	pos.fullMoves = meta.fullMoves()
}

// String implements the Stringer interface.
//
// Returns a FEN formatted string.
func (pos Position) String() string {
	sq := "-"
	if pos.enPassant != NoSquare {
		sq = pos.enPassant.String()
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		pos.board.String(),
		pos.turn.String(),
		pos.castlingRights.String(),
		sq,
		pos.halfMoveClock,
		pos.fullMoves,
	)
}

// moveCastlingRights computes the new castling rights after a move.
func moveCastlingRights(cr castlingRights, m Move) castlingRights {
	switch p1, s1, s2 := m.P1(), m.S1(), m.S2(); {
	case p1 == WhiteKing:
		return cr & ^(castleWhiteKing | castleWhiteQueen)
	case p1 == BlackKing:
		return cr & ^(castleBlackKing | castleBlackQueen)
	case (p1 == WhiteRook && s1 == A1) || s2 == A1:
		return cr & ^castleWhiteQueen
	case (p1 == WhiteRook && s1 == H1) || s2 == H1:
		return cr & ^castleWhiteKing
	case (p1 == BlackRook && s1 == A8) || s2 == A8:
		return cr & ^castleBlackQueen
	case (p1 == BlackRook && s1 == H8) || s2 == H8:
		return cr & ^castleBlackKing
	default:
		return cr
	}
}

// moveEnPassant computes the en passant square after a move.
func moveEnPassant(m Move) Square {
	if m.P1().Type() != Pawn {
		return NoSquare
	}

	switch c, s1, s2 := m.P1().Color(), m.S1(), m.S2(); {
	case c == White &&
		s1.bitboard()&bbRank2 > 0 &&
		s2.bitboard()&bbRank4 > 0:
		return s2 - 8
	case c == Black &&
		s1.bitboard()&bbRank7 > 0 &&
		s2.bitboard()&bbRank5 > 0:
		return s2 + 8
	default:
		return NoSquare
	}
}
