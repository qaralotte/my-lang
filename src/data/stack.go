package data

type Element interface{}

type Stack struct {
	list []Element
}

func NewStack() *Stack {
	return &Stack{
		list: make([]Element, 0),
	}
}

func (s *Stack) Len() int {
	return len(s.list)
}

func (s *Stack) IsEmpty() bool {
	return len(s.list) == 0
}

func (s *Stack) Push(x interface{}) {
	s.list = append(s.list, x)
}

func (s *Stack) Pop() Element {
	if s.IsEmpty() {
		return nil
	}
	ret := s.list[len(s.list)-1]
	s.list = s.list[:len(s.list)-1]
	return ret
}

func (s *Stack) Top() Element {
	if s.IsEmpty() {
		return nil
	}
	return s.list[len(s.list)-1]
}

func (s *Stack) Clear() {
	if s.IsEmpty() {
	}
	for i := 0; i < s.Len(); i++ {
		s.list[i] = nil
	}
	s.list = make([]Element, 0)
}

func (s *Stack) Bottom() Element {
	return s.list[0]
}

func (s *Stack) Copy() *Stack {
	newList := make([]Element, s.Len())
	copy(newList, s.list)
	return &Stack{
		list: newList,
	}
}

func (s *Stack) Load(newStack *Stack) {
	newList := make([]Element, newStack.Len())
	copy(newList, newStack.list)
	s.list = newList
}
