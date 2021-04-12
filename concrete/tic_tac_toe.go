package concrete

import "board-game/core"

type TicTacToe struct {
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{}
}

func (t *TicTacToe) Board() core.Board {
	return core.Board{
		Width:        3,
		Height:       3,
		MoveLocStyle: core.MoveLocStyle_InCell,
	}
}

func (t *TicTacToe) Players() []*core.Player {
	return []*core.Player{
		core.NewPlayer("X"),
		core.NewPlayer("O"),
	}
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
