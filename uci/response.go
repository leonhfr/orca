package uci

import (
	"fmt"
	"strings"
	"time"

	"github.com/leonhfr/orca/chess"
)

// response is the interface implemented by objects that represent
// UCI responses from the Engine to the GUI.
type response interface {
	fmt.Stringer
}

// responseID represents a "id" command.
//
//	name <x>
//
// This must be sent after receiving the "uci" command to identify the engine,
// e.g. "id name Shredder X.Y\n".
//
//	author <x>
//
// This must be sent after receiving the "uci" command to identify the engine,
// e.g. "id author Stefan MK\n".
type responseID struct {
	name   string
	author string
}

func (r responseID) String() string {
	return fmt.Sprintf("id name %s\nid author %s", r.name, r.author)
}

// responseUCIOK represents a "uciok" command.
//
// Must be sent after the id and optional options to tell the GUI that the engine
// has sent all infos and is ready in uci mode.
type responseUCIOK struct{}

func (r responseUCIOK) String() string {
	return "uciok"
}

// responseReadyOK represents a "readyok" command.
//
// This must be sent when the engine has received an "isready" command and has
// processed all input and is ready to accept new commands now.
//
// It is usually sent after a command that can take some time to be able to wait for the engine,
// but it can be used anytime, even when the engine is searching,
// and must always be answered with "isready".
type responseReadyOK struct{}

func (r responseReadyOK) String() string {
	return "readyok"
}

// responseBestMove represents a "bestmove" command.
//
//	bestmove <move1> [ ponder <move2> ]
//
// The engine has stopped searching and found the move <move> best in this position.
// The engine can send the move it likes to ponder on. The engine must not start pondering automatically.
//
// This command must always be sent if the engine stops searching, also in pondering mode if there is a
// "stop" command, so for every "go" command a "bestmove" command is needed!
// Directly before that the engine should send a final "info" command with the final search information,
// the the GUI has the complete statistics about the last search.
type responseBestMove struct {
	move chess.Move
}

func (r responseBestMove) String() string {
	return fmt.Sprintf("bestmove %s", r.move.String())
}

// responseInfo represents an "info" command.
//
// The engine wants to send information to the GUI. This should be done whenever one of the info has changed.
// The engine can send only selected infos or multiple infos with one info command, e.g.
//
//	info currmove e2e4 currmovenumber 1
//	info depth 12 nodes 123456 nps 100000
//
// Also all infos belonging to the pv should be sent together e.g.
//
//	info depth 2 score cp 214 time 1242 nodes 2124 nps 34928 pv e2e4 e7e5 g1f3
//
// I suggest to start sending "currmove", "currmovenumber", "currline" and "refutation" only after one second
// to avoid too much traffic.
// Additional info:
//
//	depth <x>
//
// Search depth in plies.
//
//	seldepth <x>
//
// Selective search depth in plies,
// if the engine sends seldepth there must also be a "depth" present in the same string.
//
//	time <x>
//
// The time searched in ms, this should be sent together with the pv.
//
//	nodes <x>
//
// x nodes searched, the engine should send this info regularly.
//
//	pv <move1> ... <move i>
//
// The best line found.
//
//	multipv <num>
//
// This for the multi pv mode.
// For the best move/pv add "multipv 1" in the string when you send the pv.
// In k-best mode always send all k variants in k strings together.
//
//	score
//
// There are four possibilities:
//
//   - cp <x>: the score from the engine's point of view in centipawns
//
//   - mate <y>: mate in y moves, not plies. If the engine is getting mated use negative values for y
//
//   - lowerbound: the score is just a lower bound
//
//   - upperbound: the score is just an upper bound
//
//     currmove <move>
//
// Currently searching this move.
//
//	currmovenumber <x>
//
// Currently searching move number x, for the first move x should be 1 not 0.
//
//	hashfull <x>
//
// The hash is x permill full, the engine should send this info regularly.
//
//	nps <x>
//
// x nodes per second searched, the engine should send this info regularly.
//
//	tbhits <x>
//
// x positions where found in the endgame table bases.
//
//	sbhits <x>
//
// x positions where found in the shredder endgame databases.
//
//	cpuload <x>
//
// The cpu usage of the engine is x permill.
//
//	string <str>
//
// Any string str which will be displayed be the engine,
// if there is a string command the rest of the line will be interpreted as <str>.
//
//	refutation <move1> <move2> ... <move i>
//
// Move <move1> is refuted by the line <move2> ... <move i>, i can be any number >= 1.
//
// Example: after move d1h5 is searched, the engine can send
// "info refutation d1h5 g6h5"
// if g6h5 is the best answer after d1h5 or if g6h5 refutes the move d1h5.
// If there is no refutation for d1h5 found, the engine should just send
// "info refutation d1h5".
// The engine should only send this if the option "UCI_ShowRefutations" is set to true.
//
//	currline <cpunr> <move1> ... <move i>
//
// This is the current line the engine is calculating. <cpunr> is the number of the cpu if
// the engine is running on more than one cpu. <cpunr> = 1,2,3...
// If the engine is just using one cpu, <cpunr> can be omitted.
// If <cpunr> is greater than 1, always send all k lines in k strings together.
// The engine should only send this if the option "UCI_ShowCurrLine" is set to true.
type responseOutput struct {
	Output
	time time.Duration
}

