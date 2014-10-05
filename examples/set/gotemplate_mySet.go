// Template Set type
//
// Tries to be similar to Python's set type
package main

// template type Set(A)

// Store struct{} which are 0 size as the members in the map
type mySetNothing struct{}

type mySet struct {
	m map[string]mySetNothing
}

// Returns a new empty set with the given capacity
func newSizedmySet(capacity int) *mySet {
	return &mySet{
		m: make(map[string]mySetNothing, capacity),
	}
}

// Returns a new empty set
func newmySet() *mySet {
	return newSizedmySet(0)
}

// Returns the number of elements in the set
func (s *mySet) Len() int {
	return len(s.m)
}

// Contains returns whether elem is in the set or not
func (s *mySet) Contains(elem string) bool {
	_, found := s.m[elem]
	return found
}

// Add adds elem to the set, returning the set
//
// If the element already exists then it has no effect
func (s *mySet) Add(elem string) *mySet {
	s.m[elem] = mySetNothing{}
	return s
}

// AddList adds a list of elems to the set
//
// If the elements already exists then it has no effect
func (s *mySet) AddList(elems []string) *mySet {
	for _, elem := range elems {
		s.m[elem] = mySetNothing{}
	}
	return s
}

// Discard removes elem from the set
//
// If it wasn't in the set it does nothing
//
// It returns the set
func (s *mySet) Discard(elem string) *mySet {
	delete(s.m, elem)
	return s
}

// Remove removes elem from the set
//
// It returns whether the elem was in the set or not
func (s *mySet) Remove(elem string) bool {
	_, found := s.m[elem]
	if found {
		delete(s.m, elem)
	}
	return found
}

// Pop removes elem from the set and returns it
//
// It also returns whether the elem was found or not
func (s *mySet) Pop(elem string) (string, bool) {
	_, found := s.m[elem]
	if found {
		delete(s.m, elem)
	}
	return elem, found
}

// AsList returns all the elements as a slice
func (s *mySet) AsList() []string {
	elems := make([]string, len(s.m))
	i := 0
	for elem := range s.m {
		elems[i] = elem
		i++
	}
	return elems
}

// Clear removes all the elements
func (s *mySet) Clear() *mySet {
	s.m = make(map[string]mySetNothing)
	return s
}

// Copy returns a shallow copy of the Set
func (s *mySet) Copy() *mySet {
	newSet := newSizedmySet(len(s.m))
	for elem := range s.m {
		newSet.m[elem] = mySetNothing{}
	}
	return newSet
}

// Difference returns a new set with all the elements that are in this
// set but not in the other
func (s *mySet) Difference(other *mySet) *mySet {
	newSet := newSizedmySet(len(s.m))
	for elem := range s.m {
		if _, found := other.m[elem]; !found {
			newSet.m[elem] = mySetNothing{}
		}
	}
	return newSet
}

// DifferenceUpdate removes all the elements that are in the other set
// from this set.  It returns the set.
func (s *mySet) DifferenceUpdate(other *mySet) *mySet {
	m := s.m
	for elem := range other.m {
		delete(m, elem)
	}
	return s
}

// Intersection returns a new set with all the elements that are only in this
// set and the other set. It returns the new set.
func (s *mySet) Intersection(other *mySet) *mySet {
	newSet := newSizedmySet(len(s.m) + len(other.m))
	for elem := range s.m {
		if _, found := other.m[elem]; found {
			newSet.m[elem] = mySetNothing{}
		}
	}
	for elem := range other.m {
		if _, found := s.m[elem]; found {
			newSet.m[elem] = mySetNothing{}
		}
	}
	return newSet
}

// IntersectionUpdate changes this set so that it only contains
// elements that are in both this set and the other set.  It returns
// the set.
func (s *mySet) IntersectionUpdate(other *mySet) *mySet {
	for elem := range s.m {
		if _, found := other.m[elem]; !found {
			delete(s.m, elem)
		}
	}
	return s
}

// Union returns a new set with all the elements that are in either
// set. It returns the new set.
func (s *mySet) Union(other *mySet) *mySet {
	newSet := newSizedmySet(len(s.m) + len(other.m))
	for elem := range s.m {
		newSet.m[elem] = mySetNothing{}
	}
	for elem := range other.m {
		newSet.m[elem] = mySetNothing{}
	}
	return newSet
}

// Update adds all the elements from the other set to this set.
// It returns the set.
func (s *mySet) Update(other *mySet) *mySet {
	for elem := range other.m {
		s.m[elem] = mySetNothing{}
	}
	return s
}

/*
 |  isdisjoint(...)
 |      Return True if two sets have a null intersection.
 |
 |  issubset(...)
 |      Report whether another set contains this set.
 |
 |  issuperset(...)
 |      Report whether this set contains another set.
 |
 |  symmetric_difference(...)
 |      Return the symmetric difference of two sets as a new set.
 |
 |      (i.e. all elements that are in exactly one of the sets.)
 |
 |  symmetric_difference_update(...)
 |      Update a set with the symmetric difference of itself and another.
*/
