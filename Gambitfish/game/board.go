// Board is a representation of a chess board and the operations possible on it.
package game

import "fmt"

type Board struct {
	Squares  [64]Piece
	Active   Color
	Winner   Color
	PieceSet map[Piece]Square
}

type Move struct {
	piece  Piece
	square Square
}

func (m Move) String() string {
	return fmt.Sprintf("%v%v to %v", m.piece, m.piece.Square(), &m.square)
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
		b.Squares[blackPawnSquare.Index()] = &Pawn{&BasePiece{color: BLACK, board: b}}
		b.Squares[whitePawnSquare.Index()] = &Pawn{&BasePiece{color: WHITE, board: b}}
	}
	// Add rooks.
	b.Squares[0] = &Rook{&BasePiece{color: WHITE, board: b}, false}
	b.Squares[7] = &Rook{&BasePiece{color: WHITE, board: b}, false}
	b.Squares[56] = &Rook{&BasePiece{color: BLACK, board: b}, false}
	b.Squares[63] = &Rook{&BasePiece{color: BLACK, board: b}, false}
	// Add &Knights.
	b.Squares[1] = &Knight{&BasePiece{color: WHITE, board: b}}
	b.Squares[6] = &Knight{&BasePiece{color: WHITE, board: b}}
	b.Squares[57] = &Knight{&BasePiece{color: BLACK, board: b}}
	b.Squares[62] = &Knight{&BasePiece{color: BLACK, board: b}}
	// Add &Bishops.
	b.Squares[2] = &Bishop{&BasePiece{color: WHITE, board: b}}
	b.Squares[5] = &Bishop{&BasePiece{color: WHITE, board: b}}
	b.Squares[58] = &Bishop{&BasePiece{color: BLACK, board: b}}
	b.Squares[61] = &Bishop{&BasePiece{color: BLACK, board: b}}
	// Add queens
	b.Squares[3] = &Queen{&BasePiece{color: WHITE, board: b}}
	b.Squares[59] = &Queen{&BasePiece{color: BLACK, board: b}}
	// Add Kings
	b.Squares[4] = &King{&BasePiece{color: WHITE, board: b}, false}
	b.Squares[60] = &King{&BasePiece{color: BLACK, board: b}, false}
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
			newline += " "
			// Print a new line every 8 Squares, and reset.
			if ((i + 1) % 8) == 0 {
				fmt.Println(newline)
				newline = ""
			}
		}
	}
}

// ApplyMove changes the state of the board for any given move.
// TODO: This won't work for castling, we will need to special case that move.
func (b *Board) ApplyMove(m Move) {
	p := m.piece
	s := m.square
	// Check for victory
	if _, ok := b.Squares[s.Index()].(*King); ok {
		b.Winner = b.Active
	}
	// Reset the original square for this piece.
	oldS := b.PieceSet[p]
	b.Squares[oldS.Index()] = nil
	// Place the piece on the new squares.
	b.PieceSet[p] = s

	b.Squares[s.Index()] = p
	// Change the active player.
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
			moves = append(moves, Move{piece, square})
		}
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
