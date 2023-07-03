package chess

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// ShredderFEN is the Shredder Forsyth-Edwards Notation.
//
// Compatible with Chess960.
type ShredderFEN struct{}

// Encode encodes a Position into a Shredder-FEN string.
//
// Implements the Notation interface.
func (ShredderFEN) Encode(pos *Position) string {
	sq := "-"
	if pos.enPassant != NoSquare {
		sq = pos.enPassant.String()
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		pos.board.String(),
		pos.turn.String(),
		shredderFenCastling(pos.castling),
		sq,
		pos.halfMoveClock,
		pos.fullMoves,
	)
}

// Decode decodes a Shredder-FEN string into a Position.
//
// Implements the Notation interface.
func (ShredderFEN) Decode(s string) (*Position, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid fen (%s), must have 6 fields", s)
	}

	var err error
	pos := &Position{}

	pos.board, err = fenBoard(fields[0])
	if err != nil {
		return nil, err
	}

	pos.turn, err = fenTurn(fields[1])
	if err != nil {
		return nil, err
	}

	pos.castling, err = shredderFenCastlingField(fields[2], pos.board.sqKings)
	if err != nil {
		return nil, err
	}

	for c := Black; c <= White; c++ {
		for s := aSide; s <= hSide; s++ {
			pos.castleChecks[2*uint8(c)+uint8(s)] = newCastleCheck(c, s, pos.board.sqKings, pos.castling.files, pos.castling.rights)
		}
	}

	pos.enPassant, err = fenEnPassantSquare(fields[3])
	if err != nil {
		return nil, err
	}

	pos.halfMoveClock, err = fenHalfMoveClock(fields[4])
	if err != nil {
		return nil, err
	}

	pos.fullMoves, err = fenFullMoves(fields[5])
	if err != nil {
		return nil, err
	}

	pos.hash = newZobristHash(pos)
	pos.pawnHash = newPawnZobristHash(pos)

	return pos, nil
}

// shredderFenCastling formats castling to string
// and follows the Shredder-FEN notation.
func shredderFenCastling(c castling) string {
	if c.rights == noCastle {
		return "-"
	}

	var rights string
	if c.rights.canCastle(White, hSide) {
		rights += strings.ToUpper(c.files[hSide].String())
	}
	if c.rights.canCastle(White, aSide) {
		rights += strings.ToUpper(c.files[aSide].String())
	}
	if c.rights.canCastle(Black, hSide) {
		rights += c.files[hSide].String()
	}
	if c.rights.canCastle(Black, aSide) {
		rights += c.files[aSide].String()
	}
	return rights
}

// shredderFenCastlingField parses a castling field that
// follows the Shredder-FEN notation.
func shredderFenCastlingField(field string, kings [2]Square) (castling, error) {
	if field == "-" {
		return castling{}, nil
	}

	files := [2]File{FileA, FileH}
	rights := noCastle

	runes := []rune(field)
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	for _, r := range runes {
		lr := unicode.ToLower(r)
		if lr < 'a' || 'h' < lr {
			return castling{}, fmt.Errorf("invalid fen castling field (%s)", field)
		}

		file := File(lr - 'a')
		c, s := Black, aSide
		if r == unicode.ToUpper(r) {
			c = White
		}
		if File(kings[c].File().String()[0]-'a') < file {
			s = hSide
		}

		rights |= castlingRightsMap[2*uint8(c)+uint8(s)]
		files[s] = file
	}

	return castling{files: files, rights: rights}, nil
}

var castlingRightsMap = [4]castlingRights{castleBlackA, castleBlackH, castleWhiteA, castleWhiteH} // indexed by 2*Color+side.
