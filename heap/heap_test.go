// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	"math/rand"
	"testing"
)

func verify(t *testing.T, h *Heap, i int) {
	hs := *h
	n := len(hs)
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if hs[j1] < hs[i] {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, hs[i], j1, hs[j1])
			return
		}
		verify(t, h, j1)
	}
	if j2 < n {
		if hs[j2] < hs[i] {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, hs[i], j1, hs[j2])
			return
		}
		verify(t, h, j2)
	}
}

func TestInit0(t *testing.T) {
	h := new(Heap)
	for i := 20; i > 0; i-- {
		h.Push(0) // all elements are the same
	}
	h.Init()
	verify(t, h, 0)

	for i := 1; len(*h) > 0; i++ {
		x := h.Pop()
		verify(t, h, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	h := new(Heap)
	for i := 20; i > 0; i-- {
		h.Push(A(i)) // all elements are different
	}
	h.Init()
	verify(t, h, 0)

	for i := 1; len(*h) > 0; i++ {
		x := h.Pop()
		verify(t, h, 0)
		if int(x) != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	h := new(Heap)
	verify(t, h, 0)

	for i := 20; i > 10; i-- {
		h.Push(A(i))
	}
	h.Init()
	verify(t, h, 0)

	for i := 10; i > 0; i-- {
		h.Push(A(i))
		verify(t, h, 0)
	}

	for i := 1; len(*h) > 0; i++ {
		x := h.Pop()
		if i < 20 {
			h.Push(A(20 + i))
		}
		verify(t, h, 0)
		if int(x) != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	h := new(Heap)
	for i := 0; i < 10; i++ {
		h.Push(A(i))
	}
	verify(t, h, 0)

	for len(*h) > 0 {
		i := len(*h) - 1
		x := h.Remove(i)
		if int(x) != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		verify(t, h, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := new(Heap)
	for i := 0; i < 10; i++ {
		h.Push(A(i))
	}
	verify(t, h, 0)

	for i := 0; len(*h) > 0; i++ {
		x := h.Remove(0)
		if int(x) != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		verify(t, h, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := new(Heap)
	for i := 0; i < N; i++ {
		h.Push(A(i))
	}
	verify(t, h, 0)

	m := make(map[int]bool)
	for len(*h) > 0 {
		m[int(h.Remove((len(*h)-1)/2))] = true
		verify(t, h, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := make(Heap, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(0) // all elements are the same
		}
	}
}

func TestFix(t *testing.T) {
	h := new(Heap)
	verify(t, h, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(A(i))
	}
	verify(t, h, 0)

	if (*h)[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", (*h)[0])
	}
	(*h)[0] = 210
	h.Fix(0)
	verify(t, h, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(len(*h))
		if i&1 == 0 {
			(*h)[elem] *= 2
		} else {
			(*h)[elem] /= 2
		}
		h.Fix(elem)
		verify(t, h, 0)
	}
}
