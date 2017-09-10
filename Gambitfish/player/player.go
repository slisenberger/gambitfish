package player

import "errors"
import "fmt"
import "../game"
import "../engine/evaluate"
import "../engine/search"

type Player interface {
	MakeMove(*game.Board) error
}

// AIPlayer is a player that makes moves according to AI.
type AIPlayer struct {
	Evaluator evaluate.Evaluator
	Color     game.Color
}

func (p *AIPlayer) MakeMove(b *game.Board) error {
	move := search.AlphaBetaSearch(b, p.Evaluator, 1, 0, 0)
	if move == nil {
		return errors.New("no move could be made")
	}
	fmt.Println("AI Player making best move: " + move.String())

	b.ApplyMove(*move)
	return nil
}
