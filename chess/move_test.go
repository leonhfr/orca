package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMove_Valid(t *testing.T) {
	t.Parallel()
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			t.Parallel()
			move, _ := NewMove(unsafeFEN(tt.preFEN), tt.moveUCI)
			assert.Equal(t, tt.move, move)
			for _, tag := range tt.tags {
				assert.True(t, move.HasTag(tag))
			}
		})
	}
}

func TestNewMove_Invalid(t *testing.T) {
	t.Parallel()
	type (
		args struct {
			fen string
			uci string
		}
		want struct {
			move Move
			err  error
		}
	)

	tests := []struct {
		args args
		want want
	}{
		{
			args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2"},
			want{0, errInvalidMove},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args.uci, func(t *testing.T) {
			t.Parallel()
			move, err := NewMove(unsafeFEN(tt.args.fen), tt.args.uci)
			assert.Equal(t, tt.want.move, move)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
