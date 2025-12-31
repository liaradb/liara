package entity

import "github.com/liaradb/liaradb/domain/value"

type Row struct {
	Metadata Metadata
	Data     value.Data
}
