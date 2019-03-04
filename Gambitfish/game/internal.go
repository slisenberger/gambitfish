// Internal manages the internal datastructures useful by all parts of
// the engine, such as an opening book, zobrist hash codes, and legal
// move positions.
package game

import "math/rand"
import "math/bits"
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

// The preprocessed conversion between a normally indexed bitboard and
// the corresponding index on a rotated bitboard.
var ROT_90_LEFT_CONVERSION [64]uint64
var ROT_45_RIGHT_CONVERSION [64]uint64
var ROT_45_LEFT_CONVERSION [64]uint64

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

// The precomputed relevant occupancies to determine
// blockers for ray attacks.
var BLOCKERMASKBISHOP [64]uint64
var BLOCKERMASKROOK [64]uint64

// The hash key for a particular square for rooks/bishops
// in looking up precomputed moves.
var MAGICNUMBERBISHOP [64]uint64
var MAGICNUMBERROOK [64]uint64

// The precomputed attack squares for bishops and rooks. Each
// Square has a table, and each of those has a hash table indexed
// By a magic number computation.
var BISHOPATTACKS [64][]uint64
var ROOKATTACKS [64][]uint64

// TODO(slisenberger): include en passant in zobrist.

func InitInternalData() {
	// This needs to happen before magic bitboard generation
	// because we use this as a convenience for building attack
	// boards.
	InitRayAttacks()
	InitBlockerMasks()
	InitMagicBitboards()
	LEGALKINGMOVES = LegalKingMovesDict()
	LEGALKNIGHTMOVES = LegalKnightMovesDict()
	InitPawnAttacks()
	InitZobristNumbers()
	BuildByteLookupTable()
}

// In order to use rotated bitboards, we'll
// need a way to translate our usual coordinate
// system to find the index in a rotated bitboard.
// These maps convert from a given square to its
// rotated equivalent.
func InitRotatedBitboardConversions() {
	ROT_90_LEFT_CONVERSION = [64]uint64{
		7, 15, 23, 31, 39, 47, 55, 63,
		6, 14, 22, 30, 38, 46, 54, 62,
		5, 13, 21, 29, 37, 45, 53, 61,
		4, 12, 20, 28, 36, 44, 52, 60,
		3, 11, 19, 27, 35, 43, 51, 59,
		2, 10, 18, 26, 34, 42, 50, 58,
		1, 9, 17, 25, 33, 41, 49, 57,
		0, 8, 16, 24, 32, 40, 48, 56,
	}
	ROT_45_LEFT_CONVERSION = [64]uint64{
		0,
		8, 1,
		16, 9, 2,
		24, 17, 10, 3,
		32, 25, 18, 11, 4,
		40, 33, 26, 19, 12, 5,
		48, 41, 34, 27, 20, 13, 6,
		56, 49, 42, 35, 28, 21, 14, 7,
		57, 50, 43, 36, 29, 22, 15,
		58, 51, 44, 37, 30, 23,
		59, 52, 45, 38, 31,
		60, 53, 46, 39,
		61, 54, 47,
		62, 55,
		63,
	}

	ROT_45_RIGHT_CONVERSION = [64]uint64{
		7,
		6, 15,
		5, 14, 23,
		4, 13, 22, 31,
		3, 12, 21, 30, 39,
		2, 11, 20, 29, 38, 47,
		1, 10, 19, 28, 37, 46, 55,
		0, 9, 18, 27, 36, 45, 54, 63,
		8, 17, 26, 35, 44, 53, 62,
		16, 25, 34, 43, 52, 61,
		24, 33, 42, 51, 60,
		32, 41, 50, 59,
		40, 49, 58,
		48, 57,
		56,
	}
}

// Generates the set of important bits for a square
// that can block rook movement. Returns a bitboard with
// the relevant occupancy. We ignore squares on the edges of
// the board, since they can't block further movement.
func BlockerMaskRook(s Square) uint64 {
	var bb uint64
	bb = 0
	for col := 2; col < 8; col++ {
		if col == s.Col() {
			continue
		}
		bb = SetBitOnBoard(bb, GetSquare(s.Row(), col))
	}
	for row := 2; row < 8; row++ {
		if row == s.Row() {
			continue
		}
		bb = SetBitOnBoard(bb, GetSquare(row, s.Col()))
	}
	return bb
}

// Generates the set of important bits for a bishop that
// can block bishop movement. We ignore squares on the edges
// of the board, since they can't block further movement.
func BlockerMaskBishop(s Square) uint64 {
	var bb uint64
	bb = 0
	for col := 2; col < 8; col++ {
		deltaCol := col - s.Col()
		// If it's our column, this is our square, it's irrelevant
		// for blockers.
		if deltaCol == 0 {
			continue
		}
		// For all other columns, look for the row above and below, and
		// set these on the board.
		if s.Row()+deltaCol >= 1 && s.Row()+deltaCol <= 8 {
			bb = SetBitOnBoard(bb, GetSquare(s.Row()+deltaCol, col))

		}
		if s.Row()-deltaCol >= 1 && s.Row()-deltaCol <= 8 {
			bb = SetBitOnBoard(bb, GetSquare(s.Row()+deltaCol, col))
		}
	}
	return bb
}

