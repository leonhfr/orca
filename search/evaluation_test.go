package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

var testPositions = []struct {
	fen   string
	score int32
}{
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0},
	{"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23", 10},
	{"r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9", -25},
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0},
	{"r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/R3K2R w KQkq - 3 10", -13},
	{"r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", 0},
	{"r1bqkbnr/ppp1p1pp/2n5/3pPp2/8/5N2/PPPP1PPP/RNBQKB1R w KQkq f6 0 4", 24},
	{"r1bqkbnr/ppp1p1pp/2n5/3pPp2/3N4/8/PPPP1PPP/RNBQKB1R b KQkq - 1 4", -20},
	{"r7/1Pp5/2P3p1/8/6pb/4p1kB/4P1p1/6K1 w - - 0 1", -723},
}

func TestEvaluate(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.fen, func(t *testing.T) {
			si := newSearchInfo(noTable{}, noPawnTable{})
			pos := unsafeFEN(tt.fen)
			assert.Equal(t, tt.score, si.evaluate(pos))
		})
	}
}

func BenchmarkEvaluate(b *testing.B) {
	for _, bb := range testPositions {
		b.Run(bb.fen, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			si := newSearchInfo(noTable{}, noPawnTable{})
			for n := 0; n < b.N; n++ {
				si.evaluate(pos)
			}
		})
	}
}

func unsafeFEN(fen string) *chess.Position {
	p, err := chess.NewPosition(fen)
	if err != nil {
		panic(err)
	}
	return p
}
