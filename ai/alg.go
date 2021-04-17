package ai

import "board-game/core"

type IAlg interface {
	ICompress
	ISimilar
	ISelect
}

type ICompress interface {
	Compress(mat [][]*core.PlaySignal) interface{}
}

type ISimilar interface {
	GenSimilar(base [][]*core.PlaySignal) (interface{}, error)
}

type NextRates struct {
	nextZip interface{}
	rates   [3]int // {win, draw, lose}
}

type ISelect interface {
	Select(curr [][]*core.PlaySignal) (i, j int, err error)
}
