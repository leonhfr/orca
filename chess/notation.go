package chess

// Notation is the interface implemented by objects that can
// encode and decode positions.
type Notation interface {
	// Encode encodes a position into a string. Does not validate the position.
	Encode(pos *Position) string
	// Decode decodes a string into a position. Does not validate the position.
	Decode(s string) (*Position, error)
}
