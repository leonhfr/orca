package search

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/leonhfr/orca/chess"
)

var wait = 10 * time.Millisecond

func TestNewRistrettoTable(t *testing.T) {
	table, err := newRistrettoTable(1)
	defer table.close()

	require.Nil(t, err)
	require.Equal(t, int64(87381), table.cache.MaxCost())
}

func TestTableGet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	table, err := newRistrettoTable(1)
	defer table.close()

	require.Nil(t, err)

	_, ok := table.get(hash)
	require.False(t, ok)
}

func TestTableSet(t *testing.T) {
	//nolint:gosec
	hash := chess.Hash(rand.Uint64())
	//nolint:gosec
	want := searchEntry{
		score:    rand.Int31(),
		depth:    uint8(rand.Uint32()),
		nodeType: exact,
	}

	table, err := newRistrettoTable(1)
	defer table.close()

	require.Nil(t, err)

	table.set(hash, want)
	time.Sleep(wait)
	entry, ok := table.get(hash)

	require.True(t, ok)
	require.Equal(t, want, entry)
}
