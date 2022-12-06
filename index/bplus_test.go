package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
	"github.com/go-playground/assert/v2"
	"math/rand"
	"testing"
	"time"
)

func TestBPlus(t *testing.T) {
	pool := buffer.NewLRUBufferPool("duckweed")
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
	node = FromPage(node.page, pool).(*IndexNode)
	assert.Equal(t, node.keys, keys)
	assert.Equal(t, node.children, children)
}

func TestBPlusTree_Put1(t *testing.T) {
	tree := NewBPlusTree(9)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 114514; i++ {
		bytes := databox.IntToBytes(int64(rand.Int()))
		tree.Put(i, bytes[:])
	}
	tree.bf.Flush()
}

func TestBPlusTree_Put2(t *testing.T) {
	tree := NewBPlusTree(8)
	rand.Seed(time.Now().Unix())
	for i, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(i, bytes[:])
	}
	tree.bf.Flush()
}

func TestBPlusTree_Get(t *testing.T) {
	tree := NewBPlusTree(8)
	rand.Seed(time.Now().Unix())
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.bf.Flush()
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		b, f := tree.Get(v)
		assert.Equal(t, f, true)
		assert.Equal(t, b, bytes[:])
	}
}
