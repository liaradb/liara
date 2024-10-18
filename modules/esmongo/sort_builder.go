package esmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SortBuilder struct {
	data  []primitive.E
	limit int
	skip  int
}

func (sb *SortBuilder) SetLimit(limit int) {
	sb.limit = limit
}

func (sb *SortBuilder) SetSkip(skip int) {
	sb.skip = skip
}

func (sb *SortBuilder) Asc(key string) {
	sb.data = append(sb.data, primitive.E{
		Key:   key,
		Value: 1,
	})
}

func (sb *SortBuilder) Desc(key string) {
	sb.data = append(sb.data, primitive.E{
		Key:   key,
		Value: -1,
	})
}

func (sb *SortBuilder) IfAsc(key string, test bool) {
	if test {
		sb.data = append(sb.data, primitive.E{
			Key:   key,
			Value: 1,
		})
	}
}

func (sb *SortBuilder) IfDesc(key string, test bool) {
	if test {
		sb.data = append(sb.data, primitive.E{
			Key:   key,
			Value: -1,
		})
	}
}

func (sb *SortBuilder) IfAscElseDesc(key string, test bool) {
	if test {
		sb.data = append(sb.data, primitive.E{
			Key:   key,
			Value: 1,
		})
	} else {
		sb.data = append(sb.data, primitive.E{
			Key:   key,
			Value: -1,
		})
	}
}

func (sb *SortBuilder) Build() *options.FindOptions {
	o := options.Find()
	if sb.data != nil {
		o = o.SetSort(sb.data)
	}
	if sb.limit != 0 {
		o = o.SetLimit(int64(sb.limit))
	}
	if sb.skip != 0 {
		o = o.SetSkip(int64(sb.skip))
	}
	return o
}
