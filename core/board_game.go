package core

type BoardGame interface {
	IBoard
	IPlayerCollection
	IGameRule
	ICompress
}

type ICompress interface {
	Compress(mat [][]*PlaySignal) interface{}
	GenSimilar(base [][]*PlaySignal) interface{}
}

type IBoard interface {
	Board() *Board
}

type IPlayerCollection interface {
	Players() []*Player
}

type IGameRule interface {
	RoundEnd(snapshot *MoveSnapshot) bool
	GameEnd(snapshot *MoveSnapshot) (end bool, winner *Player)
}
