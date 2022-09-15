package game

import "fmt"

import "sort"

// Moves have Pieces and squares
type Move struct {
	Piece           Piece
	Square          Square
	Old             Square
	Capture         Piece
	EnPassant       bool
	KSCastle        bool
	QSCastle        bool
	Promotion       Piece // Applicable for only Pawn moves
	TwoPawnAdvance  bool // For En Passant Management.
	NoMove bool
	Score float64
}

// Define a new move type that uses bits on an int
// Bit map:
// 1-4: Move piece
// 5-10: Old square
// 11-16: New Square
// 17-20: Capture Piece
// 21-24: Promotion Piece  // Can we remove this to make room for other bits?
// 25: Two pawn advance
// 26: En passant
// 27: Castle Kingside
// 28: Castle Queenside
type EfficientMove uint32

func NewEfficientMove(p Piece, s Square, o Square) EfficientMove {
	m := uint32(0)
	m = m | uint32(p) << 28
	m = m | uint32(o) << 22
	m = m | uint32(s) << 16
	return EfficientMove(m)
}

func (m EfficientMove) AddCapture (p Piece) EfficientMove {
	em := uint32(m) | uint32(p) << 12
	m = EfficientMove(em)
	return m
}

func (m EfficientMove) AddPromotion (p Piece) EfficientMove {
	return EfficientMove(uint32(m) | uint32(p) << 8)
}

func (m EfficientMove) AddTwoPawnAdvance() EfficientMove {
	return EfficientMove(uint32(m) | uint32(1) << 7)
}

func (m EfficientMove) AddEnPassant() EfficientMove {
	return EfficientMove(uint32(m) | uint32(1) << 6)
}

func (m EfficientMove) AddKSCastle() EfficientMove {
	return EfficientMove(uint32(m) | uint32(1) << 5)
}
func (m EfficientMove) AddQSCastle() EfficientMove {
	return EfficientMove(uint32(m) | uint32(1) << 4)
}

func MoveToEfficientMove(m Move) EfficientMove {
	em := NewEfficientMove(Piece(m.Piece), Square(m.Square), Square(m.Old))
	if m.Capture != NULLPIECE {
		em = em.AddCapture(Piece(m.Capture))
	}
	if m.Promotion != NULLPIECE {
		em = em.AddPromotion(m.Promotion)
	}
	if m.TwoPawnAdvance {
		em = em.AddTwoPawnAdvance()
	}
	if m.EnPassant {
		em = em.AddEnPassant()
	}
	if m.KSCastle {
		em = em.AddKSCastle()
	}
	if m.QSCastle {
		em = em.AddQSCastle()
	}
	return em
}

func EfficientMoveToMove(e EfficientMove) Move {
	p := (e & 0xF0000000) >> 28
	fmt.Printf("piece %v", Piece(p))
	o := (e & 0x0FC00000) >> 22
	fmt.Printf("square %v", Square(o))
	s := (e & 0x003F0000) >> 16
	fmt.Printf("old %v", Square(s))
	m := Move{
		Piece: Piece(p),
		Square: Square(s),
		Old: Square(o),
	}

	cp := (e & 0x0000F000 >> 12)
	if Piece(cp) != NULLPIECE {
		m.Capture = Piece(cp)
	}

	pp := (e & 0x00000F00 >> 8)
	tpa := (e & 0x00000080 >> 7)
	ep := (e & 0x00000040 >> 6)
	kc := (e & 0x00000020 >> 5)
	qc := (e & 0x00000010 >> 4)
	m.Promotion = Piece(pp)
	m.TwoPawnAdvance = (tpa == 1)
	m.EnPassant = (ep == 1)
	m.KSCastle = (kc == 1)
	m.QSCastle = (qc == 1)
	return m
}

func (e EfficientMove) Piece() Piece {
	return Piece((e & 0xF0000000) >> 28)
}

