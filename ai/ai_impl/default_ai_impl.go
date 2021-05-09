package ai_impl

import (
	"board-game/ai"
	"board-game/core"
	"github.com/jroimartin/gocui"
	"math/rand"
	"time"
)

type DefaultAIImpl struct {
	players []*core.Player

	allAI        bool
	strategies   map[ai.Strategy]ai.CalcFunc
	currStrategy ai.Strategy
}

const AIStrategyRandom ai.Strategy = 1 + iota

func NewDefaultAIImpl(players []*core.Player) *DefaultAIImpl {
	res := &DefaultAIImpl{
		players:      players,
		allAI:        true,
		currStrategy: AIStrategyRandom,
	}

	for _, p := range res.players {
		res.allAI = res.allAI && p.AI
	}

	res.strategies = map[ai.Strategy]ai.CalcFunc{
		AIStrategyRandom: res.randomStrategy,
	}

	return res
}

func (t *DefaultAIImpl) RegisterStrategy(key ai.Strategy, strategy ai.CalcFunc) bool {
	if _, ok := t.strategies[key]; ok {
		return false
	}
	t.strategies[key] = strategy
	return true
}

func (t *DefaultAIImpl) SetCurrentStrategy(key ai.Strategy) bool {
	if _, ok := t.strategies[key]; !ok {
		return false
	}
	t.currStrategy = key
	return true
}

func (t *DefaultAIImpl) randomStrategy(curr [][]*core.PlaySignal) (i, j int, err error) {
	optionals := make([][2]int, 0)
	for i := range curr {
		for j := range curr[i] {
			if curr[i][j] == nil {
				optionals = append(optionals, [2]int{i, j})
			}
		}
	}
	if len(optionals) == 0 {
		return -1, -1, ai.ErrCannotMove
	}
	res := optionals[rand.Intn(len(optionals))]
	return res[0], res[1], nil
}

func (t *DefaultAIImpl) AIMove(snapshot *core.MoveSnapshot, moveFn core.MoveFn) core.UpdateFn {
	return func(g *gocui.Gui) error {
		i, j, err := t.Calculate(snapshot.Board)
		if err == ai.ErrCannotMove {
			return nil
		} else if err != nil {
			return err
		}

		v, err := g.View(core.CellToName(i, j))
		if err != nil {
			return err
		}

		if !t.allAI {
			time.Sleep(100 * time.Millisecond)
		}

		return moveFn(g, v)
	}
}

func (t *DefaultAIImpl) Calculate(curr [][]*core.PlaySignal) (i, j int, err error) {
	return t.strategies[t.currStrategy](curr)
}
