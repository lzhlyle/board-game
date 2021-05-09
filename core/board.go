package core

// 棋盘
type Board struct {
	Width, Height int
	MoveLocStyle  MoveLocStyle
}

// 落子样式
type MoveLocStyle int8

const (
	MoveLocStyle_InCell  MoveLocStyle = 1 + iota // 格内落子，如井字棋、飞行棋
	MoveLocStyle_InCross                         // 交界落子，如五子棋、围棋
)

type MoveLocStylePainting struct {
	defStr, locStrFmt string
	side              int
	frame             bool
}
