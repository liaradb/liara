package esmongo

import "github.com/cardboardrobots/baseerror"

var (
	ErrNotFound = baseerror.ErrNotFound
	ErrNoMatch  = baseerror.ErrInvalidArgument.Wrap("no match")
)
