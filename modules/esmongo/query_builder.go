package esmongo

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type QueryBuilder struct {
	Filter FilterBuilder
	Sort   SortBuilder
}

func Query(
	filter FilterBuilder,
	sort SortBuilder,
) QueryBuilder {
	return QueryBuilder{
		Filter: filter,
		Sort:   sort,
	}
}

type QueryType interface {
	Filter() map[string]any
	Sort() map[string]Sort
	Projection() map[string]Projection
	Offset() int
	Limit() int
}

func Field(names ...string) string {
	return strings.Join(names, ".")
}

type elementA interface {
	build() bson.A
}

type elementD interface {
	build() bson.D
}

type elementE interface {
	build() bson.E
}

type elementM interface {
	build() bson.M
}

func a(values ...any) bson.A {
	return bson.A(values)
}

func d(elements ...bson.E) bson.D {
	return bson.D(elements)
}

func e(key string, value any) bson.E {
	return bson.E{
		Key:   key,
		Value: value,
	}
}
