package core

type BoardGame interface {
	IBoard
	IPlayerCollection
	GameRule
}

type IBoard interface {
	Board() Board
}

type IPlayerCollection interface {
	Players() []*Player
}

type GameRule interface {
	RoundEnd(snapshot *MoveSnapshot) bool
	GameEnd(snapshot *MoveSnapshot) (end bool, winner *Player)
}
