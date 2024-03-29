package concrete

import (
	"board-game/core"
)

type Gobang struct {
	board   *core.Board
	players []*core.Player
}

func NewGobang() *Gobang {
	return &Gobang{
		board: &core.Board{
			Width:        15,
			Height:       15,
			MoveLocStyle: core.MoveLocStyle_InCross,
		},
		players: []*core.Player{
			core.NewPlayer("X"),
			core.NewAIPlayer("O"),
		},
	}
}

func (g *Gobang) StartPlayerSequence(lastStarter *core.Player, winner *core.Player, players []*core.Player) []*core.Player {
	return players
}

func (g *Gobang) Board() *core.Board {
	return g.board
}

func (g *Gobang) Players() []*core.Player {
	return g.players
}

func (g *Gobang) NextPlayer(last *core.Player) *core.Player {
	if last == g.players[0] {
		return g.players[1]
	}
	return g.players[0]
}

func (g *Gobang) RoundEnd(snapshot *core.MoveSnapshot) bool {
	return true
}

var (
	dirGroup = [][][]int{
		{
			{0, 1},  // up
			{0, -1}, // down
		},
		{
			{-1, 0}, // left
			{1, 0},  // right
		},
		{
			{-1, -1}, // left-down
			{1, 1},   // right-up
		},
		{
			{-1, 1}, // left-up
			{1, -1}, // right-down
		},
	}
)

func (g *Gobang) dfs(mat [][]*core.PlaySignal, target *core.PlaySignal, i, j, dx, dy, cnt int) int {
	if i < 0 || j < 0 || i >= len(mat) || j >= len(mat[0]) {
		return cnt
	}
	if mat[i][j] != target {
		return cnt
	}
	return g.dfs(mat, target, i+dx, j+dy, dx, dy, cnt+1)
}

func (g *Gobang) GameEnd(snapshot *core.MoveSnapshot) (end bool, winner *core.Player) {
	target := snapshot.Board[snapshot.I][snapshot.J]
	for _, group := range dirGroup {
		tot := 0
		for _, dir := range group {
			tot += g.dfs(snapshot.Board, target, snapshot.I, snapshot.J, dir[0], dir[1], 0)
			if tot > 5 {
				return true, snapshot.Player
			}
		}
	}

	// 未结束
	for _, row := range snapshot.Board {
		for _, sgn := range row {
			if sgn == nil {
				return false, nil
			}
		}
	}

	// 平局
	return true, nil
}

func (g *Gobang) IsGameOver(curr [][]*core.PlaySignal) bool {
	// 检查棋盘是否填满
	isFull := true
	for i := 0; i < len(curr); i++ {
		for j := 0; j < len(curr[i]); j++ {
			if curr[i][j] == nil {
				isFull = false
				break
			}
		}
		if !isFull {
			break
		}
	}
	if isFull {
		return true
	}

	// 检查是否有一方连成了五子
	for i := 0; i < len(curr); i++ {
		for j := 0; j < len(curr[i]); j++ {
			if curr[i][j] != nil {
				if g.checkWin(curr, i, j) {
					return true
				}
			}
		}
	}

	return false
}

// 辅助函数，用于检查某个位置是否有五子连线
func (g *Gobang) checkWin(curr [][]*core.PlaySignal, x, y int) bool {
	tag := curr[x][y].Tag
	dirs := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for _, dir := range dirs {
		count := 1
		for i := 1; i < 5; i++ {
			x1 := x + i*dir[0]
			y1 := y + i*dir[1]
			if x1 < 0 || x1 >= len(curr) || y1 < 0 || y1 >= len(curr[x1]) || curr[x1][y1] == nil || curr[x1][y1].Tag != tag {
				break
			}
			count++
		}
		if count >= 5 {
			return true
		}
	}
	return false
}
