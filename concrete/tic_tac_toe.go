package concrete

import (
	"board-game/ai"
	"board-game/core"
	"github.com/jroimartin/gocui"
	"math/rand"
	"sort"
	"time"
)

type TicTacToe struct {
	board   *core.Board
	players []*core.Player

	// Alg
	curr2Ld      map[int32]int32
	ld2Similar   map[int32][]int32
	ld2NextRates map[int32][]ai.NextRates
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{
		board: &core.Board{
			Width:        3,
			Height:       3,
			MoveLocStyle: core.MoveLocStyle_InCell,
		},
		players: []*core.Player{
			core.NewAIPlayer("X"),
			core.NewAIPlayer("O"),
		},
	}
}

func (t *TicTacToe) GameStart(lastStarter *core.Player, winner *core.Player, players []*core.Player, player2Idx map[*core.Player]int) *core.Player {
	if lastStarter == nil {
		return players[0]
	}

	// 轮流开局
	return players[(player2Idx[lastStarter]+1)%len(players)]
}

func (t *TicTacToe) Select(curr [][]*core.PlaySignal) (i, j int, err error) {
	optionals := make([][2]int, 0)
	for i := range curr {
		for j := range curr[i] {
			if curr[i][j] == nil {
				optionals = append(optionals, [2]int{i, j})
			}
		}
	}
	if len(optionals) == 0 {
		return -1, -1, ai.ErrCannotMove
	}
	res := optionals[rand.Intn(len(optionals))]
	return res[0], res[1], nil
}

func (t *TicTacToe) AIMove(snapshot *core.MoveSnapshot, moveFn core.MoveFn) core.UpdateFn {
	return func(g *gocui.Gui) error {
		i, j, err := t.Select(snapshot.Board)
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
func (t *TicTacToe) Compress(mat [][]*core.PlaySignal) interface{} {
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

// GenSimilar
// 旋转 90度，共 4 种
// 左右翻转后，90度，共 4 种
// 最多共 8 种，还需从中去重
func (t *TicTacToe) GenSimilar(base [][]*core.PlaySignal) (interface{}, error) {
	m := make(map[int32]int8, 8)          // 去重用
	mats := [2][][]*core.PlaySignal{base} // 翻转前后 2 种

	mats[1] = ai.FlipLR(base) // 翻转

	for _, mat := range mats {
		m[t.Compress(mat).(int32)] = 0
		for i := 0; i < 3; i++ {
			mat, _ = ai.SpinSquare90(mat)
			m[t.Compress(mat).(int32)] = 0
		}
	}

	res := make([]int32, 0, len(m))
	for sml := range m {
		res = append(res, sml)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res, nil
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
