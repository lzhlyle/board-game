package core

import (
	"testing"
)

func Test_clone(t *testing.T) {
	ori := [][]PlaySignal{
		{
			PlaySignal{Tag: "x"},
			PlaySignal{Tag: "z"},
			PlaySignal{Tag: "o"},
			PlaySignal{Tag: "z"},
			PlaySignal{Tag: "o"},
			PlaySignal{Tag: "x"},
		},
	}

	res := (&MoveSnapshot{Board: ori}).cloneBoard()
	res[0][1] = PlaySignal{Tag: "r"}

	if ori[0][1] == res[0][1] {
		t.Log("未完全拷贝，会影响")
	} else {
		t.Log("已完全拷贝，不会影响")
	}
}
