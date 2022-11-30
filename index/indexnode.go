package index

import "Duckweed/page"

type IndexNode struct {
}

func (node *IndexNode) GetPage() *page.Page {
	//TODO implement me
	panic("implement me")
}

func (node *IndexNode) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func IndexNodeFromPage(page *page.Page) *IndexNode {
	return nil
}
