package ast

import "fmt"

type Type int

const (
	VOLATILE Type = iota
	INT
	FLOAT
	STRING
	BOOL
)

func TypeString(t Type) string {
	switch t {
	case VOLATILE:
		return "volatile"
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	}
	panic(fmt.Sprintf("错误: 未知类型 %d", t))
}

// Merge 如果可以的话，类型应该向下合并
func tryMerge(typ1 *Type, typ2 *Type) {

	// 如果两个类型其中一个是float类型，而另一个是int类型，则将int转换成float
	if *typ1 == INT && *typ2 == FLOAT {
		*typ1 = *typ2
	}
	if *typ1 == FLOAT && *typ2 == INT {
		*typ2 = *typ1
	}

}

// 判断两个类型之间是否可以计算
func canCalc(op int, typ1 Type, typ2 Type) bool {

	// INT FLOAT STRING BOOL
	// 从小到大排序，为了方便后续的判断
	if typ1 > typ2 {
		typ3 := typ1
		typ1 = typ2
		typ2 = typ3
	}

	// 如果两个类型中有一个是 none，则到运行时才可以判定是否可以计算，编译期暂时通过
	if typ1 == VOLATILE || typ2 == VOLATILE {
		return true
	}

	// 如果两个类型相等且为int或者float，可以计算任意运算符
	if typ1 == typ2 && (typ1 == INT || typ1 == FLOAT) {
		return true
	}

	// 如果两个类型相等且为string，则只能计算+运算符
	if typ1 == typ2 && typ1 == STRING && op == ADD {
		return true
	}

	// 如果两个类型其中一个是string，另一个是int，则只能计算*运算符
	if typ1 == INT && typ2 == STRING && op == MUL {
		return true
	}

	return false
}
