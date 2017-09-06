// An implementation of a Knight
package game

type Knight struct {
	BasePiece
}

func (p Knight) LegalMoves() []Square {
	return nil
}

func (p Knight) String() string {
	switch p.color {
	case WHITE:
		return "N"
	case BLACK:
		return "n"
	}
	return ""
}
