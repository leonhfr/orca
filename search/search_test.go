package search

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/orca/chess"
)

func TestNew(t *testing.T) {
	t.Parallel()
	e := NewEngine()
	assert.Equal(t, 64, e.tableSize)
}

func TestWithTableSize(t *testing.T) {
	t.Parallel()
	e := NewEngine(WithTableSize(128))
	assert.Equal(t, 128, e.tableSize)
}

func TestWithOwnBook(t *testing.T) {
	t.Parallel()
	e := NewEngine(WithOwnBook(true))
	assert.True(t, e.ownBook)
}

func TestInit(t *testing.T) {
	t.Parallel()
	engine := NewEngine()
	assert.Equal(t, noTable{}, engine.table)

	err := engine.Init()

	assert.Nil(t, err)
	assert.IsType(t, &arrayTable{}, engine.table)
}

func TestSearch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		fen     string
		depth   int
		nodes   int
		book    bool
		outputs []Output
	}{
		{
			name:  "mate in 1",
			fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
			depth: 2,
			outputs: []Output{
				{Depth: 1, Nodes: 91, Score: mate - 1, Mate: 1, PV: []chess.Move{0xd206c1836d}},
				{Depth: 2, Nodes: 338, Score: mate - 1, Mate: 1, PV: []chess.Move{0xd206c1836d}},
			},
		},
		{
			name:  "lasker trap without opening book",
			fen:   "rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/5N2/PP2PPPP/RNBQKB1R b KQkq - 1 3",
			depth: 2,
			outputs: []Output{
				{Depth: 1, Nodes: 115, Score: 0, Mate: 0, PV: []chess.Move{0x13602c106a3}},
				{Depth: 2, Nodes: 924, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc2ab9, 0x12c02c018da}},
			},
		},
		{
			name:    "lasker trap with opening book",
			fen:     "rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/5N2/PP2PPPP/RNBQKB1R b KQkq - 1 3",
			depth:   2,
			book:    true,
			outputs: []Output{{PV: []chess.Move{0x1cc2b7e}, Depth: 1, Nodes: 1, Score: 0, Mate: 0}},
		},
		{
			name:  "nodes limit",
			fen:   "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			nodes: 16384,
			depth: 5,
			outputs: []Output{
				{Depth: 1, Nodes: 747, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc38d2}},
				{Depth: 2, Nodes: 5805, Score: 23, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed}},
				{Depth: 3, Nodes: 19798, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed, 0x19002c85d26}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			engine := NewEngine()
			_ = engine.Init()
			if tt.book {
				engine.ownBook = true
			}
			engine.table = newHashMapTable()
			engine.pawnTable = noPawnTable{}
			output := engine.Search(context.Background(), unsafeFEN(tt.fen), tt.depth, tt.nodes)
			outputs := make([]Output, 0, tt.depth)
			for o := range output {
				outputs = append(outputs, o)
			}

			assert.Equal(t, tt.outputs, outputs)
		})
	}
}

func TestCachedSearch(t *testing.T) {
	t.Parallel()
	fen := "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	depth := 8

	tests := []struct {
		name   string
		cached bool
		want   []Output
	}{
		{
			"not cached",
			false,
			[]Output{
				{Depth: 1, Nodes: 747, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc38d2}},
				{Depth: 2, Nodes: 5805, Score: 23, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed}},
				{Depth: 3, Nodes: 19798, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed, 0x19002c85d26}},
				{Depth: 4, Nodes: 59913, Score: 4, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c50b76, 0x6401cc15cf, 0x6401cc26ea}},
				{Depth: 5, Nodes: 184981, Score: 3, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c50b76, 0x6401cc1649, 0x12702c3455e, 0x14f02c4954c}},
				{Depth: 6, Nodes: 634537, Score: 3, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c50b76, 0x6401cc1649, 0x12702c3455e, 0x14f02c4954c}},
				{Depth: 7, Nodes: 3300539, Score: 1, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c58b74, 0x6401cc38d2, 0x6401cc2e6a, 0x19006c83b63, 0x14a02c30b76, 0x11802c03915}},
				{Depth: 8, Nodes: 18846865, Score: 1, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed, 0x19002c85d26, 0x14f02c52d23, 0x11302c05a1a, 0x14f02c50a31, 0x6401cc92cc}},
			},
		},
		{
			"cached",
			true,
			[]Output{
				{Depth: 1, Nodes: 747, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc38d2}},
				{Depth: 2, Nodes: 5805, Score: 23, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed}},
				{Depth: 3, Nodes: 19798, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed, 0x19002c85d26}},
				{Depth: 4, Nodes: 42379, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed, 0x19002c85d26}},
				{Depth: 5, Nodes: 167447, Score: 3, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c50b76, 0x6401cc1649, 0x12702c3455e, 0x14f02c4954c}},
				{Depth: 6, Nodes: 613541, Score: 3, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c50b76, 0x6401cc1649, 0x12702c3455e, 0x14f02c4954c}},
				{Depth: 7, Nodes: 3279541, Score: 1, Mate: 0, PV: []chess.Move{0x12702c25b66, 0x14f02c58b74, 0x6401cc38d2, 0x6401cc2e6a, 0x19006c83b63, 0x14a02c30b76, 0x11802c03915}},
				{Depth: 8, Nodes: 18821824, Score: 1, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x12c02c328ed, 0x19002c85d26, 0x14f02c52d23, 0x11302c05a1a, 0x14f02c50a31, 0x6401cc92cc}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			engine := NewEngine()
			_ = engine.Init()
			if !tt.cached {
				engine.table = newHashMapTable()
				engine.pawnTable = noPawnTable{}
			}
			pos := unsafeFEN(fen)
			output := engine.Search(context.Background(), pos, depth, math.MaxInt)
			outputs := make([]Output, 0, depth)
			for o := range output {
				outputs = append(outputs, o)
			}

			require.Equal(t, len(tt.want), len(outputs))
			for i, o := range outputs {
				wantMoves := make([]string, 0, len(tt.want[i].PV))
				gotMoves := make([]string, 0, len(o.PV))
				for _, m := range tt.want[i].PV {
					wantMoves = append(wantMoves, m.String())
				}
				for _, m := range o.PV {
					gotMoves = append(gotMoves, m.String())
				}
				assert.Equal(t, tt.want[i], o, fmt.Sprintf("want moves %s, got %s", strings.Join(wantMoves, ", "), strings.Join(gotMoves, ", ")))
			}
		})
	}
}

func BenchmarkCachedSearch(b *testing.B) {
	fen := "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	depth := 6

	benchs := []struct {
		name   string
		cached bool
	}{
		{"not cached", false},
		{"cached", true},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				engine := NewEngine()
				_ = engine.Init()
				if !bb.cached {
					engine.table = noTable{}
					engine.pawnTable = noPawnTable{}
				}
				pos := unsafeFEN(fen)
				b.StartTimer()

				output := engine.Search(context.Background(), pos, depth, math.MaxInt)
				for o := range output {
					_ = o
				}
			}
		})
	}
}
