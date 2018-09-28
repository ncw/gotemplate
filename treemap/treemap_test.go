package treemap

import (
	"testing"
)

func less(x, y Key) bool { return x.(int) < y.(int) }

func value(x string) Value { return Value(x) }

func TestNew(t *testing.T) {
	tr := New(less)
	if tr.Len() != 0 {
		t.Error("count should be zero")
	}
	if tr.endNode.left != nil {
		t.Error("root should be zero")
	}
}

func TestSet(t *testing.T) {
	x := value("x")
	tr := New(less)
	tr.Set(0, x)
	if tr.endNode.left.key != 0 {
		t.Errorf("wrong key, expected 0, got %d", tr.endNode.left.key)
	}
	if v := tr.endNode.left.value; v != x {
		t.Errorf("wrong returned value, expected '%s', got '%s'", x, v)
	}
	if tr.Len() != 1 {
		t.Errorf("wrong count, expected 1, got %d", tr.Len())
	}
}

func TestDel(t *testing.T) {
	tr := New(less)
	tr.Set(0, "x")
	tr.Del(0)
	if tr.Len() != 0 {
		t.Errorf("wrong count after deletion, expected 0, got %d", tr.Len())
	}
	if tr.endNode.left != nil {
		t.Error("wrong tree state after deletion")
	}
}

func TestGet(t *testing.T) {
	x := value("x")
	tr := New(less)
	tr.Set(0, x)
	v, ok := tr.Get(0)
	if v != x || !ok {
		t.Errorf("wrong returned value, expected 'x', got '%s'", v)
	}
	if tr.Len() != 1 {
		t.Errorf("wrong count, expected 1, got %d", tr.Len())
	}
	if v, ok := tr.Get(2); v != nil || ok {
		t.Errorf("wrong returned value, expected nil, got '%v'", v)
	}
	if tr.Len() != 1 {
		t.Errorf("wrong count, expected 1, got %d", tr.Len())
	}
}

func TestContains(t *testing.T) {
	tr := New(less)
	tr.Set(0, "x")
	val := tr.Contains(0)
	if !val {
		t.Error("existing is not exist")
	}
	val = tr.Contains(1)
	if val {
		t.Error("not existing is exist")
	}
}

func TestLen(t *testing.T) {
	tr := New(less)
	if tr.Len() != 0 {
		t.Errorf("wrong count, expected 0, got %d", tr.Len())
	}
	tr.Set(0, "x")
	if tr.Len() != 1 {
		t.Errorf("wrong count, expected 1, got %d", tr.Len())
	}
	tr.Set(1, "x")
	if tr.Len() != 2 {
		t.Errorf("wrong count, expected 2, got %d", tr.Len())
	}
	tr.Del(1)
	if tr.Len() != 1 {
		t.Errorf("wrong count, expected 1, got %d", tr.Len())
	}
	tr.Del(0)
	if tr.Len() != 0 {
		t.Errorf("wrong count, expected 0, got %d", tr.Len())
	}
}

func TestClear(t *testing.T) {
	tr := New(less)
	tr.Set(0, "x")
	tr.Set(1, "y")
	tr.Set(2, "z")
	tr.Clear()
	if tr.Len() != 0 {
		t.Error("count is not zero")
	}
	if tr.endNode.left != nil {
		t.Error("root is not nil")
	}
}

func testRange(t *testing.T, it, end ForwardIterator, exp []Value) {
	var got []Value
	for ; it != end; it.Next() {
		got = append(got, it.Value())
	}
	if len(got) != len(exp) {
		t.Errorf("wrong range length, expected %d, got %d", len(exp), len(got))
	}
	for i, v := range exp {
		if got[i] != v {
			t.Errorf("wrong value, expected '%s', got '%s'", exp[i], got[i])
		}
	}
}

func TestRange(t *testing.T) {
	tr := New(less)
	tr.Set(0, "x")
	tr.Set(1, "y")
	tr.Set(2, "z")
	tr.Set(3, "m")
	tr.Set(4, "n")
	it, end := tr.Range(1, 3)
	testRange(t, it, end, []Value{"y", "z", "m"})
	it, end = tr.Range(1, 9)
	testRange(t, it, end, []Value{"y", "z", "m", "n"})
}

func TestLowerBound(t *testing.T) {
	tr := New(less)
	it := tr.LowerBound(0)
	if it.Valid() {
		t.Error("lower bound should not exists")
		return
	}
	tr.Set(2, "a")
	tr.Set(4, "b")
	tr.Set(6, "c")
	tr.Set(8, "d")
	tr.Set(10, "e")
	tr.Set(12, "e")
	tr.Set(14, "e")
	tr.Set(16, "e")
	tr.Set(18, "e")
	tr.Set(20, "e")

	tbl := [][2]int{
		{0, 2},
		{2, 2},
		{3, 4},
		{4, 4},
		{9, 10},
		{10, 10},
		{11, 12},
		{19, 20},
		{20, 20},
	}

	for _, tb := range tbl {
		it = tr.LowerBound(tb[0])
		if !it.Valid() {
			t.Error("lower bound should exists")
			return
		}
		if k := it.Key(); k != tb[1] {
			t.Errorf("lower bound should be %v", tb[1])
			return
		}
	}

	it = tr.LowerBound(21)
	if it.Valid() {
		t.Error("lower bound should not exists")
		return
	}
}

