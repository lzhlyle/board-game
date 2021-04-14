package core

// 单个玩家
type Player struct {
	Signal *PlaySignal
	state  PlayState
}

func NewPlayer(tag string) *Player {
	return &Player{
		Signal: &PlaySignal{
			Tag: tag,
		},
		state: PlayState_Ready,
	}
}

// 玩家标识
type PlaySignal struct {
	Tag string
}

// 玩家状态
type PlayState int8

const (
	PlayState_Ready     PlayState = 1 + iota // 已就绪
	PlayState_InProcess                      // 进行中
	PlayState_Victory                        // 已胜利
	PlayState_Over                           // 已失败
)
