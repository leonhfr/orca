package chess

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPosition(t *testing.T) {
	t.Parallel()
	type args struct {
		s string
		n Notation
	}
	tests := []struct {
		args args
		want error
	}{
		{args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", FEN{}}, nil},
		{args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP/RNBQKBNR w KQkq - 0 1", FEN{}}, errors.New("invalid fen rank field (PPPPPPP)")},
	}

	for _, tt := range tests {
		t.Run(tt.args.s, func(t *testing.T) {
			t.Parallel()
			pos, err := NewPosition(tt.args.s)
			if tt.want == nil {
				assert.Equal(t, pos.String(), tt.args.s)
			}
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestStartingPosition(t *testing.T) {
	t.Parallel()
	assert.Equal(t, startFEN, StartingPosition().String())
}

func TestPosition_MakeMove(t *testing.T) {
	t.Parallel()
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			t.Parallel()
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
		meta := pos.Metadata()
		hash := pos.Hash()
		pawnHash := pos.PawnHash()
		b.Run(bb.moveUCI, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = pos.MakeMove(bb.move)
				pos.UnmakeMove(bb.move, meta, hash, pawnHash)
			}
		})
	}
}

func TestPosition_UnmakeMove(t *testing.T) {
	t.Parallel()
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			t.Parallel()
			pos := unsafeFEN(tt.preFEN)
			pre := unsafeFEN(tt.preFEN)
			meta := pos.Metadata()
			hash := pos.Hash()
			pawnHash := pos.PawnHash()
			_ = pos.MakeMove(tt.move)
			pos.UnmakeMove(tt.move, meta, hash, pawnHash)
			assert.Equal(t, tt.preFEN, pos.String())
			assert.Equal(t, pre.hash, pos.hash)
		})
	}
}
