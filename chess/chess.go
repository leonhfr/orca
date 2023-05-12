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
	for _, s := range []side{kingSide, queenSide} {
		data := castles[2*uint8(player)+uint8(s)]
		if pos.castlingRights.canCastle(player, s) && bbOccupancy&data.bbTravel == 0 {
			moves = append(moves, newMove(King.color(player), NoPiece, data.s1, data.s2, NoSquare, NoPiece))
		}
	}

	// Pawn moves
	bbUpOne, bbUpTwo := pawnMoveBitboard(bbPawn, bbOccupancy, player)
	for _, dest := range []bbDir{
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

	// Pawn Captures
	bbCaptureR, bbCaptureL := pawnCaptureBitboard(bbPawn, player)
	bbEnPassant := pos.enPassant.bitboard()
	for _, dest := range []bbDir{
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
	for _, origin := range []bbPt{
		{Knight, bbPlayer & pos.board.bbKnight},
		{Bishop, bbPlayer & pos.board.bbBishop},
		{Rook, bbPlayer & pos.board.bbRook},
		{Queen, bbPlayer & pos.board.bbQueen},
		{King, bbPlayer & pos.board.bbKing},
	} {
		p1 := origin.pt.color(player)
		for ; origin.bb > 0; origin.bb = origin.bb.resetLSB() {
			s1 := origin.bb.scanForward()
			for bbs2 := pieceBitboard(s1, p1.Type(), bbOccupancy) & ^bbPlayer; bbs2 > 0; bbs2 = bbs2.resetLSB() {
				s2 := bbs2.scanForward()
				p2 := pos.board.pieceByColor(s2, opponent)
				moves = append(moves, newMove(p1, p2, s1, s2, pos.enPassant, NoPiece))
			}
		}
	}

	return moves
}

// bbDir associates a bitboard with a direction
type bbDir struct {
	bb  bitboard
	dir direction
}

// bbPt associates a bitboard with a PieceType
type bbPt struct {
	pt PieceType
	bb bitboard
}

// castles contains the castles' data
//
// indexed by 2*Color+side
var castles = []castleData{
	{1<<F8 | 1<<G8, E8, G8},         // black, king side
	{1<<B8 | 1<<C8 | 1<<D8, E8, C8}, // black, queen side
	{1<<F1 | 1<<G1, E1, G1},         // white, king side
	{1<<B1 | 1<<C1 | 1<<D1, E1, C1}, // white, queen side
}

// casteData represents a castle's data
type castleData struct {
	bbTravel bitboard // bitboard traveled by the king
	s1       Square   // king s1
	s2       Square   // king s2
}
