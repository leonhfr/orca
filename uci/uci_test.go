package uci

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerRun(t *testing.T) {
	name, author := "NAME", "AUTHOR"
	e := new(mockEngine)
	e.On("Close")
	e.On("Options").Return([]Option{})
	w := &strings.Builder{}
	s := NewController(name, author, w)

	r := strings.NewReader("uci\nfake command\nquit\n")

	s.Run(context.Background(), e, r)
	e.AssertExpectations(t)
	assert.Equal(t, "id name NAME\nid author AUTHOR\nuciok\n", w.String())
}
