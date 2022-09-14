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
// 21-26: Capture Square  // Can we remove this to make room for other bits?
// 27-30: Promotion Piece
// 31: Two pawn advance
// 32: En passant
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
	return EfficientMove(uint32(m) | uint32(p) << 2)
}

func (m EfficientMove) AddTwoPawnAdvance() EfficientMove {
	return EfficientMove(uint32(m) | uint32(1) << 1)
}

func (m EfficientMove) AddEnPassant() EfficientMove {
	return EfficientMove(uint32(m) | uint32(1))
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

	pp := (e & 0x0000003C >> 2)
	m.Promotion = Piece(pp)
	return m
}


func (m Move) String() string {
	var s string
	if m.QSCastle {
		s = "O-O-O"
	} else if m.KSCastle {
		s = "O-O"
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

func (m Move) Equals(m2 Move) bool {
	return m.Piece == m2.Piece && m.Square == m2.Square && m.Promotion == m2.Promotion && m.Capture == m.Capture

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

// Order the moves in an intelligent way for alpha beta pruning.
func OrderMoves(b *Board, moves []Move, depth int, km KillerMoves) []Move {

	var k [2]*Move
	if km != nil {
		k = km.GetKillerMoves(depth)
	} else {
		k = [2]*Move{nil, nil}
	}
	// Loop through the move list the rest of the times for other orderings.
	// Score constants
	captureScore := 1500.0
	km1Score := 1000.0
	km2Score := 999.0
	bestMoveScore := 2000.0
	e := PieceSquareEvaluator{}
	// Start with what we already believe the best move is.
	bestMove := Move{}
	if entry, ok := TranspositionTable[ZobristHash(b)]; ok && !entry.BestMove.NoMove {
		bestMove = entry.BestMove
	}

	for i := 0; i < len(moves); i++ {
		m := moves[i]
		if m.Equals(bestMove) {
			m.Score = bestMoveScore
			continue
		}
		if m.Capture != NULLPIECE {
			m.Score = captureScore + m.Capture.Value() - m.Piece.Value()
			continue
		} else {

			// If it's a killer move, order it highly.
			if k[0] != nil && k[0].Equals(m) {
				m.Score = km1Score
				continue
			}
			if k[1] != nil && k[1].Equals(m) {
				m.Score = km2Score
				continue
			}
			// Order non captures by piece value weights.
			bs := ApplyMove(b, m)
			m.Score = e.Evaluate(b)
			UndoMove(b, m, bs)
		}
		moves[i] = m
	}
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].Score > moves[j].Score
	})
	return moves
}
