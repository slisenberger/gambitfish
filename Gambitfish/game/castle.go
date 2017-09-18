// Castle.go provides utilities for managing castling.
package game

var whiteQueensideSquares = []Square{
	Square{1, 4}, Square{1, 3}, Square{1, 2},
}
var blackQueensideSquares = []Square{
	Square{8, 4}, Square{8, 3}, Square{8, 2},
}
var whiteKingsideSquares = []Square{
	Square{1, 6}, Square{1, 7},
}
var blackKingsideSquares = []Square{
	Square{8, 6}, Square{8, 7},
}

func CanCastleQueenside(b *Board, c Color) bool {
	if !b.qsCastlingRights[c] {
		return false
	}
	if IsCheck(b, c) {
		return false
	}
	var castleSquares []Square
	switch c {
	case 1:
		castleSquares = whiteQueensideSquares
		break
	case -1:
		castleSquares = blackQueensideSquares
		break

	}
	if !CanCastleGeneric(b, c, castleSquares) {
		return false
	}
	return true
}

func CanCastleKingside(b *Board, c Color) bool {
	if !b.ksCastlingRights[c] {
		return false
	}
	if IsCheck(b, c) {
		return false
	}
	var castleSquares []Square
	switch c {
	case 1:
		castleSquares = whiteKingsideSquares
		break
	case -1:
		castleSquares = blackKingsideSquares
		break

	}
	if !CanCastleGeneric(b, c, castleSquares) {
		return false
	}
	return true
}

func CanCastleGeneric(b *Board, c Color, castleSquares []Square) bool {
	// Don't castle if there are pieces between.
	for _, cs := range castleSquares {
		if b.Squares[cs.Index()] != nil {
			return false
		}
	}
	attacking := GetAttacking(b, -1*c)
	for _, s := range attacking {
		for i, cs := range castleSquares {
			// If checking the queenside knight square, stop.
			if i >= 2 {
				break
			}
			// If any attacked square is our forbidden squares, return false.
			if s.Index() == cs.Index() {
				return false
			}

		}
	}

	return true

}

// Calculates legal castle moves in the game. Should only be called on kings.
func CastlingMoves(p *King, cur Square) []Move {
	moves := []Move{}
	if CanCastleQueenside(p.Board(), p.Color()) {
		s := Square{cur.Row, cur.Col - 2}
		move := NewMove(p, s, cur)
		move.QSCastle = true
		moves = append(moves, move)
	}
	if CanCastleKingside(p.Board(), p.Color()) {
		s := Square{cur.Row, cur.Col + 2}
		move := NewMove(p, s, cur)
		move.KSCastle = true
		moves = append(moves, move)
	}
	return moves
}
