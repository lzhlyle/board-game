package ai

import (
	"board-game/core"
)

// Spin90 旋转 90 度
func Spin90(mat [][]*core.PlaySignal) ([][]*core.PlaySignal, error) {
	if len(mat) == 0 {
		return mat, nil
	}
	if len(mat) != len(mat[0]) {
		return nil, ErrNotSquare
	}

	for i, row := range mat {
		for j := range row {
			if i < j {
				mat[i][j], mat[j][i] = mat[j][i], mat[i][j]
			}
		}
		for l, r := 0, len(row)-1; l < r; {
			row[l], row[r] = row[r], row[l]
			l++
			r--
		}
	}
	return mat, nil
}
