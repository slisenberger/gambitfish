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
	for i := 0; i < 20; i++ {
		moves := b.AllLegalMoves()
		for _, move := range moves {
			fmt.Println("legal move:" + move.String())
		}

		fmt.Println("selecting move at random../n/n")
		n := rand.Int() % len(moves)
		b.ApplyMove(moves[n])

		fmt.Println("selected move: " + moves[n].String())
		fmt.Println("new board: ")
		b.Print()
	}
}
