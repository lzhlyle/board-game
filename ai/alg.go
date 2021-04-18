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

//type ISimilar interface {
//	GenSimilar(base [][]*core.PlaySignal) ([]*Similar, error)
//}

//type Similar struct {
//	Zip interface{}
//	how int8 // 3 bits used
//	// 最低两位 00:0/01:90/10:180/11:270
//	// 第三位 0:not-flip/1:flip
//}

//type Angle int8
//
//const (
//	Angle0 Angle = 0b_00 + iota
//	Angle90
//	Angle180
//	Angle270
//)
//
//func NewSimilar(zip interface{}, flip bool, angle Angle) *Similar {
//	how := int8(angle)
//	if flip {
//		how |= 1 << 2
//	}
//	return &Similar{Zip: zip, how: how}
//}

type NextRates struct {
	NextZip interface{}
	Rates   [3]int // {win, draw, lose}
}

type ICalculate interface {
	Calculate(curr [][]*core.PlaySignal) (i, j int, err error)
}
