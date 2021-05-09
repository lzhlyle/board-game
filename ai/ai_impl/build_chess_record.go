package ai_impl

import (
	"board-game/ai"
	"board-game/core"
)

type BuildChessRecord struct {
	chessRecord ai.IChessRecord
	zip         ai.IZip
	boardGame   core.BoardGame

	// (preZip, nextRates)
	Zip2NextRates map[int32][]*ai.NextRates
}

func NewBuildChessRecord(chessRecord ai.IChessRecord, zip ai.IZip, boardGame core.BoardGame) *BuildChessRecord {
	res := &BuildChessRecord{chessRecord: chessRecord, zip: zip, boardGame: boardGame, Zip2NextRates: map[int32][]*ai.NextRates{}}
	res.buildChessRecord()
	return res
}

func (bcr *BuildChessRecord) buildChessRecord() {
	// later @lzh 可考虑多协程并发处理，注意将 Zip2NextRates 改为支持并发的数据结构
	bcr.dfs(core.NewEmptyMoveSnapshot(bcr.boardGame.Board().Width, bcr.boardGame.Board().Height), bcr.Zip2NextRates)
}

func (bcr *BuildChessRecord) dfs(aSnapshot *core.MoveSnapshot, zip2NextRates map[int32][]*ai.NextRates) *ai.NextRates {
	aZip := bcr.zip.Zip(aSnapshot.Board).(int32)

	// terminator
	if end, winner := bcr.boardGame.GameEnd(aSnapshot); end {
		zip2NextRates[aZip] = []*ai.NextRates{}
		if winner == nil {
			return ai.NewNextRates(aZip, [3]int{0, 100, 0})
		}
		return ai.NewNextRates(aZip, [3]int{100, 0, 0})
	}

	// look up possibles
	b := bcr.boardGame.NextPlayer(aSnapshot.Player)
	bPossibles := make([][2]int, 0) // [2]int: {i, j}
	for i, row := range aSnapshot.Board {
		for j := 0; j < len(row); j++ {
			if row[j] == nil {
				bPossibles = append(bPossibles, [2]int{i, j})
			}
		}
	}

	allBNextRates := make([]*ai.NextRates, len(bPossibles))
	// travel possibles
	var aRates = [3]int{}
	for i, pos := range bPossibles {
		bSnapshot := core.GenSnapshot(aSnapshot.Step+1, pos[0], pos[1], b, aSnapshot)
		bNextRates := bcr.dfs(bSnapshot, zip2NextRates)
		allBNextRates[i] = bNextRates

		// cross accumulate
		// 零和游戏，交叉累计
		bRates := bNextRates.Rates
		aRates[0] += bRates[2]
		aRates[1] += bRates[1]
		aRates[2] += bRates[0]
	}
	// average
	for i := 0; i < 3; i++ {
		aRates[i] = aRates[i] / len(bPossibles)
	}

	zip2NextRates[aZip] = bcr.sort(bcr.filter(allBNextRates))
	return ai.NewNextRates(aZip, aRates)
}

func (bcr *BuildChessRecord) filter(rates []*ai.NextRates) []*ai.NextRates {
	// 能赢则赢
	res := make([]*ai.NextRates, 0, len(rates))
	for _, rate := range rates {
		if rate.Rates[0] == 100 {
			// 多几种选择，而非直接返回
			res = append(res, rate)
		}
	}
	if len(res) > 0 {
		return res
	}

	// 对方再下一步不会赢，则才可走
	for _, rate := range rates {
		nextWillWin := false
		for _, next := range bcr.Zip2NextRates[rate.NextZip.(int32)] {
			if next.Rates[0] == 100 {
				nextWillWin = true
				break
			}
		}
		if !nextWillWin {
			res = append(res, rate)
		}
	}
	if len(res) == 0 {
		// 必输，放弃治疗
		return rates
	}

	return res
}

func (bcr *BuildChessRecord) sort(rates []*ai.NextRates) []*ai.NextRates {
	if len(rates) == 0 {
		return rates
	}

	if rates[0].Rates[0] == 100 {
		return rates // 能赢，无需排序
	}

	bcr.chessRecord.SortRecords(rates)

	return rates
}
