// This tests the set package
package set

import (
	"sort"
	"testing"
)

// Check the set is equal to the literal slice
func assertEqual(t *testing.T, s *Set, b []int) {
	a := make([]int, len(s.m))
	i := 0
	for elem := range s.m {
		a[i] = int(elem)
		i++
	}
	sort.Ints(a)
	sort.Ints(b)
	if len(a) != len(b) {
		t.Fatalf("Bad lengths %v vs %v", a, b)
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("%v != %v", a, b)
		}
	}
}

func TestNewSizedSet(t *testing.T) {
	a := NewSizedSet(0)
	if a.Len() != 0 {
		t.Fatal("length nonzero")
	}
	a = NewSizedSet(100)
	if a.Len() != 0 {
		t.Fatal("length nonzero")
	}
}

func TestNewSet(t *testing.T) {
	a := NewSet()
	if a.Len() != 0 {
		t.Fatal("length nonzero")
	}
}

func TestSetLen(t *testing.T) {
	a := NewSet()
	if a.Len() != 0 {
		t.Fatal("length nonzero")
	}
	a.Add(1)
	if a.Len() != 1 {
		t.Fatal("length not 1")
	}
	a.Discard(1)
	if a.Len() != 0 {
		t.Fatal("length nonzero")
	}
}

func TestSetContains(t *testing.T) {
	a := NewSet().Add(1)
	if a.Contains(0) {
		t.Fatal("0 found in set")
	}
	if !a.Contains(1) {
		t.Fatal("1 not found in set")
	}
	if a.Contains(2) {
		t.Fatal("2 found in set")
	}
}

func TestSetAdd(t *testing.T) {
	a := NewSet().Add(1).Add(2).Add(3)
	assertEqual(t, a, []int{1, 2, 3})
}

func TestSetAddList(t *testing.T) {
	a := NewSet().AddList([]A{1, 2, 3})
	assertEqual(t, a, []int{1, 2, 3})
}

func TestSetDiscard(t *testing.T) {
	a := NewSet().Add(1).Add(3).Discard(1).Discard(2)
	assertEqual(t, a, []int{3})
}

func TestSetRemove(t *testing.T) {
	a := NewSet().Add(1).Add(3)
	assertEqual(t, a, []int{1, 3})
	if a.Remove(1) != true {
		t.Fatal("1 not in set")
	}
	if a.Remove(2) != false {
		t.Fatal("2 in set")
	}
	assertEqual(t, a, []int{3})
}

func TestSetPop(t *testing.T) {
	a := NewSet().Add(1).Add(3)
	assertEqual(t, a, []int{1, 3})
	elem, found := a.Pop(1)
	if elem != 1 || found != true {
		t.Fatal("pop existing element failed")
	}
	elem, found = a.Pop(2)
	if elem != 2 || found != false {
		t.Fatal("pop non-existing element failed")
	}
	assertEqual(t, a, []int{3})
}

func TestSetAsList(t *testing.T) {
	a := NewSet().Add(1).Add(3)
	assertEqual(t, a, []int{1, 3})
	as := a.AsList()
	if len(as) != 2 {
		t.Fatal("length != 2")
	}
	if as[1] < as[0] {
		as[0], as[1] = as[1], as[0]
	}
	if as[0] != 1 || as[1] != 3 {
		t.Fatal("set as list failed")
	}
}

func TestSetClear(t *testing.T) {
	a := NewSet().Add(1).Add(3)
	assertEqual(t, a, []int{1, 3})
	a.Clear()
	assertEqual(t, a, []int{})
}

func TestSetCopy(t *testing.T) {
	a := NewSet().Add(1).Add(3)
	assertEqual(t, a, []int{1, 3})
	b := a.Copy()
	assertEqual(t, b, []int{1, 3})
	if a == b || &a.m == &b.m {
		t.Fatal("set copy failed")
	}
}

func TestSetDifference(t *testing.T) {
	assertEqual(t, NewSet().Difference(NewSet()), []int{})
	a := NewSet().Add(1).Add(3)
	b := NewSet().Add(1).Add(2)
	assertEqual(t, a.Difference(b), []int{3})
	assertEqual(t, a, []int{1, 3})
	assertEqual(t, a.Difference(NewSet()), []int{1, 3})
	assertEqual(t, a, []int{1, 3})
	assertEqual(t, NewSet().Difference(a), []int{})
	assertEqual(t, a, []int{1, 3})
	assertEqual(t, b.Difference(a), []int{2})
	assertEqual(t, b, []int{1, 2})
}

func TestSetDifferenceUpdate(t *testing.T) {
	assertEqual(t, NewSet().DifferenceUpdate(NewSet()), []int{})
	a := NewSet().Add(1).Add(3)
	b := NewSet().Add(1).Add(2)
	assertEqual(t, a.DifferenceUpdate(b), []int{3})
	a.DifferenceUpdate(b)
	assertEqual(t, a, []int{3})
	assertEqual(t, a.DifferenceUpdate(NewSet()), []int{3})
	assertEqual(t, NewSet().DifferenceUpdate(a), []int{})
	a.Add(1)
	assertEqual(t, b.DifferenceUpdate(a), []int{2})
	assertEqual(t, b, []int{2})
}

func TestSetIntersection(t *testing.T) {
	a := NewSet()
	b := NewSet()
	assertEqual(t, a.Intersection(b), []int{})
	assertEqual(t, b.Intersection(a), []int{})
	a.Add(1)
	a.Add(2)
	b.Add(2)
	b.Add(3)
	assertEqual(t, a.Intersection(b), []int{2})
	assertEqual(t, a, []int{1, 2})
	assertEqual(t, b.Intersection(a), []int{2})
	assertEqual(t, b, []int{2, 3})
}

