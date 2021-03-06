// Package main initializes the necessary game state and prints the
// current state of the art of the engine.
package main

import "./game"
import "./player"
import "fmt"
import "log"
import "math/rand"
import "runtime/pprof"
import "os"
import "time"

func main() {
	f, err := os.Create("pprof.cpu")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	rand.Seed(time.Now().Unix())
	game.InitInternalData()
	b := game.DefaultBoard()
	e := game.CompoundEvaluator{
		Evaluators: []game.Evaluator{
			game.MaterialEvaluator{},
			// game.PieceSquareEvaluator{},
			// We'll turn this on when I like it
			// MobilityEvaluator{},

			// game.OpeningEvaluator{},
			// game.KingSafetyEvaluator{},
		},
	}
	p1 := player.CommandLinePlayer{Color: game.WHITE}
	p2 := player.AIPlayer{Evaluator: e, Depth: 4, Color: game.BLACK}
	b.Print()
	for i := 0; i < 300; i++ {
		time.Sleep(1 * time.Second)
		if over, winner := b.CalculateGameOver(); over {
			if winner != 0 {
				fmt.Println(fmt.Sprintf("WINNER: %v in %v moves", winner, b.Move))
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
