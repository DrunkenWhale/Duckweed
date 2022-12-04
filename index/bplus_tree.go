package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
)

type BPlusTree struct {
	bf        buffer.BufferPool
	root      BPlusNode
	ridLength int
}

func NewBPlusTree(ridLength int) *BPlusTree {
	bf := buffer.NewLRUBufferPool()
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

// 暂定
//     8  byte
// |root page id|

type TreeMetaNode struct {
	rootPageID int
	page       *page.Page
}

func GetTreeMetaNode(bf buffer.BufferPool) *TreeMetaNode {
	p := bf.GetPage(0)
	if p == nil {
		panic("Illegal Page")
	}
	node := &TreeMetaNode{}
	node.page = p
	node.FromBytes(p.GetBytes())
	return node
}

func (node *TreeMetaNode) GetRootNodeID() int {
	return node.rootPageID
}

func (node *TreeMetaNode) setRootNodeIDAndSync(newRootID int) {
	node.rootPageID = newRootID
	node.sync()
}

func (node *TreeMetaNode) FromBytes(bytes []byte) {
	b := [8]byte{}
	copy(b[:], bytes[:8])
	rootPageID := databox.BytesToInt(b)
	node.rootPageID = int(rootPageID)
	return
}

func (node *TreeMetaNode) ToBytes() []byte {
	bytes := make([]byte, page.PageSize)
	rootPageIDBytes := databox.IntToBytes(int64(node.rootPageID))
	copy(bytes[:8], rootPageIDBytes[:])
	return bytes
}

// 同步刷盘
func (node *TreeMetaNode) sync() {
	node.page.WriteBytes(node.ToBytes())
}
