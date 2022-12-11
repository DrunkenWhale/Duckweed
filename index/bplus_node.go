package index

import (
	"Duckweed/buffer"
	"Duckweed/page"
	"Duckweed/trans"
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
	// 关于key重复是覆盖而非报duplicate这件事
	// 就当是我很懒吧
	// 诶嘿(o゜▽゜)o☆
	Put(key int, value []byte) (int, int, bool)

	// Get
	// @return
	// []byte: 值, 下面的bool值为false时 这个值是空的
	// bool  : 是否查找到该key
	Get(key int) ([]byte, bool)

	// @return 返回删除是否成功
	Delete(key int) bool

	// 获取子节点中最右边的那个的ID
	getLeftmostNodeID() int

	// 将节点内容同步到缓冲池中的对应页
	sync()
}

func FromPage(page *page.Page, bf buffer.BufferPool, rc trans.Recovery) BPlusNode {
	flag := page.GetBytes()[0]
	var node BPlusNode = nil
	if flag == IndexNodeFlag {
		node = IndexNodeFromPage(page, bf, rc)
	} else if flag == LeafNodeFlag {
		node = LeafNodeFromPage(page, bf, rc)
	} else {
		panic("Illegal Page!")
	}
	return node
}

const (
	FillFactor    = 0.75
	IndexNodeFlag = iota + 1
	LeafNodeFlag
	MetaNodeFlag
)
