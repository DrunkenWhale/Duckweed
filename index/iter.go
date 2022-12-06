package index

import "Duckweed/buffer"

func NewIter(leftmostPageID int, bf buffer.BufferPool) *Iter {
	return &Iter{
		nextPageID: leftmostPageID,
		keys:       make([]int, 0),
		records:    make([][]byte, 0),
		index:      -1,
		bf:         bf,
	}
}

// Scan用的迭代器
type Iter struct {

	// 当前迭代器所在的page的rightSibling的ID
	nextPageID int

	// 当前所在的节点的kv
	keys    []int
	records [][]byte

	// 遍历到kv的第几个下标
	index int

	// B+树所在的内存池
	bf buffer.BufferPool
}

// 正常返回kv键值
// 但是在没有下一个的时候会返回(-1,nil)
func (i *Iter) Next() (int, []byte) {
	if i.index == len(i.keys)-1 {
		if i.nextPageID == -1 {
			// 没有右兄弟
			return -1, nil
		}
		for i.nextPageID != -1 {
			p := i.bf.GetPage(i.nextPageID)
			// 怎么说扫的都应该是叶子节点
			node := FromPage(p, i.bf).(*LeafNode)
			i.nextPageID = node.rightSibling
			i.keys = node.keys
			i.records = node.rids
			i.index = -1
			if len(i.keys) != 0 {
				// 这地方写的很丑陋
				// 不过没办法
				// 访问到空节点的时候就只能这样哦
				// 万一它后面还有节点呢
				break
			}
		}
		// 如果已经是最后一页了 并且key一滴都没有了
		// 就不会有下一次了 呜呜
		if i.nextPageID == -1 && len(i.keys) == 0 {
			return -1, nil
		}
	}
	i.index++
	resKey := i.keys[i.index]
	resBytes := i.records[i.index]
	return resKey, resBytes
}
