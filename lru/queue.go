package lru

import "container/list"

// Queue 队列
// 只能存放int类型的数据
// 且数据不能重复
type Queue struct {
	size int
	l    *list.List
	m    map[int]*list.Element
}

func NewQueue(size int) *Queue {
	return &Queue{
		size: size,
		l:    (&list.List{}).Init(),
		m:    make(map[int]*list.Element, 0),
	}
}

func (que *Queue) Push(e int) {
	element := que.l.PushFront(e)
	que.m[e] = element

}

func (que *Queue) Pop() int {
	res := que.l.Back()
	que.l.Remove(res)
	delete(que.m, res.Value.(int))
	return res.Value.(int)
}

func (que *Queue) Size() int {
	return que.l.Len()
}

func (que *Queue) IsOverflow() bool {
	return que.Size() > que.size
}

func (que *Queue) Contains(e int) bool {
	_, flag := que.m[e]
	return flag
}

func (que *Queue) Remove(e int) {
	element, flag := que.m[e]
	if !flag {
		// 元素不存在于队列中
		return
	}
	que.l.Remove(element)
	delete(que.m, e)
}
