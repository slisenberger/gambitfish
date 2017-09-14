package game

import "fmt"

type Capture struct {
	Piece  Piece
	Square Square
}

// Moves have Pieces and squares
type Move struct {
	Piece       Piece
	Square      Square
	Old         Square
	Capture     *Capture
	EnPassant   bool
	CastleLong  bool
	CastleShort bool
	Promotion   Piece // Applicable for only Pawn moves
}

func (m Move) String() string {
	return fmt.Sprintf("%v%v to %v", m.Piece, m.Old, m.Square)
}

func NewMove(p Piece, square Square, old Square) Move {
	return Move{
		Piece:       p,
		Square:      square,
		Old:         old,
		EnPassant:   false,
		CastleLong:  false,
		CastleShort: false,
		Promotion:   nil,
	}
}

// Order the moves in an intelligent way for alpha beta pruning.
func OrderMoves(moves []Move) []Move {
	captures := []Move{}
	nonCaptures := []Move{}
	for _, m := range moves {
		if m.Capture != nil {
			captures = append(captures, m)
		} else {
			nonCaptures = append(nonCaptures, m)
		}
	}
	results := []Move{}
	results = append(results, captures...)
	results = append(results, nonCaptures...)
	return results
}
