package chess

import (
	"fmt"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board          board
	hash           Hash
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

	pos.hash = newZobristHash(pos)

	return pos, nil
}

// StartingPosition returns the starting position.
func StartingPosition() *Position {
	pos, _ := NewPosition(startFEN)
	return pos
}

// Hash returns the position Zobrist hash.
func (pos *Position) Hash() Hash {
	return pos.hash
}

// Turn returns the color of the next player to move in this position.
func (pos Position) Turn() Color {
	return pos.turn
}

// MakeMove makes a move.
//
// Checks the legality of the resulting position.
// Returns true if the move was legal and has been made.
//
// The returned metadata can be used to unmake the move and
// restore the position to the previous state.
func (pos *Position) MakeMove(m Move) (Metadata, bool) {
	switch {
	case m.HasTag(KingSideCastle) || m.HasTag(QueenSideCastle):
		if !pos.isCastleLegal(m) {
			return NoMetadata, false
		}
	case m.P1() == WhiteKing || m.P1() == BlackKing:
		if pos.isSquareAttacked(m.S2()) {
			return NoMetadata, false
		}
	default:
		if pos.isDiscoveredCheck(m) {
			return NoMetadata, false
		}
	}

	metadata := newMetadata(pos.turn, pos.castlingRights,
		pos.halfMoveClock, pos.fullMoves, pos.enPassant)
	cr := pos.castlingRights

	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbWhite&pos.board.bbPawn, pos.board.bbBlack&pos.board.bbPawn)
	}

	pos.board.makeMove(m)
	if pos.isSquareAttacked((pos.board.bbKing & pos.board.getColor(pos.turn)).scanForward()) {
		pos.board.makeMove(m)
		return NoMetadata, false
	}

	pos.turn = pos.turn.other()
	pos.castlingRights = moveCastlingRights(pos.castlingRights, m)
	pos.enPassant = moveEnPassant(m)
	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbWhite&pos.board.bbPawn, pos.board.bbBlack&pos.board.bbPawn)
	}
	pos.hash ^= xorHashPartialMove(m, cr, pos.castlingRights)

	if m.P1().Type() == Pawn || m.HasTag(Capture) {
		pos.halfMoveClock = 0
	} else {
		pos.halfMoveClock++
	}

	if pos.turn == White {
		pos.fullMoves++
	}

	return metadata, true
}

// UnmakeMove unmakes a move and restores the previous position.
func (pos *Position) UnmakeMove(m Move, meta Metadata) {
	cr := pos.castlingRights
	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbWhite&pos.board.bbPawn, pos.board.bbBlack&pos.board.bbPawn)
	}
	pos.board.makeMove(m)
	pos.turn = meta.turn()
	pos.castlingRights = meta.castleRights()
	pos.enPassant = meta.enPassant()
	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbWhite&pos.board.bbPawn, pos.board.bbBlack&pos.board.bbPawn)
	}
	pos.hash ^= xorHashPartialMove(m, cr, pos.castlingRights)
	pos.halfMoveClock = meta.halfMoveClock()
	pos.fullMoves = meta.fullMoves()
}

// InCheck returns whether the current player is in check.
func (pos *Position) InCheck() bool {
	bbKing := pos.board.getColor(pos.turn) & pos.board.bbKing
	return pos.isSquareAttacked(bbKing.scanForward())
}

// PieceMap executes the callback for each piece on the board, passing the piece
// and its square as arguments. Intended to be used in evaluation functions.
func (pos *Position) PieceMap(cb func(p Piece, sq Square)) {
	for p, bb := range [12]bitboard{
		pos.board.bbBlack & pos.board.bbPawn,
		pos.board.bbWhite & pos.board.bbPawn,
		pos.board.bbBlack & pos.board.bbKnight,
		pos.board.bbWhite & pos.board.bbKnight,
		pos.board.bbBlack & pos.board.bbBishop,
		pos.board.bbWhite & pos.board.bbBishop,
		pos.board.bbBlack & pos.board.bbRook,
		pos.board.bbWhite & pos.board.bbRook,
		pos.board.bbBlack & pos.board.bbQueen,
		pos.board.bbWhite & pos.board.bbQueen,
		pos.board.bbBlack & pos.board.bbKing,
		pos.board.bbWhite & pos.board.bbKing,
	} {
		for ; bb > 0; bb = bb.resetLSB() {
			cb(Piece(p), bb.scanForward())
		}
	}
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
	p1 := m.P1()
	if pt := p1.Type(); pt != King && pt != Rook {
		return cr
	}

	switch s1, s2 := m.S1(), m.S2(); {
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
