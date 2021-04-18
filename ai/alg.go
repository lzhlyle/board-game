package ai

import "board-game/core"

type IAlg interface {
	IZip
	ISimilar
	ICalculate
}

type IZip interface {
	Zip(mat [][]*core.PlaySignal) interface{}
}

type ISimilar interface {
	GenSimilar(base [][]*core.PlaySignal) (interface{}, error)
}

type NextRates struct {
	NextZip interface{}
	Rates   [3]int // {win, draw, lose}
}

type ICalculate interface {
	Calculate(curr [][]*core.PlaySignal) (i, j int, err error)
}
