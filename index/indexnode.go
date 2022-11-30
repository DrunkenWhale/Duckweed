package index

import "Duckweed/page"

type IndexNode struct {
	page     *page.Page
	keys     []int // 键后续可能会扩展(多种类型) 但我要想先做个int的试试
	children []int
}

// TODO 想出有哪些meta信息

func (node *IndexNode) GetPage() *page.Page {
	return node.page
}

func (node *IndexNode) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func IndexNodeFromPage(page *page.Page) *IndexNode {
	node := &IndexNode{
		page: page,
	}
	return node
}
