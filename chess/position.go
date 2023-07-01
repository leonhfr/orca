package chess

import (
	"fmt"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board         board
	hash          Hash
	pawnHash      Hash
	castleChecks  [4]castleCheck
	castling      castling
	turn          Color
	enPassant     Square
	halfMoveClock uint8
	fullMoves     uint8
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

	files, err := fenCastlingFiles(fields[2])
	if err != nil {
		return nil, err
	}

	rights, err := fenCastlingRights(fields[2])
	if err != nil {
		return nil, err
	}

	pos.castling = castling{files, rights}

	for c := Black; c <= White; c++ {
		for s := aSide; s <= hSide; s++ {
			pos.castleChecks[2*uint8(c)+uint8(s)] = newCastleCheck(c, s, pos.board.sqKings, files, rights)
		}
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
	pos.pawnHash = newPawnZobristHash(pos)

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

// PawnHash returns the position pawn Zobrist hash.
func (pos *Position) PawnHash() Hash {
	return pos.pawnHash
}

// Turn returns the color of the next player to move in this position.
func (pos Position) Turn() Color {
	return pos.turn
}

// FullMoves returns the number of full moves.
func (pos Position) FullMoves() uint8 {
	return pos.fullMoves
}

// MakeMove makes a move.
//
// Checks the legality of the resulting position.
// Returns true if the move was legal and has been made.
//
// The returned metadata can be used to unmake the move and
// restore the position to the previous state.
func (pos *Position) MakeMove(m Move) bool {
	if (m.HasTag(ASideCastle) || m.HasTag(HSideCastle)) && !pos.isCastleLegal(m) {
		return false
	}

	cr := pos.castling.rights

	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbColors[White]&pos.board.bbPieces[Pawn], pos.board.bbColors[Black]&pos.board.bbPieces[Pawn])
	}

	pos.board.makeMove(m, pos.castling.files)
	if pos.isSquareAttacked(pos.board.sqKings[pos.turn]) {
		pos.board.unmakeMove(m, pos.castling.files)
		return false
	}

	pos.turn = pos.turn.Other()
	pos.castling.rights = moveCastlingRights(pos.castling.rights, pos.castling.files, m)
	pos.enPassant = moveEnPassant(m)
	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbColors[White]&pos.board.bbPieces[Pawn], pos.board.bbColors[Black]&pos.board.bbPieces[Pawn])
	}

	partialHash, partialPawnHash := xorHashPartialMove(m, cr, pos.castling.rights, pos.castling.files)
	pos.hash ^= partialHash
	pos.pawnHash ^= partialPawnHash

	if m.P1().Type() == Pawn || m.HasTag(Capture) {
		pos.halfMoveClock = 0
	} else {
		pos.halfMoveClock++
	}

	if pos.turn == White {
		pos.fullMoves++
	}

	return true
}

// UnmakeMove unmakes a move and restores the previous position.
func (pos *Position) UnmakeMove(m Move, meta Metadata, hash, pawnHash Hash) {
	pos.board.unmakeMove(m, pos.castling.files)
	pos.turn = meta.turn()
	pos.castling.rights = meta.castleRights()
	pos.enPassant = meta.enPassant()
	pos.halfMoveClock = meta.halfMoveClock()
	pos.fullMoves = meta.fullMoves()
	pos.hash = hash
	pos.pawnHash = pawnHash
}

// MakeNullMove makes a null (passing) move.
func (pos *Position) MakeNullMove() {
	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbColors[White]&pos.board.bbPieces[Pawn],
			pos.board.bbColors[Black]&pos.board.bbPieces[Pawn])
	}

	pos.turn = pos.turn.Other()
	pos.enPassant = NoSquare
	pos.hash ^= polyTurn
}

// UnmakeNullMove unmakes a null move.
func (pos *Position) UnmakeNullMove(meta Metadata, hash Hash) {
	pos.turn = meta.turn()
	pos.enPassant = meta.enPassant()
	pos.hash = hash
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
		pos.castling.String(),
		sq,
		pos.halfMoveClock,
		pos.fullMoves,
	)
}

// moveCastlingRights computes the new castling rights after a move.
func moveCastlingRights(cr castlingRights, cf [2]File, m Move) castlingRights {
	p1 := m.P1()
	if pt := p1.Type(); pt != King && pt != Rook && !m.HasTag(Capture) {
		return cr
	}

	blackRookA := newSquare(cf[aSide], Rank8)
	blackRookH := newSquare(cf[hSide], Rank8)
	whiteRookA := newSquare(cf[aSide], Rank1)
	whiteRookH := newSquare(cf[hSide], Rank1)

	switch s1, s2 := m.S1(), m.S2(); {
	case p1 == BlackKing:
		return cr & ^(castleBlackA | castleBlackH)
	case p1 == WhiteKing:
		return cr & ^(castleWhiteA | castleWhiteH)
	case (p1 == BlackRook && s1 == blackRookA) || s2 == blackRookA:
		return cr & ^castleBlackA
	case (p1 == BlackRook && s1 == blackRookH) || s2 == blackRookH:
		return cr & ^castleBlackH
	case (p1 == WhiteRook && s1 == whiteRookA) || s2 == whiteRookA:
		return cr & ^castleWhiteA
	case (p1 == WhiteRook && s1 == whiteRookH) || s2 == whiteRookH:
		return cr & ^castleWhiteH
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
