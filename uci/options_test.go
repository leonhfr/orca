package uci

import (
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
)

// compile time check that booleanUCIOption implements uciOption.
var _ uciOption = booleanUCIOption{}

// compile time check that integerSearchOption implements searchOption.
var _ searchOption = integerSearchOption{}

// compile time check that booleanSearchOption implements searchOption.
var _ searchOption = booleanSearchOption{}

func TestOptionBooleanString(t *testing.T) {
	assert.Equal(t, chess960Option.name, chess960Option.String())
}

func TestOptionBooleanUCI(t *testing.T) {
	assert.Equal(t, responseOption{
		Type:    booleanOptionType,
		Name:    chess960Option.name,
		Default: strconv.FormatBool(chess960Option.def),
	}, chess960Option.response())
}

// optionBoolean.defaultFunc tested in New

func TestOptionBooleanOptionFunc(t *testing.T) {
	type want struct {
		notation     chess.Notation
		moveNotation chess.MoveNotation
		err          string
	}

	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "value cannot be parsed as bool",
			args: "foobar",
			want: want{chess.FEN{}, chess.UCI{}, "strconv.ParseBool: parsing \"foobar\": invalid syntax"},
		},
		{
			name: "true",
			args: "true",
			want: want{chess.ShredderFEN{}, chess.UCIChess960{}, ""},
		},
		{
			name: "false",
			args: "false",
			want: want{chess.FEN{}, chess.UCI{}, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := chess960Option.optionFunc(tt.args)
			if err != nil {
				assert.Equal(t, tt.want.err, err.Error())
			}

			c := NewController("", "", io.Discard)
			fn(c)
			assert.Equal(t, tt.want.notation, c.notation)
			assert.Equal(t, tt.want.moveNotation, c.moveNotation)
		})
	}
}
