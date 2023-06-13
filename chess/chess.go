// Package chess provides types and functions to handle chess positions.
package chess

// CheckData contains check data.
type CheckData bitboard

// InCheck returns check data and whether the king is in check.
func (pos *Position) InCheck() (CheckData, bool) {
	bbAttackedBy := pos.attackedByBitboard(pos.board.kingSquare(pos.turn), pos.turn)
	return CheckData(bbAttackedBy), bbAttackedBy > 0
}

// HasInsufficientMaterial returns true if there is insufficient material to achieve a mate.
//
// Combinations include:
//
//	king versus king
//	king and bishop versus king
//	king and knight versus king
//	king and bishop versus king and bishop with the bishops on the same color
func (pos *Position) HasInsufficientMaterial() bool {
	bbOccupancy := pos.board.bbColors[White] ^ pos.board.bbColors[Black]
	pieces := bbOccupancy.ones()
	if pieces > 4 {
		return false
	}

	knights := pos.board.bbPieces[Knight].ones()
	bishops := pos.board.bbPieces[Bishop].ones()

	if pieces == 2 || pieces == 3 && (knights == 1 || bishops == 1) {
		return true
	}

	if bbBlack := pos.board.bbColors[Black] & pos.board.bbPieces[Bishop]; pieces == 4 && bishops == 2 && bbBlack.ones() == 1 {
		bbWhite := pos.board.bbColors[White] & pos.board.bbPieces[Bishop]
		sqBlack := bbBlack.scanForward()
		sqWhite := bbWhite.scanForward()
		return sqBlack.sameColor(sqWhite)
	}

	return false
}

// PseudoMoves returns the list of pseudo moves.
//
// Some moves may be putting the moving player's king in check and therefore be illegal.
func (pos *Position) PseudoMoves(data CheckData) []Move {
	bbAttackedBy := bitboard(data)
	switch count := bbAttackedBy.ones(); {
	case count > 1:
		return pos.pseudoMoves(bbFull, true, false)
	case count == 1:
		s1 := bbAttackedBy.scanForward()
		s2 := pos.board.kingSquare(pos.turn)
		bbInterference := bbInBetweens[s1][s2] | bbAttackedBy
		return pos.pseudoMoves(bbInterference, false, false)
	default:
		return pos.pseudoMoves(bbFull, false, false)
	}
}

// LoudMoves returns the list of pseudo loud moves.
// Loud moves are moves that capture an opponent piece.
//
// Some moves may be putting the moving player's king in check and therefore be illegal.
func (pos *Position) LoudMoves() []Move {
	return pos.pseudoMoves(bbFull, false, true)
}

