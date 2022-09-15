package player

import "bufio"
import "errors"
import "fmt"
import "math"
import "os"
import "strings"
import "time"
import "../game"
import "../engine/search"

type Player interface {
	MakeMove(*game.Board) error
}

// AIPlayer is a player that makes moves according to AI.
type AIPlayer struct {
	Evaluator game.Evaluator
	Depth     int
	Color     game.Color
}

func (p *AIPlayer) MakeMove(b *game.Board) error {
	start := time.Now()
	km := game.NewKillerMoves()
	// Use iterative deepening to try and find good paths early. It's likely that
	// the best move on ply 1 is the best on ply 2. This fills the transposition table
	// to lead with the best move on future plies.
	var eval float64
	var move game.EfficientMove
	var nodes int
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	d := 1
	for d <= p.Depth {
		eval, move, nodes = search.AlphaBetaSearch(b, p.Evaluator, d, alpha, beta, false, p.Color, km)
		fmt.Println(fmt.Sprintf("iteration %v: best move is %v (%v nodes searched)", d, move, nodes))
		d++
	}
	t := time.Since(start)
	fmt.Println(fmt.Sprintf("evaluation over in: %v", t))
	if move == game.EfficientMove(0) {
		return errors.New("no move could be made")
	}
	// Convert eval to + for white, - for black.
	if p.Color == game.BLACK {
		eval = -1 * eval
	}
	fmt.Println(fmt.Sprintf("AI Player making best move with depth %v: %v, eval %v", p.Depth, move, eval))

	PrintPrincipalVariation(b)
	game.ApplyMove(b, move)
	return nil
}

// CommandLinePlayer is a player that makes moves according to input from the command line.
type CommandLinePlayer struct {
	Color game.Color
}

func (p *CommandLinePlayer) MakeMove(b *game.Board) error {
	reader := bufio.NewReader(os.Stdin)
	// Compare legal moves against the input choice.
	moves := b.AllLegalMoves()
	var move game.EfficientMove
	foundMove := false
	for !foundMove {
		candidates := []game.EfficientMove{}
		fmt.Println("Please input a move. What square is the piece you would like to move? (for castling, start with the king)")
		b, _, _ := reader.ReadLine()
		from := string(b)
		for _, m := range moves {
			if m.Old().String() == from {
				candidates = append(candidates, m)
			}
		}
		if len(candidates) == 0 {
			fmt.Println(fmt.Sprintf("No legal moves start with square %v. Please try again.", from))
			continue
		}
		fmt.Println("What square would you like to move to? (for castling, move the king)")
		b, _, _ = reader.ReadLine()
		to := string(b)
		for _, c := range candidates {
			if c.Square().String() == to {
				move = c
				foundMove = true
			}
		}
		if foundMove {
			break
		}
		fmt.Println(fmt.Sprintf("No legal moves from %v to %v. Please try again.", from, to))
	}
	game.ApplyMove(b, move)
	return nil
}

// Print principal variation prints the expected best continuation
// from a given board.
func PrintPrincipalVariation(b *game.Board) {
	moves := []game.EfficientMove{}
	bs := []game.BoardState{}
	seenMoves := make(map[game.EfficientMove]bool)
	// Get the principal variation, change board state.
	for {
		entry, ok := game.TranspositionTable[game.ZobristHash(b)]
		if !ok || entry.BestMove == game.EfficientMove(0) {
			break
		}
		moves = append(moves, entry.BestMove)
		// Break after first repetition
		if seenMoves[entry.BestMove] {
			break
		}

		seenMoves[entry.BestMove] = true
		b.SwitchActivePlayer()
		bs = append(bs, game.ApplyMove(b, entry.BestMove))
	}
	// Undo board state.
	for i := len(moves) - 1; i >= 0; i-- {
		game.UndoMove(b, moves[i], bs[i])
		b.SwitchActivePlayer()
	}
	// Print the principal variation.
	entry := game.TranspositionTable[game.ZobristHash(b)]
	fmt.Printf("\nEvaluation at Depth %v: %v\n", entry.Depth, entry.Eval)
	fmt.Println("Principal Variation: ")
	pvStrings := []string{}
	for _, m := range moves {
		pvStrings = append(pvStrings, m.String())
	}
	fmt.Println(strings.Join(pvStrings, " "))
}
