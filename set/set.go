// Template Set type
//
// Tries to be similar to Python's set type
package set

// template type Set(A)
type A int

// Store struct{} which are 0 size as the members in the map
type SetNothing struct{}

type Set struct {
	m map[A]SetNothing
}

// Returns a new empty set with the given capacity
func NewSizedSet(capacity int) *Set {
	return &Set{
		m: make(map[A]SetNothing, capacity),
	}
}

// Returns a new empty set
func NewSet() *Set {
	return NewSizedSet(0)
}

// Returns the number of elements in the set
func (s *Set) Len() int {
	return len(s.m)
}

// Contains returns whether elem is in the set or not
func (s *Set) Contains(elem A) bool {
	_, found := s.m[elem]
	return found
}

// Add adds elem to the set, returning the set
//
// If the element already exists then it has no effect
func (s *Set) Add(elem A) *Set {
	s.m[elem] = SetNothing{}
	return s
}

// AddList adds a list of elems to the set
//
// If the elements already exists then it has no effect
func (s *Set) AddList(elems []A) *Set {
	for _, elem := range elems {
		s.m[elem] = SetNothing{}
	}
	return s
}

// Discard removes elem from the set
//
// If it wasn't in the set it does nothing
//
// It returns the set
func (s *Set) Discard(elem A) *Set {
	delete(s.m, elem)
	return s
}

// Remove removes elem from the set
//
// It returns whether the elem was in the set or not
func (s *Set) Remove(elem A) bool {
	_, found := s.m[elem]
	if found {
		delete(s.m, elem)
	}
	return found
}

// Pop removes elem from the set and returns it
//
// It also returns whether the elem was found or not
func (s *Set) Pop(elem A) (A, bool) {
	_, found := s.m[elem]
	if found {
		delete(s.m, elem)
	}
	return elem, found
}

// AsList returns all the elements as a slice
func (s *Set) AsList() []A {
	elems := make([]A, len(s.m))
	i := 0
	for elem := range s.m {
		elems[i] = elem
		i++
	}
	return elems
}

// Clear removes all the elements
func (s *Set) Clear() *Set {
	s.m = make(map[A]SetNothing)
	return s
}

// Copy returns a shallow copy of the Set
func (s *Set) Copy() *Set {
	newSet := NewSizedSet(len(s.m))
	for elem := range s.m {
		newSet.m[elem] = SetNothing{}
	}
	return newSet
}

// Difference returns a new set with all the elements that are in this
// set but not in the other
func (s *Set) Difference(other *Set) *Set {
	newSet := NewSizedSet(len(s.m))
	for elem := range s.m {
		if _, found := other.m[elem]; !found {
			newSet.m[elem] = SetNothing{}
		}
	}
	return newSet
}

// DifferenceUpdate removes all the elements that are in the other set
// from this set.  It returns the set.
func (s *Set) DifferenceUpdate(other *Set) *Set {
	m := s.m
	for elem := range other.m {
		delete(m, elem)
	}
	return s
}

// Intersection returns a new set with all the elements that are only in this
// set and the other set. It returns the new set.
func (s *Set) Intersection(other *Set) *Set {
	newSet := NewSizedSet(len(s.m) + len(other.m))
	for elem := range s.m {
		if _, found := other.m[elem]; found {
			newSet.m[elem] = SetNothing{}
		}
	}
	for elem := range other.m {
		if _, found := s.m[elem]; found {
			newSet.m[elem] = SetNothing{}
		}
	}
	return newSet
}

// IntersectionUpdate changes this set so that it only contains
// elements that are in both this set and the other set.  It returns
// the set.
func (s *Set) IntersectionUpdate(other *Set) *Set {
	for elem := range s.m {
		if _, found := other.m[elem]; !found {
			delete(s.m, elem)
		}
	}
	return s
}

// Union returns a new set with all the elements that are in either
// set. It returns the new set.
func (s *Set) Union(other *Set) *Set {
	newSet := NewSizedSet(len(s.m) + len(other.m))
	for elem := range s.m {
		newSet.m[elem] = SetNothing{}
	}
	for elem := range other.m {
		newSet.m[elem] = SetNothing{}
	}
	return newSet
}

// Update adds all the elements from the other set to this set.
// It returns the set.
func (s *Set) Update(other *Set) *Set {
	for elem := range other.m {
		s.m[elem] = SetNothing{}
	}
	return s
}

// IsSuperset returns a bool indicating whether this set is a superset of other set.
func (s *Set) IsSuperset(strict bool, other *Set) bool {
	if strict && len(other.m) >= len(s.m) {
		return false
	}
A:
	for v := range other.m {
		for i := range s.m {
			if v == i {
				continue A
			}
		}
		return false
	}
	return true
}

// IsSubset returns a bool indicating whether this set is a subset of other set.
func (s *Set) IsSubset(strict bool, other *Set) bool {
	if strict && len(s.m) >= len(other.m) {
		return false
	}
A:
	for v := range s.m {
		for i := range other.m {
			if v == i {
				continue A
			}
		}
		return false
	}
	return true
}

// IsDisjoint returns a bool indicating whether this set and other set have any elements in common.
func (s *Set) IsDisjoint(other *Set) bool {
	for v := range s.m {
		if other.Contains(v) {
			return false
		}
	}
	return true
}

// SymmetricDifference returns a new set of all elements that are a member of exactly
// one of this set and other set(elements which are in one of the sets, but not in both).
func (s *Set) SymmetricDifference(other *Set) *Set {
	work1 := s.Union(other)
	work2 := s.Intersection(other)
	for v := range work2.m {
		delete(work1.m, v)
	}
	return work1
}

// SymmetricDifferenceUpdate modifies this set to be a set of all elements that are a member
// of exactly one of this set and other set(elements which are in one of the sets,
// but not in both) and returns this set.
func (s *Set) SymmetricDifferenceUpdate(other *Set) *Set {
	work := s.SymmetricDifference(other)
	*s = *work
	return s
}