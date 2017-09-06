// Piepce is an interface that defines the operations possible for a piece on the board.
package game

// Define the possible colors of a piece as an enum
type Color int

const (
	WHITE Color = 0
	BLACK Color = 1
)

type Piece interface {
	// Returns an array of all the legal positions this piece can move to.
	LegalMoves() []Square
	// Returns a string representation of this piece.
	String() string
	// Returns the square for this piece.
	Square() *Square
	// Returns the color of this piece.
	Color() Color
}

type BasePiece struct {
	color  Color
	square *Square
	board  *Board
}

func (bp BasePiece) Square() *Square {
	return bp.square
}

func (bp BasePiece) Color() Color {
	return bp.color
}

// TargetLegal returns true if a candidate piece can move to the
// desired square. Returns a
func (p BasePiece) TargetLegal(s Square) bool {
	if !s.InPlay() {
		return false
	}
	occupant := p.board.Squares[s.Index()]
	if occupant == nil {
		return true
	} else {
		if p.Color() != occupant.Color() {
			return true
		}
	}
	return false
}

// Returns true if we should stop searching along a ray, because a piece has
// encountered a blockage.
func (p BasePiece) Stop(s Square) bool {
	if !s.InPlay() {
		return true
	}
	if occupant := p.board.Squares[s.Index()]; occupant != nil {
		return true
	}
	return false
}

func (p BasePiece) ColumnAndRowMoves() []Square {
	moves := []Square{}
	// Move in each direction, checking for blocking pieces.
	// Left in row.
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row, col: p.square.col - i}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row, col: p.square.col + i}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row + i, col: p.square.col}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row - i, col: p.square.col}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	return moves
}

func (p BasePiece) DiagonalMoves() []Square {
	moves := []Square{}
	// Move in a diagonal, checking for blocking pieces.
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row + i, col: p.square.col + i}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row - i, col: p.square.col + i}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row - i, col: p.square.col - i}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: p.square.row + i, col: p.square.col - i}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	return moves
}

func (p BasePiece) KnightMoves() []Square {
	moves := []Square{}
	// try all knight squares
	squares := []Square{
		{row: p.square.row + 2, col: p.square.col + 1},
		{row: p.square.row + 2, col: p.square.col - 1},
		{row: p.square.row - 2, col: p.square.col + 1},
		{row: p.square.row - 2, col: p.square.col - 1},
		{row: p.square.row + 1, col: p.square.col + 2},
		{row: p.square.row + 1, col: p.square.col - 2},
		{row: p.square.row - 1, col: p.square.col + 2},
		{row: p.square.row - 1, col: p.square.col - 2},
	}
	for _, s := range squares {
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
	}
	return moves
}

func (p BasePiece) KingMoves() []Square {
	moves := []Square{}
	// try all knight squares
	squares := []Square{
		{row: p.square.row + 1, col: p.square.col - 1},
		{row: p.square.row + 1, col: p.square.col},
		{row: p.square.row + 1, col: p.square.col + 1},
		{row: p.square.row + -1, col: p.square.col - 1},
		{row: p.square.row + -1, col: p.square.col},
		{row: p.square.row + -1, col: p.square.col + 1},
		{row: p.square.row, col: p.square.col - 1},
		{row: p.square.row, col: p.square.col + 1},
	}
	for _, s := range squares {
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
	}
	return moves

}

func (p BasePiece) PawnMoves() []Square {
	// Check if the piece can move two squares.
	var isStartPawn bool
	switch p.color {
	case BLACK:
		isStartPawn = p.square.row == 7
		break
	case WHITE:
		isStartPawn = p.square.row == 2
		break
	}
	moves := []Square{}
	s := Square{row: p.square.row + 1, col: p.square.col}
	if p.TargetLegal(s) {
		moves = append(moves, s)
	}
	if isStartPawn {
		s := Square{row: p.square.row + 2, col: p.square.col}
		if p.TargetLegal(s) {
			moves = append(moves, s)
		}
		//TODO(slisenberger): build in en passant.

	}
	return moves
}
