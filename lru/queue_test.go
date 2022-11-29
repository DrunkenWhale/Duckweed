package lru

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestQueue(t *testing.T) {
	queueNumber := 17
	queue := NewQueue(queueNumber)
	for i := 0; i < queueNumber; i++ {
		queue.Push(i)
		assert.Equal(t, queue.Size(), i+1)
	}
	assert.Equal(t, queue.IsOverflow(), true)
	for i := 0; i < queueNumber; i++ {
		queue.Contains(queueNumber - i)
		assert.Equal(t, i, queue.Pop())
		assert.Equal(t, queue.Size(), queueNumber-i-1)
	}
}
