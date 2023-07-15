package utils

func Max(nums ...int) int {
	var maxNum = 0
	for _, num := range nums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}
