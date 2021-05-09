package concrete

import (
	"board-game/ai"
	ai_impl2 "board-game/ai_impl"
	"board-game/core"
	"math"
	"math/rand"
	"sort"
)

type TicTacToeAI struct {
	*TicTacToe
	*ai_impl2.ChessRecordGenerator
	*ai_impl2.DefaultAIImpl
}

const AIStrategySmart ai.Strategy = 100 + iota

func NewTicTacToeAI() *TicTacToeAI {
	res := &TicTacToeAI{
		TicTacToe: NewTicTacToe(),
	}

	res.ChessRecordGenerator = ai_impl2.NewChessRecordGenerator(res, res, res)

	res.DefaultAIImpl = ai_impl2.NewDefaultAIImpl(res.TicTacToe.players)
	res.DefaultAIImpl.RegisterStrategy(AIStrategySmart, res.smartStrategy)
	res.DefaultAIImpl.SetCurrentStrategy(AIStrategySmart)

	return res
}

func (t *TicTacToeAI) smartStrategy(curr [][]*core.PlaySignal) (i, j int, err error) {
	currZip := t.Zip(curr).(int32)
	rates := t.ChessRecordGenerator.Zip2NextRates[currZip]
	if len(rates) == 0 {
		return -1, -1, ai.ErrCannotMove
	}

	rate := rates[0]
	nextZip := rate.NextZip.(int32) // 默认

	// 相近结果的多种走法，增加趣味性
	if rates[0].Rates[0] < 100 {
		// 非必赢
		for i := 1; i < len(rates); i++ {
			// 输率在 10% 以内，则可任选
			if math.Abs(float64(rates[i].Rates[2])-float64(rate.Rates[2])) < 10 {
				nextZip = rates[rand.Intn(i)].NextZip.(int32)
			}
		}
	}

	// 确定差异
	diff := currZip ^ nextZip
	// 求差异的格子值
	cell := -1
	for diff > 0 { // 最多 9 次
		diff >>= 2 // 2 = len(players)
		cell++
	}
	// 确定格子的点位
	return cell / 3, cell % 3, nil
}

func (t *TicTacToeAI) Zip(mat [][]*core.PlaySignal) interface{} {
	// 位运算思路
	// 共9个格子，每个格子3种情况
	// 2个位可表示4种情况(00, 01, 10, 11)
	// 18个位即可表示一个棋盘，返回 int32
	res := int32(0b_00_00_00_00_00_00_00_00_00) // 最低2位(右)表示(0, 0)格子
	for i, row := range mat {
		for j, sgn := range row {
			switch sgn {
			case t.TicTacToe.players[0].Signal:
				res |= 1 << (2 * (3*i + j))
			case t.TicTacToe.players[1].Signal:
				res |= 1 << ((2 * (3*i + j)) + 1)
			}
		}
	}
	return res
}

func (t *TicTacToeAI) SortRecords(rates []*ai.NextRates) {
	sort.Slice(rates, func(i, j int) bool {
		// 若输率差在 20% 以内，则优先胜率更高的
		if math.Abs(float64(rates[i].Rates[2])-float64(rates[i-1].Rates[2])) < 20 {
			return rates[i].Rates[0] > rates[j].Rates[0]
		}

		// 优先输率更低的
		return rates[i].Rates[2] < rates[j].Rates[2]
	})
}