// pseudoMoves returns the pseudo moves depending on some options.
//
// bbInterference passes the bitboard in which all pieces (except the king) must move to.
// Use it when the king is in check so that pieces can either attack the checking piece or
// interfere in case of an attack by a sliding piece.
//
// onlyKing returns only the king moves, bypassing all others. Use it when the king is
// in double check.
//
// loud returns only moves that capture enemy pieces. Use it in quiescence search.
func (pos *Position) pseudoMoves(bbInterference bitboard, onlyKing, loud bool) []Move {
	size := 50
	if loud {
		size = 20
	}
	moves := make([]Move, 0, size)

	// Setting up variables
	player, opponent := pos.turn, pos.turn.Other()
	pawn := WhitePawn
	king := WhiteKing
	bbOccupancy := pos.board.bbColors[White] ^ pos.board.bbColors[Black]
	bbPlayer, bbOpponent := pos.board.bbColors[White], pos.board.bbColors[Black]
	upOne, upTwo := north, doubleNorth
	captureR, captureL := northEast, northWest
	if pos.turn == Black {
		pawn = BlackPawn
		king = BlackKing
		bbPlayer, bbOpponent = bbOpponent, bbPlayer
		upOne, upTwo = south, doubleSouth
		captureR, captureL = southEast, southWest
	}
	bbPawn := pos.board.bbPieces[Pawn] & bbPlayer

	// King moves
	sqKing, sqEnemyKing := pos.board.kingSquare(player), pos.board.kingSquare(opponent)
	bbKing := bbKingMoves[sqKing] & ^bbKingMoves[sqEnemyKing] & ^bbPlayer
	if loud {
		bbKing &= bbOpponent
	}
	for bbs2 := bbKing; bbs2 > 0; bbs2 = bbs2.resetLSB() {
		s2, p2 := bbs2.scanForward(), NoPiece
		if s2.bitboard()&bbOpponent > 0 {
			p2 = pos.board.pieceByColor(s2, opponent)
		}
		moves = append(moves, newPieceMove(king, p2, sqKing, s2, false))
	}

	if onlyKing {
		return moves
	}

	bbChecks := pos.attackBitboards(sqEnemyKing, opponent)

	// Castles
	if !loud {
		if data := castles[2*uint8(player)+uint8(kingSide)]; pos.castlingRights.canCastle(player, kingSide) && bbOccupancy&data.bbTravel == 0 {
			check := data.rook2.bitboard()&bbChecks[Rook] > 0
			moves = append(moves, newCastleMove(king, data.king1, data.king2, kingSide, check))
		}
		if data := castles[2*uint8(player)+uint8(queenSide)]; pos.castlingRights.canCastle(player, queenSide) && bbOccupancy&data.bbTravel == 0 {
			check := data.rook2.bitboard()&bbChecks[Rook] > 0
			moves = append(moves, newCastleMove(king, data.king1, data.king2, queenSide, check))
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
				s2bb := s2.bitboard()

				if s2bb&(bbRank1^bbRank8) == 0 {
					check := s2bb&bbChecks[Pawn] > 0
					moves = append(moves, newPawnMove(pawn, NoPiece, s1, s2, NoSquare, NoPiece, check))
					continue
				}

				moves = append(moves,
					newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Queen.color(player), s2bb&bbChecks[Queen] > 0),
					newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Rook.color(player), s2bb&bbChecks[Rook] > 0),
					newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Bishop.color(player), s2bb&bbChecks[Bishop] > 0),
					newPawnMove(pawn, NoPiece, s1, s2, NoSquare, Knight.color(player), s2bb&bbChecks[Knight] > 0),
				)
			}
		}
	}

	// Pawn captures
	bbCaptureR, bbCaptureL := pawnCaptureBitboard(bbPawn, player)
	bbEnPassant := pos.enPassant.bitboard()

	bbPawnInterference := bbInterference
	if pos.enPassant != NoSquare &&
		(bbInterference&pos.board.bbPieces[Pawn]).scanForward()+Square(upOne) == pos.enPassant {
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
			s2bb := s2.bitboard()

			if s2bb&(bbRank1^bbRank8) == 0 {
				check := s2bb&bbChecks[Pawn] > 0
				moves = append(moves, newPawnMove(pawn, p2, s1, s2, pos.enPassant, NoPiece, check))
				continue
			}

			moves = append(moves,
				newPawnMove(pawn, p2, s1, s2, NoSquare, Queen.color(player), s2bb&bbChecks[Queen] > 0),
				newPawnMove(pawn, p2, s1, s2, NoSquare, Rook.color(player), s2bb&bbChecks[Rook] > 0),
				newPawnMove(pawn, p2, s1, s2, NoSquare, Bishop.color(player), s2bb&bbChecks[Bishop] > 0),
				newPawnMove(pawn, p2, s1, s2, NoSquare, Knight.color(player), s2bb&bbChecks[Knight] > 0),
			)
		}
	}

	// Other pieces
	for pt := Knight; pt <= Queen; pt++ {
		p1 := pt.color(player)
		for bbs1 := pos.board.bbPieces[pt] & bbPlayer; bbs1 > 0; bbs1 = bbs1.resetLSB() {
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
				check := s2.bitboard()&bbChecks[pt] > 0
				moves = append(moves, newPieceMove(p1, p2, s1, s2, check))
			}
		}
	}

	return moves
}

// isSquareAttacked checks whether the square is attacked by
// an enemy piece.
func (pos *Position) isSquareAttacked(sq Square) bool {
	return pos.attackedByBitboard(sq, pos.turn) > 0
}

