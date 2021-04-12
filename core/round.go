package core

// 回合
type Round struct {
	state RoundState
}

// 回合状态
type RoundState int8

const (
	RoundState_Ready     RoundState = 1 + iota // 已就绪
	RoundState_InProcess                       // 进行中
	RoundState_End                             // 已结束
)
