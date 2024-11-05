package esmongo

import (
	"context"
	"testing"
)

func TestRepository(t *testing.T) {
	t.Skip()
	r := NewRepository(nil, bookMapper{})
	r.Insert(context.Background(), Book{})
}

type BookID string

func (b BookID) String() string { return string(b) }

type Book struct {
	id      BookID
	version Version
	title   string
}

func (b Book) ID() BookID       { return b.id }
func (b Book) Version() Version { return b.version }
func (b Book) Events() []Event  { return nil }

type BookModel struct {
	Title string `bson:"title"`
}

type bookMapper struct{}

func (b bookMapper) FromModel(m Model[BookModel]) Book {
	return Book{}
}

func (bookMapper) ToModel(b Book) BookModel {
	return BookModel{
		Title: b.title,
	}
}
