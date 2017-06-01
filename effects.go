package mselector

import (
	"sort"

	"github.com/fwhappy/util"
)

// 有效牌算法
func calcEffects(s []int) (int, []int) {
	var effects = []int{}

	step := getTilesStep(s)

	// 可能是一类有效的牌
	maybeFirstEffects := getMaybeTing(slice2map(s))
	maybeFirstEffectsSlice := map2slice(maybeFirstEffects)
	sort.Ints(maybeFirstEffectsSlice)
	for tile := range maybeFirstEffects {
		fillingStep := getTilesStep(append(util.SliceCopy(s), tile))
		if fillingStep > step {
			effects = append(effects, tile)
		}
	}

	return step, effects
}

// 计算牌型所属步骤
// 暂未考虑7对
func getTilesStep(s []int) int {
	var step int
	var tmpCard = util.SliceCopy(s)
	// 升序排列
	sort.Ints(tmpCard)

	// 找到所有的对
	var pos = findDuiPos(tmpCard)
	for _, v := range pos {
		var tmp = removeDui(tmpCard, v)
		cnt, _ := findKeOrShunCnt(tmp)
		if cnt >= step {
			step = cnt + 1
		}
	}
	// 找出无对的步骤
	if cnt, _ := findKeOrShunCnt(s); cnt > step {
		step = cnt
	}

	return step
}

// 有效牌算法
func calcEffectsAndRemainWeight(s []int) (int, []int, int) {
	var effects = []int{}

	step := getTilesStep(s)
	totalRemainWeight := 0

	// 可能是一类有效的牌
	maybeFirstEffects := getMaybeTing(slice2map(s))
	maybeFirstEffectsSlice := map2slice(maybeFirstEffects)
	sort.Ints(maybeFirstEffectsSlice)
	for tile := range maybeFirstEffects {
		fillingStep, remainWeight := getTilesStepAndRemainWeight(append(util.SliceCopy(s), tile))
		if fillingStep > step {
			effects = append(effects, tile)
			totalRemainWeight += remainWeight
		}
	}

	return step, effects, totalRemainWeight
}

// 计算牌型所属步骤
// 暂未考虑7对
func getTilesStepAndRemainWeight(s []int) (int, int) {
	var step int
	var remainWeight int
	var tmpCard = util.SliceCopy(s)
	// 升序排列
	sort.Ints(tmpCard)

	// 找到所有的对
	var pos = findDuiPos(tmpCard)
	for _, v := range pos {
		var tmp = removeDui(tmpCard, v)
		cnt, remainTiles := findKeOrShunCnt(tmp)
		if cnt >= step {
			step = cnt + 1
			remainWeight = getTilesWeightSum(remainTiles, nil)
		}
	}
	// 找出无对的步骤
	if cnt, remainTiles := findKeOrShunCnt(s); cnt > step {
		step = cnt
		remainWeight = getTilesWeightSum(remainTiles, nil)
	}

	return step, remainWeight
}

// 计算牌型分数
// 每张牌一分
func getTilesScore(s []int) int {
	var score7Dui, scoreWithDui, scoreWithoutDui int
	var tmpCard = util.SliceCopy(s)
	// 升序排列
	sort.Ints(tmpCard)

	// 找到所有的对
	var pos = findDuiPos(tmpCard)
	score7Dui = len(pos) * 2
	showDebug("牌型7对分:%v", score7Dui)

	for _, v := range pos {
		var tmp = removeDui(tmpCard, v)
		score, remainTiles := getTilesScoreWithoutDui(tmp)
		showDebug("牌型带对积分:%v,对:%v,剩:%v", score+2, tmpCard[v], remainTiles)
		if score+2 > scoreWithDui {
			scoreWithDui = score + 2
		}
	}

	// 不取对的有效分
	scoreWithoutDui, remainTiles := getTilesScoreWithoutDui(tmpCard)
	showDebug("牌型不带对积分:%v,剩:%v", scoreWithoutDui, remainTiles)

	return util.SliceMaxInt([]int{score7Dui, scoreWithDui, scoreWithoutDui})
}

// 计算牌型分数：不带对子
func getTilesScoreWithoutDui(s []int) (int, []int) {
	// 得分
	score := 0
	// 剩余的牌
	remainTiles := []int{}
	// 不断从头到尾遍历，先找刻字，在找顺子(如果是刻字一定是刻字，不需要在组合为顺子)
	for {
		if len(s) <= 2 {
			remainTiles = append(remainTiles, s...)
			break
		}
		if firstIsKe(s) || firstIsShun(s) {
			s = s[3:]
			score += 3
		} else {
			remainTiles = append(remainTiles, s[0])
			s = s[1:]
		}
	}

	return score, remainTiles
}

// 抽排算法，计算刻或者顺的个数
func findKeOrShunCnt(s []int) (int, []int) {
	// showDebug("s:%v", s)
	var cnt = 0
	remain := []int{}
	for {
		if len(s) <= 2 {
			remain = append(remain, s...)
			break
		}
		sort.Ints(s)
		if firstIsKe(s) {
			// showDebug("找到刻:%v", s[:3])
			s = s[3:]
			cnt++
		} else if firstIsShun(s) {
			// showDebug("找到顺:%v,%v,%v", s[0], s[0]+1, s[0]+2)
			s = util.SliceDel(s, s[0], s[0]+1, s[0]+2)
			cnt++
		} else {
			remain = append(remain, s[0])
			s = s[1:]

		}
	}
	// showDebug("findKeOrShunCnt:%v", cnt)

	return cnt, remain
}

func firstIsKe(s []int) bool {
	return s[0] == s[1] && s[1] == s[2]
}
func firstIsShun(s []int) bool {
	return util.IntInSlice(s[0]+1, s) && util.IntInSlice(s[0]+2, s)
}
