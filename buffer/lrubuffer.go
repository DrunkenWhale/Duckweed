package buffer

import (
	"Duckweed/disk"
	"Duckweed/lru"
	"Duckweed/page"
	"Duckweed/trans"
)

type LRUBufferPool struct {
	nextFreePageID int
	pageNumber     int
	pool           map[int]*page.Page
	lru2q          *lru.LRU2Q
	disk           disk.DiskManager
	recovery       trans.Recovery
}

func (bf *LRUBufferPool) FetchNewPage() *page.Page {
	p := page.NewPage(bf.nextFreePageID, make([]byte, 0))
	bf.pool[p.GetPageID()] = p
	bf.nextFreePageID++
	return p
}

func NewLRUBufferPool(diskManager disk.DiskManager, recovery trans.Recovery) *LRUBufferPool {
	bf := &LRUBufferPool{
		pageNumber: MaxPageNumber,
		pool:       make(map[int]*page.Page),
		lru2q:      lru.NewLRU2Q(MaxPageNumber/4*3, MaxPageNumber/4),
		disk:       diskManager,
		recovery:   recovery,
	}
	bf.nextFreePageID = bf.disk.GetNextFreePageID()
	return bf
}

//func (bf *LRUBufferPool) PutPage(page *page.Page) {
//	bf.lru2q.Push(page.GetPageID())
//	bf.pool[page.GetPageID()] = page
//}

// GetPage
// 当某个页面非法的时候 会返回nil
// 即 不是从磁盘中读取
// 或是不从buffer pool中获取的page
func (bf *LRUBufferPool) GetPage(pageID int) *page.Page {
	if bf.disk.GetNextFreePageID() <= pageID {
		// 该页还没进磁盘
		_, flag := bf.pool[pageID]
		if !flag {
			// 如果它也不存在于内存
			// 那么这是一个未被分配的非法页
			return nil
		}
	}
	// 先走LRU
	// 更新置换策略
	if out := bf.lru2q.Push(pageID); out != -1 {
		// 有页面要被刷写到磁盘上
		p, flag := bf.pool[out]
		if !flag {
			// 写map的操作和lru的push一定是一起完成的
			// 不可能不存在于map中
			panic("Illegal BPlusNode Request")
		}
		// 先刷备份文件
		bf.recovery.WriteBackups(p.GetPageID())
		// 把页刷到磁盘上
		bf.disk.Write(p.GetPageID(), p)
		// 释放空间
		delete(bf.pool, p.GetPageID())
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
	i := 0
	for _, p := range bf.pool {
		pages[i] = p
		i++
	}
	bf.disk.BatchWrite(pages)
	return
}
