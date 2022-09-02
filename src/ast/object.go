package ast

import (
	"my-lang/token"
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
		Name  string
		Value interface{}
	}

	// Function 方法
	Function struct {
		Name string
		Args []string      // 局部变量
		Body []token.Token // 内容
	}

	// Channel 通道 (建立两个对象表的联系)
	Channel struct {
		Previous *ObjectList // 上一层对象表
	}
)

func (*Variable) obj() {}
func (*Function) obj() {}
func (*Channel) obj()  {}

func NewObjectList(previous *ObjectList) *ObjectList {
	objs := make([]Object, 1)
	objs[0] = &Channel{
		Previous: previous,
	}
	return &ObjectList{
		Objects: &objs,
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

// 往对象表里插入新对象
func (objs *ObjectList) addBatch(object []Object) {
	*objs.Objects = append(*objs.Objects, object...)
}

func (objs *ObjectList) Len() int {
	return len(*objs.Objects)
}

// FindObject 倒序查找当前对象表，如果没有找到对象，则往上一层继续找
func (objs *ObjectList) FindObject(name string) Object {

	// 从后往前遍历
	for i := objs.Len() - 1; i > 0; i-- {
		// 根据对象类型来判断
		fieldName, isExsit := getObjectField(objs.Get(i), "Name")
		if !isExsit {
			// 如果不存在 Name 字段，则跳过检查
			continue
		}

		if name == fieldName.String() {
			return objs.Get(i)
		}
	}

	// 如果没有找到的话就往上一层查找
	channel := objs.Get(0).(*Channel)
	if channel.Previous != nil {
		return channel.Previous.FindObject(name)
	}

	return nil
}

// Get 获取对象表指定下标的对象
func (objs *ObjectList) Get(index int) Object {
	return (*objs.Objects)[index]
}

// Add 往对象表里插入新对象
func (objs *ObjectList) Add(object Object) {
	*objs.Objects = append(*objs.Objects, object)
}

func (objs *ObjectList) Clear() {
	*objs.Objects = (*objs.Objects)[:1]
}
