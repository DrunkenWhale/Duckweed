package index

func insertSliceWithIndex[T any](arr []T, index int, value T) []T {
	if index > len(arr) {
		panic("You, over step")
	}
	return append(arr[:index], append([]T{value}, arr[index:]...)...)
}

func numLessThanEqual(arr []int, num int) int {
	if arr[0] > num {
		return 0
	}
	if arr[len(arr)-1] < num {
		return len(arr)
	}
	return upperBoundSearch(arr, num)
}

func upperBoundSearch(arr []int, num int) int {
	i := 0
	j := len(arr) - 1
	for i <= j {
		mid := (i + j) >> 1
		if arr[mid] > num {
			j = mid - 1
		} else {
			i = mid + 1
		}
	}
	return i
}
