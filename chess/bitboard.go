package chess

import (
	"fmt"
	"math/bits"
)

// bitboard is a board representation encoded in an unsigned 64-bit integer.
// The 64 squares board has A1 as the least significant bit and H8 as the most.
type bitboard uint64

// emptyBitboard is an empty bitboard.
const emptyBitboard bitboard = 0

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

// ones returns the number of one bits in the bitboard.
func (b bitboard) ones() int {
	return bits.OnesCount64(uint64(b))
}

// String returns a 64 character string of 1s and 0s
// with the most significant bit.
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

var (
	bbFiles         [64]bitboard // bbFiles contains the file bitboards indexed by Square.
	bbRanks         [64]bitboard // bbRanks contains the rank bitboards indexed by Square.
	bbDiagonals     [64]bitboard // bbDiagonals contains the diagonal bitboards indexed by Square.
	bbAntiDiagonals [64]bitboard // bbAntiDiagonals contains the anti diagonal bitboards indexed by Square.
)

// initializes bbFiles.
func initBBFiles() {
	files := [8]bitboard{bbFileA, bbFileB, bbFileC, bbFileD, bbFileE, bbFileF, bbFileG, bbFileH}
	for sq := A1; sq <= H8; sq++ {
		bbFiles[sq] = files[sq.File()]
	}
}

// initBBRanks bbRanks.
func initBBRanks() {
	ranks := [8]bitboard{bbRank1, bbRank2, bbRank3, bbRank4, bbRank5, bbRank6, bbRank7, bbRank8}
	for sq := A1; sq <= H8; sq++ {
		bbRanks[sq] = ranks[sq.Rank()]
	}
}

// initializes bbDiagonals.
func initBBDiagonals() {
	for sq := A1; sq <= H8; sq++ {
		var bb bitboard
		for upLeft := sq; upLeft.Rank() != Rank8 && upLeft.File() != FileA; upLeft += 7 {
			bb |= upLeft.bitboard()
		}
		for downRight := sq; downRight.Rank() != Rank1 && downRight.File() != FileH; downRight -= 7 {
			bb |= downRight.bitboard()
		}
		bbDiagonals[sq] = bb
	}
}

// initializes bbAntiDiagonals.
func initBBAntiDiagonals() {
	for sq := A1; sq <= H8; sq++ {
		var bb bitboard
		for upRight := sq; upRight.Rank() != Rank8 && upRight.File() != FileH; upRight += 9 {
			bb |= upRight.bitboard()
		}
		for downLeft := sq; downLeft.Rank() != Rank1 && downLeft.File() != FileA; downLeft -= 9 {
			bb |= downLeft.bitboard()
		}
		bbAntiDiagonals[sq] = bb
	}
}

const (
	bbWhiteKingCastle  = 1<<F1 | 1<<H1 // bbWhiteKingCastle is the rook swap bitboard of a white king side castle.
	bbWhiteQueenCastle = 1<<A1 | 1<<D1 // bbWhiteQueenCastle is the rook swap bitboard of a white queen side castle.
	bbBlackKingCastle  = 1<<F8 | 1<<H8 // bbBlackKingCastle is the rook swap bitboard of a black king side castle.
	bbBlackQueenCastle = 1<<A8 | 1<<D8 // bbBlackQueenCastle is the rook swap bitboard of a black queen side castle.
)
