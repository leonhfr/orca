package chess

// PawnMap executes the callback for each pawn on the board.
// Intended to be used in evaluation functions.
func (pos *Position) PawnMap(cb func(p Piece, sq Square)) {
	for c := Black; c <= White; c++ {
		bbOwnPawn := pos.board.bbColors[c] & pos.board.bbPieces[Pawn]
		pawn := Pawn.color(c)

		for bb := bbOwnPawn; bb > 0; bb = bb.resetLSB() {
			sq := bb.scanForward()

			cb(pawn, sq)
		}
	}
}
