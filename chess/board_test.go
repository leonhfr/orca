package chess

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	startingSquareMap = map[Square]Piece{
		A8: BlackRook, B8: BlackKnight, C8: BlackBishop, D8: BlackQueen,
		E8: BlackKing, F8: BlackBishop, G8: BlackKnight, H8: BlackRook,
		A7: BlackPawn, B7: BlackPawn, C7: BlackPawn, D7: BlackPawn,
		E7: BlackPawn, F7: BlackPawn, G7: BlackPawn, H7: BlackPawn,
		A2: WhitePawn, B2: WhitePawn, C2: WhitePawn, D2: WhitePawn,
		E2: WhitePawn, F2: WhitePawn, G2: WhitePawn, H2: WhitePawn,
		A1: WhiteRook, B1: WhiteKnight, C1: WhiteBishop, D1: WhiteQueen,
		E1: WhiteKing, F1: WhiteBishop, G1: WhiteKnight, H1: WhiteRook,
	}

	startingBoard = board{
		bbPieces: [6]bitboard{
			71776119061282560,
			4755801206503243842,
			2594073385365405732,
			9295429630892703873,
			576460752303423496,
			1152921504606846992,
		},
		bbColors: [2]bitboard{
			18446462598732840960,
			65535,
		},
		sqKings: [2]Square{E8, E1},
	}

	testPositions = []struct {
		move    Move
		moveUCI string
		tags    []MoveTag
		preFEN  string
		postFEN string
	}{
		{
			move:    newMove(WhitePawn, NoPiece, E2, E4, NoSquare, NoPiece),
			moveUCI: "e2e4",
			tags:    []MoveTag{Quiet},
			preFEN:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			postFEN: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			move:    newMove(BlackPawn, NoPiece, A4, B3, B3, NoPiece),
			moveUCI: "a4b3",
			tags:    []MoveTag{Capture},
			preFEN:  "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
			postFEN: "2r3k1/1q1nbppp/r3p3/3pP3/2pP4/PpQ2N2/2RN1PPP/2R4K w - - 0 24",
		},
		{
			move:    newMove(WhiteKing, NoPiece, E1, G1, NoSquare, NoPiece),
			moveUCI: "e1g1",
			tags:    []MoveTag{HSideCastle},
			preFEN:  "r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9",
			postFEN: "r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B2RK1 b kq - 2 9",
		},
		{
			move:    newMove(WhiteKnight, NoPiece, G1, F3, NoSquare, NoPiece),
			moveUCI: "g1f3",
			tags:    []MoveTag{Quiet},
			preFEN:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			postFEN: "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
		},
		{
			move:    newMove(WhiteKing, NoPiece, E1, C1, NoSquare, NoPiece),
			moveUCI: "e1c1",
			tags:    []MoveTag{ASideCastle},
			preFEN:  "r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/R3K2R w KQkq - 3 10",
			postFEN: "r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/2KR3R b kq - 4 10",
		},
		{
			move:    newMove(WhitePawn, BlackPawn, E4, D5, NoSquare, NoPiece),
			moveUCI: "e4d5",
			tags:    []MoveTag{Capture},
			preFEN:  "r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
			postFEN: "r1bqkbnr/ppp1pppp/2n5/3P4/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 3",
		},
		{
			move:    newMove(WhitePawn, NoPiece, E5, F6, F6, NoPiece),
			moveUCI: "e5f6",
			tags:    []MoveTag{Capture},
			preFEN:  "r1bqkbnr/ppp1p1pp/2n5/3pPp2/8/5N2/PPPP1PPP/RNBQKB1R w KQkq f6 0 4",
			postFEN: "r1bqkbnr/ppp1p1pp/2n2P2/3p4/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 4",
		},
		{
			move:    newMove(BlackKnight, WhiteKnight, C6, D4, NoSquare, NoPiece),
			moveUCI: "c6d4",
			tags:    []MoveTag{Capture},
			preFEN:  "r1bqkbnr/ppp1p1pp/2n5/3pPp2/3N4/8/PPPP1PPP/RNBQKB1R b KQkq - 1 4",
			postFEN: "r1bqkbnr/ppp1p1pp/8/3pPp2/3n4/8/PPPP1PPP/RNBQKB1R w KQkq - 0 5",
		},
		{
			move:    newMove(WhitePawn, BlackRook, B7, A8, NoSquare, WhiteQueen),
			moveUCI: "b7a8q",
			tags:    []MoveTag{Promotion, Capture},
			preFEN:  "r7/1Pp5/2P3p1/8/6pb/4p1kB/4P1p1/6K1 w - - 0 1",
			postFEN: "Q7/2p5/2P3p1/8/6pb/4p1kB/4P1p1/6K1 b - - 0 1",
		},
	}
)

func TestNewBoard(t *testing.T) {
	assert.Equal(t, startingBoard, newBoard(startingSquareMap))
}

func TestBoard_MakeMoveBoard(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.move.String(), func(t *testing.T) {
			pos, post := unsafeFEN(tt.preFEN), unsafeFEN(tt.postFEN)
			pos.board.makeMove(tt.move, pos.castling.files)
			want := strings.Fields(tt.postFEN)[0]
			assert.Equal(t, want, pos.board.String())
			assert.Equal(t, post.board, pos.board)
		})
	}
}

func TestBoard_String(t *testing.T) {
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	assert.Equal(t, expected, startingBoard.String())
}

func BenchmarkBoard_PieceAt(b *testing.B) {
	squares := []Square{E1, E4, B7}
	pos := unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	for _, sq := range squares {
		b.Run(sq.String(), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				pos.board.pieceAt(sq)
			}
		})
	}
}

func BenchmarkBoard_PieceByColor(b *testing.B) {
	squares := []Square{E1, E4, B7}
	pos := unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	for _, sq := range squares {
		b.Run(sq.String(), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				pos.board.pieceByColor(sq, White)
			}
		})
	}
}

// unsafeFEN returns a position without error checking, only meant for tests.
func unsafeFEN(fen string) *Position {
	p, err := NewPosition(fen, FEN{})
	if err != nil {
		panic(err)
	}
	return p
}
