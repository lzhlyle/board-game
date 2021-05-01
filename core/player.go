package core

// Player 单个玩家
type Player struct {
	Signal *PlaySignal
	AI     bool
	state  PlayState
}

func NewPlayer(tag string) *Player {
	return &Player{
		Signal: &PlaySignal{
			Tag: tag,
		},
		AI:    false,
		state: PlayState_Ready,
	}
}

func NewAIPlayer(tag string) *Player {
	return &Player{
		Signal: &PlaySignal{
			Tag: tag,
		},
		AI:    true,
		state: PlayState_Ready,
	}
}

// PlaySignal 玩家标识
type PlaySignal struct {
	Tag string
}

// PlayState 玩家状态
type PlayState int8

const (
	PlayState_Ready     PlayState = 1 + iota // 已就绪
	PlayState_InProcess                      // 进行中
	PlayState_Victory                        // 已胜利
	PlayState_Over                           // 已失败
)

func CloneMatrix(mat [][]*PlaySignal) [][]*PlaySignal {
	res := make([][]*PlaySignal, len(mat))
	for i, row := range mat {
		res[i] = make([]*PlaySignal, len(row))
		for j, ps := range row {
			res[i][j] = ps
		}
	}
	return res
}
