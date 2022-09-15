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
	LastMove    EfficientMove
	EPSquare    Square // The square a pawn was just pushed two forward.
	AllMoves    []EfficientMove
	History     []Position
}

type BoardState struct {
	LastMove EfficientMove
	WKSCastling bool
	WQSCastling bool
	BKSCastling bool
	BQSCastling bool
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
	b.History = make([]Position, 50)
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
func ApplyMove(b *Board, m EfficientMove) BoardState {
	p := m.Piece()
	o := m.Old()
	s := m.Square()
	c := m.Capture()
	// Save old state.
	bs := BoardState{
		LastMove: b.LastMove,
		WKSCastling: b.WKSCastling,
		BKSCastling: b.BKSCastling,
		WQSCastling: b.WQSCastling,
		BQSCastling: b.BQSCastling,
		EPSquare: b.EPSquare,
	}
	if p == NULLPIECE {
		b.Print()
		fmt.Println(b.AllLegalMoves())
		panic("nil piece: " + m.String())
	}
	if m.Capture() == NULLPIECE && b.Squares[s] != NULLPIECE {
		b.Print()
		fmt.Println(m)
		fmt.Println("don't have capture in a move to occupied sq..")
	}
	if m.Capture() != NULLPIECE && s > H8 {
		b.Print()
		fmt.Println(m)
		fmt.Println("capturing offboard sq..")
	}
	// If there's a capture: remove the captured piece.
	if c != NULLPIECE {
		// In en passant, the piece is not on the square we move to.
		if m.EnPassant() {
			b.Position = UnSetPiece(b.Position, c, b.EPSquare)
			b.Squares[b.EPSquare] = NULLPIECE
		} else {
			b.Position = UnSetPiece(b.Position, c, s)
			b.Squares[s] = NULLPIECE
		}
	}
	// Then, move the piece to its new square.
	// Check for promotion of a pawn.
	if m.Promotion() != NULLPIECE {
		b.Squares[s] = m.Promotion()
		b.Position = SetPiece(b.Position, m.Promotion(), s)
	} else {
		b.Squares[s] = p
		b.Position = SetPiece(b.Position, p, s)
	}
	// Check for castling and modify rook state if so.  // New rook squares are relative to king.
	if m.QSCastle() || m.KSCastle() {
		var newRookSquare Square
		var oldRookSquare Square
		if m.QSCastle() {
			newRookSquare = GetSquare(o.Row(), o.Col()-1)
			oldRookSquare = GetSquare(o.Row(), 1)
		}
		if m.KSCastle() {
			newRookSquare = GetSquare(o.Row(), o.Col()+1)
			oldRookSquare = GetSquare(o.Row(), 8)
		}
		b.Squares[newRookSquare] = b.Squares[oldRookSquare]
		b.Position = SetPiece(b.Position, b.Squares[oldRookSquare], newRookSquare)
		b.Position = UnSetPiece(b.Position, b.Squares[oldRookSquare], oldRookSquare)
		b.Squares[oldRookSquare] = NULLPIECE
	}
	// Then, remove the piece from its old square.
	b.Position = UnSetPiece(b.Position, p, o)
	b.Squares[o] = NULLPIECE
	// Modify castling state from rook and king moves.
	// We know any piece moving from e8, e1, a8, h8, a1, or h1 must
	// change castling rights.
	switch o {
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
	if c != NULLPIECE {
		switch s {
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
	if m.TwoPawnAdvance() {
		b.EPSquare = s
	} else {
		b.EPSquare = OFFBOARD_SQUARE
	}

	// Advance the move counter.
	if b.Active == BLACK {
		b.Move++
	}
	b.LastMove = m
	// Update bitboard representations.
	b.Position = UpdateBitboards(b.Position)

	// Update this board's move history.
	b.History = append(b.History, b.Position)
	return bs
}

// UndoMove returns a board to the state it was at prior to
// applying a move m.
func UndoMove(b *Board, m EfficientMove, bs BoardState) {
	p := m.Piece()
	o := m.Old()
	s := m.Square()
	c := m.Capture()

	// Remove ourselves from our square.
	b.Squares[s] = NULLPIECE

	// Undo promotions and change our candidate piece to a pawn.
	if m.Promotion() != NULLPIECE {
		b.Position = UnSetPiece(b.Position, m.Promotion(), s)
	} else {
		b.Position = UnSetPiece(b.Position, p, s)
	}
	// Return our piece to its original place.
	b.Position = SetPiece(b.Position, p, o)
	b.Squares[o] = p
	// Return a captured piece.
	if c != NULLPIECE {
		if m.EnPassant(){
			b.Squares[bs.EPSquare] = c
		        b.Position = SetPiece(b.Position, c, bs.EPSquare)
		} else {
			b.Squares[s] = c
		        b.Position = SetPiece(b.Position, c, s)
		}
	}

	// Undo rook moves from castling.
	if m.QSCastle() || m.KSCastle() {
		var newRookSquare Square
		var oldRookSquare Square
		if m.QSCastle() {
			newRookSquare = GetSquare(o.Row(), o.Col()-1)
			oldRookSquare = GetSquare(o.Row(), 1)
		}
		if m.KSCastle() {
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
	b.WQSCastling = bs.WQSCastling
	b.BQSCastling = bs.BQSCastling
	b.WKSCastling = bs.WKSCastling
	b.BKSCastling = bs.BKSCastling

	// Reapply original en passant column.
	b.EPSquare = bs.EPSquare
	// Reverse the move counter.
	if b.Active == BLACK {
		b.Move--
	}
	b.LastMove = bs.LastMove


	b.Position = UpdateBitboards(b.Position)
	b.History = b.History[:len(b.History) - 1]

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
func (b *Board) AllLegalMoves() []EfficientMove {
	var moves []EfficientMove
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalMoves(b, p, Square(s))
		for i := 0; i < len(m); i++ {
			bs := ApplyMove(b, m[i])
			if !IsCheck(b, b.Active) {
				moves = append(moves, m[i])
			}
			UndoMove(b, m[i], bs)
		}
	}
	return moves
}

// AllLegalCaptures enumerates all of the legal captures currently available to the
// active player.
func (b *Board) AllLegalCaptures() []EfficientMove {
	var moves []EfficientMove
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalCaptures(b, p, Square(s))
		for _, move := range m {
			bs := ApplyMove(b, move)
			if !IsCheck(b, b.Active) {
				moves = append(moves, move)
			}
			UndoMove(b, move, bs)
		}
	}
	return moves
}

// AllLegalChecks enumerates all of the legal checks currently available to the
// active player.
func (b *Board) AllLegalChecks() []EfficientMove {
	var moves []EfficientMove
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalMoves(b, p, Square(s))
		for _, move := range m {
			bs := ApplyMove(b, move)
			if !IsCheck(b, b.Active) && IsCheck(b, -1 * b.Active) {
				moves = append(moves, move)
			}
			UndoMove(b, move, bs)
		}
	}
	return moves
}

// AllLegalChecksAndCaptures enumerates all of the loud moves currently available to the
// active player.
func (b *Board) AllLegalChecksAndCaptures() []EfficientMove {
	var moves []EfficientMove
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalMoves(b, p, Square(s))
		for _, move := range m {
			bs := ApplyMove(b, move)
			if !IsCheck(b, b.Active) && (IsCheck(b, -1 * b.Active) || move.Capture() != NULLPIECE){
				moves = append(moves, move)
			}
			UndoMove(b, move, bs)
		}
	}
	return moves
}

func (b *Board) AllQuiescenceMoves() []EfficientMove {
	var moves []EfficientMove
	for s, p := range b.Squares {
		if p == NULLPIECE {
			continue
		}
		if p.Color() != b.Active {
			continue
		}
		m := LegalMoves(b, p, Square(s))
		for _, move := range m {
			// In check, all moves should be searched.
			if IsCheck(b, b.Active) {
				moves = append(moves, move)
				continue

			}
			// Otherwise, see if it's a check or capture.
			bs := ApplyMove(b, move)
			if !IsCheck(b, b.Active) && (IsCheck(b, -1 * b.Active) || move.Capture() != NULLPIECE){
				moves = append(moves, move)
			}
			UndoMove(b, move, bs)
		}
	}
	return moves

}

// Returns true if the game is drawn, won, or lost. The integer
// returned is the winner (or draw). This function should be called on the beginning
// of a move.
func (b *Board) CalculateGameOver(lm []EfficientMove) (bool, Color) {
	if len(lm) == 0 {
		// If the active player is in check, they lose.
		if IsCheck(b, b.Active) {
			return true, -1 * b.Active
		} else {
			// Otherwise, it's stalemate!
			return true, 0
		}

	}

	posCount := 0
	for _, h := range b.History {
		if h == b.Position {
			posCount += 1
		}
		if posCount >= 3 {
			// Draw by repetition!
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
	var p Piece
	res = 0
	for i := 0; i < 64; i++ {
		p = b.Squares[i]
		if p != NULLPIECE && p.Color() == c {
			res = res | AttackBitboard(b, p, Square(i))
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
			hash = hash ^ ZOBRISTPIECES[p][s]
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
