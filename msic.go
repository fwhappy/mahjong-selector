package mselector

func map2slice(m map[int]int) []int {
	tiles := []int{}
	for tile, cnt := range m {
		for i := 0; i < cnt; i++ {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

func slice2map(s []int) map[int]int {
	m := make(map[int]int)
	for _, tile := range s {
		if _, exists := m[tile]; exists {
			m[tile]++
		} else {
			m[tile] = 1
		}
	}
	return m
}

func mergeMap(m1, m2 map[int]int) map[int]int {
	for tile, cnt := range m2 {
		if _, exists := m1[tile]; exists {
			m1[tile] += cnt
		} else {
			m1[tile] = cnt
		}
	}
	return m1
}

func getMinValueSlice(m map[int]int) ([]int, int) {
	// value到keys的对应关系
	valueKeys := map[int][]int{}
	// 最小值
	minValue := -1

	for k, v := range m {
		keys, exists := valueKeys[v]
		if !exists {
			keys = []int{k}
		} else {
			keys = append(keys, k)
		}
		valueKeys[v] = keys

		if minValue == -1 || v <= minValue {
			minValue = v
		}
	}

	return valueKeys[minValue], minValue
}

func getMaxValueSlice(m map[int]int) ([]int, int) {
	// value到keys的对应关系
	valueKeys := map[int][]int{}
	// 最大值
	maxValue := -1

	for k, v := range m {
		keys, exists := valueKeys[v]
		if !exists {
			keys = []int{k}
		} else {
			keys = append(keys, k)
		}
		valueKeys[v] = keys

		if v >= maxValue {
			maxValue = v
		}
	}

	return valueKeys[maxValue], maxValue
}
