package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/uci"
)

func TestNew(t *testing.T) {
	e := New()
	assert.Equal(t, 64, e.tableSize)
}

func TestWithTableSize(t *testing.T) {
	e := New(WithTableSize(128))
	assert.Equal(t, 128, e.tableSize)
}

func TestWithOwnBook(t *testing.T) {
	e := New(WithOwnBook(true))
	assert.True(t, e.ownBook)
}

func TestInit(t *testing.T) {
	engine := New()
	assert.Equal(t, noTable{}, engine.table)

	err := engine.Init()

	assert.Nil(t, err)
	assert.IsType(t, &arrayTable{}, engine.table)
}

func TestOptions(t *testing.T) {
	e := New()
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
			e := New()
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
		book    bool
		outputs []uci.Output
	}{
		{
			name:  "mate in 2",
			fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
			depth: 2,
			outputs: []uci.Output{
				{Depth: 1, Nodes: 304, Score: 21, Mate: 0, PV: []chess.Move{0x2c1235c}},
				{Depth: 2, Nodes: 275, Score: mate - 1, Mate: 1, PV: []chess.Move{0x2c1836d}},
			},
		},
		{
			name:  "lasker trap without opening book",
			fen:   "rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/5N2/PP2PPPP/RNBQKB1R b KQkq - 1 3",
			depth: 2,
			outputs: []uci.Output{
				{Depth: 1, Nodes: 29, Score: 38, Mate: 0, PV: []chess.Move{0x2c106a3}},
				{Depth: 2, Nodes: 44, Score: 5, Mate: 0, PV: []chess.Move{0x2c106a3, 0x1cc3481}},
			},
		},
		{
			name:    "lasker trap with opening book",
			fen:     "rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/5N2/PP2PPPP/RNBQKB1R b KQkq - 1 3",
			depth:   2,
			book:    true,
			outputs: []uci.Output{{PV: []chess.Move{0x1cc2b7e}, Depth: 1, Nodes: 1, Score: 1, Mate: 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := New()
			if tt.book {
				engine.ownBook = true
			}
			engine.table = noTable{}
			output := engine.Search(context.Background(), unsafeFEN(tt.fen), tt.depth)
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
	depth := 5

	tests := []struct {
		name   string
		cached bool
		want   []uci.Output
	}{
		{
			"not cached",
			false,
			[]uci.Output{
				{Depth: 1, Nodes: 1637, Score: 15, Mate: 0, PV: []chess.Move{0x1cc58da}},
				{Depth: 2, Nodes: 3627, Score: 0, Mate: 0, PV: []chess.Move{0x1cc17cf, 0x1cc49de}},
				{Depth: 3, Nodes: 103459, Score: 3, Mate: 0, PV: []chess.Move{0x2c25b66, 0x2c50b76, 0x1cc521a}},
				{Depth: 4, Nodes: 170757, Score: 1, Mate: 0, PV: []chess.Move{0x2c25b66, 0x2c50b76, 0x1cc38d2, 0x2c3455e}},
				{Depth: 5, Nodes: 4096282, Score: 0, Mate: 0, PV: []chess.Move{0x2c25b66, 0x2c14362, 0x2c47345, 0x2c58b74, 0x2c03915}},
			},
		},
		{
			"cached",
			true,
			[]uci.Output{
				{Depth: 1, Nodes: 1637, Score: 15, Mate: 0, PV: []chess.Move{0x1cc58da}},
				{Depth: 2, Nodes: 3986, Score: 0, Mate: 0, PV: []chess.Move{0x1cc17cf, 0x1cc49de}},
				{Depth: 3, Nodes: 98708, Score: 3, Mate: 0, PV: []chess.Move{0x2c25b66, 0x2c50b76, 0x1cc521a}},
				{Depth: 4, Nodes: 153294, Score: 1, Mate: 0, PV: []chess.Move{0x2c25b66, 0x2c50b76, 0x1cc38d2, 0x2c3455e}},
				{Depth: 5, Nodes: 3226640, Score: 0, Mate: 0, PV: []chess.Move{0x2c25b66, 0x2c14362, 0x2c47345, 0x2c58b74, 0x2c03915}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := New()
			_ = engine.Init()
			if !tt.cached {
				engine.table = noTable{}
			}
			pos := unsafeFEN(fen)
			var outputs []uci.Output
			output := engine.Search(context.Background(), pos, depth)
			for o := range output {
				outputs = append(outputs, o)
			}
			assert.Equal(t, tt.want, outputs)
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
				engine := New()
				_ = engine.Init()
				if !bb.cached {
					engine.table = noTable{}
				}
				pos := unsafeFEN(fen)
				b.StartTimer()

				output := engine.Search(context.Background(), pos, depth)
				for o := range output {
					_ = o
				}
			}
		})
	}
}
