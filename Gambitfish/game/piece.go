// Piepce is an interface that defines the operations possible for a piece on the board.
package game

// Define the possible Colors of a piece as an enum
type Color int

const (
	WHITE Color = 1
	BLACK Color = -1
)

func (c Color) String() string {
	if c == WHITE {
		return "WHITE"
	} else {
		return "BLACK"
	}
}

type Piece interface {
	// Returns an array of all the legal positions this piece can move to.
	LegalMoves() []Move
	// Returns an array of all the squares this piece is attacking. Only open
	// squares or squares containing opponent pieces are considered attacked.

	Attacking() []Square
	// Returns a string representation of this piece.
	String() string
	// Returns a unicode graphic representation of this piece.
	Graphic() string
	// Returns the Color of this piece.
	Color() Color
	Board() *Board
	ApplyMove(Move)
	Value() float64
}

type BasePiece struct {
	C Color
	B *Board
}

func (bp *BasePiece) Color() Color {
	return bp.C
}

// TargetLegal returns true if a candidate piece can move to the
// desired square. It also optionally returns a piece that will
// be captured.
func TargetLegal(p Piece, s Square, capture bool) (bool, Piece) {
	if s == OFFBOARD_SQUARE {
		return false, nil
	}
	occupant := p.Board().Squares[s]
	if occupant == nil {
		return true, nil
	} else {
		if capture && p.Color() != occupant.Color() {
			return true, occupant
		}
	}
	return false, nil
}

// Returns true if we should stop searching along a ray, because a piece has
// encountered a blockage.
func Stop(p Piece, s Square) bool {
	if s == OFFBOARD_SQUARE {
		return true
	}
	if occupant := p.Board().Squares[s]; occupant != nil {
		return true
	}
	return false
}

func ColumnAndRowMoves(p Piece, cur Square) []Move {
	moves := []Move{}
	// Move in each direction, checking for blocking pieces.
	// Left in Row.
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row(), cur.Col()-i)
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row(), cur.Col()+i)
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row()+i, cur.Col())
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row()-i, cur.Col())
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	return moves
}

func DiagonalMoves(p Piece, cur Square) []Move {
	moves := []Move{}
	// Move in a diagonal, checking for blocking pieces.
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row()+i, cur.Col()+i)
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row()-i, cur.Col()+i)
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row()-i, cur.Col()-i)
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	for i := 1; i <= 7; i++ {
		s := GetSquare(cur.Row()+i, cur.Col()-i)
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
		if Stop(p, s) {
			break
		}
	}
	return moves
}

func KnightMoves(p Piece, cur Square) []Move {
	moves := []Move{}
	// try all knight squares
	squares := []Square{
		GetSquare(cur.Row()+2, cur.Col()+1),
		GetSquare(cur.Row()+2, cur.Col()-1),
		GetSquare(cur.Row()-2, cur.Col()+1),
		GetSquare(cur.Row()-2, cur.Col()-1),
		GetSquare(cur.Row()+1, cur.Col()+2),
		GetSquare(cur.Row()+1, cur.Col()-2),
		GetSquare(cur.Row()-1, cur.Col()+2),
		GetSquare(cur.Row()-1, cur.Col()-2),
	}
	for _, s := range squares {
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
	}
	return moves
}

func KingMoves(p Piece, cur Square) []Move {
	moves := []Move{}
	// try all knight squares
	squares := []Square{
		GetSquare(cur.Row()+1, cur.Col()-1),
		GetSquare(cur.Row()+1, cur.Col()),
		GetSquare(cur.Row()+1, cur.Col()+1),
		GetSquare(cur.Row()+-1, cur.Col()-1),
		GetSquare(cur.Row()+-1, cur.Col()),
		GetSquare(cur.Row()+-1, cur.Col()+1),
		GetSquare(cur.Row(), cur.Col()-1),
		GetSquare(cur.Row(), cur.Col()+1),
	}
	for _, s := range squares {
		if l, capture := TargetLegal(p, s, true); l {
			move := NewMove(p, s, cur)
			if capture != nil {
				move.Capture = &Capture{Piece: capture, Square: s}
			}
			moves = append(moves, move)
		}
	}
	return moves

}

