package uci

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

func TestResponseString(t *testing.T) {
	controller := NewController("", "", io.Discard)

	m1 := chess.Move(chess.B1) ^ chess.Move(chess.A3)<<6 ^ chess.Move(chess.NoPiece)<<20
	m2 := chess.Move(chess.E6) ^ chess.Move(chess.E7)<<6 ^ chess.Move(chess.NoPiece)<<20

	tests := []struct {
		name string
		args response
		want string
	}{
		{name: "id", args: responseID{name: "NAME", author: "AUTHOR"}, want: "id name NAME\nid author AUTHOR"},
		{name: "uciok", args: responseUCIOK{}, want: "uciok"},
		{name: "readyok", args: responseReadyOK{}, want: "readyok"},
		{name: "bestmove", args: responseBestMove{m1}, want: "bestmove b1a3"},
		{
			name: "info score positive",
			args: responseOutput{
				Output{
					Depth: 8,
					Nodes: 1024,
					Score: 3000,
					PV:    []chess.Move{m1, m2},
				},
				time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score cp 3000 pv b1a3 e6e7 time 5000",
		},
		{
			name: "info score negative",
			args: responseOutput{
				Output{
					Depth: 8,
					Nodes: 1024,
					Score: -3000,
					PV:    []chess.Move{m1, m2},
				},
				time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score cp -3000 pv b1a3 e6e7 time 5000",
		},
		{
			name: "info mate positive",
			args: responseOutput{
				Output{
					Depth: 8,
					Score: 3000,
					Nodes: 1024,
					Mate:  5,
					PV:    []chess.Move{m1, m2},
				},
				time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score mate 5 pv b1a3 e6e7 time 5000",
		},
		{
			name: "info mate negative",
			args: responseOutput{
				Output{
					Depth: 8,
					Nodes: 1024,
					Score: -3000,
					Mate:  -5,
					PV:    []chess.Move{m1, m2},
				},
				time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score mate -5 pv b1a3 e6e7 time 5000",
		},
		{
			name: "integer option",
			args: testOptions[OptionInteger],
			want: "option name INTEGER OPTION type spin default 32 min 2 max 1024",
		},
		{
			name: "boolean option",
			args: testOptions[OptionBoolean],
			want: "option name BOOLEAN OPTION type check default false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.format(controller))
		})
	}
}
