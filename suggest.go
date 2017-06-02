package mselector

import (
	"sort"

	"github.com/fwhappy/util"
)

// GetSuggest 获取建议出的牌
// depth 查找深度, -1 表示找出最优解，直至没有牌
func (ms *MSelector) GetSuggest(depth int) int {
	// 排序手牌
	s := map2slice(ms.handTiles)
	sort.Ints(s)

	// 如果有定缺的牌，则推荐缺的牌
	// 因为大的牌在右手边，所以优先推荐，便于用户操作
	if ms.lack > 0 {
		for i := len(s) - 1; i >= 0; i-- {
			if isLack(s[i], ms.lack) {
				showDebug("此牌型有缺")
				return s[i]
			}
		}
	}

	// 如果当前牌型已经叫牌，推荐剩余张数最多的牌
	if tingMap := getTingMap(s); len(tingMap) > 0 {
		showDebug("此牌型已听牌")
		selectTile := 0
		selectTileTingCnt := -1

		for tile, tingTiles := range tingMap {
			tingCnt := ms.getRemainTilesCnt(tingTiles)
			showDebug("打出:%v, 可听:%v, 剩余(张):%v", tile, tingTiles, tingCnt)
			if tingCnt > selectTileTingCnt {
				selectTileTingCnt = tingCnt
				selectTile = tile
			}
		}
		return selectTile
	}

	// 如果当前牌未叫牌，则根据权重选牌
	switch ms.aiLevel {
	case AI_BRASS: // 英勇黄铜
		return ms.suggestByAIBrass(s)
	case AI_SLIVER: // 不屈白银
		return ms.suggestByAISliver(s)
	case AI_GOLD: // 荣耀黄金
		return ms.suggestByAIGold(s)
	case AI_PLATINUM: // 华贵铂金
		return ms.suggestByAIPlatinum(s)
	case AI_DIAMOND: // 璀璨钻石
		return ms.suggestByAIDiamond(s)
	case AI_MASTER: // 非凡大师
		return ms.suggestByAIMaster(s)
	case AI_KING: // 最强王者
		return ms.suggestByAIKing(s)
	default:
		return ms.suggestByAIPlatinum(s)
	}
}

// GetEffects 计算牌型所在的牌阶以及一类有效牌数量
func (ms *MSelector) GetEffects(s []int) (int, int) {
	step, effects := calcEffects(s)
	return step, len(effects)
}

// 英勇黄铜
// 打手牌的第一张
func (ms *MSelector) suggestByAIBrass(s []int) int {
	showDebug("--------------------------------英勇黄铜-------------------------------")
	return s[0]
}

// 不屈白银
// 根据权重来计算出什么牌
// 相同权重的话，计算剩余牌的张数
func (ms *MSelector) suggestByAISliver(s []int) int {
	showDebug("--------------------------------不屈白银-------------------------------")

	// 生成牌的权重列表
	tilesWeight := getTilesWeight(s, nil)

	return ms.selectByWeightAndRemain(tilesWeight)
}

// 不屈白银
// 根据权重来计算出什么牌
// 相同权重的话，计算剩余牌的张数
func (ms *MSelector) suggestByAIGold(s []int) int {
	showDebug("-------------------------------荣耀黄金--------------------------------")

	// 计算当前手牌所处的阶段
	unPlayStep := getTilesStep(s)
	showDebug("当前手牌牌阶:%v", unPlayStep)

	// 最多一类有效牌张数
	maxEffectTileCnt := 0
	// 最多一类有效牌列表
	maxEffectTiles := []int{}
	maxEffectList := map[int][]int{}

	// 循环删除某一张手牌，计算一类有效牌的数量
	// 如果某些牌的有效张
	for _, playTile := range util.SliceUniqueInt(s) {
		tiles := util.SliceDel(s, playTile)
		step, effects := calcEffects(tiles)
		if step > unPlayStep {
			sort.Ints(effects)
			showDebug("打出:%v,手牌:%v, 一类有效牌:%v(%v)(remain:%v)------------", playTile, tiles, effects, len(effects), ms.getRemainTilesCnt(effects))

			effectsLen := len(effects)
			if effectsLen > maxEffectTileCnt {
				maxEffectTileCnt = effectsLen
				maxEffectTiles = []int{playTile}
				maxEffectList = map[int][]int{}
				maxEffectList[playTile] = effects
			} else if len(effects) == maxEffectTileCnt {
				maxEffectTiles = append(maxEffectTiles, playTile)
				maxEffectList[playTile] = effects
			}
		}
	}

	// 如果存在相同的一类有效牌，则根据权重再取一次
	if len(maxEffectTiles) > 1 {
		showDebug("存在多张有效牌相同的打法，根据权重重新筛选一次:%v", maxEffectTiles)
		return maxEffectTiles[0]
		// return ms.selectByWeightAndRemain(s, maxEffectTiles)
	}

	return maxEffectTiles[0]
}

