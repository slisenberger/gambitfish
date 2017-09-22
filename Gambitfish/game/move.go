package game

import "fmt"
import "math/rand"

type Capture struct {
	Piece  Piece
	Square Square
}

// Moves have Pieces and squares
type Move struct {
	Piece                Piece
	Square               Square
	Old                  Square
	Capture              *Capture
	EnPassant            bool
	KSCastle             bool
	QSCastle             bool
	Promotion            Piece // Applicable for only Pawn moves
	PrevQSCastlingRights map[Color]bool
	PrevKSCastlingRights map[Color]bool
	PrevEPSquare         Square
	TwoPawnAdvance       bool // For En Passant Management.
}

func (m Move) String() string {
	if m.QSCastle {
		return "O-O-O"
	} else if m.KSCastle {
		return "O-O"
	}
	mv := fmt.Sprintf("%v%v", m.Piece, m.Old)
	if m.Capture != nil {
		mv += "x"
	} else {
		mv += "-"
	}
	mv += m.Square.String()
	if m.Promotion != nil {
		mv = fmt.Sprintf("%v=%v", mv, m.Promotion)
	}
	return mv
}

func NewMove(p Piece, square Square, old Square) Move {
	qr := p.Board().qsCastlingRights
	kr := p.Board().ksCastlingRights
	return Move{
		Piece:     p,
		Square:    square,
		Old:       old,
		EnPassant: false,
		KSCastle:  false,
		QSCastle:  false,
		Promotion: nil,
		PrevQSCastlingRights: map[Color]bool{
			WHITE: qr[WHITE],
			BLACK: qr[BLACK],
		},
		PrevKSCastlingRights: map[Color]bool{
			WHITE: kr[WHITE],
			BLACK: kr[BLACK],
		},
		PrevEPSquare:   p.Board().EPSquare,
		TwoPawnAdvance: false,
	}
}

// Order the moves in an intelligent way for alpha beta pruning.
func OrderMoves(moves []Move) []Move {
	// We now want to return the moves in a smart order (e.g., try checks and captures.)
	// Just to get things off the ground, we'll shuffle the moves, just to get some variety
	// in the AI vs AI games.
	for i := range moves {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
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
