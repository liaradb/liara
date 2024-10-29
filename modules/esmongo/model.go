package esmongo

type Entity[I EntityID] interface {
	ID() I
	Version() Version
	Events() []Event
}

type Event interface {
	ID() EventID
	Type() EventType
	EntityID() EntityID
	Version() Version
}

type EntityID interface {
	String() string
}

type Version interface {
	Value() int
}

type EventID interface {
	String() string
}

type EventType interface {
	String() string
}

type Model[T any] struct {
	ModelData `bson:"inline"`
	Value     T `bson:"inline"`
}

type ModelData struct {
	ID      string       `bson:"_id"`
	Version int          `bson:"version"`
	Events  []ModelEvent `bson:"events"`
}

type Mapper[I EntityID, E Entity[I], M any] interface {
	FromModel(*M) *E
	ToModel(*E) *M
}

type ModelEvent struct {
	ID       string `bson:"id"`
	Type     string `bson:"type"`
	EntityID string `bson:"entityId"`
	Version  string `bson:"version"`
	Date     []byte `bson:"data"`
}
