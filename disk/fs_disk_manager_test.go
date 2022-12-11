package disk

import (
	"Duckweed/page"
	"testing"
)

func TestFSDiskManager_Clear(t *testing.T) {
	manager := NewFSDiskManager("114514")
	manager.Write(page.NewPage(114, make([]byte, 4096)))
	manager.Write(page.NewPage(3, []byte("1919810")))
	manager.Write(page.NewPage(1, []byte("1919810")))
	manager.Clear()
}
