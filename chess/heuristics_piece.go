package chess

// CountPieces returns the count of knights, bishops, rooks, and queens.
func (pos *Position) CountPieces() (int, int, int, int) {
	return pos.board.bbPieces[Knight].ones(),
		pos.board.bbPieces[Bishop].ones(),
		pos.board.bbPieces[Rook].ones(),
		pos.board.bbPieces[Queen].ones()
}

// PieceMap executes the callback for each piece on the board, passing the piece
// and its square as arguments.
//
// Does not take pawns or kings into account.
//
// Intended to be used in evaluation functions.
func (pos *Position) PieceMap(cb func(p Piece, sq Square, mobility int)) {
	bbOccupancy := pos.board.bbColors[Black] | pos.board.bbColors[White]

	for p := BlackKnight; p <= WhiteQueen; p++ {
		c := p.Color()
		pt := p.Type()
		bbPiece := pos.board.bbColors[c] & pos.board.bbPieces[pt]

		for ; bbPiece > 0; bbPiece = bbPiece.resetLSB() {
			sq := bbPiece.scanForward()

			mobility := pieceMobility(sq, pt, pos.board.bbColors[c], bbOccupancy)

			cb(p, sq, mobility)
		}
	}
}

// pieceMobility computes the mobility of the piece.
//
// May include illegal moves.
func pieceMobility(sq Square, pt PieceType, bbPlayer, bbOccupancy bitboard) int {
	bb := pieceBitboard(sq, pt, bbOccupancy) & ^bbPlayer
	return bb.ones()
}
