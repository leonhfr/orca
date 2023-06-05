package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var zobristTests = []struct {
	name     string
	args     string
	move     Move
	hash     Hash
	pawnHash Hash
}{
	// Variation 1
	{
		"starting position",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		newMove(WhitePawn, NoPiece, E2, E4, NoSquare, NoPiece),
		0x463b96181691fc9c,
		0x91af2b0abcd2875d,
	},
	{
		"position after e2e4",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		newMove(BlackPawn, NoPiece, D7, D5, E3, NoPiece),
		0x823c9b50fd114196,
		0xad7e00e8f875bf5e,
	},
	{
		"position after e2e4 d7d5",
		"rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
		newMove(WhitePawn, NoPiece, E4, E5, D6, NoPiece),
		0x0756b94461c50fb0,
		0xd0c20456cb867471,
	},
	{
		"position after e2e4 d7d5 e4e5",
		"rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		newMove(BlackPawn, NoPiece, F7, F5, NoSquare, NoPiece),
		0x662fafb965db29d4,
		0x496d340160bfd71c,
	},
	{
		"position after e2e4 d7d5 e4e5 f7f5",
		"rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
		newMove(WhiteKing, NoPiece, E1, E2, F6, NoPiece),
		0x22a48b5a8e47ff78,
		0x25d4743271107fcb,
	},
	{
		"position after e2e4 d7d5 e4e5 f7f5 e1e2",
		"rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR b kq - 0 3",
		newMove(BlackKing, NoPiece, E8, F7, NoSquare, NoPiece),
		0x652a607ca3f242c1,
		0x8ada538d1dadfe89,
	},
	{
		"position after e2e4 d7d5 e4e5 f7f5 e1e2 e8f7",
		"rnbq1bnr/ppp1pkpp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR w - - 0 4",
		NoMove,
		0x00fdd303c946bdd9,
		0xac5343bab48469f1,
	},
	// Variation 2
	{
		"position after a2a4 b7b5 h2h4 b5b4 c2c4",
		"rnbqkbnr/p1pppppp/8/8/PpP4P/8/1P1PPPP1/RNBQKBNR b KQkq c3 0 3",
		newMove(BlackPawn, WhitePawn, B4, C3, C3, NoPiece),
		0x3c8123ea7b067637,
		0x13f92b8acce2e19d,
	},
	{
		"position after a2a4 b7b5 h2h4 b5b4 c2c4 b4c3",
		"rnbqkbnr/p1pppppp/8/8/P6P/2p5/1P1PPPP1/RNBQKBNR w KQkq - 0 3",
		newMove(WhiteRook, NoPiece, A1, A3, NoSquare, NoPiece),
		0x93d32682782edfae,
		0x44479b90d26da46f,
	},
	{
		"position after a2a4 b7b5 h2h4 b5b4 c2c4 b4c3 a1a3",
		"rnbqkbnr/p1pppppp/8/8/P6P/R1p5/1P1PPPP1/1NBQKBNR b Kkq - 0 4",
		NoMove,
		0x5c3f9b829b279560,
		0x44479b90d26da46f,
	},
	// Variation 3
	{
		"edge cases",
		"rnbqkbnr/ppp1pppp/8/8/2B1P3/3p1N2/PPPP1PPP/RNBQK2R w KQkq - 0 1",
		newMove(WhiteKing, NoPiece, E1, G1, NoSquare, NoPiece),
		0xc1d8a7be522690d4,
		0xae9b4272b8c05c8e,
	},
	{
		"edge cases after e1g1",
		"rnbqkbnr/ppp1pppp/8/8/2B1P3/3p1N2/PPPP1PPP/RNBQ1RK1 b kq - 0 1",
		newMove(BlackPawn, WhitePawn, D3, C2, NoSquare, NoPiece),
		0xbfaddc8a2ed54bc4,
		0xf3190569a888ae18,
	},
	{
		"edge cases after e1g1 d3c2",
		"rnbqkbnr/ppp1pppp/8/8/2B1P3/5N2/PPpP1PPP/RNBQ1RK1 w kq - 0 1",
		newMove(WhiteQueen, NoPiece, D1, E2, NoSquare, NoPiece),
		0x76d0c8a9800cf271,
		0xc2b237e0a97692a4,
	},
	{
		"edge cases after e1g1 d3c2 d1e2",
		"rnbqkbnr/ppp1pppp/8/8/2B1P3/5N2/PPpPQPPP/RNB2RK1 b kq - 0 1",
		newMove(BlackPawn, WhiteKnight, C2, B1, NoSquare, BlackQueen),
		0xe0d6478bce48a9f5,
		0xc2b237e0a97692a4,
	},
	{
		"edge cases after e1g1 d3c2 d1e2 c2b1",
		"rnbqkbnr/ppp1pppp/8/8/2B1P3/5N2/PP1PQPPP/RqB2RK1 w kq - 0 1",
		NoMove,
		0xee91778b335c5907,
		0x54b7e2104b285114,
	},
	// From panics
	{
		"panic 1",
		"r1b1kb1r/pppp1ppp/2n1p3/8/3P3P/2P5/PP1BPP1P/R2QKB1R b KQkq - 0 8",
		newMove(BlackKnight, WhitePawn, C6, D4, NoSquare, NoPiece),
		0x9fa15cde880e1a00,
		0x8595253885ed0948,
	},
	{
		"panic 1 after c6d4",
		"r1b1kb1r/pppp1ppp/4p3/8/3n3P/2P5/PP1BPP1P/R2QKB1R w KQkq - 0 8",
		NoMove,
		0x51fcf072716c211c,
		0x9cb5e575c2cb72f5,
	},
}

func TestZobristHash(t *testing.T) {
	for _, tt := range zobristTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.args)
			key := newZobristHash(pos)
			assert.Equal(t, tt.hash, key)
		})
	}
}

func TestPawnZobristHash(t *testing.T) {
	for _, tt := range zobristTests {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.args)
			key := newPawnZobristHash(pos)
			assert.Equal(t, tt.pawnHash, key)
		})
	}
}

func TestIncrementalZobristHash(t *testing.T) {
	for i, tt := range zobristTests {
		if tt.move != NoMove {
			t.Run(fmt.Sprintf("%s (%s)", tt.name, tt.move.String()), func(t *testing.T) {
				pos := unsafeFEN(tt.args)
				pos.MakeMove(tt.move)

				assert.Equal(t, zobristTests[i+1].hash, pos.Hash())
				assert.Equal(t, zobristTests[i+1].pawnHash, pos.PawnHash())
			})
		}
	}
}

func BenchmarkZobristHash(b *testing.B) {
	pos := unsafeFEN("rnbqkbnr/p1pppppp/8/8/PpP4P/8/1P1PPPP1/RNBQKBNR b KQkq c3 0 3")
	for n := 0; n < b.N; n++ {
		newZobristHash(pos)
	}
}
