package testutil

import (
	"errors"
	"reflect"
	"testing"
)

func Equals[T comparable](t *testing.T, a T, b T, name string) bool {
	t.Helper()

	equals := a == b
	if !equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, a)
	}

	return equals
}

func EqualsArray[T comparable](t *testing.T, a []T, b []T, name string) bool {
	t.Helper()

	equals := reflect.DeepEqual(a, b)
	if !equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, a)
	}

	return equals
}

func Getter[T comparable](t *testing.T, a func() T, b T, name string) bool {
	t.Helper()

	value := a()
	equals := value == b
	if !equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, value)
	}

	return equals
}

func GetterArray[T comparable](t *testing.T, a func() []T, b []T, name string) bool {
	t.Helper()

	value := a()
	equals := reflect.DeepEqual(value, b)
	if !equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, value)
	}

	return equals
}

func True(t *testing.T, value bool, message string) bool {
	t.Helper()

	if !value {
		t.Error(message)
	}

	return value
}

func False(t *testing.T, value bool, message string) bool {
	t.Helper()

	if value {
		t.Error(message)
	}

	return !value
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func MustError(t *testing.T, err error, message string) {
	t.Helper()

	if err == nil {
		t.Fatal(message)
	}
}

func ErrorIs(t *testing.T, err error, target error, message string) bool {
	t.Helper()

	equals := errors.Is(err, target)
	if !equals {
		t.Error(message)
	}

	return equals
}

func NotEquals[T comparable](t *testing.T, a T, b T, name string) bool {
	t.Helper()

	notEquals := a != b
	if !notEquals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, a)
	}

	return notEquals
}

func NotEqualsEval[T comparable](t *testing.T, a func() T, b T, name string) bool {
	t.Helper()

	value := a()
	equals := value == b
	if equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, value)
	}

	return !equals
}
