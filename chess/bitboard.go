package chess

import (
	"fmt"
	"math"
	"math/bits"
)

// bitboard is a board representation encoded in an unsigned 64-bit integer.
// The 64 squares board has A1 as the least significant bit and H8 as the most.
type bitboard uint64

const (
	// bbEmpty is an empty bitboard.
	bbEmpty bitboard = 0
	// bbFull is a full bitboard.
	bbFull bitboard = math.MaxUint64
)

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

// northOne shifts the bitboard toward the north.
func (b bitboard) northOne() bitboard {
	return b << 8
}

// southOne shifts the bitboard toward the south.
func (b bitboard) southOne() bitboard {
	return b >> 8
}

// eastOne shifts the bitboard toward the east.
func (b bitboard) eastOne() bitboard {
	return b << 1 & ^bbFileA
}

// westOne shifts the bitboard toward the west.
func (b bitboard) westOne() bitboard {
	return b >> 1 & ^bbFileH
}

// northEastOne shifts the bitboard toward the north east.
func (b bitboard) northEastOne() bitboard {
	return b << 9 & ^bbFileA
}

// northWestOne shifts the bitboard toward the north west.
func (b bitboard) northWestOne() bitboard {
	return b << 7 & ^bbFileH
}

// southEastOne shifts the bitboard toward the south east.
func (b bitboard) southEastOne() bitboard {
	return b >> 7 & ^bbFileA
}

// southWestOne shifts the bitboard toward the south west.
func (b bitboard) southWestOne() bitboard {
	return b >> 9 & ^bbFileH
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
	bbFiles         [64]bitboard     // bbFiles contains the file bitboards indexed by Square.
	bbRanks         [64]bitboard     // bbRanks contains the rank bitboards indexed by Square.
	bbDiagonals     [64]bitboard     // bbDiagonals contains the diagonal bitboards indexed by Square.
	bbAntiDiagonals [64]bitboard     // bbAntiDiagonals contains the anti diagonal bitboards indexed by Square.
	bbInBetweens    [64][64]bitboard // bbInBetween contains the in between bitboards indexed by from and to Squares.
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
		// up right
		for file, rank := int(sq.File()), int(int(sq.Rank())); file <= int(FileH) && rank <= int(Rank8); file, rank = file+1, rank+1 {
			bb |= newSquare(File(file), Rank(rank)).bitboard()
		}
		// down left
		for file, rank := int(sq.File()), int(sq.Rank()); file >= int(FileA) && rank >= int(Rank1); file, rank = file-1, rank-1 {
			bb |= newSquare(File(file), Rank(rank)).bitboard()
		}
		bbDiagonals[sq] = bb
	}
}

// initializes bbAntiDiagonals.
func initBBAntiDiagonals() {
	for sq := A1; sq <= H8; sq++ {
		var bb bitboard
		// up left
		for file, rank := int(sq.File()), int(sq.Rank()); file >= int(FileA) && rank <= int(Rank8); file, rank = file-1, rank+1 {
			bb |= newSquare(File(file), Rank(rank)).bitboard()
		}
		// down right
		for file, rank := int(sq.File()), int(sq.Rank()); file <= int(FileH) && rank >= int(Rank1); file, rank = file+1, rank-1 {
			bb |= newSquare(File(file), Rank(rank)).bitboard()
		}
		bbAntiDiagonals[sq] = bb
	}
}

// initializes bbInBetween.
func initBBInBetweens() {
	m1 := uint64(bbFull)
	var a2a7 uint64 = 0x0001010101010100
	var b2g7 uint64 = 0x0040201008040200
	var h1b7 uint64 = 0x0002040810204080

	for s1 := uint64(A1); s1 <= uint64(H8); s1++ {
		for s2 := uint64(A1); s2 <= uint64(H8); s2++ {
			btwn := m1<<s1 ^ m1<<s2
			file := s2&7 - s1&7
			rank := (s2 | 7 - s1) >> 3
			line := (file&7 - 1) & a2a7
			line += 2 * ((rank&7 - 1) >> 58)
			line += ((rank-file)&15 - 1) & b2g7
			line += ((rank+file)&15 - 1) & h1b7
			line *= btwn & -btwn
			bbInBetweens[s1][s2] = bitboard(line & btwn)
		}
	}
}

const (
	bbWhiteKingCastle  = 1<<F1 | 1<<H1 // bbWhiteKingCastle is the rook swap bitboard of a white king side castle.
	bbWhiteQueenCastle = 1<<A1 | 1<<D1 // bbWhiteQueenCastle is the rook swap bitboard of a white queen side castle.
	bbBlackKingCastle  = 1<<F8 | 1<<H8 // bbBlackKingCastle is the rook swap bitboard of a black king side castle.
	bbBlackQueenCastle = 1<<A8 | 1<<D8 // bbBlackQueenCastle is the rook swap bitboard of a black queen side castle.
)
