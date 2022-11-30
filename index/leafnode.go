package index

import "Duckweed/page"

type LeafNode struct {
}

func (node *LeafNode) GetPage() *page.Page {
	//TODO implement me
	panic("implement me")
}

func (node *LeafNode) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func LeafNodeFromPage(page *page.Page) *LeafNode {
	return nil
}
