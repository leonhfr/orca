package search

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/leonhfr/orca/chess"
)

var wait = 10 * time.Millisecond

func TestNewTable(t *testing.T) {
	table, err := newTable(1)
	defer table.close()

	require.Nil(t, err)
	require.Equal(t, int64(43690), table.cache.MaxCost())
}

func TestTableGet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	table, err := newTable(1)
	defer table.close()

	require.Nil(t, err)

	_, ok := table.get(hash)
	require.False(t, ok)
}

func TestTableSet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	//nolint:gosec
	want := tableEntry{
		score:    rand.Int(),
		depth:    rand.Int(),
		nodeType: exact,
		best:     true,
	}

	table, err := newTable(1)
	defer table.close()

	require.Nil(t, err)

	table.set(hash, want)
	time.Sleep(wait)
	entry, ok := table.get(hash)

	require.True(t, ok)
	require.Equal(t, want, entry)
}