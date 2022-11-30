package page

const PageSize = 4096

type Page struct {
	pageID int
	bytes  []byte
}

func NewPage(pageID int, bytes []byte) *Page {
	return &Page{
		pageID: pageID,
		bytes:  bytes,
	}
}

func (p *Page) GetPageID() int {
	return p.pageID
}

// GetBytes 不能修改这个Byte数组
func (p *Page) GetBytes() []byte {
	return p.bytes
}

func (p *Page) WriteBytes(bytes []byte) {
	b := make([]byte, len(bytes))
	copy(b, bytes)
	p.bytes = b
}
