// Piepce is an interface that defines the operations possible for a piece on the board.
package game

import "fmt"

// Define the possible Colors of a piece as an enum
type Color int

const (
	WHITE     Color = 1
	NULLCOLOR Color = 0
	BLACK     Color = -1
)

// Piece represents the color and value combinations of pieces.
type Piece int

const (
	NULLPIECE = Piece(iota)
	WHITEPAWN
	WHITEKNIGHT
	WHITEBISHOP
	WHITEROOK
	WHITEQUEEN
	WHITEKING
	BLACKPAWN
	BLACKKNIGHT
	BLACKBISHOP
	BLACKROOK
	BLACKQUEEN
	BLACKKING
)

// PieceType represents the color-agnostic classification for a piece.
type PieceType int

const (
	NULLPIECETYPE = PieceType(iota)
	PAWN
	BISHOP
	KNIGHT
	ROOK
	QUEEN
	KING
)

func (c Color) String() string {
	if c == WHITE {
		return "WHITE"
	} else {
		return "BLACK"
	}
}

// TargetLegal returns true if a candidate piece can move to the
// desired square. It also optionally returns a piece that will
// be captured.
func TargetLegal(b *Board, p Piece, s Square, capture bool) (bool, Piece) {
	if s == OFFBOARD_SQUARE {
		return false, NULLPIECE
	}
	occupant := b.Squares[s]
	if occupant == NULLPIECE {
		return true, NULLPIECE
	} else {
		if capture && p.Color() != occupant.Color() {
			return true, occupant
		}
	}
	return false, NULLPIECE
}
func LegalKnightMoves(b *Board, p Piece, cur Square) []EfficientMove {
	var moves []EfficientMove
	pos := b.Position
	km := LEGALKNIGHTMOVES[cur]
	// Iterate through legal non captures
	sq := SquaresFromBitBoard(km &^ pos.Occupied)
	for i := 0; i < len(sq); i++ {
		moves = append(moves, NewEfficientMove(p, sq[i], cur))
	}
	// Iterate through legal captures
	var opp uint64
	switch p.Color() {
	case WHITE:
		opp = pos.BlackPieces
	case BLACK:
		opp = pos.WhitePieces
	}
	sq = SquaresFromBitBoard(km & opp)
	for i := 0; i < len(sq); i++ {
		s := sq[i]
		move := NewEfficientMove(p, s, cur)
		move = move.AddCapture(b.Squares[s])
		if b.Squares[s] == NULLPIECE {
			fmt.Println("Last move: " + b.LastMove.String())
			fmt.Println(SquaresFromBitBoard(km))
			fmt.Println(SquaresFromBitBoard(opp))
			b.Print()
			panic("some knight capture is nil. abort! " + s.String())

		}
		moves = append(moves, move)
	}
	return moves
}

func RayMoves(b *Board, p Piece, cur Square, bishop, rook bool) []EfficientMove {
	var moves []EfficientMove
	pos := b.Position
	var allAtk uint64
	allAtk = 0
	if bishop {
		bm := BLOCKERMASKBISHOP[cur]
		magic := MAGICNUMBERBISHOP[cur]
		key := ((bm & pos.Occupied) * magic) >> SHIFTSIZEBISHOP[cur]
		allAtk |= BISHOPATTACKS[cur][key]
	}
	if rook {
		bm := BLOCKERMASKROOK[cur]
		magic := MAGICNUMBERROOK[cur]
		key := ((bm & pos.Occupied) * magic) >> SHIFTSIZEROOK[cur]
		allAtk |= ROOKATTACKS[cur][key]
	}

	// TODO(slisenberger)
	// THIS IS ALL COPIED BOILERPLATE.. FACTOR THIS OUT.
	// Iterate through legal non captures
	for _, s := range SquaresFromBitBoard(allAtk &^ pos.Occupied) {
		moves = append(moves, NewEfficientMove(p, s, cur))
	}
	// Iterate through legal captures
	var opp uint64
	switch p.Color() {
	case WHITE:
		opp = pos.BlackPieces
	case BLACK:
		opp = pos.WhitePieces
	}
	for _, s := range SquaresFromBitBoard(allAtk & opp) {
		move := NewEfficientMove(p, s, cur)
		move = move.AddCapture(b.Squares[s])
		if b.Squares[s] == NULLPIECE {
			fmt.Println("Last move")
			fmt.Println("----------")
			fmt.Println(SquaresFromBitBoard(allAtk))
			fmt.Println(SquaresFromBitBoard(opp))
			fmt.Println(SquaresFromBitBoard(pos.Occupied))
			fmt.Println(SquaresFromBitBoard(pos.WhitePawns))
			fmt.Println(SquaresFromBitBoard(pos.WhiteBishops))
			fmt.Println(SquaresFromBitBoard(pos.WhiteKnights))
			fmt.Println(SquaresFromBitBoard(pos.BlackPawns))
			fmt.Println(SquaresFromBitBoard(pos.BlackBishops))
			fmt.Println(SquaresFromBitBoard(pos.BlackKnights))
			b.Print()
			panic("some ray capture is nil. abort! " + s.String())

		}
		moves = append(moves, move)
	}
	return moves
}

