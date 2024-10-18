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

func Sort() *SortBuilder {
	return &SortBuilder{}
}

func (sb *SortBuilder) SetLimit(limit int) *SortBuilder {
	sb.limit = limit
	return sb
}

func (sb *SortBuilder) SetSkip(skip int) *SortBuilder {
	sb.skip = skip
	return sb
}

func (sb *SortBuilder) Asc(key string) *SortBuilder {
	sb.data = append(sb.data, primitive.E{
		Key:   key,
		Value: 1,
	})
	return sb
}

func (sb *SortBuilder) Desc(key string) *SortBuilder {
	sb.data = append(sb.data, primitive.E{
		Key:   key,
		Value: -1,
	})
	return sb
}

func (sb *SortBuilder) IfAsc(key string, test bool) *SortBuilder {
	if test {
		sb.data = append(sb.data, primitive.E{
			Key:   key,
			Value: 1,
		})
	}
	return sb
}

func (sb *SortBuilder) IfDesc(key string, test bool) *SortBuilder {
	if test {
		sb.data = append(sb.data, primitive.E{
			Key:   key,
			Value: -1,
		})
	}
	return sb
}

func (sb *SortBuilder) IfAscElseDesc(key string, test bool) *SortBuilder {
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
	return sb
}

func (sb *SortBuilder) Build() *options.FindOptions {
	if sb == nil {
		return options.Find()
	}

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
