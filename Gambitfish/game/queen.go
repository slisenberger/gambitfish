// An implementation of a Queen.
package game

type Queen struct {
	BasePiece
}

func (p Queen) LegalMoves() []Square {
	return nil
}

func (p Queen) String() string {
	switch p.color {
	case WHITE:
		return "Q"
	case BLACK:
		return "q"
	}
	return ""
}
