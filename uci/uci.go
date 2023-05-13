package uci

import (
	"fmt"
	"io"
)

// State contains the UCI engine's state.
type State struct {
	name    string
	author  string
	debug   bool
	respond responder
}

// New creates a new State.
func New(name, author string, w io.Writer) *State {
	return &State{
		name:    name,
		author:  author,
		respond: func(r response) { fmt.Fprintln(w, r.String()) },
	}
}
