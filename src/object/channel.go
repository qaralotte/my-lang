package object

type Channel struct {
	Next []Object
}

func (c *Channel) obj() {}

func NewChannel(next []Object) *Channel {
	return &Channel{
		Next: next,
	}
}
