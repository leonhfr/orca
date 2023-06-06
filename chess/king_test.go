package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type kingCallbackArgs struct {
	sq Square
	sd int
}

var kingPieceMapTests = []struct {
	name  string
	fen   string
	kings [2]kingCallbackArgs // indexed by color
}{
	{
		"perfect shields queen side",
		"2k5/1ppp4/8/8/8/8/1PPP4/2K5 w - - 0 1",
		[2]kingCallbackArgs{{C8, 0}, {C1, 0}},
	},
	{
		"perfect shields king side",
		"6k1/5ppp/8/8/8/8/5PPP/6K1 w - - 0 1",
		[2]kingCallbackArgs{{G8, 0}, {G1, 0}},
	},
	{
		"imperfect shields",
		"6k1/4p3/6p1/8/6P1/8/5P1P/6K1 w - - 0 1",
		[2]kingCallbackArgs{{G8, 2}, {G1, 1}},
	},
	{
		"does not need shield",
		"4k3/3pp2p/8/8/6P1/8/5P1P/4K3 w - - 0 1",
		[2]kingCallbackArgs{{E8, 0}, {E1, 0}},
	},
}

func TestKingMap(t *testing.T) {
	for _, tt := range kingPieceMapTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			pos.KingMap(func(p Piece, sq Square, shieldDefects int) {
				assert.Equal(t, tt.kings[p.Color()].sq, sq)
				assert.Equal(t, tt.kings[p.Color()].sd, shieldDefects)
			})
		})
	}
}

func BenchmarKingMap(b *testing.B) {
	for _, bb := range kingPieceMapTests {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				pos.KingMap(func(p Piece, sq Square, shieldDefects int) {
					_ = 1
				})
			}
		})
	}
}
