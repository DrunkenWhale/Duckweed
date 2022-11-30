package buffer

import (
	"Duckweed/disk"
	"Duckweed/lru"
	"Duckweed/page"
)

type LRUBufferPool struct {
	pageNumber int
	pool       map[int]*page.Page
	lru2q      *lru.LRU2Q
	disk       disk.DiskManager
}

func (bf *LRUBufferPool) GetPage(pageID int) *page.Page {
	// 先走LRU
	// 更新置换策略
	if out := bf.lru2q.Push(pageID); out != -1 {
		// 有页面要被刷写到磁盘上
		p, flag := bf.pool[out]
		if !flag {
			// 写map的操作和lru的push一定是一起完成的
			// 不可能不存在于map中
			panic("Illegal Page Request")
		}
		// 把页刷到磁盘上
		bf.disk.Write(p)
		// 释放空间
		delete(bf.pool, pageID)
	}
	p, flag := bf.pool[pageID]
	if flag {
		// 如果在缓存池里
		return p
	}
	// 不在缓存池中
	// 从磁盘上读取页
	pg := bf.disk.Read(pageID)
	// 写入缓存池
	bf.pool[pageID] = pg
	return pg
}

// Flush 这个操作会把所有的页都刷盘
func (bf *LRUBufferPool) Flush() {
	pages := make([]*page.Page, len(bf.pool))
	for i, p := range bf.pool {
		pages[i] = p
	}
	bf.disk.BatchWrite(pages)
	return
}

func (bf *LRUBufferPool) Pin(pageID int) {
	//TODO implement me
	panic("implement me")
}

func (bf *LRUBufferPool) UnPin(pageID int) {
	//TODO implement me
	panic("implement me")
}
