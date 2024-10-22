package esmongo

type Entity[I EntityID] interface {
	ID() I
	Version() Version
}

type EntityID interface {
	String() string
}

type Version interface {
	Value() int
}

type Model[T any] struct {
	ModelData `bson:"inline"`
	Value     T `bson:"inline"`
}

type ModelData struct {
	ID      string `bson:"_id"`
	Version int    `bson:"version"`
}

type Mapper[I EntityID, E Entity[I], M any] interface {
	FromModel(*M) *E
	ToModel(*E) *M
}
