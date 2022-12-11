package index

import (
	"Duckweed/buffer"
	"Duckweed/disk"
	"Duckweed/page"
	"Duckweed/trans"
)

type BPlusTree struct {
	bf        buffer.BufferPool
	dm        disk.DiskManager
	rc        trans.Recovery
	root      BPlusNode
	ridLength int
}

func NewBPlusTree(name string, ridLength int) *BPlusTree {
	journalDiskManager := disk.NewFSDiskManager(name + "-journal")
	dbDiskManager := disk.NewFSDiskManager(name)
	rc := trans.NewJournalRecovery(dbDiskManager, journalDiskManager)
	bf := buffer.NewLRUBufferPool(dbDiskManager, rc)
	tree := &BPlusTree{
		bf:        bf,
		dm:        dbDiskManager,
		rc:        rc,
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
		newRoot := NewIndexNode(tree.bf, tree.rc, newRootKeys, newRootChildren)
		tree.root = newRoot
		tree.root.sync()

		// 储存 meta 信息的节点也要刷新一遍
		GetTreeMetaNode(tree.bf, tree.rc).
			setRootNodeIDAndSync(newRoot.GetPage().GetPageID())
		return
	}
	tree.root.sync()
}

func (tree *BPlusTree) Scan() *Iter {
	leftmostNodeID := tree.root.getLeftmostNodeID()
	return NewIter(leftmostNodeID, tree.bf, tree.rc)
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

func (tree *BPlusTree) StartTransaction() {
	tree.rc.StartTransaction()
	GetTreeMetaNode(tree.bf, tree.rc).sync()
	tree.rc.WriteBackups(1)
}

func (tree *BPlusTree) Commit() {
	tree.rc.Commit()
	tree.bf.Flush()
}

func (tree *BPlusTree) Rollback() {
	tree.rc.Rollback()
	tree.bf.Clear()
	tree.root = nil
	tree.init()
}

func (tree *BPlusTree) init() {
	p := tree.bf.GetPage(0)
	var root BPlusNode

	if p == nil {
		// pageID=0的页不存在且不合法 或是没有下一个节点
		// 说明这棵树是空的
		node := &TreeMetaNode{
			rootPageID: 0,
			bf:         tree.bf,
			rc:         tree.rc,
			page:       tree.bf.FetchNewPage(),
		}
		if node.page.GetPageID() != 0 {
			panic("Not an empty B+ Tree, please check race condition!")
		}
		root = NewLeafNode(tree.bf, tree.rc, tree.ridLength, -1, make([]int, 0), make([][]byte, 0))
		node.setRootNodeIDAndSync(root.GetPage().GetPageID())
	} else {
		node := GetTreeMetaNode(tree.bf, tree.rc)
		// 树非空
		if node.GetRootNodeID() != -1 {
			rootNodePage := tree.bf.GetPage(node.GetRootNodeID())
			root = FromPage(rootNodePage, tree.bf, tree.rc)
		} else {
			node := &TreeMetaNode{
				rootPageID: 0,
				bf:         tree.bf,
				rc:         tree.rc,
				page:       tree.bf.FetchNewPage(),
			}
			// 这边属于又双叒叕的抽象泄露了
			// 由于page和page ID绑死了
			// 这里会给出错误的page ID
			// 得手动设置为0
			// 哦 我知道你想问其他页
			// 那些不需要
			// 因为页中有存储真正的page ID
			// 但是meta node page一定是在index=0的页中
			// 所以需要区别对待
			if node.page.GetPageID() != 0 {
				node.page = page.NewPage(0, node.page.Copy().GetBytes())
			}
			root = NewLeafNode(tree.bf, tree.rc, tree.ridLength, -1, make([]int, 0), make([][]byte, 0))
			node.setRootNodeIDAndSync(root.GetPage().GetPageID())
		}
	}
	tree.root = root
}
