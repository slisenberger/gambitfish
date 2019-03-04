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
// in looking up precomputed moves. I tried computing these but it was
// taking a really long time. Graciously stolen from Crafty's open source
// engine.
var MAGICNUMBERROOK = [64]uint64{
	0x0080001020400080, 0x0040001000200040, 0x0080081000200080, 0x0080040800100080,
	0x0080020400080080, 0x0080010200040080, 0x0080008001000200, 0x0080002040800100,
	0x0000800020400080, 0x0000400020005000, 0x0000801000200080, 0x0000800800100080,
	0x0000800400080080, 0x0000800200040080, 0x0000800100020080, 0x0000800040800100,
	0x0000208000400080, 0x0000404000201000, 0x0000808010002000, 0x0000808008001000,
	0x0000808004000800, 0x0000808002000400, 0x0000010100020004, 0x0000020000408104,
	0x0000208080004000, 0x0000200040005000, 0x0000100080200080, 0x0000080080100080,
	0x0000040080080080, 0x0000020080040080, 0x0000010080800200, 0x0000800080004100,
	0x0000204000800080, 0x0000200040401000, 0x0000100080802000, 0x0000080080801000,
	0x0000040080800800, 0x0000020080800400, 0x0000020001010004, 0x0000800040800100,
	0x0000204000808000, 0x0000200040008080, 0x0000100020008080, 0x0000080010008080,
	0x0000040008008080, 0x0000020004008080, 0x0000010002008080, 0x0000004081020004,
	0x0000204000800080, 0x0000200040008080, 0x0000100020008080, 0x0000080010008080,
	0x0000040008008080, 0x0000020004008080, 0x0000800100020080, 0x0000800041000080,
	0x00FFFCDDFCED714A, 0x007FFCDDFCED714A, 0x003FFFCDFFD88096, 0x0000040810002101,
	0x0001000204080011, 0x0001000204000801, 0x0001000082000401, 0x0001FFFAABFAD1A2,
}

var MAGICNUMBERBISHOP = [64]uint64{
	0x0002020202020200, 0x0002020202020000, 0x0004010202000000, 0x0004040080000000,
	0x0001104000000000, 0x0000821040000000, 0x0000410410400000, 0x0000104104104000,
	0x0000040404040400, 0x0000020202020200, 0x0000040102020000, 0x0000040400800000,
	0x0000011040000000, 0x0000008210400000, 0x0000004104104000, 0x0000002082082000,
	0x0004000808080800, 0x0002000404040400, 0x0001000202020200, 0x0000800802004000,
	0x0000800400A00000, 0x0000200100884000, 0x0000400082082000, 0x0000200041041000,
	0x0002080010101000, 0x0001040008080800, 0x0000208004010400, 0x0000404004010200,
	0x0000840000802000, 0x0000404002011000, 0x0000808001041000, 0x0000404000820800,
	0x0001041000202000, 0x0000820800101000, 0x0000104400080800, 0x0000020080080080,
	0x0000404040040100, 0x0000808100020100, 0x0001010100020800, 0x0000808080010400,
	0x0000820820004000, 0x0000410410002000, 0x0000082088001000, 0x0000002011000800,
	0x0000080100400400, 0x0001010101000200, 0x0002020202000400, 0x0001010101000200,
	0x0000410410400000, 0x0000208208200000, 0x0000002084100000, 0x0000000020880000,
	0x0000001002020000, 0x0000040408020000, 0x0004040404040000, 0x0002020202020000,
	0x0000104104104000, 0x0000002082082000, 0x0000000020841000, 0x0000000000208800,
	0x0000000010020200, 0x0000000404080200, 0x0000040404040400, 0x0002020202020200,
}

// The amount to shift by when calculating keys for magic numbers, indexed by square.
var SHIFTSIZEROOK = [64]uint64{
	52, 53, 53, 53, 53, 53, 53, 52,
	53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 53, 53, 53, 53, 53,
}

var SHIFTSIZEBISHOP = [64]uint64{
	58, 59, 59, 59, 59, 59, 59, 58,
	59, 59, 59, 59, 59, 59, 59, 59,
	59, 59, 57, 57, 57, 57, 59, 59,
	59, 59, 57, 55, 55, 57, 59, 59,
	59, 59, 57, 55, 55, 57, 59, 59,
	59, 59, 57, 57, 57, 57, 59, 59,
	59, 59, 59, 59, 59, 59, 59, 59,
	58, 59, 59, 59, 59, 59, 59, 58,
}

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
	BuildByteLookupTable()
	InitRayAttacks()
	InitBlockerMasks()
	InitMagicBitboards()
	LEGALKINGMOVES = LegalKingMovesDict()
	LEGALKNIGHTMOVES = LegalKnightMovesDict()
	InitPawnAttacks()
	InitZobristNumbers()
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
		// If it's our row, this is our square, it's irrelevant
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
			bb = SetBitOnBoard(bb, GetSquare(s.Row()-deltaCol, col))
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
	ROOKATTACKS = [64][]uint64{}
	BISHOPATTACKS = [64][]uint64{}
	// Get all the possible blockerboards, and their resulting
	// legal moves.

	var i uint64
	// Loop once for rooks.
	for i = 0; i < 64; i++ {
		mask := BLOCKERMASKROOK[i]
		bitCount := bits.OnesCount64(mask)
		perms := GenerateBlockerPermutations(Square(i), mask, mask)
		rookAttacks := make([]uint64, 1<<uint64(bitCount))
		for _, perm := range perms {
			moves := RookMovesOnBoard(Square(i), perm)
			key := (perm * MAGICNUMBERROOK[i]) >> uint64(64-bitCount)
			rookAttacks[key] = moves
		}
		ROOKATTACKS[i] = rookAttacks
	}

	// And again for bishops.
	for i = 0; i < 64; i++ {
		mask := BLOCKERMASKBISHOP[i]
		bitCount := bits.OnesCount64(mask)
		perms := GenerateBlockerPermutations(Square(i), mask, mask)
		bishopAttacks := make([]uint64, 1<<uint64(bitCount))
		for _, perm := range perms {
			moves := BishopMovesOnBoard(Square(i), perm)
			key := (perm * MAGICNUMBERBISHOP[i]) >> uint64(64-bitCount)
			bishopAttacks[key] = moves
		}
		BISHOPATTACKS[i] = bishopAttacks
	}

}

