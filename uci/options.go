package uci

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/leonhfr/orca/search"
)

// optionType represent an option type.
type optionType uint8

const (
	integerOptionType optionType = iota // OptionInteger represents an integer option.
	booleanOptionType                   // OptionBoolean represents a boolean option.
)

var (
	// availableUCIOptions holds all the uci available options.
	availableUCIOptions = []uciOption{chess960Option}

	// availableSearchOptions holds all the search available options.
	availableSearchOptions = []searchOption{tableSizeOption, ownBookOption}

	// chess960Option represents the chess mode, classic or Chess960.
	chess960Option = booleanUCIOption{
		name: "UCI_Chess960",
		def:  false,
		fn:   WithChess960,
	}

	// tableSizeOption represents the size of the transposition hash table.
	tableSizeOption = integerSearchOption{
		name: "Hash",
		def:  64,
		min:  1,
		max:  16 * 1024,
		fn:   search.WithTableSize,
	}

	// ownBook represents whether the search engine should use its own opening book.
	ownBookOption = booleanSearchOption{
		name: "OwnBook",
		def:  false,
		fn:   search.WithOwnBook,
	}

	errOptionName   = errors.New("option name not found")
	errOutsideBound = errors.New("option value outside bounds")
)

// option is the interface implemented by all options.
type option interface {
	fmt.Stringer
	response() responseOption
}

// uciOption is the interface implemented by each option that modifies the Controller.
type uciOption interface {
	option
	defaultFunc() func(*Controller)
	optionFunc(value string) (func(*Controller), error)
}

// booleanUCIOption represents a boolean option.
//
//nolint:govet
type booleanUCIOption struct {
	name string
	def  bool
	fn   func(bool) Option
}

// String implements the uciOption interface.
func (o booleanUCIOption) String() string {
	return o.name
}

// response implements the uciOption interface.
func (o booleanUCIOption) response() responseOption {
	return responseOption{
		Type:    booleanOptionType,
		Name:    o.name,
		Default: strconv.FormatBool(o.def),
	}
}

// defaultFunc implements the uciOption interface.
func (o booleanUCIOption) defaultFunc() func(*Controller) {
	return o.fn(o.def)
}

// optionFunc implements the uciOption interface.
func (o booleanUCIOption) optionFunc(value string) (func(*Controller), error) {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return func(_ *Controller) {}, err
	}

	return o.fn(v), nil
}

// searchOption is the interface implemented by each option that modifies search.Engine.
type searchOption interface {
	option
	defaultFunc() func(*search.Engine)
	optionFunc(value string) (func(*search.Engine), error)
}

// integerSearchOption represents an integer option.
//
//nolint:govet
type integerSearchOption struct {
	name          string
	def, min, max int
	fn            func(int) search.Option
}

// String implements the searchOption interface.
func (o integerSearchOption) String() string {
	return o.name
}

// response implements the searchOption interface.
func (o integerSearchOption) response() responseOption {
	return responseOption{
		Type:    integerOptionType,
		Name:    o.name,
		Default: strconv.Itoa(o.def),
		Min:     strconv.Itoa(o.min),
		Max:     strconv.Itoa(o.max),
	}
}

// defaultFunc implements the searchOption interface.
func (o integerSearchOption) defaultFunc() func(*search.Engine) {
	return o.fn(o.def)
}

// optionFunc implements the searchOption interface.
func (o integerSearchOption) optionFunc(value string) (func(*search.Engine), error) {
	v, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return func(_ *search.Engine) {}, err
	}

	if int(v) < o.min || int(v) > o.max {
		return func(_ *search.Engine) {}, errOutsideBound
	}

	return o.fn(int(v)), nil
}

// booleanSearchOption represents a boolean option.
//
//nolint:govet
type booleanSearchOption struct {
	name string
	def  bool
	fn   func(bool) search.Option
}

// String implements the searchOption interface.
func (o booleanSearchOption) String() string {
	return o.name
}

// response implements the searchOption interface.
func (o booleanSearchOption) response() responseOption {
	return responseOption{
		Type:    booleanOptionType,
		Name:    o.name,
		Default: strconv.FormatBool(o.def),
	}
}

// defaultFunc implements the searchOption interface.
func (o booleanSearchOption) defaultFunc() func(*search.Engine) {
	return o.fn(o.def)
}

// optionFunc implements the searchOption interface.
func (o booleanSearchOption) optionFunc(value string) (func(*search.Engine), error) {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return func(_ *search.Engine) {}, err
	}

	return o.fn(v), nil
}
