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

func TestNewTable(t *testing.T) {
	table := newArrayTable(1)
	defer table.close()

	require.Equal(t, uint64(43690), table.length)
}

func TestTableGet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	table := newArrayTable(1)
	defer table.close()

	_, ok := table.get(hash)
	require.False(t, ok)
}

func TestTableSet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	//nolint:gosec
	want := searchEntry{
		hash:     hash,
		score:    rand.Int31(),
		depth:    uint8(rand.Uint32()),
		nodeType: exact,
	}

	table := newArrayTable(1)
	defer table.close()

	table.set(hash, want)
	entry, ok := table.get(hash)

	require.True(t, ok)
	require.Equal(t, want, entry)
}