func (o responseOutput) String() string {
	var res []string

	if o.Depth > 0 {
		res = append(res, "depth", fmt.Sprint(o.Depth))
	}
	if o.Nodes > 0 {
		res = append(res, "nodes", fmt.Sprint(o.Nodes))
	}
	if o.Mate != 0 {
		res = append(res, "score mate", fmt.Sprint(o.Mate))
	} else {
		res = append(res, "score cp", fmt.Sprint(o.Score))
	}
	if len(o.PV) > 0 {
		res = append(res, "pv")
		for _, move := range o.PV {
			res = append(res, move.String())
		}
	}
	if o.time > 0 {
		res = append(res, "time", fmt.Sprint(o.time.Milliseconds()))
	}

	return fmt.Sprintf("info %s", strings.Join(res, " "))
}

// Option represents an "option" command.
//
// This command tells the GUI which parameters can be changed in the engine.
// This should be sent once at engine startup after the "uci" and the "id" commands
// if any parameter can be changed in the engine.
// The GUI should parse this and build a dialog for the user to change the settings.
//
// Note that not every option needs to appear in this dialog as some options like
// "Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
//
// Note that the GUI need not send the setoption command when starting the engine for every option if
// it doesn't want to change the default value.
// For all allowed combinations see the examples below,
// as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
//
//	name <id>
//
// The option has the name id.
//
// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
// Usually those options should not be displayed in the normal engine options window of the GUI but
// get a special treatment. "Pondering" for example should be set automatically when pondering is
// enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
// automatically by the GUI. All those certain options have the prefix "UCI_" except for the
// first 6 options below. If the GUI gets an unknown Option with the prefix "UCI_", it should just
// ignore it and not display it in the engine's options dialog.
//
//	type <t>
//
// The option has type t.
// There are 5 different types of options the engine can send
//
//   - check: a checkbox that can either be true or false
//
//   - spin: a spin wheel that can be an integer in a certain range
//
//   - combo: a combo box that can have different predefined strings as a value
//
//   - button: a button that can be pressed to send a command to the engine
//
//   - string: a text field that has a string as a value, an empty string has the value "<empty>"
//
//     default <x>
//
// The default value of this parameter is x.
//
//	min <x>
//
// The minimum value of this parameter is x.
//
//	max <x>
//
// The maximum value of this parameter is x.
//
//	var <x>
//
// A predefined value of this parameter is x.
//
// Examples:
//
// Here are 5 strings for each of the 5 possible types of options
//   - "option name Nullmove type check default true\n"
//   - "option name Selectivity type spin default 2 min 0 max 4\n"
//   - "option name Style type combo default Normal var Solid var Normal var Risky\n"
//   - "option name NalimovPath type string default c:\\n"
//   - "option name Clear Hash type button\n"
func (o Option) String() string {
	switch o.Type {
	case OptionInteger:
		var min, max string
		if len(o.Min) > 0 {
			min = fmt.Sprintf(" min %s", o.Min)
		}
		if len(o.Max) > 0 {
			max = fmt.Sprintf(" max %s", o.Max)
		}
		return fmt.Sprintf(
			"option name %s type spin default %s%s%s",
			o.Name, o.Default, min, max,
		)
	case OptionBoolean:
		return fmt.Sprintf(
			"option name %s type check default %s",
			o.Name, o.Default,
		)
	default:
		return ""
	}
}
