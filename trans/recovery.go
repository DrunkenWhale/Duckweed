package trans

import "Duckweed/page"

type Recovery interface {
	StartTransaction()
	End()
	Commit()
	Rollback()

	/***************** 这下面的接口其实是不对的 实际上 应该开一个transaction的struct来处理下面的调用 ****************/
	// 记录某个page
	// 这里我们只允许单进程单事务嘛
	// 所以recovery manager状态是唯一的捏
	// 写进这里的数据会在commit的时候保存
	// 并且保证不会掉盘
	Record(page *page.Page)

	// 通过pageID获取指定的page
	// 并且 如果有必要
	// 刷写到磁盘上(指journal)
	WriteBackups(pageID int)
}
