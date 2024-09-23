package esmongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type SortBuilder struct {
	data []primitive.E
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

func (sb *SortBuilder) Build() primitive.D {
	return primitive.D(sb.data)
}
