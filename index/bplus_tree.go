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
