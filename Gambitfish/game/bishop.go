// An implementation of a Bishop
package game

type Bishop struct {
	BasePiece
}

func (p Bishop) LegalMoves() []Square {
	return nil
}

func (p Bishop) String() string {
	switch p.color {
	case WHITE:
		return "B"
	case BLACK:
		return "b"
	}
	return ""
}
