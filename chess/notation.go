package chess

// Notation is the interface implemented by objects that can
// encode and decode positions.
type Notation interface {
	// Encode encodes a position into a string. Does not validate the position.
	Encode(pos *Position) string
	// Decode decodes a string into a position. Does not validate the position.
	Decode(s string) (*Position, error)
}

// MoveNotation is the interface implemented by objects that can
// encode and decode moves.
type MoveNotation interface {
	// Encode encodes a move into a string. Does not validate the move.
	Encode(pos *Position, m Move) string
	// Decode decodes a string into a move. Does not validate the move.
	Decode(pos *Position, s string) (Move, error)
}
