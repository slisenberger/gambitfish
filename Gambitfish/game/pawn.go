// An implementation of a Pawn.
package game

type Pawn struct {
	BasePiece
}

func (p Pawn) LegalMoves() []Square {
	return nil
}

func (p Pawn) String() string {
	switch p.color {
	case WHITE:
		return "P"
	case BLACK:
		return "p"
	}
	return ""
}
