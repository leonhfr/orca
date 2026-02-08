package chess

import "errors"

// Square is one of the 64 squares on a chess board.
type Square uint8

//nolint:revive
const (
	A1, B1, C1, D1, E1, F1, G1, H1 Square = 8*iota + 0, 8*iota + 1, 8*iota + 2,
		8*iota + 3, 8*iota + 4, 8*iota + 5, 8*iota + 6, 8*iota + 7
	A2, B2, C2, D2, E2, F2, G2, H2
	A3, B3, C3, D3, E3, F3, G3, H3
	A4, B4, C4, D4, E4, F4, G4, H4
	A5, B5, C5, D5, E5, F5, G5, H5
	A6, B6, C6, D6, E6, F6, G6, H6
	A7, B7, C7, D7, E7, F7, G7, H7
	A8, B8, C8, D8, E8, F8, G8, H8
	NoSquare Square = 64
)

// newSquare create a new square from a file and rank.
func newSquare(f File, r Rank) Square {
	return Square(f) + Square(8*r)
}

// uciSquare parses a square from UCI notation.
func uciSquare(uci string) (Square, error) {
	b := []rune(uci)
	f := File(b[0] - 'a')
	r := Rank(b[1] - '1')
	if f > FileH || r > Rank8 {
		return NoSquare, errors.New("invalid uci square")
	}
	return newSquare(f, r), nil
}

// File returns the square's file.
func (sq Square) File() File {
	return File(sq & 7)
}

// Rank returns the square's rank.
func (sq Square) Rank() Rank {
	return Rank(sq & 56 >> 3)
}

// bitboard returns the square's bitboard.
func (sq Square) bitboard() bitboard {
	return 1 << sq
}

// sameColor checks whether the two squares are of the same color.
func (sq Square) sameColor(other Square) bool {
	return (9*uint16(sq^other))&8 == 0
}

// String implements the Stringer interface.
//
// Returns a UCI-compatible representation.
func (sq Square) String() string {
	return sq.File().String() + sq.Rank().String()
}

const (
	fileChars = "abcdefgh"
	rankChars = "12345678"
)

// A File is the file of a square.
type File uint8

//nolint:revive
const (
	FileA File = iota // FileA is the file A.
	FileB             // FileB is the file B.
	FileC             // FileC is the file C.
	FileD             // FileD is the file D.
	FileE             // FileE is the file E.
	FileF             // FileF is the file F.
	FileG             // FileG is the file G.
	FileH             // FileH is the file H.
)

// String implements the Stringer interface.
//
// Returns a UCI-compatible representation.
func (f File) String() string {
	return fileChars[f : f+1]
}

// Rank is the rank of a square.
type Rank uint8

//nolint:revive
const (
	Rank1 Rank = iota // Rank1 is the rank 1.
	Rank2             // Rank2 is the rank 2.
	Rank3             // Rank3 is the rank 3.
	Rank4             // Rank4 is the rank 4.
	Rank5             // Rank5 is the rank 5.
	Rank6             // Rank6 is the rank 6.
	Rank7             // Rank7 is the rank 7.
	Rank8             // Rank8 is the rank 8.
)

// String implements the Stringer interface.
//
// Returns a UCI-compatible representation.
func (r Rank) String() string {
	return rankChars[r : r+1]
}

// direction represent a step in a direction on the chess board.
type direction int

const (
	north     direction = 8
	northEast direction = 9
	east      direction = 1
	southEast direction = -7
	south     direction = -8
	southWest direction = -9
	west      direction = -1
	northWest direction = 7

	doubleNorth direction = 16
	doubleSouth direction = -16
)
