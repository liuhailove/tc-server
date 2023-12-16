package utils

import "sort"

// MedianFloat32 获取 float32 数组的中值
func MedianFloat32(input []float32) float32 {
	num := len(input)
	if num == 0 {
		return 0
	} else if num == 1 {
		return input[0]
	}
	sort.Slice(input, func(i, j int) bool {
		return input[i] < input[j]
	})

	if num%2 != 0 {
		return input[num/2]
	}
	left := input[num/2-1]
	right := input[num/2]
	return (left + right) / 2
}
