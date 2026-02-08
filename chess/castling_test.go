package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastling_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args castling
		want string
	}{
		{castling{[2]File{FileA, FileH}, 0}, "-"},
		{castling{[2]File{FileA, FileH}, castleWhiteH | castleWhiteA}, "KQ"},
		{castling{[2]File{FileA, FileH}, castleWhiteH | castleWhiteA | castleBlackH | castleBlackA}, "KQkq"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.args), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.args.String())
		})
	}
}

func TestNewCastleCheck(t *testing.T) {
	t.Parallel()
	type args struct {
		c     Color
		s     side
		kings [2]Square
		cf    [2]File
		cr    castlingRights
	}

	all := castleBlackA | castleBlackH | castleWhiteA | castleWhiteH

	tests := []struct {
		name string
		fen  string
		args
		want castleCheck
	}{
		{
			"Chess960 580",
			"qbb1rkrn/1ppppppp/p7/7n/8/P2P4/1PP1PPPP/QBBRNKRN w Gg - 0 9",
			args{White, hSide, [2]Square{F8, F1}, [2]File{FileD, FileG}, castleWhiteH},
			castleCheck{
				bbKingTravel:    bbEmpty,
				bbRookTravel:    bbEmpty,
				bbNoEnemyPawn:   61440,
				bbNoEnemyKnight: 15767552,
				bbNoEnemyKing:   61680,
				bbNoCheck:       F1.bitboard() | G1.bitboard(),
				king1:           F1,
				king2:           G1,
				rook1:           G1,
				rook2:           F1,
			},
		},
		{
			"Chess960 865",
			"bqkr1rnn/1ppp1ppp/p4b2/4p3/P7/3PP2N/1PP2PPP/BQRBKR1N w FC - 3 9",
			args{White, hSide, [2]Square{C8, E1}, [2]File{FileC, FileF}, castleWhiteA | castleWhiteH},
			castleCheck{
				bbKingTravel:    G1.bitboard(),
				bbRookTravel:    bbEmpty,
				bbNoEnemyPawn:   63488,
				bbNoEnemyKnight: 16309248,
				bbNoEnemyKing:   63736,
				bbNoCheck:       E1.bitboard() | F1.bitboard() | G1.bitboard(),
				king1:           E1,
				king2:           G1,
				rook1:           F1,
				rook2:           F1,
			},
		},
		{
			"Chess960 877",
			"qrk1rnb1/p1pp1ppp/1p2Bbn1/8/4P3/6P1/PPPP1P1P/QRK1RNBN w EBeb - 1 9",
			args{White, aSide, [2]Square{C8, C1}, [2]File{FileB, FileE}, all},
			castleCheck{
				bbKingTravel:    bbEmpty,
				bbRookTravel:    D1.bitboard(),
				bbNoEnemyPawn:   2560,
				bbNoEnemyKnight: 659712,
				bbNoEnemyKing:   3594,
				bbNoCheck:       C1.bitboard(),
				king1:           C1,
				king2:           C1,
				rook1:           B1,
				rook2:           D1,
			},
		},
		{
			"Chess960 944",
			"b2krn1q/p1rppppp/1Q3n2/2p1b3/1P4P1/8/P1PPPP1P/BBRKRNN1 w - - 3 9",
			args{White, aSide, [2]Square{D8, D1}, [2]File{FileC, FileE}, castleBlackH | castleWhiteA | castleWhiteH},
			castleCheck{
				bbKingTravel:    bbEmpty,
				bbRookTravel:    bbEmpty,
				bbNoEnemyPawn:   7680,
				bbNoEnemyKnight: 1979136,
				bbNoEnemyKing:   7710,
				bbNoCheck:       C1.bitboard() | D1.bitboard(),
				king1:           D1,
				king2:           C1,
				rook1:           C1,
				rook2:           D1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := newCastleCheck(tt.c, tt.s, tt.kings, tt.cf, tt.cr)
			assert.Equal(t, tt.want, got)
		})
	}
}
