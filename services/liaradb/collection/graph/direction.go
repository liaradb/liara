package graph

type Direction byte

const (
	BiDirection      Direction = 0
	ForwardDirection Direction = 1
	ReverseDirection Direction = 2
)

const DirectionSize = 1

func (d *Direction) Size() int { return DirectionSize }
