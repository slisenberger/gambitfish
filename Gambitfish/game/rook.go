// An implementation of a Rook
package game

type Rook struct {
	BasePiece
}

func (p Rook) LegalMoves() []Square {
	return nil
}

func (p Rook) String() string {
	switch p.color {
	case WHITE:
		return "R"
	case BLACK:
		return "r"
	}
	return ""
}
