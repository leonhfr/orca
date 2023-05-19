package search

import (
	"unsafe"

	"github.com/dgraph-io/ristretto"

	"github.com/leonhfr/orca/chess"
)

// transpositionTable is the interface that transposition tables should implement.
// Allows the storing of results of previously performed searches by mapping
// chess.Hash to tableEntry structs.
type transpositionTable interface {
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
	score    int
	depth    int
	best     chess.Move
	nodeType nodeType
}

// nodeType represents the score bounds for this entry.
type nodeType uint8

const (
	noBounds   nodeType = iota // score with undefined bounds
	lowerBound                 // lower bound score (cut node)
	upperBound                 // upper bound score (all node)
	exact                      // exact score (pv node)
)

// noTable does not store anything at all.
type noTable struct{}

func (noTable) get(_ chess.Hash) (searchEntry, bool) { return searchEntry{}, false } // implements transpositionTable.
func (noTable) set(_ chess.Hash, _ searchEntry)      {}                              // implements transpositionTable.
func (noTable) close()                               {}                              // implements transpositionTable.

// ristrettoTable uses dgraph-io/ristretto as backend.
//
// Implements the transpositionTable interface.
type ristrettoTable struct {
	cache *ristretto.Cache
}

// newRistrettoTable returns a new ristrettoTable.
//
// Takes the desired table size in Megabytes as argument.
func newRistrettoTable(size int) (*ristrettoTable, error) {
	entrySize := uint64(unsafe.Sizeof(searchEntry{}))
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

	return &ristrettoTable{cache: cache}, nil
}

// Implements the transpositionTable interface.
func (tt *ristrettoTable) get(key chess.Hash) (searchEntry, bool) {
	entry, found := tt.cache.Get(key)
	if !found {
		return searchEntry{}, false
	}
	return entry.(searchEntry), true
}

// Implements the transpositionTable interface.
func (tt *ristrettoTable) set(key chess.Hash, entry searchEntry) {
	tt.cache.Set(key, entry, 1)
}

// Implements the transpositionTable interface.
func (tt *ristrettoTable) close() {
	tt.cache.Close()
}
