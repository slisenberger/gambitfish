package evaluate

import "../../game"

type MaterialEvaluator struct{}

func (m MaterialEvaluator) Evaluate(b *game.Board) float {
	eval := 0
	for _, piece := range b.Squares() {
		if piece == nil {
			continue
		}
		if piece.Color() == b.active {
			eval += GetPieceValue(piece)
		} else {
			eval -= GetPieceValue(piece)
		}

	}

}

func GetPieceValue(p game.Piece) float {
	switch p.(type) {
	case game.Pawn:
		return 1
		break
	case game.Bishop:
		return 3
		break
	case game.Knight:
		return 3
		break
	case game.Rook:
		return 5
		break
	case game.King:
		return 100
		break
	case game.Queen:
		return 9
		break
	}
}
