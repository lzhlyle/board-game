package ai

import "board-game/core"

type IAlg interface {
	ICompress
}

type ICompress interface {
	Compress(mat [][]*core.PlaySignal) interface{}
	GenSimilar(base [][]*core.PlaySignal) (interface{}, error)
}
