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

	expected := concatenate([]response{
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

func concatenate(responses []response) string {
	s := make([]string, len(responses))
	for i, r := range responses {
		s[i] = r.String() + "\n"
	}
	return strings.Join(s, "")
}
