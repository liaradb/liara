package domain

type (
	Entity[I id, V version] struct {
		id      I
		version V
	}

	id interface {
		String() string
	}

	version interface {
		Value() int
		Increment()
	}
)

func (e Entity[I, V]) ID() I      { return e.id }
func (e Entity[I, V]) Version() V { return e.version }

func NewEntity[I id, V version](
	id I,
	version V,
) Entity[I, V] {
	return Entity[I, V]{
		id:      id,
		version: version,
	}
}

func (e *Entity[I, V]) Increment() {
	e.version.Increment()
}
