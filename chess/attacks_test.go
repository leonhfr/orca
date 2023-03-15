package chess

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPieceBitboard(t *testing.T) {
	type args struct {
		fen string
		sq  Square
		pt  PieceType
	}
	tests := []struct {
		args args
		want []Square
	}{
		{args{"k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1", A1, King}, []Square{
			A2, B2,
			B1, // will be removed
		}},
		{args{"k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1", B1, Queen}, []Square{
			A2, B2, B3, B4, B5,
			B6, B7, B8, C2, D3,
			E4, F5, G6, H7,
			A1, C1, // will be removed
		}},
		{args{"k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1", C1, Rook}, []Square{
			C2, C3, C4, C5, C6,
			C7, C8,
			B1, D1, // will be removed
		}},
		{args{"k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1", D1, Bishop}, []Square{
			A4, B3, C2, E2, F3,
			G4, H5,
		}},
		{args{"k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1", E1, Knight}, []Square{C2, D3, F3, G2}},
		{args{"k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1", F2, NoPieceType}, []Square{}},
		{args{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", A1, Rook}, []Square{B1, C1, D1, E1, A2}},
		{args{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", H1, Rook}, []Square{E1, F1, G1, H2}},
		{args{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", D2, Bishop}, []Square{C1, E1, C3, E3, F4, G5, H6}},
		{args{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", E2, Bishop}, []Square{D1, F1, D3, C4, B5, A6, F3}},
	}

	for _, tt := range tests {
		t.Run(tt.args.pt.String()+" "+tt.args.sq.String(), func(t *testing.T) {
			pos := unsafeFEN(tt.args.fen)
			occupancy := pos.board.getColor(White) ^ pos.board.getColor(Black)
			got := pieceBitboard(tt.args.sq, tt.args.pt, occupancy)
			assert.ElementsMatch(t, tt.want, got.mapping())
		})
	}
}

func TestPawnMoveBitboard(t *testing.T) {
	tests := []struct {
		fen   string
		upOne []Square
		upTwo []Square
	}{
		{"k7/3p4/4p3/8/2pP4/8/5P2/K7 w - - 0 1", []Square{F3, D5}, []Square{F4}},
		{"k7/3p4/4p3/5N2/2pP4/8/5P2/K7 b - d3 0 1", []Square{C3, D6, E5}, []Square{D5}},
	}

	for _, tt := range tests {
		t.Run(tt.fen, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			pawn := pos.board.getBitboard(Pawn.color(pos.turn))
			occupancy := pos.board.getColor(White) ^ pos.board.getColor(Black)
			upOne, upTwo := pawnMoveBitboard(pawn, occupancy, pos.turn)
			assert.ElementsMatch(t, tt.upOne, upOne.mapping())
			assert.ElementsMatch(t, tt.upTwo, upTwo.mapping())
		})
	}
}

func TestPawnCaptureBitboard(t *testing.T) {
	tests := []struct {
		fen      string
		captureR []Square
		captureL []Square
	}{
		{"k7/3p4/4p3/8/2pP4/8/5P2/K7 w - - 0 1", []Square{E5, G3}, []Square{C5, E3}},
		{"k7/3p4/4p3/5N2/2pP4/8/5P2/K7 b - d3 0 1", []Square{D3, E6, F5}, []Square{B3, C6, D5}},
	}

	for _, tt := range tests {
		t.Run(tt.fen, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			pawn := pos.board.getBitboard(Pawn.color(pos.turn))
			captureR, captureL := pawnCaptureBitboard(pawn, pos.turn)
			assert.ElementsMatch(t, tt.captureR, captureR.mapping())
			assert.ElementsMatch(t, tt.captureL, captureL.mapping())
		})
	}
}

func TestRookBishopMoves(t *testing.T) {
	magics := [][64]Magic{rookMagics, bishopMagics}
	moves := [][]bitboard{bbMagicRookMoves, bbMagicBishopMoves}

	for pieceIndex, pt := range []PieceType{Rook, Bishop} {
		t.Run(pt.String(), func(t *testing.T) {
			for i := 0; i < 20; i++ {
				bb, sq := randomPosition()
				index := magics[pieceIndex][sq].index(bb)
				actual := moves[pieceIndex][index]
				expected := slowMoves(pt, sq, bb)

				assert.Equal(t, expected, actual)
			}
		})
	}
}

func BenchmarkRookBishopMoves(b *testing.B) {
	magics := [][64]Magic{rookMagics, bishopMagics}
	moves := [][]bitboard{bbMagicRookMoves, bbMagicBishopMoves}
	bb, sq := randomPosition()

	for pieceIndex, pt := range []PieceType{Rook, Bishop} {
		b.Run(pt.String()+"-magic", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				index := magics[pieceIndex][sq].index(bb)
				_ = moves[pieceIndex][index]
			}
		})
		b.Run(pt.String()+"-slow", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				slowMoves(pt, sq, bb)
			}
		})
	}
}

func randomPosition() (bitboard, Square) {
	//nolint:gosec
	bb, sq := bitboard(rand.Uint64()), Square(rand.Intn(64))
	return bb, sq
}
