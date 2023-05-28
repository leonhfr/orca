#!/bin/bash
### perft - run perft test
###
### Usage:
###   $ ./perft.sh <depth> <fen> [moves]
###
### With perftree:
###   $ perftree ./perft.sh
###
### Requires:
###  - perftree (/github.com/agausmann/perftree)
###  - stockfish
###
### Options:
###   <depth>   Maximum depth of the evaluation.
###   <fen>     FEN string of the base position.
###   [moves]   Optional space-separated list of moves
###             from the base position to the position
###             to be evaluated.
###   -h        Show this message.

if [[ $# = 0 ]] || [[ "$1" = "-h" ]]; then
  sed -rn 's/^### ?//p' "$0"
  exit 0
fi

go run . "$1" "$2" "$3"
