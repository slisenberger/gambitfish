// Internal manages the internal datastructures useful by all parts of
// the engine, such as an opening book, zobrist hash codes, and legal
// move positions.
package game

// LEGALKINGMOVES is an array of bitboards for all the legal king
// moves for each square on the board.
var LEGALKINGMOVES [64]uint64

func InitInternalData() {
	LEGALKINGMOVES = LegalKingMovesDict()
}
