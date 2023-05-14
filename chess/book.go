package chess

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

// Book holds the opening book data.
type Book struct {
	m map[Hash][]openingMove
}

// openingMove represent a single move openingMove.
type openingMove struct {
	s1, s2 Square
	promo  PieceType
	weight int
}

// WeightedMove is a single weighted move.
// The weight is positive but no bounds are enforced.
type WeightedMove struct {
	Move   Move
	Weight int
}

// NewBook returns a new empty book.
func NewBook() *Book {
	return &Book{
		m: make(map[Hash][]openingMove),
	}
}

// Init takes a reader to binary data (Polyglot files .bin) and initializes
// the opening book. It may be called several times with data from
// different books and data will be merged. Weights will be returned as is
// so it is up to the user to ensures they are consistent across books.
func (b *Book) Init(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanEntries)

	for scanner.Scan() {
		data := scanner.Bytes()
		hash := Hash(binary.BigEndian.Uint64(data[0:8]))
		move := binary.BigEndian.Uint16(data[8:10])
		weight := binary.BigEndian.Uint16(data[10:12])

		if moves, ok := b.m[hash]; ok {
			b.m[hash] = append(moves, parseRawMove(move, weight))
		} else {
			b.m[hash] = []openingMove{parseRawMove(move, weight)}
		}
	}

	return scanner.Err()
}

// Lookup takes a position and returns a sorted list of weighted moves.
// If the position is not found, nil is returned.
func (b *Book) Lookup(pos *Position) []WeightedMove {
	moves, ok := b.m[pos.hash]
	if !ok {
		return nil
	}

	weightedMoves := make([]WeightedMove, 0, len(moves))
	for _, move := range moves {
		s1, s2 := move.s1, move.s2
		p1 := pos.board.pieceAt(s1)
		s2 = castlingDestination(s1, s2, p1)
		p2 := pos.board.pieceAt(s2)

		weightedMoves = append(weightedMoves, WeightedMove{
			Move:   newMove(p1, p2, s1, s2, NoSquare, move.promo.color(p1.Color())),
			Weight: move.weight,
		})
	}
	return weightedMoves
}

// parseRawMove parses a raw move.
//
// A raw move is a bit field with the following meaning (bit 0 is the least
// significant bit):
//
//	===================================
//	0,1,2               to file
//	3,4,5               to row
//	6,7,8               from file
//	9,10,11             from row
//	12,13,14            promotion piece
func parseRawMove(move, weight uint16) openingMove {
	file2 := move & 7
	rank2 := (move >> 3) & 7
	file1 := (move >> 6) & 7
	rank1 := (move >> 9) & 7
	promo := (move >> 12) & 7

	return openingMove{
		s1:     Square(8*rank1 + file1),
		s2:     Square(8*rank2 + file2),
		promo:  promotionCodes[promo],
		weight: int(weight),
	}
}

// promotionCodes is an array of piece types indexed by promotion codes,
// determined as follow:
//
//	none       0
//	knight     1
//	bishop     2
//	rook       3
//	queen      4
var promotionCodes = [5]PieceType{NoPieceType, Knight, Bishop, Rook, Queen}

// castlingDestination returns a new destination square if the move is a castling.
//
// Castling moves are unconventionally represented as follow:
//
//	white short      e1h1
//	white long       e1a1
//	black short      e8h8
//	black long       e8a8
//
// If the move is not castling, the original destination square is returned.
func castlingDestination(from, to Square, fromPiece Piece) Square {
	switch {
	case from == E1 && to == H1 && fromPiece == WhiteKing:
		return G1
	case from == E1 && to == A1 && fromPiece == WhiteKing:
		return C1
	case from == E8 && to == H8 && fromPiece == WhiteKing:
		return G8
	case from == E8 && to == A8 && fromPiece == WhiteKing:
		return C8
	default:
		return to
	}
}

// scanEntries implements the bufio.SplitFunc interface.
//
// Returns the data by chunks of 16 bytes.
func scanEntries(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch {
	case atEOF && len(data) == 0:
		// expected EOF
		return 0, nil, nil
	case len(data) >= 16:
		// return raw entry bytes
		return 16, data[0:16], nil
	case atEOF:
		// if at EOF and we still have data, the data is malformed
		return len(data), data, errors.New("expected data to be multiple of 16 bytes")
	default:
		// request more data
		return 0, nil, nil
	}
}