func InitBlockerMasks() {
	var i uint64
	BLOCKERMASKBISHOP = [64]uint64{}
	BLOCKERMASKROOK = [64]uint64{}
	for i = 0; i < 64; i++ {
		BLOCKERMASKBISHOP[i] = BlockerMaskBishop(Square(i))
		BLOCKERMASKROOK[i] = BlockerMaskRook(Square(i))
	}
}

// Initialize the preset arrays containing all possible
// boards that could block rook or bishop movement, and the
// resulting moves that can be made because of them.
func InitMagicBitboards() {
	// Store the magic number keys.
	MAGICNUMBERROOK = [64]uint64{}
	MAGICNUMBERBISHOP = [64]uint64{}
	ROOKATTACKS = [64][]uint64{}
	BISHOPATTACKS = [64][]uint64{}
	// Get all the possible blockerboards, and their resulting
	// legal moves.

	// And find the magics for all the squares.
	var i uint64
	var lookingForNumber bool
	var magic uint64
	// Loop once for rooks.
	for i = 0; i < 64; i++ {
		lookingForNumber = true
		bitCount := bits.OnesCount64(BLOCKERMASKROOK[i])
		rookAttacks := make([]uint64, 1<<uint64(bitCount))
		for lookingForNumber {
			magic = rand.Uint64()

			// Stop looking for magic numbers if we found one for this square.
			lookingForNumber = false
		}
		ROOKATTACKS[i] = rookAttacks
		MAGICNUMBERROOK[i] = magic
	}

	// And again for bishops.
	for i = 0; i < 64; i++ {
		bitCount := bits.OnesCount64(BLOCKERMASKBISHOP[i])
		bishopAttacks := make([]uint64, 1<<uint64(bitCount))
		lookingForNumber = true
		for lookingForNumber {
			magic = rand.Uint64()

			// Stop looking for magic numbers if we found one for this square.
			lookingForNumber = false
		}
		BISHOPATTACKS[i] = bishopAttacks
		MAGICNUMBERBISHOP[i] = magic
	}

}

// Returns a set of legal moves for a rook on Square s with blockers on
// bb.
func RookMovesOnBoard(s Square, bb uint64) uint64 {
	var movesbb uint64
	// Check north.
	for {
		if (s.Row() + 1) > 8 {
			break // off the board
		}
		moveNorth := GetSquare(s.Row()+1, s.Col())
		movesbb |= (uint64(1) << uint64(moveNorth))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveNorth))&bb > 0 {
			break
		}
	}

	// Check south.
	for {
		if (s.Row() - 1) < 1 {
			break // off the board
		}
		moveSouth := GetSquare(s.Row()-1, s.Col())
		movesbb |= (uint64(1) << uint64(moveSouth))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveSouth))&bb > 0 {
			break
		}
	}

	// Check west.
	for {
		if (s.Col() - 1) < 1 {
			break // off the board
		}
		moveWest := GetSquare(s.Row(), s.Col()-1)
		movesbb |= (uint64(1) << uint64(moveWest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveWest))&bb > 0 {
			break
		}
	}
	// Check east.
	for {
		if (s.Col() + 1) > 8 {
			break // off the board
		}
		moveEast := GetSquare(s.Row(), s.Col()-1)
		movesbb |= (uint64(1) << uint64(moveEast))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveEast))&bb > 0 {
			break
		}
	}
	return movesbb
}

func BishopMovesOnBoard(s Square, bb uint64) uint64 {
	var movesbb uint64
	// Check northwest.
	for {
		if (s.Row() + 1) > 8 {
			break // off the board
		}
		if (s.Col() - 1) < 1 {
			break // off the board
		}
		moveNorthwest := GetSquare(s.Row()+1, s.Col()-1)
		movesbb |= (uint64(1) << uint64(moveNorthwest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveNorthwest))&bb > 0 {
			break
		}
	}

	// Check northeast.
	for {
		if (s.Row() + 1) > 8 {
			break // off the board
		}
		if (s.Col() + 1) > 8 {
			break // off the board
		}
		moveNortheast := GetSquare(s.Row()+1, s.Col()+1)
		movesbb |= (uint64(1) << uint64(moveNortheast))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveNortheast))&bb > 0 {
			break
		}
	}

	// Check southwest.
	for {
		if (s.Row() - 1) < 1 {
			break // off the board
		}
		if (s.Col() - 1) < 1 {
			break // off the board
		}
		moveSouthwest := GetSquare(s.Row()-1, s.Col()-1)
		movesbb |= (uint64(1) << uint64(moveSouthwest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveSouthwest))&bb > 0 {
			break
		}
	}

	// Check southeast.
	for {
		if (s.Row() - 1) < 1 {
			break // off the board
		}
		if (s.Col() + 1) > 8 {
			break // off the board
		}
		moveSouthwest := GetSquare(s.Row()-1, s.Col()+1)
		movesbb |= (uint64(1) << uint64(moveSouthwest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveSouthwest))&bb > 0 {
			break
		}
	}
	return movesbb
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
					// Since positive directions can wrap around, we might need to end the loop here.
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