// 华贵铂金
func (ms *MSelector) suggestByAIPlatinum(s []int) int {
	showDebug("-------------------------------华贵铂金--------------------------------")

	// 计算当前手牌所处的阶段
	unPlayStep := getTilesStep(s)
	showDebug("当前手牌牌阶:%v", unPlayStep)

	if unPlayStep < 3 {
		showDebug("当前牌阶低于3级, 直接按照权重选牌")
		tilesWeight := getTilesWeight(s, nil)
		return ms.selectByWeightAndRemain(tilesWeight)
	}

	// 最多一类有效牌张数
	maxEffectTileCnt := 0
	// 最多一类有效牌列表
	maxEffectTiles := []int{}
	maxEffectList := map[int][]int{}
	maxEffectTotalWeights := map[int]int{}

	// 循环删除某一张手牌，计算一类有效牌的数量
	// 如果某些牌的有效张
	for _, playTile := range util.SliceUniqueInt(s) {
		tiles := util.SliceDel(s, playTile)
		step, effects, totalWeight := calcEffectsAndRemainWeight(tiles)
		if step >= unPlayStep {
			sort.Ints(effects)
			showDebug("打出:%v,手牌:%v, 一类有效牌:%v(%v)(remain:%v)------------", playTile, tiles, effects, len(effects), ms.getRemainTilesCnt(effects))

			effectsLen := len(effects)
			if effectsLen > maxEffectTileCnt {
				maxEffectTileCnt = effectsLen
				maxEffectTiles = []int{playTile}
				maxEffectList = map[int][]int{}
				maxEffectList[playTile] = effects
				maxEffectTotalWeights = map[int]int{}
				maxEffectTotalWeights[playTile] = totalWeight
			} else if len(effects) == maxEffectTileCnt {
				maxEffectTiles = append(maxEffectTiles, playTile)
				maxEffectList[playTile] = effects
				maxEffectTotalWeights[playTile] = totalWeight
			}
		}
	}

	// 如果存在相同的一类有效牌，则根据权重再取一次
	if len(maxEffectTiles) > 1 {
		showDebug("存在多张有效牌相同的打法，根据剩余牌权重重新筛选一次:%v", maxEffectTotalWeights)

		// 读取权重最大的牌
		maxWeightTiles, _ := getMaxValueSlice(maxEffectTotalWeights)
		showDebug("权重最大的牌:%v", maxWeightTiles)

		// 找出权重最小的牌中，关联牌最少的一张
		maxRemainCnt := 0
		var maxRemainTile int
		for _, tile := range maxWeightTiles {
			remainCnt := ms.getRemainTilesCnt([]int{tile})
			if maxRemainCnt < remainCnt {
				maxRemainCnt = remainCnt
				maxRemainTile = tile
			}
		}
		return maxRemainTile
	}

	return maxEffectTiles[0]
}

// 璀璨钻石
func (ms *MSelector) suggestByAIDiamond(s []int) int {
	showDebug("-------------------------------璀璨钻石--------------------------------")

	// // 未打出时的阶段
	// unPlayStep := getTilesStep(s)
	// showDebug("未打牌的阶段:%v", unPlayStep)

	// times := 0
	// start := time.Now().UnixNano()

	// for _, playTile := range util.SliceUniqueInt(s) {
	// 	tiles := util.SliceDel(s, playTile)
	// 	step, effects := calcEffects(tiles)
	// 	if step >= unPlayStep {
	// 		sort.Ints(effects)
	// 		showDebug("打出:%v,手牌:%v, 一类有效牌:%v(%v)(remain:%v)------------", playTile, tiles, effects, len(effects), ms.getRemainTilesCnt(effects))

	// 		// 计算二类有效牌

	// 		// 取出所有有关系的牌
	// 		relationTiles := ms.getRelationTiles(tiles...)
	// 		// showDebug("可能是二类的有效牌:%v", relationTiles)
	// 		secondEffects := []int{}

	// 		for _, tile := range relationTiles {
	// 			// 跳过一类有效牌
	// 			if util.IntInSlice(tile, effects) {
	// 				continue
	// 			}

	// 			for k, replacedTile := range tiles {
	// 				times++
	// 				replacedTiles := append(append([]int{tile}, tiles[:k]...), tiles[k+1:]...)
	// 				curStep, curEffects := calcEffects(replacedTiles)
	// 				if curStep == step && len(curEffects) > len(effects) {
	// 					if !util.IntInSlice(tile, secondEffects) {
	// 						showDebug("%v替换成%v, 一类有效牌变为:%v(%v)(remain:%v)", replacedTile, tile, curEffects, len(curEffects), ms.getRemainTilesCnt(curEffects))
	// 						secondEffects = append(secondEffects, tile)
	// 					}
	// 				}
	// 			}
	// 		}

	// 		showDebug("二类有效牌:%v(%v)(remain:%v)", secondEffects, len(secondEffects), ms.getRemainTilesCnt(secondEffects))
	// 	}
	// 	// break
	// }

	// showDebug("总计算次数:%v, 耗时:%v", times, (time.Now().UnixNano() - start))

	// return 0
	return ms.suggestByAIPlatinum(s)
}

// 非凡大师
func (ms *MSelector) suggestByAIMaster(s []int) int {
	return ms.suggestByAIPlatinum(s)
}

// 最强王者
func (ms *MSelector) suggestByAIKing(s []int) int {
	return ms.suggestByAIPlatinum(s)
}

// 根据权重筛选
func (ms *MSelector) selectByWeightAndRemain(tilesWeight map[int]int) int {
	// 读取权重最小的牌
	minWeightTiles, _ := getMinValueSlice(tilesWeight)
	showDebug("权重最小的牌:%v", minWeightTiles)

	// 找出权重最小的牌中，关联牌最少的一张
	minRelationCnt := 1000000
	minRelationTile := 0
	for _, tile := range minWeightTiles {
		relationCnt := ms.getRemainTilesCnt(ms.getRelationTiles(tile))
		if relationCnt < minRelationCnt {
			minRelationCnt = relationCnt
			minRelationTile = tile
		}
	}
	return minRelationTile
}

// 根据剩余张数筛选
func (ms *MSelector) selectByRemainCnt(s []int) int {
	remainCnt := 0
	var selectedTile int

	for _, tile := range s {
		if cnt := ms.getRemainTilesCnt([]int{tile}); cnt >= remainCnt {
			remainCnt = cnt
			selectedTile = tile
		}
	}

	return selectedTile
}
