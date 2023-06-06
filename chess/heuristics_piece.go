package chess

import "math/bits"

// CountPieces returns the count of knights, bishops, rooks, and queens.
func (pos *Position) CountPieces() (int, int, int, int) {
	return pos.board.bbPieces[Knight].ones(), pos.board.bbPieces[Bishop].ones(), pos.board.bbPieces[Rook].ones(), pos.board.bbPieces[Queen].ones()
}

// PieceMap executes the callback for each piece on the board, passing the piece
// and its square as arguments.
//
// Does not take pawns into account.
//
// Intended to be used in evaluation functions.
func (pos *Position) PieceMap(cb func(p Piece, sq Square)) {
	for p, bb := range [8]bitboard{
		pos.board.bbColors[Black] & pos.board.bbPieces[Knight],
		pos.board.bbColors[White] & pos.board.bbPieces[Knight],
		pos.board.bbColors[Black] & pos.board.bbPieces[Bishop],
		pos.board.bbColors[White] & pos.board.bbPieces[Bishop],
		pos.board.bbColors[Black] & pos.board.bbPieces[Rook],
		pos.board.bbColors[White] & pos.board.bbPieces[Rook],
		pos.board.bbColors[Black] & pos.board.bbPieces[Queen],
		pos.board.bbColors[White] & pos.board.bbPieces[Queen],
	} {
		piece := Piece(p + 2)
		for ; bb > 0; bb = bb.resetLSB() {
			cb(piece, bb.scanForward())
		}
	}
}

// UniquePieceMap executes the callback for each piece on the board that do not
// have an opponent mirrored piece, passing the piece and its square as arguments.
//
// Does not take pawns into account.
//
// Intended to be used in evaluation functions.
func (pos *Position) UniquePieceMap(cb func(p Piece, sq Square)) {
	bbBlack, bbWhite := pos.board.bbColors[Black], pos.board.bbColors[White]
	bbQueen := pos.board.bbPieces[Queen]
	bbRook := pos.board.bbPieces[Rook]
	bbBishop := pos.board.bbPieces[Bishop]
	bbKnight := pos.board.bbPieces[Knight]

	for p, bb := range [8]bitboard{
		^bitboard(bits.ReverseBytes64(uint64(bbWhite&bbKnight))) & bbBlack & bbKnight,
		^bitboard(bits.ReverseBytes64(uint64(bbBlack&bbKnight))) & bbWhite & bbKnight,
		^bitboard(bits.ReverseBytes64(uint64(bbWhite&bbBishop))) & bbBlack & bbBishop,
		^bitboard(bits.ReverseBytes64(uint64(bbBlack&bbBishop))) & bbWhite & bbBishop,
		^bitboard(bits.ReverseBytes64(uint64(bbWhite&bbRook))) & bbBlack & bbRook,
		^bitboard(bits.ReverseBytes64(uint64(bbBlack&bbRook))) & bbWhite & bbRook,
		bbBlack & bbQueen,
		bbWhite & bbQueen,
	} {
		piece := Piece(p + 2)
		for ; bb > 0; bb = bb.resetLSB() {
			cb(piece, bb.scanForward())
		}
	}
}
