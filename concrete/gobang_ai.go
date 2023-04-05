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

func (g *GobangAI) evaluate(curr [][]*core.PlaySignal, tag string) float64 {
	var score float64
	var opponentTag string
	if tag == "X" {
		opponentTag = "O"
	} else {
		opponentTag = "X"
	}
	// 计算当前执子方的得分
	score += evaluateDirection(curr, tag, 0, 1)  // 水平方向
	score += evaluateDirection(curr, tag, 1, 0)  // 垂直方向
	score += evaluateDirection(curr, tag, 1, 1)  // 正斜方向
	score += evaluateDirection(curr, tag, -1, 1) // 反斜方向

	// 计算对手的得分
	score -= evaluateDirection(curr, opponentTag, 0, 1)  // 水平方向
	score -= evaluateDirection(curr, opponentTag, 1, 0)  // 垂直方向
	score -= evaluateDirection(curr, opponentTag, 1, 1)  // 正斜方向
	score -= evaluateDirection(curr, opponentTag, -1, 1) // 反斜方向

	return score
}

var countToScore = map[int]float64{
	1: 0.1,
	2: 1,
	3: 10,
	4: 1000,
}

// 评估某个方向上的得分
func evaluateDirection(curr [][]*core.PlaySignal, tag string, dx, dy int) float64 {
	var score float64
	rows, cols := len(curr), len(curr[0])
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if curr[i][j] != nil && curr[i][j].Tag == tag {
				// 找到了当前执子方的棋子
				count := 1
				// 统计连续的相同颜色的棋子数
				for k := 1; k < 5; k++ {
					x, y := i+k*dx, j+k*dy
					if x < 0 || x >= rows || y < 0 || y >= cols {
						break
					}
					if curr[x][y] == nil || curr[x][y].Tag != tag {
						break
					}
					count++
				}
				// 根据连续相同颜色的棋子数计算得分
				if s, ok := countToScore[count]; ok {
					score += s
				}
			}
		}
	}
	return score
}
