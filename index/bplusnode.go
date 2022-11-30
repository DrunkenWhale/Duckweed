package index

import "Duckweed/page"

type BPlusNode interface {
	ToBytes() []byte
	// GetPage 存储该节点的页
	GetPage() *page.Page
	// FetchNode 从buffer pool中拿取页转化为node
	FetchNode(pageID int) BPlusNode
}

func FromPage(page *page.Page) BPlusNode {
	flag := page.GetBytes()[0]
	if flag == IndexNodeFlag {
		return IndexNodeFromPage(page)
	} else if flag == LeafNodeFlag {
		return LeafNodeFromPage(page)
	} else {
		panic("Illegal Page!")
	}
}

const (
	IndexNodeFlag = iota + 1
	LeafNodeFlag
)
