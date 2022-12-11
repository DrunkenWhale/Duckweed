package trans

type Recovery interface {
	StartTransaction()
	EndTransaction()
	Commit()
	Abort()
	Rollback()
}
