// Board is a representation of a chess Board and the operations possible on it.
package game

import "fmt"
import "math/rand"

type Board struct {
	Squares  [64]Piece
	Active   Color
	Winner   Color
	PieceSet map[Piece]Square
}

type Move struct {
	Piece  Piece
	Square Square
	Old    Square
}

func (m Move) String() string {
	return fmt.Sprintf("%v%v to %v", m.Piece, m.Old, m.Square)
}

func (b *Board) InitPieceSet() {
	b.PieceSet = make(map[Piece]Square, 32)
	for i, piece := range b.Squares {
		if piece != nil {
			b.PieceSet[piece] = SquareFromIndex(i)
		}
	}
}

func DefaultBoard() *Board {
	b := &Board{Active: WHITE}
	// Add pawns
	for i := 1; i <= 8; i++ {
		blackPawnSquare := &Square{7, i}
		whitePawnSquare := &Square{2, i}
		b.Squares[blackPawnSquare.Index()] = &Pawn{&BasePiece{C: BLACK, B: b}}
		b.Squares[whitePawnSquare.Index()] = &Pawn{&BasePiece{C: WHITE, B: b}}
	}
	// Add rooks.
	b.Squares[0] = &Rook{&BasePiece{C: WHITE, B: b}, false}
	b.Squares[7] = &Rook{&BasePiece{C: WHITE, B: b}, false}
	b.Squares[56] = &Rook{&BasePiece{C: BLACK, B: b}, false}
	b.Squares[63] = &Rook{&BasePiece{C: BLACK, B: b}, false}
	// Add &Knights.
	b.Squares[1] = &Knight{&BasePiece{C: WHITE, B: b}}
	b.Squares[6] = &Knight{&BasePiece{C: WHITE, B: b}}
	b.Squares[57] = &Knight{&BasePiece{C: BLACK, B: b}}
	b.Squares[62] = &Knight{&BasePiece{C: BLACK, B: b}}
	// Add &Bishops.
	b.Squares[2] = &Bishop{&BasePiece{C: WHITE, B: b}}
	b.Squares[5] = &Bishop{&BasePiece{C: WHITE, B: b}}
	b.Squares[58] = &Bishop{&BasePiece{C: BLACK, B: b}}
	b.Squares[61] = &Bishop{&BasePiece{C: BLACK, B: b}}
	// Add queens
	b.Squares[3] = &Queen{&BasePiece{C: WHITE, B: b}}
	b.Squares[59] = &Queen{&BasePiece{C: BLACK, B: b}}
	// Add Kings
	b.Squares[4] = &King{&BasePiece{C: WHITE, B: b}, false}
	b.Squares[60] = &King{&BasePiece{C: BLACK, B: b}, false}
	b.InitPieceSet()
	return b
}

func (b *Board) Print() {
	switch b.Active {
	case WHITE:
		fmt.Println("White to play:")
	case BLACK:
		fmt.Println("Black to play:")
	}
	fmt.Println("")

	// We want to print a row at a time, but in backwards order to how this is stored
	// Since a1 occurs at the bottom left of the Board.
	for row := 7; row >= 0; row-- {
		newline := fmt.Sprintf("%v | ", row+1)
		// Get 8 consecutive Squares
		Squares := b.Squares[8*row : 8*row+8]
		for i, piece := range Squares {
			if piece == nil {
				newline += "·"
			} else {
				newline += piece.Graphic()
			}
			newline += " "
			// Print a new line every 8 Squares, and reset.
			if ((i + 1) % 8) == 0 {
				fmt.Println(newline)
				newline = ""
			}
		}
	}
	// Print the letters at the bottom.
	fmt.Println("   ________________")
	fmt.Println("    a b c d e f g h")

}

// ApplyMove changes the state of the Board for any given move.
// TODO: This won't work for castling, we will need to special case that move.
func (b *Board) ApplyMove(m Move) {
	p := m.Piece
	s := m.Square
	// Check for victory
	// TODO(slisenberger): this needs some serious work.
	// Currently checks if we are capturing a king, and if so, loser is
	// opposite Color.
	if king, ok := b.Squares[s.Index()].(*King); ok {
		b.Winner = -1 * king.Color()
	}
	b.Squares[m.Old.Index()] = nil
	// Place the piece on the new squares.
	b.PieceSet[p] = s

	b.Squares[s.Index()] = p
}

func (b *Board) SwitchActivePlayer() {
	switch b.Active {
	case WHITE:
		b.Active = BLACK
		break
	case BLACK:
		b.Active = WHITE
		break
	}
}

// AllLegalMoves enumerates all of the legal moves currently available to the
// active player.
func (b *Board) AllLegalMoves() []Move {
	var moves []Move
	for _, piece := range b.Squares {
		if piece == nil {
			continue
		}
		if piece.Color() != b.Active {
			continue
		}
		squares := piece.LegalMoves()
		for _, square := range squares {
			moves = append(moves, Move{piece, square, b.PieceSet[piece]})
		}
	}
	// We now want to return the moves in a smart order (e.g., try checks and captures.)
	// Just to get things off the ground, we'll shuffle the moves, just to get some variety
	// in the AI vs AI games.
	for i := range moves {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	return moves
}

// Returns true if the game is drawn, won, or lost.
func (b *Board) Finished() bool {
	if len(b.AllLegalMoves()) == 0 {
		return true
	}
	return b.Winner != 0
}

// Returns a copy of this Board with different references.
func (b *Board) Copy() *Board {
	newB := &Board{Squares: b.Squares, Active: b.Active, Winner: b.Winner}
	newB.InitPieceSet()
	return newB
}
