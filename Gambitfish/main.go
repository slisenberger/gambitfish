// Package main initializes the necessary game state and prints the
// current state of the art of the engine.
package main

import "./game"
import "fmt"
import "math/rand"
import "time"

func main() {
	rand.Seed(time.Now().Unix())
	b := game.DefaultBoard()
	b.Print()
	for i := 0; i < 300; i++ {
		if b.Winner != 0 {
			fmt.Println("WINNER")
			break
		}
		moves := b.AllLegalMoves()
		if len(moves) == 0 {
			fmt.Println("GAME ends in STALEMATE! no legal moves!")
			break
		}
		for _, move := range moves {
			fmt.Println("legal move:" + move.String())
		}

		fmt.Println("selecting move at random..\n")
		n := rand.Int() % len(moves)
		fmt.Println("selected move: " + moves[n].String())
		b.ApplyMove(moves[n])

		fmt.Println("new board: ")
		b.Print()
	}
}
