package chess

import "math/bits"

// PseudoMoves returns the list of pseudo moves.
//
// Some moves may be putting the moving player's king in check and therefore be illegal.
func (pos *Position) PseudoMoves() ([]Move, bool) {
	bbAttackedBy := pos.attackedByBitboard(pos.board.kingSquare(pos.turn))
	switch bits.OnesCount64(uint64(bbAttackedBy)) {
	case 0:
		return pos.pseudoMoves(bbFull, false, false, false), false
	case 1:
		s1 := bbAttackedBy.scanForward()
		s2 := pos.board.kingSquare(pos.turn)
		bbInterference := bbInBetweens[s1][s2] | bbAttackedBy
		return pos.pseudoMoves(bbInterference, false, false, false), true
	default:
		return pos.pseudoMoves(bbFull, false, true, false), true
	}
}

// LoudMoves returns the list of pseudo loud moves.
// Loud moves are moves that capture an opponent piece.
//
// Some moves may be putting the moving player's king in check and therefore be illegal.
func (pos *Position) LoudMoves() []Move {
	return pos.pseudoMoves(bbFull, false, false, true)
}

// pseudoMoves returns the pseudo moves depending on some options.
//
// bbInterference passes the bitboard in which all pieces (except the king) must move to.
// Use it when the king is in check so that pieces can either attack the checking piece or
// interfere in case of an attack by a sliding piece.
//
// allPromos adds all promotion moves. Intending to be used in perft tests.
// Setting it to false only adds queen promotions.
//
// onlyKing returns only the king moves, bypassing all others. Use it when the king is
// in double check.
//
// loud returns only moves that capture enemy pieces. Use it in quiescence search.
func (pos *Position) pseudoMoves(bbInterference bitboard, allPromos, onlyKing, loud bool) []Move {
	size := 50
	if loud {
		size = 20
	}
	moves := make([]Move, 0, size)

	// Setting up variables
	player, opponent := pos.turn, pos.turn.other()
	pawn := WhitePawn
	king := WhiteKing
	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	bbPlayer, bbOpponent := pos.board.bbWhite, pos.board.bbBlack
	upOne, upTwo := north, doubleNorth
	captureR, captureL := northEast, northWest
	if pos.turn == Black {
		pawn = BlackPawn
		king = BlackKing
		bbPlayer, bbOpponent = pos.board.bbBlack, pos.board.bbWhite
		upOne, upTwo = south, doubleSouth
		captureR, captureL = southEast, southWest
	}
	bbPawn := pos.board.bbPawn & bbPlayer

	// King moves
	sqPlayer, sqOpponent := pos.board.kingSquare(player), pos.board.kingSquare(opponent)
	bbKing := bbKingMoves[sqPlayer] & ^bbKingMoves[sqOpponent] & ^bbPlayer
	if loud {
		bbKing &= bbOpponent
	}
	for bbs2 := bbKing; bbs2 > 0; bbs2 = bbs2.resetLSB() {
		s2, p2 := bbs2.scanForward(), NoPiece
		if s2.bitboard()&bbOpponent > 0 {
			p2 = pos.board.pieceByColor(s2, opponent)
		}
		moves = append(moves, newMove(king, p2, sqPlayer, s2, NoSquare, NoPiece))
	}

	if onlyKing {
		return moves
	}

	// Castles
	if !loud {
		if data := castles[2*uint8(player)+uint8(kingSide)]; pos.castlingRights.canCastle(player, kingSide) && bbOccupancy&data.bbTravel == 0 {
			moves = append(moves, newCastleMove(king, data.s1, data.s2, kingSide))
		}
		if data := castles[2*uint8(player)+uint8(queenSide)]; pos.castlingRights.canCastle(player, queenSide) && bbOccupancy&data.bbTravel == 0 {
			moves = append(moves, newCastleMove(king, data.s1, data.s2, queenSide))
		}

	}

	// Pawn moves
	if !loud {
		bbUpOne, bbUpTwo := pawnMoveBitboard(bbPawn, bbOccupancy, player)
		for _, dest := range [2]bbDir{
			{bbUpOne & bbInterference, upOne},
			{bbUpTwo & bbInterference, upTwo},
		} {
			for ; dest.bb > 0; dest.bb = dest.bb.resetLSB() {
				s2 := dest.bb.scanForward()
				s1 := s2 - Square(dest.dir)

				if s2.bitboard()&(bbRank1^bbRank8) == 0 {
					moves = append(moves, newPawnMove(pawn, NoPiece, s1, s2, NoSquare, NoPiece))
					continue
				}

				if allPromos {
					moves = append(moves,
						newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Queen.color(player)),
						newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Rook.color(player)),
						newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Bishop.color(player)),
						newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Knight.color(player)),
					)
					continue
				}

				moves = append(moves, newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Queen.color(player)))
			}
		}
	}

	// Pawn captures
	bbCaptureR, bbCaptureL := pawnCaptureBitboard(bbPawn, player)
	bbEnPassant := pos.enPassant.bitboard()

	bbPawnInterference := bbInterference
	if pos.enPassant != NoSquare &&
		(bbInterference&pos.board.bbPawn).scanForward()+Square(upOne) == pos.enPassant {
		bbPawnInterference |= bbEnPassant
	}

	for _, dest := range [2]bbDir{
		{bbCaptureR & (bbOpponent | bbEnPassant) & bbPawnInterference, captureR},
		{bbCaptureL & (bbOpponent | bbEnPassant) & bbPawnInterference, captureL},
	} {
		for ; dest.bb > 0; dest.bb = dest.bb.resetLSB() {
			s2 := dest.bb.scanForward()
			s1 := s2 - Square(dest.dir)
			p2 := pos.board.pieceByColor(s2, opponent)

			if s2.bitboard()&(bbRank1^bbRank8) == 0 {
				moves = append(moves, newPawnMove(pawn, p2, s1, s2, pos.enPassant, NoPiece))
				continue
			}

			if allPromos {
				moves = append(moves,
					newPawnMove(pawn, p2, s1, s2, NoSquare, Queen.color(player)),
					newPawnMove(pawn, p2, s1, s2, NoSquare, Rook.color(player)),
					newPawnMove(pawn, p2, s1, s2, NoSquare, Bishop.color(player)),
					newPawnMove(pawn, p2, s1, s2, NoSquare, Knight.color(player)),
				)
				continue
			}

			moves = append(moves, newPawnMove(pawn, p2, s1, s2, NoSquare, Queen.color(player)))
		}
	}

	// Other pieces
	for pt := Knight; pt <= Queen; pt += 2 {
		p1 := pt.color(player)
		for bbs1 := pos.board.getBitboard(pt, player); bbs1 > 0; bbs1 = bbs1.resetLSB() {
			s1 := bbs1.scanForward()
			bbs2 := pieceBitboard(s1, pt, bbOccupancy) & ^bbPlayer & bbInterference
			if loud {
				bbs2 &= bbOpponent
			}
			for ; bbs2 > 0; bbs2 = bbs2.resetLSB() {
				s2, p2 := bbs2.scanForward(), NoPiece
				if s2.bitboard()&bbOpponent > 0 {
					p2 = pos.board.pieceByColor(s2, opponent)
				}
				moves = append(moves, newPieceMove(p1, p2, s1, s2))
			}
		}
	}

	return moves
}

