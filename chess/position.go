package chess

import (
	"fmt"
	"math/bits"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board          board
	hash           Hash
	pawnHash       Hash
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
func (pos *Position) MakeMove(m Move) (Metadata, bool) {
	if (m.HasTag(KingSideCastle) || m.HasTag(QueenSideCastle)) && !pos.isCastleLegal(m) {
		return NoMetadata, false
	}

	metadata := newMetadata(pos.turn, pos.castlingRights,
		pos.halfMoveClock, pos.fullMoves, pos.enPassant)
	cr := pos.castlingRights

	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbColors[White]&pos.board.bbPieces[Pawn], pos.board.bbColors[Black]&pos.board.bbPieces[Pawn])
	}

	pos.board.makeMove(m)
	if pos.isSquareAttacked(pos.board.kingSquare(pos.turn)) {
		pos.board.makeMove(m)
		return NoMetadata, false
	}

	pos.turn = pos.turn.other()
	pos.castlingRights = moveCastlingRights(pos.castlingRights, m)
	pos.enPassant = moveEnPassant(m)
	if pos.enPassant != NoSquare {
		pos.hash ^= enPassantHash(pos.enPassant, pos.turn,
			pos.board.bbColors[White]&pos.board.bbPieces[Pawn], pos.board.bbColors[Black]&pos.board.bbPieces[Pawn])
	}

	partialHash, partialPawnHash := xorHashPartialMove(m, cr, pos.castlingRights)
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

	return metadata, true
}

// UnmakeMove unmakes a move and restores the previous position.
func (pos *Position) UnmakeMove(m Move, meta Metadata, hash, pawnHash Hash) {
	pos.board.makeMove(m)
	pos.turn = meta.turn()
	pos.castlingRights = meta.castleRights()
	pos.enPassant = meta.enPassant()
	pos.halfMoveClock = meta.halfMoveClock()
	pos.fullMoves = meta.fullMoves()
	pos.hash = hash
	pos.pawnHash = pawnHash
}

// CountPieces returns the count of knights, bishops, rooks, and queens.
func (pos *Position) CountPieces() (int, int, int, int) {
	return pos.board.bbPieces[Knight].ones(), pos.board.bbPieces[Bishop].ones(), pos.board.bbPieces[Rook].ones(), pos.board.bbPieces[Queen].ones()
}

// PieceMap executes the callback for each piece on the board, passing the piece
// and its square as arguments.
//
// Does not take pawns into account.
//
// Intended to be used in evaluation functions.
func (pos *Position) PieceMap(cb func(p Piece, sq Square)) {
	for p, bb := range [8]bitboard{
		pos.board.bbColors[Black] & pos.board.bbPieces[Knight],
		pos.board.bbColors[White] & pos.board.bbPieces[Knight],
		pos.board.bbColors[Black] & pos.board.bbPieces[Bishop],
		pos.board.bbColors[White] & pos.board.bbPieces[Bishop],
		pos.board.bbColors[Black] & pos.board.bbPieces[Rook],
		pos.board.bbColors[White] & pos.board.bbPieces[Rook],
		pos.board.bbColors[Black] & pos.board.bbPieces[Queen],
		pos.board.bbColors[White] & pos.board.bbPieces[Queen],
	} {
		piece := Piece(p + 2)
		for ; bb > 0; bb = bb.resetLSB() {
			cb(piece, bb.scanForward())
		}
	}

	cb(BlackKing, (pos.board.bbColors[Black] & pos.board.bbPieces[King]).scanForward())
	cb(WhiteKing, (pos.board.bbColors[White] & pos.board.bbPieces[King]).scanForward())
}

// UniquePieceMap executes the callback for each piece on the board that do not
// have an opponent mirrored piece, passing the piece and its square as arguments.
//
// Does not take pawns into account.
//
// Intended to be used in evaluation functions.
func (pos *Position) UniquePieceMap(cb func(p Piece, sq Square)) {
	bbBlack, bbWhite := pos.board.bbColors[Black], pos.board.bbColors[White]
	bbKing := pos.board.bbPieces[King]
	bbQueen := pos.board.bbPieces[Queen]
	bbRook := pos.board.bbPieces[Rook]
	bbBishop := pos.board.bbPieces[Bishop]
	bbKnight := pos.board.bbPieces[Knight]

	for p, bb := range [8]bitboard{
		^bitboard(bits.ReverseBytes64(uint64(bbWhite&bbKnight))) & bbBlack & bbKnight,
		^bitboard(bits.ReverseBytes64(uint64(bbBlack&bbKnight))) & bbWhite & bbKnight,
		^bitboard(bits.ReverseBytes64(uint64(bbWhite&bbBishop))) & bbBlack & bbBishop,
		^bitboard(bits.ReverseBytes64(uint64(bbBlack&bbBishop))) & bbWhite & bbBishop,
		^bitboard(bits.ReverseBytes64(uint64(bbWhite&bbRook))) & bbBlack & bbRook,
		^bitboard(bits.ReverseBytes64(uint64(bbBlack&bbRook))) & bbWhite & bbRook,
		bbBlack & bbQueen,
		bbWhite & bbQueen,
	} {
		piece := Piece(p + 2)
		for ; bb > 0; bb = bb.resetLSB() {
			cb(piece, bb.scanForward())
		}
	}

	cb(BlackKing, (bbBlack & bbKing).scanForward())
	cb(WhiteKing, (bbWhite & bbKing).scanForward())
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
	if pt := p1.Type(); pt != King && pt != Rook && !m.HasTag(Capture) {
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
