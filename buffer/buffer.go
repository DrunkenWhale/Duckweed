package buffer

import "Duckweed/page"

const (
	// MaxPageNumber 最大缓存 16 * 4096 Byte
	MaxPageNumber = 16
)

type BufferPool interface {
	// GetPage from disk
	GetPage(pageID int) *page.Page
	// Flush all dirty page to disk
	Flush()
	//// Pin 固定一个页面 防止其释放
	//Pin(pageID int)
	//// UnPin 放开一个页面捏
	//UnPin(pageID int)
}
