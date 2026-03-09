package schema

type Orderer[T Orderer[T]] interface {
	Equaler[T]
	Greaterer[T]
	GreaterOrEqualer[T]
	Lesser[T]
	LessOrEqualer[T]
}

type Equaler[T Equaler[T]] interface {
	Equal(T) bool
}

type Greaterer[T Greaterer[T]] interface {
	Greater(T) bool
}

type GreaterOrEqualer[T GreaterOrEqualer[T]] interface {
	GreaterOrEqual(T) bool
}

type Lesser[T Lesser[T]] interface {
	Less(T) bool
}

type LessOrEqualer[T LessOrEqualer[T]] interface {
	LessOrEqual(T) bool
}
