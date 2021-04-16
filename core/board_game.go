package core

type BoardGame interface {
	IBoard
	IPlayerCollection
	IGameRule
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
