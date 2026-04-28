package iterator

import (
	"container/list"
	"slices"
	"testing"
)

func TestForward(t *testing.T) {
	t.Parallel()

	l := list.New()
	values := []string{"a", "b", "c"}

	for _, value := range values {
		l.PushBack(value)
	}

	for range Forward[string](l) {
		break
	}

	var result []string
	for value := range Forward[string](l) {
		result = append(result, value)
	}

	if !slices.Equal(result, values) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()

	l := list.New()
	values := []string{"a", "b", "c"}

	for _, value := range values {
		l.PushBack(value)
	}

	for range Reverse[string](l) {
		break
	}

	var result []string
	for value := range Reverse[string](l) {
		result = append(result, value)
	}

	slices.Reverse(values)
	if !slices.Equal(result, values) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestItem(t *testing.T) {
	t.Parallel()

	l := list.New()

	if _, ok := Item[string](l.Front()); ok {
		t.Error("should not return value")
	}

	want := "a"
	l.PushBack(want)

	if v, ok := Item[string](l.Front()); !ok {
		t.Error("should return value")
	} else if v != want {
		t.Errorf("incorrect value: %v, expected: %v", v, want)
	}
}
