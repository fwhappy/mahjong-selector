package mselector

import (
	"sort"

	"github.com/fwhappy/util"
)

// 是否是定缺的牌
func isLack(tile, lack int) bool {
	if lack == 0 {
		return false
	}
	return tile/10 == lack/10
}

func checkHu(s []int) bool {
	var tmpCard = util.SliceCopy(s)
	// 升序排列
	sort.Ints(tmpCard)
	// 找到所有的对
	var pos = findDuiPos(tmpCard)
	// 找不到对，无法胡牌
	if len(pos) == 0 {
		return false
	}
	// 7对(目前版本只有手中为7个对才可以胡)
	if len(pos) == 7 {
		return true
	}
	// 遍历所有对，因为胡必须有对
	for _, v := range pos {
		var tmp = removeDui(tmpCard, v)
		if allIsKeOrShun(tmp) {
			return true
		}
	}
	return false
}

func removeDui(card []int, duiPos int) []int {
	var tmp = make([]int, 0)
	tmp = append(tmp, card[:duiPos]...)
	tmp = append(tmp, card[duiPos+2:]...)
	return tmp
}

// 算出所有可听的列表
// 到这里就认为手牌已经无缺了
func getTingMap(s []int) map[int][]int {
	tingMap := make(map[int][]int)

	// 可能会听的牌
	maybeTing := getMaybeTing(slice2map(s))
	for i, tile := range s {
		// 跳过重复的牌
		if i != 0 && s[i] == s[i-1] {
			continue
		}
		tmp := util.SliceCopy(s)
		// 能听的牌
		tingTiles := []int{}

		for mayTile := range maybeTing {
			if tmp[i] = mayTile; checkHu(tmp) {
				tingTiles = append(tingTiles, mayTile)
			}
		}
		if len(tingTiles) > 0 {
			tingMap[tile] = tingTiles
		}
	}

	return tingMap
}

// 可能听的牌只能是几种情况：
// 1、手牌
// 2、手牌的临牌
func getMaybeTing(m map[int]int) map[int]int {
	var maybeTing = map[int]int{}
	for i := range m {
		switch {
		case i == MAHJONG_CRAK1 || i == MAHJONG_BAM1 || i == MAHJONG_DOT1:
			maybeTing[i] = 1
			maybeTing[i+1] = 1
		case i == MAHJONG_CRAK9 || i == MAHJONG_BAM9 || i == MAHJONG_DOT9:
			maybeTing[i-1] = 1
			maybeTing[i] = 1
		default:
			maybeTing[i-1] = 1
			maybeTing[i] = 1
			maybeTing[i+1] = 1
		}
	}
	return maybeTing
}

// 出去一个对后，所有的牌必须是刻或者顺
func allIsKeOrShun(card []int) bool {
	// 不断从头到尾遍历，先找刻字，在找顺子(如果是刻字一定是刻字，不需要在组合为顺子)
	var count = len(card)
	for i := 0; i < count/3; i++ {
		find := findAndRemoveKeOrShun(&card)
		if !find {
			return false
		}
	}
	return len(card) == 0
}

func findAndRemoveKeOrShun(card *[]int) bool {
	var find = findAndRemoveKe(card)
	if find {
		return true
	}
	return findAndRemoveShun(card)
}

func findAndRemoveKe(card *[]int) bool {
	var v = *card
	if v[0] == v[1] && v[1] == v[2] {
		var tmp = make([]int, 0)
		tmp = append(tmp, v[3:]...)
		*card = tmp
		return true
	}
	return false
}

func findAndRemoveShun(card *[]int) bool {
	var v = *card
	var tmp = make([]int, 0)
	for i := 1; i < len(v); i++ {
		switch {
		case v[i] == v[i-1]:
			tmp = append(tmp, v[i])
		case v[i] == v[i-1]+1:
			if v[i]-v[0] == 2 {
				tmp = append(tmp, v[i+1:]...)
				*card = tmp
				return true
			}
		default:
			return false
		}
	}
	return false
}

func findDuiPos(card []int) []int {
	var pos = []int{}
	for i := 0; i < len(card)-1; i++ {
		if card[i] == card[i+1] {
			pos = append(pos, i)
			i++
		}
	}
	return pos
}

// 获取牌的权重列表
func getTilesWeight(tiles []int, specified []int) map[int]int {
	m := slice2map(tiles)
	tilesWeight := map[int]int{}

	// 统计出出牌的优先级
	for tile, cnt := range m {
		// 指定只统计某些牌的权重
		if specified != nil && !util.IntInSlice(tile, specified) {
			continue
		}
		score := 0

		// 自身的分
		if util.IntInSlice(tile, []int{1, 9, 11, 19, 21, 29}) {
			// 最边张不加分
		} else if util.IntInSlice(tile, []int{2, 8, 12, 18, 22, 28}) {
			// 靠近边张的加10分
			score += 10
		} else {
			// 中间的牌+20分
			score += 20
		}
		// 如果有上下张，每张40分
		if m[tile-1] > 0 {
			score += 40
		}
		if m[tile+1] > 0 {
			score += 40
		}
		// 如果有隔张，再加30分
		if !util.IntInSlice(tile, []int{1, 11, 21}) && m[tile-2] > 0 {
			score += 30
		}
		if !util.IntInSlice(tile, []int{9, 19, 29}) && m[tile+2] > 0 {
			score += 30
		}

		// 看同样的牌的张数，如果同样的牌比较多，每张20分
		if cnt >= 2 {
			score += 50 * (cnt - 1)
		}
		tilesWeight[tile] = score
		// showDebug("牌:%v,权重:%v", tile, score)
	}
	return tilesWeight
}

// 获取牌的权重列表
func getTilesWeightSum(tiles []int, specified []int) int {
	tilesWeight := getTilesWeight(tiles, specified)
	weight := 0
	for _, score := range tilesWeight {
		weight += score
	}
	return weight
}
