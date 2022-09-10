// Internal manages the internal datastructures useful by all parts of
// the engine, such as an opening book, zobrist hash codes, and legal
// move positions.
package game

import "math/rand"
import "math/bits"
import "time"
import "fmt"

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
var BLOCKERMASKBISHOP = [64]uint64{
	0x0040201008040200, 0x0000402010080400, 0x0000004020100A00, 0x0000000040221400,
	0x0000000002442800, 0x0000000204085000, 0x0000020408102000, 0x0002040810204000,
	0x0020100804020000, 0x0040201008040000, 0x00004020100A0000, 0x0000004022140000,
	0x0000000244280000, 0x0000020408500000, 0x0002040810200000, 0x0004081020400000,
	0x0010080402000200, 0x0020100804000400, 0x004020100A000A00, 0x0000402214001400,
	0x0000024428002800, 0x0002040850005000, 0x0004081020002000, 0x0008102040004000,
	0x0008040200020400, 0x0010080400040800, 0x0020100A000A1000, 0x0040221400142200,
	0x0002442800284400, 0x0004085000500800, 0x0008102000201000, 0x0010204000402000,
	0x0004020002040800, 0x0008040004081000, 0x00100A000A102000, 0x0022140014224000,
	0x0044280028440200, 0x0008500050080400, 0x0010200020100800, 0x0020400040201000,
	0x0002000204081000, 0x0004000408102000, 0x000A000A10204000, 0x0014001422400000,
	0x0028002844020000, 0x0050005008040200, 0x0020002010080400, 0x0040004020100800,
	0x0000020408102000, 0x0000040810204000, 0x00000A1020400000, 0x0000142240000000,
	0x0000284402000000, 0x0000500804020000, 0x0000201008040200, 0x0000402010080400,
	0x0002040810204000, 0x0004081020400000, 0x000A102040000000, 0x0014224000000000,
	0x0028440200000000, 0x0050080402000000, 0x0020100804020000, 0x0040201008040200,
}
var BLOCKERMASKROOK = [64]uint64{
	0x000101010101017E, 0x000202020202027C, 0x000404040404047A, 0x0008080808080876,
	0x001010101010106E, 0x002020202020205E, 0x004040404040403E, 0x008080808080807E,
	0x0001010101017E00, 0x0002020202027C00, 0x0004040404047A00, 0x0008080808087600,
	0x0010101010106E00, 0x0020202020205E00, 0x0040404040403E00, 0x0080808080807E00,
	0x00010101017E0100, 0x00020202027C0200, 0x00040404047A0400, 0x0008080808760800,
	0x00101010106E1000, 0x00202020205E2000, 0x00404040403E4000, 0x00808080807E8000,
	0x000101017E010100, 0x000202027C020200, 0x000404047A040400, 0x0008080876080800,
	0x001010106E101000, 0x002020205E202000, 0x004040403E404000, 0x008080807E808000,
	0x0001017E01010100, 0x0002027C02020200, 0x0004047A04040400, 0x0008087608080800,
	0x0010106E10101000, 0x0020205E20202000, 0x0040403E40404000, 0x0080807E80808000,
	0x00017E0101010100, 0x00027C0202020200, 0x00047A0404040400, 0x0008760808080800,
	0x00106E1010101000, 0x00205E2020202000, 0x00403E4040404000, 0x00807E8080808000,
	0x007E010101010100, 0x007C020202020200, 0x007A040404040400, 0x0076080808080800,
	0x006E101010101000, 0x005E202020202000, 0x003E404040404000, 0x007E808080808000,
	0x7E01010101010100, 0x7C02020202020200, 0x7A04040404040400, 0x7608080808080800,
	0x6E10101010101000, 0x5E20202020202000, 0x3E40404040404000, 0x7E80808080808000,
}

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
	52, 53, 53, 53, 53, 53, 53, 52,
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
	InitMagicBitboards()
	LEGALKINGMOVES = LegalKingMovesDict()
	LEGALKNIGHTMOVES = LegalKnightMovesDict()
	InitPawnAttacks()
	InitZobristNumbers()
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
		perms := GenerateBlockerPermutations(mask, mask)
		fmt.Printf("found %d permutations\n", len(perms))
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
		perms := GenerateBlockerPermutations(mask, mask)
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
		moveNorth = GetSquare(moveNorth.Row()+1, moveNorth.Col())
		if moveNorth == OFFBOARD_SQUARE {
			break
		}
		movesbb = SetBitOnBoard(movesbb, moveNorth)
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveNorth)) & bb) != 0 {
			break
		}
	}

	// Check south.
	moveSouth := s
	for {
		moveSouth = GetSquare(moveSouth.Row()-1, moveSouth.Col())
		if moveSouth == OFFBOARD_SQUARE {
			break
		}
		movesbb = SetBitOnBoard(movesbb, moveSouth)
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveSouth)) & bb) != 0 {
			break
		}
	}

	// Check west.
	moveWest := s
	for {
		moveWest = GetSquare(moveWest.Row(), moveWest.Col()-1)
		if moveWest == OFFBOARD_SQUARE {
			break
		}
		movesbb = SetBitOnBoard(movesbb, moveWest)
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveWest)) & bb) != 0 {
			break
		}
	}
	// Check east.
	moveEast := s
	for {
		moveEast = GetSquare(moveEast.Row(), moveEast.Col()+1)
		if moveEast == OFFBOARD_SQUARE {
			break
		}
		movesbb = SetBitOnBoard(movesbb, moveEast)
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveEast)) & bb) != 0 {
			break
		}
	}
	if uint64(s) == 63 {
		fmt.Println(SquaresFromBitBoard(movesbb))
	}
	return movesbb
}

