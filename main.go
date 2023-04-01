package main

import (
	"board-game/concrete"
	"board-game/core"
)

func main() {
	//game := concrete.NewGobang()
	game := concrete.NewGobangAI()
	//game := concrete.NewTicTacToe()
	//game := concrete.NewTicTacToeAI()
	core.NewPlay(game).Play()
}
