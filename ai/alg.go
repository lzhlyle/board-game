package ai

import "board-game/core"

type IAlg interface {
	ICalculate
}

// ICalculate 走棋计算
type ICalculate interface {
	Calculate(curr [][]*core.PlaySignal) (i, j int, err error)
}

type CalcFunc func(curr [][]*core.PlaySignal) (i, j int, err error)

// IZip 棋盘压缩
type IZip interface {
	Zip(mat [][]*core.PlaySignal) interface{}
	Diff4Cell(curr, next interface{}) (i, j int)
}

// IChessRecord 棋谱
type IChessRecord interface {
	SortRecords(records []*NextRates)
}

// NextRates 下一步胜率集
type NextRates struct {
	NextZip interface{}
	Rates   [3]int // {win, draw, lose}
}

func NewNextRates(nextZip interface{}, rates [3]int) *NextRates {
	return &NextRates{NextZip: nextZip, Rates: rates}
}

type Strategy int
