package main

import (
	"board-game/concrete"
	"board-game/core"
)

func main() {
	bg := concrete.NewTicTacToe()
	core.NewPlay(bg).Play()
}
