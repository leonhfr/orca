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
	get(hash chess.Hash) (pawnEntry, bool)
	// set adds an entry to the table for the given hash.
	// If an entry already exists, it is replaced.
	set(hash chess.Hash, mg, eg int32)
	// close initiates a graceful shutdown of the pawn transposition table.
	close()
}

// pawnEntry hols a pawn evaluation.
type pawnEntry struct {
	hash uint64
	data uint64
}

// serializePawnData serializes a pawn entry data.
func serializePawnData(mg, eg int32) uint64 {
	return uint64(eg)<<32 + uint64(mg)
}

// mg returns the middle game score.
func (pe pawnEntry) mg() int32 {
	return int32(uint32(pe.data))
}

// eg returns the endgame score.
func (pe pawnEntry) eg() int32 {
	return int32(uint32(uint64(pe.data+1<<31) >> 32))
}

// noPawnTable does not store anything at all.
type noPawnTable struct{}

func (noPawnTable) get(_ chess.Hash) (pawnEntry, bool) { return pawnEntry{}, false } // implements transpositionPawnTable.
func (noPawnTable) set(_ chess.Hash, _, _ int32)       {}                            // implements transpositionPawnTable.
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
func (ar *arrayPawnTable) get(hash chess.Hash) (pawnEntry, bool) {
	entry := ar.table[ar.hash(hash)]
	return entry, chess.Hash(entry.hash^entry.data) == hash
}

// Implements the transpositionPawnTable interface.
func (ar *arrayPawnTable) set(hash chess.Hash, mg, eg int32) {
	index := ar.hash(hash)
	data := serializePawnData(mg, eg)

	ar.table[index] = pawnEntry{
		uint64(hash) ^ data,
		data,
	}
}

// Implements the transpositionPawnTable interface.
func (ar *arrayPawnTable) close() {
	ar.table = nil
}

// hash is the hash function used by the array table.
func (ar *arrayPawnTable) hash(hash chess.Hash) uint64 {
	index, _ := bits.Mul64(uint64(hash), ar.length)
	return index
}
