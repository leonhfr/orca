package uci

// State contains the UCI engine's state.
type State struct {
	name   string
	author string
}

// New creates a new State.
func New(name, author string) *State {
	return &State{
		name:   name,
		author: author,
	}
}
