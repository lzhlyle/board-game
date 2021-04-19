package concrete

import (
	"board-game/ai"
	"board-game/core"
	"github.com/jroimartin/gocui"
	"math"
	"math/rand"
	"sort"
	"time"
)

type TicTacToe struct {
	board   *core.Board
	players []*core.Player

	// Alg
	// (preZip, nextRates)
	zip2NextRates map[int32][]*ai.NextRates
}

func NewTicTacToe() *TicTacToe {
	return (&TicTacToe{
		board: &core.Board{
			Width:        3,
			Height:       3,
			MoveLocStyle: core.MoveLocStyle_InCell,
		},
		players: []*core.Player{
			core.NewPlayer("X"),
			core.NewAIPlayer("O"),
		},
		zip2NextRates: make(map[int32][]*ai.NextRates),
	}).buildChessRecord()
}

func (t *TicTacToe) buildChessRecord() *TicTacToe {
	// later @lzh 可考虑多协程并发处理，注意将 zip2NextRates 改为支持并发的数据结构
	t.dfs(core.NewEmptyMoveSnapshot(t.board.Width, t.board.Height), t.zip2NextRates)
	return t
}

// dfs 递归搜索棋谱
// 返回 aSnapShot 的胜率情况
func (t *TicTacToe) dfs(aSnapshot *core.MoveSnapshot, zip2NextRates map[int32][]*ai.NextRates) *ai.NextRates {
	aZip := t.Zip(aSnapshot.Board).(int32)

	// terminator
	if end, winner := t.GameEnd(aSnapshot); end {
		zip2NextRates[aZip] = []*ai.NextRates{}
		if winner == nil {
			return ai.NewNextRates(aZip, [3]int{0, 100, 0})
		}
		return ai.NewNextRates(aZip, [3]int{100, 0, 0})
	}

	// look up possibles
	b := t.NextPlayer(aSnapshot.Player)
	bPossibles := make([][2]int, 0) // [2]int: {i, j}
	for i, row := range aSnapshot.Board {
		for j := 0; j < len(row); j++ {
			if row[j] == nil {
				bPossibles = append(bPossibles, [2]int{i, j})
			}
		}
	}

	allBNextRates := make([]*ai.NextRates, len(bPossibles))
	// travel possibles
	var aRates = [3]int{}
	for i, pos := range bPossibles {
		bSnapshot := core.GenSnapshot(aSnapshot.Step+1, pos[0], pos[1], b, aSnapshot)
		bNextRates := t.dfs(bSnapshot, zip2NextRates)
		allBNextRates[i] = bNextRates

		// cross accumulate
		// 零和游戏，交叉累计
		bRates := bNextRates.Rates
		aRates[0] += bRates[2]
		aRates[1] += bRates[1]
		aRates[2] += bRates[0]
	}
	// average
	for i := 0; i < 3; i++ {
		aRates[i] = aRates[i] / len(bPossibles)
	}

	zip2NextRates[aZip] = t.sort(t.filter(allBNextRates))
	return ai.NewNextRates(aZip, aRates)
}

func (t *TicTacToe) filter(rates []*ai.NextRates) []*ai.NextRates {
	// 能赢则赢
	res := make([]*ai.NextRates, 0, len(rates))
	for _, rate := range rates {
		if rate.Rates[0] == 100 {
			// 多几种选择，而非直接返回
			res = append(res, rate)
		}
	}
	if len(res) > 0 {
		return res
	}

	// 对方再下一步不会赢，则才可走
	for _, rate := range rates {
		nextWillWin := false
		for _, next := range t.zip2NextRates[rate.NextZip.(int32)] {
			if next.Rates[0] == 100 {
				nextWillWin = true
				break
			}
		}
		if !nextWillWin {
			res = append(res, rate)
		}
	}
	if len(res) == 0 {
		// 必输，放弃治疗
		return rates
	}

	return res
}

func (t *TicTacToe) sort(rates []*ai.NextRates) []*ai.NextRates {
	if len(rates) == 0 {
		return rates
	}

	if rates[0].Rates[0] == 100 {
		return rates // 能赢，无需排序
	}

	// 策略优先级排序
	sort.Slice(rates, func(i, j int) bool {
		// 若输率差在 20% 以内，则优先胜率更高的
		if math.Abs(float64(rates[i].Rates[2])-float64(rates[i-1].Rates[2])) < 20 {
			return rates[i].Rates[0] > rates[j].Rates[0]
		}

		// 优先输率更低的
		return rates[i].Rates[2] < rates[j].Rates[2]
	})
	return rates
}

func (t *TicTacToe) NextPlayer(last *core.Player) *core.Player {
	if last == t.players[0] {
		return t.players[1]
	}
	return t.players[0]
}

type AIStrategy int

