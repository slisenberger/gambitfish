package game

import "fmt"

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
