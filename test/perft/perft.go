// This package provides the capability to run perft tests.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/leonhfr/orca/chess"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("expected arguments to be provided")
		os.Exit(1)
	}

	r, err := perft(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(r)
}

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

			if ok := pos.MakeMove(m); !ok {
				return chess.PerftResult{}, fmt.Errorf("could not play move %v", m.String())
			}
		}
	}

	return pos.Perft(int(depth)), nil
}
