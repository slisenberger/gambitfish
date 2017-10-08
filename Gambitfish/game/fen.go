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
	for _, char := range split[2] {
		switch char {
		case 'K':
			b.WKSCastling = true
		case 'k':
			b.BKSCastling = true
		case 'Q':
			b.WQSCastling = true
		case 'q':
			b.BQSCastling = true
		}
	}

	// TODO(slisenberger): Add En-passant square
	b.EPSquare = OFFBOARD_SQUARE
	// TODO(slisenberger): add move count!
	for s, p := range b.Squares {
		if p != NULLPIECE {
			b.Position = SetPiece(b.Position, p, Square(s))
		}
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
			b.Squares[square] = BLACKPAWN
		case 'P':
			b.Squares[square] = WHITEPAWN
		case 'b':
			b.Squares[square] = BLACKBISHOP
		case 'B':
			b.Squares[square] = WHITEBISHOP
		case 'n':
			b.Squares[square] = BLACKKNIGHT
		case 'N':
			b.Squares[square] = WHITEKNIGHT
		case 'q':
			b.Squares[square] = BLACKQUEEN
		case 'Q':
			b.Squares[square] = WHITEQUEEN
		case 'k':
			b.Squares[square] = BLACKKING
		case 'K':
			b.Squares[square] = WHITEKING
		case 'r':
			b.Squares[square] = BLACKROOK
		case 'R':
			b.Squares[square] = WHITEROOK
		default:
			return fmt.Errorf("fen notation has unrecognized char: %v", string(char))
		}

		if colNum > 8 {
			return fmt.Errorf("fen board contained row with >8 chars: %v", row)
		}
	}
	return nil
}
