package command

import (
	"errors"

	"github.com/liaradb/liaradb/domain/value"
)

type AppendEvent struct {
	TenantID    value.TenantID
	Options     AppendOptions
	PartitionID value.PartitionID
	Events      []EventOptions
}

func (ae *AppendEvent) Valid() error {
	if len(ae.Events) == 0 {
		return value.ErrNoEvents
	}

	errs := make([]error, 0)
	for _, em := range ae.Events {
		if err := em.Valid(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
