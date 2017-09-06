// Piece is an interface that defines the operations possible for a piece on the board.
package game

// Define the possible colors of a piece as an enum
type Color int

const (
	WHITE Color = 0
	BLACK Color = 1
)

type Piece interface {
	// Returns an array of all the legal positions this piece can move to.
	LegalMoves() []Square
	// Returns a string representation of this piece.
	String() string
}

type BasePiece struct {
	color  Color
	square Square
}
