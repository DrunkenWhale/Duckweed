package index

import "Duckweed/buffer"

type BPlusTree struct {
	bf        buffer.BufferPool
	root      BPlusNode
	ridLength int
}

// TODO 确保磁盘上有数据时 能够找到正确的根节点

func NewBPlusTree(ridLength int) *BPlusTree {
	bf := buffer.NewLRUBufferPool()
	p := bf.GetPage(1)
	var root BPlusNode
	if p == nil {
		// 改页不存在且不合法
		// 说明这棵树是空的
		root = NewLeafNode(bf, ridLength, -1, make([]int, 0), make([][]byte, 0))
	} else {
		// 该页存在 是有效的
		root = FromPage(p, bf)
	}
	return &BPlusTree{
		bf:        bf,
		root:      root,
		ridLength: ridLength,
	}
}

func (tree *BPlusTree) Put(key int, bytes []byte) {
	newNodeID, splitKey, isSplit := tree.root.Put(key, bytes)
	if isSplit {
		// 如果分裂了
		// 说明根节点要往上一层
		newRootKeys := make([]int, 1)
		newRootKeys[0] = splitKey
		newRootChildren := make([]int, 2)
		newRootChildren[0] = tree.root.GetPage().GetPageID()
		newRootChildren[1] = newNodeID
		newRoot := NewIndexNode(tree.bf, newRootKeys, newRootChildren)
		tree.root = newRoot
		tree.root.sync()
		return
	}
	tree.root.sync()
}
