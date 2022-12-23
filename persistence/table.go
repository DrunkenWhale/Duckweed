package persistence

import (
	"Duckweed/index"
	"reflect"
	"strings"
)

type Table[T any] struct {
	name   string
	fields []*Field
	tpe    reflect.Type
	tree   *index.BPlusTree
}

func NewTable[T any](name string, model *T) {
	table := Table[T]{
		name: name,
	}
	tpe := reflect.TypeOf(model)
	table.tpe = tpe
	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		// 如果不是public的 则加入fields中
		if field.IsExported() {
			table.fields = append(
				table.fields,
				// 字段用小写表示
				NewField(
					strings.ToLower(field.Name),
					FieldTypeLength(
						ReflectTypeToFieldType(field.Type),
					),
				),
			)
		}
	}
	length := 0
	for _, field := range table.fields {
		length += FieldTypeLength(field.Tpe)
	}
	tree := index.NewBPlusTree(name, length)
	table.tree = tree
}

func (t *Table[T]) Get(key int) (T, bool) {
	bytes, flag := t.tree.Get(key)
	if !flag {
		return nil, false
	}

}

func (t *Table[T]) generateStructFromBytes() {
	res := reflect.New(t.tpe)
	ele := res.Elem()
	for i := 0; i < ele.NumField(); i++ {
		injectDataToField(ele.Field(i), t.fields[i].Tpe, nil) // 各种类型的转换我还没写好... 刚想起来就写了个int的
	}
}

func (t *Table[T]) Put(key int, value T) {

}

func (t *Table[T]) Update(key int, value T) {

}

func (t *Table[T]) Delete(key int) {

}
