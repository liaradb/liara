package indexed

import "github.com/liaradb/liaradb/encoder/base"

type baseID = base.ID

type TestID struct {
	baseID
}

func NewTestID() TestID {
	return TestID{
		base.NewID(),
	}
}

func NewTestIDFromString(value string) (TestID, error) {
	if id, err := base.NewIDFromString(value); err != nil {
		return TestID{}, err
	} else {
		return TestID{id}, nil
	}
}

const TestIDSize = base.BaseIDSize
