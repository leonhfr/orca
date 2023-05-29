package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/leonhfr/orca/chess"
)

// perft expects the arguments to be depth, then fen,
// and finally an optional space-separated move list.
func perft(args []string) (chess.PerftResult, error) {
	if len(args) < 2 || 3 < len(args) {
		return chess.PerftResult{}, errors.New("expected arguments to be depth, then fen")
	}

	depth, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return chess.PerftResult{}, err
	}

	pos, err := chess.NewPosition(args[1])
	if err != nil {
		return chess.PerftResult{}, err
	}

	if len(args) == 3 {
		for _, move := range strings.Fields(args[2]) {
			m, err := chess.NewMove(pos, move)
			if err != nil {
				return chess.PerftResult{}, err
			}

			_, ok := pos.MakeMove(m)
			if !ok {
				return chess.PerftResult{}, fmt.Errorf("could not play move %v", m.String())
			}
		}
	}

	return pos.Perft(int(depth)), nil
}