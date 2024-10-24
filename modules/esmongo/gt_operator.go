package esmongo

import "go.mongodb.org/mongo-driver/bson"

func GT[V ~int | ~float32](
	key string,
	value V,
) gt[V] {
	return gt[V]{
		key:   key,
		value: value,
	}

}

type gt[V ~int | ~float32] struct {
	key   string
	value V
}

func (g gt[V]) build() bson.D {
	return bson.D{{
		Key: g.key,
		Value: bson.D{{
			Key:   string(OperatorGT),
			Value: g.value}}}}
}
