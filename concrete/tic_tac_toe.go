package concrete

import (
	"board-game/core"
)

type TicTacToe struct {
	board   *core.Board
	players []*core.Player
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

func (t *TicTacToe) Board() *core.Board {
	return t.board
}

func (t *TicTacToe) Players() []*core.Player {
	return t.players
}

func (t *TicTacToe) StartPlayerSequence(lastStarter *core.Player, winner *core.Player, players []*core.Player) []*core.Player {
	if lastStarter == nil {
		return players
	}
	// 轮流开局
	players[0], players[1] = players[1], players[0]
	return players
}

func (t *TicTacToe) NextPlayer(last *core.Player) *core.Player {
	if last == t.players[0] {
		return t.players[1]
	}
	return t.players[0]
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
