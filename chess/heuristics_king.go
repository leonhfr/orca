package chess

// KingMap executes the callback for both kings on the board.
//
// Intended to be used in evaluation functions.
func (pos *Position) KingMap(cb func(p Piece, sq Square, shieldDefects, openFiles, halfOpenFiles int)) {
	bbOpenFiles, bbHalfOpenFiles := pos.fileData()

	for c := Black; c <= White; c++ {
		king := King.color(c)
		bb := pos.board.bbColors[c] & pos.board.bbPieces[King]
		sq := bb.scanForward()

		bbFiles := bb | bb.eastOne() | bb.westOne()
		openFiles := (bbOpenFiles & bbFiles).ones()
		halfOpenFiles := (bbHalfOpenFiles[c] & bbFiles).ones()

		var shieldDefects int
		for _, psm := range pawnShieldMasks[c] {
			if psm[0]&bb > 0 {
				bbPawns := psm[1] & pos.board.bbColors[c] & pos.board.bbPieces[Pawn]
				pawnCount := bbPawns.ones()

				pawnCount = min(pawnCount, pawnShieldCount)
				shieldDefects = pawnShieldCount - pawnCount
				break
			}
		}

		cb(king, sq, shieldDefects, openFiles, halfOpenFiles)
	}
}

// fileData computes open files and half open files, indexed by color.
func (pos *Position) fileData() (bitboard, [2]bitboard) {
	bbBlackPawn := pos.board.bbColors[Black] & pos.board.bbPieces[Pawn]
	bbWhitePawn := pos.board.bbColors[White] & pos.board.bbPieces[Pawn]
	bbBlackFileFill := bbBlackPawn.fileFill()
	bbWhiteFileFill := bbWhitePawn.fileFill()

	bbBlackHalfOpenFile := ^bbBlackFileFill
	bbWhiteHalfOpenFile := ^bbWhiteFileFill
	bbOpenFiles := bbBlackHalfOpenFile & bbWhiteHalfOpenFile

	return bbOpenFiles, [2]bitboard{bbBlackHalfOpenFile, bbWhiteHalfOpenFile}
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
