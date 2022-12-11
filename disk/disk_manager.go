package disk

import (
	"Duckweed/page"
)

type DiskManager interface {
	// 根据pageID读取page
	Read(pageID int) *page.Page
	// 写一个page到硬盘上
	Write(page *page.Page)
	// 批量写入page
	// 话说这里的效率和单页写入有什么差吗
	// 写的和随机IO似的
	BatchWrite(pages []*page.Page)
	// 获取下一个可以分配的pageID
	GetNextFreePageID() int
	// 清空磁盘上属于该disk manager的全部内容
	// 请谨慎调用！
	Clear()
}
