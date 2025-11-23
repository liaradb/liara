package raw

type Orderer[T any] interface {
	Equaler[T]
	Greaterer[T]
	GreaterOrEqualer[T]
	Lesser[T]
	LessOrEqualer[T]
}

type Equaler[T any] interface {
	Equal(T) bool
}

type Greaterer[T any] interface {
	Greater(T) bool
}

type GreaterOrEqualer[T any] interface {
	GreaterOrEqual(T) bool
}

type Lesser[T any] interface {
	Less(T) bool
}

type LessOrEqualer[T any] interface {
	LessOrEqual(T) bool
}
