package goART

import "bytes"

type Tree struct {
	size int
	root artNode
}

func (t *Tree) Search(key []byte) (val any, exist bool) {
	if t.root == nil {
		return
	}
	return doSearch(&t.root, key, 0)
}

func (t *Tree) Insert(key []byte, val any) (old any, replaced bool) {
	if t.root == nil {
		t.root = makeLeaf(key, val)
		t.size++
		return
	}
	old, replaced = doInsert(&t.root, key, 0, val)
	if !replaced {
		t.size++
	}
	return
}

func (t *Tree) Delete(key []byte) (old any, deleted bool) {
	if t.root == nil {
		return
	} else if leaf, ok := parseLeaf(t.root); ok {
		if 0 != bytes.Compare(leaf.key, key) {
			return
		}
		old = leaf.val
		deleted = true
		t.root = nil
		t.size--
		return
	}
	old, deleted = doDelete(&t.root, key, 0)
	if deleted {
		t.size--
	}
	return
}

func (t *Tree) Size() int {
	return t.size
}
