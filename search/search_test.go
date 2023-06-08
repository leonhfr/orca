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
	"github.com/leonhfr/orca/uci"
)

// compile time check that Engine implements uci.Engine.
var _ uci.Engine = (*Engine)(nil)

func TestNew(t *testing.T) {
	e := NewEngine()
	assert.Equal(t, 64, e.tableSize)
}

func TestWithTableSize(t *testing.T) {
	e := NewEngine(WithTableSize(128))
	assert.Equal(t, 128, e.tableSize)
}

func TestWithOwnBook(t *testing.T) {
	e := NewEngine(WithOwnBook(true))
	assert.True(t, e.ownBook)
}

func TestInit(t *testing.T) {
	engine := NewEngine()
	assert.Equal(t, noTable{}, engine.table)

	err := engine.Init()

	assert.Nil(t, err)
	assert.IsType(t, &arrayTable{}, engine.table)
}

func TestOptions(t *testing.T) {
	e := NewEngine()
	options := e.Options()
	assert.Equal(t, []uci.Option{
		{
			Type:    uci.OptionInteger,
			Name:    "Hash",
			Default: "64",
			Min:     "1",
			Max:     "16384",
		},
		{
			Type:    uci.OptionBoolean,
			Name:    "OwnBook",
			Default: "false",
		},
	}, options)
}

func TestSetOption(t *testing.T) {
	type args struct {
		name, value string
		tableSize   int
		err         error
	}

	tests := []struct {
		name string
		args args
	}{
		{
			"option exists",
			args{"Hash", "128", 128, nil},
		},
		{
			"option outside bounds",
			args{"Hash", "0", 64, errOutsideBound},
		},
		{
			"option does not exist",
			args{"Whatever", "Whatever", 64, errOptionName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine()
			err := e.SetOption(tt.args.name, tt.args.value)
			assert.Equal(t, tt.args.tableSize, e.tableSize)
			assert.Equal(t, tt.args.err, err)
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		depth   int
		nodes   int
		book    bool
		outputs []uci.Output
	}{
		{
			name:  "mate in 1",
			fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
			depth: 2,
			outputs: []uci.Output{
				{Depth: 1, Nodes: 98, Score: mate - 1, Mate: 1, PV: []chess.Move{0x13006c1836d}},
				{Depth: 2, Nodes: 347, Score: mate - 1, Mate: 1, PV: []chess.Move{0x13006c1836d}},
			},
		},
		{
			name:  "lasker trap without opening book",
			fen:   "rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/5N2/PP2PPPP/RNBQKB1R b KQkq - 1 3",
			depth: 2,
			outputs: []uci.Output{
				{Depth: 1, Nodes: 112, Score: 0, Mate: 0, PV: []chess.Move{0x13502c106a3}},
				{Depth: 2, Nodes: 915, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc2ab9, 0x13502c018da}},
			},
		},
		{
			name:    "lasker trap with opening book",
			fen:     "rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/5N2/PP2PPPP/RNBQKB1R b KQkq - 1 3",
			depth:   2,
			book:    true,
			outputs: []uci.Output{{PV: []chess.Move{0x1cc2b7e}, Depth: 1, Nodes: 1, Score: 0, Mate: 0}},
		},
		{
			name:  "nodes limit",
			fen:   "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			nodes: 16384,
			depth: 5,
			outputs: []uci.Output{
				{Depth: 1, Nodes: 728, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc38d2}},
				{Depth: 2, Nodes: 4520, Score: 23, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x13e02c328ed}},
				{Depth: 3, Nodes: 19256, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x13e02c328ed, 0x16502c85d26}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()
			_ = engine.Init()
			if tt.book {
				engine.ownBook = true
			}
			engine.table = newHashMapTable()
			engine.pawnTable = noPawnTable{}
			output := engine.Search(context.Background(), unsafeFEN(tt.fen), tt.depth, tt.nodes)
			var outputs []uci.Output
			for o := range output {
				outputs = append(outputs, o)
			}

			assert.Equal(t, tt.outputs, outputs)
		})
	}
}

func TestCachedSearch(t *testing.T) {
	fen := "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	depth := 6

	tests := []struct {
		name   string
		cached bool
		want   []uci.Output
	}{
		{
			"not cached",
			false,
			[]uci.Output{
				{Depth: 1, Nodes: 728, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc38d2}},
				{Depth: 2, Nodes: 4520, Score: 23, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x13e02c328ed}},
				{Depth: 3, Nodes: 19256, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x13e02c328ed, 0x16502c85d26}},
				{Depth: 4, Nodes: 107882, Score: 4, Mate: 0, PV: []chess.Move{0x13d02c25b66, 0x14902c50b76, 0x6401cc15cf, 0x6401cc26ea}},
				{Depth: 5, Nodes: 827378, Score: 3, Mate: 0, PV: []chess.Move{0x13d02c25b66, 0x14902c50b76, 0x6401cc1649, 0x13d02c3455e, 0x14402c4954c}},
				{Depth: 7, Nodes: 3176849, Score: 1, Mate: 0, PV: []chess.Move{0x13d02c25b66, 0x14402c58b74, 0x6401cc38d2, 0x6401cc2e6a, 0x16606c83b63, 0x13f02c30b76, 0x13402c03915}},
			},
		},
		{
			"cached",
			true,
			[]uci.Output{
				{Depth: 1, Nodes: 728, Score: 5, Mate: 0, PV: []chess.Move{0x6401cc38d2}},
				{Depth: 2, Nodes: 4520, Score: 23, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x13e02c328ed}},
				{Depth: 3, Nodes: 19256, Score: 25, Mate: 0, PV: []chess.Move{0x6401cc38d2, 0x13e02c328ed, 0x16502c85d26}},
				{Depth: 4, Nodes: 107882, Score: 4, Mate: 0, PV: []chess.Move{0x13d02c25b66, 0x14902c50b76, 0x6401cc15cf, 0x6401cc26ea}},
				{Depth: 5, Nodes: 826478, Score: 3, Mate: 0, PV: []chess.Move{0x13d02c25b66, 0x14902c50b76, 0x6401cc1649, 0x13d02c3455e, 0x14402c4954c}},
				{Depth: 7, Nodes: 2857556, Score: 1, Mate: 0, PV: []chess.Move{0x13d02c25b66, 0x14402c58b74, 0x6401cc38d2, 0x6401cc2e6a, 0x16606c83b63, 0x13f02c30b76, 0x13402c03915}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()
			_ = engine.Init()
			if !tt.cached {
				engine.table = newHashMapTable()
				engine.pawnTable = noPawnTable{}
			}
			pos := unsafeFEN(fen)
			var outputs []uci.Output
			output := engine.Search(context.Background(), pos, depth, math.MaxInt)
			for o := range output {
				outputs = append(outputs, o)
			}

			require.Equal(t, len(tt.want), len(outputs))
			for i, o := range outputs {
				var wantMoves, gotMoves []string
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
	depth := 5

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
