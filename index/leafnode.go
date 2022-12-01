package index

import (
	"Duckweed/buffer"
	"Duckweed/databox"
	"Duckweed/page"
)

// leaf node in disk
//
//
// 	 	  1 byte      8 byte	   8 byte     8 byte   8 byte
// |head|IsLeafNode|rightSibling|maxKVNumber|kvNumber|ridLength|
// |body|slot(int)|key|key|key...|=>
//							 *****
// 		  			  <=|rid([]byte])|rid([]byte])|rid([]byte])|
//			   			ridLength byte  ridLength byte

type LeafNode struct {
	bf           buffer.BufferPool
	maxKVNumber  int
	ridLength    int
	page         *page.Page
	keys         []int // 键后续可能会扩展(多种类型) 但我要想先做个int的试试
	rids         [][]byte
	rightSibling int
}

func (node *LeafNode) GetPage() *page.Page {
	return node.page
}

func (node *LeafNode) FetchNode(pageID int) BPlusNode {
	p := node.bf.GetPage(pageID)
	n := FromPage(p, node.bf)
	return n
}

func (node *LeafNode) ToBytes() []byte {
	header := make([]byte, 1)
	header[0] = LeafNodeFlag
	if len(node.keys) != len(node.rids) {
		panic("Keys Number Should equal Rids Number (o゜▽゜)o☆")
	}
	rightSiblingBytes := databox.IntToBytes(int64(node.rightSibling))
	maxKeysNumberBytes := databox.IntToBytes(int64(node.maxKVNumber))
	keysNumberBytes := databox.IntToBytes(int64(len(node.keys)))
	ridLengthBytes := databox.IntToBytes(int64(node.ridLength))
	header = append(header, rightSiblingBytes[:]...)
	header = append(header, maxKeysNumberBytes[:]...)
	header = append(header, keysNumberBytes[:]...)
	header = append(header, ridLengthBytes[:]...)
	keysBytes := make([]byte, 8*len(node.keys))
	for i := 0; i < len(node.keys); i++ {
		b := databox.IntToBytes(int64(node.keys[i]))
		copy(keysBytes[i*8:(i+1)*8], b[:])
	}
	ridsBytes := make([]byte, node.ridLength*len(node.rids))
	for i := 0; i < len(node.rids); i++ {
		b := node.rids[len(node.rids)-1-i]
		copy(ridsBytes[i*node.ridLength:(i+1)*node.ridLength], b[:])
	}
	blankBytesSize := page.PageSize - (len(header) + len(keysBytes) + len(ridsBytes))
	blankBytes := make([]byte, blankBytesSize)
	bytes := append(header, keysBytes...)
	bytes = append(bytes, blankBytes...)
	bytes = append(bytes, ridsBytes...)
	return bytes
}

func LeafNodeFromPage(p *page.Page, bf buffer.BufferPool) *LeafNode {
	bytes := p.GetBytes()
	rightSiblingBytes := [8]byte{}
	copy(rightSiblingBytes[:], bytes[1:9])
	rightSibling := int(databox.BytesToInt(rightSiblingBytes))
	maxKeysNumberBytes := [8]byte{}
	copy(maxKeysNumberBytes[:], bytes[9:17])
	maxKeysNumber := int(databox.BytesToInt(maxKeysNumberBytes))
	kvNumberBytes := [8]byte{}
	copy(kvNumberBytes[:], bytes[17:25])
	kvNumber := int(databox.BytesToInt(kvNumberBytes))
	ridLengthBytes := [8]byte{}
	copy(ridLengthBytes[:], bytes[9:17])
	ridLength := int(databox.BytesToInt(ridLengthBytes))

	headerOffset := 4*8 + 1

	keys := make([]int, kvNumber)
	rids := make([][]byte, kvNumber)

	for i := 0; i < kvNumber; i++ {
		b := [8]byte{}
		copy(b[:], bytes[headerOffset+i*8:headerOffset+(i+1)*8])
		num := databox.BytesToInt(b)
		keys[i] = int(num)
	}
	for i := 0; i < kvNumber; i++ {
		b := make([]byte, ridLength)
		copy(b[:], bytes[page.PageSize-(i+1)*ridLength:page.PageSize-i*ridLength])
		rids[i] = b
	}
	node := &LeafNode{
		bf:           bf,
		maxKVNumber:  maxKeysNumber,
		ridLength:    ridLength,
		page:         p,
		keys:         keys,
		rids:         rids,
		rightSibling: rightSibling,
	}
	return node
}
