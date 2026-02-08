package uci

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/search"
)

// response is the interface implemented by objects that represent
// UCI responses from the Engine to the GUI.
type response interface {
	format(c *Controller) string
}

// responseID represents a "id" command.
type responseID struct {
	name   string
	author string
}

func (r responseID) format(_ *Controller) string {
	return fmt.Sprintf("id name %s\nid author %s", r.name, r.author)
}

// responseUCIOK represents a "uciok" command.
type responseUCIOK struct{}

func (r responseUCIOK) format(_ *Controller) string {
	return "uciok"
}

// responseReadyOK represents a "readyok" command.
type responseReadyOK struct{}

func (r responseReadyOK) format(_ *Controller) string {
	return "readyok"
}

// responseBestMove represents a "bestmove" command.
type responseBestMove struct {
	move chess.Move
}

func (r responseBestMove) format(c *Controller) string {
	return "bestmove " + c.moveNotation.Encode(c.position, r.move)
}

// responseInfo represents an "info" command.
type responseOutput struct {
	search.Output
	time time.Duration
}

func (o responseOutput) format(c *Controller) string {
	var res []string

	if o.Depth > 0 {
		res = append(res, "depth", strconv.Itoa(o.Depth))
	}
	if o.Nodes > 0 {
		res = append(res, "nodes", strconv.Itoa(o.Nodes))
	}
	if o.Mate != 0 {
		res = append(res, "score mate", strconv.Itoa(o.Mate))
	} else {
		res = append(res, "score cp", strconv.Itoa(o.Score))
	}
	if len(o.PV) > 0 {
		res = append(res, "pv")
		for _, move := range o.PV {
			res = append(res, c.moveNotation.Encode(c.position, move))
		}
	}
	if o.time > 0 {
		res = append(res, "time", strconv.FormatInt(o.time.Milliseconds(), 10))
	}

	return "info " + strings.Join(res, " ")
}

// responseOption represents an "option" command.
//
//nolint:govet
type responseOption struct {
	Type    optionType
	Name    string
	Default string
	Min     string
	Max     string
}

func (o responseOption) format(_ *Controller) string {
	switch o.Type {
	case integerOptionType:
		var min, max string
		if len(o.Min) > 0 {
			min = " min " + o.Min
		}
		if len(o.Max) > 0 {
			max = " max " + o.Max
		}
		return fmt.Sprintf(
			"option name %s type spin default %s%s%s",
			o.Name, o.Default, min, max,
		)
	case booleanOptionType:
		return fmt.Sprintf(
			"option name %s type check default %s",
			o.Name, o.Default,
		)
	default:
		return ""
	}
}
