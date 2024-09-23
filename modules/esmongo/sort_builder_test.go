package esmongo

import (
	"testing"
)

func TestSortBuilder_Asc(t *testing.T) {
	sb := SortBuilder{}
	sb.Asc("name")
	s := sb.Build()
	value := s[0]
	if value.Key != "name" || value.Value != 1 {
		t.Error("value is incorrect")
	}
}

func TestSortBuilder_Desc(t *testing.T) {
	sb := SortBuilder{}
	sb.Desc("name")
	s := sb.Build()
	value := s[0]
	if value.Key != "name" || value.Value != -1 {
		t.Error("value is incorrect")
	}
}