func RayAttackBitboard(b *Board, cur Square, bishop, rook bool) uint64 {
	var res uint64
	res = 0
	pos := b.Position
	if bishop {
		mask := BLOCKERMASKBISHOP[cur]
		magic := MAGICNUMBERBISHOP[cur]
		key := ((mask & pos.Occupied) * magic) >> SHIFTSIZEBISHOP[cur]
		res |= BISHOPATTACKS[cur][key]
	}
	if rook {
		mask := BLOCKERMASKROOK[cur]
		magic := MAGICNUMBERROOK[cur]
		key := ((mask & pos.Occupied) * magic) >> SHIFTSIZEROOK[cur]
		res |= ROOKATTACKS[cur][key]
	}
	return res
}

func LegalKingMoves(b *Board, p Piece, cur Square) []EfficientMove {
	moves := []EfficientMove{}
	pos := b.Position
	km := LEGALKINGMOVES[cur]
	// Iterate through legal non captures
	for _, s := range SquaresFromBitBoard(km &^ pos.Occupied) {
		moves = append(moves, NewEfficientMove(p, s, cur))
	}
	// Iterate through legal captures
	var opp uint64
	switch p.Color() {
	case WHITE:
		opp = pos.BlackPieces
	case BLACK:
		opp = pos.WhitePieces
	}
	for _, s := range SquaresFromBitBoard(km & opp) {
		move := NewEfficientMove(p, s, cur)
		move = move.AddCapture(b.Squares[s])
		if b.Squares[s] == NULLPIECE {
			fmt.Println("Last move: " + b.LastMove.String())
			fmt.Println(SquaresFromBitBoard(km))
			fmt.Println(SquaresFromBitBoard(opp))
			b.Print()
			panic("some king capture is nil. abort! " + s.String())

		}
		moves = append(moves, move)
	}
	return moves

}

// Returns the set of squares a pawn is attacking.
func PawnAttackingSquares(p Piece, cur Square) []Square {
	var res uint64
	res = 0
	switch p.Color() {
	case WHITE:
		res = WHITEPAWNATTACKS[cur]
	case BLACK:
		res = BLACKPAWNATTACKS[cur]
	}
	return SquaresFromBitBoard(res)
}

