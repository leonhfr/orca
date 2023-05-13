package uci

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/orca/chess"
)

func TestRun(t *testing.T) {
	name, author := "NAME", "AUTHOR"
	e := new(mockEngine)
	w := &strings.Builder{}
	s := NewState(name, author, w)

	r := strings.NewReader("uci\nfake command\nquit")

	Run(context.Background(), e, r, s)
	e.AssertExpectations(t)
	assert.Equal(t, "id name NAME\nid author AUTHOR\nuciok\n", w.String())
}

// mockEngine is a mock that implements the engine interface.
type mockEngine struct {
	mock.Mock
}

// Search implements the engine interface.
func (m *mockEngine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan *Output {
	args := m.Called(ctx, pos, maxDepth)
	return args.Get(0).(chan *Output)
}
