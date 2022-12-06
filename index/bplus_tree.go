package index

import (
	"Duckweed/buffer"
)

type BPlusTree struct {
	bf        buffer.BufferPool
	root      BPlusNode
	ridLength int
}

func NewBPlusTree(ridLength int) *BPlusTree {
	bf := buffer.NewLRUBufferPool("duckweed")
	tree := &BPlusTree{
		bf:        bf,
		root:      nil,
		ridLength: ridLength,
	}
	tree.init()
	return tree
}

func (tree *BPlusTree) Get(key int) ([]byte, bool) {
	return tree.root.Get(key)
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

		// 储存 meta 信息的节点也要刷新一遍
		GetTreeMetaNode(tree.bf).setRootNodeIDAndSync(newRoot.GetPage().GetPageID())
		return
	}
	tree.root.sync()
}

func (tree *BPlusTree) Scan() *Iter {
	leftmostNodeID := tree.root.getLeftmostNodeID()
	return NewIter(leftmostNodeID, tree.bf)
}

func (tree *BPlusTree) Update(key int, value []byte) {
	if _, ok := tree.Get(key); ok {
		// 键存在
		tree.root.Put(key, value)
	}
}

func (tree *BPlusTree) Delete(key int) {
	tree.root.Delete(key)
}

func (tree *BPlusTree) init() {
	p := tree.bf.GetPage(0)
	var root BPlusNode
	if p == nil {
		// pageID=0的页不存在且不合法
		// 说明这棵树是空的
		node := &TreeMetaNode{
			rootPageID: 0,
			page:       tree.bf.FetchNewPage(),
		}
		if node.page.GetPageID() != 0 {
			panic("Not an empty B+ Tree, please check race condition!")
		}
		root = NewLeafNode(tree.bf, tree.ridLength, -1, make([]int, 0), make([][]byte, 0))
		node.setRootNodeIDAndSync(root.GetPage().GetPageID())
	} else {
		// 该页存在 是有效的
		node := GetTreeMetaNode(tree.bf)
		rootNodePage := tree.bf.GetPage(node.GetRootNodeID())
		root = FromPage(rootNodePage, tree.bf)
	}
	tree.root = root
}
