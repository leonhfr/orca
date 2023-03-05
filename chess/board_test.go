package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	startingSquareMap = map[Square]Piece{
		A8: BlackRook, B8: BlackKnight, C8: BlackBishop, D8: BlackQueen,
		E8: BlackKing, F8: BlackBishop, G8: BlackKnight, H8: BlackRook,
		A7: BlackPawn, B7: BlackPawn, C7: BlackPawn, D7: BlackPawn,
		E7: BlackPawn, F7: BlackPawn, G7: BlackPawn, H7: BlackPawn,
		A2: WhitePawn, B2: WhitePawn, C2: WhitePawn, D2: WhitePawn,
		E2: WhitePawn, F2: WhitePawn, G2: WhitePawn, H2: WhitePawn,
		A1: WhiteRook, B1: WhiteKnight, C1: WhiteBishop, D1: WhiteQueen,
		E1: WhiteKing, F1: WhiteBishop, G1: WhiteKnight, H1: WhiteRook,
	}

	startingBoard = board{
		bbKing:   1152921504606846992,
		bbQueen:  576460752303423496,
		bbRook:   9295429630892703873,
		bbBishop: 2594073385365405732,
		bbKnight: 4755801206503243842,
		bbPawn:   71776119061282560,
		bbWhite:  65535,
		bbBlack:  18446462598732840960,
	}
)

func TestNewBoard(t *testing.T) {
	assert.Equal(t, startingBoard, newBoard(startingSquareMap))
}

func TestBoard_String(t *testing.T) {
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	assert.Equal(t, expected, startingBoard.String())
}