func TestUpperBound(t *testing.T) {
	tr := New(less)
	it := tr.UpperBound(0)
	if it.Valid() {
		t.Error("upper bound should not exists")
		return
	}
	tr.Set(2, "a")
	tr.Set(4, "b")
	tr.Set(6, "c")
	tr.Set(8, "d")
	tr.Set(10, "e")
	tr.Set(12, "e")
	tr.Set(14, "e")
	tr.Set(16, "e")
	tr.Set(18, "e")
	tr.Set(20, "e")

	tbl := [][2]int{
		{0, 2},
		{2, 4},
		{3, 4},
		{4, 6},
		{9, 10},
		{10, 12},
		{11, 12},
		{19, 20},
	}

	for _, tb := range tbl {
		it = tr.UpperBound(tb[0])
		if !it.Valid() {
			t.Error("lower bound should exists")
			return
		}
		if k := it.Key(); k != tb[1] {
			t.Errorf("upper bound should be %v", tb[1])
			return
		}
	}

	it = tr.UpperBound(20)
	if it.Valid() {
		t.Error("upper bound should not exists")
		return
	}
	it = tr.UpperBound(21)
	if it.Valid() {
		t.Error("upper bound should not exists")
		return
	}
}

func TestEmptyRange(t *testing.T) {
	tr := New(less)
	tr.Set(0, "x")
	tr.Set(1, "y")
	tr.Set(2, "z")
	tr.Set(3, "m")
	tr.Set(4, "n")
	if rng, end := tr.Range(5, 10); rng != end {
		t.Error("range should be empty")
	}
}

func TestDelNil(t *testing.T) {
	x := "x"
	tr := New(less)
	tr.Set(0, value(x))
	tr.Del(1)
	if tr.Len() != 1 {
		t.Errorf("wrong count after del, expected 1, got %d", tr.Len())
	}
}

func TestIteration(t *testing.T) {
	kvs := []struct {
		key   Key
		value Value
	}{
		{0, "a"},
		{1, "b"},
		{2, "c"},
		{3, "d"},
		{4, "e"},
	}
	tr := New(less)
	for _, kv := range kvs {
		tr.Set(kv.key, kv.value)
	}
	assert := func(expKey Key, expValue Value, gotKey Key, gotValue Value) {
		if expKey != gotKey || expValue != gotValue {
			t.Errorf("expected %v, %s, got %v, %s", expKey, expValue, gotKey, gotValue)
		}
	}
	count := 0
	fwd := tr.Iterator()
	for ; fwd.Valid(); fwd.Next() {
		assert(kvs[count].key, kvs[count].value, fwd.Key(), fwd.Value())
		count++
	}
	for fwd != tr.Iterator() {
		fwd.Prev()
		count--
		assert(kvs[count].key, kvs[count].value, fwd.Key(), fwd.Value())
	}
	count = len(kvs)
	rev := tr.Reverse()
	for ; rev.Valid(); rev.Next() {
		count--
		assert(kvs[count].key, kvs[count].value, rev.Key(), rev.Value())
	}
	rbegin := tr.Reverse()
	for rev != rbegin {
		rev.Prev()
		assert(kvs[count].key, kvs[count].value, rev.Key(), rev.Value())
		count++
	}
}

func TestOutOfBoundsForwardIterationNext(t *testing.T) {
	tr := New(less)
	tr.Set(0, "a")
	tr.Set(1, "b")
	tr.Set(2, "c")
	tr.Set(3, "d")
	tr.Set(4, "e")
	it := tr.Iterator()
	for ; it.Valid(); it.Next() {
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("should have panicked!")
		}
	}()
	it.Next()
}

func TestOutOfBoundsForwardIterationPrev(t *testing.T) {
	tr := New(less)
	tr.Set(0, "a")
	tr.Set(1, "b")
	tr.Set(2, "c")
	tr.Set(3, "d")
	tr.Set(4, "e")
	it := tr.Iterator()
	defer func() {
		if r := recover(); r == nil {
			t.Error("should have panicked!")
		}
	}()
	it.Prev()
}

func TestOutOfBoundsReverseIterationNext(t *testing.T) {
	tr := New(less)
	tr.Set(0, "a")
	tr.Set(1, "b")
	tr.Set(2, "c")
	tr.Set(3, "d")
	tr.Set(4, "e")
	it := tr.Reverse()
	for ; it.Valid(); it.Next() {
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("should have panicked!")
		}
	}()
	it.Next()
}

func TestOutOfBoundsReverseIterationPrev(t *testing.T) {
	tr := New(less)
	tr.Set(0, "a")
	tr.Set(1, "b")
	tr.Set(2, "c")
	tr.Set(3, "d")
	tr.Set(4, "e")
	it := tr.Reverse()
	defer func() {
		if r := recover(); r == nil {
			t.Error("should have panicked!")
		}
	}()
	it.Prev()
}

func TestRangeSingle(t *testing.T) {
	tr := New(less)
	tr.Set(0, "a")
	tr.Set(1, "b")
	tr.Set(2, "c")
	visited := false
	for it, end := tr.Range(1, 1); it != end; it.Next() {
		if visited || it.Value() != "b" {
			t.Error("only single element 'b' should be found")
		}
		visited = true
	}
	if !visited {
		t.Error("single element 'b' should be found")
	}
}
