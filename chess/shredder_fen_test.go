package chess

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// compile time check that ShredderFEN implements Notation.
var _ Notation = ShredderFEN{}

func TestShredderFEN(t *testing.T) {
	tests := []struct {
		args string
		want error
	}{
		{"bqnb1rkr/pp3ppp/3ppn2/2p5/5P2/P2P4/NPP1P1PP/BQ1BNRKR w HFhf - 2 9", nil},
		{"2nnrbkr/p1qppppp/8/1ppb4/6PP/3PP3/PPP2P2/BQNNRBKR w HEhe - 1 9", nil},
		{"b1q1rrkb/pppppppp/3nn3/8/P7/1PPP4/4PPPP/BQNNRKRB w GE - 1 9", nil},
		{"qbbnnrkr/2pp2pp/p7/1p2pp2/8/P3PP2/1PPP1KPP/QBBNNR1R w hf - 0 9", nil},
		{"1nbbnrkr/p1p1ppp1/3p4/1p3P1p/3Pq2P/8/PPP1P1P1/QNBBNRKR w HFhf - 0 9", nil},
		{"bqnb1rkr/pp3ppp/3ppn2/2p5/5P2/P2P4/NPP1P1PP/BQ1BNRKR w ABc - 2 9", errors.New("invalid fen castling files (ABc)")},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			pos, err := ShredderFEN{}.Decode(tt.args)
			if tt.want == nil {
				assert.Equal(t, tt.args, ShredderFEN{}.Encode(pos))
			}
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestShredderFENCastling(t *testing.T) {
	tests := []struct {
		args castling
		want string
	}{
		{castling{[2]File{FileA, FileH}, noCastle}, "-"},
		{castling{[2]File{FileA, FileH}, castleBlackA | castleBlackH | castleWhiteA | castleWhiteH}, "HAha"},
		{castling{[2]File{FileA, FileH}, castleWhiteA | castleWhiteH}, "HA"},
		{castling{[2]File{FileF, FileH}, castleBlackA | castleBlackH | castleWhiteA | castleWhiteH}, "HFhf"},
		{castling{[2]File{FileE, FileG}, castleWhiteA | castleWhiteH}, "GE"},
		{castling{[2]File{FileD, FileH}, castleBlackA | castleBlackH | castleWhiteA | castleWhiteH}, "HDhd"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := shredderFenCastling(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestShredderFENCastlingFiles(t *testing.T) {
	type want struct {
		cf  [2]File
		err error
	}

	tests := []struct {
		args string
		want
	}{
		{"-", want{[2]File{FileA, FileH}, nil}},
		{"HAha", want{[2]File{FileA, FileH}, nil}},
		{"HA", want{[2]File{FileA, FileH}, nil}},
		{"HFhf", want{[2]File{FileF, FileH}, nil}},
		{"GE", want{[2]File{FileE, FileG}, nil}},
		{"HDhd", want{[2]File{FileD, FileH}, nil}},
		{"ABc", want{[2]File{FileA, FileH}, errors.New("invalid fen castling files (ABc)")}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			cf, err := shredderFenCastlingFiles(tt.args)
			assert.Equal(t, tt.want.cf, cf)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestShredderFENCastlingRights(t *testing.T) {
	type args struct {
		field string
		files [2]File
	}
	type want struct {
		cr  castlingRights
		err error
	}

	tests := []struct {
		args
		want
	}{
		{args{"-", [2]File{FileA, FileH}}, want{noCastle, nil}},
		{args{"HAha", [2]File{FileA, FileH}}, want{castleBlackA | castleBlackH | castleWhiteA | castleWhiteH, nil}},
		{args{"HA", [2]File{FileA, FileH}}, want{castleWhiteA | castleWhiteH, nil}},
		{args{"HFhf", [2]File{FileF, FileH}}, want{castleBlackA | castleBlackH | castleWhiteA | castleWhiteH, nil}},
		{args{"GE", [2]File{FileE, FileG}}, want{castleWhiteA | castleWhiteH, nil}},
		{args{"HDhd", [2]File{FileD, FileH}}, want{castleBlackA | castleBlackH | castleWhiteA | castleWhiteH, nil}},
		{args{"AB", [2]File{FileA, FileH}}, want{0, errors.New("invalid fen castling rights (AB)")}},
	}

	for _, tt := range tests {
		t.Run(tt.args.field, func(t *testing.T) {
			cr, err := shredderFenCastlingRights(tt.args.field, tt.args.files)
			assert.Equal(t, tt.want.cr, cr)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