// attackedByBitboard returns the bitboard of the pieces that attack the square.
func (pos *Position) attackedByBitboard(sq Square, c Color) bitboard {
	bbOpponent := pos.board.bbColors[c.Other()]
	bbOccupancy := pos.board.bbColors[White] ^ pos.board.bbColors[Black]
	bbRookMoves := bbMagicRookMoves[rookMagics[sq].index(bbOccupancy)]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[sq].index(bbOccupancy)]

	var bb bitboard
	bb |= singlePawnCaptureBitboard(sq, c) & pos.board.bbPieces[Pawn]
	bb |= bbKingMoves[sq] & pos.board.bbPieces[King]
	bb |= bbKnightMoves[sq] & pos.board.bbPieces[Knight]
	bb |= (pos.board.bbPieces[Queen] | pos.board.bbPieces[Rook]) & bbRookMoves
	bb |= (pos.board.bbPieces[Queen] | pos.board.bbPieces[Bishop]) & bbBishopMoves
	return bb & bbOpponent
}

// attackedAndDefendedByBitboard returns the bitboard of the pieces that attack and defend the square.
func (pos *Position) attackedAndDefendedByBitboard(sq Square, bbOccupancy bitboard) bitboard {
	bbRookMoves := bbMagicRookMoves[rookMagics[sq].index(bbOccupancy)]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[sq].index(bbOccupancy)]

	var bb bitboard
	bb |= singlePawnCaptureBitboard(sq, Black) & pos.board.bbPieces[Pawn]
	bb |= singlePawnCaptureBitboard(sq, White) & pos.board.bbPieces[Pawn]
	bb |= bbKingMoves[sq] & pos.board.bbPieces[King]
	bb |= bbKnightMoves[sq] & pos.board.bbPieces[Knight]
	bb |= (pos.board.bbPieces[Queen] | pos.board.bbPieces[Rook]) & bbRookMoves
	bb |= (pos.board.bbPieces[Queen] | pos.board.bbPieces[Bishop]) & bbBishopMoves
	return bb
}

// attackBitboards returns the bitboards where pieces would attack the square.
func (pos *Position) attackBitboards(sq Square, c Color) [5]bitboard {
	bbOccupancy := pos.board.bbColors[White] ^ pos.board.bbColors[Black]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[sq].index(bbOccupancy)]
	bbRookMoves := bbMagicRookMoves[rookMagics[sq].index(bbOccupancy)]

	return [5]bitboard{
		singlePawnCaptureBitboard(sq, c),
		bbKnightMoves[sq],
		bbBishopMoves,
		bbRookMoves,
		bbRookMoves | bbBishopMoves,
	}
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
	bbOpponent := pos.board.bbColors[pos.turn.Other()]
	cc := castleChecks[2*uint8(pos.turn)+uint8(s)]

	if cc.bbPawn&pos.board.bbPieces[Pawn]&bbOpponent > 0 ||
		cc.bbKnight&pos.board.bbPieces[Knight]&bbOpponent > 0 ||
		cc.bbKing&pos.board.bbPieces[King]&bbOpponent > 0 {
		return false
	}

	var bbBishopAttacks, bbRookAttacks bitboard
	bbOccupancy := pos.board.bbColors[White] ^ pos.board.bbColors[Black]
	for _, sq := range cc.squares {
		index := bishopMagics[sq].index(bbOccupancy)
		bbBishopAttacks |= bbMagicBishopMoves[index]
	}

	if bb := pos.board.bbPieces[Bishop] | pos.board.bbPieces[Queen]; bbBishopAttacks&bbOpponent&bb > 0 {
		return false
	}

	for _, sq := range cc.squares {
		index := rookMagics[sq].index(bbOccupancy)
		bbRookAttacks |= bbMagicRookMoves[index]
	}

	return bbRookAttacks&(pos.board.bbPieces[Rook]|pos.board.bbPieces[Queen])&bbOpponent == 0
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
	{1<<F8 | 1<<G8, E8, G8, F8},         // black, king side
	{1<<B8 | 1<<C8 | 1<<D8, E8, C8, D8}, // black, queen side
	{1<<F1 | 1<<G1, E1, G1, E1},         // white, king side
	{1<<B1 | 1<<C1 | 1<<D1, E1, C1, D1}, // white, queen side
}

// casteData represents a castle's data.
type castleData struct {
	bbTravel bitboard // bitboard traveled by the king
	king1    Square   // king s1
	king2    Square   // king s2
	rook2    Square   // rook s2
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
