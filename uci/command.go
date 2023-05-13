package uci

import (
	"context"
	"fmt"
	"time"

	"github.com/leonhfr/orca/chess"
)

// command is the interface implemented by objects that represent
// UCI commands from the GUI to the Engine.
type command interface {
	run(ctx context.Context, e engine, s *State)
}

// commandUCI represents a "uci" command.
//
// Tell engine to use the uci (universal chess interface),
// this will be sent once as a first command after program boot
// to tell the engine to switch to uci mode.
//
// After receiving the uci command the engine must identify itself with the "id" command
// and send the "option" commands to tell the GUI which engine settings the engine supports if any.
//
// After that the engine should send "uciok" to acknowledge the uci mode.
// If no "uciok" is sent within a certain time period, the engine task will be killed by the GUI.
type commandUCI struct{}

// run implements the command interface.
func (commandUCI) run(_ context.Context, _ engine, s *State) {
	s.respond(responseID{
		name:   s.name,
		author: s.author,
	})

	s.respond(responseUCIOK{})
}

// commandDebug represents a "debug" command.
//
// Switch the debug mode of the engine on and off.
// In debug mode the engine should send additional infos to the GUI, e.g. with the "info string" command,
// to help debugging, e.g. the commands that the engine has received etc.
//
// This mode should be switched off by default and this command can be sent
// any time, also when the engine is thinking.
type commandDebug struct {
	on bool
}

// run implements the command interface.
func (c commandDebug) run(_ context.Context, _ engine, s *State) {
	s.debug = c.on
	s.logDebug("debug set to ", c.on)
}

// commandIsReady represents an "isready" command.
//
// This is used to synchronize the engine with the GUI. When the GUI has sent a command or
// multiple commands that can take some time to complete,
// this command can be used to wait for the engine to be ready again or
// to ping the engine to find out if it is still alive.
// E.g. this should be sent after setting the path to the tablebases as this can take some time.
//
// This command is also required once before the engine is asked to do any search
// to wait for the engine to finish initializing.
//
// This command must always be answered with "readyok" and can be sent also when the engine is calculating
// in which case the engine should also immediately answer with "readyok" without stopping the search.
type commandIsReady struct{}

// run implements the command interface.
func (commandIsReady) run(_ context.Context, _ engine, _ *State) {
}

// commandSetOption represents a "setoption" command.
//
// This is sent to the engine when the user wants to change the internal parameters
// of the engine. For the "button" type no value is needed.
//
// One string will be sent for each parameter and this will only be sent when the engine is waiting.
// The name and value of the option in <id> should not be case sensitive and can include spaces.
//
// The substrings "value" and "name" should be avoided in <id> and <x> to allow unambiguous parsing,
// for example do not use <name> = "draw value".
//
// Here are some strings for the example below:
//
//	setoption name Nullmove value true\n
//	setoption name Selectivity value 3\n
//	setoption name Style value Risky\n
//	setoption name Clear Hash\n
//	setoption name NalimovPath value c:\chess\tb\4;c:\chess\tb\5\n
type commandSetOption struct {
	name  string
	value string
}

// run implements the command interface.
func (c commandSetOption) run(_ context.Context, _ engine, _ *State) {
}

// commandUCINewGame represents a "ucinewgame" command.
//
// This is sent to the engine when the next search (started with "position" and "go") will be from
// a different game. This can be a new game the engine should play or a new game it should analyze but
// also the next position from a testsuite with positions only.
//
// If the GUI hasn't sent a "ucinewgame" before the first "position" command, the engine shouldn't
// expect any further ucinewgame commands as the GUI is probably not supporting the ucinewgame command.
// So the engine should not rely on this command even though all new GUIs should support it.
//
// As the engine's reaction to "ucinewgame" can take some time the GUI should always send "isready"
// after "ucinewgame" to wait for the engine to finish its operation.
type commandUCINewGame struct{}

// run implements the command interface.
func (commandUCINewGame) run(_ context.Context, _ engine, _ *State) {
}

// commandPosition represents a "position" command.
//
//	position [fen <fenstring> | startpos ] moves <move1> ... <move i>
//
// Set up the position described in fenstring on the internal board and
// play the moves on the internal chess board.
//
// If the game was played  from the start position the string "startpos" will be sent.
//
// Note: no "new" command is needed. However, if this position is from a different game than
// the last position sent to the engine, the GUI should have sent a "ucinewgame" in between.
type commandPosition struct {
	fen      string
	startPos bool
	moves    []string
}

// run implements the command interface.
func (c commandPosition) run(_ context.Context, _ engine, s *State) {
	if c.startPos {
		s.position = chess.StartingPosition()
	} else if len(c.fen) > 0 {
		pos, err := chess.NewPosition(c.fen)
		if err != nil {
			s.logError(err)
			return
		}
		s.position = pos
	}

	for _, move := range c.moves {
		m, err := chess.NewMove(s.position, move)
		if err != nil {
			s.logError(err)
			return
		}

		if _, ok := s.position.MakeMove(m); !ok {
			s.logError(fmt.Errorf("failed to play move %s", move))
			return
		}
	}

	s.logDebug("position set to FEN ", s.position.String())
}

// commandGo represents a "go" command.
//
// Start calculating on the current position set up with the "position" command.
//
// There are a number of commands that can follow this command, all will be sent in the same string.
// If one command is not sent its value should be interpreted as it would not influence the search.
//
//	searchmoves <move1> ... <move i>
//
// Restrict search to this moves only.
//
// Example: After "position startpos" and "go infinite searchmoves e2e4 d2d4"
// the engine should only search the two moves e2e4 and d2d4 in the initial position.
//
//	ponder
//
// Start searching in pondering mode.
//
// Do not exit the search in ponder mode, even if it's mate!
// This means that the last move sent in in the position string is the ponder move.
// The engine can do what it wants to do, but after a "ponderhit" command
// it should execute the suggested move to ponder on. This means that the ponder move sent by
// the GUI can be interpreted as a recommendation about which move to ponder. However, if the
// engine decides to ponder on a different move, it should not display any mainlines as they are
// likely to be misinterpreted by the GUI because the GUI expects the engine to ponder
// on the suggested move.
//
//	wtime <x>
//
// White has x msec left on the clock.
//
//	btime <x>
//
// Black has x msec left on the clock.
//
//	winc <x>
//
// White increment per move in milliseconds if x > 0.
//
//	binc <x>
//
// Black increment per move in milliseconds if x > 0.
//
//	movestogo <x>
//
// There are x moves to the next time control, this will only be sent if x > 0.
// If you don't get this and get the wtime and btime it's sudden death.
//
//	depth <x>
//
// Search x plies only.
//
//	nodes <x>
//
// Search x nodes only.
//
//	mate <x>
//
// Search for a mate in x moves.
//
//	movetime <x>
//
// Search exactly x milliseconds.
//
//	infinite
//
// Search until the "stop" command. Do not exit the search without being told so in this mode!
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
func (c commandGo) run(_ context.Context, _ engine, _ *State) {
}

// commandStop represents a "stop" command.
//
// Stop calculating as soon as possible,
// don't forget the "bestmove" and possibly the "ponder" token when finishing the search.
type commandStop struct{}

// run implements the command interface.
func (commandStop) run(_ context.Context, _ engine, _ *State) {
}

// commandQuit represents a "quit" command.
//
// Quit the program as soon as possible.
type commandQuit struct{}

// run implements the command interface.
func (commandQuit) run(_ context.Context, _ engine, _ *State) {
}
