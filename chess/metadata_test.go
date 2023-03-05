package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastlingRight_String(t *testing.T) {
	tests := []struct {
		args castlingRights
		want string
	}{
		{0, "-"},
		{castleWhiteKing | castleWhiteQueen, "KQ"},
		{castleWhiteKing | castleWhiteQueen | castleBlackKing | castleBlackQueen, "KQkq"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.args), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.String())
		})
	}
}
