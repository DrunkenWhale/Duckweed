package index

import "Duckweed/page"

type BPlusNode interface {
	ToBytes() []byte
	// GetPage 存储该节点的页
	GetPage() *page.Page
}

func FromPage(page *page.Page) BPlusNode {
	flag := page.GetBytes()[0]
	if flag == 1 {
		return IndexNodeFromPage(page)
	} else if flag == 2 {
		return LeafNodeFromPage(page)
	} else {
		panic("Illegal Page!")
	}
}