const (
	AIStrategyRandom AIStrategy = 1 + iota
	AIStrategyHighestWin
	AIStrategyLowestLose
)

func (t *TicTacToe) StartPlayerSequence(lastStarter *core.Player, winner *core.Player, players []*core.Player) []*core.Player {
	if lastStarter == nil {
		return players
	}
	// 轮流开局
	players[0], players[1] = players[1], players[0]
	return players
}

func (t *TicTacToe) Calculate(curr [][]*core.PlaySignal) (i, j int, err error) {
	currZip := t.Zip(curr).(int32)
	rates := t.zip2NextRates[currZip]
	if len(rates) == 0 {
		return -1, -1, ai.ErrCannotMove
	}

	rate := rates[0]
	nextZip := rate.NextZip.(int32) // 默认

	// 相近结果的多种走法，增加趣味性
	if rates[0].Rates[0] < 100 {
		// 非必赢
		for i := 1; i < len(rates); i++ {
			// 输率在 10% 以内，则可任选
			if math.Abs(float64(rates[i].Rates[2])-float64(rate.Rates[2])) < 10 {
				nextZip = rates[rand.Intn(i)].NextZip.(int32)
			}
		}
	}

	// 确定差异
	diff := currZip ^ nextZip
	// 求差异的格子值
	cell := -1
	for diff > 0 { // 最多 9 次
		diff >>= 2 // 2 = len(players)
		cell++
	}
	// 确定格子的点位
	return cell / 3, cell % 3, nil

	//// random strategy
	//optionals := make([][2]int, 0)
	//for i := range curr {
	//	for j := range curr[i] {
	//		if curr[i][j] == nil {
	//			optionals = append(optionals, [2]int{i, j})
	//		}
	//	}
	//}
	//if len(optionals) == 0 {
	//	return -1, -1, ai.ErrCannotMove
	//}
	//res := optionals[rand.Intn(len(optionals))]
	//return res[0], res[1], nil
}

func (t *TicTacToe) AIMove(snapshot *core.MoveSnapshot, moveFn core.MoveFn) core.UpdateFn {
	return func(g *gocui.Gui) error {
		i, j, err := t.Calculate(snapshot.Board)
		if err == ai.ErrCannotMove {
			return nil
		} else if err != nil {
			return err
		}

		v, err := g.View(core.CellToName(i, j))
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return moveFn(g, v)
	}
}

// Compress
// 位运算思路
// 共9个格子，每个格子3种情况
// 2个位可表示4种情况(00, 01, 10, 11)
// 18个位即可表示一个棋盘，返回 int32
func (t *TicTacToe) Zip(mat [][]*core.PlaySignal) interface{} {
	res := int32(0b_00_00_00_00_00_00_00_00_00) // 最低2位(右)表示(0, 0)格子
	for i, row := range mat {
		for j, sgn := range row {
			switch sgn {
			case t.players[0].Signal:
				res |= 1 << (2 * (3*i + j))
			case t.players[1].Signal:
				res |= 1 << ((2 * (3*i + j)) + 1)
			}
		}
	}
	return res
}

func (t *TicTacToe) Board() *core.Board {
	return t.board
}

func (t *TicTacToe) Players() []*core.Player {
	return t.players
}

func (t *TicTacToe) RoundEnd(_ *core.MoveSnapshot) bool {
	return true
}

func (t *TicTacToe) GameEnd(snapshot *core.MoveSnapshot) (end bool, winner *core.Player) {
	if snapshot.Player == nil {
		return false, nil
	}

	validRow := func(mat [][]*core.PlaySignal, i int) bool {
		for j := 0; j < 3; j++ {
			if mat[i][j] != mat[i][0] {
				return false
			}
		}
		return true
	}

	validCol := func(mat [][]*core.PlaySignal, j int) bool {
		for i := 0; i < 3; i++ {
			if mat[i][j] != mat[0][j] {
				return false
			}
		}
		return true
	}

	validRToL := func(mat [][]*core.PlaySignal, i, j int) bool {
		if i+j != 2 {
			// 不在「撇」对角线上
			return false
		}
		return mat[0][2] == mat[1][1] && mat[1][1] == mat[2][0]
	}

	validLToR := func(mat [][]*core.PlaySignal, i, j int) bool {
		if i != j {
			// 不在「捺」对角线上
			return false
		}
		return mat[0][0] == mat[1][1] && mat[1][1] == mat[2][2]
	}

	// 校验当前(i, j)的横竖撇捺是否三连
	if validRow(snapshot.Board, snapshot.I) ||
		validCol(snapshot.Board, snapshot.J) ||
		validRToL(snapshot.Board, snapshot.I, snapshot.J) ||
		validLToR(snapshot.Board, snapshot.I, snapshot.J) {
		return true, snapshot.Player
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
