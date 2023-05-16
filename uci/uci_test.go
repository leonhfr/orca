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
	e.On("Close")
	e.On("Options").Return([]Option{})
	w := &strings.Builder{}
	s := NewState(name, author, w)

	r := strings.NewReader("uci\nfake command\nquit\n")

	Run(context.Background(), e, r, s)
	e.AssertExpectations(t)
	assert.Equal(t, "id name NAME\nid author AUTHOR\nuciok\n", w.String())
}

// mockEngine is a mock that implements the Engine interface.
type mockEngine struct {
	mock.Mock
}

// Init implements the Engine interface.
func (me *mockEngine) Init() error {
	args := me.Called()
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

// Close implements the Engine interface.
func (me *mockEngine) Close() {
	me.Called()
}

// Options implements the Engine interface.
func (me *mockEngine) Options() []Option {
	args := me.Called()
	return args.Get(0).([]Option)
}

// SetOption implements the Engine interface.
func (me *mockEngine) SetOption(name, value string) error {
	args := me.Called(name, value)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

// Search implements the Engine interface.
func (me *mockEngine) Search(ctx context.Context, pos *chess.Position, maxDepth int) <-chan Output {
	args := me.Called(ctx, pos, maxDepth)
	return args.Get(0).(chan Output)
}

// testOptions contains test options.
var testOptions = []Option{
	{
		Type:    OptionInteger,
		Name:    "INTEGER OPTION",
		Default: "32",
		Min:     "2",
		Max:     "1024",
	},
	{
		Type:    OptionBoolean,
		Name:    "BOOLEAN OPTION",
		Default: "false",
	},
}
