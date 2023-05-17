package chess

// PseudoMoves returns the list of pseudo moves.
//
// Some moves may be putting the moving player's king in check and therefore be illegal.
func (pos *Position) PseudoMoves() []Move {
	moves := make([]Move, 0, 50)

	// Setting up variables
	player, opponent := pos.turn, pos.turn.other()
	pawn := WhitePawn
	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	bbPlayer, bbOpponent := pos.board.bbWhite, pos.board.bbBlack
	upOne, upTwo := north, doubleNorth
	captureR, captureL := northEast, northWest
	if pos.turn == Black {
		pawn = BlackPawn
		bbPlayer, bbOpponent = pos.board.bbBlack, pos.board.bbWhite
		upOne, upTwo = south, doubleSouth
		captureR, captureL = southEast, southWest
	}
	bbPawn := pos.board.bbPawn & bbPlayer

	// Castles
	for _, s := range [2]side{kingSide, queenSide} {
		data := castles[2*uint8(player)+uint8(s)]
		if pos.castlingRights.canCastle(player, s) && bbOccupancy&data.bbTravel == 0 {
			moves = append(moves, newMove(King.color(player), NoPiece, data.s1, data.s2, NoSquare, NoPiece))
		}
	}

	// Pawn moves
	bbUpOne, bbUpTwo := pawnMoveBitboard(bbPawn, bbOccupancy, player)
	for _, dest := range [2]bbDir{
		{bbUpOne, upOne},
		{bbUpTwo, upTwo},
	} {
		for ; dest.bb > 0; dest.bb = dest.bb.resetLSB() {
			s2 := dest.bb.scanForward()
			s1 := s2 - Square(dest.dir)
			if pawn == WhitePawn && s2.Rank() == Rank8 || pawn == BlackPawn && s2.Rank() == Rank1 {
				moves = append(moves,
					newMove(pawn, NoPiece, s1, s2, NoSquare, Queen.color(player)),
					newMove(pawn, NoPiece, s1, s2, NoSquare, Rook.color(player)),
					newMove(pawn, NoPiece, s1, s2, NoSquare, Bishop.color(player)),
					newMove(pawn, NoPiece, s1, s2, NoSquare, Knight.color(player)),
				)
			} else {
				moves = append(moves, newMove(pawn, NoPiece, s1, s2, NoSquare, NoPiece))
			}
		}
	}

	// Pawn captures
	bbCaptureR, bbCaptureL := pawnCaptureBitboard(bbPawn, player)
	bbEnPassant := pos.enPassant.bitboard()
	for _, dest := range [2]bbDir{
		{bbCaptureR & (bbOpponent | bbEnPassant), captureR},
		{bbCaptureL & (bbOpponent | bbEnPassant), captureL},
	} {
		for ; dest.bb > 0; dest.bb = dest.bb.resetLSB() {
			s2 := dest.bb.scanForward()
			s1 := s2 - Square(dest.dir)
			p2 := pos.board.pieceByColor(s2, opponent)

			if pawn == WhitePawn && s2.Rank() == Rank8 || pawn == BlackPawn && s2.Rank() == Rank1 {
				moves = append(moves,
					newMove(pawn, p2, s1, s2, NoSquare, Queen.color(player)),
					newMove(pawn, p2, s1, s2, NoSquare, Rook.color(player)),
					newMove(pawn, p2, s1, s2, NoSquare, Bishop.color(player)),
					newMove(pawn, p2, s1, s2, NoSquare, Knight.color(player)),
				)
			} else {
				moves = append(moves, newMove(pawn, p2, s1, s2, pos.enPassant, NoPiece))
			}
		}
	}

	// Other pieces
	for _, origin := range [5]bbPt{
		{Knight, bbPlayer & pos.board.bbKnight},
		{Bishop, bbPlayer & pos.board.bbBishop},
		{Rook, bbPlayer & pos.board.bbRook},
		{Queen, bbPlayer & pos.board.bbQueen},
		{King, bbPlayer & pos.board.bbKing},
	} {
		p1 := origin.pt.color(player)
		for ; origin.bb > 0; origin.bb = origin.bb.resetLSB() {
			s1 := origin.bb.scanForward()
			for bbs2 := pieceBitboard(s1, origin.pt, bbOccupancy) & ^bbPlayer; bbs2 > 0; bbs2 = bbs2.resetLSB() {
				s2, p2 := bbs2.scanForward(), NoPiece
				if s2.bitboard()&bbOpponent > 0 {
					p2 = pos.board.pieceByColor(s2, opponent)
				}
				moves = append(moves, newMove(p1, p2, s1, s2, NoSquare, NoPiece))
			}
		}
	}

	return moves
}

