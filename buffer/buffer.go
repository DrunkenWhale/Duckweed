package buffer

import (
	"Duckweed/page"
)

const (
	// MaxPageNumber 最大缓存 16 * PageSize Byte
	MaxPageNumber = 16
)

type BufferPool interface {
	// GetPage from disk
	GetPage(pageID int) *page.Page
	//向缓存池中加入页面
	//PutPage(page *page.Page)

	// Flush all dirty index to disk
	Flush()

	// 丢弃内存中的所有页
	Clear()

	FetchNewPage() *page.Page
	//// Pin 固定一个页面 防止其释放
	//Pin(pageID int)
	//// UnPin 放开一个页面捏
	//UnPin(pageID int)
}
