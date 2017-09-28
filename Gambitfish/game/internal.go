// Internal manages the internal datastructures useful by all parts of
// the engine, such as an opening book, zobrist hash codes, and legal
// move positions.
package game

// The preprocessed set of squares a piece can move to for a given
// index.
var LEGALKINGMOVES [64]uint64
var LEGALKNIGHTMOVES [64]uint64

// The directions map to integers for directional movement on a chessboard.
type Direction int

const (
	NE = Direction(9)
	N  = Direction(8)
	NW = Direction(7)
	E  = Direction(1)
	W  = Direction(-1)
	SE = Direction(-7)
	S  = Direction(-8)
	SW = Direction(-9)
)

// The preprocessed bitboard of squares a ray piece can move to
// in a given direction. Each of the 64 squares for each direction
// has a bitboard of attacking squares.
var RAY_ATTACKS map[Direction][64]uint64

func InitInternalData() {
	LEGALKINGMOVES = LegalKingMovesDict()
	LEGALKNIGHTMOVES = LegalKnightMovesDict()
	InitRayAttacks()
}

// InitRayAttacks initializes the set of bitboards for ray movements
// in directions.
func InitRayAttacks() {
	RAY_ATTACKS = map[Direction][64]uint64{}
	dirs := []Direction{NE, N, NW, E, W, SE, S, SW}
	for _, dir := range dirs {
		RAY_ATTACKS[dir] = [64]uint64{}
		// Create an entry for each square.
		var bb uint64
		var cur uint64
		var i uint64
		var j uint64
		vectors := [64]uint64{}
		for i = 0; i < 64; i++ {
			bb = 0
			cur = 1 << i
			// We will extend in each direction at most 7 times.
			for j = 1; j <= 7; j++ {
				if dir > 0 {
					cur = cur << uint64(dir)
					s := Square(i + j*uint64(dir))
					// Since positive directions can wrap around, we might nee dto end the loop here.
					// Shouldn't be on the first column.
					if dir == NE || dir == E {
						// Find the square we are on.
						if s.Col() == 1 {
							break
						}
						// For NW, we shouldn't be on the 8th column.
					} else {
						if s.Col() == 8 {
							break
						}
					}
				} else {
					cur = cur >> uint64(-dir)
				}
				bb = bb | cur
			}
			vectors[i] = bb
		}
		RAY_ATTACKS[dir] = vectors
	}
}
