package disk

import (
	"Duckweed/index"
	"fmt"
)

type DummyDiskManager struct {
}

func NewDummyDiskManager() *DummyDiskManager {
	return &DummyDiskManager{}
}

func (dm *DummyDiskManager) Write(page *index.Page) {
	fmt.Printf("Write Page(ID=%d)To \n", page.GetPageID())
}

func (dm *DummyDiskManager) BatchWrite(pages []*index.Page) {
	fmt.Println("Start Batch Write")
	for _, p := range pages {
		dm.Write(p)
	}
	fmt.Println("End Batch Write")
	return
}

func (dm *DummyDiskManager) Read(pageID int) *index.Page {
	return nil
}
