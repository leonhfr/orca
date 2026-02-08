package search

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leonhfr/orca/chess"
)

// compile time check that noTable implements transpositionTable.
var _ transpositionTable = noTable{}

// compile time check that arrayTable implements transpositionTable.
var _ transpositionTable = (*arrayTable)(nil)

// compile time check that hashMapTable implements transpositionTable.
var _ transpositionTable = (*hashMapTable)(nil)

func TestNewTable(t *testing.T) {
	t.Parallel()
	table := newArrayTable(1)
	defer table.close()

	require.Equal(t, uint64(43690), table.length)
}

func TestTableGet(t *testing.T) {
	t.Parallel()
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	table := newArrayTable(1)
	defer table.close()

	_, ok := table.get(hash)
	require.False(t, ok)
}

func TestTableSet(t *testing.T) {
	t.Parallel()
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	//nolint:gosec
	score, depth := rand.Int31(), uint8(rand.Uint32())
	best, nt := chess.NoMove, exact

	table := newArrayTable(1)
	defer table.close()

	table.set(hash, best, score, nt, depth)
	entry, ok := table.get(hash)

	require.True(t, ok)
	require.Equal(t, best, entry.best)
	require.Equal(t, score, entry.score())
	require.Equal(t, nt, entry.nodeType())
	require.Equal(t, depth, entry.depth())
}

// hashMapTable uses a map as backend. Intended to be used for tests.
//
// Implements the transpositionTable interface.
type hashMapTable struct {
	table map[chess.Hash]searchEntry
}

// newHashMapTable returns a new HashMapTable.
// Does not return any cached entries, but permits to collect the principal variation.
// Use in tests only.
func newHashMapTable() *hashMapTable {
	return &hashMapTable{
		table: make(map[chess.Hash]searchEntry),
	}
}

// Implements the transpositionTable interface.
func (hm *hashMapTable) inc() {
}

// Implements the transpositionTable interface.
func (hm *hashMapTable) get(_ chess.Hash) (searchEntry, bool) {
	return searchEntry{}, false
}

// Implements the transpositionTable interface.
func (hm *hashMapTable) set(hash chess.Hash, best chess.Move, score int32, nt nodeType, depth uint8) {
	hm.table[hash] = searchEntry{
		uint64(hash),
		best,
		serializeSearchData(score, nt, depth, 0),
	}
}

// Implements the transpositionTable interface.
func (hm *hashMapTable) principalVariation(pos *chess.Position) []chess.Move {
	type unmakeMove struct {
		move     chess.Move
		meta     chess.Metadata
		hash     chess.Hash
		pawnHash chess.Hash
	}

	pv := make([]chess.Move, 0, 10)
	unmakeMoveStack := make([]unmakeMove, 0, 10)

	for hash := pos.Hash(); ; hash = pos.Hash() {
		entry, inCache := hm.table[hash]
		if !inCache || entry.best == chess.NoMove {
			break
		}

		meta := pos.Metadata()
		pawnHash := pos.PawnHash()
		if ok := pos.MakeMove(entry.best); !ok {
			break
		}

		pv = append(pv, entry.best)
		unmakeMoveStack = append(unmakeMoveStack, unmakeMove{
			move:     entry.best,
			meta:     meta,
			hash:     hash,
			pawnHash: pawnHash,
		})
	}

	length := len(unmakeMoveStack)
	for i := range length {
		um := unmakeMoveStack[length-i-1]
		pos.UnmakeMove(um.move, um.meta, um.hash, um.pawnHash)
	}

	return pv
}

// Implements the transpositionTable interface.
func (hm *hashMapTable) close() {
	hm.table = nil
}