func PawnMoves(b *Board, p Piece, cur Square) []EfficientMove {
	// Check if the piece can move two squares.
	var isStartPawn bool
	var direction int // Which way pawns move.
	switch p.Color() {
	case BLACK:
		isStartPawn = cur.Row() == 7
		direction = -1
		break
	case WHITE:
		isStartPawn = cur.Row() == 2
		direction = 1
		break
	}
	moves := []EfficientMove{}
	s := GetSquare(cur.Row()+direction, cur.Col())
	if l, _ := TargetLegal(b, p, s, false); l {
		// Check for promotion and add the promotions.
		if s.Row() == 1 || s.Row() == 8 {
			move := NewEfficientMove(p, s, cur)
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEBISHOP)
			case BLACK:
				move = move.AddPromotion(BLACKBISHOP)
			}
			moves = append(moves, move)
			move = NewEfficientMove(p, s, cur)
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEKNIGHT)
			case BLACK:
				move = move.AddPromotion(BLACKKNIGHT)
			}
			moves = append(moves, move)
			move = NewEfficientMove(p, s, cur)
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEROOK)
			case BLACK:
				move = move.AddPromotion(BLACKROOK)
			}
			moves = append(moves, move)
			move = NewEfficientMove(p, s, cur)
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEQUEEN)
			case BLACK:
				move = move.AddPromotion(BLACKQUEEN)
			}
			moves = append(moves, move)

		} else {
			moves = append(moves, NewEfficientMove(p, s, cur))
		}
		// We only can move forward two if we can also move forward one.
		if isStartPawn {
			s := GetSquare(cur.Row()+2*direction, cur.Col())
			if l, _ := TargetLegal(b, p, s, false); l {
				m := NewEfficientMove(p, s, cur)
				m = m.AddTwoPawnAdvance()
				moves = append(moves, m)
			}
		}
	}
	// Check for side captures.
	sq := PawnAttackingSquares(p, cur)
	for i := 0; i < len(sq); i++ {
		s := sq[i]
		occupant := b.Squares[s]
		if occupant != NULLPIECE && occupant.Color() != p.Color() {
			if s.Row() == 1 || s.Row() == 8 {
				move := NewEfficientMove(p, s, cur)
				move = move.AddCapture(occupant)
				switch p.Color() {
				case WHITE:
					move = move.AddPromotion(WHITEQUEEN)
				case BLACK:
					move = move.AddPromotion(BLACKQUEEN)
				}
				moves = append(moves, move)
				move = NewEfficientMove(p, s, cur)
				move = move.AddCapture(occupant)
				switch p.Color() {
				case WHITE:
					move = move.AddPromotion(WHITEKNIGHT)
				case BLACK:
					move = move.AddPromotion(BLACKKNIGHT)
				}
				moves = append(moves, move)
				move = NewEfficientMove(p, s, cur)
				move = move.AddCapture(occupant)
				switch p.Color() {
				case WHITE:
					move = move.AddPromotion(WHITEBISHOP)
				case BLACK:
					move = move.AddPromotion(BLACKBISHOP)
				}
				moves = append(moves, move)
				move = NewEfficientMove(p, s, cur)
				move = move.AddCapture(occupant)
				switch p.Color() {
				case WHITE:
					move = move.AddPromotion(WHITEROOK)
				case BLACK:
					move = move.AddPromotion(BLACKROOK)
				}
				moves = append(moves, move)

			} else {
				move := NewEfficientMove(p, s, cur)
				move = move.AddCapture(occupant)
				moves = append(moves, move)
			}
		}
	}
	// Check for en passants
	epSquare := b.EPSquare
	// If en passant is legal, we migh be able to capture.
	if epSquare != OFFBOARD_SQUARE {
		adjToEP := cur.Col()-1 == epSquare.Col() || cur.Col()+1 == epSquare.Col()
		if p.Color() == WHITE && cur.Row() == 5 && adjToEP && epSquare.Row() == 6 {
			move := NewEfficientMove(p, GetSquare(6, epSquare.Col()), cur)
			move = move.AddEnPassant()
			capturedPiece := b.Squares[epSquare]
			if capturedPiece == NULLPIECE {
				panic(fmt.Sprintf("capture on %v is nil", epSquare))
			}
			move = move.AddCapture(capturedPiece)
			moves = append(moves, move)
		}
		if p.Color() == BLACK && cur.Row() == 4 && adjToEP && epSquare.Row() == 3 {
			move := NewEfficientMove(p, GetSquare(3, epSquare.Col()), cur)
			move = move.AddEnPassant()
			capturedPiece := b.Squares[epSquare]
			if capturedPiece == NULLPIECE {
				panic(fmt.Sprintf("capture on %v is nil", epSquare))
			}
			move = move.AddCapture(capturedPiece)
			moves = append(moves, move)
		}
	}
	return moves
}



var colors = [13]Color{
	NULLCOLOR, WHITE, WHITE, WHITE, WHITE, WHITE, WHITE, 
	BLACK, BLACK, BLACK, BLACK, BLACK, BLACK,
}

func (p Piece) Color() Color {
	return colors[p]
}

func (p Piece) Type() PieceType {
	switch p {
	case WHITEPAWN, BLACKPAWN:
		return PAWN
	case WHITEKNIGHT, BLACKKNIGHT:
		return KNIGHT
	case WHITEBISHOP, BLACKBISHOP:
		return BISHOP
	case WHITEROOK, BLACKROOK:
		return ROOK
	case WHITEQUEEN, BLACKQUEEN:
		return QUEEN
	case WHITEKING, BLACKKING:
		return KING
	}
	return NULLPIECETYPE
}

func (p Piece) Graphic() string {
	switch p {
	case BLACKPAWN:
		return "♙"
	case WHITEPAWN:
		return "♟"
	case BLACKBISHOP:
		return "♗"
	case WHITEBISHOP:
		return "♝"
	case BLACKKNIGHT:
		return "♘"
	case WHITEKNIGHT:
		return "♞"
	case BLACKROOK:
		return "♖"
	case WHITEROOK:
		return "♜"
	case BLACKQUEEN:
		return "♕"
	case WHITEQUEEN:
		return "♛"
	case BLACKKING:
		return "♔"
	case WHITEKING:
		return "♚"
	}
	return ""
}

func (p Piece) String() string {
	switch p {
	case BLACKPAWN, WHITEPAWN:
		return "P"
	case BLACKBISHOP, WHITEBISHOP:
		return "B"
	case BLACKKNIGHT, WHITEKNIGHT:
		return "N"
	case BLACKROOK, WHITEROOK:
		return "R"
	case BLACKQUEEN, WHITEQUEEN:
		return "Q"
	case BLACKKING, WHITEKING:
		return "K"
	}
	return ""
}

