package esmongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

type FilterBuilder map[string]any

func Filter() FilterBuilder {
	return FilterBuilder{}
}

func (fb FilterBuilder) If(condition bool, a func(fb FilterBuilder)) FilterBuilder {
	if condition {
		a(fb)
	}
	return fb
}

func (fb FilterBuilder) IfElse(condition bool, a func(fb FilterBuilder), b func(fb FilterBuilder)) FilterBuilder {
	if condition {
		a(fb)
	} else {
		b(fb)
	}
	return fb
}

func (fb FilterBuilder) Property(key string, value any) FilterBuilder {
	fb[key] = value
	return fb
}

func (fb FilterBuilder) PropertyIf(condition bool, key string, value any) FilterBuilder {
	if condition {
		fb[key] = value
	}
	return fb
}

func (fb FilterBuilder) PropertyIfElse(condition bool, key string, a any, b any) FilterBuilder {
	if condition {
		fb[key] = a
	} else {
		fb[key] = b
	}
	return fb
}

func (fb FilterBuilder) Build() bson.M {
	return bson.M(fb)
}
