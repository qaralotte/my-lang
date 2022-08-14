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
