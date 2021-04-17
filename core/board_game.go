package core

import "github.com/jroimartin/gocui"

type BoardGame interface {
	IBoard
	IPlayerCollection
	IGameRule
}

type IBoard interface {
	Board() *Board
}

type IPlayerCollection interface {
	Players() []*Player
}

type IGameRule interface {
	GameStart(lastStarter *Player, winner *Player, players []*Player, player2Idx map[*Player]int) *Player
	RoundEnd(snapshot *MoveSnapshot) bool
	GameEnd(snapshot *MoveSnapshot) (end bool, winner *Player)
}

type Hook interface {
	AIMove(snapshot *MoveSnapshot, moveFn MoveFn) UpdateFn
}

type (
	UpdateFn func(gui *gocui.Gui) error
	MoveFn   func(g *gocui.Gui, v *gocui.View) error
)
