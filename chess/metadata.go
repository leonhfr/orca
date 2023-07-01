package chess

// Metadata represents a position's metadata.
//
//	32 bits
//	__ square fullMove halfMove __ cccc _ T
//	square      en passant square
//	fullMove    full moves
//	halfMove    half move clock
//	cccc        castle rights
//	T           turn color
//	_           unused bit
type Metadata uint32

// NoMetadata represents the absence of metadata.
const NoMetadata Metadata = 0

// Metadata returns the position metadata.
func (pos *Position) Metadata() Metadata {
	return newMetadata(pos.turn, pos.castling.rights,
		pos.halfMoveClock, pos.fullMoves, pos.enPassant)
}

// newMetadata create a new metadata.
func newMetadata(c Color, cr castlingRights, halfMoveClock, fullMoves uint8, enPassant Square) Metadata {
	return Metadata(c) |
		Metadata(cr)<<2 |
		Metadata(halfMoveClock)<<8 |
		Metadata(fullMoves)<<16 |
		Metadata(enPassant)<<24
}

// turn returns the turn color.
func (m Metadata) turn() Color {
	return Color(m & 1)
}

// castleRights returns the castle rights.
func (m Metadata) castleRights() castlingRights {
	return castlingRights((m >> 2) & 15)
}

// halfMoveClock returns the half move clock.
func (m Metadata) halfMoveClock() uint8 {
	return uint8(m >> 8)
}

// fullMoves returns the full moves.
func (m Metadata) fullMoves() uint8 {
	return uint8(m >> 16)
}

// enPassant returns the en passant square.
func (m Metadata) enPassant() Square {
	return Square(m >> 24)
}
