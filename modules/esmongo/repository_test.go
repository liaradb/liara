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

type BookVersion int

func (bv BookVersion) Value() int { return int(bv) }

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

func (b bookMapper) FromModel(d Metadata, m BookModel) Book {
	return Book{
		id:      BookID(d.ID),
		version: BookVersion(d.Version),
		title:   m.Title,
	}
}

func (bookMapper) ToModel(b Book) BookModel {
	return BookModel{
		Title: b.title,
	}
}
