package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticExchange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		fen  string
		move Move
		want []PieceType
	}{
		{
			"obvious case w",
			"1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1",
			newMove(WhiteRook, BlackPawn, E1, E5, NoSquare, NoPiece),
			[]PieceType{Rook},
		},
		{
			"obvious case b",
			"1k2r3/1pp4p/p7/4P3/8/P5P1/1PP4P/2K2R2 b - - 0 1",
			newMove(BlackRook, WhitePawn, E8, E5, NoSquare, NoPiece),
			[]PieceType{Rook},
		},
		{
			"xrays w",
			"1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1",
			newMove(WhiteKnight, BlackPawn, D3, E5, NoSquare, NoPiece),
			[]PieceType{Knight, Knight, Rook, Bishop, Queen, Queen},
		},
		{
			"xrays b",
			"1k2r2q/1ppn3p/p4b2/4P3/8/P2N2P1/1PP1R1BP/2K1Q3 b - - 0 1",
			newMove(BlackKnight, WhitePawn, D7, E5, NoSquare, NoPiece),
			[]PieceType{Knight, Knight, Bishop, Rook, Rook, Queen, Queen},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pos := unsafeFEN(tt.fen)
			var got []PieceType
			pos.StaticExchange(tt.move, func(pt PieceType) bool {
				got = append(got, pt)
				return false
			})
			assert.Equal(t, tt.want, got)
		})
	}
}
