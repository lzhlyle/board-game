package ai

import "board-game/core"

type IAlg interface {
	IZip
	//ISimilar
	ICalculate
}

type IZip interface {
	Zip(mat [][]*core.PlaySignal) interface{}
}

type NextRates struct {
	NextZip interface{}
	Rates   [3]int // {win, draw, lose}
}

func NewNextRates(nextZip interface{}, rates [3]int) *NextRates {
	return &NextRates{NextZip: nextZip, Rates: rates}
}

type ICalculate interface {
	Calculate(curr [][]*core.PlaySignal) (i, j int, err error)
}
