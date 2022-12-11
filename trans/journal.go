package trans

import "Duckweed/disk"

type JournalRecovery struct {
	disk disk.FSDiskManager
}

func NewJournalRecovery(disk disk.FSDiskManager) *JournalRecovery {
	return &JournalRecovery{disk: disk}
}

func (r *JournalRecovery) StartTransaction() {
	//TODO implement me
	panic("implement me")
}

func (r *JournalRecovery) EndTransaction() {
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
