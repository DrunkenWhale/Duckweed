package disk

import (
	"Duckweed/page"
	"io/fs"
	"os"
)

const (
	StoragePath = "data"
)

type FSDiskManager struct {
	pageSize    int
	storagePath string
	file        *os.File
}

func NewFSDiskManager() *FSDiskManager {
	dm := &FSDiskManager{
		pageSize:    page.PageSize,
		storagePath: StoragePath,
		file:        nil,
	}
	dm.open()
	return dm
}

func (dm *FSDiskManager) Write(page *page.Page) {
	pageID := page.GetPageID()
	bytes := page.GetBytes()
	dm.writePageBytes(bytes, pageID)
}

func (dm *FSDiskManager) Read(pageID int) *page.Page {
	bytes := dm.readPageBytes(pageID)
	return page.NewPage(pageID, bytes)
}

func (dm *FSDiskManager) BatchWrite(pages []*page.Page) {
	for _, p := range pages {
		dm.Write(p)
	}
}

func (dm *FSDiskManager) open() {
	_, err := os.ReadDir(dm.storagePath)
	if err != nil {
		// 父目录不存在则递归创建
		err := os.MkdirAll(dm.storagePath, fs.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	f, err := os.OpenFile(dm.storagePath+string(os.PathSeparator)+"duckweed", os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		panic(err)
	}
	dm.file = f
}

func (dm *FSDiskManager) close() {
	err := dm.file.Close()
	if err != nil {
		panic(err)
	}
}

// 没用DirectIO
// 后续可以考虑做一个试试
func (dm *FSDiskManager) readPageBytes(pageID int) []byte {
	offset := pageID * dm.pageSize
	bytes := make([]byte, dm.pageSize)
	_, err := dm.file.ReadAt(bytes, int64(offset))
	if err != nil {
		panic(err)
	}
	return bytes
}

func (dm *FSDiskManager) writePageBytes(bytes []byte, pageID int) {
	offset := pageID * dm.pageSize
	_, err := dm.file.WriteAt(bytes, int64(offset))
	if err != nil {
		panic(err)
	}
}
