package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pawnPieceMapTests = []struct {
	name         string
	fen          string
	uniquePieces map[Square]Piece
	allPieces    map[Square]Piece
}{
	{
		name:         "starting position",
		fen:          startFEN,
		uniquePieces: map[Square]Piece{},
		allPieces: map[Square]Piece{
			A7: BlackPawn, B7: BlackPawn, C7: BlackPawn, D7: BlackPawn,
			E7: BlackPawn, F7: BlackPawn, G7: BlackPawn, H7: BlackPawn,
			A2: WhitePawn, B2: WhitePawn, C2: WhitePawn, D2: WhitePawn,
			E2: WhitePawn, F2: WhitePawn, G2: WhitePawn, H2: WhitePawn,
		},
	},
	{
		name: "partial mirror",
		fen:  "r1bq1rk1/pppp1ppp/2nb1n2/1B2p3/4P3/P1NP1N2/1PP2PPP/R1BQK2R w KQ - 0 1",
		uniquePieces: map[Square]Piece{
			A3: WhitePawn, A7: BlackPawn,
			D3: WhitePawn, D7: BlackPawn,
		},
		allPieces: map[Square]Piece{
			A7: BlackPawn, B7: BlackPawn, C7: BlackPawn, D7: BlackPawn,
			E5: BlackPawn, F7: BlackPawn, G7: BlackPawn, H7: BlackPawn,
			A3: WhitePawn, B2: WhitePawn, C2: WhitePawn, D3: WhitePawn,
			E4: WhitePawn, F2: WhitePawn, G2: WhitePawn, H2: WhitePawn,
		},
	},
}

func TestPawnMap(t *testing.T) {
	for _, tt := range pawnPieceMapTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var pieces int
			pos.PawnMap(func(p Piece, sq Square) {
				pieces++
				assert.Equal(t, tt.allPieces[sq], p, fmt.Sprintf("%v:%v", sq.String(), p.String()))
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
				pos.PawnMap(func(p Piece, sq Square) {
					_ = 1
				})
			}
		})
	}
}
