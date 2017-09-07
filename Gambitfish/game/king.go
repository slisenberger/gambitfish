// An implementation of a King.
package game

type King struct {
	BasePiece
	HasMoved bool
}

func (p King) LegalMoves() []Square {
	return p.KingMoves()
}

func (p King) String() string {
	switch p.color {
	case WHITE:
		return "K"
	case BLACK:
		return "k"
	}
	return ""
}
