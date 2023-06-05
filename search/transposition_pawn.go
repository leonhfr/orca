package search

import (
	"math/bits"
	"unsafe"

	"github.com/leonhfr/orca/chess"
)

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
	hash chess.Hash
	mg   int32
	eg   int32
}

// noPawnTable does not store anything at all.
type noPawnTable struct{}

func (noPawnTable) get(_ chess.Hash) (pawnEntry, bool) { return pawnEntry{}, false } // implements transpositionPawnTable.
func (noPawnTable) set(_ chess.Hash, _ pawnEntry)      {}                            // implements transpositionPawnTable.
func (noPawnTable) close()                             {}                            // implements transpositionPawnTable.

// arrayPawnTable uses an array as backend.
//
// Implements the transpositionPawnTable interface.
type arrayPawnTable struct {
	table  []pawnEntry
	length uint64
}

// newArrayPawnTable returns a new arrayPawnTable.
//
// Takes the desired table size in Kilobytes as argument.
func newArrayPawnTable(size int) *arrayPawnTable {
	entrySize := uint64(unsafe.Sizeof(pawnEntry{}))
	length := 1024 * uint64(size) / entrySize

	return &arrayPawnTable{
		table:  make([]pawnEntry, length),
		length: length,
	}
}

// Implements the transpositionPawnTable interface.
func (ar *arrayPawnTable) get(key chess.Hash) (pawnEntry, bool) {
	entry := ar.table[ar.hash(key)]
	return entry, entry.hash == key
}

// Implements the transpositionPawnTable interface.
func (ar *arrayPawnTable) set(key chess.Hash, entry pawnEntry) {
	index := ar.hash(key)
	ar.table[index] = entry
}

// Implements the transpositionPawnTable interface.
func (ar *arrayPawnTable) close() {
	ar.table = nil
}

// hash is the hash function used by the array table.
func (ar *arrayPawnTable) hash(key chess.Hash) uint64 {
	index, _ := bits.Mul64(uint64(key), ar.length)
	return index
}
