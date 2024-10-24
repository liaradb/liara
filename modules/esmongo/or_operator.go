package esmongo

import "go.mongodb.org/mongo-driver/bson"

func Or(elements ...elementD) or {
	return or{
		elements: elements,
	}
}

type or struct {
	elements []elementD
}

func (o *or) build() bson.E {
	elements := make(bson.A, len(o.elements))
	for _, element := range o.elements {
		elements = append(elements, element.build())
	}
	return bson.E{
		Key:   string(OperatorOr),
		Value: elements}
}
