package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastling_String(t *testing.T) {
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
			assert.Equal(t, tt.want, tt.args.String())
		})
	}
}