func TestSetIntersectionUpdate(t *testing.T) {
	a := NewSet()
	b := NewSet()
	assertEqual(t, a.IntersectionUpdate(b), []int{})
	assertEqual(t, b.IntersectionUpdate(a), []int{})
	a.Add(1)
	a.Add(2)
	b.Add(2)
	b.Add(3)
	assertEqual(t, a.IntersectionUpdate(b), []int{2})
	assertEqual(t, a, []int{2})
	a.Add(1)
	a.Add(2)
	b.Add(2)
	b.Add(3)
	assertEqual(t, b.IntersectionUpdate(a), []int{2})
	assertEqual(t, b, []int{2})
}

func TestSetUnion(t *testing.T) {
	a := NewSet()
	b := NewSet()
	assertEqual(t, a.Union(b), []int{})
	assertEqual(t, b.Union(a), []int{})
	a.Add(1)
	a.Add(2)
	b.Add(2)
	b.Add(3)
	assertEqual(t, a.Union(b), []int{1, 2, 3})
	assertEqual(t, a, []int{1, 2})
	a.Clear().Add(1).Add(2)
	b.Clear().Add(2).Add(3)
	assertEqual(t, b.Union(a), []int{1, 2, 3})
	assertEqual(t, b, []int{2, 3})
}

func TestSetUpdate(t *testing.T) {
	a := NewSet()
	b := NewSet()
	assertEqual(t, a.Update(b), []int{})
	assertEqual(t, b.Update(a), []int{})
	a.Add(1)
	a.Add(2)
	b.Add(2)
	b.Add(3)
	assertEqual(t, a.Update(b), []int{1, 2, 3})
	assertEqual(t, a, []int{1, 2, 3})
	a.Clear().Add(1).Add(2)
	b.Clear().Add(2).Add(3)
	assertEqual(t, b.Update(a), []int{1, 2, 3})
	assertEqual(t, b, []int{1, 2, 3})
}

func TestSetIsSuperset(t *testing.T) {
	a := NewSet().Add(1).Add(2).Add(3)
	b := NewSet().Add(1).Add(2)
	assertEqual(t, a, []int{1, 2, 3})
	assertEqual(t, b, []int{1, 2})
	//test if superset returns correctly with strict true
	if a.IsSuperset(true, b) == false {
		t.Fatal("strict superset failed")
	}
	//test if superset returns correctly with strict false
	if a.IsSuperset(false, b) == false {
		t.Fatal("non-strict superset failed")
	}

	b.Add(3)
	assertEqual(t, a, []int{1, 2, 3})
	assertEqual(t, b, []int{1, 2, 3})
	//test if superset returns correctly with strict true
	if a.IsSuperset(true, b) == true {
		t.Fatal("strict superset failed")
	}
	//test if superset returns correctly with strict false
	if a.IsSuperset(false, b) == false {
		t.Fatal("non-strict superset failed")
	}
}

func TestSetIsSubset(t *testing.T) {
	a := NewSet().Add(1).Add(2).Add(3)
	b := NewSet().Add(1).Add(2)
	assertEqual(t, a, []int{1, 2, 3})
	assertEqual(t, b, []int{1, 2})
	//test if subset returns correctly with strict true
	if b.IsSubset(true, a) == false {
		t.Fatal("strict subset failed")
	}
	//test if subset returns correctly with strict false
	if b.IsSubset(false, a) == false {
		t.Fatal("non-strict subset failed")
	}
	b.Add(3)
	assertEqual(t, a, []int{1, 2, 3})
	assertEqual(t, b, []int{1, 2, 3})
	//test if subset returns correctly with strict true
	if b.IsSubset(true, a) == true {
		t.Fatal("strict subset failed")
	}
	//test if subset returns correctly with strict false
	if b.IsSubset(false, a) == false {
		t.Fatal("non-strict subset failed")
	}
}

func TestSetIsDisjoint(t *testing.T) {
	a := NewSet().Add(1)
	b := NewSet().Add(2)
	assertEqual(t, a, []int{1})
	assertEqual(t, b, []int{2})
	if a.IsDisjoint(b) == false || b.IsDisjoint(a) == false {
		t.Fatal("disjoint failed #1")
	}

	c := NewSet().Add(1).Add(2)
	d := NewSet().Add(2)
	assertEqual(t, c, []int{1, 2})
	assertEqual(t, d, []int{2})
	if c.IsDisjoint(d) == true || d.IsDisjoint(c) == true {
		t.Fatal("disjoint failed #2")
	}

	e := NewSet().Add(1)
	f := NewSet().Add(1).Add(2)
	assertEqual(t, e, []int{1})
	assertEqual(t, f, []int{1, 2})
	if e.IsDisjoint(f) == true || f.IsDisjoint(e) == true {
		t.Fatal("disjoint failed #3")
	}
}

func TestSetSymmetricDifference(t *testing.T) {
	a := NewSet().Add(1).Add(2)
	b := NewSet().Add(2).Add(3)
	c := a.SymmetricDifference(b)
	assertEqual(t, a, []int{1, 2})
	assertEqual(t, b, []int{2, 3})
	assertEqual(t, c, []int{1, 3})
}

func TestSetSymmetricDifferenceUpdate(t *testing.T) {
	a := NewSet().Add(1).Add(2)
	b := NewSet().Add(2).Add(3)
	assertEqual(t, a, []int{1, 2})
	assertEqual(t, b, []int{2, 3})
	a.SymmetricDifferenceUpdate(b)
	assertEqual(t, a, []int{1, 3})
}
