package uci

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandUCI(t *testing.T) {
	name, author := "NAME", "AUTHOR"
	w := &strings.Builder{}
	state := New(name, author, w)

	expected := concatenateResponses([]response{
		responseID{name, author},
		responseUCIOK{},
	})

	commandUCI{}.run(context.Background(), state)
	assert.Equal(t, expected, w.String())
}

func TestCommandDebug(t *testing.T) {
	state := New("", "", io.Discard)

	for _, tt := range []bool{true, false} {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			commandDebug{on: tt}.run(context.Background(), state)
			assert.Equal(t, tt, state.debug)
		})
	}
}

func TestCommandPosition(t *testing.T) {
	tests := []struct {
		c    commandPosition
		r    []string
		want string
	}{
		{
			commandPosition{startPos: true},
			[]string{},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			commandPosition{fen: "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23"},
			[]string{},
			"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
		},
		{
			commandPosition{fen: "bad fen"},
			[]string{"info string invalid fen (bad fen), must have 6 fields"},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			commandPosition{
				fen:   "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
				moves: []string{"a4b3"},
			},
			[]string{},
			"2r3k1/1q1nbppp/r3p3/3pP3/2pP4/PpQ2N2/2RN1PPP/2R4K w - - 0 24",
		},
		{
			commandPosition{
				fen:   "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
				moves: []string{"bad move"},
			},
			[]string{"info string invalid move in UCI notation"},
			"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			w := &strings.Builder{}
			state := New("", "", w)
			tt.c.run(context.Background(), state)
			assert.Equal(t, tt.want, state.position.String())
			assert.Equal(t, concatenateStrings(tt.r), w.String())
		})
	}
}

func concatenateStrings(ss []string) string {
	res := make([]string, len(ss))
	for i, s := range ss {
		res[i] = s + "\n"
	}
	return strings.Join(res, "")
}

func concatenateResponses(responses []response) string {
	s := make([]string, len(responses))
	for i, r := range responses {
		s[i] = r.String() + "\n"
	}
	return strings.Join(s, "")
}
