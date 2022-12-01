package index

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestNumLessThan(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	assert.Equal(t, 0, numLessThanEqual(arr, 0))
	assert.Equal(t, 1, numLessThanEqual(arr, 1))
	assert.Equal(t, 2, numLessThanEqual(arr, 2))
	assert.Equal(t, 3, numLessThanEqual(arr, 3))
	assert.Equal(t, 4, numLessThanEqual(arr, 4))
	assert.Equal(t, 5, numLessThanEqual(arr, 5))
	assert.Equal(t, 5, numLessThanEqual(arr, 6))
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
