package buffer

import "Duckweed/index"

const (
	// MaxPageNumber 最大缓存 16 * 4096 Byte
	MaxPageNumber = 16
)

type BufferPool interface {
	// GetPage from disk
	GetPage(pageID int) *index.Page
	// Flush all dirty index to disk
	Flush()
	//// Pin 固定一个页面 防止其释放
	//Pin(pageID int)
	//// UnPin 放开一个页面捏
	//UnPin(pageID int)
}
