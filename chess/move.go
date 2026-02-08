package chess

import "errors"

var (
	errInvalidMove     = errors.New("invalid move in UCI notation")
	errMissingPosition = errors.New("missing position")
)

// MoveTag represents a notable consequence of a move.
type MoveTag uint64

const (
	// Quiet indicates that the move is a priori quiet.
	Quiet MoveTag = 1 << (iota + 24)
	// Capture indicates that the move captures a piece.
	Capture
	// Check indicates that the move puts the enemy king in check.
	// Only informative, moves that do not have the tag may put the
	// enemy king in check.
	Check
	// EnPassant indicates that the move captures a piece via en passant.
	EnPassant
	// Promotion indicates that the move is a promotion.
	Promotion
	// ASideCastle indicates that the move is a A side castle (queen side in classic chess).
	ASideCastle
	// HSideCastle indicates that the move is a H side castle (king side in classic chess).
	HSideCastle
)

// Move represents a move from one square to another.
//
//	32 bits
//	uint32 with move score
//
//	32 bits
//	xxxxxxxx pppp tttt ffff TTTTTT FFFFFF
//
//	xxxxxxxx   move tags
//	pppp       promo piece
//	tttt       to piece
//	ffff       from piece
//	TTTTTT     to square
//	FFFFFF     from square
type Move uint64

// NoMove represents the absence of a move.
const NoMove Move = 0

// newMove creates a new move.
//
// Expects the classic chess castling convention (king jumps two squares to castle).
//
// Example: e1g1 for H side white castle.
func newMove(p1, p2 Piece, s1, s2, enPassant Square, promo Piece) Move {
	var tags MoveTag
	switch pt := p1.Type(); {
	case pt == King && ((s1 == E1 && s2 == C1) || (s1 == E8 && s2 == C8)):
		tags ^= ASideCastle
	case pt == King && ((s1 == E1 && s2 == G1) || (s1 == E8 && s2 == G8)):
		tags ^= HSideCastle
	case pt == Pawn && s2 == enPassant:
		tags ^= EnPassant
		p2 = Pawn.color(p1.Color().Other())
	case promo != NoPiece:
		tags ^= Promotion
	}

	if p2 != NoPiece {
		tags |= Capture
	}

	if tags == 0 {
		tags ^= Quiet
	}

	return Move(s1) ^ Move(s2)<<6 ^
		Move(p1)<<12 ^ Move(p2)<<16 ^
		Move(promo)<<20 ^ Move(tags)
}

// newChess960Move creates a new move.
//
// Expects the Chess960 castling convention (king takes own rook).
//
// Example: e1h1 for H side white castle in the classic position.
func newChess960Move(p1, p2 Piece, s1, s2, enPassant Square, promo Piece, files [2]File) Move {
	var tags MoveTag
	switch pt, rook := p1.Type(), Rook.color(p1.Color()); {
	case pt == King && p2 == rook && s2.File() == files[aSide]:
		tags ^= ASideCastle
		p2 = NoPiece
		s2 = newSquare(kingFinalFile[aSide], s1.Rank())
	case pt == King && p2 == rook && s2.File() == files[hSide]:
		tags ^= HSideCastle
		p2 = NoPiece
		s2 = newSquare(kingFinalFile[hSide], s1.Rank())
	case pt == Pawn && s2 == enPassant:
		tags ^= EnPassant
		p2 = Pawn.color(p1.Color().Other())
	case promo != NoPiece:
		tags ^= Promotion
	}

	if p2 != NoPiece {
		tags |= Capture
	}

	if tags == 0 {
		tags ^= Quiet
	}

	return Move(s1) ^ Move(s2)<<6 ^
		Move(p1)<<12 ^ Move(p2)<<16 ^
		Move(promo)<<20 ^ Move(tags)
}

// newCastleMove creates a new castle move.
func newCastleMove(p1 Piece, s1, s2 Square, s side, check bool) Move {
	var tags MoveTag

	if s == aSide {
		tags ^= ASideCastle
	} else {
		tags ^= HSideCastle
	}

	if check {
		tags ^= Check
	}

	return Move(s1) ^ Move(s2)<<6 ^
		Move(p1)<<12 ^ Move(NoPiece)<<16 ^
		Move(NoPiece)<<20 ^ Move(tags)
}

// newPawnMove creates a new pawn move.
func newPawnMove(p1, p2 Piece, s1, s2, enPassant Square, promo Piece, check bool) Move {
	var tags MoveTag

	if s2 == enPassant {
		tags ^= EnPassant
		p2 = Pawn.color(p1.Color().Other())
	} else if promo != NoPiece {
		tags ^= Promotion
	}

	if p2 != NoPiece {
		tags ^= Capture
	}

	if check {
		tags ^= Check
	}

	if tags == 0 {
		tags ^= Quiet
	}

	return Move(s1) ^ Move(s2)<<6 ^
		Move(p1)<<12 ^ Move(p2)<<16 ^
		Move(promo)<<20 ^ Move(tags)
}

// newPieceMove creates a new piece move.
func newPieceMove(p1, p2 Piece, s1, s2 Square, check bool) Move {
	var tags MoveTag

	if p2 != NoPiece {
		tags ^= Capture
	}

	if check {
		tags ^= Check
	}

	if tags == 0 {
		tags ^= Quiet
	}

	return Move(s1) ^ Move(s2)<<6 ^
		Move(p1)<<12 ^ Move(p2)<<16 ^
		Move(NoPiece)<<20 ^ Move(tags)
}

// NewMove creates a new move from a UCI string.
//
// Shorthand for:
//
//	UCI{}.Decode(pos, move)
func NewMove(pos *Position, move string) (Move, error) {
	return UCI{}.Decode(pos, move)
}

// S1 returns the origin square of the move.
func (m Move) S1() Square {
	return Square(m & 63)
}

// S2 returns the destination square of the move.
func (m Move) S2() Square {
	return Square((m >> 6) & 63)
}

// P1 returns the piece in the origin square.
func (m Move) P1() Piece {
	return Piece((m >> 12) & 15)
}

// P2 returns the piece in the destination square.
func (m Move) P2() Piece {
	return Piece((m >> 16) & 15)
}

// Promo returns the promotion piece of the move.
func (m Move) Promo() Piece {
	return Piece((m >> 20) & 15)
}

// HasTag checks whether the move has the given MoveTag.
func (m Move) HasTag(tag MoveTag) bool {
	return tag&MoveTag(m) > 0
}

// Score returns the score of the move.
func (m Move) Score() uint32 {
	return uint32(m >> 32)
}

// WithScore returns a new move with the score set.
func (m Move) WithScore(score uint32) Move {
	return (Move(score) << 32) ^ m
}

// String implements the Stringer interface.
//
// Returns a UCI-compatible representation.
//
// Shorthand for:
//
//	UCI{}.Encode(nil, move)
func (m Move) String() string {
	return UCI{}.Encode(nil, m)
}
