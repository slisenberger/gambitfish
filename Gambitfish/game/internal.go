// Internal manages the internal datastructures useful by all parts of
// the engine, such as an opening book, zobrist hash codes, and legal
// move positions.
package game

// The preprocessed set of squares a piece can move to for a given
// index.
var LEGALKINGMOVES [64]uint64
var LEGALKNIGHTMOVES [64]uint64

func InitInternalData() {
	LEGALKINGMOVES = LegalKingMovesDict()
	LEGALKNIGHTMOVES = LegalKnightMovesDict()
}
