// book.go is a way of applying book lines to moves.
package search

import "../../game"

// BookMove will return a new book move for the given board.
// For now, it simply hardcodes a SLAV for the first two moves.
func BookMove(b *game.Board) *game.Move {
	// This is a hack to have the slav.
	if b.Move == 1 && b.Active == game.WHITE {
		// 1.d4
		move := game.NewMove(b.Squares[game.D2], game.D4, game.D2, b)
		return &move
	}
	if b.Move == 1 && b.Active == game.BLACK {
		// 1.c6
		move := game.NewMove(b.Squares[game.C7], game.C6, game.C7, b)
		return &move
	}
	if b.Move == 2 && b.Active == game.WHITE {
		// 2.c4
		move := game.NewMove(b.Squares[game.C2], game.C4, game.C2, b)
		return &move
	}
	if b.Move == 2 && b.Active == game.BLACK {
		// 2.d5
		move := game.NewMove(b.Squares[game.D7], game.D5, game.D7, b)
		return &move
	}
	return nil
}
