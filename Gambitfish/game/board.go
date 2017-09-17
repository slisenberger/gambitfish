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
func ApplyMove(b *Board, m Move) {
	p := m.Piece
	s := m.Square
	// If there's a capture: remove the captured piece.
	if m.Capture != nil {
		delete(b.PieceSet, m.Capture.Piece)
	}
	// Then, move the piece to its new square.
	// Check for promotion of a pawn.
	if m.Promotion != nil {
		b.PieceSet[m.Promotion] = s
		b.Squares[s.Index()] = m.Promotion
		delete(b.PieceSet, p)
	} else {
		b.PieceSet[p] = s
		b.Squares[s.Index()] = p
	}
	// Then, remove the piece from its old square.
	b.Squares[m.Old.Index()] = nil
}

// UndoMove returns a board to the state it was at prior to
// applying a move m.
func UndoMove(b *Board, m Move) {
	p := m.Piece
	o := m.Old
	s := m.Square
	// Return the piece to its old square, and undo promotion.
	if m.Promotion != nil {
		delete(b.PieceSet, m.Promotion)
	}

	b.Squares[o.Index()] = p
	b.PieceSet[p] = o
	// Return the square we were on to its old state.
	b.Squares[s.Index()] = nil
	// Return a captured piece.
	if m.Capture != nil {
		b.PieceSet[m.Capture.Piece] = m.Capture.Square
		b.Squares[m.Capture.Square.Index()] = m.Capture.Piece
	}
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
		m := piece.LegalMoves()
		for _, move := range m {
			// Remove pseudolegal moves
			ApplyMove(b, move)
			isCheck := IsCheck(b, b.Active)
			UndoMove(b, move)
			if isCheck {
				fmt.Println(fmt.Sprintf("not making a move %v: puts king in check!", move))

			} else {
				moves = append(moves, move)
			}
		}
	}
	// We now want to return the moves in a smart order (e.g., try checks and captures.)
	// Just to get things off the ground, we'll shuffle the moves, just to get some variety
	// in the AI vs AI games.
	for i := range moves {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	return OrderMoves(moves)
}

// Returns true if the game is drawn, won, or lost. The integer
// returned is the winner (or draw). This function should be called on the beginning
// of a move.
func (b *Board) CalculateGameOver() (bool, Color) {
	if len(b.AllLegalMoves()) == 0 {
		// If the active player is in check, they lose.
		if IsCheck(b, b.Active) {
			return true, -1 * b.Active
		} else {
			// Otherwise, it's stalemate!
			return true, 0
		}

	}
	return false, 0
}

// Returns a copy of this Board with different references.
func CopyBoard(b *Board) *Board {
	var s [64]Piece
	copy(s[:], b.Squares[:])
	newB := &Board{Squares: s, Active: b.Active, Winner: b.Winner}
	newB.InitPieceSet()
	return newB
}

// Returns true if the board state results in the color c's king being in check.
func IsCheck(b *Board, c Color) bool {
	for piece, _ := range b.PieceSet {
		// If it's our piece, we don't care.
		if piece.Color() == c {
			continue
		}
		attacking := piece.Attacking()
		for _, s := range attacking {
			occupant := b.Squares[s.Index()]
			// check if it's us under attack
			if occupant != nil && occupant.Color() == c {
				// If our king is under attack, it's check.
				if _, ok := occupant.(*King); ok {
					fmt.Println(fmt.Sprintf("found king on square %v", s.Index()))
					return true
				}
			}
		}
	}
	return false
}
