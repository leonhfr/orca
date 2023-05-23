package uci

import (
	"context"

	"github.com/leonhfr/orca/chess"
)

// Engine is the interface that should be implemented by the search engine.
type Engine interface {
	// Init initializes the search engine.
	Init() error
	// Close shuts down the resources used by the search engine.
	Close()
	// Options lists the available options.
	Options() []Option
	// SeOption sets an option.
	SetOption(name, value string) error
	// Search runs a search on the given position until the given depth.
	// Cancelling the context stops the search.
	Search(ctx context.Context, pos *chess.Position, maxDepth, maxNodes int) <-chan Output
}

// Output holds a search output.
type Output struct {
	PV    []chess.Move // Principal variation, best line found.
	Depth int          // Search depth in plies.
	Nodes int          // Number of nodes searched.
	Score int          // Score from the engine's point of view in centipawns.
	Mate  int          // Number of moves before mate. Positive for the current player to mate, negative for the current player to be mated.
}

// OptionType represent an option type.
type OptionType uint8

const (
	OptionInteger OptionType = iota // OptionInteger represents an integer option.
	OptionBoolean                   // OptionBoolean represents a boolean option.
)

// Option represents an available option.
//
//nolint:govet
type Option struct {
	Type    OptionType
	Name    string
	Default string
	Min     string
	Max     string
}
