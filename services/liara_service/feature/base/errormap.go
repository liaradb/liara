package base

import (
	"context"
	"errors"
	"net/http"

	"github.com/cardboardrobots/baseerror"
	"google.golang.org/grpc/codes"
)

func GetStatusCodeGRPC(err error) codes.Code {
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	if errors.Is(err, baseerror.ErrAlreadyExists) {
		return codes.AlreadyExists
	}
	if errors.Is(err, baseerror.ErrFailedPrecondition) {
		return codes.FailedPrecondition
	}
	if errors.Is(err, baseerror.ErrInvalidArgument) {
		return codes.InvalidArgument
	}
	if errors.Is(err, baseerror.ErrNotFound) {
		return codes.NotFound
	}
	if errors.Is(err, baseerror.ErrOutOfRange) {
		return codes.OutOfRange
	}
	if errors.Is(err, baseerror.ErrPermissionDenied) {
		return codes.PermissionDenied
	}
	if errors.Is(err, baseerror.ErrUnauthenticated) {
		return codes.Unauthenticated
	}
	return codes.Unknown
}

func GetStatusCodeHTTP(err error) int {
	if errors.Is(err, context.Canceled) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, baseerror.ErrAlreadyExists) {
		return http.StatusNotModified
	}
	if errors.Is(err, baseerror.ErrFailedPrecondition) {
		return http.StatusPreconditionFailed
	}
	if errors.Is(err, baseerror.ErrInvalidArgument) {
		return http.StatusBadRequest
	}
	if errors.Is(err, baseerror.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, baseerror.ErrOutOfRange) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, baseerror.ErrPermissionDenied) {
		return http.StatusForbidden
	}
	if errors.Is(err, baseerror.ErrUnauthenticated) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
