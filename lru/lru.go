package lru

type LRU struct {
	queue *Queue
}

func NewLRU(size int) *LRU {
	return &LRU{
		queue: NewQueue(size),
	}
}

// Push 放入元素
// 如果元素已经存在
// 则移动到lru列表的头部
func (lru *LRU) Push(e int) {
	if lru.queue.Contains(e) {
		lru.queue.Remove(e)
		lru.queue.Push(e)
	} else {
		lru.queue.Push(e)
	}
}

func (lru *LRU) Pop() int {
	return lru.queue.Pop()
}

func (lru *LRU) Contains(e int) bool {
	return lru.queue.Contains(e)
}

func (lru *LRU) IsOverflow() bool {
	return lru.queue.IsOverflow()
}
