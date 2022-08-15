package variable

type Type int

const (
	NONE Type = iota
	NUMBER
	STRING
)

func TypeString(t Type) string {
	switch t {
	case NONE:
		return "none"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	}
	return "undef"
}

// Merge 如果可以的话，类型应该向下合并
func Merge(typ1 *Type, typ2 *Type) {

	// 如果两个类型中有一个是none类型，则none类型自动变成另一个类型
	if *typ1 == NONE {
		typ1 = typ2
	}
	if *typ2 == NONE {
		typ2 = typ1
	}

}

// CanCalc 两个类型之间是否可以计算
func CanCalc(typ1 Type, typ2 Type) bool {
	// 两个数字当然可以计算
	if typ1 == NUMBER && typ2 == NUMBER {
		return true
	}
	return false
}
