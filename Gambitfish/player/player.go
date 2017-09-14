package player

import "errors"
import "fmt"
import "math"
import "../game"
import "../engine/evaluate"
import "../engine/search"

type Player interface {
	MakeMove(*game.Board) error
}

// AIPlayer is a player that makes moves according to AI.
type AIPlayer struct {
	Evaluator evaluate.Evaluator
	Depth     int
	Color     game.Color
}

func (p *AIPlayer) MakeMove(b *game.Board) error {
	eval, move := search.AlphaBetaSearch(b, p.Evaluator, p.Depth, math.Inf(-1), math.Inf(1))
	if move == nil {
		return errors.New("no move could be made")
	}
	fmt.Println(fmt.Sprintf("AI Player making best move with depth %v: %v, eval %v", p.Depth, move, eval))

	game.ApplyMove(b, *move)
	return nil
}
