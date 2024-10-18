package esmongo

import (
	"testing"
)

// TODO: Change tests to use public interface

func TestSortBuilder_Asc(t *testing.T) {
	sb := SortBuilder{}
	sb.Asc("name")
	value := sb.data[0]
	if value.Key != "name" || value.Value != 1 {
		t.Error("value is incorrect")
	}
}

func TestSortBuilder_Desc(t *testing.T) {
	sb := SortBuilder{}
	sb.Desc("name")
	value := sb.data[0]
	if value.Key != "name" || value.Value != -1 {
		t.Error("value is incorrect")
	}
}
