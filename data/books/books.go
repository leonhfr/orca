// Package books exports embedded data from opening books.
package books

import _ "embed"

// Performance opening book in Polyglot format.
// Made by Marc Lacrosse.
// Source: http://wbec-ridderkerk.nl/html/download.htm
//
// Positions: 92.954
//
// Depth: 20
//
// Saturation (percentage of analyzed positions): 95.6%
//
//go:embed performance.bin
var Performance []byte
