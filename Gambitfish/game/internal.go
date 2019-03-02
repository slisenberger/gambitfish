// Internal manages the internal datastructures useful by all parts of
// the engine, such as an opening book, zobrist hash codes, and legal
// move positions.
package game

import "math/rand"
import "time"

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
var RAY_ATTACKS_N [64]uint64
var RAY_ATTACKS_NW [64]uint64
var RAY_ATTACKS_NE [64]uint64
var RAY_ATTACKS_E [64]uint64
var RAY_ATTACKS_W [64]uint64
var RAY_ATTACKS_SE [64]uint64
var RAY_ATTACKS_S [64]uint64
var RAY_ATTACKS_SW [64]uint64

// TODO(slisenberger): We want to use rotated bitboards. Here's what
// we need:
// 1.) Prepopulate, for all 64 possible squares, for all possible
// occupancies on their rank, what the vector of movements is.
// 2.) Do the same for diagonals.
// 3.) When calculating sliding attacks, look up the

// The preprocessed set of squares a pawn can move to in a capture.
var WHITEPAWNATTACKS [64]uint64
var BLACKPAWNATTACKS [64]uint64

// The preprocessed random numbers to be used for zobrist hash keys.
// It is a map of color-> piece -> square
var ZOBRISTPIECES map[Color]map[PieceType][64]uint64

// A random number indicating turn to move.
var ZOBRISTTURN uint64

// Random numbers for castling rights.
var ZOBRISTWKS uint64
var ZOBRISTWQS uint64
var ZOBRISTBKS uint64
var ZOBRISTBQS uint64

// TODO(slisenberger): include en passant in zobrist.

func InitInternalData() {
	LEGALKINGMOVES = LegalKingMovesDict()
	LEGALKNIGHTMOVES = LegalKnightMovesDict()
	InitRayAttacks()
	InitPawnAttacks()
	InitZobristNumbers()
	BuildByteLookupTable()
}

// LegalKingMovesDict returns a 64-indexed set of bitboards
// for the legal moves a king can make. This is intended to be used
// as a precomputed index for constructing valid king moves.
func LegalKingMovesDict() [64]uint64 {
	var legalMoves [64]uint64
	// For each square on the board, prepopulate the legal moves.
	var i uint64
	var bb uint64
	var kp uint64
	for i = 0; i < 64; i++ {
		// king position
		kp = 1 << i
		bb = 0
		s := Square(i)
		// Check for going off the left side
		if s.Col() != 1 {
			bb = bb | (kp << 7) | (kp >> 9) | (kp >> 1)
		}
		// Check for going off the right side
		if s.Col() != 8 {
			bb = bb | (kp >> 7) | (kp << 9) | (kp << 1)
		}
		bb = bb | (kp >> 8) | (kp << 8)
		legalMoves[i] = bb
	}
	return legalMoves
}

// LegalKnightMovesDict returns a 64-indexed set of bitboards
// for the legal moves a knight can make. This is intended to be used
// as a precomputed index for detecting legal knight moves.
func LegalKnightMovesDict() [64]uint64 {
	var legalMoves [64]uint64
	// For each square on the board, prepopulate the legal moves.
	var i uint64
	var bb uint64
	var kp uint64
	for i = 0; i < 64; i++ {
		// knight position
		kp = 1 << i
		bb = 0
		s := Square(i)
		// Check for going off the left side
		if s.Col() != 1 {
			bb = bb | (kp << 15) | (kp >> 17)
		}
		if s.Col() > 2 {
			bb = bb | (kp << 6) | (kp >> 10)
		}
		// Check for going off the right side
		if s.Col() < 7 {
			bb = bb | (kp >> 6) | (kp << 10)
		}
		if s.Col() != 8 {
			bb = bb | (kp >> 15) | (kp << 17)
		}
		legalMoves[i] = bb
	}
	return legalMoves
}

