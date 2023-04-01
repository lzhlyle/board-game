package concrete

import (
	"math"

	"board-game/ai"
	ai_impl2 "board-game/ai_impl"
	"board-game/core"
)

type GobangAI struct {
	*Gobang
	*ai_impl2.DefaultAIImpl
}

func NewGobangAI() *GobangAI {
	res := &GobangAI{
		Gobang: NewGobang(),
	}

	res.DefaultAIImpl = ai_impl2.NewDefaultAIImpl(res.Gobang.players)

	const AIStrategySmart ai.Strategy = 100 + iota
	res.DefaultAIImpl.RegisterStrategy(AIStrategySmart, res.alphaBetaStrategy)
	res.DefaultAIImpl.SetCurrentStrategy(AIStrategySmart)

	return res
}

var nextHand = map[string]string{
	"O": "X",
	"X": "O",
}

// AI holds "O" from ChatGPT
func (g *GobangAI) alphaBetaStrategy(curr [][]*core.PlaySignal) (i, j int, err error) {
	// 用正负无穷表示极大值和极小值
	maxVal, minVal := math.Inf(-1), math.Inf(1)
	var bestMove [2]int

	// 搜索最优解
	for i := 0; i < len(curr); i++ {
		for j := 0; j < len(curr[i]); j++ {
			if curr[i][j] == nil {
				// 假设当前空位下棋，计算估值
				curr[i][j] = &core.PlaySignal{Tag: "O"}
				deep = 0
				val := g.minValue(curr, maxVal, minVal, "O")
				curr[i][j] = nil

				// 更新最优解
				if val > maxVal {
					maxVal = val
					bestMove = [2]int{i, j}
				}
			}
		}
	}

	return bestMove[0], bestMove[1], nil
}

func (g *GobangAI) maxValue(curr [][]*core.PlaySignal, alpha, beta float64, tag string) float64 {
	// 检查游戏是否结束，如果是，返回估值
	if g.IsGameOver(curr) {
		return g.evaluate(curr, tag)
	}

	val := math.Inf(-1)

	// 递归搜索所有可能的落子
	for i := 0; i < len(curr); i++ {
		for j := 0; j < len(curr[i]); j++ {
			if curr[i][j] == nil {
				curr[i][j] = &core.PlaySignal{Tag: nextHand[tag]}
				val = math.Max(val, g.minValue(curr, alpha, beta, nextHand[tag]))
				curr[i][j] = nil

				// alpha-beta剪枝
				if val >= beta {
					return val
				}
				alpha = math.Max(alpha, val)
			}
		}
	}

	return val
}

var deep = 0

func (g *GobangAI) minValue(curr [][]*core.PlaySignal, alpha, beta float64, tag string) float64 {
	deep++
	// 检查游戏是否结束，如果是，返回估值
	if g.IsGameOver(curr) {
		return g.evaluate(curr, tag)
	}

	val := math.Inf(1)

	// 递归搜索所有可能的落子
	for i := 0; i < len(curr); i++ {
		for j := 0; j < len(curr[i]); j++ {
			if curr[i][j] == nil {
				curr[i][j] = &core.PlaySignal{Tag: nextHand[tag]}
				val = math.Min(val, g.maxValue(curr, alpha, beta, nextHand[tag]))
				curr[i][j] = nil

				// alpha-beta剪枝
				if val <= alpha {
					return val
				}
				beta = math.Min(beta, val)
			}
		}
	}

	return val
}

// 这特么算数量有什么用？！围棋也不是这么算的好吧
func (g *GobangAI) evaluate(curr [][]*core.PlaySignal, tag string) float64 {
	// 计算当前局面中 X 和 O 的棋子数量
	xCount, oCount := 0, 0
	for i := 0; i < len(curr); i++ {
		for j := 0; j < len(curr[i]); j++ {
			if curr[i][j] != nil {
				if curr[i][j].Tag == "X" {
					xCount++
				} else if curr[i][j].Tag == "O" {
					oCount++
				}
			}
		}
	}

	// 计算当前局面得分
	score := 0
	if tag == "X" {
		score = xCount - oCount
	} else {
		score = oCount - xCount
	}

	// 加上特殊位置得分
	specialScores := map[[2]int]int{
		{7, 7}: 50, {7, 8}: 50, {8, 7}: 50, {8, 8}: 50, // 四个角
		{6, 7}: 10, {7, 6}: 10, {7, 9}: 10, {9, 7}: 10, {8, 8}: 10, // 中心点和中心点周围
		{6, 6}: 5, {6, 8}: 5, {8, 6}: 5, {8, 9}: 5, {9, 6}: 5, {9, 8}: 5, // 四个边角和边角周围
	}
	for pos, s := range specialScores {
		if curr[pos[0]][pos[1]] != nil && curr[pos[0]][pos[1]].Tag == tag {
			score += s
		}
	}

	return float64(score)
}
