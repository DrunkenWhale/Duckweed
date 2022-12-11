package trans

import (
	"Duckweed/disk"
	"Duckweed/page"
)

type JournalRecovery struct {
	// 磁盘管理器是肯定的要的捏
	disk disk.DiskManager
	// 事务是否已经开启
	isInTransaction bool

	// pageID map to page
	// 注意 这里的page应当是copy产生的
	// 即被修改之前的page内容
	dirtyPageTable map[int]*page.Page
}

func NewJournalRecovery(disk disk.DiskManager) *JournalRecovery {
	return &JournalRecovery{
		disk:            disk,
		isInTransaction: false,
	}
}

func (r *JournalRecovery) StartTransaction() {
	if r.isInTransaction {
		panic("You can't nesting transaction")
	}
	r.isInTransaction = true
}

func (r *JournalRecovery) End() {

	//TODO implement me
	panic("implement me")
}

func (r *JournalRecovery) Commit() {
	//TODO implement me
	panic("implement me")
}

func (r *JournalRecovery) Abort() {
	//TODO implement me
	panic("implement me")
}

func (r *JournalRecovery) Rollback() {
	//TODO implement me
	panic("implement me")
}

// 写入的原页需要满足
// 1: 原数据页不能是空的
// 2: 已经写入磁盘过
func (r *JournalRecovery) Record(p *page.Page) {
	_, flag := r.dirtyPageTable[p.GetPageID()]
	if !flag {
		// 如果不存在于脏页表中
		// 那么要写进去捏
		// 这样在回滚的时候才能恢复到最初的状态
		r.dirtyPageTable[p.GetPageID()] = p.Copy()
		return
	}
	// 如果已经存在哩
	// 那么覆写就没有意义了
	// 我们应该保持它最初的状态捏
	return
}
