package test

import (
	"fmt"
	"reflect"
	"testing"
)

type Student struct {
	Name   string
	age    int64
	Gender bool
}

func TestRef(t *testing.T) {
	NewTable[Student](Student{})
}

func NewTable[T any](s T) {
	tpe := reflect.TypeOf(s)
	l := reflect.New(tpe).Elem()
	l.Field(0).SetString("1")
	l.Field(2).SetBool(false)
	fmt.Println(l)
}
