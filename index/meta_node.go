package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
)

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
