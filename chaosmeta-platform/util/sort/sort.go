package sort

import "sort"

func RemoveDuplicates(nums []int) []int {
	sort.Ints(nums)
	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}
	return nums[:i+1]
}
