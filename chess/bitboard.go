package chess

import "fmt"

// bitboard is a board representation encoded in an unsigned 64-bit integer.
// The 64 squares board has A1 as the least significant bit and H8 as the most.
type bitboard uint64

// deBruijnMagicTable contains a lookup table of squares indexed by the result
// of the de Bruijn multiplication.
var deBruijnMagicTable = [64]Square{
	A1, H6, B1, A8, A7, D4, C1, E8,
	B8, B7, B6, F5, E4, A3, D1, F8,
	G7, C8, D5, E7, C7, C6, F3, E6,
	G5, A5, F4, H3, B3, D2, E1, G8,
	G6, H7, C4, D8, A6, E5, H2, F7,
	C5, D7, E3, D6, H4, G3, C2, F6,
	B4, H5, G2, B5, D3, G4, B2, A4,
	F2, C3, A2, E2, H1, G1, F1, H8,
}

// deBruijnMagic is the magic de Bruijn number associated with the lookup table.
const deBruijnMagic = 0x03f79d71b4cb0a89

// scanForward returns the square in the lowest significant bit position.
//
// bitboard can't be 0.
//
// Uses de Bruijn forward scanning:
// https://www.chessprogramming.org/BitScan#De_Bruijn_Multiplication
func (b bitboard) scanForward() Square {
	index := ((b ^ (b - 1)) * deBruijnMagic) >> 58
	return deBruijnMagicTable[index]
}

// resetLSB resets the lowest significant bit.
func (b bitboard) resetLSB() bitboard {
	return b & (b - 1)
}

// String returns a 64 character string of 1s and 0s
// with the most significant bit.
//
// Returns an UCI-compatible representation.
func (b bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}

const (
	bbRank1 bitboard = (1<<A1 | 1<<B1 | 1<<C1 | 1<<D1 | 1<<E1 | 1<<F1 | 1<<G1 | 1<<H1) << (8 * iota)
	bbRank2
	bbRank3
	bbRank4
	bbRank5
	bbRank6
	bbRank7
	bbRank8
)

const (
	bbFileA bitboard = (1<<A1 | 1<<A2 | 1<<A3 | 1<<A4 | 1<<A5 | 1<<A6 | 1<<A7 | 1<<A8) << iota
	bbFileB
	bbFileC
	bbFileD
	bbFileE
	bbFileF
	bbFileG
	bbFileH
)