// InitRayAttacks initializes the set of bitboards for ray movements
// in directions.
func InitRayAttacks() {
	var ra [64]uint64
	dirs := []Direction{NE, N, NW, E, W, SE, S, SW}
	for _, dir := range dirs {
		ra = [64]uint64{}
		// Create an entry for each square.
		var bb uint64
		var cur uint64
		var i uint64
		var j uint64
		for i = 0; i < 64; i++ {
			bb = 0
			cur = 1 << i
			// We will extend in each direction at most 7 times.
			for j = 1; j <= 7; j++ {
				s := Square(i + j*uint64(dir))
				if dir > 0 {
					cur = cur << uint64(dir)
					// Since positive directions can wrap around, we might nee dto end the loop here.
					// Shouldn't be on the first column.
					if dir == NE || dir == E {
						// Find the square we are on.
						if s.Col() == 1 {
							break
						}
						// For NW, we shouldn't be on the 8th column.
					} else if dir == NW {
						if s.Col() == 8 {
							break
						}
					}
				} else {
					cur = cur >> uint64(-dir)
					// It's also possible for negative directions to wrap around.
					// Western movement should never be on col 8
					if dir == W || dir == SW {
						if s.Col() == 8 {
							break
						}
					} else if dir == SE {
						if s.Col() == 1 {
							break
						}

					}
				}
				bb = bb | cur
			}
			ra[i] = bb
		}
		switch dir {
		case N:
			RAY_ATTACKS_N = ra
		case NE:
			RAY_ATTACKS_NE = ra
		case NW:
			RAY_ATTACKS_NW = ra
		case S:
			RAY_ATTACKS_S = ra
		case SE:
			RAY_ATTACKS_SE = ra
		case SW:
			RAY_ATTACKS_SW = ra
		case E:
			RAY_ATTACKS_E = ra
		case W:
			RAY_ATTACKS_W = ra
		}
	}
}

func InitPawnAttacks() {
	var i uint64
	var wbb uint64
	var bbb uint64
	var pos uint64
	for i = 0; i < 64; i++ {
		wbb = 0
		bbb = 0
		pos = 1 << i
		s := Square(i)
		if s.Col() > 1 {
			wbb = wbb | (pos << 7)
			bbb = bbb | (pos >> 9)
		}
		if s.Col() < 8 {
			wbb = wbb | (pos << 9)
			bbb = bbb | (pos >> 7)
		}
		WHITEPAWNATTACKS[i] = wbb
		BLACKPAWNATTACKS[i] = bbb
	}
}

// Initializes the set of random numbers necessary to hash positions.
// See https://chessprogramming.wikispaces.com/Zobrist+Hashing
func InitZobristNumbers() {
	rand.Seed(time.Now().UTC().UnixNano())
	ZOBRISTPIECES = map[Color]map[PieceType][64]uint64{}
	colors := []Color{WHITE, BLACK}
	pieces := []PieceType{PAWN, BISHOP, KNIGHT, ROOK, QUEEN, KING}
	for _, c := range colors {
		piecemap := map[PieceType][64]uint64{}
		for _, p := range pieces {
			squares := [64]uint64{}
			for i := 0; i < 63; i++ {
				squares[i] = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())

			}
			piecemap[p] = squares
		}
		ZOBRISTPIECES[c] = piecemap
	}
	ZOBRISTTURN = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	ZOBRISTWKS = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())

	ZOBRISTWQS = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())

	ZOBRISTBKS = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())

	ZOBRISTBQS = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())

}

// Returns the ray attack bitboard for a given direction and square
func RayAttacks(d Direction, s Square) uint64 {
	switch d {
	case N:
		return RAY_ATTACKS_N[s]
	case NE:
		return RAY_ATTACKS_NE[s]
	case NW:
		return RAY_ATTACKS_NW[s]
	case S:
		return RAY_ATTACKS_S[s]
	case SE:
		return RAY_ATTACKS_SE[s]
	case SW:
		return RAY_ATTACKS_SW[s]
	case W:
		return RAY_ATTACKS_W[s]
	case E:
		return RAY_ATTACKS_E[s]
	}
	return 0

}
