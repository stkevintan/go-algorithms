package main

import "fmt"

func sort(arr []int) {
	if len(arr) == 0 {
		return
	}
	var (
		left  = 0
		right = len(arr) - 1
	)

	for left < right {
		if arr[left+1] < arr[left] {
			arr[left+1], arr[left] = arr[left], arr[left+1]
			left++
		} else {
			arr[right], arr[left+1] = arr[left+1], arr[right]
			right--
		}
	}

	sort(arr[0:left])
	sort(arr[left+1 : len(arr)])
}

func main() {
	raw := []int{0, 3, 5, 2, 1, 7, 6, 8, 4, 9, 9}
	sort(raw)
	fmt.Println(raw)
}
