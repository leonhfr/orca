package uci

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/orca/chess"
)

// mockEngine is a mock that implements the engine interface.
type mockEngine struct {
	mock.Mock
}

// Search implements the engine interface.
func (m *mockEngine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan *Output {
	args := m.Called(ctx, pos, maxDepth)
	return args.Get(0).(chan *Output)
}
