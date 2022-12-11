package trans

import (
	"Duckweed/disk"
	"Duckweed/page"
)

type JournalRecovery struct {
	// journal日志的disk manager
	journalDiskManager disk.DiskManager

	// 数据库的disk manager
	dbDiskManager disk.DiskManager

	// 事务是否已经开启
	isInTransaction bool

	// pageID map to page
	// 注意 这里的page应当是copy产生的
	// 即被修改之前的page内容
	dirtyPageTable map[int]*dirtyPageInfo
}

func NewJournalRecovery(dbDiskManager, journalDiskManager disk.DiskManager) *JournalRecovery {
	return &JournalRecovery{
		journalDiskManager: journalDiskManager,
		dbDiskManager:      dbDiskManager,
		isInTransaction:    false,
		dirtyPageTable:     make(map[int]*dirtyPageInfo),
	}
}

func (r *JournalRecovery) StartTransaction() {
	if r.isInTransaction {
		panic("You can't nesting transaction")
	}
	r.isInTransaction = true
}

func (r *JournalRecovery) End() {
	r.dirtyPageTable = make(map[int]*dirtyPageInfo)
	r.isInTransaction = false
	r.journalDiskManager.Clear()

}

// 这个居然要配合外部实现 零昏！
func (r *JournalRecovery) Commit() {
	r.End()
	return
}

func (r *JournalRecovery) Rollback() {
	pageNum := r.journalDiskManager.GetNextFreePageID()
	for i := 0; i < pageNum; i++ {
		p := r.journalDiskManager.Read(i)
		r.dbDiskManager.Write(p.GetPageID(), p)
	}
	r.End()
	return
}

// 写入的原页需要满足
// 1: 原数据页不能是空的
// 2: 已经写入磁盘过
func (r *JournalRecovery) Record(p *page.Page) {
	if !r.isInTransaction && p.GetPageID() != 0 {
		// 如果不处于事务状态中则不记录
		// 因为初始化的时候也会用到这个方法 所以得特判一下
		//log.Println("operation not in transaction")
		return
	}
	_, flag := r.dirtyPageTable[p.GetPageID()]
	if !flag {
		// 如果不存在于脏页表中
		// 那么要写进去捏
		// 这样在回滚的时候才能恢复到最初的状态
		r.dirtyPageTable[p.GetPageID()] =
			&dirtyPageInfo{
				hasFlashed: false,
				page:       p.Copy(),
			}
		return
	}
	// 如果已经存在哩
	// 那么覆写就没有意义了
	// 我们应该保持它最初的状态捏
	return
}

func (r *JournalRecovery) WriteBackups(pageID int) {
	if !r.isInTransaction {
		// 如果不处于事务状态中则不写入
		//log.Println("operation not in transaction")
		return
	}
	info, flag := r.dirtyPageTable[pageID]
	if !flag {
		// 不存在这个page
		// 无需刷写
		return
	}
	if info.hasFlashed {
		// 如果它曾经被刷写过
		// 没有必要重复刷写备份文件
		return
	}
	// 刷写日志文件到磁盘
	// 并且标记为已刷写过
	// 这里的磁盘上page的位置未必和page里的id相照应
	r.journalDiskManager.Write(r.journalDiskManager.GetNextFreePageID(), info.page)
	info.hasFlashed = true
}

type dirtyPageInfo struct {
	page *page.Page

	// 是否已刷写过
	// 刷写过的备份页面无需再次刷写
	hasFlashed bool
}
