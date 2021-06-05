package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

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

func CellToName(i, j int) string {
	return fmt.Sprintf("cell-%d-%d", i, j)
}

func NameToCell(name string) (i, j int) {
	arr := strings.Split(name, "-")
	i, _ = strconv.Atoi(arr[1])
	j, _ = strconv.Atoi(arr[2])
	return
}

func IsCell(name string) bool {
	arr := strings.Split(name, "-")
	return len(arr) == 3
}

func calcWinViewLocation(w, h int, g *gocui.Gui) (wx0, wy0, wx1, wy1 int, err error) {
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
