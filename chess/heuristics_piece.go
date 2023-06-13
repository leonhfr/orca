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
func (pos *Position) PieceMap(cb func(p Piece, sq Square, mobility int, trapped bool)) {
	bbOccupancy := pos.board.bbColors[Black] | pos.board.bbColors[White]

	for p := BlackKnight; p <= WhiteQueen; p++ {
		c, op := p.Color(), p.Color().Other()
		pt := p.Type()
		bbPiece := pos.board.bbColors[c] & pos.board.bbPieces[pt]

		for ; bbPiece > 0; bbPiece = bbPiece.resetLSB() {
			sq := bbPiece.scanForward()
			bb := sq.bitboard()

			mobility := pieceMobility(sq, pt, pos.board.bbColors[c], bbOccupancy)

			var trapped bool
			for _, ts := range trappedSituations[p] {
				if ts.bbTrapped&bb > 0 {
					if ts.bbTrapping&pos.board.bbColors[op]&pos.board.bbPieces[ts.pt] > 0 {
						trapped = true
					}
					break
				}
			}

			cb(p, sq, mobility, trapped)
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

// trappedSituation represents a trapped piece situation.
type trappedSituation struct {
	bbTrapped  bitboard
	bbTrapping bitboard
	pt         PieceType
}

// trappedSituations contains the trapped situations, indexed by piece type.
var trappedSituations = [10][]trappedSituation{
	{}, // black pawns
	{}, // white pawns
	{ // black knights, same as below
		{1 << H1, 1 << H2, Pawn},
		{1 << H1, 1 << F2, Pawn},
		{1 << H2, 1<<G2 | 1<<H3, Pawn},
		{1 << H2, 1<<G2 | 1<<F3, Pawn},
	},
	{ // white knights
		{1 << H8, 1 << H7, Pawn},       // a white knight on h8 with a black pawn on h7
		{1 << H8, 1 << F7, Pawn},       // a white knight on h8 with a black pawn on f7
		{1 << H7, 1<<G7 | 1<<H6, Pawn}, // a white knight on h7 with black pawns on g7 and h6
		{1 << H7, 1<<G7 | 1<<F6, Pawn}, // a white knight on h7 with black pawns on g7 and f6
	},
	{ // black bishops, same as below
		{1 << H2, 1<<F2 | 1<<G3, Pawn},
		{1 << H3, 1<<G4 | 1<<F3, Pawn},
	},
	{ // white bishops
		{1 << H7, 1<<F7 | 1<<G6, Pawn},
		{1 << H6, 1<<G5 | 1<<F6, Pawn},
	},
	{ // black rooks, same as below
		{1<<G8 | 1<<H8 | 1<<G7 | 1<<H7, 1<<F8 | 1<<G8, King},
	},
	{ // white rooks
		{1<<G1 | 1<<H1 | 1<<G2 | 1<<H2, 1<<F1 | 1<<G1, King}, // a white rook on h1/g1/h2/g2 with a white king on f1 or g1
	},
	{}, // black queens
	{}, // white queens
}
