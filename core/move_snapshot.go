package core

type MoveSnapshot struct {
	Step      int             // 最后步数，从 0 起计
	I, J      int             // 最后落子处
	Player    *Player         // 最后落子玩家
	Board     [][]*PlaySignal // 最后棋盘
	Pre, Next *MoveSnapshot   // 上一步、下一步
}

func NewMoveSnapshot(width, height int) *MoveSnapshot {
	board := make([][]*PlaySignal, height)
	for i := range board {
		board[i] = make([]*PlaySignal, width)
	}
	return &MoveSnapshot{Step: 0, Board: board}
}

func NewGameSnapshot(step, i, j int, player *Player, pre *MoveSnapshot) *MoveSnapshot {
	curr := &MoveSnapshot{
		Step:   step,
		I:      i,
		J:      j,
		Board:  CloneMatrix(pre.Board),
		Player: player,
		Pre:    pre,
	}
	pre.Next = curr
	curr.Board[i][j] = player.Signal
	return curr
}
