// package bitboard is an implementation of a bitboard representation
// of a chess board state and its corresponding utilities.
package game

import "fmt"

import "math/bits"

type Position struct {
	// White piece locations.
	WhiteKing    uint64
	WhiteQueens  uint64
	WhiteRooks   uint64
	WhiteKnights uint64
	WhiteBishops uint64
	WhitePawns   uint64
	// Black piece locations.
	BlackKing    uint64
	BlackQueens  uint64
	BlackRooks   uint64
	BlackKnights uint64
	BlackBishops uint64
	BlackPawns   uint64

	// Convenience maps
	WhitePieces uint64
	BlackPieces uint64
	Occupied    uint64
}

func UpdateBitboards(bb Position) Position {
	bb.WhitePieces = bb.WhiteKing | bb.WhiteQueens | bb.WhiteRooks | bb.WhiteKnights | bb.WhiteBishops | bb.WhitePawns
	bb.BlackPieces = bb.BlackKing | bb.BlackQueens | bb.BlackRooks | bb.BlackKnights | bb.BlackBishops | bb.BlackPawns
	bb.Occupied = bb.WhitePieces | bb.BlackPieces
	return bb
}

// SetPiece updates a bitboard to move a piece to a given square.
func SetPiece(bb Position, p Piece, s Square) Position {
	switch p.Type() {
	case KING:
		if p.Color() == WHITE {
			bb.WhiteKing = SetBitOnBoard(bb.WhiteKing, s)
		} else {
			bb.BlackKing = SetBitOnBoard(bb.BlackKing, s)
		}
	case QUEEN:
		if p.Color() == WHITE {
			bb.WhiteQueens = SetBitOnBoard(bb.WhiteQueens, s)
		} else {
			bb.BlackQueens = SetBitOnBoard(bb.BlackQueens, s)
		}
	case ROOK:
		if p.Color() == WHITE {
			bb.WhiteRooks = SetBitOnBoard(bb.WhiteRooks, s)
		} else {
			bb.BlackRooks = SetBitOnBoard(bb.BlackRooks, s)
		}
	case BISHOP:
		if p.Color() == WHITE {
			bb.WhiteBishops = SetBitOnBoard(bb.WhiteBishops, s)
		} else {
			bb.BlackBishops = SetBitOnBoard(bb.BlackBishops, s)
		}
	case KNIGHT:
		if p.Color() == WHITE {
			bb.WhiteKnights = SetBitOnBoard(bb.WhiteKnights, s)
		} else {
			bb.BlackKnights = SetBitOnBoard(bb.BlackKnights, s)
		}
	case PAWN:
		if p.Color() == WHITE {
			bb.WhitePawns = SetBitOnBoard(bb.WhitePawns, s)
		} else {
			bb.BlackPawns = SetBitOnBoard(bb.BlackPawns, s)
		}
	}
	return bb
}

func UnSetPiece(bb Position, p Piece, s Square) Position {
	switch p.Type() {
	case KING:
		if p.Color() == WHITE {
			bb.WhiteKing = UnSetBitOnBoard(bb.WhiteKing, s)
		} else {
			bb.BlackKing = UnSetBitOnBoard(bb.BlackKing, s)
		}
	case QUEEN:
		if p.Color() == WHITE {
			bb.WhiteQueens = UnSetBitOnBoard(bb.WhiteQueens, s)
		} else {
			bb.BlackQueens = UnSetBitOnBoard(bb.BlackQueens, s)
		}
	case ROOK:
		if p.Color() == WHITE {
			bb.WhiteRooks = UnSetBitOnBoard(bb.WhiteRooks, s)
		} else {
			oldRooks := bb.BlackRooks
			bb.BlackRooks = UnSetBitOnBoard(bb.BlackRooks, s)
			if oldRooks == bb.BlackRooks {
				fmt.Println("black rook was not unset..")
				fmt.Println(s)
				fmt.Println(SquaresFromBitBoard(bb.BlackRooks))
			}
		}
	case BISHOP:
		if p.Color() == WHITE {
			bb.WhiteBishops = UnSetBitOnBoard(bb.WhiteBishops, s)
		} else {
			bb.BlackBishops = UnSetBitOnBoard(bb.BlackBishops, s)
		}
	case KNIGHT:
		if p.Color() == WHITE {
			bb.WhiteKnights = UnSetBitOnBoard(bb.WhiteKnights, s)
		} else {
			bb.BlackKnights = UnSetBitOnBoard(bb.BlackKnights, s)
		}
	case PAWN:
		if p.Color() == WHITE {
			bb.WhitePawns = UnSetBitOnBoard(bb.WhitePawns, s)
		} else {
			bb.BlackPawns = UnSetBitOnBoard(bb.BlackPawns, s)
		}
	}
	return bb
}

