package search

import (
	"unsafe"

	"github.com/dgraph-io/ristretto"

	"github.com/leonhfr/orca/chess"
)

// transpositionTable holds a transposition table.
// Stores results of previously performed searches.
//
// Maps chess.Hash to tableEntry structs.
type transpositionTable struct {
	cache *ristretto.Cache
}

// tableEntry holds a search result entry.
// Only essential information is retained.
type tableEntry struct {
	score    int
	depth    int
	nodeType nodeType
	best     bool
}

// nodeType represents the score bounds for this entry.
type nodeType uint8

const (
	noBounds   nodeType = iota // score with undefined bounds
	lowerBound                 // lower bound score (cut node)
	upperBound                 // upper bound score (all node)
	exact                      // exact score (pv node)
)

// newTable returns a new transpositionTable.
//
// Takes the desired table size in Megabytes as argument.
func newTable(size int) (*transpositionTable, error) {
	entrySize := uint64(unsafe.Sizeof(tableEntry{}))
	maxCost := int64(1024 * 1024 * uint64(size) / entrySize)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10 * maxCost,
		MaxCost:     maxCost,
		BufferItems: 64,
		KeyToHash: func(key any) (uint64, uint64) {
			return uint64(key.(chess.Hash)), 0
		},
	})
	if err != nil {
		return nil, err
	}

	return &transpositionTable{cache: cache}, nil
}

// get returns the entry (if any) for the given hash.
// The boolean is true if the entry was found, false if not.
//
// Assumes that the table has been initialized.
func (tt *transpositionTable) get(key chess.Hash) (tableEntry, bool) {
	entry, found := tt.cache.Get(key)
	if !found {
		return tableEntry{}, false
	}
	return entry.(tableEntry), true
}

// set adds an entry to the table for the given hash.
// If an entry already exists, it is replaced.
// The addition is not guaranteed.
//
// Assumes that the table has been initialized.
func (tt *transpositionTable) set(key chess.Hash, entry tableEntry) {
	tt.cache.Set(key, entry, 1)
}

// close initiates a graceful shutdown of the transposition table.
//
// Assumes that the table has been initialized.
func (tt *transpositionTable) close() {
	tt.cache.Close()
}
