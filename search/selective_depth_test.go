package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestLateMoveReduction(t *testing.T) {
	t.Parallel()
	type args struct {
		validMoves int
		inCheck    bool
		depth      uint8
		move       chess.Move
	}

	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"validMoves", args{4, false, 4, chess.Move(chess.Quiet)}, 0},
		{"in check", args{5, true, 4, chess.Move(chess.Quiet)}, 0},
		{"depth", args{5, false, 3, chess.Move(chess.Quiet)}, 0},
		{"move", args{5, false, 4, chess.Move(chess.Capture)}, 0},
		{"reduction", args{5, false, 4, chess.Move(chess.Quiet)}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := lateMoveReduction(tt.args.validMoves, tt.args.inCheck, tt.args.depth, tt.args.move)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestShouldNullMovePrune(t *testing.T) {
	t.Parallel()
	type args struct {
		fen     string
		inCheck bool
		depth   uint8
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"in check", args{"8/8/8/8/5R2/2pk4/5K2/8 w - - 0 1", true, 3}, false},
		{"depth", args{"8/8/8/8/5R2/2pk4/5K2/8 w - - 0 1", false, 2}, false},
		{"pieces", args{"8/8/8/8/5R2/2pk4/5K2/8 b - - 0 1", false, 3}, false},
		{"pruning", args{"8/8/8/8/5R2/2pk4/5K2/8 w - - 0 1", false, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pos := unsafeFEN(tt.args.fen)
			got := shouldNullMovePrune(pos, tt.args.inCheck, tt.args.depth)
			assert.Equal(t, tt.want, got)
		})
	}
}
