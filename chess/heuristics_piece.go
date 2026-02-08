package chess

// CountPieces returns the counts of knights, bishops, rooks, and queens.
func (pos *Position) CountPieces() (int, int, int, int) {
	return pos.board.bbPieces[Knight].ones(),
		pos.board.bbPieces[Bishop].ones(),
		pos.board.bbPieces[Rook].ones(),
		pos.board.bbPieces[Queen].ones()
}

// CountOwnPieces returns the count of knights, bishops, rooks, and queens.
func (pos *Position) CountOwnPieces() int {
	bbPieces := pos.board.bbPieces[Knight] ^
		pos.board.bbPieces[Bishop] ^
		pos.board.bbPieces[Rook] ^
		pos.board.bbPieces[Queen]
	return (pos.board.bbColors[pos.turn] & bbPieces).ones()
}

// FileData contains data on half open and open files for a particular pawn structure.
type FileData [3]bitboard

// FileData computes data  on closed, half open, and open files
// for a particular pawn structure.
func (pos *Position) FileData() FileData {
	bbBlackPawn := pos.board.bbColors[Black] & pos.board.bbPieces[Pawn]
	bbWhitePawn := pos.board.bbColors[White] & pos.board.bbPieces[Pawn]
	bbBlackFileFill := bbBlackPawn.fileFill()
	bbWhiteFileFill := bbWhitePawn.fileFill()

	bbBlackHalfOpenFile := ^bbBlackFileFill
	bbWhiteHalfOpenFile := ^bbWhiteFileFill
	bbOpenFiles := bbBlackHalfOpenFile & bbWhiteHalfOpenFile

	return FileData{
		bbBlackHalfOpenFile,
		bbWhiteHalfOpenFile,
		bbOpenFiles,
	}
}

// OnHalfOpenFile returns true if the square is on an half file of the given color.
func (fd FileData) OnHalfOpenFile(sq Square, c Color) bool {
	return fd[c]&sq.bitboard() > 0
}

// OnOpenFile returns true if the square is on an open file.
func (fd FileData) OnOpenFile(sq Square) bool {
	return fd[2]&sq.bitboard() > 0
}

// PieceProperty represents different piece properties.
type PieceProperty uint8

// NoPieceProperty represents the absence of properties.
const NoPieceProperty PieceProperty = 0

const (
	// Trapped represents a trapped piece. It should be hard for the piece to escape.
	Trapped PieceProperty = 1 << iota
	// Lost represents a lost piece. It should be impossible for the piece to escape.
	Lost
)

// HasProperty checks the presence of the given property.
func (pp PieceProperty) HasProperty(p PieceProperty) bool {
	return pp&p > 0
}

// PieceMap executes the callback for each piece on the board, passing the piece
// and its square as arguments.
//
// Does not take pawns or kings into account.
//
// Intended to be used in evaluation functions.
func (pos *Position) PieceMap(cb func(p Piece, sq Square, mobility int, properties PieceProperty)) {
	bbOccupancy := pos.board.bbColors[Black] | pos.board.bbColors[White]

	for p := BlackKnight; p <= WhiteQueen; p++ {
		c, op := p.Color(), p.Color().Other()
		pt := p.Type()
		bbPiece := pos.board.bbColors[c] & pos.board.bbPieces[pt]

		for ; bbPiece > 0; bbPiece = bbPiece.resetLSB() {
			sq := bbPiece.scanForward()
			bb := sq.bitboard()

			mobility := pieceMobility(sq, pt, pos.board.bbColors[c], bbOccupancy)

			var properties PieceProperty
			for _, ts := range pieceSituations[p] {
				if ts.bbPlayerPiece&bb > 0 {
					if ts.bbOpponentPiece&pos.board.bbColors[op]&pos.board.bbPieces[ts.opponentPieceType] > 0 {
						properties ^= Trapped
					}
					break
				}
			}

			cb(p, sq, mobility, properties)
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

// pieceSituation represents a piece situation.
type pieceSituation struct {
	bbPlayerPiece     bitboard
	bbOpponentPiece   bitboard
	opponentPieceType PieceType
	property          PieceProperty
}

// pieceSituations contains the piece situations, indexed by piece type.
var pieceSituations = [10][]pieceSituation{
	{}, // black pawns
	{}, // white pawns
	{ // black knights, same as below
		{1 << H1, 1 << H2, Pawn, Trapped},
		{1 << H1, 1 << F2, Pawn, Trapped},
		{1 << H2, 1<<G2 | 1<<H3, Pawn, Trapped},
		{1 << H2, 1<<G2 | 1<<F3, Pawn, Trapped},
	},
	{ // white knights
		{1 << H8, 1 << H7, Pawn, Trapped},       // a white knight on h8 with a black pawn on h7
		{1 << H8, 1 << F7, Pawn, Trapped},       // a white knight on h8 with a black pawn on f7
		{1 << H7, 1<<G7 | 1<<H6, Pawn, Trapped}, // a white knight on h7 with black pawns on g7 and h6
		{1 << H7, 1<<G7 | 1<<F6, Pawn, Trapped}, // a white knight on h7 with black pawns on g7 and f6
	},
	{ // black bishops, same as below
		{1 << H2, 1<<F2 | 1<<G3, Pawn, Lost},
		{1 << H3, 1<<G4 | 1<<F3, Pawn, Trapped},
	},
	{ // white bishops
		{1 << H7, 1<<F7 | 1<<G6, Pawn, Lost},
		{1 << H6, 1<<G5 | 1<<F6, Pawn, Trapped},
	},
	{ // black rooks, same as below
		{1<<G8 | 1<<H8 | 1<<G7 | 1<<H7, 1<<F8 | 1<<G8, King, Trapped},
	},
	{ // white rooks
		{1<<G1 | 1<<H1 | 1<<G2 | 1<<H2, 1<<F1 | 1<<G1, King, Trapped}, // a white rook on h1/g1/h2/g2 with a white king on f1 or g1
	},
	{}, // black queens
	{}, // white queens
}
