package chess

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFENBoard(t *testing.T) {
	type want struct {
		b   board
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP", want{board{}, errors.New("invalid fen board (rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP)")}},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", want{startingBoard, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			b, err := fenBoard(tt.args)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.b, b)
		})
	}
}

func TestFENFileField(t *testing.T) {
	type want struct {
		fm  map[File]Piece
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"", want{nil, errors.New("invalid fen rank field ()")}},
		{"-", want{nil, errors.New("invalid fen rank field (-)")}},
		{"rnbqkbnr", want{map[File]Piece{
			FileA: BlackRook,
			FileB: BlackKnight,
			FileC: BlackBishop,
			FileD: BlackQueen,
			FileE: BlackKing,
			FileF: BlackBishop,
			FileG: BlackKnight,
			FileH: BlackRook,
		}, nil}},
		{"RNBQKBNR", want{map[File]Piece{
			FileA: WhiteRook,
			FileB: WhiteKnight,
			FileC: WhiteBishop,
			FileD: WhiteQueen,
			FileE: WhiteKing,
			FileF: WhiteBishop,
			FileG: WhiteKnight,
			FileH: WhiteRook,
		}, nil}},
		{"2p5", want{map[File]Piece{FileC: BlackPawn}, nil}},
		{"4P3", want{map[File]Piece{FileE: WhitePawn}, nil}},
		{"8", want{map[File]Piece{}, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			fm, err := fenFileField(tt.args)
			assert.Equal(t, tt.want.fm, fm)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestFENTurn(t *testing.T) {
	type want struct {
		turn Color
		err  error
	}

	tests := []struct {
		args string
		want
	}{
		{"-", want{White, errors.New("invalid fen turn (-)")}},
		{"w", want{White, nil}},
		{"b", want{Black, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			turn, err := fenTurn(tt.args)
			assert.Equal(t, tt.want.turn, turn)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestFENCastlingFiles(t *testing.T) {
	type want struct {
		cf  [2]File
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"-", want{[2]File{FileA, FileH}, nil}},
		{"KQkq", want{[2]File{FileA, FileH}, nil}},
		{"KQ", want{[2]File{FileA, FileH}, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			cf, err := fenCastlingFiles(tt.args)
			assert.Equal(t, tt.want.cf, cf)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestFENCastlingRights(t *testing.T) {
	type want struct {
		cr  castlingRights
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"-", want{0, nil}},
		{"KQkq", want{castleBlackA | castleBlackH | castleWhiteA | castleWhiteH, nil}},
		{"KQ", want{castleWhiteA | castleWhiteH, nil}},
		{"A", want{0, errors.New("invalid fen castling rights (A)")}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			cr, err := fenCastlingRights(tt.args)
			assert.Equal(t, tt.want.cr, cr)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestFENEnPassantSquare(t *testing.T) {
	type want struct {
		sq  Square
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"", want{NoSquare, errors.New("invalid fen en passant square ()")}},
		{"-", want{NoSquare, nil}},
		{"e3", want{E3, nil}},
		{"c6", want{C6, nil}},
		{"e4", want{NoSquare, errors.New("invalid fen en passant square (e4)")}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			sq, err := fenEnPassantSquare(tt.args)
			assert.Equal(t, tt.want.sq, sq)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestFENHalfMoveClock(t *testing.T) {
	type want struct {
		hmc uint8
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"-1", want{0, errors.New("invalid fen full moves count (-1)")}},
		{"0", want{0, nil}},
		{"1", want{1, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			hmc, err := fenHalfMoveClock(tt.args)
			assert.Equal(t, tt.want.hmc, hmc)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestFENFullMoves(t *testing.T) {
	type want struct {
		fm  uint8
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"0", want{0, errors.New("invalid fen full moves count (0)")}},
		{"1", want{1, nil}},
		{"2", want{2, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			fm, err := fenFullMoves(tt.args)
			assert.Equal(t, tt.want.fm, fm)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
