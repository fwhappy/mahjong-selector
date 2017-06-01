package mselector

import (
	"sort"

	"github.com/fwhappy/util"
)

// MSelector 麻将选牌器
type MSelector struct {
	aiLevel      int         // AI等级
	tiles        map[int]int // 所有牌
	handTiles    map[int]int // 手牌
	discardTiles map[int]int // 明牌
	showTiles    map[int]int // 弃牌
	remainTiles  map[int]int // 剩余牌
	lack         int         // 缺的牌
}

// NewMSelector 生成一个选牌器
func NewMSelector() *MSelector {
	selector := &MSelector{}
	selector.tiles = make(map[int]int, 0)
	selector.clean()
	return selector
}

// 清空选牌器
func (ms *MSelector) clean() {
	ms.handTiles = make(map[int]int, 0)
	ms.discardTiles = make(map[int]int, 0)
	ms.showTiles = make(map[int]int, 0)
	ms.remainTiles = make(map[int]int, 0)
}

// SetAILevel 设置AI级别
func (ms *MSelector) SetAILevel(level int) {
	ms.aiLevel = level
}

// SetLack 设置定缺
func (ms *MSelector) SetLack(tile int) {
	ms.lack = tile
}

// GetLack 设置定缺
func (ms *MSelector) GetLack() int {
	return ms.lack
}

// GetTiles 设置所有的牌
func (ms *MSelector) GetTiles() map[int]int {
	return ms.tiles
}

// GetShuffleTiles 获取打乱排序的牌
func (ms *MSelector) GetShuffleTiles() []int {
	return util.ShuffleSliceInt(map2slice(ms.tiles))
}

// SetTiles 设置所有的牌
func (ms *MSelector) SetTiles(s []int) {
	ms.tiles = slice2map(s)
}

// AddHandTilesMap 添加手牌
func (ms *MSelector) AddHandTilesMap(m map[int]int) {
	ms.handTiles = mergeMap(ms.handTiles, m)
}

// AddHandTilesSlice 添加手牌
func (ms *MSelector) AddHandTilesSlice(s []int) {
	ms.handTiles = mergeMap(ms.handTiles, slice2map(s))
}

// SetHandTilesSlice 设置手牌
func (ms *MSelector) SetHandTilesSlice(s []int) {
	ms.handTiles = slice2map(s)
}

// SetHandTilesMap 设置手牌
func (ms *MSelector) SetHandTilesMap(m map[int]int) {
	ms.handTiles = m
}

// ShowHandTiles 显示手牌
func (ms *MSelector) ShowHandTiles() []int {
	s := map2slice(ms.handTiles)
	sort.Ints(s)
	return s
}

// AddShowTilesSlice 添加明牌
func (ms *MSelector) AddShowTilesSlice(s []int) {
	ms.showTiles = mergeMap(ms.showTiles, slice2map(s))
}

// AddShowTilesMap 添加明牌
func (ms *MSelector) AddShowTilesMap(m map[int]int) {
	ms.showTiles = mergeMap(ms.showTiles, m)
}

// SetShowTilesSlice 添加明牌
func (ms *MSelector) SetShowTilesSlice(s []int) {
	ms.showTiles = slice2map(s)
}

// SetShowTilesMap 添加明牌
func (ms *MSelector) SetShowTilesMap(m map[int]int) {
	ms.showTiles = m
}

// ShowShowTiles 显示明牌
func (ms *MSelector) ShowShowTiles() []int {
	s := map2slice(ms.showTiles)
	sort.Ints(s)
	return s
}

// AddDiscardTilesSlice 添加弃牌
func (ms *MSelector) AddDiscardTilesSlice(s []int) {
	ms.discardTiles = mergeMap(ms.discardTiles, slice2map(s))
}

// AddDiscardTilesMap 添加弃牌
func (ms *MSelector) AddDiscardTilesMap(m map[int]int) {
	ms.discardTiles = mergeMap(ms.discardTiles, m)
}

// SetDiscardTilesSlice 设置弃牌
func (ms *MSelector) SetDiscardTilesSlice(s []int) {
	ms.discardTiles = slice2map(s)
}

// SetDiscardTilesMap 设置弃牌
func (ms *MSelector) SetDiscardTilesMap(m map[int]int) {
	ms.discardTiles = m
}

// ShowDiscardTiles 显示明牌
func (ms *MSelector) ShowDiscardTiles() []int {
	s := map2slice(ms.discardTiles)
	sort.Ints(s)
	return s
}

// SetRemainTilesSlice 设置剩余的牌
func (ms *MSelector) SetRemainTilesSlice(s []int) {
	ms.remainTiles = slice2map(s)
}

// SetRemainTilesMap 添加剩余的牌
func (ms *MSelector) SetRemainTilesMap(m map[int]int) {
	ms.remainTiles = m
}

// DeductRemainTiles 扣除剩余的牌
func (ms *MSelector) DeductRemainTiles(tiles ...int) {
	for _, tile := range tiles {
		cnt, exists := ms.remainTiles[tile]
		if !exists {
			showError("扣除剩余牌的数量错误, tile:%v", tile)
		} else if cnt == 1 {
			delete(ms.remainTiles, tile)
		} else {
			ms.remainTiles[tile]--
		}
	}
}

// ShowRemainTiles 显示剩余的牌
func (ms *MSelector) ShowRemainTiles() []int {
	s := map2slice(ms.remainTiles)
	sort.Ints(s)
	return s
}

// CalcRemaimTiles 计算剩余的牌
func (ms *MSelector) CalcRemaimTiles() {
	ms.remainTiles = make(map[int]int, 0)
	for tile, cnt := range ms.tiles {
		cnt -= ms.handTiles[tile]
		cnt -= ms.showTiles[tile]
		cnt -= ms.discardTiles[tile]

		if cnt > 0 {
			ms.remainTiles[tile] = cnt
		} else if cnt < 0 {
			showError("计算剩余牌的数量错误, tile:%v", tile)
		}
	}
}

// 读取给予的牌的剩余张数之和
func (ms *MSelector) getRemainTilesCnt(tiles []int) int {
	cnt := 0
	for _, tile := range util.SliceUniqueInt(tiles) {
		cnt += ms.remainTiles[tile]
	}
	return cnt
}

func (ms *MSelector) hasTile(tile int) bool {
	_, exists := ms.tiles[tile]
	return exists
}

// 获取跟某张牌有关系的牌
func (ms *MSelector) getRelationTiles(tiles ...int) []int {
	relationTiles := []int{}

	for _, tile := range tiles {
		relationTiles = append(relationTiles, tile)
		if ms.hasTile(tile - 1) {
			relationTiles = append(relationTiles, tile-1)
		}
		if ms.hasTile(tile + 1) {
			relationTiles = append(relationTiles, tile+1)
		}
		if ms.hasTile(tile - 2) {
			relationTiles = append(relationTiles, tile-2)
		}
		if ms.hasTile(tile + 2) {
			relationTiles = append(relationTiles, tile+2)
		}
	}
	return util.SliceUniqueInt(relationTiles)
}
