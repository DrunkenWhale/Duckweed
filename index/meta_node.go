package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
	"Duckweed/trans"
)

// 暂定
//     1 byte		8  byte
// meta root flag|root page id|

type TreeMetaNode struct {
	rootPageID int
	rc         trans.Recovery
	bf         buffer.BufferPool
	page       *page.Page
}

func GetTreeMetaNode(bf buffer.BufferPool, rc trans.Recovery) *TreeMetaNode {
	p := bf.GetPage(0)
	if p == nil {
		panic("Illegal Page")
	}
	node := &TreeMetaNode{
		bf: bf,
		rc: rc,
	}
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
	copy(b[:], bytes[1:9])
	rootPageID := databox.BytesToInt(b)
	node.rootPageID = int(rootPageID)
	return
}

func (node *TreeMetaNode) ToBytes() []byte {
	bytes := make([]byte, page.PageSize)
	bytes[0] = MetaNodeFlag
	rootPageIDBytes := databox.IntToBytes(int64(node.rootPageID))
	copy(bytes[1:9], rootPageIDBytes[:])
	return bytes
}

// 同步刷盘
func (node *TreeMetaNode) sync() {
	if node.page.GetBytes() != nil || len(node.page.GetBytes()) != 0 {
		bs := make([]byte, 4096)
		bs[0] = MetaNodeFlag
		b := databox.IntToBytes(int64(-1))
		// 如果是空的 就暂时初始化为-1 代表是空树
		copy(bs[1:9], b[:])
		node.page.WriteBytes(bs)
	}
	node.rc.Record(node.page)

	node.page.WriteBytes(node.ToBytes())
	// 设为脏页
	node.page.Defile()
}
