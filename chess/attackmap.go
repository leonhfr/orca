package chess

// attackMap contains a bitboard map of all pieces that attack a square.
type attackMap [64]bitboard

// newAttackMap creates a new attackMap from a position.
func newAttackMap(pos *Position) attackMap {
	var am attackMap
	occupancy := pos.board.bbWhite ^ pos.board.bbBlack
	for p := BlackPawn; p <= WhiteKing; p++ {
		for bb := pos.board.getBitboard(p); bb > 0; bb = bb.resetLSB() {
			am.setAttacks(bb.scanForward(), p, occupancy)
		}
	}
	return am
}

// updateAttackMap performs a xor operation on the attack map bitboards for a move.
//
// s1, s2, p1, p2 are different than the move definition.
// They should reflect the final state.
// Hence, p1 typically is NoPiece as the piece was moved from the spot.
// We use the p1 parameter when unmaking a move.
//
// Both ranks and bishops bitboards are expected to contain the queens as well.
func (am *attackMap) updateAttackMap(s1, s2 Square, p1, p2 Piece, occupancy, rooks, bishops bitboard) {
	bbRooks := (am[s1] | am[s2]) & rooks
	bbBishops := (am[s1] | am[s2]) & bishops
	unset := s1.bitboard() | s2.bitboard() | bbRooks | bbBishops
	for sq := A1; sq <= H8; sq++ {
		am[sq] &= ^unset
	}
	if p1 != NoPiece {
		am.setAttacks(s1, p1, occupancy)
	}
	if p2 != NoPiece {
		am.setAttacks(s2, p2, occupancy)
	}
	for ; bbRooks > 0; bbRooks = bbRooks.resetLSB() {
		am.setAttacks(bbRooks.scanForward(), WhiteRook, occupancy)
	}
	for ; bbBishops > 0; bbBishops = bbBishops.resetLSB() {
		am.setAttacks(bbBishops.scanForward(), WhiteBishop, occupancy)
	}
}

// setAttacks performs a xor operation on the attack map bitboards for a piece on a given square.
func (am *attackMap) setAttacks(sq Square, p Piece, occupancy bitboard) {
	for bb, attacks := sq.bitboard(), attackBitboard(sq, p, occupancy); attacks > 0; attacks = attacks.resetLSB() {
		am[attacks.scanForward()] |= bb
	}
}
