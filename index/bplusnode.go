package index

type BPlusNode interface {
	ToBytes() []byte
}

func FromBytes(bytes []byte) BPlusNode {
	if bytes[0] == 1 {
		return IndexNodeFromBytes(bytes)
	} else if bytes[0] == 2 {
		return LeafNodeFromBytes(bytes)
	} else {
		panic("Illegal Page!")
	}
}
