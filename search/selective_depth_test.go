package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestLateMoveReduction(t *testing.T) {
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
		{"check", args{5, true, 4, chess.Move(chess.Quiet)}, 0},
		{"depth", args{5, false, 3, chess.Move(chess.Quiet)}, 0},
		{"move", args{5, false, 4, chess.Move(chess.Capture)}, 0},
		{"reduction", args{5, false, 4, chess.Move(chess.Quiet)}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lateMoveReduction(tt.args.validMoves, tt.args.inCheck, tt.args.depth, tt.args.move)
			assert.Equal(t, tt.want, got)
		})
	}
}
