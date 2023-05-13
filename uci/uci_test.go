package uci

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/orca/chess"
)

// mockEngine is a mock that implements the Engine interface
type mockEngine struct {
	mock.Mock
}

func (m *mockEngine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan Output {
	args := m.Called(ctx, pos, maxDepth)
	return args.Get(0).(chan Output)
}
