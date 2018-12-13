package treemap

import (
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

const NumIters = 10000
const RandMax = 40

func min(kv map[Key]Value) Key {
	if len(kv) == 0 {
		return nil
	}
	var key Key
	for k := range kv {
		if key == nil || k.(int) < key.(int) {
			key = k
		}
	}
	return key
}

func max(kv map[Key]Value) Key {
	if len(kv) == 0 {
		return nil
	}
	var key Key
	for k := range kv {
		if key == nil || k.(int) > key.(int) {
			key = k
		}
	}
	return key
}

type pair struct {
	k int
	v interface{}
}

func testRandomData() []pair {
	var kv []pair
	for i := 0; i < NumIters; i++ {
		k := int(rand.Int63n(RandMax))
		v := value(strconv.Itoa(int(rand.Int63n(RandMax))))
		kv = append(kv, pair{k, v})
	}
	return kv
}

func testKeys(t *testing.T, mp map[Key]Value, tr *TreeMap) {
	var gotKeys []int
	for it := tr.Iterator(); it.Valid(); it.Next() {
		gotKeys = append(gotKeys, it.Key().(int))
	}

	var expKeys []int
	for k := range mp {
		expKeys = append(expKeys, k.(int))
	}
	sort.Ints(expKeys)

	if !reflect.DeepEqual(gotKeys, expKeys) {
		t.Errorf("wrong keys, expected %v, got %v", expKeys, gotKeys)
	}
}

func testMinMax(t *testing.T, mp map[Key]Value, tr *TreeMap) {
	exp := min(mp)
	var got Key
	if it := tr.Iterator(); it.Valid() {
		got = it.Key()
	}
	if exp != got {
		t.Errorf("wrong min, expected %d, got %d", exp, got)
	}

	exp = max(mp)
	got = nil
	if it := tr.Reverse(); it.Valid() {
		got = it.Key()
	}
	if exp != got {
		t.Errorf("wrong max, expected %d, got %d", exp, got)
	}
}

func testReverse(t *testing.T, mp map[Key]Value, tr *TreeMap) {
	for it := tr.Reverse(); it.Valid(); it.Next() {
		if mp[it.Key()] != it.Value() {
			t.Errorf("wrong value, expected %s, got %s", mp[it.Key()], it.Value())
		}
	}
}

func TestRandom(t *testing.T) {
	tr := New(less)
	mp := make(map[Key]Value)
	kvs := testRandomData()
	for i, kv := range kvs {
		k, v := kv.k, kv.v
		exp, expOK := mp[k]
		if got, gotOK := tr.Get(k); got != exp || gotOK != expOK {
			t.Errorf("wrong returned value, expected %s, got %s", exp, got)
		}

		if i%3 == 0 && (i/200)%2 == 0 {
			tr.Set(k, v)
			mp[k] = v
		} else {
			delete(mp, k)
			tr.Del(k)
		}

		if len(mp) != tr.Len() {
			t.Errorf("wrong count, expected %d, got %d", len(mp), tr.Len())
			return
		}

		testKeys(t, mp, tr)
		testMinMax(t, mp, tr)
		testReverse(t, mp, tr)

		if !treeInvariant(tr.endNode.left) {
			t.Errorf("invariant error")
		}
	}
}
