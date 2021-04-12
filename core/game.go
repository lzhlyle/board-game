package core

// 整盘棋
type Game struct {
	state GameState
}

// 游戏状态
type GameState int8

const (
	GameState_Ready     GameState = 1 + iota // 已就绪
	GameState_InProcess                      // 进行中
	GameState_End                            // 已结束
)
