package ast

import (
	"reflect"
)

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
		Name    string
		Objects *ObjectList // 方法局部对象表
		Stmts   []Stmt      // 方法内语句
		Return  Expr        // 返回值
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

func NewFunction(name string, parent *ObjectList) *Function {
	return &Function{
		Name:    name,
		Objects: NewObjectList(parent),
		Stmts:   nil,
		Return:  nil,
	}
}

func NewChannel(next *ObjectList) *Channel {
	return &Channel{
		Next: next,
	}
}

// 获取对象名称
func getObjectField(obj Object, field string) (reflect.Value, bool) {
	f := reflect.ValueOf(obj).Elem().FieldByName(field)
	if f.Kind() == reflect.Invalid {
		return f, false
	}
	return f, true
}

// 获取对象表指定下标的对象
func (objs *ObjectList) get(index int) Object {
	return (*objs.Objects)[index]
}

// 往对象表里插入新对象
func (objs *ObjectList) add(object Object) {
	*objs.Objects = append(*objs.Objects, object)
}

// 往对象表里插入新对象
func (objs *ObjectList) addBatch(object []Object) {
	*objs.Objects = append(*objs.Objects, object...)
}

func (objs *ObjectList) size() int {
	return len(*objs.Objects)
}

// 倒序查找当前对象表，如果没有找到对象，则往上一层继续找
func (objs *ObjectList) findObject(name string) Object {
	// 如果表里没有任何变量，直接返回 nil
	if objs.size() == 0 {
		return nil
	}

	// 从后往前遍历
	for i := objs.size() - 1; i > 0; i-- {
		// 根据对象类型来判断
		fieldName, isExsit := getObjectField(objs.get(i), "Name")
		if !isExsit {
			// 如果不存在 Name 字段，则跳过检查
			continue
		}

		if name == fieldName.String() {
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

// 获取父对象 (外部对象)
func (objs *ObjectList) getParentObject() Object {
	parentObjs := objs.get(0).(*Channel).Next
	if parentObjs == nil {
		return nil
	}
	return parentObjs.get(parentObjs.size() - 1)
}