// isDiscoveredCheck checks whether the moving piece uncovers
// a check given by an enemy piece.
func (pos *Position) isDiscoveredCheck(m Move) bool {
	s1, s2 := m.S1(), m.S2()
	bb := bbFiles[s1] | bbRanks[s1] | bbDiagonals[s1] | bbAntiDiagonals[s1]
	if bb&pos.board.bbKing&pos.board.getColor(pos.turn) == 0 {
		return false
	}

	kingSq := pos.board.kingSquare(pos.turn)
	bbCaptured := s2.bitboard()
	bbOpponent := pos.board.getColor(pos.turn.other()) & ^bbCaptured

	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	bbOccupancy &= ^s1.bitboard()
	bbOccupancy |= bbCaptured
	if m.HasTag(EnPassant) {
		if pos.turn == White {
			bbOccupancy &= ^(bbCaptured >> 8)
		} else {
			bbOccupancy &= ^(bbCaptured << 8)
		}
	}
	bbRookMoves := bbMagicRookMoves[rookMagics[kingSq].index(bbOccupancy)]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[kingSq].index(bbOccupancy)]

	return (pos.board.bbQueen^pos.board.bbRook)&bbRookMoves&bbOpponent > 0 ||
		(pos.board.bbQueen^pos.board.bbBishop)&bbBishopMoves&bbOpponent > 0
}

