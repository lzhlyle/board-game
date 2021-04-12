package main

import (
	"board-game/concrete"
	"board-game/core"
)

func main() {
	param := concrete.NewTicTacToe()
	core.NewPlay(param).Play()
}
