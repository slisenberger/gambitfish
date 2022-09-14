package game

const MAX_KILLER_MOVES = 2

type KillerMoves map[int][MAX_KILLER_MOVES]*Move

// Add a new Killergame.Move to the queue.
func (k KillerMoves) AddKillerMove(ply int, move *Move) {
	// Get the Killer game.Moves for this ply
	km, prs := k[ply]
	if !prs {
		km = [MAX_KILLER_MOVES]*Move{nil, nil}
	}
	// Store in the first killer move section, and move existing moves down.
	km[1] = km[0]
	km[0] = move
	k[ply] = km
}

func (k KillerMoves) GetKillerMoves(ply int) [2]*Move  {
	return k[ply]
}

func NewKillerMoves() KillerMoves {
	return make(map[int][MAX_KILLER_MOVES]*Move)
}
