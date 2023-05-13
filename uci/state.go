package uci

import (
	"fmt"
	"io"

	"github.com/leonhfr/orca/chess"
)

// State contains the UCI engine's state.
type State struct {
	name     string
	author   string
	debug    bool
	position *chess.Position
	writer   io.Writer
}

// New creates a new State.
func New(name, author string, writer io.Writer) *State {
	return &State{
		name:     name,
		author:   author,
		position: chess.StartingPosition(),
		writer:   writer,
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
