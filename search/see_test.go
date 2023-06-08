package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestSEE(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		move chess.Move
		want int32
	}{
		{
			"obvious case w",
			"1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1",
			newMove(chess.E1, chess.E5, chess.WhiteRook, chess.BlackPawn, chess.NoPiece, chess.Quiet),
			100,
		},
		{
			"obvious case b",
			"1k2r3/1pp4p/p7/4P3/8/P5P1/1PP4P/2K2R2 b - - 0 1",
			newMove(chess.E8, chess.E5, chess.BlackRook, chess.WhitePawn, chess.NoPiece, chess.Quiet),
			100,
		},
		{
			"xrays w",
			"1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1",
			newMove(chess.D3, chess.E5, chess.WhiteKnight, chess.BlackPawn, chess.NoPiece, chess.Quiet),
			-225,
		},
		{
			"xrays b",
			"1k2r2q/1ppn3p/p4b2/4P3/8/P2N2P1/1PP1R1BP/2K1Q3 b - - 0 1",
			newMove(chess.D7, chess.E5, chess.BlackKnight, chess.WhitePawn, chess.NoPiece, chess.Quiet),
			100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			got := see(pos, tt.move)
			assert.Equal(t, tt.want, got)
		})
	}
}
