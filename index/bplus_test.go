package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/disk"
	"Duckweed/page"
	"Duckweed/trans"
	"github.com/go-playground/assert/v2"
	"math/rand"
	"testing"
	"time"
)

func TestBPlus(t *testing.T) {
	pool := buffer.NewLRUBufferPool(
		disk.NewFSDiskManager("duckweed"),
		trans.NewJournalRecovery(disk.NewFSDiskManager("duckweed-journal")),
	)
	keys := make([]int, 2)
	keys[0] = 1
	keys[1] = 2
	children := make([]int, 3)
	children[0] = 114
	children[1] = 514
	children[2] = 1919810
	node := &IndexNode{
		bf:          pool,
		rc:          nil,
		maxKVNumber: 114514,
		page:        &page.Page{},
		keys:        keys,
		children:    children,
	}
	node.sync()
	node = FromPage(node.page, pool, node.rc).(*IndexNode)
	assert.Equal(t, node.keys, keys)
	assert.Equal(t, node.children, children)
}

func TestBPlusTree_Put1(t *testing.T) {
	tree := NewBPlusTree("duckweed", 9)
	tree.StartTransaction()
	rand.Seed(time.Now().Unix())
	for i := 0; i < 114514; i++ {
		bytes := databox.IntToBytes(int64(rand.Int()))
		tree.Put(i, bytes[:])
	}
	tree.bf.Flush()
}

func TestBPlusTree_Put2(t *testing.T) {
	tree := NewBPlusTree("duckweed", 8)
	rand.Seed(time.Now().Unix())
	for i, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(i, bytes[:])
	}
	tree.bf.Flush()
}

func TestBPlusTree_Get(t *testing.T) {
	tree := NewBPlusTree("duckweed", 8)
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

func TestBPlusTree_Scan(t *testing.T) {
	tree := NewBPlusTree("duckweed", 8)
	rand.Seed(time.Now().Unix())
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.bf.Flush()
	it := tree.Scan()
	for true {
		k, v := it.Next()
		if k == -1 {
			break
		}
		res, f := tree.Get(k)
		assert.Equal(t, f, true)
		assert.Equal(t, res, v)
	}
}

func TestBPlusTree_Update(t *testing.T) {
	tree := NewBPlusTree("duckweed", 8)
	rand.Seed(time.Now().Unix())
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.bf.Flush()
	for i, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(i))
		tree.Update(v, bytes[:])
		b, ok := tree.Get(v)
		assert.Equal(t, true, ok)
		assert.Equal(t, b, bytes[:])
	}

}

func TestBPlusTree_Delete(t *testing.T) {
	tree := NewBPlusTree("duckweed", 8)
	rand.Seed(time.Now().Unix())
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.bf.Flush()
	for _, v := range rand.Perm(114514) {
		tree.Delete(v)
	}
	tree.bf.Flush()
	for _, v := range rand.Perm(114514) {
		b, ok := tree.Get(v)
		assert.Equal(t, false, ok)
		assert.Equal(t, b, nil)
	}

}