// Returns a set of legal moves for a rook on Square s with blockers on
// bb.
func RookMovesOnBoard(s Square, bb uint64) uint64 {
	var movesbb uint64
	// Check north.
	moveNorth := s
	for {
		if (moveNorth.Row() + 1) > 8 {
			break // off the board
		}
		moveNorth = GetSquare(moveNorth.Row()+1, moveNorth.Col())
		movesbb = movesbb | uint64((uint64(1) << uint64(moveNorth)))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveNorth))&bb > 0 {
			break
		}
	}

	// Check south.
	moveSouth := s
	for {
		if (moveSouth.Row() - 1) < 1 {
			break // off the board
		}
		moveSouth = GetSquare(moveSouth.Row()-1, moveSouth.Col())
		movesbb |= (uint64(1) << uint64(moveSouth))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveSouth))&bb > 0 {
			break
		}
	}

	// Check west.
	moveWest := s
	for {
		if (moveWest.Col() - 1) < 1 {
			break // off the board
		}
		moveWest = GetSquare(moveWest.Row(), moveWest.Col()-1)
		movesbb |= (uint64(1) << uint64(moveWest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveWest))&bb > 0 {
			break
		}
	}
	// Check east.
	moveEast := s
	for {
		if (moveEast.Col() + 1) > 8 {
			break // off the board
		}
		moveEast = GetSquare(moveEast.Row(), moveEast.Col()+1)
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
	moveNorthwest := s
	for {
		if (moveNorthwest.Row() + 1) > 8 {
			break // off the board
		}
		if (moveNorthwest.Col() - 1) < 1 {
			break // off the board
		}
		moveNorthwest = GetSquare(moveNorthwest.Row()+1, moveNorthwest.Col()-1)
		movesbb |= (uint64(1) << uint64(moveNorthwest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveNorthwest))&bb > 0 {
			break
		}
	}

	// Check northeast.
	moveNortheast := s
	for {
		if (moveNortheast.Row() + 1) > 8 {
			break // off the board
		}
		if (moveNortheast.Col() + 1) > 8 {
			break // off the board
		}
		moveNortheast = GetSquare(moveNortheast.Row()+1, moveNortheast.Col()+1)
		movesbb |= (uint64(1) << uint64(moveNortheast))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveNortheast))&bb > 0 {
			break
		}
	}

	// Check southwest.
	moveSouthwest := s
	for {
		if (moveSouthwest.Row() - 1) < 1 {
			break // off the board
		}
		if (moveSouthwest.Col() - 1) < 1 {
			break // off the board
		}
		moveSouthwest = GetSquare(moveSouthwest.Row()-1, moveSouthwest.Col()-1)
		movesbb |= (uint64(1) << uint64(moveSouthwest))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveSouthwest))&bb > 0 {
			break
		}
	}

	// Check southeast.
	moveSoutheast := s
	for {
		if (moveSoutheast.Row() - 1) < 1 {
			break // off the board
		}
		if (moveSoutheast.Col() + 1) > 8 {
			break // off the board
		}
		moveSoutheast = GetSquare(moveSoutheast.Row()-1, moveSoutheast.Col()+1)
		movesbb |= (uint64(1) << uint64(moveSoutheast))
		// If the new move intersects the blocker board, break.
		if (uint64(1)<<uint64(moveSoutheast))&bb > 0 {
			break
		}
	}
	return movesbb
}

func GenerateBlockerPermutations(s Square, blockerMask uint64, board uint64) []uint64 {
	allPerms := []uint64{}
	// If we've twiddled all the bits, then this board is a permutation.
	if blockerMask == 0 {
		return []uint64{board}
	}
	// Determine the next bit we are toggling in permutations, and remove it from the blocker mask.
	nextToggledBit := bits.TrailingZeros64(blockerMask)
	blockerMask = UnSetBitOnBoard(blockerMask, Square(nextToggledBit))
	blockerLeftOn := board
	blockerRemoved := board ^ (uint64(1) << uint64(nextToggledBit))

	// Now enumerate all permutations with the toggled bit left on, and the toggled bit left off.
	allPerms = append(allPerms, GenerateBlockerPermutations(s, blockerMask, blockerLeftOn)...)
	allPerms = append(allPerms, GenerateBlockerPermutations(s, blockerMask, blockerRemoved)...)
	return allPerms
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
