package page

const PageSize int = 4096

type Page struct {
	isDirty bool
	pageID  int
	bytes   []byte
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

// 从现有的page中复制出一个一模一样的Page
// 但是他们之间的修改是独立的
func (p *Page) Copy() *Page {
	bytes := make([]byte, PageSize)
	copy(bytes, p.GetBytes())
	page := NewPage(p.pageID, bytes)
	return page
}

// 查看某一页是否为脏页
func (p *Page) IsDirty() bool {
	return p.isDirty
}

// 把某一页置为脏页
func (p *Page) Defile() {
	p.isDirty = true
	return
}
