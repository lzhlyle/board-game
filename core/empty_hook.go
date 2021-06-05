package core

import "github.com/jroimartin/gocui"

type EmptyHook struct {
}

func (*EmptyHook) AIMove(snapshot *MoveSnapshot, moveFn MoveFn) UpdateFn {
	return nil
}

func (*EmptyHook) BeforeRound(g *gocui.Gui) {

}

func (*EmptyHook) AfterRound(g *gocui.Gui, snapshot *MoveSnapshot) {

}
