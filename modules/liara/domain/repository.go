package domain

import (
	"context"
	"iter"
)

type (
	CommandRepository[I id, T Entity[I, Version]] interface {
		GetByID(context.Context, I) (*T, error)
		Insert(context.Context, *T) error
		Replace(context.Context, *T) error
	}

	QueryRepository[I ~string, T any, F any] interface {
		GetByID(context.Context, I) (*T, error)
		Search(context.Context, Page, F) iter.Seq2[*T, error]
	}

	// Entity[I ~string] interface {
	// 	ID() I
	// 	Version() Version
	// }

	Page struct {
		Offset int
		Limit  int
	}
)
