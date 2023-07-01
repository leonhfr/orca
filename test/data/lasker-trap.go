// Package testdata exports test polyglot data.
package testdata

import _ "embed"

//nolint:revive
//go:embed lasker-trap.bin
var LaskerTrap []byte
