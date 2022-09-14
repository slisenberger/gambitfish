package game

const MAX_KILLER_MOVES = 2

type KillerMoves map[int][MAX_KILLER_MOVES]EfficientMove

// Add a new Killergame.Move to the queue.
func (k KillerMoves) AddKillerMove(ply int, move EfficientMove) {
	// Get the Killer game.Moves for this ply
	km, prs := k[ply]
	if !prs {
		km = [MAX_KILLER_MOVES]EfficientMove{EfficientMove(0), EfficientMove(0)}
	}
	// Store in the first killer move section, and move existing moves down.
	km[1] = km[0]
	km[0] = move
	k[ply] = km
}

func (k KillerMoves) GetKillerMoves(ply int) [2]EfficientMove  {
	return k[ply]
}

func NewKillerMoves() KillerMoves {
	return make(map[int][MAX_KILLER_MOVES]EfficientMove)
}
