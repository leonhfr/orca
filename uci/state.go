package uci

import (
	"fmt"
	"io"
	"sync"

	"github.com/leonhfr/orca/chess"
)

// State contains the UCI engine's state.
type State struct {
	name     string
	author   string
	debug    bool
	position *chess.Position
	writer   io.Writer
	mu       sync.Mutex
	stop     chan struct{}
}

// NewState creates a new State.
func NewState(name, author string, writer io.Writer) *State {
	return &State{
		name:     name,
		author:   author,
		position: chess.StartingPosition(),
		writer:   writer,
		mu:       sync.Mutex{},
		stop:     make(chan struct{}),
	}
}

// logError logs an error to the output.
func (s *State) logError(err error) {
	fmt.Fprintln(s.writer, "info string", err.Error())
}

// logDebug logs debug info to the output.
func (s *State) logDebug(v ...any) {
	if s.debug {
		fmt.Fprintln(s.writer, "info string", fmt.Sprint(v...))
	}
}

// respond processes responses.
func (s *State) respond(r response) {
	fmt.Fprintln(s.writer, r.String())
}
