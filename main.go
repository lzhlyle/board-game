package main

import (
	"board-game/concrete"
	"board-game/core"
)

func main() {
	game := concrete.NewTicTacToe()
	core.NewPlay(game).Play()
}
