// Board is a representation of a chess Board and the operations possible on it.
package game

import "fmt"
import "math/rand"

type Board struct {
	Squares          [64]Piece
	Active           Color
	Winner           Color
	PieceSet         map[Piece]Square
	ksCastlingRights map[Color]bool
	qsCastlingRights map[Color]bool
	Move             int
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
		b.Squares[blackPawnSquare.Index()] = &Pawn{&BasePiece{C: BLACK, B: b}}
		b.Squares[whitePawnSquare.Index()] = &Pawn{&BasePiece{C: WHITE, B: b}}
	}
	// Add rooks.
	b.Squares[0] = &Rook{&BasePiece{C: WHITE, B: b}, false, true}
	b.Squares[7] = &Rook{&BasePiece{C: WHITE, B: b}, true, false}
	b.Squares[56] = &Rook{&BasePiece{C: BLACK, B: b}, false, true}
	b.Squares[63] = &Rook{&BasePiece{C: BLACK, B: b}, true, false}
	// Add &Knights.
	b.Squares[1] = &Knight{&BasePiece{C: WHITE, B: b}}
	b.Squares[6] = &Knight{&BasePiece{C: WHITE, B: b}}
	b.Squares[57] = &Knight{&BasePiece{C: BLACK, B: b}}
	b.Squares[62] = &Knight{&BasePiece{C: BLACK, B: b}}
	// Add &Bishops.
	b.Squares[2] = &Bishop{&BasePiece{C: WHITE, B: b}}
	b.Squares[5] = &Bishop{&BasePiece{C: WHITE, B: b}}
	b.Squares[58] = &Bishop{&BasePiece{C: BLACK, B: b}}
	b.Squares[61] = &Bishop{&BasePiece{C: BLACK, B: b}}
	// Add queens
	b.Squares[3] = &Queen{&BasePiece{C: WHITE, B: b}}
	b.Squares[59] = &Queen{&BasePiece{C: BLACK, B: b}}
	// Add Kings
	b.Squares[4] = &King{&BasePiece{C: WHITE, B: b}, false}
	b.Squares[60] = &King{&BasePiece{C: BLACK, B: b}, false}
	b.InitPieceSet()
	b.ksCastlingRights = map[Color]bool{WHITE: true, BLACK: true}
	b.qsCastlingRights = map[Color]bool{WHITE: true, BLACK: true}
	b.Move = 1
	return b
}

func (b *Board) Print() {
	fmt.Println(fmt.Sprintf("Move %v: %v to play", b.Move, b.Active))
	fmt.Println(fmt.Sprintf("Castling Rights:\n KINGSIDE: %v\n QUEENSIDE: %v", b.ksCastlingRights, b.qsCastlingRights))
	fmt.Println("")

	// We want to print a Row at a time, but in backwards order to how this is stored
	// Since a1 occurs at the bottom left of the Board.
	for Row := 7; Row >= 0; Row-- {
		newline := fmt.Sprintf("%v | ", Row+1)
		// Get 8 consecutive Squares
		Squares := b.Squares[8*Row : 8*Row+8]
		for i, piece := range Squares {
			if piece == nil {
				newline += "·"
			} else {
				newline += piece.Graphic()
			}
			newline += " "
			// Print a new line every 8 Squares, and reset.
			if ((i + 1) % 8) == 0 {
				fmt.Println(newline)
				newline = ""
			}
		}
	}
	// Print the letters at the bottom.
	fmt.Println("   ________________")
	fmt.Println("    a b c d e f g h")

}

// ApplyMove changes the state of the Board for any given move.
func ApplyMove(b *Board, m Move) {
	p := m.Piece
	o := m.Old
	s := m.Square
	// If there's a capture: remove the captured piece.
	if m.Capture != nil {
		delete(b.PieceSet, m.Capture.Piece)
	}
	// Then, move the piece to its new square.
	// Check for promotion of a pawn.
	if m.Promotion != nil {
		b.PieceSet[m.Promotion] = s
		b.Squares[s.Index()] = m.Promotion
		delete(b.PieceSet, p)
	} else {
		b.PieceSet[p] = s
		b.Squares[s.Index()] = p
	}
	// Check for castling and modify rook state if so.
	// New rook squares are relative to king.
	if m.QSCastle || m.KSCastle {
		var newRookSquare Square
		var oldRookSquare Square
		if m.QSCastle {
			newRookSquare = Square{o.Row, o.Col - 1}
			oldRookSquare = Square{o.Row, 1}
		}
		if m.KSCastle {
			newRookSquare = Square{o.Row, o.Col + 1}
			oldRookSquare = Square{o.Row, 8}
		}
		rook := b.Squares[oldRookSquare.Index()]
		b.PieceSet[rook] = newRookSquare
		b.Squares[newRookSquare.Index()] = rook
		b.Squares[oldRookSquare.Index()] = nil
	}
	// Then, remove the piece from its old square.
	b.Squares[m.Old.Index()] = nil
	// Modify castling state from rook and king moves.
	if _, ok := p.(*King); ok {
		b.qsCastlingRights[p.Color()] = false
		b.ksCastlingRights[p.Color()] = false
	}
	if r, ok := p.(*Rook); ok {
		if r.QS {
			b.qsCastlingRights[p.Color()] = false
		} else if r.KS {
			b.ksCastlingRights[p.Color()] = false
		}
	}
	// Affect castling state of captured rooks.
	if m.Capture != nil {
		if r, ok := m.Capture.Piece.(*Rook); ok {
			if r.QS {
				b.qsCastlingRights[r.Color()] = false
			} else if r.KS {
				b.ksCastlingRights[r.Color()] = false
			}
		}
	}
}

