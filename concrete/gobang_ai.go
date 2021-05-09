package concrete

import (
	"board-game/ai/ai_impl"
)

type GobangAI struct {
	*Gobang
	*ai_impl.DefaultAIImpl
}

func NewGobangAI() *GobangAI {
	res := &GobangAI{
		Gobang: NewGobang(),
	}

	res.DefaultAIImpl = ai_impl.NewDefaultAIImpl(res.Gobang.players)

	return res
}
