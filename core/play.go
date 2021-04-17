package core

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

type Play struct {
	step       int
	gameState  GameState
	currPlayer *Player
	snapshot   *MoveSnapshot
	init       bool // 是否已初始化

	starter *Player
	winner  *Player

	rule       IGameRule
	board      *Board
	players    []*Player
	player2Idx map[*Player]int

	hook Hook
}

func NewPlay(bg BoardGame) *Play {
	res := &Play{
		rule:       bg,
		board:      bg.Board(),
		players:    bg.Players(),
		player2Idx: make(map[*Player]int, len(bg.Players())),
	}
	for i, player := range res.players {
		res.player2Idx[player] = i
	}
	if v, ok := bg.(Hook); ok {
		res.hook = v
	}
	return res.reset()
}

func (p *Play) reset() *Play {
	p.step = 0
	p.gameState = GameState_Ready

	p.currPlayer = p.rule.GameStart(p.starter, p.winner, p.players, p.player2Idx)
	p.starter = p.currPlayer
	p.winner = nil

	p.snapshot = NewMoveSnapshot(p.board.Width, p.board.Height)
	p.init = false
	return p
}

func (p *Play) Play() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(p.layout)
	g.Mouse = true
	g.InputEsc = true

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, p.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, p.restart); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, p.move); err != nil {
		log.Panicln(err)
	}

	p.moveIfAI(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (p *Play) moveIfAI(g *gocui.Gui) {
	if p.currPlayer.ai && p.hook != nil {
		if fn := p.hook.AIMove(p.snapshot, p.move); fn != nil {
			g.Update(fn)
		}
	}
}

func (p *Play) layout(g *gocui.Gui) error {
	if p.init {
		return nil
	}

	painting := moveLocStyleMap[p.board.MoveLocStyle]

	// x坐标：左右，y坐标：上下
	side := painting.side
	maxX, maxY := g.Size()
	x0, y0 := maxX/2-6*(p.board.Width/2), maxY/2-2*(p.board.Height/2)

	for x := 0; x < p.board.Width; x++ {
		for y := 0; y < p.board.Height; y++ {
			name := CellToName(y, x) // (y, x) 即 (i, j)
			v, err := g.SetView(name, x0+x*3*side, y0+y*side, x0+x*3*side+3*side, y0+y*side+side)
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Frame = painting.frame
			_, _ = fmt.Fprint(v, painting.defStr)
		}
	}
	p.init = true
	return nil
}

func (p *Play) move(g *gocui.Gui, v *gocui.View) error {
	if p.gameState == GameState_End {
		return nil
	} else {
		p.gameState = GameState_InProcess
	}

	painting := moveLocStyleMap[p.board.MoveLocStyle]
	if v.Buffer() != painting.defStr+"\n" {
		return nil
	}

	signal := p.currPlayer.Signal.Tag
	v.Clear()
	_, _ = fmt.Fprintf(v, painting.locStrFmt, signal)

	i, j := NameToCell(v.Name())
	p.snapshot = NewGameSnapshot(p.step, i, j, p.currPlayer, p.snapshot)

	g.Update(func(g *gocui.Gui) error {
		end, winner := p.rule.GameEnd(p.snapshot)
		if end {
			return p.win(g, winner)
		}
		return nil
	})

	p.step++
	if p.rule.RoundEnd(p.snapshot) {
		p.currPlayer = p.players[(p.player2Idx[p.currPlayer]+1)%len(p.players)]
	}

	p.moveIfAI(g)

	return nil
}

func (p *Play) win(g *gocui.Gui, winner *Player) error {
	p.winner = winner
	p.gameState = GameState_End
	time.Sleep(500 * time.Millisecond)

	x0, y0, x1, y1, err := calcWinViewLocation(p.board.Width, p.board.Height, g)
	if err != nil {
		return err
	}
	v, err := g.SetView("win", x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if winner != nil {
		_, _ = fmt.Fprintf(v, "%s win!", winner.Signal.Tag)
	} else {
		_, _ = fmt.Fprint(v, "draw!")
	}

	_ = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, p.restart)
	_ = g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_ = g.DeleteView(v.Name())
		_ = g.DeleteKeybinding("", gocui.KeyEsc, gocui.ModNone)
		return nil
	})
	_, _ = g.SetCurrentView(v.Name())
	return nil
}

func (p *Play) restart(g *gocui.Gui, v *gocui.View) error {
	_ = g.DeleteView(v.Name())
	for _, v := range g.Views() {
		v.Clear()
	}
	_ = g.DeleteKeybinding("", gocui.KeyEnter, gocui.ModNone)
	p.reset()

	p.moveIfAI(g)

	return nil
}

func (p *Play) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
