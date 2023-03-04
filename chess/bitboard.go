package chess

import "fmt"

// bitboard is a board representation encoded in an unsigned 64-bit integer.
// The 64 squares board has A1 as the least significant bit and H8 as the most.
type bitboard uint64

// String returns a 64 character string of 1s and 0s
// with the most significant bit.
//
// Returns an UCI-compatible representation.
func (b bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}
