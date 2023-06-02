package search

import "github.com/leonhfr/orca/chess"

// transpositionPawnTable is the interface that pawn transposition tables should implement.
//
// Allows the storing of pawn evaluation results by mapping chess.Hash to pawnEntry structs.
type transpositionPawnTable interface {
	// get returns the entry (if any) for the given hash
	// and a boolean representing whether the value was found or not.
	get(key chess.Hash) (pawnEntry, bool)
	// set adds an entry to the table for the given hash.
	// If an entry already exists, it is replaced.
	set(key chess.Hash, entry pawnEntry)
	// close initiates a graceful shutdown of the pawn transposition table.
	close()
}

// pawnEntry hols a pawn evaluation.
type pawnEntry struct {
	// hash chess.Hash
	mg int32
	eg int32
}

// noPawnTable does not store anything at all.
type noPawnTable struct{}

func (noPawnTable) get(_ chess.Hash) (pawnEntry, bool) { return pawnEntry{}, false } // implements transpositionPawnTable.
func (noPawnTable) set(_ chess.Hash, _ pawnEntry)      {}                            // implements transpositionPawnTable.
func (noPawnTable) close()                             {}                            // implements transpositionPawnTable.
