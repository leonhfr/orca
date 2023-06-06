package search

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leonhfr/orca/chess"
)

// compile time check that noPawnTable implements transpositionPawnTable.
var _ transpositionPawnTable = noPawnTable{}

// compile time check that arrayPawnTable implements transpositionPawnTable.
var _ transpositionPawnTable = (*arrayPawnTable)(nil)

func TestNewPawnTable(t *testing.T) {
	table := newArrayPawnTable(1)
	defer table.close()

	require.Equal(t, uint64(64), table.length)
}

func TestPawnTableGet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	table := newArrayTable(1)
	defer table.close()

	_, ok := table.get(hash)
	require.False(t, ok)
}

func TestPawnTableSet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	//nolint:gosec
	want := pawnEntry{
		hash: hash,
		mg:   rand.Int31(),
		eg:   rand.Int31(),
	}

	table := newArrayPawnTable(1)
	defer table.close()

	table.set(hash, want.mg, want.eg)

	entry, ok := table.get(hash)

	require.True(t, ok)
	require.Equal(t, want.mg, entry.mg)
	require.Equal(t, want.eg, entry.eg)
}
