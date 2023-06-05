package chess

// PawnProperty represents different pawn properties.
type PawnProperty uint8

// NoProperty represents the absence of properties.
const NoProperty PawnProperty = 0

const (
	// Doubled represents a doubled pawn.
	// A pawn is doubled if there are more pawns of the same color on a given file.
	Doubled PawnProperty = 1 << iota
	// Isolani represents an isolated pawn.
	// An isolani has no friendly pawns on adjacent files.
	Isolani
	// HalfIsolani represents a pawn isolated on only one adjacent file.
	HalfIsolani
	// Passed represents a passed pawn.
	Passed
)

// PawnMap executes the callback for each pawn on the board.
// Intended to be used in evaluation functions.
func (pos *Position) PawnMap(cb func(p Piece, sq Square, properties PawnProperty)) {
	for c := Black; c <= White; c++ {
		bbPlayerPawn := pos.board.bbColors[c] & pos.board.bbPieces[Pawn]
		bbOpponentPawn := pos.board.bbColors[c.other()] & pos.board.bbPieces[Pawn]
		bbOpponentFrontSpans := bbOpponentPawn.frontSpans(c.other())

		bbPawnsBehindOwn := bbPlayerPawn & bbPlayerPawn.rearSpans(c)
		bbPawnsInFrontOwn := bbPlayerPawn & bbPlayerPawn.frontSpans(c)
		bbDoubled := bbPawnsBehindOwn | bbPawnsInFrontOwn

		bbNoNeighborOnEastFile := bbPlayerPawn.noNeighborOnEastFile()
		bbNoNeighborOnWestFile := bbPlayerPawn.noNeighborOnWestFile()

		bbIsolanis := bbNoNeighborOnEastFile & bbNoNeighborOnWestFile
		bbHalfIsolanis := bbNoNeighborOnEastFile ^ bbNoNeighborOnWestFile

		bbPassed := passedPawns(bbPlayerPawn, bbOpponentFrontSpans)

		for pawn := Pawn.color(c); bbPlayerPawn > 0; bbPlayerPawn = bbPlayerPawn.resetLSB() {
			sq := bbPlayerPawn.scanForward()
			bb := sq.bitboard()
			var properties PawnProperty

			if bb&bbDoubled > 0 {
				properties ^= Doubled
			}

			if bb&bbIsolanis > 0 {
				properties ^= Isolani
			} else if bb&bbHalfIsolanis > 0 {
				properties ^= HalfIsolani
			}

			if bb&bbPassed > 0 {
				properties ^= Passed
			}

			cb(pawn, sq, properties)
		}
	}
}

func passedPawns(bbPlayerPawn, bbOpponentFrontSpans bitboard) bitboard {
	bbAllFrontSpans := bbOpponentFrontSpans
	bbAllFrontSpans |= bbAllFrontSpans.eastOne() | bbAllFrontSpans.westOne()
	return bbPlayerPawn & ^bbAllFrontSpans
}

func (b bitboard) noNeighborOnEastFile() bitboard {
	return b & ^b.westAttackFileFill()
}

func (b bitboard) noNeighborOnWestFile() bitboard {
	return b & ^b.eastAttackFileFill()
}

// frontSpans computes the front spans.
// Front spans are front fill shifted one step further in the fill direction.
func (b bitboard) frontSpans(c Color) bitboard {
	if c == Black {
		return b.southFill().southOne()
	}
	return b.northFill().northOne()
}

// rearSpans computes the rear spans.
// Rear spans are rear fill shifted one step further in the fill direction.
func (b bitboard) rearSpans(c Color) bitboard {
	if c == Black {
		return b.northFill().northOne()
	}
	return b.southFill().southOne()
}

// eastAttackFileFill computes the east attack file fill.
func (b bitboard) eastAttackFileFill() bitboard {
	return b.fileFill().eastOne()
}

// westAttackFileFill computes the west attack file fill.
func (b bitboard) westAttackFileFill() bitboard {
	return b.fileFill().westOne()
}

// fileFill computes the file fill.
func (b bitboard) fileFill() bitboard {
	return b.northFill() | b.southFill()
}

// northFill computes the north fill.
func (b bitboard) northFill() bitboard {
	b |= b << 8
	b |= b << 16
	b |= b << 32
	return b
}

// southFill computes the south fill.
func (b bitboard) southFill() bitboard {
	b |= b >> 8
	b |= b >> 16
	b |= b >> 32
	return b
}

// northOne shifts the bitboard toward the north.
func (b bitboard) northOne() bitboard {
	return b << 8
}

// southOne shifts the bitboard toward the south.
func (b bitboard) southOne() bitboard {
	return b >> 8
}

// eastOne shifts the bitboard toward the east.
func (b bitboard) eastOne() bitboard {
	return b << 1 & ^bbFileA
}

// westOne shifts the bitboard toward the west.
func (b bitboard) westOne() bitboard {
	return b >> 1 & ^bbFileH
}
