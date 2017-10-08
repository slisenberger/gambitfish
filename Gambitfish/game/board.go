// Board is a representation of a chess Board and the operations possible on it.
package game

import "fmt"

type Board struct {
	Squares     [64]Piece
	Position    Position
	Active      Color
	Winner      Color
	WKSCastling bool
	WQSCastling bool
	BKSCastling bool
	BQSCastling bool
	Move        int
	EPSquare    Square // The square a pawn was just pushed two forward.
}

func DefaultBoard() *Board {
	b := &Board{Active: WHITE}
	// Add pawns
	for i := 1; i <= 8; i++ {
		blackPawnSquare := GetSquare(7, i)
		whitePawnSquare := GetSquare(2, i)
		b.Squares[blackPawnSquare] = BLACKPAWN

		b.Squares[whitePawnSquare] = WHITEPAWN
	}
	// Add rooks.
	b.Squares[0] = WHITEROOK
	b.Squares[7] = WHITEROOK
	b.Squares[56] = BLACKROOK
	b.Squares[63] = BLACKROOK
	// Add &Knights.
	b.Squares[1] = WHITEKNIGHT
	b.Squares[6] = WHITEKNIGHT
	b.Squares[57] = BLACKKNIGHT
	b.Squares[62] = BLACKKNIGHT
	// Add &Bishops.
	b.Squares[2] = WHITEBISHOP
	b.Squares[5] = WHITEBISHOP
	b.Squares[58] = BLACKBISHOP
	b.Squares[61] = BLACKBISHOP
	// Add queens
	b.Squares[3] = WHITEQUEEN
	b.Squares[59] = BLACKQUEEN
	// Add Kings
	b.Squares[4] = WHITEKING
	b.Squares[60] = BLACKKING
	b.Position = Position{}
	for s, p := range b.Squares {
		if p != NULLPIECE {
			b.Position = SetPiece(b.Position, p, Square(s))
		}
	}
	b.Position = UpdateBitboards(b.Position)
	b.WKSCastling = true
	b.WQSCastling = true
	b.BKSCastling = true
	b.BQSCastling = true
	b.Move = 1
	b.EPSquare = OFFBOARD_SQUARE
	return b
}

