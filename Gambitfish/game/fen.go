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
	b.EPSquare = OFFBOARD_SQUARE
	// TODO(slisenberger): add move count!
	// Basic board initialization.
	b.InitPieceSet()
	for p, s := range b.PieceSet {
		b.Position = SetPiece(b.Position, p, s)
	}
	b.Position = UpdateBitboards(b.Position)
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
		square := GetSquare(rowNum, colNum)
		switch char {
		case 'p':
			b.Squares[square] = &Pawn{&BasePiece{C: BLACK}}
		case 'P':
			b.Squares[square] = &Pawn{&BasePiece{C: WHITE}}
		case 'b':
			b.Squares[square] = &Bishop{&BasePiece{C: BLACK}}
		case 'B':
			b.Squares[square] = &Bishop{&BasePiece{C: WHITE}}
		case 'n':
			b.Squares[square] = &Knight{&BasePiece{C: BLACK}}
		case 'N':
			b.Squares[square] = &Knight{&BasePiece{C: WHITE}}
		case 'q':
			b.Squares[square] = &Queen{&BasePiece{C: BLACK}}
		case 'Q':
			b.Squares[square] = &Queen{&BasePiece{C: WHITE}}
		case 'k':
			b.Squares[square] = &King{&BasePiece{C: BLACK}}
		case 'K':
			b.Squares[square] = &King{&BasePiece{C: WHITE}}
		case 'r':
			ks := false
			qs := false
			if square.Col() == 1 {
				qs = true
			}
			if square.Col() == 8 {
				ks = true
			}
			b.Squares[square] = &Rook{&BasePiece{C: BLACK}, qs, ks}
		case 'R':
			ks := false
			qs := false
			if square.Col() == 1 {
				qs = true
			}
			if square.Col() == 8 {
				ks = true
			}
			b.Squares[square] = &Rook{&BasePiece{C: WHITE}, qs, ks}
		default:
			return fmt.Errorf("fen notation has unrecognized char: %v", string(char))
		}

		if colNum > 8 {
			return fmt.Errorf("fen board contained row with >8 chars: %v", row)
		}
	}
	return nil
}
