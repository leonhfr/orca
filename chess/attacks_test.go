package chess

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoves(t *testing.T) {
	magics := [][64]Magic{rookMagics, bishopMagics}
	moves := [][]bitboard{bbMagicRookMoves, bbMagicBishopMoves}

	for pieceIndex, pt := range []PieceType{Rook, Bishop} {
		t.Run(pt.String(), func(t *testing.T) {
			for i := 0; i < 20; i++ {
				bb, sq := randomPosition()
				index := magics[pieceIndex][sq].index(bb)
				actual := moves[pieceIndex][index]
				expected := slowMoves(pt, sq, bb)

				assert.Equal(t, expected, actual)
			}
		})
	}
}

func BenchmarkMoves(b *testing.B) {
	magics := [][64]Magic{rookMagics, bishopMagics}
	moves := [][]bitboard{bbMagicRookMoves, bbMagicBishopMoves}
	bb, sq := randomPosition()

	for pieceIndex, pt := range []PieceType{Rook, Bishop} {
		b.Run(pt.String()+"-magic", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				index := magics[pieceIndex][sq].index(bb)
				_ = moves[pieceIndex][index]
			}
		})
		b.Run(pt.String()+"-slow", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				slowMoves(pt, sq, bb)
			}
		})
	}
}

func randomPosition() (bitboard, Square) {
	//nolint:gosec
	bb, sq := bitboard(rand.Uint64()), Square(rand.Intn(64))
	return bb, sq
}
