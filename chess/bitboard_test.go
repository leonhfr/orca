package chess

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitboard_String(t *testing.T) {
	t.Parallel()
	var bb bitboard
	for _, sq := range []Square{A1, H1, A8, H8} {
		bb ^= sq.bitboard()
	}
	expected := "10000001" + strings.Repeat("0", 48) + "10000001"
	assert.Equal(t, bb.String(), expected)
}

// mapping returns the list of squares set to 1.
func (b bitboard) mapping() []Square {
	if b == 0 {
		return nil
	}
	squares := make([]Square, 0, 8)
	for b > 0 {
		squares = append(squares, b.scanForward())
		b = b.resetLSB()
	}
	return squares
}
