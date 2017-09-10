package player

type Player interface {
	MakeMove(*Board) error
}

// AIPlayer is a player that makes moves according to AI.
type AIPlayer struct {
	Evaluator evaluate.Evaluator
	Color     game.Color
}

func (p *AIPlayer) MakeMove(b *Board) error {
	move := search.AlphaBetaSearch(b, p.Evaluator, 0, 0, 0)
	b.ApplyMove(move)
	return nil
}
