package index

import "Duckweed/page"

type LeafNode struct {
	maxKVNumber  int
	page         *page.Page
	keys         []int // 键后续可能会扩展(多种类型) 但我要想先做个int的试试
	rids         []string
	rightSibling int
}

func (node *LeafNode) GetPage() *page.Page {
	return node.page
}

func (node *LeafNode) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func LeafNodeFromPage(page *page.Page) *LeafNode {
	return nil
}
