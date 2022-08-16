package object

// 结构示意图
// [nil, V1, F1]
// 		   	  |
// 			 [Channel, Arg1, Arg2, V1, V2, V3]

type Function struct {
	Variable
	Args []Object
	Next *Channel
}

func (fn *Function) obj() {}

func NewFunction(name string, args []Object, parent []Object) *Function {
	return &Function{
		Variable: Variable{
			Type: NONE,
			Name: name,
		},
		Args: args,
		Next: NewChannel(InitObjects(parent)),
	}
}
