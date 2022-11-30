package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
)

// index node in disk
//
//
// 	 	  1 byte      8 byte	   8 byte
// |head|IsLeafNode|maxKeysNumber|keyNumber|
// |body|slot(int)|key|key|key...|=>
//			*****
// 		  <=|children(int)|children(int)|
//			   	8 byte        8 byte

// 哦对了 理论上来说 keys实际上要比children小1
// 我们还得写一个numLessThan

const FillFactor = 0.75

type IndexNode struct {
	bf          buffer.BufferPool
	maxKVNumber int
	page        *page.Page
	keys        []int // 键后续可能会扩展(多种类型) 但我要想先做个int的试试
	children    []int
}

// TODO 想出有哪些meta信息

func GetIndexNode(pageID int) *IndexNode {
	// TODO 获取索引节点
	return nil
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
	return node.maxKVNumber >= int(FillFactor*float64(len(node.keys)))
}

// 同步这个节点的数据到缓存上的页面上
// 因为传的是指针
// 所以修改的时候会修改缓存中的页面
// (前提是这个页是从缓存中拿取的)
func (node *IndexNode) sync() {
	node.page.WriteBytes(node.ToBytes())
}

// IndexNodeFromPage
// TODO: Test
func IndexNodeFromPage(p *page.Page) *IndexNode {
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
		maxKVNumber: int(maxKeysNumber),
		page:        p,
		keys:        keys,
		children:    children,
	}
	return node
}
