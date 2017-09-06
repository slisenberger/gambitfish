// Package main initializes the necessary game state and prints the
// current state of the art of the engine.
package main

import "./game"
import "fmt"
import "math/rand"

func main() {
	b := game.DefaultBoard()
	b.Print()
	for _, piece := range b.Squares {
		moves := piece.LegalMoves()
		if moves != nil && len(moves) > 0 {
			for _, move := range moves {

				fmt.Println(fmt.Printf("legal move: piece %v from %v to %v", piece, piece.Square().String(), move.String()))

			}
		}
	}
}
