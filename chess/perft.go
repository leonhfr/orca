package chess

import (
	"fmt"
	"sort"
	"strings"
)

// PerftResult contains the result of a perft test.
type PerftResult struct {
	moves []perftMove
	nodes int
}

// perfMove contains a single move and the number of its descendant nodes.
type perftMove struct {
	move  Move
	nodes int
}

// Perft performs a perft test.
func (pos *Position) Perft(depth int) PerftResult {
	hash := pos.Hash()
	pawnHash := pos.PawnHash()

	var result PerftResult
	checkData, _ := pos.InCheck()
	for _, m := range pos.PseudoMoves(checkData) {
		if meta, ok := pos.MakeMove(m); ok {
			nodes := perft(pos, depth-1)
			result.moves = append(result.moves, perftMove{
				move:  m,
				nodes: nodes,
			})
			result.nodes += nodes
			pos.UnmakeMove(m, meta, hash, pawnHash)
		}
	}

	sort.Slice(result.moves, func(i, j int) bool {
		return strings.Compare(result.moves[i].String(), result.moves[j].String()) < 0
	})

	return result
}

// perft returns the number of nodes until the given depth.
func perft(pos *Position, depth int) int {
	hash := pos.Hash()
	pawnHash := pos.PawnHash()

	if depth <= 0 {
		return 1
	}

	checkData, _ := pos.InCheck()
	moves := pos.PseudoMoves(checkData)

	if depth == 1 {
		var nodes int
		for _, m := range moves {
			if meta, ok := pos.MakeMove(m); ok {
				nodes++
				pos.UnmakeMove(m, meta, hash, pawnHash)
			}
		}
		return nodes
	}

	var nodes int
	for _, m := range moves {
		if meta, ok := pos.MakeMove(m); ok {
			nodes += perft(pos, depth-1)
			pos.UnmakeMove(m, meta, hash, pawnHash)
		}
	}
	return nodes
}

// String implements fmt.Stringer.
//
// Returns the perft result as an output accepted by the perftree cli.
func (pr PerftResult) String() string {
	b := &strings.Builder{}
	for _, m := range pr.moves {
		fmt.Fprintln(b, m.String())
	}
	fmt.Fprintf(b, "\n%d", pr.nodes)
	return b.String()
}

// String implements fmt.Stringer.
//
// Returns the move in UCI notation followed by a space, then the number of nodes
// of which this move is an ancestor.
func (pm perftMove) String() string {
	return fmt.Sprintf("%v %d", pm.move.String(), pm.nodes)
}
