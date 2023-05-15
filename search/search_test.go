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

func TestInit(t *testing.T) {
	engine := New()
	assert.Equal(t, noTable{}, engine.table)

	err := engine.Init()

	assert.Nil(t, err)
	assert.IsType(t, &ristrettoTable{}, engine.table)
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
		name  string
		fen   string
		depth int
		oo    []uci.Output
	}{
		{
			name:  "mate in 2",
			fen:   "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
			depth: 2,
			oo: []uci.Output{
				{Depth: 1, Nodes: 46, Score: 357, Mate: 0, PV: []chess.Move{0x2c322dc}},
				{Depth: 2, Nodes: 58, Score: 9223372036854775806, Mate: 1, PV: []chess.Move{0x2c1836d}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := New()
			output := engine.Search(context.Background(), unsafeFEN(tt.fen), tt.depth)
			var outputs []uci.Output
			for o := range output {
				outputs = append(outputs, *o)
			}

			assert.Equal(t, tt.oo, outputs)
		})
	}
}
