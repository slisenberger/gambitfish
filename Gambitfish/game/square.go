// Square is a convenience type for representing squares on the board.
package game

import "fmt"

type Square struct {
	Row, Col int
}

// Prints this square as a string.
func (s Square) String() string {
	return fmt.Sprintf("%v%v", string('a'+s.Col-1), s.Row)
}

// Returns an index for this square in a one-dimensional array.
func (s Square) Index() int {
	return 8*(s.Row-1) + s.Col - 1
}

// InPlay returns true if a square is on the 8 by 8 chess board.
func (s Square) InPlay() bool {
	return s.Row >= 1 && s.Row <= 8 && s.Col >= 1 && s.Col <= 8
}

// SquareFromIndex returns a new square from an index into a single dimension array.
func SquareFromIndex(i int) Square {
	return Square{Row: i/8 + 1, Col: (i % 8) + 1}
}
