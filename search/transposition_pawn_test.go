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
	t.Parallel()
	table := newArrayPawnTable(1)
	defer table.close()

	require.Equal(t, uint64(64), table.length)
}

func TestPawnTableGet(t *testing.T) {
	t.Parallel()
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	table := newArrayTable(1)
	defer table.close()

	_, ok := table.get(hash)
	require.False(t, ok)
}

func TestPawnTableSet(t *testing.T) {
	t.Parallel()
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	//nolint:gosec
	mg, eg := rand.Int31(), rand.Int31()

	table := newArrayPawnTable(1)
	defer table.close()

	table.set(hash, mg, eg)

	entry, ok := table.get(hash)

	require.True(t, ok)
	require.Equal(t, mg, entry.mg())
	require.Equal(t, eg, entry.eg())
}