func (e EfficientMove) Old() Square {
	return Square((e & 0x0FC00000) >> 22)

}
func (e EfficientMove) Square() Square {
	return Square((e & 0x003F0000) >> 16)

}
func (e EfficientMove) Capture() Piece {
	return Piece((e & 0x0000F000 >> 12))

}
func (e EfficientMove) Promotion() Piece {
	return Piece((e & 0x00000F00 >> 8))

}
func (e EfficientMove) TwoPawnAdvance() bool {
	return (e & 0x00000080 >> 7) == 1

}
func (e EfficientMove) EnPassant() bool {
	return (e & 0x00000040 >> 6) == 1

}
func (e EfficientMove) KSCastle() bool {
	return (e & 0x00000020 >> 5) == 1

}
func (e EfficientMove) QSCastle() bool {
	return (e & 0x00000010 >> 4) == 1
}


func (m Move) String() string {
	var s string
	if m.QSCastle {
		return "O-O-O"
	} else if m.KSCastle {
		return "O-O"
	} else {
		s = fmt.Sprintf("%v%v", m.Piece, m.Old)
		if m.Capture != NULLPIECE {
			s += "x"
		} else {
			s += "-"
		}
		s += m.Square.String()
		if m.Promotion != NULLPIECE {
			s += "%v=%v" + m.Promotion.String()
		}

		if m.EnPassant {
			s += "(en passant)"
		}
	}
	return s
}

func (m EfficientMove) String() string {
	var s string
	if m.QSCastle() {
		return "O-O-O"
	} else if m.KSCastle() {
		return "O-O"
	} else {
		s = fmt.Sprintf("%v%v", m.Piece(), m.Old())
		if m.Capture() != NULLPIECE {
			s += "x"
		} else {
			s += "-"
		}
		s += m.Square().String()
		if m.Promotion() != NULLPIECE {
			s += "%v=%v" + m.Promotion().String()
		}

		if m.EnPassant() {
			s += "(en passant)"
		}
	}
	return s
}
func (m Move) Equals(m2 Move) bool {
	return m.Piece == m2.Piece && m.Square == m2.Square && m.Promotion == m2.Promotion && m.Capture == m2.Capture && m.TwoPawnAdvance == m2.TwoPawnAdvance && m.EnPassant == m2.EnPassant

}

func NewMove(p Piece, square Square, old Square, b *Board) Move {
	m := Move{
		Piece:           p,
		Square:          square,
		Old:             old,
		EnPassant:       false,
		KSCastle:        false,
		QSCastle:        false,
		Promotion:       NULLPIECE,
		TwoPawnAdvance:  false,
	}
	return m
}

type MoveRank struct {
	Move EfficientMove
	Score float64
}

// Order the moves in an intelligent way for alpha beta pruning.
func OrderMoves(b *Board, moves []EfficientMove, depth int, km KillerMoves) []EfficientMove {

	var k [2]EfficientMove
	if km != nil {
		k = km.GetKillerMoves(depth)
	} else {
		k = [2]EfficientMove{0, 0}
	}
	// Loop through the move list the rest of the times for other orderings.
	// Score constants
	captureScore := 1000.0
	km1Score := 1000.0
	km2Score := 999.0
	bestMoveScore := 2000.0
	e := PieceSquareEvaluator{}
	// Start with what we already believe the best move is.
	bestMove := EfficientMove(0)
	if entry, ok := TranspositionTable[ZobristHash(b)]; ok && entry.BestMove != EfficientMove(0) {
		bestMove = entry.BestMove
	}
	moveScores := make(map[EfficientMove]float64, len(moves))

	for i := 0; i < len(moves); i++ {
		m := moves[i]
		if m == bestMove {
			moveScores[m] = bestMoveScore
			continue
		}
		if m.Capture() != NULLPIECE {
			moveScores[m] = captureScore + m.Capture().Value() - m.Piece().Value()
			continue
		} else {

			// If it's a killer move, order it highly.
			if k[0] != EfficientMove(0) && k[0] == m {
				moveScores[m] = km1Score
				continue
			}
			if k[1] != EfficientMove(0) && k[1] == m {
				moveScores[m] = km2Score
				continue
			}
			// Order non captures by piece value weights.
			bs := ApplyMove(b, m)
			moveScores[m] = e.Evaluate(b)
	          	UndoMove(b, m, bs)
		}
	}
	sort.Slice(moves, func(i, j int) bool {
		return moveScores[moves[i]] > moveScores[moves[j]]
	})
	return moves
}
