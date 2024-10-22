package esmongo

type Model[T any] struct {
	ModelData `bson:"inline"`
	Value     T `bson:"inline"`
}

type ModelData struct {
	ID      string `bson:"_id"`
	Version int    `bson:"version"`
}
