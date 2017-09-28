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
	return (board ^ (1 << uint64(s)))
}

// SquaresFromBitBoard returns a list of squares represented by the bits
// in a bitboard.
func SquaresFromBitBoard(board uint64) []Square {
	s := []Square{}
	i := 0
	// Count the bits, as long as they exist.
	for board > 0 {
		// And the board with 1, if so rsb is set, get the square.
		if board&uint64(1) > 0 {
			s = append(s, Square(i))
		}
		i++
		board = board >> 1
	}
	return s
}
