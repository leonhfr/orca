package chess

import (
	"errors"
	"math/rand"
)

//go:generate go run magics_gen.go

// Magic represents a magic bitboard.
type Magic struct {
	Mask   uint64
	Magic  uint64
	Shift  int
	Offset int
}

// index computes the index of the move bitboard.
func (m *Magic) index(blockers bitboard) int {
	hash := (uint64(blockers) & m.Mask) * m.Magic
	return m.Offset + int(hash>>m.Shift)
}

// FindMagics finds magics for all squares.
//
// This function is intended to be used during generation of magics.
func FindMagics(pt PieceType) [64]Magic {
	var magics [64]Magic
	var index int
	for sq := A1; sq <= H8; sq++ {
		indexBits := slowMasks(pt, sq).ones()
		magic, moves := findMagic(pt, sq, indexBits)
		magic.Offset = index
		magics[sq] = magic
		index += moves
	}
	return magics
}

// findMagic finds a magic for a square.
//
// This function is intended to be used during generation of magics.
func findMagic(pt PieceType, sq Square, indexBits int) (Magic, int) {
	mask := uint64(slowMasks(pt, sq))
	shift := 64 - indexBits
	for {
		magic := Magic{Mask: mask, Magic: randomMagic(), Shift: shift}
		moves, err := slowMoveTable(pt, sq, magic)
		if err == nil {
			return magic, len(moves)
		}
	}
}

// slowMoveTable computes a move table.
//
// This function is intended to be used during initialization of move tables.
func slowMoveTable(pt PieceType, sq Square, magic Magic) ([]bitboard, error) {
	table := make([]bitboard, 1<<(64-magic.Shift))
	for blockers := bbEmpty; ; {
		index := magic.index(blockers) - magic.Offset
		moves := slowMoves(pt, sq, blockers)

		if table[index] == bbEmpty {
			table[index] = moves
		} else if table[index] != moves {
			return nil, errors.New("hash collision")
		}

		blockers = (blockers - bitboard(magic.Mask)) & bitboard(magic.Mask)
		if blockers == bbEmpty {
			break
		}
	}
	return table, nil
}

// randomMagic returns a random uint64 with a low bit count.
func randomMagic() uint64 {
	//nolint:gosec
	m1, m2, m3 := rand.Uint64(), rand.Uint64(), rand.Uint64()
	return m1 & m2 & m3
}