func (b *Board) Print() {
	fmt.Println(fmt.Sprintf("Move %v: %v to play", b.Move, b.Active))
	fmt.Println(fmt.Sprintf("Castling Rights:\n KINGSIDE: %v %v\n QUEENSIDE: %v %v", b.WKSCastling, b.BKSCastling, b.WQSCastling, b.BQSCastling))
	fmt.Println("")

	// We want to print a row at a time, but in backwards order to how this is stored
	// Since a1 occurs at the bottom left of the Board.
	for row := 7; row >= 0; row-- {
		newline := fmt.Sprintf("%v | ", row+1)
		// Get 8 consecutive Squares
		squares := b.Squares[8*row : 8*row+8]
		for i, piece := range squares {
			if piece == NULLPIECE {
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
	if p == NULLPIECE {
		b.Print()
		fmt.Println(b.AllLegalMoves())
		panic("nil piece: " + m.String())
	}
	// If there's a capture: remove the captured piece.
	if m.Capture != nil {
		b.Squares[m.Capture.Square] = NULLPIECE
		b.Position = UnSetPiece(b.Position, m.Capture.Piece, m.Capture.Square)
	}
	// Then, move the piece to its new square.
	// Check for promotion of a pawn.
	if m.Promotion != NULLPIECE {
		b.Squares[s] = m.Promotion
		b.Position = SetPiece(b.Position, m.Promotion, s)
	} else {
		b.Squares[s] = p
		b.Position = SetPiece(b.Position, p, s)
	}
	// Check for castling and modify rook state if so.  // New rook squares are relative to king.
	if m.QSCastle || m.KSCastle {
		var newRookSquare Square
		var oldRookSquare Square
		if m.QSCastle {
			newRookSquare = GetSquare(o.Row(), o.Col()-1)
			oldRookSquare = GetSquare(o.Row(), 1)
		}
		if m.KSCastle {
			newRookSquare = GetSquare(o.Row(), o.Col()+1)
			oldRookSquare = GetSquare(o.Row(), 8)
		}
		rook := b.Squares[oldRookSquare]
		b.Squares[newRookSquare] = rook
		b.Position = SetPiece(b.Position, rook, newRookSquare)
		b.Squares[oldRookSquare] = NULLPIECE
		b.Position = UnSetPiece(b.Position, rook, oldRookSquare)
	}
	// Then, remove the piece from its old square.
	b.Position = UnSetPiece(b.Position, p, m.Old)
	b.Squares[m.Old] = NULLPIECE
	// Modify castling state from rook and king moves.
	// We know any piece moving from e8, e1, a8, h8, a1, or h1 must
	// change castling rights.
	switch m.Old {
	case A8:
		b.BQSCastling = false
	case H8:
		b.BKSCastling = false
	case E8:
		b.BKSCastling = false
		b.BQSCastling = false
	case A1:
		b.WQSCastling = false
	case H1:
		b.WKSCastling = false
	case E1:
		b.WKSCastling = false
		b.WQSCastling = false
	}
	// Affect castling state of captured rooks.
	if m.Capture != nil {
		switch m.Capture.Square {
		case A8:
			b.BQSCastling = false
		case H8:
			b.BKSCastling = false
		case A1:
			b.WQSCastling = false
		case H1:
			b.WKSCastling = false
		}
	}
	// Apply En Passant state
	if m.TwoPawnAdvance {
		b.EPSquare = s
	} else {
		b.EPSquare = OFFBOARD_SQUARE
	}

	// Advance the move counter.
	if b.Active == BLACK {
		b.Move++
	}

	// Update bitboard representations.
	b.Position = UpdateBitboards(b.Position)
}

// UndoMove returns a board to the state it was at prior to
// applying a move m.
func UndoMove(b *Board, m Move) {
	p := m.Piece
	o := m.Old
	s := m.Square

	b.Squares[o] = p
	b.Position = SetPiece(b.Position, p, o)
	// Return the square we were on to its old state.
	b.Squares[s] = NULLPIECE
	if m.Promotion != NULLPIECE {
		b.Position = UnSetPiece(b.Position, m.Promotion, s)
	} else {
		b.Position = UnSetPiece(b.Position, p, s)
	}
	// Return a captured piece.
	if m.Capture != nil {
		b.Squares[m.Capture.Square] = m.Capture.Piece
		b.Position = SetPiece(b.Position, m.Capture.Piece, m.Capture.Square)
	}

	// Undo rook moves from castling.
	if m.QSCastle || m.KSCastle {
		var newRookSquare Square
		var oldRookSquare Square
		if m.QSCastle {
			newRookSquare = GetSquare(o.Row(), o.Col()-1)
			oldRookSquare = GetSquare(o.Row(), 1)
		}
		if m.KSCastle {
			newRookSquare = GetSquare(o.Row(), o.Col()+1)
			oldRookSquare = GetSquare(o.Row(), 8)
		}
		rook := b.Squares[newRookSquare]
		b.Squares[newRookSquare] = NULLPIECE
		b.Position = UnSetPiece(b.Position, rook, newRookSquare)
		b.Squares[oldRookSquare] = rook
		b.Position = SetPiece(b.Position, rook, oldRookSquare)
	}

	// Reapply original castling rights.
	b.WQSCastling = m.PrevWQSCastling
	b.BQSCastling = m.PrevBQSCastling
	b.WKSCastling = m.PrevWKSCastling
	b.BKSCastling = m.PrevBKSCastling

	// Reapply original en passant column.
	b.EPSquare = m.PrevEPSquare
	// Reverse the move counter.
	if b.Active == BLACK {
		b.Move--
	}

	b.Position = UpdateBitboards(b.Position)
}

func (b *Board) SwitchActivePlayer() {
	switch b.Active {
	case WHITE:
		b.Active = BLACK
	case BLACK:
		b.Active = WHITE
	}
}

// AllLegalMoves enumerates all of the legal moves currently available to the
// active player.
func (b *Board) AllLegalMoves() []Move {
	var moves []Move
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalMoves(b, p, Square(s))
		for _, move := range m {
			ApplyMove(b, move)
			if !IsCheck(b, b.Active) {
				moves = append(moves, move)
			}
			UndoMove(b, move)
		}
	}
	return moves
}

// AllLegalCaptures enumerates all of the legal moves currently available to the
// active player.
func (b *Board) AllLegalCaptures() []Move {
	var moves []Move
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalCaptures(b, p, Square(s))
		for _, move := range m {
			ApplyMove(b, move)
			if !IsCheck(b, b.Active) {
				moves = append(moves, move)
			}
			UndoMove(b, move)
		}
	}
	return moves
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

// Returns true if the board state results in the Color c's king being in check.
func IsCheck(b *Board, c Color) bool {
	atk := GetAttackBitboard(b, -1*c)
	switch c {
	case WHITE:
		return atk&b.Position.WhiteKing > 0
	case BLACK:
		return atk&b.Position.BlackKing > 0
	}

	return false
}

// Returns a bitboard of all squares currently being attacked by c.
// Attacked means a piece could be captured.
func GetAttackBitboard(b *Board, c Color) uint64 {
	var res uint64
	res = 0
	for s, p := range b.Squares {
		if p != NULLPIECE && p.Color() == c {
			res = res | AttackBitboard(b, p, Square(s))
		}
	}
	return res
}

// Returns the hash value for this board.
func ZobristHash(b *Board) uint64 {
	var hash uint64
	hash = 0
	for s, p := range b.Squares {
		if p != NULLPIECE {
			hash = hash ^ ZOBRISTPIECES[p.Color()][p.Type()][s]
		}
	}
	if b.Active == WHITE {
		hash = hash ^ ZOBRISTTURN
	}
	if b.WQSCastling {
		hash = hash ^ ZOBRISTWQS
	}
	if b.BQSCastling {
		hash = hash ^ ZOBRISTBQS
	}
	if b.WKSCastling {
		hash = hash ^ ZOBRISTWKS
	}
	if b.BKSCastling {
		hash = hash ^ ZOBRISTBKS
	}
	return hash
}