func (p Piece) Value() float64 {
	switch p {
	case BLACKPAWN, WHITEPAWN:
		return 1.0
	case BLACKBISHOP, WHITEBISHOP:
		return 3.0
	case BLACKKNIGHT, WHITEKNIGHT:
		return 3.0
	case BLACKROOK, WHITEROOK:
		return 5.0
	case BLACKQUEEN, WHITEQUEEN:
		return 9.0
	case BLACKKING, WHITEKING:
		return 100.0
	}
	return 0.0
}

// Legal moves returns the pseudolegal moves for a given piece on a
// board on the starting square cur.
func LegalMoves(b *Board, p Piece, cur Square) []EfficientMove {
	m := []EfficientMove{}
	switch p {
	case WHITEKNIGHT, BLACKKNIGHT:
		return KnightMoves(b, p, cur)
	case WHITEBISHOP, BLACKBISHOP:
		return BishopMoves(b, p, cur)
	case WHITEROOK, BLACKROOK:
		return RookMoves(b, p, cur)
	case WHITEPAWN, BLACKPAWN:
		return PawnMoves(b, p, cur)
	case WHITEKING, BLACKKING:
		return KingMoves(b, p, cur)
	case WHITEQUEEN, BLACKQUEEN:
		return QueenMoves(b, p, cur)
	}
	return m
}

func LegalCaptures(b *Board, p Piece, cur Square) []EfficientMove {
	moves := []EfficientMove{}
	atk := AttackBitboard(b, p, cur)
	var opp uint64
	switch p.Color() {
	case WHITE:
		opp = b.Position.BlackPieces
	case BLACK:
		opp = b.Position.WhitePieces
	}

	sq := SquaresFromBitBoard(atk & opp)
	for i := 0; i < len(sq); i++ {
		s := sq[i]
		isPromotion := p.Type() == PAWN && (s.Row() == 1 || s.Row() == 8)
		if !isPromotion {
			move := NewEfficientMove(p, s, cur)
			move = move.AddCapture(b.Squares[s])
			moves = append(moves, move)
		} else {
			move := NewEfficientMove(p, s, cur)
			move.AddCapture(b.Squares[s])
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEQUEEN)
			case BLACK:
				move = move.AddPromotion(BLACKQUEEN)
			}
			moves = append(moves, move)
			move = NewEfficientMove(p, s, cur)
			move = move.AddCapture(b.Squares[s])
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEBISHOP)
			case BLACK:
				move = move.AddPromotion(BLACKBISHOP)
			}
			moves = append(moves, move)
			move = NewEfficientMove(p, s, cur)
			move = move.AddCapture(b.Squares[s])
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEKNIGHT)
			case BLACK:
				move = move.AddPromotion(BLACKKNIGHT)
			}
			moves = append(moves, move)
			move = NewEfficientMove(p, s, cur)
			move = move.AddCapture(b.Squares[s])
			switch p.Color() {
			case WHITE:
				move = move.AddPromotion(WHITEROOK)
			case BLACK:
				move = move.AddPromotion(BLACKROOK)
			}
			moves = append(moves, move)
		}
	}
	// Do En Passant captures for pawns.
	if p.Type() == PAWN {
		// If attacking the en passant square, we can capture there.
		sq = SquaresFromBitBoard(atk)
		for i := 0; i < len(sq); i++ {
			s := sq[i]
			if s == b.EPSquare {
				var move EfficientMove
				switch p.Color() {
				case WHITE:
					move = NewEfficientMove(p, GetSquare(6, b.EPSquare.Col()), cur)
				case BLACK:
					move = NewEfficientMove(p, GetSquare(3, b.EPSquare.Col()), cur)
				}
				capturedPiece := b.Squares[b.EPSquare]
				move = move.AddCapture(capturedPiece)
				moves = append(moves, move)
			}
		}
	}
	return moves
}

// Returns a list of squares under attack by piece p on square cur.
func AttackBitboard(b *Board, p Piece, cur Square) uint64 {
	switch p {
	case WHITEPAWN:
		return WHITEPAWNATTACKS[cur]
	case BLACKPAWN:
		return BLACKPAWNATTACKS[cur]
	case WHITEKNIGHT, BLACKKNIGHT:
		return LEGALKNIGHTMOVES[cur]
	case WHITEBISHOP, BLACKBISHOP:
		return BishopAttackBitboard(b, cur)
	case WHITEROOK, BLACKROOK:
		return RookAttackBitboard(b, cur)
	case WHITEQUEEN, BLACKQUEEN:
		return QueenAttackBitboard(b, cur)
	case WHITEKING, BLACKKING:
		return LEGALKINGMOVES[cur]
	}

	return 0
}
