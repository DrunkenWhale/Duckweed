package buffer

import (
	"Duckweed/lru"
	"Duckweed/page"
)

type LRUBufferPool struct {
	pageNumber int
	lru2q      *lru.LRU2Q
}

func (bf *LRUBufferPool) GetPage(pageID int) *page.Page {
	//TODO implement me
	panic("implement me")
}

func (bf *LRUBufferPool) Flush() {
	//TODO implement me
	panic("implement me")
}

func (bf *LRUBufferPool) flushToDisk(page page.Page) {

}

func (bf *LRUBufferPool) Pin(pageID int) {
	//TODO implement me
	panic("implement me")
}

func (bf *LRUBufferPool) UnPin(pageID int) {
	//TODO implement me
	panic("implement me")
}
