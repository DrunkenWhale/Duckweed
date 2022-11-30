package index

type Page struct {
	pageID int
}

func (page *Page) GetPageID() int {
	return page.pageID
}
func (page *Page) ToBytes() []byte {
	panic("Implement me!")
}

func FromBytes(bytes []byte) *Page {
	panic("Implement me!")
}
