package main

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	s := []string{"b", "c", "a", "e"}
	Sort(s)
	expected := []string{"a", "b", "c", "e"}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Got %v, Expected %v", s, expected)
	}
}

func TestSortF(t *testing.T) {
	s := []float64{0.1, -1.6, 1.5, +9, -8, 0.01, 1E99}
	SortF(s)
	expected := []float64{-8, -1.6, 0.01, 0.1, 1.5, +9, 1E99}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Got %v, Expected %v", s, expected)
	}
}

func TestSortGt(t *testing.T) {
	s := []string{"b", "c", "a", "e"}
	SortGt(s)
	expected := []string{"e", "c", "b", "a"}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Got %v, Expected %v", s, expected)
	}
}
