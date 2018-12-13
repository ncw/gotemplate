package treemap

// nolint: gocyclo
func treeSubInvariant(x *node) int {
	if x == nil {
		return 1
	}
	if x.left != nil && x.left.parent != x {
		return 0
	}
	if x.right != nil && x.right.parent != x {
		return 0
	}
	if x.left == x.right && x.left != nil {
		return 0
	}
	if !x.isBlack {
		if x.left != nil && !x.left.isBlack {
			return 0
		}
		if x.right != nil && !x.right.isBlack {
			return 0
		}
	}
	h := treeSubInvariant(x.left)
	if h == 0 {
		return 0
	}
	if h != treeSubInvariant(x.right) {
		return 0
	}
	if x.isBlack {
		h++
	}
	return h
}

func treeInvariant(root *node) bool {
	if root == nil {
		return true
	}
	if root.parent == nil {
		return false
	}
	if root != root.parent.left {
		return false
	}
	if !root.isBlack || !root.parent.isBlack {
		return false
	}
	return treeSubInvariant(root) != 0
}
