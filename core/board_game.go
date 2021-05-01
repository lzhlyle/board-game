package core

import "github.com/jroimartin/gocui"

type BoardGame interface {
	IBoard
	IPlayerCollection
	IGameRule
}

// IBoard 棋盘
type IBoard interface {
	Board() *Board
}

// IPlayerCollection 玩家集
type IPlayerCollection interface {
	Players() []*Player
	StartPlayerSequence(lastStarter *Player, winner *Player, players []*Player) []*Player
	NextPlayer(last *Player) *Player
}

// IGameRule 游戏规则
type IGameRule interface {
	RoundEnd(snapshot *MoveSnapshot) bool
	GameEnd(snapshot *MoveSnapshot) (end bool, winner *Player)
}

// Hook 钩子
type Hook interface {
	AIMove(snapshot *MoveSnapshot, moveFn MoveFn) UpdateFn
}

type (
	UpdateFn func(gui *gocui.Gui) error
	MoveFn   func(g *gocui.Gui, v *gocui.View) error
)
