// Board is a representation of a chess Board and the operations possible on it.
package game

import "fmt"

type Board struct {
	Squares          [64]Piece
	Position         Position
	Active           Color
	Winner           Color
	PieceSet         map[Piece]Square
	ksCastlingRights map[Color]bool
	qsCastlingRights map[Color]bool
	Move             int
	EPSquare         Square // The square a pawn was just pushed two forward.
}

func (b *Board) InitPieceSet() {
	b.PieceSet = make(map[Piece]Square, 32)
	for i, piece := range b.Squares {
		if piece != nil {
			b.PieceSet[piece] = Square(i)
		}
	}
}

func DefaultBoard() *Board {
	b := &Board{Active: WHITE}
	// Add pawns
	for i := 1; i <= 8; i++ {
		blackPawnSquare := GetSquare(7, i)
		whitePawnSquare := GetSquare(2, i)
		b.Squares[blackPawnSquare] = &Pawn{C: BLACK}

		b.Squares[whitePawnSquare] = &Pawn{C: WHITE}
	}
	// Add rooks.
	b.Squares[0] = &Rook{WHITE, false, true}
	b.Squares[7] = &Rook{WHITE, true, false}
	b.Squares[56] = &Rook{BLACK, false, true}
	b.Squares[63] = &Rook{BLACK, true, false}
	// Add &Knights.
	b.Squares[1] = &Knight{C: WHITE}
	b.Squares[6] = &Knight{C: WHITE}
	b.Squares[57] = &Knight{C: BLACK}
	b.Squares[62] = &Knight{C: BLACK}
	// Add &Bishops.
	b.Squares[2] = &Bishop{C: WHITE}
	b.Squares[5] = &Bishop{C: WHITE}
	b.Squares[58] = &Bishop{C: BLACK}
	b.Squares[61] = &Bishop{C: BLACK}
	// Add queens
	b.Squares[3] = &Queen{C: WHITE}
	b.Squares[59] = &Queen{C: BLACK}
	// Add Kings
	b.Squares[4] = &King{C: WHITE}
	b.Squares[60] = &King{C: BLACK}
	b.InitPieceSet()
	b.Position = Position{}
	for p, s := range b.PieceSet {
		b.Position = SetPiece(b.Position, p, s)
	}
	b.Position = UpdateBitboards(b.Position)
	b.ksCastlingRights = map[Color]bool{WHITE: true, BLACK: true}
	b.qsCastlingRights = map[Color]bool{WHITE: true, BLACK: true}
	b.Move = 1
	b.EPSquare = OFFBOARD_SQUARE
	return b
}