// Returns the set of squares a pawn is attacking.
func PawnAttackingSquares(p Piece, cur Square) []Square {
	var direction int // Which way pawns move.
	switch p.Color() {
	case BLACK:
		direction = -1
		break
	case WHITE:
		direction = 1
		break
	}

	// Check for side attacks.
	attacks := []Square{
		GetSquare(cur.Row()+direction, cur.Col()+1),
		GetSquare(cur.Row()+direction, cur.Col()-1),
	}
	var results []Square
	for _, attack := range attacks {
		// If it's not legal to move there, we aren't attacking it.
		if attack == OFFBOARD_SQUARE {
			continue
		}
		// We aren't attacking if it is our own piece.
		occupant := p.Board().Squares[attack]
		if occupant != nil && occupant.Color() == p.Color() {
			continue
		}
		// Else, it's fine.
		results = append(results, attack)
	}
	return results
}

func PawnMoves(p Piece, cur Square) []Move {
	// Check if the piece can move two squares.
	var isStartPawn bool
	var direction int // Which way pawns move.
	switch p.Color() {
	case BLACK:
		isStartPawn = cur.Row() == 7
		direction = -1
		break
	case WHITE:
		isStartPawn = cur.Row() == 2
		direction = 1
		break
	}
	moves := []Move{}
	s := GetSquare(cur.Row()+direction, cur.Col())
	if l, _ := TargetLegal(p, s, false); l {
		// Check for promotion and add the promotions.
		if s.Row() == 1 || s.Row() == 8 {
			move := NewMove(p, s, cur)
			move.Promotion = &Queen{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
			moves = append(moves, move)
			move = NewMove(p, s, cur)
			move.Promotion = &Bishop{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
			moves = append(moves, move)
			move = NewMove(p, s, cur)
			move.Promotion = &Knight{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
			moves = append(moves, move)
			move = NewMove(p, s, cur)
			move.Promotion = &Rook{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
			moves = append(moves, move)

		} else {
			moves = append(moves, NewMove(p, s, cur))
		}
		// We only can move forward two if we can also move forward one.
		if isStartPawn {
			s := GetSquare(cur.Row()+2*direction, cur.Col())
			if l, _ := TargetLegal(p, s, false); l {
				m := NewMove(p, s, cur)
				m.TwoPawnAdvance = true
				moves = append(moves, m)
			}
		}
	}
	// Check for side captures.
	captures := []Square{
		GetSquare(cur.Row()+direction, cur.Col()+1),
		GetSquare(cur.Row()+direction, cur.Col()-1),
	}
	for _, s := range captures {
		if s == OFFBOARD_SQUARE {
			continue
		}
		occupant := p.Board().Squares[s]
		if occupant != nil && occupant.Color() != p.Color() {
			if s.Row() == 1 || s.Row() == 8 {
				move := NewMove(p, s, cur)
				move.Capture = &Capture{occupant, s}
				move.Promotion = &Queen{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
				moves = append(moves, move)
				move = NewMove(p, s, cur)
				move.Capture = &Capture{occupant, s}
				move.Promotion = &Bishop{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
				moves = append(moves, move)
				move = NewMove(p, s, cur)
				move.Capture = &Capture{occupant, s}
				move.Promotion = &Knight{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
				moves = append(moves, move)
				move = NewMove(p, s, cur)
				move.Capture = &Capture{occupant, s}
				move.Promotion = &Rook{BasePiece: &BasePiece{B: p.Board(), C: p.Color()}}
				moves = append(moves, move)

			} else {
				move := NewMove(p, s, cur)
				move.Capture = &Capture{occupant, s}
				moves = append(moves, move)
			}
		}
	}
	// Check for en passants
	epCol := p.Board().EPCol
	adjToEP := cur.Col()-1 == epCol || cur.Col()+1 == epCol
	if epCol != 0 {
		if p.Color() == WHITE && cur.Row() == 5 && adjToEP {
			move := NewMove(p, GetSquare(6, epCol), cur)
			epSquare := GetSquare(5, epCol)
			move.Capture = &Capture{p.Board().Squares[epSquare], epSquare}
			moves = append(moves, move)
		}
		if p.Color() == BLACK && cur.Row() == 4 && adjToEP {
			move := NewMove(p, GetSquare(3, epCol), cur)
			epSquare := GetSquare(4, epCol)
			move.Capture = &Capture{p.Board().Squares[epSquare], epSquare}
			moves = append(moves, move)
		}
	}
	return moves
}

// The default is no-op; only kings, rooks, and pawns have fancy accounting.
func (bp *BasePiece) ApplyMove(m Move) {
	return
}

// The nil piece can't move.
func (bp *BasePiece) LegalMoves() []Move {
	return []Move{}
}

// The nil piece.
func (bp *BasePiece) String() string {
	return "x"
}

func (bp *BasePiece) Board() *Board {
	return bp.B
}

func (bp *BasePiece) Value() float64 {
	return 0.0
}
func (bp *BasePiece) Graphic() string {
	return ""
}
