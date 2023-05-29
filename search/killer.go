package search

import "github.com/leonhfr/orca/chess"

// killerList contains a list of killer moves ordered by depth (plies).
//
// Killer moves are quiet moves that caused a beta cutoff in a sibling node.
type killerList struct {
	entries [][2]chess.Move
}

// newKillerList returns a new killerList.
func newKillerList() *killerList {
	return &killerList{}
}

// increaseDepth sets the depth. Intended to be used during iterative deepening.
func (kl *killerList) increaseDepth(depth int) {
	add := depth - len(kl.entries)
	for i := 0; i < add; i++ {
		kl.entries = append(kl.entries, [2]chess.Move{})
	}
}

// get returns the killer moves at this depth.
func (kl *killerList) get(depth uint8) [2]chess.Move {
	index := len(kl.entries) - int(depth)
	if index < 0 || len(kl.entries) <= index {
		return [2]chess.Move{}
	}
	return kl.entries[index]
}

// set inserts the killer move at the given depth.
//
// If moves are already present, the older ones are pushed out.
// Guarantees that the moves will be different.
func (kl *killerList) set(move chess.Move, depth uint8) {
	index := len(kl.entries) - int(depth)
	if index < 0 || len(kl.entries) <= index {
		return
	}

	if kl.entries[index][0] == move || kl.entries[index][1] == move {
		return
	}

	kl.entries[index][0], kl.entries[index][1] = move, kl.entries[index][0]
}
