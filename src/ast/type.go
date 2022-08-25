package ast

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Type int

const (
	INVALID Type = iota
	INT
	FLOAT
	STRING
	BOOL
)

var types = map[string]Type{
	"int64":   INT,
	"float64": FLOAT,
	"string":  STRING,
	"bool":    BOOL,
}

func TypeString(t Type) string {
	switch t {
	case INVALID:
		return "invalid"
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	}
	panic(fmt.Sprintf("错误: 未知类型 %v", t))
}

// GetType 反射并转换成规定的类型
func GetType(val interface{}) Type {
	return types[reflect.TypeOf(val).String()]
}

// ProcessType 根据两个type与运算符进行加工
// 1. type按照顺序摆放: INT FLOAT STRING BOOL
// 2. 类型自动隐式转换
func ProcessType(typ1 Type, typ2 Type) (Type, Type) {

	// 从小到大排序
	if typ1 > typ2 {
		typ3 := typ1
		typ1 = typ2
		typ2 = typ3
	}

	// 类型自动隐式转换
	if typ1 == INT && typ2 == FLOAT {
		return FLOAT, FLOAT
	}

	return typ1, typ2
}

// SortByType 根据类型排序值
func SortByType(v1 interface{}, v2 interface{}) (interface{}, interface{}) {
	t1, t2 := GetType(v1), GetType(v2)
	if t1 > t2 {
		return v2, v1
	}
	return v1, v2
}

// SameType 判断两个类型是否等于类型 p
func SameType(typ1 Type, typ2 Type, p Type) bool {
	return typ1 == typ2 && typ1 == p
}

// Int64ToInt 将 int64 转换成 int 类型
func Int64ToInt(i int64) int {
	return *(*int)(unsafe.Pointer(&i))
}
