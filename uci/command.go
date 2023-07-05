package uci

import (
	"context"
	"fmt"
	"time"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/search"
)

// command is the interface implemented by objects that represent
// UCI commands from the GUI to the *search.Engine.
type command interface {
	run(ctx context.Context, e *search.Engine, c *Controller)
}

// commandUCI represents a "uci" command.
type commandUCI struct{}

// run implements the command interface.
func (commandUCI) run(_ context.Context, _ *search.Engine, c *Controller) {
	c.respond(responseID{
		name:   c.name,
		author: c.author,
	})

	for _, option := range availableOptions {
		c.respond(option.uci())
	}

	c.respond(responseUCIOK{})
}

// commandDebug represents a "debug" command.
type commandDebug struct {
	on bool
}

// run implements the command interface.
func (cmd commandDebug) run(_ context.Context, _ *search.Engine, s *Controller) {
	s.debug = cmd.on
	s.logDebug("debug set to ", cmd.on)
}

// commandIsReady represents an "isready" command.
type commandIsReady struct{}

// run implements the command interface.
func (commandIsReady) run(_ context.Context, e *search.Engine, c *Controller) {
	go func() {
		err := e.Init()
		if err != nil {
			c.logError(err)
		}
		c.respond(responseReadyOK{})
	}()
}

// commandSetOption represents a "setoption" command.
type commandSetOption struct {
	name  string
	value string
}

// run implements the command interface.
func (cmd commandSetOption) run(_ context.Context, e *search.Engine, c *Controller) {
	for _, option := range availableOptions {
		if option.String() == cmd.name {
			fn, err := option.optionFunc(cmd.value)
			if err != nil {
				c.logError(err)
				return
			}
			fn(e)
			return
		}
	}

	c.logError(errOptionName)
}

// commandUCINewGame represents a "ucinewgame" command.
type commandUCINewGame struct{}

// run implements the command interface.
func (commandUCINewGame) run(_ context.Context, e *search.Engine, c *Controller) {
	go func() {
		c.position = chess.StartingPosition()
		err := e.Init()
		if err != nil {
			c.logError(err)
		}
	}()
}

// commandPosition represents a "position" command.
type commandPosition struct {
	fen      string
	moves    []string
	startPos bool
}

// run implements the command interface.
func (cmd commandPosition) run(_ context.Context, _ *search.Engine, c *Controller) {
	if cmd.startPos {
		c.position = chess.StartingPosition()
	} else if len(cmd.fen) > 0 {
		pos, err := c.notation.Decode(cmd.fen)
		if err != nil {
			c.logError(err)
			return
		}
		c.position = pos
	}

	for _, move := range cmd.moves {
		m, err := c.moveNotation.Decode(c.position, move)
		if err != nil {
			c.logError(err)
			return
		}

		if ok := c.position.MakeMove(m); !ok {
			c.logError(fmt.Errorf("failed to play move %s", move))
			return
		}
	}

	c.logDebug("position set to FEN ", c.position.String())
}

// commandGo represents a "go" command.
//
//nolint:govet
type commandGo struct {
	whiteTime      time.Duration // White has <x> ms left on the clock.
	blackTime      time.Duration // Black has <x> ms left on the clock.
	whiteIncrement time.Duration // White increment per move in ms if <x> > 0.
	blackIncrement time.Duration // Black increment per move in ms if <x> > 0.
	movesToGo      int           // Number of moves until the next time control.
	searchMoves    []string      // Restrict search to those moves only.
	depth          int           // Search <x> plies only.
	nodes          int           // Search <x> nodes only.
	moveTime       time.Duration // Search exactly <x> ms.
	infinite       bool          // Search until the stop command. Do not exit before.
}

// run implements the command interface.
func (cmd commandGo) run(ctx context.Context, e *search.Engine, c *Controller) {
	c.mu.Lock()
	start := time.Now()
	ctx, cancel := searchContext(ctx, c.stop, cmd.moveTime)

	outputs := e.Search(ctx, c.position, cmd.depth, cmd.nodes)

	go func() {
		defer c.mu.Unlock()
		defer cancel()

		var output search.Output
		for output = range outputs {
			c.respond(responseOutput{
				Output: output,
				time:   time.Since(start),
			})
		}
		if len(output.PV) > 0 {
			c.respond(responseBestMove{output.PV[0]})
		}
	}()
}

// searchContext creates a new context that is cancelled when
// a struct is emitted on the stop channel.
func searchContext(ctx context.Context, stop <-chan struct{}, limit time.Duration) (context.Context, context.CancelFunc) {
	if limit == 0 {
		limit = time.Hour
	}
	ctx, cancel := context.WithTimeout(ctx, limit)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			cancel()
		}
	}()

	return ctx, cancel
}

// commandStop represents a "stop" command.
type commandStop struct{}

// run implements the command interface.
func (commandStop) run(_ context.Context, _ *search.Engine, c *Controller) {
	select {
	case c.stop <- struct{}{}:
	default:
	}
}

// commandQuit represents a "quit" command.
type commandQuit struct{}

// run implements the command interface.
func (commandQuit) run(ctx context.Context, e *search.Engine, c *Controller) {
	commandStop{}.run(ctx, e, c)

	e.Close()

	// prevents future searches and ensures all search routines have been shut down
	c.mu.Lock()
}