func BishopMovesOnBoard(s Square, bb uint64) uint64 {
	var movesbb uint64
	// Check northwest.
	moveNorthwest := s
	for {
		moveNorthwest = GetSquare(moveNorthwest.Row()+1, moveNorthwest.Col()-1)
		if moveNorthwest == OFFBOARD_SQUARE {
			break
		}
		movesbb |= (uint64(1) << uint64(moveNorthwest))
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveNorthwest)) & bb) != 0 {
			break
		}
	}

	// Check northeast.
	moveNortheast := s
	for {
		moveNortheast = GetSquare(moveNortheast.Row()+1, moveNortheast.Col()+1)
		if moveNortheast == OFFBOARD_SQUARE {
			break
		}
		movesbb |= (uint64(1) << uint64(moveNortheast))
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveNortheast)) & bb) != 0 {
			break
		}
	}

	// Check southwest.
	moveSouthwest := s
	for {
		moveSouthwest = GetSquare(moveSouthwest.Row()-1, moveSouthwest.Col()-1)
		if moveSouthwest == OFFBOARD_SQUARE {
			break
		}
		movesbb |= (uint64(1) << uint64(moveSouthwest))
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveSouthwest)) & bb) != 0 {
			break
		}
	}

	// Check southeast.
	moveSoutheast := s
	for {
		moveSoutheast = GetSquare(moveSoutheast.Row()-1, moveSoutheast.Col()+1)
		if moveSoutheast == OFFBOARD_SQUARE {
			break
		}
		movesbb |= (uint64(1) << uint64(moveSoutheast))
		// If the new move intersects the blocker board, break.
		if ((uint64(1) << uint64(moveSoutheast)) & bb) != 0 {
			break
		}
	}
	return movesbb
}

func GenerateBlockerPermutations(blockerMask uint64, board uint64) []uint64 {
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
	allPerms = append(allPerms, GenerateBlockerPermutations(blockerMask, blockerLeftOn)...)
	allPerms = append(allPerms, GenerateBlockerPermutations(blockerMask, blockerRemoved)...)
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
			for i := 0; i < 64; i++ {
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
