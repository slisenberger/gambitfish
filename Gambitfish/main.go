// Package main initializes the necessary game state and prints the
// current state of the art of the engine.
package main

import "./game"

func main() {
	b := game.DefaultBoard()
	b.Print()
}
