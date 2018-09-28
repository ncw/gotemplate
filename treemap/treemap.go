// Package treemap provides a generic key-sorted map. It uses red-black tree under the hood.
// You can use it as a template to generate a sorted map with specific key and value types.
// Iterators are designed after C++.
//
// Example:
//
//     package main
//
//     import "fmt"
//
//     //go:generate gotemplate "github.com/ncw/gotemplate/treemap" "intStringTreeMap(int, string)"
//
//     func less(x, y int) bool { return x < y }
//
//     func main() {
//         tr := newIntStringTreeMap(less)
//         tr.Set(0, "Hello")
//         tr.Set(1, "World")
//
//         for it := tr.Iterator(); it.Valid(); it.Next() {
//             fmt.Println(it.Key(), it.Value())
//         }
//     }
package treemap

// template type TreeMap(Key, Value)

// Key is a generic key type of the map
type Key interface{}

// Value is a generic value type of the map
type Value interface{}

// TreeMap is the red-black tree based map
type TreeMap struct {
	endNode   *node
	beginNode *node
	count     int
	// Less returns a < b
	Less func(a Key, b Key) bool
}

type node struct {
	right   *node
	left    *node
	parent  *node
	isBlack bool
	key     Key
	value   Value
}

// New creates and returns new TreeMap.
// Parameter less is a function returning a < b.
func New(less func(a Key, b Key) bool) *TreeMap {
	endNode := &node{isBlack: true}
	return &TreeMap{beginNode: endNode, endNode: endNode, Less: less}
}

// Len returns total count of elements in a map.
// Complexity: O(1).
func (t *TreeMap) Len() int { return t.count }

// Set sets the value and silently overrides previous value if it exists.
// Complexity: O(log N).
func (t *TreeMap) Set(key Key, value Value) {
	parent := t.endNode
	current := parent.left
	less := true
	for current != nil {
		parent = current
		switch {
		case t.Less(key, current.key):
			current = current.left
			less = true
		case t.Less(current.key, key):
			current = current.right
			less = false
		default:
			current.value = value
			return
		}
	}
	x := &node{parent: parent, value: value, key: key}
	if less {
		parent.left = x
	} else {
		parent.right = x
	}
	if t.beginNode.left != nil {
		t.beginNode = t.beginNode.left
	}
	t.insertFixup(x)
	t.count++
}

// Del deletes the value.
// Complexity: O(log N).
func (t *TreeMap) Del(key Key) {
	z := t.findNode(key)
	if z == nil {
		return
	}
	if t.beginNode == z {
		if z.right != nil {
			t.beginNode = z.right
		} else {
			t.beginNode = z.parent
		}
	}
	t.count--
	removeNode(t.endNode.left, z)
}

// Clear clears the map.
// Complexity: O(1).
func (t *TreeMap) Clear() {
	t.count = 0
	t.beginNode = t.endNode
	t.endNode.left = nil
}

// Get retrieves a value from a map for specified key and reports if it exists.
// Complexity: O(log N).
func (t *TreeMap) Get(id Key) (Value, bool) {
	node := t.findNode(id)
	if node == nil {
		node = t.endNode
	}
	return node.value, node != t.endNode
}

// Contains checks if key exists in a map.
// Complexity: O(log N)
func (t *TreeMap) Contains(id Key) bool { return t.findNode(id) != nil }

