package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestSet(t *testing.T) {
	s := newmySet()
	s.Add("Sausage")
	s.Add("Bacon")
	actual := s.AsList()
	sort.Strings(actual)
	expected := []string{"Bacon", "Sausage"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %v, Expected %v", actual, expected)
	}
}
