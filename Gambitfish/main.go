// Package main initializes the necessary game state and prints the
// current state of the art of the engine.
package main

import "./game"
import "./engine/evaluate"
import "./player"
import "fmt"
import "math/rand"
import "time"

func main() {
	rand.Seed(time.Now().Unix())
	b := game.DefaultBoard()
	e := evaluate.MaterialEvaluator{}
	p1 := player.AIPlayer{Evaluator: e, Depth: 1, Color: game.WHITE}
	p2 := player.AIPlayer{Evaluator: e, Depth: 1, Color: game.BLACK}
	b.Print()
	for i := 0; i < 300; i++ {
		time.Sleep(1 * time.Second)
		if over, winner := b.CalculateGameOver(); over {
			if winner != 0 {
				fmt.Println("WINNER")
			} else {
				fmt.Println("GAME ends in STALEMATE! no legal moves!")
			}
			break
		}
		if b.Active == p1.Color {
			p1.MakeMove(b)
		} else {
			p2.MakeMove(b)
		}
		b.SwitchActivePlayer()
		fmt.Println("new board: ")
		b.Print()
	}
}
