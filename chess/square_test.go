package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquare_File(t *testing.T) {
	tests := []struct {
		name string
		args []Square
		want File
	}{
		{"A", []Square{A1, A2, A3, A4, A5, A6, A7, A8}, FileA},
		{"B", []Square{B1, B2, B3, B4, B5, B6, B7, B8}, FileB},
		{"C", []Square{C1, C2, C3, C4, C5, C6, C7, C8}, FileC},
		{"D", []Square{D1, D2, D3, D4, D5, D6, D7, D8}, FileD},
		{"E", []Square{E1, E2, E3, E4, E5, E6, E7, E8}, FileE},
		{"F", []Square{F1, F2, F3, F4, F5, F6, F7, F8}, FileF},
		{"G", []Square{G1, G2, G3, G4, G5, G6, G7, G8}, FileG},
		{"H", []Square{H1, H2, H3, H4, H5, H6, H7, H8}, FileH},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, sq := range tt.args {
				assert.Equal(t, tt.want, sq.File())
			}
		})
	}
}

func TestSquare_Rank(t *testing.T) {
	tests := []struct {
		name string
		args []Square
		want Rank
	}{
		{"1", []Square{A1, B1, C1, D1, E1, F1, G1, H1}, Rank1},
		{"2", []Square{A2, B2, C2, D2, E2, F2, G2, H2}, Rank2},
		{"3", []Square{A3, B3, C3, D3, E3, F3, G3, H3}, Rank3},
		{"4", []Square{A4, B4, C4, D4, E4, F4, G4, H4}, Rank4},
		{"5", []Square{A5, B5, C5, D5, E5, F5, G5, H5}, Rank5},
		{"6", []Square{A6, B6, C6, D6, E6, F6, G6, H6}, Rank6},
		{"7", []Square{A7, B7, C7, D7, E7, F7, G7, H7}, Rank7},
		{"8", []Square{A8, B8, C8, D8, E8, F8, G8, H8}, Rank8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, sq := range tt.args {
				assert.Equal(t, tt.want, sq.Rank())
			}
		})
	}
}
