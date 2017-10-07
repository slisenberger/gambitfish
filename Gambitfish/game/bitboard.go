// package bitboard is an implementation of a bitboard representation
// of a chess board state and its corresponding utilities.
package game

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
			bb.BlackRooks = UnSetBitOnBoard(bb.BlackRooks, s)
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
	return (board &^ (1 << uint64(s)))
}

// SquaresFromBitBoard returns a list of squares represented by the bits
// in a bitboard.
func SquaresFromBitBoard(board uint64) []Square {
	s := []Square{}
	for board > 0 {
		idx := BitScanForward(board)
		s = append(s, idx)
		board = board & (board - 1)
	}
	return s
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
	if board > uint64(0xF) {
		board = board >> 4
		i += 4
	}
	if board > uint64(0x3) {
		board = board >> 2
		i += 2
	}
	if board > uint64(0x3) {
		board = board >> 2
		i += 2
	}
	if board > uint64(1) {
		board = board >> 1
		i += 1
	}
	return Square(i)
}
