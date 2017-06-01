package main

import (
	"fmt"
	mselector "mahjong-selector"
)

func main() {
	// 初始化一个选牌器
	ms := mselector.NewMSelector()
	ms.SetTiles(mselector.MahjongGYTile36)
	tiles := ms.GetShuffleTiles()
	// 保护未使用变量，防止编译错误
	tiles[0] = tiles[0]

	// 预设参数
	lack := 0
	// handTiles := []int{4, 5, 5, 11, 12, 13, 14, 15, 16, 17, 17, 18, 26, 26}
	// handTiles := []int{22, 24, 24, 25, 27, 29, 11, 13, 14, 14, 15, 16, 17, 17}
	handTiles := tiles[:14] // 5 8 11 14
	// 明牌
	showTiles := []int{}
	// 弃牌
	discardTiles := []int{}
	// discardTiles := tiles[5:30]
	// remainTiles := []int{9, 8, 7}
	// remainTiles := tiles[5:30]
	aiLevel := mselector.AI_PLATINUM
	remainTiles := []int{}

	// 设置参数
	ms.SetAILevel(aiLevel)
	ms.SetLack(lack)
	ms.SetHandTilesSlice(handTiles)
	ms.SetShowTilesSlice(showTiles)
	ms.SetDiscardTilesSlice(discardTiles)
	if len(remainTiles) > 0 {
		ms.SetRemainTilesSlice(remainTiles)
	} else {
		ms.CalcRemaimTiles()
	}

	fmt.Println("定缺:", ms.GetLack())
	fmt.Println("手牌:", ms.ShowHandTiles())
	fmt.Println("明牌:", ms.ShowShowTiles())
	fmt.Println("弃牌:", ms.ShowDiscardTiles())
	fmt.Println("剩余的牌:", ms.ShowRemainTiles())

	// 选牌
	tile := ms.GetSuggest(0)

	fmt.Println("选牌结果:", tile)
}
