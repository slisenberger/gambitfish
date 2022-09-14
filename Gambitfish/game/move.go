package game

import "fmt"

import "sort"

type Capture struct {
	Piece  Piece
	Square Square
}

// Moves have Pieces and squares
type Move struct {
	Piece           Piece
	Square          Square
	Old             Square
	Capture         *Capture
	EnPassant       bool
	KSCastle        bool
	QSCastle        bool
	Promotion       Piece // Applicable for only Pawn moves
	PrevWQSCastling bool
	PrevWKSCastling bool
	PrevBQSCastling bool
	PrevBKSCastling bool
	PrevCheck       bool
	PrevLastMove    *Move
	PrevEPSquare    Square
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
// 21-26: Capture Square // 26-29: Promotion Piece
// 30: Two pawn push?
type EfficientMove uint32

func NewEfficientMove(p Piece, s Square, o Square) EfficientMove {
	m := uint32(0)
	m = m | uint32(p) << 28
	m = m | uint32(o) << 22
	m = m | uint32(s) << 16
	return EfficientMove(m)
}

func (m EfficientMove) AddCapture (p Piece, s Square) {
	em := uint32(m) | uint32(p) << 12
	em = em | uint32(s) << 6
	m = EfficientMove(em)
}

func MoveToEfficientMove(m Move) EfficientMove {
	em := NewEfficientMove(Piece(m.Piece), Square(m.Square), Square(m.Old))
	if m.Capture != nil {
		em.AddCapture(Piece(m.Capture.Piece), Square(m.Capture.Square))
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
	m := NewMoveNoBoard(Piece(p), Square(s), Square(o))
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
		if m.Capture != nil {
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
	return m.Piece == m2.Piece && m.Square == m2.Square && m.Promotion == m2.Promotion && m.Capture == m2.Capture
}

func NewMove(p Piece, square Square, old Square, b *Board) Move {
	m := Move{
		Piece:           p,
		Square:          square,
		Old:             old,
		EnPassant:       false,
		KSCastle:        false,
		QSCastle:        false,
		PrevBQSCastling: b.BQSCastling,
		PrevWQSCastling: b.WQSCastling,
		PrevBKSCastling: b.BKSCastling,
		PrevWKSCastling: b.WKSCastling,
		PrevLastMove:    b.LastMove,
		Promotion:       NULLPIECE,
		PrevEPSquare:    b.EPSquare,
		TwoPawnAdvance:  false,
	}
	return m
}

// Used for testing conversions.
func NewMoveNoBoard(p Piece, square Square, old Square) Move {
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
		if m.Capture != nil {
			m.Score = captureScore + m.Capture.Piece.Value() - m.Piece.Value()
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
			ApplyMove(b, m)
			m.Score = e.Evaluate(b)
			UndoMove(b, m)
		}
		moves[i] = m
	}
	// Sort moves by scores.
	//fmt.Println("Presort")
	//fmt.Println("-------")
	//fmt.Println(moves)
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].Score > moves[j].Score
	})
	//fmt.Println("Postsort")
	//fmt.Println("-------")
	//fmt.Println(moves)
	return moves
}
