package variable

// 结构示意图
// [nil, V2, V3]
// 		   	  |
// 			 [Parent, V1, V2]
// 					   |
// 					  [Parent, V1]

type Variable struct {
	Type
	Name string
	Next []*Variable
}

func NewVariable(name string) *Variable {
	return &Variable{
		Type: NONE,
		Name: name,
		Next: nil,
	}
}

func FindVariable(vas []*Variable, name string) *Variable {
	// 如果表里没有任何变量，直接返回 nil
	if len(vas) == 0 {
		return nil
	}

	// 从后往前遍历
	for i := len(vas) - 1; i > 0; i-- {
		va := vas[i]
		if va.Name == name {
			return va
		}
	}

	// 如果没有找到的话就往上一层查找
	if vas[0] != nil {
		return FindVariable(vas[0].Next, name)
	}
	return nil
}
