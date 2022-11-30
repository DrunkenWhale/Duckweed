package disk

import "Duckweed/index"

type DiskManager interface {
	Write(page *index.Page)
	Read(pageID int) *index.Page
	BatchWrite(pages []*index.Page)
}
