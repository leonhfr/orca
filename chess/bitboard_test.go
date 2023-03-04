package chess

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitboard_String(t *testing.T) {
	var bb bitboard
	for _, sq := range []Square{A1, H1, A8, H8} {
		bb ^= sq.bitboard()
	}
	expected := "10000001" + strings.Repeat("0", 48) + "10000001"
	assert.Equal(t, bb.String(), expected)
}
