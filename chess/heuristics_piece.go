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
func (pos *Position) PieceMap(cb func(p Piece, sq Square)) {
	for p := BlackKnight; p <= WhiteQueen; p++ {
		bbPiece := pos.board.bbColors[p.Color()] & pos.board.bbPieces[p.Type()]
		for ; bbPiece > 0; bbPiece = bbPiece.resetLSB() {
			sq := bbPiece.scanForward()

			cb(p, sq)
		}
	}
}
