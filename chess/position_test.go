package chess

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPosition(t *testing.T) {
	tests := []struct {
		args string
		want error
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", nil},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP/RNBQKBNR w KQkq - 0 1", errors.New("invalid fen rank field (PPPPPPP)")},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			pos, err := NewPosition(tt.args)
			if tt.want == nil {
				assert.Equal(t, pos.String(), tt.args)
			}
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestStartingPosition(t *testing.T) {
	assert.Equal(t, startFEN, StartingPosition().String())
}

func TestPosition_MakeMove(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			pos := unsafeFEN(tt.preFEN)
			exp := unsafeFEN(tt.postFEN)
			pos.MakeMove(tt.move)
			assert.Equal(t, tt.postFEN, pos.String())
			assert.Equal(t, exp.hash, pos.hash)
		})
	}
}

func BenchmarkPosition_MakeMove(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		hash := pos.Hash()
		b.Run(bb.moveUCI, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				meta, _ := pos.MakeMove(bb.move)
				pos.UnmakeMove(bb.move, meta, hash)
			}
		})
	}
}

func TestPosition_UnmakeMove(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			pos := unsafeFEN(tt.preFEN)
			pre := unsafeFEN(tt.preFEN)
			hash := pos.Hash()
			meta, _ := pos.MakeMove(tt.move)
			pos.UnmakeMove(tt.move, meta, hash)
			assert.Equal(t, tt.preFEN, pos.String())
			assert.Equal(t, pre.hash, pos.hash)
		})
	}
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
	name         string
	fen          string
	uniquePieces map[Square]Piece
	allPieces    map[Square]Piece
}{
	{
		name: "starting position",
		fen:  startFEN,
		uniquePieces: map[Square]Piece{
			E1: WhiteKing, E8: BlackKing,
			D1: WhiteQueen, D8: BlackQueen,
		},
		allPieces: startingSquareMap,
	},
	{
		name: "partial mirror",
		fen:  "r1bq1rk1/pppp1ppp/2nb1n2/1B2p3/4P3/P1NP1N2/1PP2PPP/R1BQK2R w KQ - 0 1",
		uniquePieces: map[Square]Piece{
			E1: WhiteKing, G8: BlackKing,
			D1: WhiteQueen, D8: BlackQueen,
			H1: WhiteRook, F8: BlackRook,
			A3: WhitePawn, A7: BlackPawn,
			D3: WhitePawn, D7: BlackPawn,
			B5: WhiteBishop, D6: BlackBishop,
		},
		allPieces: map[Square]Piece{
			A8: BlackRook, C6: BlackKnight, C8: BlackBishop, D8: BlackQueen,
			G8: BlackKing, D6: BlackBishop, F6: BlackKnight, F8: BlackRook,
			A7: BlackPawn, B7: BlackPawn, C7: BlackPawn, D7: BlackPawn,
			E5: BlackPawn, F7: BlackPawn, G7: BlackPawn, H7: BlackPawn,
			A3: WhitePawn, B2: WhitePawn, C2: WhitePawn, D3: WhitePawn,
			E4: WhitePawn, F2: WhitePawn, G2: WhitePawn, H2: WhitePawn,
			A1: WhiteRook, C3: WhiteKnight, C1: WhiteBishop, D1: WhiteQueen,
			E1: WhiteKing, B5: WhiteBishop, F3: WhiteKnight, H1: WhiteRook,
		},
	},
}

func TestPieceMap(t *testing.T) {
	for _, tt := range pieceMapTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var pieces int
			pos.PieceMap(func(p Piece, sq Square) {
				pieces++
				assert.Equal(t, tt.allPieces[sq], p, fmt.Sprintf("%v:%v", sq.String(), p.String()))
			})
			assert.Equal(t, len(tt.allPieces), pieces)
		})
	}
}

func BenchmarkPieceMap(b *testing.B) {
	for _, bb := range pieceMapTests {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				pos.PieceMap(func(p Piece, sq Square) {
					_ = 1
				})
			}
		})
	}
}

func TestUniquePieceMap(t *testing.T) {
	for _, tt := range pieceMapTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			var pieces int
			pos.UniquePieceMap(func(p Piece, sq Square) {
				pieces++
				assert.Equal(t, tt.uniquePieces[sq], p, fmt.Sprintf("%v:%v", sq.String(), p.String()))
			})
			assert.Equal(t, len(tt.uniquePieces), pieces)
		})
	}
}

func BenchmarkUniquePieceMap(b *testing.B) {
	for _, bb := range pieceMapTests {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				pos.UniquePieceMap(func(p Piece, sq Square) {
					_ = 1
				})
			}
		})
	}
}
