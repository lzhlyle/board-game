package main

import (
	"board-game/concrete"
	"board-game/core"
)

func main() {
	game := concrete.NewGobang()
	core.NewPlay(game).Play()
}
