package disk

import "Duckweed/page"

type DiskManager interface {
	Write(page *page.Page)
	BatchWrite(pages []*page.Page)
	Read(pageID int) *page.Page
}
