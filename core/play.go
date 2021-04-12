package core

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type Play struct {
	step       int
	gameState  GameState
	currPlayer *Player
	snapshot   *MoveSnapshot
	init       bool // 是否已初始化

	param   BoardGame
	board   Board
	players []*Player
}

func NewPlay(param BoardGame) *Play {
	return (&Play{
		param:   param,
		board:   param.Board(),
		players: param.Players(),
	}).reset()
}

func (p *Play) reset() *Play {
	p.step = 0
	p.gameState = GameState_Ready
	p.currPlayer = p.param.Players()[0]
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

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, p.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, p.move); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

var (
	moveLocStyleMap = map[MoveLocStyle]*MoveLocStylePainting{
		MoveLocStyle_InCell: {
			defStr:    "   ",
			locStrFmt: "  %s",
			side:      2,
			frame:     true,
		},
		MoveLocStyle_InCross: {
			defStr:    "- | -",
			locStrFmt: "- %s -",
			side:      2,
			frame:     false,
		},
	}
)

func (p *Play) layout(g *gocui.Gui) error {
	if p.init {
		return nil
	}

	painting := moveLocStyleMap[p.board.MoveLocStyle]

	// x坐标：左右，y坐标：上下
	side := painting.side
	maxX, maxY := g.Size()
	x0, y0 := maxX/2-6*side, maxY/2-2*side

	for x := 0; x < p.board.Width; x++ {
		for y := 0; y < p.board.Height; y++ {
			name := p.cellToName(y, x) // (y, x) 即 (i, j)
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

func (p *Play) cellToName(i, j int) string {
	return fmt.Sprintf("cell-%d-%d", i, j)
}

func (p *Play) nameToCell(name string) (i, j int) {
	arr := strings.Split(name, "-")
	i, _ = strconv.Atoi(arr[1])
	j, _ = strconv.Atoi(arr[2])
	return
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

	signal := p.currPlayer.signal.Tag
	v.Clear()
	_, _ = fmt.Fprintf(v, painting.locStrFmt, signal)

	i, j := p.nameToCell(v.Name())
	p.snapshot = NewGameSnapshot(p.step, i, j, p.currPlayer, p.snapshot)

	g.Update(func(gui *gocui.Gui) error {
		end, winner := p.param.GameEnd(p.snapshot)
		if end {
			return p.win(g, winner)
		}
		return nil
	})

	if p.param.RoundEnd(p.snapshot) {
		p.step++
		p.currPlayer = p.players[p.step%len(p.players)]
	}

	return nil
}

func (p *Play) win(g *gocui.Gui, winner *Player) error {
	p.gameState = GameState_End
	time.Sleep(500 * time.Millisecond)
	x0, y0, x1, y1, err := g.ViewPosition("cell-1-1")
	if err != nil {
		return err
	}

	dx, dy := (x1-x0)/2, (y1-y0)/2
	v, err := g.SetView("win", x0-dx, y0-dy, x1+dx, y1+dy)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if winner != nil {
		_, _ = fmt.Fprintf(v, "%s win!", winner.signal.Tag)
	} else {
		_, _ = fmt.Fprint(v, "draw!")
	}

	_ = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, p.restart)
	_, _ = g.SetCurrentView(v.Name())
	return nil
}

func (p *Play) restart(g *gocui.Gui, v *gocui.View) error {
	_ = g.DeleteView(v.Name())
	for _, v := range g.Views() {
		v.Clear()
	}
	p.reset()
	_ = g.DeleteKeybinding("", gocui.KeyEnter, gocui.ModNone)
	return nil
}

func (p *Play) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
