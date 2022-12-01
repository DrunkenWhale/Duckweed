package index

import (
	"github.com/go-playground/assert/v2"
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
