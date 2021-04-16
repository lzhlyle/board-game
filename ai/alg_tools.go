package ai

import (
	"board-game/core"
)

// SpinSquare90 旋转 90 度
func SpinSquare90(mat [][]*core.PlaySignal) ([][]*core.PlaySignal, error) {
	if len(mat) == 0 {
		return mat, nil
	}
	if len(mat) != len(mat[0]) {
		return nil, ErrNotSquare
	}

	res := core.CloneMatrix(mat)

	for i, row := range res {
		for j := range row {
			if i < j {
				res[i][j], res[j][i] = res[j][i], res[i][j]
			}
		}
		for l, r := 0, len(row)-1; l < r; {
			row[l], row[r] = row[r], row[l]
			l++
			r--
		}
	}
	return res, nil
}

// FlipLR 左右翻转
func FlipLR(mat [][]*core.PlaySignal) [][]*core.PlaySignal {
	res := core.CloneMatrix(mat)
	for _, row := range res {
		for l, r := 0, len(row)-1; l < r; {
			row[l], row[r] = row[r], row[l]
			l++
			r--
		}
	}
	return res
}
