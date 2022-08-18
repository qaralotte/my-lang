package ast

type (
	Object interface {
		obj()
	}

	// ObjectList 对象表 (指针数组指针)
	ObjectList struct {
		Objects *[]Object
	}

	// Variable 变量
	Variable struct {
		Type
		Name string
	}

	// Function 方法
	Function struct {
		Variable             // 基本信息
		Args     []Object    // 方法参数
		Objects  *ObjectList // 方法局部对象表
		Stmts    []Stmt      // 方法内语句
	}

	// Channel 通道 (建立两个对象表的联系)
	Channel struct {
		Next *ObjectList // 对象表
	}
)

func (*Variable) obj() {}
func (*Function) obj() {}
func (*Channel) obj()  {}

func NewObjectList(previous *ObjectList) *ObjectList {
	objs := make([]Object, 1)
	objs[0] = NewChannel(previous)
	return &ObjectList{
		Objects: &objs,
	}
}

func NewVariable(name string) *Variable {
	return &Variable{
		Type: NONE,
		Name: name,
	}
}

func NewFunction(name string, args []Object, parent *ObjectList) *Function {
	return &Function{
		Variable: Variable{
			Type: VOID,
			Name: name,
		},
		Args:    args,
		Objects: NewObjectList(parent),
		Stmts:   nil,
	}
}

func NewChannel(next *ObjectList) *Channel {
	return &Channel{
		Next: next,
	}
}

// 获取对象表指定下标的对象
func (objs *ObjectList) get(index int) Object {
	return (*objs.Objects)[index]
}

// 往对象表里插入新对象
func (objs *ObjectList) add(object Object) {
	*objs.Objects = append(*objs.Objects, object)
}

// 倒序查找当前对象表，如果没有找到对象，则往上一层继续找
func (objs *ObjectList) findObject(name string) Object {
	// 如果表里没有任何变量，直接返回 nil
	if len(*objs.Objects) == 0 {
		return nil
	}

	// 从后往前遍历
	for i := len(*objs.Objects) - 1; i > 0; i-- {
		// 根据对象类型来判断
		if name == getObjectName(objs.get(i)) {
			return objs.get(i)
		}
	}

	// 如果没有找到的话就往上一层查找
	channel := objs.get(0).(*Channel)
	if channel.Next != nil {
		return channel.Next.findObject(name)
	}

	return nil

}

// 获取对象名称
func getObjectName(obj Object) string {
	switch obj.(type) {
	case *Variable:
		return obj.(*Variable).Name
	case *Function:
		return obj.(*Function).Name
	}
	panic("错误: 非法的对象名称获取")
}

// 获取对象类型
func getObjectType(obj Object) Type {
	switch obj.(type) {
	case *Variable:
		return obj.(*Variable).Type
	case *Function:
		return obj.(*Function).Type
	}
	panic("错误: 非法的对象类型获取")
}
