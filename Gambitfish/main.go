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
			game.PieceSquareEvaluator{},
			// Calculating legal moves may be slowing this down. 
			//game.MobilityEvaluator{},
			// game.KingSafetyEvaluator{},
			// game.OpeningEvaluator{},
		},
	}
       //	p1 := player.CommandLinePlayer{Color: game.WHITE}
	p1 := player.AIPlayer{Evaluator: e, Depth: 10, Color: game.WHITE}
	p2 := player.AIPlayer{Evaluator: e, Depth: 8, Color: game.BLACK}
	b.Print()
	for i := 0; i < 300; i++ {
		//time.Sleep(1 * time.Second)
		if over, winner := b.CalculateGameOver(b.AllLegalMoves()); over {
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

		// Every 4 turns flush the table of entries that haven't been used.
		if i % 4 == 0 {
			game.EraseOldTableEntries()
                }

		fmt.Println("new board: ")
		b.Print()
	}
}
