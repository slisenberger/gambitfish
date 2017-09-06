// Square is a convenience type for representing squares on the board.
package game

import "fmt"

type Square struct {
	row, col int
}

// Prints this square as a string.
func (s *Square) String() string {
	return fmt.Sprintf("%v%v", string('a'+s.col-1), s.row)
}

// Returns an index for this square in a one-dimensional array.
func (s *Square) Index() int {
	return 8*(s.row-1) + s.col - 1
}

// InPlay returns true if a square is on the 8 by 8 chess board.
func (s *Square) InPlay() bool {
	return s.row >= 1 && s.row <= 8 && s.col >= 1 && s.col <= 8
}
