package main

import (
	"board-game/concrete"
	"board-game/core"
)

func main() {
	param := concrete.NewGobang()
	core.NewPlay(param).Play()
}
