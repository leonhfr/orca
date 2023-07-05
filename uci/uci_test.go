package uci

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/search"
)

func TestControllerRun(t *testing.T) {
	name, author := "NAME", "AUTHOR"
	e := search.NewEngine()
	w := &strings.Builder{}
	c := NewController(name, author, w)

	r := strings.NewReader("uci\nfake command\nquit\n")

	expected := concatenateResponses(c, []response{
		responseID{name, author},
		availableSearchOptions[0].response(),
		availableSearchOptions[1].response(),
		availableUCIOptions[0].response(),
		responseUCIOK{},
	})

	c.Run(context.Background(), e, r)
	assert.Equal(t, expected, w.String())
}
