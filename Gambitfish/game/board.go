// Board is a representation of a chess board and the operations possible on it.
package game

import "fmt"

type Board struct {
	squares [64]Piece
	active  Color
}

func DefaultBoard() *Board {
	b := &Board{active: WHITE}
	// Add pawns
	for i := 1; i <= 8; i++ {
		blackPawnSquare := &Square{7, i}
		whitePawnSquare := &Square{2, i}
		b.squares[blackPawnSquare.Index()] = Pawn{BasePiece{color: BLACK}}
		b.squares[whitePawnSquare.Index()] = Pawn{BasePiece{color: WHITE}}
	}
	// Add rooks.
	b.squares[0] = Rook{BasePiece{color: WHITE}}
	b.squares[7] = Rook{BasePiece{color: WHITE}}
	b.squares[56] = Rook{BasePiece{color: BLACK}}
	b.squares[63] = Rook{BasePiece{color: BLACK}}
	// Add Knights.
	b.squares[1] = Knight{BasePiece{color: WHITE}}
	b.squares[6] = Knight{BasePiece{color: WHITE}}
	b.squares[57] = Knight{BasePiece{color: BLACK}}
	b.squares[62] = Knight{BasePiece{color: BLACK}}
	// Add Bishops.
	b.squares[2] = Bishop{BasePiece{color: WHITE}}
	b.squares[5] = Bishop{BasePiece{color: WHITE}}
	b.squares[58] = Bishop{BasePiece{color: BLACK}}
	b.squares[61] = Bishop{BasePiece{color: BLACK}}
	// Add queens
	b.squares[3] = Queen{BasePiece{color: WHITE}}
	b.squares[59] = Queen{BasePiece{color: BLACK}}
	// Add Kings
	b.squares[4] = King{BasePiece{color: WHITE}}
	b.squares[60] = King{BasePiece{color: BLACK}}
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
		// Get 8 consecutive squares
		squares := b.squares[8*row : 8*row+8]
		for i, square := range squares {
			if square == nil {
				newline += "-"
			} else {
				newline += square.String()
			}
			// Print a new line every 8 squares, and reset.
			if ((i + 1) % 8) == 0 {
				fmt.Println(newline)
				newline = ""
			}
		}
	}
}
