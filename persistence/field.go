package persistence

import "reflect"

type FieldType = int

const (
	IntField FieldType = iota
	LongField
	FloatField
	DoubleField
	BoolField
	StringField // 长度可能得特殊处理 先定为128算了 可以通过注解操作
)

func FieldTypeLength(tpe FieldType) int {
	switch tpe {
	case IntField:
		return 8
	case LongField:
		return 8
	case FloatField:
		return 8
	case DoubleField:
		return 8
	case BoolField:
		return 1
	case StringField:
		return 128
	default:
		panic("Unknown Error")
	}
}

func ReflectTypeToFieldType(tpe reflect.Type) FieldType {
	fieldType := -1
	switch tpe.Name() {
	case "int", "int64":
		fieldType = LongField
	case "int32":
		fieldType = IntField
	case "float":
		fieldType = FloatField
	case "double":
		fieldType = DoubleField
	case "bool":
		fieldType = BoolField
	case "string":
		fieldType = StringField
	}
	return fieldType
}

type Field struct {
	Name string
	Tpe  FieldType
}

func NewField(name string, tpe FieldType) *Field {
	return &Field{Name: name, Tpe: tpe}
}

func injectDataToField(field reflect.Value, fieldType FieldType, data any) {
	switch fieldType {
	case IntField:
		field.SetInt(data.(int64))
	case LongField:
		field.SetInt(data.(int64))
	case FloatField:
		field.SetFloat(data.(float64))
	case DoubleField:
		field.SetFloat(data.(float64))
	case BoolField:
		field.SetBool(data.(bool))
	case StringField:
		field.SetString(data.(string))
	default:
		panic("Unknown Error")
	}
}
