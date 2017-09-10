// Piepce is an interface that defines the operations possible for a piece on the board.
package game

// Define the possible colors of a piece as an enum
type Color int

const (
	WHITE Color = -1
	BLACK Color = 1
)

type Piece interface {
	// Returns an array of all the legal positions this piece can move to.
	LegalMoves() []Square
	// Returns a string representation of this piece.
	String() string
	// Returns a unicode graphic representation of this piece.
	Graphic() string
	// Returns the color of this piece.
	Color() Color
	Board() *Board
	ApplyMove(Move)
	Value() float64
}

type BasePiece struct {
	color Color
	board *Board
}

func (bp *BasePiece) Color() Color {
	return bp.color
}

// TargetLegal returns true if a candidate piece can move to the
// desired square. Capture indicates whether the piece is intending to capture or not.
func (p *BasePiece) TargetLegal(s Square, capture bool) bool {
	if !s.InPlay() {
		return false
	}
	occupant := p.board.Squares[s.Index()]
	if occupant == nil {
		return true
	} else {
		if capture && p.Color() != occupant.Color() {
			return true
		}
	}
	return false
}

// Returns true if we should stop searching along a ray, because a piece has
// encountered a blockage.
func (p *BasePiece) Stop(s Square) bool {
	if !s.InPlay() {
		return true
	}
	if occupant := p.board.Squares[s.Index()]; occupant != nil {
		return true
	}
	return false
}

func (p *BasePiece) ColumnAndRowMoves(cur Square) []Square {
	moves := []Square{}
	// Move in each direction, checking for blocking pieces.
	// Left in row.
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row, col: cur.col - i}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row, col: cur.col + i}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row + i, col: cur.col}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row - i, col: cur.col}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	return moves
}

func (p *BasePiece) DiagonalMoves(cur Square) []Square {
	moves := []Square{}
	// Move in a diagonal, checking for blocking pieces.
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row + i, col: cur.col + i}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row - i, col: cur.col + i}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row - i, col: cur.col - i}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := Square{row: cur.row + i, col: cur.col - i}
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
		if p.Stop(s) {
			break
		}
	}
	return moves
}

func (p *BasePiece) KnightMoves(cur Square) []Square {
	moves := []Square{}
	// try all knight squares
	squares := []Square{
		{row: cur.row + 2, col: cur.col + 1},
		{row: cur.row + 2, col: cur.col - 1},
		{row: cur.row - 2, col: cur.col + 1},
		{row: cur.row - 2, col: cur.col - 1},
		{row: cur.row + 1, col: cur.col + 2},
		{row: cur.row + 1, col: cur.col - 2},
		{row: cur.row - 1, col: cur.col + 2},
		{row: cur.row - 1, col: cur.col - 2},
	}
	for _, s := range squares {
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
	}
	return moves
}

func (p *BasePiece) KingMoves(cur Square) []Square {
	moves := []Square{}
	// try all knight squares
	squares := []Square{
		{row: cur.row + 1, col: cur.col - 1},
		{row: cur.row + 1, col: cur.col},
		{row: cur.row + 1, col: cur.col + 1},
		{row: cur.row + -1, col: cur.col - 1},
		{row: cur.row + -1, col: cur.col},
		{row: cur.row + -1, col: cur.col + 1},
		{row: cur.row, col: cur.col - 1},
		{row: cur.row, col: cur.col + 1},
	}
	for _, s := range squares {
		if p.TargetLegal(s, true) {
			moves = append(moves, s)
		}
	}
	return moves

}

func (p *BasePiece) PawnMoves(cur Square) []Square {
	// Check if the piece can move two squares.
	var isStartPawn bool
	var direction int // Which way pawns move.
	switch p.color {
	case BLACK:
		isStartPawn = cur.row == 7
		direction = -1
		break
	case WHITE:
		isStartPawn = cur.row == 2
		direction = 1
		break
	}
	moves := []Square{}
	s := Square{row: cur.row + direction, col: cur.col}
	if p.TargetLegal(s, false) {
		moves = append(moves, s)
		// We only can move forward two if we can also move forward one.
		if isStartPawn {
			s := Square{row: cur.row + 2*direction, col: cur.col}
			if p.TargetLegal(s, false) {
				moves = append(moves, s)
			}
		}
	}
	// Check for side captures.
	captures := []Square{
		{row: cur.row + direction, col: cur.col + 1},
		{row: cur.row + direction, col: cur.col - 1},
	}
	for _, square := range captures {
		if !square.InPlay() {
			continue
		}
		occupant := p.board.Squares[square.Index()]
		if occupant != nil && occupant.Color() != p.Color() {
			moves = append(moves, square)
		}
	}
	// TODO: build en passant.
	return moves
}

// The default is no-op; only kings, rooks, and pawns have fancy accounting.
func (bp *BasePiece) ApplyMove(m Move) {
	return
}

// The nil piece can't move.
func (bp *BasePiece) LegalMoves() []Square {
	return []Square{}
}

// The nil piece.
func (bp *BasePiece) String() string {
	return "x"
}

func (bp *BasePiece) Board() *Board {
	return bp.board
}

func (bp *BasePiece) Value() float64 {
	return 0.0
}
func (bp *BasePiece) Graphic() string {
	return ""
}
