package chess

import "math/bits"

// KingMap executes the callback for both kings on the board.
//
// Intended to be used in evaluation functions.
func (pos *Position) KingMap(cb func(p Piece, sq Square, shieldDefects int)) {
	for c := Black; c <= White; c++ {
		king := King.color(c)
		bb := pos.board.bbColors[c] & pos.board.bbPieces[King]
		sq := bb.scanForward()

		var shieldDefects int
		for _, psm := range pawnShieldMasks[c] {
			if psm[0]&bb > 0 {
				bbPawns := psm[1] & pos.board.bbColors[c] & pos.board.bbPieces[Pawn]
				pawnCount := bits.OnesCount64(uint64(bbPawns))

				if pawnCount > pawnShieldCount {
					pawnCount = pawnShieldCount
				}
				shieldDefects = pawnShieldCount - pawnCount
				break
			}
		}

		cb(king, sq, shieldDefects)
	}
}

// pawnShieldMasks contains the shield masks indexed by color.
var pawnShieldMasks = [2][2][2]bitboard{
	{
		{blackKingSideKingMask, blackKingSidePawnMask},
		{blackQueenSideKingMask, blackQueenSidePawnMask},
	},
	{
		{whiteKingSideKingMask, whiteKingSidePawnMask},
		{whiteQueenSideKingMask, whiteQueenSidePawnMask},
	},
}

const (
	pawnShieldCount = 3

	whiteQueenSideKingMask = 1<<A1 | 1<<B1 | 1<<C1
	whiteKingSideKingMask  = 1<<F1 | 1<<G1 | 1<<H1
	blackQueenSideKingMask = whiteQueenSideKingMask << (8 * 7)
	blackKingSideKingMask  = whiteKingSideKingMask << (8 * 7)

	whiteQueenSidePawnMask = 1<<B2 | 1<<C2 | 1<<D2 | 1<<B3 | 1<<C3 | 1<<D3
	whiteKingSidePawnMask  = 1<<F2 | 1<<G2 | 1<<H2 | 1<<F3 | 1<<G3 | 1<<H3
	blackQueenSidePawnMask = whiteQueenSidePawnMask << (8 * 4)
	blackKingSidePawnMask  = whiteKingSidePawnMask << (8 * 4)
)
