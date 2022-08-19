package ast

import "fmt"

type Type int

const (
	NONE Type = iota
	INT
	FLOAT
	STRING
	BOOL
	VOID
)

func TypeString(t Type) string {
	switch t {
	case NONE:
		return "none"
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	case VOID:
		return "void"
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

// CanCalc 两个类型之间是否可以计算
func canCalc(typ1 Type, typ2 Type) bool {

	// 如果两个类型中有一个是 none，则到运行时才可以判定是否可以计算，编译期暂时通过
	if typ1 == NONE || typ2 == NONE {
		return true
	}

	// 如果两个类型是int或者float，可以计算
	if typ1 == typ2 && (typ1 == INT || typ1 == FLOAT) {
		return true
	}

	return false
}
