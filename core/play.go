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

	rule       IGameRule
	board      *Board
	players    []*Player
	compressor ICompress
}

func NewPlay(bg BoardGame) *Play {
	return (&Play{
		rule:       bg,
		board:      bg.Board(),
		players:    bg.Players(),
		compressor: bg,
	}).reset()
}

func (p *Play) reset() *Play {
	p.step = 0
	p.gameState = GameState_Ready
	p.currPlayer = p.players[0]
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
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, p.restart); err != nil {
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
			defStr:    "——|——",
			locStrFmt: "——%s——",
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
	x0, y0 := maxX/2-6*(p.board.Width/2), maxY/2-2*(p.board.Height/2)

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

	signal := p.currPlayer.Signal.Tag
	v.Clear()
	_, _ = fmt.Fprintf(v, painting.locStrFmt, signal)

	i, j := p.nameToCell(v.Name())
	p.snapshot = NewGameSnapshot(p.step, i, j, p.currPlayer, p.snapshot)

	g.Update(func(gui *gocui.Gui) error {
		end, winner := p.rule.GameEnd(p.snapshot)
		if end {
			return p.win(g, winner)
		}
		return nil
	})

	if p.rule.RoundEnd(p.snapshot) {
		p.step++
		p.currPlayer = p.players[p.step%len(p.players)]
	}

	return nil
}

func (p *Play) win(g *gocui.Gui, winner *Player) error {
	p.gameState = GameState_End
	time.Sleep(500 * time.Millisecond)

	x0, y0, x1, y1, err := p.calcWinViewLocation(p.board.Width, p.board.Height, g)
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
	_, _ = g.SetCurrentView(v.Name())
	return nil
}

func (p *Play) calcWinViewLocation(w, h int, g *gocui.Gui) (wx0, wy0, wx1, wy1 int, err error) {
	// center left-top
	ltname := fmt.Sprintf("cell-%d-%d", (h-1)/2, (w-1)/2)
	xc0, yc0, xc1, yc1, err := g.ViewPosition(ltname)
	if err != nil {
		return
	}
	if (w&1) == 0 || (h&1) == 0 {
		// center right-bottom
		rbname := fmt.Sprintf("cell-%d-%d", (h^1)/2, (w^1)/2)
		_, _, xc1, yc1, err = g.ViewPosition(rbname)
		if err != nil {
			return
		}
	}
	// center: from (xc0, yc0) to (xc1, yc1)

	// left-top cell
	x0, y0, _, _, err := g.ViewPosition("cell-0-0")
	if err != nil {
		return
	}

	// right-bottom cell
	_, _, x1, y1, err := g.ViewPosition(fmt.Sprintf("cell-%d-%d", h-1, w-1))
	if err != nil {
		return
	}

	// calculate delta
	dx0, dy0 := (xc0-x0)/2, (yc0-y0)/2
	// adjust
	if dx0 == xc0-x0 {
		dx0 = (xc1 - xc0) / 2
	}
	if dy0 == yc0-y0 {
		dy0 = (yc1 - yc0) / 2
	}
	dx1, dy1 := (x1-xc1)/2, (y1-yc1)/2
	// adjust
	if dx1 == x1-xc1 {
		dx1 = (xc1 - xc0) / 2
	}
	if dy1 == y1-yc1 {
		dy1 = (yc1 - yc0) / 2
	}

	return xc0 - dx0, yc0 - dy0, xc1 + dx1, yc1 + dy1, nil
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
