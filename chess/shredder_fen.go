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

	files, err := shredderFenCastlingFiles(fields[2])
	if err != nil {
		return nil, err
	}

	rights, err := shredderFenCastlingRights(fields[2], files)
	if err != nil {
		return nil, err
	}

	pos.castling = castling{files, rights}

	for c := Black; c <= White; c++ {
		for s := aSide; s <= hSide; s++ {
			pos.castleChecks[2*uint8(c)+uint8(s)] = newCastleCheck(c, s, pos.board.sqKings, files, rights)
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

// shredderFenCastlingFiles parses the castling files from FEN.
func shredderFenCastlingFiles(field string) ([2]File, error) {
	if field == "-" {
		return [2]File{FileA, FileH}, nil
	}

	set := make(map[rune]struct{})
	runes := []rune{}
	for _, r := range strings.ToLower(field) {
		if _, ok := set[r]; !ok {
			set[r] = struct{}{}
			runes = append(runes, r)
		}
	}

	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })

	switch len(runes) {
	case 0:
		return [2]File{FileA, FileH}, nil
	case 1:
		return [2]File{File(runes[0] - 'a'), FileH}, nil
	case 2:
		return [2]File{File(runes[0] - 'a'), File(runes[1] - 'a')}, nil
	default:
		return [2]File{FileA, FileH}, fmt.Errorf("invalid fen castling files (%s)", field)
	}
}

// shredderFenCastlingRights parses the castling rights from FEN.
func shredderFenCastlingRights(field string, files [2]File) (castlingRights, error) {
	if field == "-" {
		return noCastle, nil
	}

	aSideFile := []rune(files[aSide].String())[0]
	hSideFile := []rune(files[hSide].String())[0]

	var cr castlingRights
	for _, r := range field {
		switch r {
		case aSideFile:
			cr |= castleBlackA
		case hSideFile:
			cr |= castleBlackH
		case unicode.ToUpper(aSideFile):
			cr |= castleWhiteA
		case unicode.ToUpper(hSideFile):
			cr |= castleWhiteH
		default:
			return noCastle, fmt.Errorf("invalid fen castling rights (%s)", field)
		}
	}
	return cr, nil
}
