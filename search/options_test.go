package search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/uci"
)

func TestOptionIntegerString(t *testing.T) {
	assert.Equal(t, tableSizeOption.name, tableSizeOption.String())
}

func TestOptionIntegerUCI(t *testing.T) {
	assert.Equal(t, uci.Option{
		Type:    uci.OptionInteger,
		Name:    tableSizeOption.name,
		Default: fmt.Sprint(tableSizeOption.def),
		Min:     fmt.Sprint(tableSizeOption.min),
		Max:     fmt.Sprint(tableSizeOption.max),
	}, tableSizeOption.uci())
}

// optionInteger.defaultFunc tested in New

func TestOptionIntegerOptionFunc(t *testing.T) {
	type want struct {
		value int
		err   string
	}

	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "value cannot be parsed as integer",
			args: "foobar",
			want: want{64, "strconv.ParseInt: parsing \"foobar\": invalid syntax"},
		},
		{
			name: "value is outside bounds",
			args: "0",
			want: want{64, errOutsideBound.Error()},
		},
		{
			name: "value is valid",
			args: "256",
			want: want{256, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := tableSizeOption.optionFunc(tt.args)
			if err != nil {
				assert.Equal(t, tt.want.err, err.Error())
			}

			e := New()
			fn(e)
			assert.Equal(t, tt.want.value, e.tableSize)
		})
	}
}
