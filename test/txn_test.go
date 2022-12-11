package test

import (
	"Duckweed/databox"
	"Duckweed/index"
	"github.com/go-playground/assert/v2"
	"math/rand"
	"testing"
	"time"
)

func TestBPlusTreeRollback1(t *testing.T) {
	rand.Seed(time.Now().Unix())
	tree := index.NewBPlusTree("duckweed", 8)
	tree.StartTransaction()
	for i, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(i, bytes[:])
	}
	tree.Rollback()
	for _, v := range rand.Perm(114514) {
		bs, flag := tree.Get(v)
		assert.IsEqual(bs, nil)
		assert.Equal(t, flag, false)
	}
}

func TestBPlusTreeRollback2(t *testing.T) {
	rand.Seed(time.Now().Unix())
	tree := index.NewBPlusTree("duckweed", 8)
	tree.StartTransaction()
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.Commit()
	bss := databox.IntToBytes(int64(1919810))
	tree.Put(114, bss[:])
	{
		bs, flag := tree.Get(114)
		assert.Equal(t, true, flag)
		b := [8]byte{}
		copy(b[:], bs)
		i := int(databox.BytesToInt(b))
		assert.Equal(t, 1919810, i)
	}
	tree.Rollback()
	{
		bs, flag := tree.Get(114)
		assert.Equal(t, true, flag)
		b := [8]byte{}
		copy(b[:], bs)
		i := int(databox.BytesToInt(b))
		assert.Equal(t, 114, i)
	}
}

func TestBPlusTreeCommit(t *testing.T) {
	rand.Seed(time.Now().Unix())
	tree := index.NewBPlusTree("duckweed", 8)
	tree.StartTransaction()
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.Commit()
	for _, v := range rand.Perm(114514) {
		bs, flag := tree.Get(v)
		assert.Equal(t, true, flag)
		b := [8]byte{}
		copy(b[:], bs)
		i := int(databox.BytesToInt(b))
		assert.Equal(t, v, i)
	}
}

func TestBPlusTreeRollbackAndCommit(t *testing.T) {
	rand.Seed(time.Now().Unix())
	tree := index.NewBPlusTree("duckweed", 8)
	tree.StartTransaction()
	for i, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(i, bytes[:])
	}
	tree.Rollback()
	for _, v := range rand.Perm(114514) {
		bs, flag := tree.Get(v)
		assert.IsEqual(bs, nil)
		assert.Equal(t, flag, false)
	}
	tree.StartTransaction()
	for _, v := range rand.Perm(114514) {
		bytes := databox.IntToBytes(int64(v))
		tree.Put(v, bytes[:])
	}
	tree.Commit()
	for _, v := range rand.Perm(114514) {
		bs, flag := tree.Get(v)
		assert.Equal(t, true, flag)
		b := [8]byte{}
		copy(b[:], bs)
		i := int(databox.BytesToInt(b))
		assert.Equal(t, v, i)
	}
}
