package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestKillerListGet(t *testing.T) {
	t.Parallel()
	A, B := chess.Move(chess.A2), chess.Move(chess.B2)
	C, D := chess.Move(chess.C2), chess.Move(chess.D2)
	E, F := chess.Move(chess.E2), chess.Move(chess.F2)

	tests := []struct {
		name  string
		list  killerList
		depth uint8
		want  [2]chess.Move
	}{
		{
			"empty list",
			killerList{},
			1,
			[2]chess.Move{},
		},
		{
			"index greater than list depth",
			killerList{},
			4,
			[2]chess.Move{},
		},
		{
			"index greater than list length",
			killerList{entries: [maxSearchDepth][2]chess.Move{
				{A, B}, // 2
				{C, D}, // 1
			}},
			3,
			[2]chess.Move{},
		},
		{
			"index outside bounds",
			killerList{entries: [maxSearchDepth][2]chess.Move{
				{A, B}, // 2
				{C, D}, // 1
			}},
			maxSearchDepth,
			[2]chess.Move{},
		},
		{
			"index within bounds",
			killerList{entries: [maxSearchDepth][2]chess.Move{
				{A, B}, // 3
				{C, D}, // 2
				{E, F}, // 1
			}},
			2,
			[2]chess.Move{E, F},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.list.get(tt.depth)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKillerListSet(t *testing.T) {
	t.Parallel()
	A, B := chess.Move(chess.A2), chess.Move(chess.B2)
	C, D := chess.Move(chess.C2), chess.Move(chess.D2)

	type args struct {
		move  chess.Move
		depth uint8
	}

	tests := []struct {
		name  string
		depth int
		args  []args
		want  [maxSearchDepth][2]chess.Move
	}{
		{
			"depth 1",
			2,
			[]args{{A, 1}, {B, 1}, {C, 0}, {A, 1}, {D, 0}},
			[maxSearchDepth][2]chess.Move{{D, C}, {B, A}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			kl := newKillerList()
			for _, args := range tt.args {
				kl.set(args.move, args.depth)
			}
			assert.Equal(t, tt.want, kl.entries)
		})
	}
}
