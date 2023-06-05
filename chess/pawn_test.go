package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type pawnCallbackArgs struct {
	p  Piece
	pp PawnProperty
}

var pawnPieceMapTests = []struct {
	name      string
	fen       string
	allPieces map[Square]pawnCallbackArgs
}{
	{
		name: "starting position",
		fen:  startFEN,
		allPieces: map[Square]pawnCallbackArgs{
			A7: {BlackPawn, HalfIsolani}, B7: {BlackPawn, NoProperty},
			C7: {BlackPawn, NoProperty}, D7: {BlackPawn, NoProperty},
			E7: {BlackPawn, NoProperty}, F7: {BlackPawn, NoProperty},
			G7: {BlackPawn, NoProperty}, H7: {BlackPawn, HalfIsolani},
			A2: {WhitePawn, HalfIsolani}, B2: {WhitePawn, NoProperty},
			C2: {WhitePawn, NoProperty}, D2: {WhitePawn, NoProperty},
			E2: {WhitePawn, NoProperty}, F2: {WhitePawn, NoProperty},
			G2: {WhitePawn, NoProperty}, H2: {WhitePawn, HalfIsolani},
		},
	},
	{
		name: "properties",
		fen:  "4k3/p1p3p1/3p3p/1P5P/1PP1P1P1/8/8/4K3 w - - 0 1",
		allPieces: map[Square]pawnCallbackArgs{
			A7: {BlackPawn, Isolani}, C7: {BlackPawn, HalfIsolani},
			D6: {BlackPawn, HalfIsolani}, G7: {BlackPawn, HalfIsolani},
			H6: {BlackPawn, HalfIsolani},
			B5: {WhitePawn, HalfIsolani ^ Doubled}, B4: {WhitePawn, HalfIsolani ^ Doubled},
			C4: {WhitePawn, HalfIsolani}, E4: {WhitePawn, Isolani},
			G4: {WhitePawn, HalfIsolani}, H5: {WhitePawn, HalfIsolani},
		},
	},
}

func TestPawnMap(t *testing.T) {
	for _, tt := range pawnPieceMapTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var pieces int
			pos.PawnMap(func(p Piece, sq Square, properties PawnProperty) {
				pieces++
				assert.Equal(t, tt.allPieces[sq].p, p, fmt.Sprintf("%v:%v", sq.String(), p.String()))
				assert.Equal(t, tt.allPieces[sq].pp, properties, fmt.Sprintf("%v:%v", sq.String(), properties))
			})
			assert.Equal(t, len(tt.allPieces), pieces)
		})
	}
}

func BenchmarkPawnMap(b *testing.B) {
	for _, bb := range pawnPieceMapTests {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				pos.PawnMap(func(_ Piece, _ Square, _ PawnProperty) {
					_ = 1
				})
			}
		})
	}
}
