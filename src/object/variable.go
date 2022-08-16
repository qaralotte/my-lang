package object

type Variable struct {
	Type
	Name string
}

func (*Variable) obj() {}

func NewVariable(name string) *Variable {
	return &Variable{
		Type: NONE,
		Name: name,
	}
}
