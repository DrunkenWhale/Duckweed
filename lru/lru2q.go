package lru

type LRU2Q struct {
	lru   *LRU
	queue *Queue
}

func NewLRU2Q(queSize int, lruSize int) *LRU2Q {
	return &LRU2Q{
		lru:   NewLRU(lruSize),
		queue: NewQueue(queSize),
	}
}

// Push
// @return: 被刷出的值 (如果没有的话 这个值是-1)
func (lru2q *LRU2Q) Push(e int) int {
	if lru2q.queue.Contains(e) {
		// 队列中已存在的数据应当被放置到lru中
		lru2q.queue.Remove(e)
		lru2q.lru.Push(e)
		if lru2q.lru.IsOverflow() {
			return lru2q.lru.Pop()
		}
		return -1
	} else if lru2q.lru.Contains(e) {
		// 存在于LRU中
		// 移到lru链表前头即可
		lru2q.lru.Push(e)
		return -1
	} else {
		lru2q.queue.Push(e)
		if lru2q.queue.IsOverflow() {
			return lru2q.queue.Pop()
		}
		return -1
	}
}