// UndoMove returns a board to the state it was at prior to
// applying a move m.
func UndoMove(b *Board, m Move) {
	p := m.Piece
	o := m.Old
	s := m.Square
	// Return the piece to its old square, and undo promotion.
	if m.Promotion != nil {
		delete(b.PieceSet, m.Promotion)
	}

	b.Squares[o.Index()] = p
	b.PieceSet[p] = o
	// Return the square we were on to its old state.
	b.Squares[s.Index()] = nil
	// Return a captured piece.
	if m.Capture != nil {
		b.PieceSet[m.Capture.Piece] = m.Capture.Square
		b.Squares[m.Capture.Square.Index()] = m.Capture.Piece
	}

	// Undo rook moves from castling.
	if m.QSCastle || m.KSCastle {
		var newRookSquare Square
		var oldRookSquare Square
		if m.QSCastle {
			newRookSquare = Square{o.Row, o.Col + -1}
			oldRookSquare = Square{o.Row, 1}
		}
		if m.KSCastle {
			newRookSquare = Square{o.Row, o.Col + 1}
			oldRookSquare = Square{o.Row, 8}
		}
		rook := b.Squares[newRookSquare.Index()]
		b.Squares[newRookSquare.Index()] = nil
		b.Squares[oldRookSquare.Index()] = rook
		b.PieceSet[rook] = oldRookSquare

	}

	// Reapply original castling rights.
	b.qsCastlingRights = m.PrevQSCastlingRights
	b.ksCastlingRights = m.PrevKSCastlingRights
}

func (b *Board) SwitchActivePlayer() {
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
		m := piece.LegalMoves()
		for _, move := range m {
			ApplyMove(b, move)
			if !IsCheck(b, b.Active) {
				moves = append(moves, move)
			}
			UndoMove(b, move)
		}
	}
	// We now want to return the moves in a smart order (e.g., try checks and captures.)
	// Just to get things off the ground, we'll shuffle the moves, just to get some variety
	// in the AI vs AI games.
	for i := range moves {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
	return OrderMoves(moves)
}

// Returns true if the game is drawn, won, or lost. The integer
// returned is the winner (or draw). This function should be called on the beginning
// of a move.
func (b *Board) CalculateGameOver() (bool, Color) {
	if len(b.AllLegalMoves()) == 0 {
		// If the active player is in check, they lose.
		if IsCheck(b, b.Active) {
			return true, -1 * b.Active
		} else {
			// Otherwise, it's stalemate!
			return true, 0
		}

	}
	return false, 0
}

// Returns a copy of this Board with different references.
func CopyBoard(b *Board) *Board {
	var s [64]Piece
	copy(s[:], b.Squares[:])
	newB := &Board{Squares: s, Active: b.Active, Winner: b.Winner}
	newB.InitPieceSet()
	return newB
}

// Returns true if the board state results in the Color c's king being in check.
func IsCheck(b *Board, c Color) bool {
	attacking := GetAttacking(b, -1*c)
	for _, s := range attacking {
		occupant := b.Squares[s.Index()]
		// check if it's us under attack
		if occupant != nil && occupant.Color() == c {
			// If our king is under attack, it's check.
			if _, ok := occupant.(*King); ok {
				return true
			}
		}
	}
	return false
}

// Returns all attacked squares by c's pieces.
func GetAttacking(b *Board, c Color) []Square {
	results := []Square{}
	for piece, s := range b.PieceSet {
		if piece == nil {
			panic(fmt.Sprintf("have nil piece in piece set: square %v", s))
		}
		// If it's our piece, we don't care.
		if piece.Color() != c {
			continue
		}
		results = append(results, piece.Attacking()...)
	}
	return results
}
