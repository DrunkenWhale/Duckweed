package index

import (
	"Duckweed/buffer"
	"Duckweed/page"
)

type BPlusNode interface {
	IsLeafNode() bool

	IsIndexNode() bool

	ToBytes() []byte
	// GetPage 存储该节点的页
	GetPage() *page.Page
	// FetchNode 从buffer pool中拿取页转化为node
	FetchNode(pageID int) BPlusNode

	// Put
	// Go 没有Option 将就着写吧
	// @return
	// int:  分裂后的新的节点的page ID 只有在bool值为True时才应当加入children和key中
	// int:  分裂后的Key
	// bool: 子节点是否分裂 分裂为True 未分裂为False
	Put(key int, value []byte) (int, int, bool)

	// Get
	// @return
	// []byte: 值, 下面的bool值为false时 这个值是空的
	// bool  : 是否查找到该key
	Get(key int) ([]byte, bool)
	// 将节点内容同步到缓冲池中的对应页

	// 获取子节点中最右边的那个的ID
	getLeftmostNodeID() int

	sync()
}

func FromPage(page *page.Page, bf buffer.BufferPool) BPlusNode {
	flag := page.GetBytes()[0]
	if flag == IndexNodeFlag {
		return IndexNodeFromPage(page, bf)
	} else if flag == LeafNodeFlag {
		return LeafNodeFromPage(page, bf)
	} else {
		panic("Illegal Page!")
	}
}

const (
	FillFactor    = 0.75
	IndexNodeFlag = iota + 1
	LeafNodeFlag
)
