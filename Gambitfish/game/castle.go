// Castle.go provides utilities for managing castling.
package game

// bitboards representing squares relevant to castling decisions.
var wksMustBeUnoccupied = uint64(0x0000000000000060) // bit for bishop and knight
var wqsMustBeUnoccupied = uint64(0x000000000000000E) // bit for bishop, knight, queen
var bksMustBeUnoccupied = uint64(0x6000000000000000)
var bqsMustBeUnoccupied = uint64(0x0E00000000000000) // bit for bishop, knight, queen
var wksCantBeAttacked = uint64(0x0000000000000070)   // bit for king, bishop, knight
var bksCantBeAttacked = uint64(0x7000000000000000)   // bit for king, bishop, knight
var wqsCantBeAttacked = uint64(0x000000000000001C)   // bit for king, queen, bishop
var bqsCantBeAttacked = uint64(0x1C00000000000000)   // bit for bishop, queen, king

func CanCastleQueenside(b *Board, c Color) bool {
	switch c {
	case BLACK:
		if !b.BQSCastling {
			return false
		}
	case WHITE:
		if !b.WQSCastling {
			return false
		}
	}
	var unoccupied uint64
	var kingSlide uint64
	var rookSquare Square
	switch c {
	case WHITE:
		unoccupied = wqsMustBeUnoccupied
		kingSlide = wqsCantBeAttacked
		rookSquare = A1
	case BLACK:
		unoccupied = bqsMustBeUnoccupied
		kingSlide = bqsCantBeAttacked
		rookSquare = A8
	}
	return CanCastleGeneric(b, c, unoccupied, kingSlide, rookSquare)
}

func CanCastleKingside(b *Board, c Color) bool {
	switch c {
	case BLACK:
		if !b.BKSCastling {
			return false
		}
	case WHITE:
		if !b.WKSCastling {
			return false
		}
	}
	var rookSquare Square
	var unoccupied uint64
	var kingSlide uint64
	switch c {
	case WHITE:
		unoccupied = wksMustBeUnoccupied
		kingSlide = wksCantBeAttacked
		rookSquare = H1
	case BLACK:
		unoccupied = bksMustBeUnoccupied
		kingSlide = bksCantBeAttacked
		rookSquare = H8
	}
	return CanCastleGeneric(b, c, unoccupied, kingSlide, rookSquare)

}

func CanCastleGeneric(b *Board, c Color, castleOccupancy uint64, kingSlide uint64, rookSquare Square) bool {
	// Make sure there's a rook on our target square.
	r := b.Squares[rookSquare]
	if r == NULLPIECE {
		return false
	}
	if r.Type() != ROOK {
		return false
	}
	// Don't castle if there are pieces between.
	if b.Position.Occupied&castleOccupancy > 0 {
		return false
	}
	atk := GetAttackBitboard(b, -1*c)
	// Don't castle if any king square is under attack.
	if atk&kingSlide > 0 {
		return false
	}

	return true

}

// Calculates legal castle moves in the game. Should only be called on kings.
func CastlingMoves(b *Board, p Piece, cur Square) []EfficientMove {
	moves := []EfficientMove{}
	if CanCastleQueenside(b, p.Color()) {
		s := GetSquare(cur.Row(), cur.Col()-2)
		move := NewEfficientMove(p, s, cur)
		move.AddQSCastle()
		moves = append(moves, move)
	}
	if CanCastleKingside(b, p.Color()) {
		s := GetSquare(cur.Row(), cur.Col()+2)
		move := NewEfficientMove(p, s, cur)
		move.AddKSCastle()
		moves = append(moves, move)
	}
	return moves
}
