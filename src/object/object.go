package object

type (
	Object interface {
		obj()
	}
)

func FindObject(objs []Object, name string) Object {
	// 如果表里没有任何变量，直接返回 nil
	if len(objs) == 0 {
		return nil
	}

	// 从后往前遍历
	for i := len(objs) - 1; i > 0; i-- {
		// 根据对象类型来判断
		if name == GetName(objs[i]) {
			return objs[i]
		}
	}

	// 如果没有找到的话就往上一层查找
	if objs[0] != nil {
		obj := objs[0].(*Channel)
		return FindObject(obj.Next, name)
	}

	return nil
}

func GetName(obj Object) string {
	switch obj.(type) {
	case *Variable:
		return obj.(*Variable).Name
	case *Function:
		return obj.(*Function).Name
	}
	panic("错误: 非法的对象名称获取")
}

func GetType(obj Object) Type {
	switch obj.(type) {
	case *Variable:
		return obj.(*Variable).Type
	case *Function:
		return obj.(*Function).Type
	}
	panic("错误: 非法的对象类型获取")
}

func InitObjects(previous []Object) []Object {
	objs := make([]Object, 10) // todo DEFAULT_SIZE 10
	objs[0] = &Channel{
		Next: previous,
	}
	return objs
}
