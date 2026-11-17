package indexed

import "github.com/liaradb/liaradb/encoder/base"

type baseString = base.String

type TestName struct {
	baseString
}

func NewTestName(s string) TestName {
	return TestName{base.String(s)}
}
