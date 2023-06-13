package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type pieceCallbackArgs struct {
	piece    Piece
	mobility int
}

func TestCountPieces(t *testing.T) {
	pos := unsafeFEN(startFEN)
	knights, bishops, rooks, queens := pos.CountPieces()
	assert.Equal(t, 4, knights)
	assert.Equal(t, 4, bishops)
	assert.Equal(t, 4, rooks)
	assert.Equal(t, 2, queens)
}

var pieceMapTests = []struct {
	name string
	fen  string
	pm   map[Square]pieceCallbackArgs
}{
	{
		name: "starting position",
		fen:  startFEN,
		pm: map[Square]pieceCallbackArgs{
			A8: {BlackRook, 0}, B8: {BlackKnight, 2},
			C8: {BlackBishop, 0}, D8: {BlackQueen, 0},
			F8: {BlackBishop, 0}, G8: {BlackKnight, 2},
			H8: {BlackRook, 0},
			A1: {WhiteRook, 0}, B1: {WhiteKnight, 2},
			C1: {WhiteBishop, 0}, D1: {WhiteQueen, 0},
			F1: {WhiteBishop, 0}, G1: {WhiteKnight, 2},
			H1: {WhiteRook, 0},
		},
	},
	{
		name: "partial mirror",
		fen:  "r1bq1rk1/pppp1ppp/2nb1n2/1B2p3/4P3/P1NP1N2/1PP2PPP/R1BQK2R w KQ - 0 1",
		pm: map[Square]pieceCallbackArgs{
			A8: {BlackRook, 1}, C6: {BlackKnight, 5},
			C8: {BlackBishop, 0}, D8: {BlackQueen, 2},
			D6: {BlackBishop, 4}, F6: {BlackKnight, 5},
			F8: {BlackRook, 1},
			A1: {WhiteRook, 2}, C3: {WhiteKnight, 5},
			C1: {WhiteBishop, 5}, D1: {WhiteQueen, 2},
			B5: {WhiteBishop, 4}, F3: {WhiteKnight, 6},
			H1: {WhiteRook, 2},
		},
	},
}

func TestPieceMap(t *testing.T) {
	for _, tt := range pieceMapTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var pieces int
			pos.PieceMap(func(p Piece, sq Square, mobility int, trapped bool) {
				pieces++
				assert.Equal(t, tt.pm[sq].piece, p, fmt.Sprintf("%v:%v", sq.String(), p.String()))
				assert.Equal(t, tt.pm[sq].mobility, mobility, fmt.Sprintf("%v:%d", sq.String(), mobility))
			})
			assert.Equal(t, len(tt.pm), pieces)
		})
	}
}

func BenchmarkPieceMap(b *testing.B) {
	for _, bb := range pieceMapTests {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				pos.PieceMap(func(p Piece, sq Square, mobility int, trapped bool) {
					_ = 1
				})
			}
		})
	}
}
