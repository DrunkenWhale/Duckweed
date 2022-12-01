package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
)

// index node in disk
//
//
// 	 	   1 byte      8 byte	   8 byte
// |head|IsIndexNode|maxKeysNumber|keyNumber|
// |body|slot(int)|key|key|key...|=>
//			*****
// 		  <=|children(int)|children(int)|
//			   	8 byte        8 byte

// 哦对了 理论上来说 keys实际上要比children小1
// 我们还得写一个numLessThan

const (
	FillFactor      = 0.75
	indexHeaderSize = 17
	maxKeysNumber   = ((page.PageSize - indexHeaderSize) / 2 * 8) - 3
)

type IndexNode struct {
	bf          buffer.BufferPool
	maxKVNumber int
	page        *page.Page
	keys        []int // 键后续可能会扩展(多种类型) 但我要想先做个int的试试
	children    []int
}

func NewIndexNode(bf buffer.BufferPool, keys []int, children []int) *IndexNode {
	return &IndexNode{
		bf:          bf,
		maxKVNumber: maxKeysNumber,
		page:        bf.FetchNewPage(),
		keys:        keys,
		children:    children,
	}
}

func (node *IndexNode) IsLeafNode() bool {
	return false
}

func (node *IndexNode) IsIndexNode() bool {
	return true
}

func (node *IndexNode) Put(key int, value []byte) (int, int, bool) {
	index := node.numLessThanEqual(key)
	childNode := node.FetchNode(node.children[index])
	newNodePageID, splitKey, isSplit := childNode.Put(key, value)
	if isSplit {
		// 如果子节点分裂了
		node.keys = insertSliceWithIndex(node.keys, index, splitKey)
		node.children = insertSliceWithIndex(node.children, index+1, newNodePageID)
		if node.shouldSplit() {
			// 数量过多 需要进行分裂
			midIndex := len(node.keys) / 2
			returnKey := node.keys[midIndex]
			newKeys := node.keys[midIndex+1:]
			newChildren := node.children[midIndex:]
			splitNode := NewIndexNode(node.bf, newKeys, newChildren)
			node.keys = node.keys[:midIndex]
			node.children = node.children[:midIndex]
			splitNode.sync()
			node.sync()
			return splitNode.page.GetPageID(), returnKey, true
		}
	}
	return -1, -1, false
}

func (node *IndexNode) FetchNode(pageID int) BPlusNode {
	p := node.bf.GetPage(pageID)
	n := FromPage(p, node.bf)
	return n
}

func (node *IndexNode) GetPage() *page.Page {
	return node.page
}

func (node *IndexNode) ToBytes() []byte {
	header := make([]byte, 1)
	header[0] = IndexNodeFlag
	if len(node.keys)-len(node.children) != -1 {
		// 数值不对啊
		panic("Keys Should equal ChildrenNumber - 1 ╰(*°▽°*)╯ ")
	}
	maxKeysNumberBytes := databox.IntToBytes(int64(node.maxKVNumber))
	keysNumberBytes := databox.IntToBytes(int64(len(node.keys)))
	header = append(header, append(maxKeysNumberBytes[:], keysNumberBytes[:]...)...)
	keysBytes := make([]byte, 8*len(node.keys))
	for i := 0; i < len(node.keys); i++ {
		b := databox.IntToBytes(int64(node.keys[i]))
		copy(keysBytes[i*8:(i+1)*8], b[:])
	}
	childrenBytes := make([]byte, 8*len(node.children))
	for i := 0; i < len(node.children); i++ {
		b := databox.IntToBytes(int64(node.children[len(node.children)-1-i]))
		copy(childrenBytes[i*8:(i+1)*8], b[:])
	}
	blankBytesSize := page.PageSize - (len(header) + len(keysBytes) + len(childrenBytes))
	blankBytes := make([]byte, blankBytesSize)
	bytes := append(header, keysBytes...)
	bytes = append(bytes, blankBytes...)
	bytes = append(bytes, childrenBytes...)
	return bytes
}

func (node *IndexNode) shouldSplit() bool {
	return len(node.keys) >= int(FillFactor*float64(node.maxKVNumber))
}

func (node *IndexNode) splitKeys() bool {
	return len(node.keys) >= int(FillFactor*float64(node.maxKVNumber))
}

// 同步这个节点的数据到缓存上的页面上
// 因为传的是指针
// 所以修改的时候会修改缓存中的页面
// (前提是这个页是从缓存中拿取的)
func (node *IndexNode) sync() {
	node.page.WriteBytes(node.ToBytes())
}

// 返回可能包含num的子节点的下标
func (node *IndexNode) numLessThanEqual(num int) int {
	return numLessThanEqual(node.keys, num)
}

// IndexNodeFromPage
// TODO: Test
func IndexNodeFromPage(p *page.Page, bf buffer.BufferPool) *IndexNode {
	bytes := p.GetBytes()
	maxKeysNumberBytes := [8]byte{}
	copy(maxKeysNumberBytes[:], bytes[1:9])
	maxKeysNumber := databox.BytesToInt(maxKeysNumberBytes)
	keyNumberBytes := [8]byte{}
	copy(keyNumberBytes[:], bytes[9:17])
	keyNumber := databox.BytesToInt(keyNumberBytes)
	childrenNumber := keyNumber + 1
	keys := make([]int, keyNumber)
	children := make([]int, childrenNumber)
	headerOffset := 2*8 + 1
	for i := 0; i < int(keyNumber); i++ {
		b := [8]byte{}
		copy(b[:], bytes[headerOffset+i*8:headerOffset+(i+1)*8])
		num := databox.BytesToInt(b)
		keys[i] = int(num)
	}
	for i := 0; i < int(childrenNumber); i++ {
		b := [8]byte{}
		copy(b[:], bytes[page.PageSize-(i+1)*8:page.PageSize-i*8])
		num := databox.BytesToInt(b)
		children[i] = int(num)
	}
	node := &IndexNode{
		bf:          bf,
		maxKVNumber: int(maxKeysNumber),
		page:        p,
		keys:        keys,
		children:    children,
	}
	return node
}