func (b *Board) Print() {
	fmt.Println(fmt.Sprintf("Move %v: %v to play", b.Move, b.Active))
	fmt.Println(fmt.Sprintf("Castling Rights:\n KINGSIDE: %v\n QUEENSIDE: %v", b.ksCastlingRights, b.qsCastlingRights))
	fmt.Println("")

	// We want to print a row at a time, but in backwards order to how this is stored
	// Since a1 occurs at the bottom left of the Board.
	for row := 7; row >= 0; row-- {
		newline := fmt.Sprintf("%v | ", row+1)
		// Get 8 consecutive Squares
		squares := b.Squares[8*row : 8*row+8]
		for i, piece := range squares {
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
		b.Squares[m.Capture.Square] = nil
		b.Position = UnSetPiece(b.Position, m.Capture.Piece, m.Capture.Square)
	}
	// Then, move the piece to its new square.
	// Check for promotion of a pawn.
	if m.Promotion != nil {
		b.PieceSet[m.Promotion] = s
		b.Squares[s] = m.Promotion
		b.Position = SetPiece(b.Position, m.Promotion, s)
		delete(b.PieceSet, p)
	} else {
		b.PieceSet[p] = s
		b.Squares[s] = p
		if p == nil {
			b.Print()
			panic("nil piece: " + m.String())
		}
		b.Position = SetPiece(b.Position, p, s)
	}
	// Check for castling and modify rook state if so.
	// New rook squares are relative to king.
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
		b.PieceSet[rook] = newRookSquare
		b.Squares[newRookSquare] = rook
		b.Position = SetPiece(b.Position, rook, newRookSquare)
		b.Squares[oldRookSquare] = nil
		b.Position = UnSetPiece(b.Position, rook, oldRookSquare)
	}
	// Then, remove the piece from its old square.
	b.Position = UnSetPiece(b.Position, p, m.Old)
	b.Squares[m.Old] = nil
	// Modify castling state from rook and king moves.
	// We know any piece moving from e8, e1, a8, h8, a1, or h8 must
	// change castling rights.
	switch m.Old {
	case A8:
		b.qsCastlingRights[BLACK] = false
	case H8:
		b.ksCastlingRights[BLACK] = false
	case E8:
		b.qsCastlingRights[BLACK] = false
		b.ksCastlingRights[BLACK] = false
	case A1:
		b.qsCastlingRights[WHITE] = false
	case H1:
		b.ksCastlingRights[WHITE] = false
	case E1:
		b.qsCastlingRights[WHITE] = false
		b.ksCastlingRights[WHITE] = false
	}
	// Affect castling state of captured rooks.
	if m.Capture != nil {
		if m.Capture.Piece.Type() == ROOK {
			if m.Capture.Square.Col() == 1 {
				b.qsCastlingRights[m.Capture.Piece.Color()] = false
			} else if m.Capture.Square.Col() == 8 {
				b.ksCastlingRights[m.Capture.Piece.Color()] = false
			}
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
	// Return the piece to its old square, and undo promotion.
	if m.Promotion != nil {
		delete(b.PieceSet, m.Promotion)
	}

	b.Squares[o] = p
	b.PieceSet[p] = o
	b.Position = SetPiece(b.Position, p, o)
	// Return the square we were on to its old state.
	b.Squares[s] = nil
	if m.Promotion != nil {
		b.Position = UnSetPiece(b.Position, m.Promotion, s)
	} else {
		b.Position = UnSetPiece(b.Position, p, s)
	}
	// Return a captured piece.
	if m.Capture != nil {
		b.PieceSet[m.Capture.Piece] = m.Capture.Square
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
		b.Squares[newRookSquare] = nil
		b.Position = UnSetPiece(b.Position, rook, newRookSquare)
		b.Squares[oldRookSquare] = rook
		b.PieceSet[rook] = oldRookSquare
		b.Position = SetPiece(b.Position, rook, oldRookSquare)
	}

	// Reapply original castling rights.
	b.qsCastlingRights = m.PrevQSCastlingRights
	b.ksCastlingRights = m.PrevKSCastlingRights

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
	for piece, _ := range b.PieceSet {
		if piece.Color() != b.Active {
			continue
		}
		m := piece.LegalMoves(b)
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
	var pcs uint64
	switch c {
	case WHITE:
		pcs = b.Position.WhitePieces
	case BLACK:
		pcs = b.Position.BlackPieces
	}

	for _, s := range SquaresFromBitBoard(pcs) {
		p := b.Squares[s]
		if p == nil {
			b.Print()
			fmt.Println(fmt.Sprintf("pieces on %v", SquaresFromBitBoard(pcs)))
			panic("nil piece on square " + s.String())
		}
		res = res | p.AttackBitboard(b, s)
	}
	return res
}

// Returns the hash value for this board.
func ZobristHash(b *Board) uint64 {
	var hash uint64
	hash = 0
	for p, s := range b.PieceSet {
		hash = hash ^ ZOBRISTPIECES[p.Color()][p.Type()][s]
	}
	if b.Active == WHITE {
		hash = hash ^ ZOBRISTTURN
	}
	if b.qsCastlingRights[WHITE] {
		hash = hash ^ ZOBRISTWQS
	}
	if b.qsCastlingRights[BLACK] {
		hash = hash ^ ZOBRISTBQS
	}
	if b.ksCastlingRights[WHITE] {
		hash = hash ^ ZOBRISTWKS
	}
	if b.ksCastlingRights[BLACK] {
		hash = hash ^ ZOBRISTBKS
	}
	return hash
}
