package index

import (
	"Duckweed/buffer"
	"Duckweed/page"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestBPlus(t *testing.T) {
	pool := buffer.NewLRUBufferPool()
	keys := make([]int, 2)
	keys[0] = 1
	keys[1] = 2
	children := make([]int, 3)
	children[0] = 114
	children[1] = 514
	children[2] = 1919810
	node := &IndexNode{
		bf:          pool,
		maxKVNumber: 114514,
		page:        &page.Page{},
		keys:        keys,
		children:    children,
	}
	node.sync()
	node = FromPage(node.page).(*IndexNode)
	assert.Equal(t, node.keys, keys)
	assert.Equal(t, node.children, children)
}
