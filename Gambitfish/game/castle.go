// Castle.go provides utilities for managing castling.
package game

var whiteQueensideSquares = []Square{
	GetSquare(1, 4), GetSquare(1, 3), GetSquare(1, 2),
}
var blackQueensideSquares = []Square{
	GetSquare(8, 4), GetSquare(8, 3), GetSquare(8, 2),
}
var whiteKingsideSquares = []Square{
	GetSquare(1, 6), GetSquare(1, 7),
}
var blackKingsideSquares = []Square{
	GetSquare(8, 6), GetSquare(8, 7),
}

func CanCastleQueenside(b *Board, c Color) bool {
	if !b.qsCastlingRights[c] {
		return false
	}
	if IsCheck(b, c) {
		return false
	}
	var castleSquares []Square
	var rookSquare Square
	switch c {
	case 1:
		castleSquares = whiteQueensideSquares
		rookSquare = GetSquare(1, 1)
	case -1:
		castleSquares = blackQueensideSquares
		rookSquare = GetSquare(8, 1)
	}
	if !CanCastleGeneric(b, c, castleSquares, rookSquare) {
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
	var rookSquare Square
	switch c {
	case 1:
		castleSquares = whiteKingsideSquares
		rookSquare = GetSquare(1, 8)
	case -1:
		castleSquares = blackKingsideSquares
		rookSquare = GetSquare(8, 8)
	}
	if !CanCastleGeneric(b, c, castleSquares, rookSquare) {
		return false
	}
	return true
}

func CanCastleGeneric(b *Board, c Color, castleSquares []Square, rookSquare Square) bool {
	// Don't castle if there are pieces between.
	for _, cs := range castleSquares {
		if b.Squares[cs] != nil {
			return false
		}
	}
	// Make sure there's a rook on our target square.
	r := b.Squares[rookSquare]
	if r == nil {
		return false
	}
	if r.Type() != ROOK {
		return false
	}

	attacking := GetAttacking(b, -1*c)
	for _, s := range attacking {
		for i, cs := range castleSquares {
			// If checking the queenside knight square, stop.
			if i >= 2 {
				break
			}
			// If any attacked square is our forbidden squares, return false.
			if s == cs {
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
		s := GetSquare(cur.Row(), cur.Col()-2)
		move := NewMove(p, s, cur)
		move.QSCastle = true
		moves = append(moves, move)
	}
	if CanCastleKingside(p.Board(), p.Color()) {
		s := GetSquare(cur.Row(), cur.Col()+2)
		move := NewMove(p, s, cur)
		move.KSCastle = true
		moves = append(moves, move)
	}
	return moves
}
