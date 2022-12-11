package trans

import "Duckweed/page"

type Recovery interface {
	StartTransaction()
	End()
	Commit()
	Abort()
	Rollback()

	// 记录某个page
	// 这里我们只允许单进程单事务嘛
	// 所以recovery manager状态是唯一的捏
	// 写进这里的数据会在commit的时候保存
	// 并且保证不会掉盘
	Record(page *page.Page)
}
