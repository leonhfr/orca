package chess

// PawnCount returns the pawn count, indexed by color.
func (pos *Position) PawnCount() [2]int {
	return [2]int{
		(pos.board.bbColors[Black] & pos.board.bbPieces[Pawn]).ones(),
		(pos.board.bbColors[White] & pos.board.bbPieces[Pawn]).ones(),
	}
}

// PawnProperty represents different pawn properties.
type PawnProperty uint8

// NoPawnProperty represents the absence of properties.
const NoPawnProperty PawnProperty = 0

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
	// A passed pawn is a pawn with no opponent pawns in front on the same or adjacent files.
	Passed
	// Backward represents a backward pawn.
	// A pawn is backward if its stop square is not in its own front attack spans
	// but is controlled by an enemy sentry.
	Backward
)

// HasProperty checks the presence of the given property.
func (pp PawnProperty) HasProperty(p PawnProperty) bool {
	return pp&PawnProperty(p) > 0
}

// PawnMap executes the callback for each pawn on the board.
//
// Intended to be used in evaluation functions.
func (pos *Position) PawnMap(cb func(p Piece, sq Square, properties PawnProperty)) {
	for c := Black; c <= White; c++ {
		bbPlayerPawn := pos.board.bbColors[c] & pos.board.bbPieces[Pawn]
		bbOpponentPawn := pos.board.bbColors[c.Other()] & pos.board.bbPieces[Pawn]

		bbDoubled := doubledPawns(c, bbPlayerPawn)
		bbIsolanis, bbHalfIsolanis := isolatedPawns(bbPlayerPawn)
		bbPassed := passedPawns(c, bbPlayerPawn, bbOpponentPawn)
		bbBackward := backwardPawns(c, bbPlayerPawn, bbOpponentPawn)

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

			if bb&bbBackward > 0 {
				properties ^= Backward
			}

			cb(pawn, sq, properties)
		}
	}
}

func doubledPawns(c Color, bbPlayerPawn bitboard) bitboard {
	bbPawnsBehindOwn := bbPlayerPawn & bbPlayerPawn.rearSpans(c)
	bbPawnsInFrontOwn := bbPlayerPawn & bbPlayerPawn.frontSpans(c)
	return bbPawnsBehindOwn | bbPawnsInFrontOwn
}

func isolatedPawns(bbPlayerPawn bitboard) (bitboard, bitboard) {
	bbNoNeighborOnEastFile := bbPlayerPawn.noNeighborOnEastFile()
	bbNoNeighborOnWestFile := bbPlayerPawn.noNeighborOnWestFile()
	bbIsolanis := bbNoNeighborOnEastFile & bbNoNeighborOnWestFile
	bbHalfIsolanis := bbNoNeighborOnEastFile ^ bbNoNeighborOnWestFile
	return bbIsolanis, bbHalfIsolanis
}

func passedPawns(c Color, bbPlayerPawn, bbOpponentPawn bitboard) bitboard {
	bbAllFrontSpans := bbOpponentPawn.frontSpans(c.Other())
	bbAllFrontSpans |= bbAllFrontSpans.eastOne() | bbAllFrontSpans.westOne()
	return bbPlayerPawn & ^bbAllFrontSpans
}

func backwardPawns(c Color, bbPlayerPawn, bbOpponentPawns bitboard) bitboard {
	bbStops := bbPlayerPawn.stops(c)
	bbPlayerAttackSpans := bbPlayerPawn.frontSpans(c).eastAttackFileFill() |
		bbPlayerPawn.frontSpans(c).westAttackFileFill()
	bbOpponentAttacks := bbOpponentPawns.eastAttack(c.Other()) |
		bbOpponentPawns.westAttack(c.Other())
	bbIntersection := bbStops & bbOpponentAttacks & ^bbPlayerAttackSpans

	if c == Black {
		return bbIntersection.northOne()
	}
	return bbIntersection.southOne()
}

func (b bitboard) noNeighborOnEastFile() bitboard {
	return b & ^b.westAttackFileFill()
}

func (b bitboard) noNeighborOnWestFile() bitboard {
	return b & ^b.eastAttackFileFill()
}

// stops returns the stops bitboard.
func (b bitboard) stops(c Color) bitboard {
	if c == Black {
		return b.southOne()
	}
	return b.northOne()
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

// eastAttack returns the east attacks.
func (b bitboard) eastAttack(c Color) bitboard {
	if c == Black {
		return b.southEastOne()
	}
	return b.northEastOne()
}

// westAttack returns the west attacks.
func (b bitboard) westAttack(c Color) bitboard {
	if c == Black {
		return b.southWestOne()
	}
	return b.northWestOne()
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
