package search

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/leonhfr/orca/uci"
)

var (
	// availableOptions holds all the search engine available options.
	availableOptions = []option{tableSizeOption, ownBookOption}

	// tableSizeOption represents the size of the transposition hash table.
	tableSizeOption = optionInteger{
		name: "Hash",
		def:  64,
		min:  1,
		max:  16 * 1024,
		fn:   WithTableSize,
	}

	// ownBook represents whether the search engine should use its own opening book.
	ownBookOption = optionBoolean{
		name: "OwnBook",
		def:  false,
		fn:   WithOwnBook,
	}

	errOptionName   = errors.New("option name not found")
	errOutsideBound = errors.New("option value outside bounds")
)

// option is the interface implemented by each option type.
type option interface {
	fmt.Stringer
	uci() uci.Option
	defaultFunc() func(*Engine)
	optionFunc(value string) (func(*Engine), error)
}

// optionInteger represents an integer option.
//
//nolint:govet
type optionInteger struct {
	name          string
	def, min, max int
	fn            func(int) func(*Engine)
}

// String implements the option interface.
func (o optionInteger) String() string {
	return o.name
}

// uci implements the option interface.
func (o optionInteger) uci() uci.Option {
	return uci.Option{
		Type:    uci.OptionInteger,
		Name:    o.name,
		Default: fmt.Sprint(o.def),
		Min:     fmt.Sprint(o.min),
		Max:     fmt.Sprint(o.max),
	}
}

// defaultFunc implements the option interface.
func (o optionInteger) defaultFunc() func(*Engine) {
	return o.fn(o.def)
}

// optionFunc implements the option interface.
func (o optionInteger) optionFunc(value string) (func(*Engine), error) {
	v, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return func(e *Engine) {}, err
	}

	if int(v) < o.min || int(v) > o.max {
		return func(e *Engine) {}, errOutsideBound
	}

	return o.fn(int(v)), nil
}

// optionBoolean represents a boolean option.
//
//nolint:govet
type optionBoolean struct {
	name string
	def  bool
	fn   func(bool) func(*Engine)
}

// String implements the option interface.
func (o optionBoolean) String() string {
	return o.name
}

// uci implements the option interface.
func (o optionBoolean) uci() uci.Option {
	return uci.Option{
		Type:    uci.OptionBoolean,
		Name:    o.name,
		Default: fmt.Sprint(o.def),
	}
}

// defaultFunc implements the option interface.
func (o optionBoolean) defaultFunc() func(*Engine) {
	return o.fn(o.def)
}

// optionFunc implements the option interface.
func (o optionBoolean) optionFunc(value string) (func(*Engine), error) {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return func(e *Engine) {}, err
	}

	return o.fn(v), nil
}
