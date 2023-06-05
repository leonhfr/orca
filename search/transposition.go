package search

import (
	"math/bits"
	"unsafe"

	"github.com/leonhfr/orca/chess"
)

// transpositionTable is the interface that transposition tables should implement.
//
// Allows the storing of results of previously performed searches by mapping
// chess.Hash to searchEntry structs.
type transpositionTable interface {
	// inc increases the epoch.
	inc()
	// get returns the entry (if any) for the given hash
	// and a boolean representing whether the value was found or not.
	get(key chess.Hash) (searchEntry, bool)
	// set adds an entry to the table for the given hash.
	// If an entry already exists, it is replaced.
	// The addition is not guaranteed.
	set(key chess.Hash, entry searchEntry)
	// close initiates a graceful shutdown of the transposition table.
	close()
}

// searchEntry holds a search result entry.
// Only essential information is retained.
type searchEntry struct {
	best     chess.Move
	hash     chess.Hash
	score    int32
	depth    uint8
	nodeType nodeType
	epoch    uint8
}

func (se searchEntry) quality() uint8 {
	return se.epoch + se.depth/3
}

// nodeType represents the score bounds for this entry.
type nodeType uint8

const (
	noEntry    nodeType = iota // no entry exists
	lowerBound                 // lower bound score (cut node)
	upperBound                 // upper bound score (all node)
	exact                      // exact score (pv node)
)

// noTable does not store anything at all.
//
// Implements the transpositionTable interface.
type noTable struct{}

func (noTable) inc()                                 {}                              // implements transpositionTable.
func (noTable) get(_ chess.Hash) (searchEntry, bool) { return searchEntry{}, false } // implements transpositionTable.
func (noTable) set(_ chess.Hash, _ searchEntry)      {}                              // implements transpositionTable.
func (noTable) close()                               {}                              // implements transpositionTable.

// arrayTable uses an array as backend.
//
// Implements the transpositionTable interface.
type arrayTable struct {
	table  []searchEntry
	length uint64
	epoch  uint8
}

// newArrayTable returns a new arrayTable.
//
// Takes the desired table size in Megabytes as argument.
func newArrayTable(size int) *arrayTable {
	entrySize := uint64(unsafe.Sizeof(searchEntry{}))
	length := 1024 * 1024 * uint64(size) / entrySize

	return &arrayTable{
		table:  make([]searchEntry, length),
		length: length,
	}
}

// Implements the transpositionTable interface.
func (ar *arrayTable) inc() {
	ar.epoch++
}

// Implements the transpositionTable interface.
func (ar *arrayTable) get(key chess.Hash) (searchEntry, bool) {
	entry := ar.table[ar.hash(key)]
	return entry, entry.nodeType != noEntry && entry.hash == key
}

// Implements the transpositionTable interface.
func (ar *arrayTable) set(key chess.Hash, entry searchEntry) {
	index := ar.hash(key)
	cached := ar.table[index]
	entry.epoch = ar.epoch
	if entry.quality() >= cached.quality() {
		ar.table[index] = entry
	}
}

// Implements the transpositionTable interface.
func (ar *arrayTable) close() {
	ar.table = nil
}

// hash is the hash function used by the array table.
func (ar *arrayTable) hash(key chess.Hash) uint64 {
	// fast indexing function from Daniel Lemire's blog post
	// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	index, _ := bits.Mul64(uint64(key), ar.length)
	return index
}
