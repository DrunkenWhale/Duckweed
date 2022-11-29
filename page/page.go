package page

type Page interface {
	ToBytes() []byte
	FromBytes() *Page
}
