package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
)

// index node in disk
//
//
// 	 	   1 byte    8 byte   8 byte	   8 byte
// |head|IsIndexNode|pageID|maxKeysNumber|keyNumber|
// |body|slot(int)|key|key|key...|=>
//			*****
// 		  <=|children(int)|children(int)|
//			   	8 byte        8 byte

// 哦对了 理论上来说 keys实际上要比children小1
// 我们还得写一个numLessThan

const (
	indexHeaderSize         = 17
	_maxIndexNodeKeysNumber = ((page.PageSize - indexHeaderSize) / (2 * 8)) - 3
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
		maxKVNumber: _maxIndexNodeKeysNumber,
		page:        bf.FetchNewPage(),
		keys:        keys,
		children:    children,
	}
}

func (node *IndexNode) Delete(key int) bool {
	index := node.numLessThan(key)
	childNode := node.FetchNode(node.children[index])
	return childNode.Delete(key)
}

func (node *IndexNode) getLeftmostNodeID() int {
	// 因为不打算释放空间 所以索引是不会自己掉的
	// 这里应当不会越界
	// 你要是初始化当我没说
	// 但是要是能访问到这里 说明已经分裂过了不是吗
	// 所以还是不会越界
	child := FromPage(node.bf.GetPage(node.children[0]), node.bf)
	return child.getLeftmostNodeID()
}

func (node *IndexNode) IsLeafNode() bool {
	return false
}

func (node *IndexNode) IsIndexNode() bool {
	return true
}

func (node *IndexNode) Get(key int) ([]byte, bool) {
	index := node.numLessThan(key)
	childNode := node.FetchNode(node.children[index])
	return childNode.Get(key)
}

func (node *IndexNode) Put(key int, value []byte) (int, int, bool) {
	index := node.numLessThan(key)
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
			newChildren := node.children[midIndex+1:]
			splitNode := NewIndexNode(node.bf, newKeys, newChildren)
			node.keys = node.keys[:midIndex]
			node.children = node.children[:midIndex+1]
			splitNode.sync()
			node.sync()
			return splitNode.page.GetPageID(), returnKey, true
		}
		node.sync()
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
	pageIDBytes := databox.IntToBytes(int64(node.page.GetPageID()))
	maxKeysNumberBytes := databox.IntToBytes(int64(node.maxKVNumber))
	keysNumberBytes := databox.IntToBytes(int64(len(node.keys)))
	header = append(header, pageIDBytes[:]...)
	header = append(header, maxKeysNumberBytes[:]...)
	header = append(header, keysNumberBytes[:]...)
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

// 同步这个节点的数据到缓存上的页面上
// 因为传的是指针
// 所以修改的时候会修改缓存中的页面
// (前提是这个页是从缓存中拿取的)
func (node *IndexNode) sync() {
	node.page.WriteBytes(node.ToBytes())
	// 设为脏页
	node.page.Defile()
}

// 返回可能包含num的子节点的下标
func (node *IndexNode) numLessThan(num int) int {
	return numLessThan(node.keys, num)
}

// IndexNodeFromPage
// 从page中构建index node
func IndexNodeFromPage(p *page.Page, bf buffer.BufferPool) *IndexNode {
	bytes := p.GetBytes()
	pageIDBytes := [8]byte{}
	copy(pageIDBytes[:], bytes[1:9])
	pageID := databox.BytesToInt(pageIDBytes)
	if p.GetPageID() != int(pageID) {
		panic("Illegal Page ID: " + string(pageIDBytes[:]))
	}
	maxKeysNumberBytes := [8]byte{}
	copy(maxKeysNumberBytes[:], bytes[9:17])
	maxKeysNumber := databox.BytesToInt(maxKeysNumberBytes)
	keyNumberBytes := [8]byte{}
	copy(keyNumberBytes[:], bytes[17:25])
	keyNumber := databox.BytesToInt(keyNumberBytes)
	childrenNumber := keyNumber + 1
	keys := make([]int, keyNumber)
	children := make([]int, childrenNumber)
	headerOffset := 3*8 + 1
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
