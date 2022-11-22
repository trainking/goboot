package random

import (
	"math/rand"
	"sort"
)

type (
	// Choice  选择项的结构
	Choice struct {
		Item   interface{} // Item 元素
		Weight uint        // 权重值
	}

	// WeightRand 加权随机
	WeightRand struct {
		data   []Choice
		totals []int
		max    int
	}
)

// NewWeightRand 创建一个加权随机算法。
// 该算法使用将权重放入一个有序数组，然后从随机数，二分法去查找有效区间，得出结果
func NewWeightRand(data []Choice) WeightRand {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Weight < data[j].Weight
	})

	totals := make([]int, len(data))
	runningTotal := 0

	for i, c := range data {
		runningTotal += int(c.Weight)
		totals[i] = runningTotal
	}

	return WeightRand{data: data, totals: totals, max: runningTotal}
}

// Pick 选出随机数，此处利用二分查找的特点
func (wr WeightRand) Pick() interface{} {
	r := rand.Intn(wr.max) + 1
	i := sort.SearchInts(wr.totals, r)
	return wr.data[i].Item
}
