package lru

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestLRU2QPush(t *testing.T) {
	queSize := 15
	lruSize := 3
	lru2q := NewLRU2Q(queSize, lruSize)
	lru2q.Push(114514)
	out := lru2q.Push(114514)
	assert.Equal(t, -1, out)
	for i := 0; i < 114; i++ {
		lru2q.Push(i)
	}
	assert.Equal(t, 1, lru2q.lru.queue.Size())
	assert.Equal(t, queSize, lru2q.queue.Size())
	assert.Equal(t, false, lru2q.queue.IsOverflow())
	for i := 0; i < 114; i++ {
		lru2q.Push(i)
	}
	assert.Equal(t, 1, lru2q.lru.queue.Size())
	lru2q.Push(110)
	lru2q.Push(111)
	lru2q.Push(113)
	assert.Equal(t, 3, lru2q.lru.queue.Size())
	assert.Equal(t, 110, lru2q.Push(109))
}
