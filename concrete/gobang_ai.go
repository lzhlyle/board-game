package concrete

import (
	ai_impl2 "board-game/ai_impl"
)

type GobangAI struct {
	*Gobang
	*ai_impl2.DefaultAIImpl
}

func NewGobangAI() *GobangAI {
	res := &GobangAI{
		Gobang: NewGobang(),
	}

	res.DefaultAIImpl = ai_impl2.NewDefaultAIImpl(res.Gobang.players)

	return res
}
