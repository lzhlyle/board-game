package core

type BoardGame interface {
	Board() Board
	Players() []*Player
	RoundEnd(snapshot *MoveSnapshot) bool
	GameEnd(snapshot *MoveSnapshot) (end bool, winner *Player)
}
