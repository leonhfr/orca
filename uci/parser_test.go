package uci

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

	tests := []struct {
		name string
		args string
		want command
	}{
		{name: "no input", args: "", want: nil},
		{name: "uci", args: "uci", want: commandUCI{}},
		{name: "debug on", args: "debug on", want: commandDebug{on: true}},
		{name: "debug off", args: "debug off", want: commandDebug{on: false}},
		{name: "isready", args: "isready", want: commandIsReady{}},
		{name: "setoption", args: "setoption name NAME value VALUE", want: commandSetOption{name: "NAME", value: "VALUE"}},
		{name: "ucinewgame", args: "ucinewgame", want: commandUCINewGame{}},
		{name: "position", args: "position startpos", want: commandPosition{startPos: true}},
		{name: "position", args: "position fen " + fen, want: commandPosition{fen: fen}},
		{name: "position", args: "position fen " + fen + " moves b1a3 b1c3", want: commandPosition{fen: fen, moves: []string{"b1a3", "b1c3"}}},
		{name: "position", args: "position startpos moves b1a3 b1c3", want: commandPosition{startPos: true, moves: []string{"b1a3", "b1c3"}}},
		{
			name: "go movetime",
			args: "go movetime 500",
			want: commandGo{
				moveTime: 500 * time.Millisecond,
			},
		},
		{
			name: "go time control",
			args: "go wtime 1000 btime 2000 winc 3000 binc 4000 movestogo 5",
			want: commandGo{
				whiteTime:      1 * time.Second,
				blackTime:      2 * time.Second,
				whiteIncrement: 3 * time.Second,
				blackIncrement: 4 * time.Second,
				movesToGo:      5,
			},
		},
		{
			name: "go inifinte searchmoves",
			args: "go infinite searchmoves b1a3 b1c3",
			want: commandGo{
				searchMoves: []string{"b1a3", "b1c3"},
				infinite:    true,
			},
		},
		{
			name: "go depth nodes",
			args: "go depth 8 nodes 1024",
			want: commandGo{
				depth: 8,
				nodes: 1024,
			},
		},
		{name: "stop", args: "stop", want: commandStop{}},
		{name: "quit", args: "quit", want: commandQuit{}},
		{name: "unknown", args: "foo bar", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(strings.Fields(tt.args))
			assert.Equal(t, tt.want, got)
		})
	}
}
