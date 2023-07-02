package chess

import (
	"fmt"
	"strconv"
	"strings"
)

// FEN is the Forsyth-Edwards Notation.
type FEN struct{}

// Encode encodes a Position into a FEN string.
//
// Implements the Notation interface.
func (FEN) Encode(pos *Position) string {
	sq := "-"
	if pos.enPassant != NoSquare {
		sq = pos.enPassant.String()
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		pos.board.String(),
		pos.turn.String(),
		pos.castling.String(),
		sq,
		pos.halfMoveClock,
		pos.fullMoves,
	)
}

// Decode decodes a FEN string into a Position.
//
// Implements the Notation interface.
func (FEN) Decode(s string) (*Position, error) {
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

	files, err := fenCastlingFiles(fields[2])
	if err != nil {
		return nil, err
	}

	rights, err := fenCastlingRights(fields[2])
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

// fenBoard parses the board from FEN.
func fenBoard(field string) (board, error) {
	rankFields := strings.Split(field, "/")
	if len(rankFields) != 8 {
		return board{}, fmt.Errorf("invalid fen board (%s)", field)
	}

	m := map[Square]Piece{}
	for i, rankField := range rankFields {
		fileMap, err := fenFileField(rankField)
		if err != nil {
			return board{}, err
		}
		for f, p := range fileMap {
			m[newSquare(f, Rank(7-i))] = p
		}
	}
	return newBoard(m), nil
}

// fenFileFiled parses a single file field from FEN.
func fenFileField(rankField string) (map[File]Piece, error) {
	m := map[File]Piece{}
	file := FileA
	for _, r := range rankField {
		switch {
		case 'A' <= r && r <= 'z':
			m[file] = pieceTable[r-'A']
			file++
		case '1' <= r && r <= '8':
			file += File(r - '0')
		default:
			return nil, fmt.Errorf("invalid fen rank field (%s)", rankField)
		}
	}

	if file != FileH+1 {
		return nil, fmt.Errorf("invalid fen rank field (%s)", rankField)
	}
	return m, nil
}

// fenTurn parses the turn from FEN.
func fenTurn(field string) (Color, error) {
	switch field {
	case "w":
		return White, nil
	case "b":
		return Black, nil
	default:
		return White, fmt.Errorf("invalid fen turn (%s)", field)
	}
}

// fenCastlingFiles parses the castling files from FEN.
func fenCastlingFiles(_ string) ([2]File, error) {
	return [2]File{FileA, FileH}, nil
}

// fenCastlingRights parses the castling rights from FEN.
func fenCastlingRights(field string) (castlingRights, error) {
	if field == "-" {
		return noCastle, nil
	}

	var cr castlingRights
	for _, r := range field {
		switch r {
		case 'K':
			cr |= castleWhiteH
		case 'Q':
			cr |= castleWhiteA
		case 'k':
			cr |= castleBlackH
		case 'q':
			cr |= castleBlackA
		default:
			return noCastle, fmt.Errorf("invalid fen castling rights (%s)", field)
		}
	}
	return cr, nil
}

// fenEnPassantSquare parses the en passant square from FEN.
func fenEnPassantSquare(field string) (Square, error) {
	if field == "-" {
		return NoSquare, nil
	}
	if len(field) != 2 {
		return NoSquare, fmt.Errorf("invalid fen en passant square (%s)", field)
	}
	sq, err := uciSquare(field)
	if err != nil {
		return NoSquare, err
	}
	if sq == NoSquare || !(sq.Rank() == Rank3 || sq.Rank() == Rank6) {
		return NoSquare, fmt.Errorf("invalid fen en passant square (%s)", field)
	}
	return sq, nil
}

// fenHalfMoveClock parses the half move clock from FEN.
func fenHalfMoveClock(field string) (uint8, error) {
	halfMoveClock, err := strconv.ParseUint(field, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid fen full moves count (%s)", field)
	}
	return uint8(halfMoveClock), nil
}

// fenFullMoves parses the full moves count from FEN.
func fenFullMoves(field string) (uint8, error) {
	fullMoves, err := strconv.ParseUint(field, 10, 8)
	if err != nil || fullMoves < 1 {
		return 0, fmt.Errorf("invalid fen full moves count (%s)", field)
	}
	return uint8(fullMoves), nil
}