// isSquareAttacked checks whether the square is attacked by
// an enemy piece.
func (pos *Position) isSquareAttacked(sq Square) bool {
	return pos.attackedByBitboard(sq) > 0
}

// attackedByBitboard returns the bitboard of the pieces that attack teh square.
func (pos *Position) attackedByBitboard(sq Square) bitboard {
	bbOpponent := pos.board.getColor(pos.turn.other())
	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	bbRookMoves := bbMagicRookMoves[rookMagics[sq].index(bbOccupancy)]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[sq].index(bbOccupancy)]

	var bb bitboard
	bb |= singlePawnCaptureBitboard(sq, pos.turn) & pos.board.bbPawn
	bb |= bbKingMoves[sq] & pos.board.bbKing
	bb |= bbKnightMoves[sq] & pos.board.bbKnight
	bb |= (pos.board.bbQueen | pos.board.bbRook) & bbRookMoves
	bb |= (pos.board.bbQueen | pos.board.bbBishop) & bbBishopMoves
	return bb & bbOpponent
}

// isCastleLegal checks whether the castle move is legal.
//
// Assumes that the castle rights have already been checked and
// that the king's travel path is clear.
//
// Checks that the king does not leave, cross over, or finish on
// s square attacked by an enemy piece.
func (pos *Position) isCastleLegal(m Move) bool {
	s := kingSide
	if m.HasTag(QueenSideCastle) {
		s = queenSide
	}
	bbOpponent := pos.board.getColor(pos.turn.other())
	cc := castleChecks[2*uint8(pos.turn)+uint8(s)]

	if cc.bbPawn&pos.board.bbPawn&bbOpponent > 0 ||
		cc.bbKnight&pos.board.bbKnight&bbOpponent > 0 ||
		cc.bbKing&pos.board.bbKing&bbOpponent > 0 {
		return false
	}

	var bbBishopAttacks, bbRookAttacks bitboard
	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	for _, sq := range cc.squares {
		index := bishopMagics[sq].index(bbOccupancy)
		bbBishopAttacks |= bbMagicBishopMoves[index]
	}

	if bb := pos.board.bbBishop | pos.board.bbQueen; bbBishopAttacks&bbOpponent&bb > 0 {
		return false
	}

	for _, sq := range cc.squares {
		index := rookMagics[sq].index(bbOccupancy)
		bbRookAttacks |= bbMagicRookMoves[index]
	}

	return bbRookAttacks&(pos.board.bbRook|pos.board.bbQueen)&bbOpponent == 0
}

// bbDir associates a bitboard with a direction.
type bbDir struct {
	bb  bitboard
	dir direction
}

// castles contains the castles' data.
//
// indexed by `2*Color+side`.
var castles = [4]castleData{
	{1<<F8 | 1<<G8, E8, G8},         // black, king side
	{1<<B8 | 1<<C8 | 1<<D8, E8, C8}, // black, queen side
	{1<<F1 | 1<<G1, E1, G1},         // white, king side
	{1<<B1 | 1<<C1 | 1<<D1, E1, C1}, // white, queen side
}

// casteData represents a castle's data.
type castleData struct {
	bbTravel bitboard // bitboard traveled by the king
	s1       Square   // king s1
	s2       Square   // king s2
}

// castleCheck represents a castle's check.
type castleCheck struct {
	bbPawn   bitboard
	bbKnight bitboard
	bbKing   bitboard
	squares  [3]Square
}

// castleChecks contains the castle checks.
//
// indexed by `2*Color+side`.
var castleChecks [4]castleCheck

// initializes castleChecks.
//
// requires bbKingMoves, bbKnightMoves.
func initCastleChecks() {
	castles := [4]struct {
		color   Color
		squares [3]Square
	}{
		{Black, [3]Square{E8, F8, G8}},
		{Black, [3]Square{C8, D8, E8}},
		{White, [3]Square{E1, F1, G1}},
		{White, [3]Square{C1, D1, E1}},
	}

	for i, castle := range castles {
		var bbPawn, bbKnight, bbKing bitboard
		for _, sq := range castle.squares {
			bbPawn |= singlePawnCaptureBitboard(sq, castle.color)
			bbKnight |= bbKnightMoves[sq]
			bbKing |= bbKingMoves[sq]
		}

		castleChecks[i] = castleCheck{bbPawn, bbKnight, bbKing, castle.squares}
	}
}
