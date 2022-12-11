package disk

import (
	"Duckweed/page"
	"io/fs"
	"os"
)

const (
	StoragePath = "data"
	FileSuffix  = ".dued"
)

type FSDiskManager struct {
	filename    string
	pageSize    int
	storagePath string
	file        *os.File
}

func NewFSDiskManager(filename string) *FSDiskManager {
	dm := &FSDiskManager{
		pageSize:    page.PageSize,
		storagePath: StoragePath,
		file:        nil,
		filename:    filename, // 创建的时候会创建成 /data/duckweed/duckweed这样的文件
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
	storagePath := dm.storagePath + string(os.PathSeparator) + dm.filename
	_, err := os.ReadDir(dm.storagePath)
	if err != nil {
		// 父目录不存在则递归创建
		err := os.MkdirAll(storagePath, fs.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	f, err := os.OpenFile(
		storagePath+string(os.PathSeparator)+dm.filename+FileSuffix,
		os.O_CREATE|os.O_RDWR,
		fs.ModePerm)
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

// GetNextFreePageID page是从0开始的
func (dm *FSDiskManager) GetNextFreePageID() int {
	info, err := dm.file.Stat()
	if err != nil {
		panic(err)
	}
	if info.Size()%int64(page.PageSize) != 0 {
		panic("Illegal DB File (Illegal Size)")
	}
	return int(info.Size() / int64(page.PageSize))
}

func (dm *FSDiskManager) Clear() {
	dm.file.Truncate(0)
}
