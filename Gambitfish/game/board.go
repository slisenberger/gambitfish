// Board is a representation of a chess board and the operations possible on it.
package game

import "fmt"

type Board struct {
	Squares [64]Piece
	active  Color
}

type Move struct {
	piece  Piece
	square Square
}

func DefaultBoard() *Board {
	b := &Board{active: WHITE}
	// Add pawns
	for i := 1; i <= 8; i++ {
		blackPawnSquare := &Square{7, i}
		whitePawnSquare := &Square{2, i}
		b.Squares[blackPawnSquare.Index()] = Pawn{BasePiece{color: BLACK, square: blackPawnSquare, board: b}}
		b.Squares[whitePawnSquare.Index()] = Pawn{BasePiece{color: WHITE, square: whitePawnSquare, board: b}}
	}
	// Add rooks.
	b.Squares[0] = Rook{BasePiece{color: WHITE, square: &Square{1, 1}, board: b}}
	b.Squares[7] = Rook{BasePiece{color: WHITE, square: &Square{1, 8}, board: b}}
	b.Squares[56] = Rook{BasePiece{color: BLACK, square: &Square{8, 1}, board: b}}
	b.Squares[63] = Rook{BasePiece{color: BLACK, square: &Square{8, 8}, board: b}}
	// Add Knights.
	b.Squares[1] = Knight{BasePiece{color: WHITE, square: &Square{1, 2}, board: b}}
	b.Squares[6] = Knight{BasePiece{color: WHITE, square: &Square{1, 7}, board: b}}
	b.Squares[57] = Knight{BasePiece{color: BLACK, square: &Square{8, 2}, board: b}}
	b.Squares[62] = Knight{BasePiece{color: BLACK, square: &Square{8, 7}, board: b}}
	// Add Bishops.
	b.Squares[2] = Bishop{BasePiece{color: WHITE, square: &Square{1, 3}, board: b}}
	b.Squares[5] = Bishop{BasePiece{color: WHITE, square: &Square{1, 6}, board: b}}
	b.Squares[58] = Bishop{BasePiece{color: BLACK, square: &Square{8, 3}, board: b}}
	b.Squares[61] = Bishop{BasePiece{color: BLACK, square: &Square{8, 6}, board: b}}
	// Add queens
	b.Squares[3] = Queen{BasePiece{color: WHITE, square: &Square{1, 4}, board: b}}
	b.Squares[59] = Queen{BasePiece{color: BLACK, square: &Square{8, 4}, board: b}}
	// Add Kings
	b.Squares[4] = King{BasePiece{color: WHITE, square: &Square{1, 5}, board: b}}
	b.Squares[60] = King{BasePiece{color: BLACK, square: &Square{8, 5}, board: b}}
	return b
}

func (b *Board) Print() {
	switch b.active {
	case WHITE:
		fmt.Println("White to play:")
	case BLACK:
		fmt.Println("Black to play:")
	}
	fmt.Println("")

	newline := ""
	// We want to print a row at a time, but in backwards order to how this is stored
	// Since a1 occurs at the bottom left of the board.
	for row := 7; row >= 0; row-- {
		// Get 8 consecutive Squares
		Squares := b.Squares[8*row : 8*row+8]
		for i, square := range Squares {
			if square == nil {
				newline += "-"
			} else {
				newline += square.String()
			}
			// Print a new line every 8 Squares, and reset.
			if ((i + 1) % 8) == 0 {
				fmt.Println(newline)
				newline = ""
			}
		}
	}
}

// ApplyMove changes the state of the board for any given move.
func (b *Board) ApplyMove(p Piece, s Square) {
	b.Squares[p.Square().Index()] = nil
	if b.Squares[s.Color() != p.board.Color()] {
		return
	}
	p.square = s
	switch b.active {
	case WHITE:
		b.active = BLACK
		break
	case BLACK:
		b.active = WHITE
		break
	}
}

func (b *Board) AllLegalMoves() []Move {
	var moves []Move
	for _, piece := range b.squares {
		if piece == nil {
			return
		}
		if piece.color != active {
			continue
		}
		moves := piece.LegalMoves()
		for _, move := range moves {
			moves := append(moves, Move{piece, piece.square})
		}

	}

}
