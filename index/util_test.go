package index

import (
	"github.com/go-playground/assert/v2"
	"math/rand"
	"sort"
	"testing"
)

func TestNumLessThan(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	assert.Equal(t, 0, numLessThan(arr, 0))
	assert.Equal(t, 1, numLessThan(arr, 1))
	assert.Equal(t, 2, numLessThan(arr, 2))
	assert.Equal(t, 3, numLessThan(arr, 3))
	assert.Equal(t, 4, numLessThan(arr, 4))
	assert.Equal(t, 5, numLessThan(arr, 5))
	assert.Equal(t, 5, numLessThan(arr, 6))
}

func TestInsertSliceWithIndex(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	arr = insertSliceWithIndex(arr, 3, 7)
	assert.Equal(t, len(arr), 6)
	assert.Equal(t, arr[3], 7)
}

func TestUpperBoundSearch(t *testing.T) {
	arr := []int{1, 3, 7, 10, 17}
	assert.Equal(t, upperBoundSearch(arr, -1), 0)
	assert.Equal(t, upperBoundSearch(arr, 1), 1)
	assert.Equal(t, upperBoundSearch(arr, 3), 2)
	assert.Equal(t, upperBoundSearch(arr, 15), 4)
	assert.Equal(t, upperBoundSearch(arr, 19), 5)
}

func TestSearch(t *testing.T) {
	arr := rand.Perm(114514)
	sort.Ints(arr)
	for i := 0; i < 10000; i++ {
		num := rand.Int() % 114514
		assert.Equal(t, numLessThan(arr, num), numLessThan1(arr, num))
	}
}
