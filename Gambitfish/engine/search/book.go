// book.go is a way of applying book lines to moves.
package search

import "../../game"

// BookMove will return a new book move for the given board.
// For now, it simply hardcodes a SLAV for the first two moves.
func BookMove(b *game.Board) *game.Move {
	// This is a hack to have the slav.
	if b.Move == 1 && b.Active == game.WHITE {
		// 1.d4
		d2 := game.Square{2, 4}
		d4 := game.Square{4, 4}
		move := game.NewMove(b.Squares[d2.Index()], d4, d2)
		return &move
	}
	if b.Move == 1 && b.Active == game.BLACK {
		// 1.c6
		c7 := game.Square{7, 3}
		c6 := game.Square{6, 3}
		move := game.NewMove(b.Squares[c7.Index()], c6, c7)
		return &move
	}
	if b.Move == 2 && b.Active == game.WHITE {
		// 1.c4
		c2 := game.Square{2, 3}
		c4 := game.Square{4, 3}
		move := game.NewMove(b.Squares[c2.Index()], c4, c2)
		return &move
	}
	if b.Move == 2 && b.Active == game.BLACK {
		// 1.d5
		d7 := game.Square{7, 4}
		d5 := game.Square{5, 4}
		move := game.NewMove(b.Squares[d7.Index()], d5, d7)
		return &move
	}
	return nil
}
