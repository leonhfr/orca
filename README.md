# orca

[![Go Reference](https://pkg.go.dev/badge/github.com/leonhfr/orca.svg)](https://pkg.go.dev/github.com/leonhfr/orca)

Orca is a UCI-compliant chess engine written in Go that analyzes chess positions and computes the optimal moves.

## Installation

Several installation methods are available:

- using the `go` toolchain:

```sh
go install github.com/leonhfr/orca@latest
```

- compile from source (requires `go@1.19` and `make`):

```sh
git clone git@github.com:leonhfr/orca.git
cd orca
make build
```

## Quick start

Orca is not a complete chess software and requires a [UCI-compatible](https://backscattering.de/chess/uci/) graphical user interface (GUI) to be used comfortably. GUI options include [SCID](http://scid.sourceforge.net/), [CuteChess](https://github.com/cutechess/cutechess), [Arena](http://www.playwitharena.de/) and [Shredder](https://www.shredderchess.com/).

In the future, Orca will be available as a Lichess bot.

## Options

```
option name Hash type spin default 64 min 1 max 16384
option name OwnBook type check default false
```

Available options are:
- `Hash`: size in MB used for the transposition table
- `OwnBook`: allow the engine to use its own opening book
