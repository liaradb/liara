package esmongo

import (
	"context"
	"testing"
)

func TestRepository(t *testing.T) {
	t.Skip()
	r := NewRepository[Book, BookID, BookModel](nil, bookMapper{})
	r.Insert(context.Background(), "1", &Book{})
}

type BookID string

func (b BookID) String() string { return string(b) }

type Book struct {
	id      BookID
	version int
	title   string
}

type BookModel struct {
	Title string `bson:"title"`
}

type bookMapper struct{}

func (b bookMapper) FromModel(m *Model[BookModel]) *Book {
	return &Book{}
}

func (bookMapper) ToModel(b *Book) *Model[BookModel] {
	return &Model[BookModel]{
		ModelData: ModelData{
			ID:      b.id.String(),
			Version: b.version,
		},
		Value: BookModel{
			Title: b.title,
		},
	}
}