// Range returns a pair of iterators that you can use to go through all the keys in the range [from, to].
// More specifically it returns iterators pointing to lower bound and upper bound.
// Complexity: O(log N).
func (t *TreeMap) Range(from, to Key) (ForwardIterator, ForwardIterator) {
	return t.LowerBound(from), t.UpperBound(to)
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (t *TreeMap) LowerBound(key Key) ForwardIterator {
	result := t.endNode
	node := t.endNode.left
	if node == nil {
		return ForwardIterator{tree: t, node: t.endNode}
	}
	for {
		if t.Less(node.key, key) {
			if node.right != nil {
				node = node.right
			} else {
				return ForwardIterator{tree: t, node: result}
			}
		} else {
			result = node
			if node.left != nil {
				node = node.left
			} else {
				return ForwardIterator{tree: t, node: result}
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (t *TreeMap) UpperBound(key Key) ForwardIterator {
	result := t.endNode
	node := t.endNode.left
	if node == nil {
		return ForwardIterator{tree: t, node: t.endNode}
	}
	for {
		if !t.Less(key, node.key) {
			if node.right != nil {
				node = node.right
			} else {
				return ForwardIterator{tree: t, node: result}
			}
		} else {
			result = node
			if node.left != nil {
				node = node.left
			} else {
				return ForwardIterator{tree: t, node: result}
			}
		}
	}
}

// Iterator returns an iterator for tree map.
// It starts at the first element and goes to the one-past-the-end position.
// You can iterate a map at O(N) complexity.
// Method complexity: O(1)
func (t *TreeMap) Iterator() ForwardIterator { return ForwardIterator{tree: t, node: t.beginNode} }

// Reverse returns a reverse iterator for tree map.
// It starts at the last element and goes to the one-before-the-start position.
// You can iterate a map at O(N) complexity.
// Method complexity: O(log N)
func (t *TreeMap) Reverse() ReverseIterator {
	node := t.endNode.left
	if node != nil {
		node = mostRight(node)
	}
	return ReverseIterator{tree: t, node: node}
}

func (t *TreeMap) findNode(id Key) *node {
	current := t.endNode.left
	for current != nil {
		switch {
		case t.Less(id, current.key):
			current = current.left
		case t.Less(current.key, id):
			current = current.right
		default:
			return current
		}
	}
	return nil
}

func mostLeft(x *node) *node {
	for x.left != nil {
		x = x.left
	}
	return x
}

func mostRight(x *node) *node {
	for x.right != nil {
		x = x.right
	}
	return x
}

func successor(x *node) *node {
	if x.right != nil {
		return mostLeft(x.right)
	}
	for x != x.parent.left {
		x = x.parent
	}
	return x.parent
}

func predecessor(x *node) *node {
	if x.left != nil {
		return mostRight(x.left)
	}
	for x.parent != nil && x != x.parent.right {
		x = x.parent
	}
	return x.parent
}

func rotateLeft(x *node) {
	y := x.right
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func rotateRight(x *node) {
	y := x.left
	x.left = y.right
	if x.left != nil {
		x.left.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.right = x
	x.parent = y
}

func (t *TreeMap) insertFixup(x *node) {
	root := t.endNode.left
	x.isBlack = x == root
	for x != root && !x.parent.isBlack {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x != x.parent.left {
					x = x.parent
					rotateLeft(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateRight(x)
				break
			}
		} else {
			y := x.parent.parent.left
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x == x.parent.left {
					x = x.parent
					rotateRight(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateLeft(x)
				break
			}
		}
	}
}

// nolint: gocyclo
func removeNode(root *node, z *node) {
	var y *node
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = successor(z)
	}
	var x *node
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}
	var w *node
	if x != nil {
		x.parent = y.parent
	}
	if y == y.parent.left {
		y.parent.left = x
		if y != root {
			w = y.parent.right
		} else {
			root = x // w == nil
		}
	} else {
		y.parent.right = x
		w = y.parent.left
	}
	removedBlack := y.isBlack
	if y != z {
		y.parent = z.parent
		if z == z.parent.left {
			y.parent.left = y
		} else {
			y.parent.right = y
		}
		y.left = z.left
		y.left.parent = y
		y.right = z.right
		if y.right != nil {
			y.right.parent = y
		}
		y.isBlack = z.isBlack
		if root == z {
			root = y
		}
	}
	if removedBlack && root != nil {
		if x != nil {
			x.isBlack = true
		} else {
			for {
				if w != w.parent.left {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateLeft(w.parent)
						if root == w.left {
							root = w
						}
						w = w.left.right
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if x == root || !x.isBlack {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.right == nil || w.right.isBlack {
							w.left.isBlack = true
							w.isBlack = false
							rotateRight(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.right.isBlack = true
						rotateLeft(w.parent)
						break
					}
				} else {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateRight(w.parent)
						if root == w.right {
							root = w
						}
						w = w.right.left
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if !x.isBlack || x == root {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.left == nil || w.left.isBlack {
							w.right.isBlack = true
							w.isBlack = false
							rotateLeft(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.left.isBlack = true
						rotateRight(w.parent)
						break
					}
				}
			}
		}
	}
}

// ForwardIterator represents a position in a tree map.
// It is designed to iterate a map in a forward order.
// It can point to any position from the first element to the one-past-the-end element.
type ForwardIterator struct {
	tree *TreeMap
	node *node
}

// Valid reports if an iterator's position is valid.
// In other words it returns true if an iterator is not at the one-past-the-end position.
func (i ForwardIterator) Valid() bool { return i.node != i.tree.endNode }

// Next moves an iterator to the next element.
// It panics if goes out of bounds.
func (i *ForwardIterator) Next() {
	if i.node == i.tree.endNode {
		panic("out of bound iteration")
	}
	i.node = successor(i.node)
}

// Prev moves an iterator to the previous element.
// It panics if goes out of bounds.
func (i *ForwardIterator) Prev() {
	i.node = predecessor(i.node)
	if i.node == nil {
		panic("out of bound iteration")
	}
}

// Key returns a key at an iterator's position
func (i ForwardIterator) Key() Key { return i.node.key }

// Value returns a value at an iterator's position
func (i ForwardIterator) Value() Value { return i.node.value }

// ReverseIterator represents a position in a tree map.
// It is designed to iterate a map in a reverse order.
// It can point to any position from the one-before-the-start element to the last element.
type ReverseIterator struct {
	tree *TreeMap
	node *node
}

// Valid reports if an iterator's position is valid.
// In other words it returns true if an iterator is not at the one-before-the-start position.
func (i ReverseIterator) Valid() bool { return i.node != nil }

// Next moves an iterator to the next element in reverse order.
// It panics if goes out of bounds.
func (i *ReverseIterator) Next() {
	if i.node == nil {
		panic("out of bound iteration")
	}
	i.node = predecessor(i.node)
}

// Prev moves an iterator to the previous element in reverse order.
// It panics if goes out of bounds.
func (i *ReverseIterator) Prev() {
	if i.node != nil {
		i.node = successor(i.node)
	} else {
		i.node = i.tree.beginNode
	}
	if i.node == i.tree.endNode {
		panic("out of bound iteration")
	}
}

// Key returns a key at an iterator's position
func (i ReverseIterator) Key() Key { return i.node.key }

// Value returns a value at an iterator's position
func (i ReverseIterator) Value() Value { return i.node.value }
