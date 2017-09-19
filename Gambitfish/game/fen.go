package game

import "fmt"
import "strings"

// BoardFromFen returns a new board object created from
// Forsyth edwards notation.
// https://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation
func BoardFromFen(s string) (*Board, error) {
	b := &Board{}
	split := strings.Split(s, " ")
	if len(split) != 6 {
		return nil, fmt.Errorf("fen string didn't contain 6 parts: %v", s)
	}
	// Parse piece config
	rows := strings.Split(split[0], "/")
	if len(rows) != 8 {
		return nil, fmt.Errorf("fen board contained wrong number of rows: %v", s)
	}
	// Add pieces to board.
	// FEN is given 8th rank to 1st, so we iterate in reverse.
	for i := 0; i < 8; i++ {
		row := rows[i]
		rowNum := 8 - i
		if err := HandleFenBoardRow(row, b, rowNum); err != nil {
			return nil, err
		}
	}

	// Add board color
	switch split[1] {
	case "w":
		b.Active = WHITE
	case "b":
		b.Active = BLACK
	default:
		return nil, fmt.Errorf("invalid board color in fen: %v", split[1])
	}

	// Add castling rights
	b.ksCastlingRights = map[Color]bool{WHITE: false, BLACK: false}
	b.qsCastlingRights = map[Color]bool{WHITE: false, BLACK: false}
	for _, char := range split[2] {
		switch char {
		case 'K':
			b.ksCastlingRights[WHITE] = true
		case 'k':
			b.ksCastlingRights[BLACK] = true
		case 'Q':
			b.qsCastlingRights[WHITE] = true
		case 'q':
			b.qsCastlingRights[BLACK] = true
		}
	}

	// TODO(slisenberger): Add En-passant square
	// TODO(slisenberger): add move count!
	b.InitPieceSet()
	return b, nil
}

// HandleFenBoardRow assigns the proper pieces to a given board from a
// row in fen notation.
func HandleFenBoardRow(row string, b *Board, rowNum int) error {
	colNum := 0
	for _, char := range row {
		// See if this represents an int.
		if j := int(char - '0'); j > 0 && j <= 8 {
			colNum += j
			continue
		}
		colNum += 1
		square := Square{rowNum, colNum}
		switch char {
		case 'p':
			b.Squares[square.Index()] = &Pawn{&BasePiece{C: BLACK, B: b}}
		case 'P':
			b.Squares[square.Index()] = &Pawn{&BasePiece{C: WHITE, B: b}}
		case 'b':
			b.Squares[square.Index()] = &Bishop{&BasePiece{C: BLACK, B: b}}
		case 'B':
			b.Squares[square.Index()] = &Bishop{&BasePiece{C: WHITE, B: b}}
		case 'n':
			b.Squares[square.Index()] = &Knight{&BasePiece{C: BLACK, B: b}}
		case 'N':
			b.Squares[square.Index()] = &Knight{&BasePiece{C: WHITE, B: b}}
		case 'q':
			b.Squares[square.Index()] = &Queen{&BasePiece{C: BLACK, B: b}}
		case 'Q':
			b.Squares[square.Index()] = &Queen{&BasePiece{C: WHITE, B: b}}
		case 'k':
			b.Squares[square.Index()] = &King{&BasePiece{C: BLACK, B: b}}
		case 'K':
			b.Squares[square.Index()] = &King{&BasePiece{C: WHITE, B: b}}
		case 'r':
			ks := false
			qs := false
			if square.Col == 1 {
				qs = true
			}
			if square.Col == 8 {
				ks = true
			}
			b.Squares[square.Index()] = &Rook{&BasePiece{C: BLACK, B: b}, qs, ks}
		case 'R':
			ks := false
			qs := false
			if square.Col == 1 {
				qs = true
			}
			if square.Col == 8 {
				ks = true
			}
			b.Squares[square.Index()] = &Rook{&BasePiece{C: WHITE, B: b}, qs, ks}
		default:
			return fmt.Errorf("fen notation has unrecognized char: %v", string(char))
		}

		if colNum > 8 {
			return fmt.Errorf("fen board contained row with >8 chars: %v", row)
		}
	}
	return nil
}