// LoudMoves returns the list of pseudo loud moves.
// Loud moves are moves that capture an opponent piece.
//
// Some moves may be putting the moving player's king in check and therefore be illegal.
func (pos *Position) LoudMoves() []Move {
	moves := make([]Move, 0, 20)

	// Setting up variables
	player, opponent := pos.turn, pos.turn.other()
	pawn := WhitePawn
	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	bbPlayer, bbOpponent := pos.board.bbWhite, pos.board.bbBlack
	captureR, captureL := northEast, northWest
	if pos.turn == Black {
		pawn = BlackPawn
		bbPlayer, bbOpponent = pos.board.bbBlack, pos.board.bbWhite
		captureR, captureL = southEast, southWest
	}
	bbPawn := pos.board.bbPawn & bbPlayer

	// Pawn captures
	bbCaptureR, bbCaptureL := pawnCaptureBitboard(bbPawn, player)
	bbEnPassant := pos.enPassant.bitboard()
	for _, dest := range [2]bbDir{
		{bbCaptureR & (bbOpponent | bbEnPassant), captureR},
		{bbCaptureL & (bbOpponent | bbEnPassant), captureL},
	} {
		for ; dest.bb > 0; dest.bb = dest.bb.resetLSB() {
			s2 := dest.bb.scanForward()
			s1 := s2 - Square(dest.dir)
			p2 := pos.board.pieceByColor(s2, opponent)

			if pawn == WhitePawn && s2.Rank() == Rank8 || pawn == BlackPawn && s2.Rank() == Rank1 {
				moves = append(moves,
					newMove(pawn, p2, s1, s2, NoSquare, Queen.color(player)),
					newMove(pawn, p2, s1, s2, NoSquare, Rook.color(player)),
					newMove(pawn, p2, s1, s2, NoSquare, Bishop.color(player)),
					newMove(pawn, p2, s1, s2, NoSquare, Knight.color(player)),
				)
			} else {
				moves = append(moves, newMove(pawn, p2, s1, s2, pos.enPassant, NoPiece))
			}
		}
	}

	// Other pieces
	for _, origin := range [5]bbPt{
		{Knight, bbPlayer & pos.board.bbKnight},
		{Bishop, bbPlayer & pos.board.bbBishop},
		{Rook, bbPlayer & pos.board.bbRook},
		{Queen, bbPlayer & pos.board.bbQueen},
		{King, bbPlayer & pos.board.bbKing},
	} {
		p1 := origin.pt.color(player)
		for ; origin.bb > 0; origin.bb = origin.bb.resetLSB() {
			s1 := origin.bb.scanForward()
			for bbs2 := pieceBitboard(s1, origin.pt, bbOccupancy) & bbOpponent; bbs2 > 0; bbs2 = bbs2.resetLSB() {
				s2 := bbs2.scanForward()
				p2 := pos.board.pieceByColor(s2, opponent)
				moves = append(moves, newMove(p1, p2, s1, s2, NoSquare, NoPiece))
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
	bbOpponent := pos.board.getColor(pos.turn.other())
	if singlePawnCaptureBitboard(sq, pos.turn)&pos.board.bbPawn&bbOpponent > 0 ||
		bbKingMoves[sq]&pos.board.bbKing&bbOpponent > 0 ||
		bbKnightMoves[sq]&pos.board.bbKnight&bbOpponent > 0 {
		return true
	}

	bbOccupancy := pos.board.bbWhite ^ pos.board.bbBlack
	bbRookMoves := bbMagicRookMoves[rookMagics[sq].index(bbOccupancy)]
	bbBishopMoves := bbMagicBishopMoves[bishopMagics[sq].index(bbOccupancy)]

	return (pos.board.bbQueen|pos.board.bbRook)&bbRookMoves&bbOpponent > 0 ||
		(pos.board.bbQueen|pos.board.bbBishop)&bbBishopMoves&bbOpponent > 0
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

// bbPt associates a bitboard with a PieceType.
type bbPt struct {
	pt PieceType
	bb bitboard
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
