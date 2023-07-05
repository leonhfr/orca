// Package uci handles communications between Universal Chess Interface (UCI) commands
// and a chess engine.
//
// Chess engines have to implement the uci.Engine interface.
package uci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/search"
)

// Controller handles communications between UCI commands and a chess engine.
//
// Chess engines have to implement the uci.Engine interface.
//
//nolint:govet
type Controller struct {
	name         string
	author       string
	debug        bool
	position     *chess.Position
	notation     chess.Notation
	moveNotation chess.MoveNotation
	writer       io.Writer
	mu           sync.Mutex
	stop         chan struct{}
}

// NewController creates a new Controller.
func NewController(name, author string, writer io.Writer) *Controller {
	return &Controller{
		name:         name,
		author:       author,
		position:     chess.StartingPosition(),
		notation:     chess.FEN{},
		moveNotation: chess.UCI{},
		writer:       writer,
		mu:           sync.Mutex{},
		stop:         make(chan struct{}),
	}
}

// Option represents an option.
type Option func(*Controller)

// WithChess960 sets the chess mode to classic chess or Chess960.
func WithChess960(on bool) Option {
	if on {
		return func(c *Controller) {
			c.notation = chess.ShredderFEN{}
			c.moveNotation = chess.UCIChess960{}
		}
	}

	return func(c *Controller) {
		c.notation = chess.FEN{}
		c.moveNotation = chess.UCI{}
	}
}

// Run runs the controller.
//
// Run parses command from the reader, executes them with the provided
// search engine and writes the responses on the writer.
func (c *Controller) Run(ctx context.Context, e *search.Engine, r io.Reader) {
	// graceful shutdown when context canceled
	// sending EOF to the UCI scanner by closing the pipe
	pipeR, pipeW := io.Pipe()
	go func() { _, _ = io.Copy(pipeW, r) }()
	go func() { <-ctx.Done(); pipeW.Close() }()

	for scanner := bufio.NewScanner(pipeR); scanner.Scan(); {
		cmd := parse(strings.Fields(scanner.Text()))
		if cmd == nil {
			continue
		}
		cmd.run(ctx, e, c)
		if _, ok := cmd.(commandQuit); ok {
			break
		}
	}
}

// logError logs an error to the output.
func (c *Controller) logError(err error) {
	fmt.Fprintln(c.writer, "info string", err.Error())
}

// logDebug logs debug info to the output.
func (c *Controller) logDebug(v ...any) {
	if c.debug {
		fmt.Fprintln(c.writer, "info string", fmt.Sprint(v...))
	}
}

// respond processes responses.
func (c *Controller) respond(r response) {
	fmt.Fprintln(c.writer, r.format(c))
}