// SetBitOnBoard updates a single 64 bit int with a newly set bit.
func SetBitOnBoard(board uint64, s Square) uint64 {
	return (board | (1 << uint64(s)))
}

// UnSetBitOnBoard updates a single 64 bit int, removing that bit from the board.
func UnSetBitOnBoard(board uint64, s Square) uint64 {
	return (board &^ (uint64(1) << uint64(s)))
}

// SquaresFromBitBoard returns a list of squares represented by the bits
// in a bitboard.
func SquaresFromBitBoard(board uint64) []Square { // Split the board into 8 2byte words, and get the precomputed squares from all
	// 4 of these byte patterns.
	var b1, b2, b3, b4 uint64
	b1 = board & 0xffff
	b2 = board & 0xffff0000
	b3 = board & 0xffff00000000
	b4 = board & 0xffff000000000000

	// Build a slice and add the elements to it.
	s := make([]Square, 0, bits.OnesCount64(board))
	if b1 != 0 {
		s = append(s, byteLookup[b1]...)
	}
	if b2 != 0 {
		s = append(s, byteLookup[b2]...)
	}
	if b3 != 0 {
		s = append(s, byteLookup[b3]...)
	}
	if b4 != 0 {
		s = append(s, byteLookup[b4]...)
	}
	// This is the old way. Keeping it for posterity.
	//	for board > 0 {
	//		idx := BitScanForward(board)
	//		s = append(s, idx)
	//		board = board & (board - 1)
	//	}
	return s
}

// To speed up turning a bitboard into a list of squares, we look up pre-existing conversions of
// bytes to a list of squares.
var byteLookup map[uint64][]Square

func BuildByteLookupTable() {
	// We want to build a lookup table for each 2-byte word and
	// the corresponding squares that belongs to it. We can do this by
	// enumerating each word, running a slow calculation, and then adding
	// these to the map.
	byteLookup = map[uint64][]Square{}

	// For each of 8 possible byte-sized words.
	var i uint64
	var b_conf uint64
	var b uint64
	for i = 0; i < 4; i++ {
		// Take the byte configuration pattern and build the two byte word out of it.
		for b_conf = 0; b_conf < 65536; b_conf++ {
			// bitshift it to represent the real byte affected.
			b = b_conf << (16 * i)
			s := []Square{}
			for b > 0 {
				idx := BitScanForward(b)
				s = append(s, idx)
				b = b & (b - 1)
			}
			byteLookup[b_conf<<(16*i)] = s
		}
	}
}

// BITSCAN UTILITIES

// The index for a De Bruijn sequence
var dbIndex64 = [64]int{
	0, 1, 48, 2, 57, 49, 28, 3,
	61, 58, 50, 42, 38, 29, 17, 4,
	62, 55, 59, 36, 53, 51, 43, 22,
	45, 39, 33, 30, 24, 18, 12, 5,
	63, 47, 56, 27, 60, 41, 37, 16,
	54, 35, 52, 21, 44, 32, 23, 11,
	46, 26, 40, 15, 34, 20, 31, 10,
	25, 14, 19, 9, 13, 8, 7, 6}

// BitScan returns the first square encountered in a bitscan. If
// forward is true, we start at the LSB. Else, we start at the MSB.
func BitScan(board uint64, forward bool) Square {
	if forward {
		return BitScanForward(board)
	} else {
		return BitScanReverse(board)
	}
}

func BitScanForward(board uint64) Square {
	debruijn64 := uint64(0x03F79D71B4CB0A89)
	// Count the bits, as long as they exist.
	if board == 0 {
		return Square(64)
	}
	return Square(dbIndex64[((board&-board)*debruijn64)>>58])
}

var ms1btable = [256]uint8{
	0x00, 0x01, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04,
	0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
	0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
}

func BitScanReverse(board uint64) Square {
	i := 0
	if board > uint64(0xFFFFFFFF) {
		board = board >> 32
		i += 32
	}
	if board > uint64(0xFFFF) {
		board = board >> 16
		i += 16
	}
	if board > uint64(0xFF) {
		board = board >> 8
		i += 8
	}
	return Square(i + int(ms1btable[board]-1))
}
