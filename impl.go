package goART

import "bytes"

func parseLeaf(node artNode) (leaf *leafNode, ok bool) {
	if node.kind() != kindLeaf {
		return
	}
	leaf, ok = node.(*leafNode)
	if !ok {
		panic("malformed node")
	}
	return
}

func keyAt(key []byte, idx int) byte {
	if idx >= 0 && idx < len(key) {
		return key[idx]
	}
	return 0
}

func makeLeaf(key []byte, val any) artNode {
	return &leafNode{key, val}
}

func doSearch(node *artNode, key []byte, depth int) (val any, exist bool) {
	if leaf, ok := parseLeaf(*node); ok {
		if bytes.Compare(leaf.key, key) == 0 {
			return leaf.val, true
		}
		return
	}
	prefix, prefLen := (*node).prefixAndLen()
	if prefLen != checkPrefix(prefix, prefLen, key, depth) {
		return
	}
	depth += prefLen
	var next = (*node).findChild(keyAt(key, depth))
	if next == nil {
		return
	}
	return doSearch(next, key, depth)
}

func checkPrefix(prefix []byte, prefLen int, key []byte, depth int) (common int) {
	for common = 0; common < len(prefix) && depth+common < len(key); common++ {
		if prefix[common] != key[depth+common] {
			return
		}
	}
	if prefLen <= len(key)-depth {
		return prefLen
	}
	return
}

func fullPrefix(node *artNode, depth int) []byte {
	for node != nil {
		if leaf, ok := parseLeaf(*node); ok {
			return leaf.key[depth:]
		}
		node = (*node).anyChild()
	}
	panic("malformed node")
}

func doInsert(node *artNode, key []byte, depth int, val any) (old any, replaced bool) {
	if leaf, ok := parseLeaf(*node); ok {
		if bytes.Compare(leaf.key, key) == 0 {
			old = leaf.val
			replaced = true
			leaf.val = val
			return
		}
		n4 := new(node4)
		var idx = depth
		for ; idx < len(leaf.key) && idx < len(key); idx++ {
			if leaf.key[idx] != key[idx] {
				break
			}
		}
		n4.setPrefix(key[depth:idx])
		n4.addChild(keyAt(leaf.key, idx), *node)
		n4.addChild(keyAt(key, idx), makeLeaf(key, val))
		*node = n4
		return
	}
	prefix, prefLen := (*node).prefixAndLen()
	common := checkPrefix(prefix, prefLen, key, depth)
	if prefLen != common {
		n4 := new(node4)
		n4.setPrefix(key[depth : depth+common])
		n4.addChild(keyAt(key, depth+common), makeLeaf(key, val))
		if common >= len(prefix) {
			prefix = fullPrefix(node, depth)
		}
		(*node).setPrefix(prefix[common:])
		n4.addChild(keyAt(prefix, common), *node)

		*node = n4
		return
	}
	depth += prefLen
	next := (*node).findChild(keyAt(key, depth))
	if next == nil {
		*node = (*node).addChild(keyAt(key, depth), makeLeaf(key, val))
		return
	}
	return doInsert(next, key, depth, val)
}

func tryCompressPrefix(oldNode artNode, depth int, newNode artNode) {
	if oldNode.kind() != kindNode4 || newNode.kind() == kindLeaf {
		return
	}
	if oldNode == newNode {
		return
	}
	oldPref, oldLen := oldNode.prefixAndLen()
	newPref, newLen := newNode.prefixAndLen()
	if len(oldPref) == oldLen && len(newPref) == newLen {
		oldPref = append(oldPref, newPref...)
		newNode.setPrefix(oldPref)
		return
	}
	newNode.setPrefix(fullPrefix(&newNode, depth)[:oldLen+newLen])
}

func doDelete(node *artNode, key []byte, depth int) (old any, deleted bool) {
	prefix, prefLen := (*node).prefixAndLen()
	common := checkPrefix(prefix, prefLen, key, depth)
	if prefLen != common {
		return
	}
	depth += prefLen
	next := (*node).findChild(keyAt(key, depth))
	if next == nil {
		return
	}
	if leaf, ok := parseLeaf(*next); ok {
		if 0 == bytes.Compare(leaf.key, key) {
			old, deleted = leaf.val, true
			var oldNode = *node
			*node = (*node).removeChild(keyAt(key, depth))
			tryCompressPrefix(oldNode, depth-prefLen, *node)
		}
		return
	}
	return doDelete(next, key, depth)
}
