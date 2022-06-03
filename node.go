package goART

type kind uint8

const (
	kindLeaf kind = iota
	kindNode4
	kindNode16
	kindNode48
	kindNode256
)

const (
	prefLenMax = 8
)

const (
	keySizeMaxNode4   = 4
	childSizeMaxNode4 = 4
	childSizeMinNode4 = 2

	keySizeMaxNode16   = 16
	childSizeMaxNode16 = 16
	childSizeMinNode16 = 5

	keySizeMaxNode48   = 256
	childSizeMaxNode48 = 48
	childSizeMinNode48 = 17

	childSizeMaxNode256 = 256
	childSizeMinNode256 = 49
)

type artNode interface {
	kind() kind
	prefixAndLen() (prefix []byte, prefLen int)
	setPrefix(prefix []byte)
	findChild(b byte) *artNode
	addChild(b byte, child artNode) artNode
	removeChild(b byte) artNode
	anyChild() *artNode
}

type baseNode struct {
	prefLen   int
	prefix    [prefLenMax]byte
	childSize uint8
}

func (base *baseNode) prefixAndLen() (prefix []byte, prefLen int) {
	prefLen = base.prefLen
	if prefLen <= len(base.prefix) {
		prefix = base.prefix[:prefLen]
	}
	return
}

func (base *baseNode) setPrefix(prefix []byte) {
	base.prefLen = len(prefix)
	copy(base.prefix[:], prefix)
}

type leafNode struct {
	key []byte
	val any
}

func (leaf *leafNode) kind() kind {
	return kindLeaf
}

func (leaf *leafNode) prefixAndLen() ([]byte, int) {
	panic("implement me")
}

func (leaf *leafNode) setPrefix(_ []byte) {
}

func (leaf *leafNode) findChild(_ byte) *artNode {
	panic("implement me")
}

func (leaf *leafNode) addChild(_ byte, _ artNode) artNode {
	panic("implement me")
}

func (leaf *leafNode) removeChild(_ byte) artNode {
	panic("implement me")
}

func (leaf *leafNode) anyChild() *artNode {
	return nil
}

type node4 struct {
	baseNode
	keys     [keySizeMaxNode4]byte
	children [childSizeMaxNode4]*artNode
}

func (n4 *node4) kind() kind {
	return kindNode4
}

func (n4 *node4) findChild(b byte) *artNode {
	for i := uint8(0); i < n4.baseNode.childSize; i++ {
		if n4.keys[i] == b {
			return n4.children[i]
		}
	}
	return nil
}

func (n4 *node4) addChild(b byte, child artNode) artNode {
	if n4.childSize+1 <= childSizeMaxNode4 {
		n4.keys[n4.childSize] = b
		n4.children[n4.childSize] = &child
		n4.childSize++
		return n4
	}
	n16 := n4.upgrade()
	n16.addChild(b, child)
	return n16
}

func (n4 *node4) removeChild(b byte) artNode {
	var idx uint8 = 0
	for ; idx < n4.childSize; idx++ {
		if n4.keys[idx] != b {
			continue
		}
		n4.keys[idx] = n4.keys[n4.childSize-1]
		n4.children[idx], n4.children[n4.childSize-1] = n4.children[n4.childSize-1], nil
		n4.childSize--
		if n4.childSize >= childSizeMinNode4 {
			break
		}
		return *n4.children[0]
	}
	return n4
}

func (n4 *node4) anyChild() *artNode {
	return n4.children[0]
}

func (n4 *node4) upgrade() (n16 *node16) {
	n16 = &node16{baseNode: n4.baseNode}
	copy(n16.keys[:], n4.keys[:])
	copy(n16.children[:], n4.children[:])
	return
}

func (n16 *node16) downgrade() (n4 *node4) {
	n4 = &node4{baseNode: n16.baseNode}
	copy(n4.keys[:], n16.keys[:])
	copy(n4.children[:], n16.children[:])
	return
}

type node16 struct {
	baseNode
	keys     [keySizeMaxNode16]byte
	children [childSizeMaxNode16]*artNode
}

func (n16 *node16) kind() kind {
	return kindNode16
}

func (n16 *node16) findChild(b byte) *artNode {
	for i := uint8(0); i < n16.baseNode.childSize; i++ {
		if n16.keys[i] == b {
			return n16.children[i]
		}
	}
	return nil
}

func (n16 *node16) addChild(b byte, child artNode) artNode {
	if n16.childSize+1 <= childSizeMaxNode16 {
		n16.keys[n16.childSize] = b
		n16.children[n16.childSize] = &child
		n16.childSize++
		return n16
	}
	n48 := n16.upgrade()
	n48.addChild(b, child)
	return n48
}

func (n16 *node16) removeChild(b byte) artNode {
	var idx uint8 = 0
	for ; idx < n16.childSize; idx++ {
		if n16.keys[idx] != b {
			continue
		}
		n16.keys[idx] = n16.keys[n16.childSize-1]
		n16.children[idx], n16.children[n16.childSize-1] = n16.children[n16.childSize-1], nil
		n16.childSize--
		if n16.childSize >= childSizeMinNode16 {
			break
		}
		return n16.downgrade()
	}
	return n16
}

func (n16 *node16) anyChild() *artNode {
	return n16.children[0]
}

func (n16 *node16) upgrade() (n48 *node48) {
	n48 = &node48{baseNode: n16.baseNode}
	for i := uint8(0); i < n16.childSize; i++ {
		n48.childIndex[n16.keys[i]] = i + 1
	}
	copy(n48.children[:], n16.children[:])
	return
}

func (n48 *node48) downgrade() (n16 *node16) {
	n16 = &node16{baseNode: n48.baseNode}
	for i := range n48.childIndex {
		child := n48.findChild(uint8(i))
		if child != nil {
			n16.addChild(byte(i), *child)
		}
	}
	return
}

type node48 struct {
	baseNode
	childIndex [keySizeMaxNode48]uint8
	children   [childSizeMaxNode48]*artNode
}

func (n48 *node48) kind() kind {
	return kindNode48
}

func (n48 *node48) findChild(b byte) *artNode {
	index := n48.childIndex[b] - 1
	if index >= childSizeMaxNode48 {
		return nil
	}
	return n48.children[index]
}

func (n48 *node48) addChild(b byte, child artNode) artNode {
	if n48.childSize+1 <= childSizeMaxNode48 {
		for i, child := range n48.children {
			if child == nil {
				n48.childIndex[b] = 1 + uint8(i)
				n48.children[i] = child
				break
			}
		}
		n48.childSize++
		return n48
	}
	n256 := n48.upgrade()
	n256.addChild(b, child)
	return n256
}

func (n48 *node48) removeChild(b byte) artNode {
	index := n48.childIndex[b] - 1
	if index >= childSizeMaxNode48 || n48.children[index] == nil {
		return n48
	}
	n48.childIndex[b] = 0
	n48.children[index] = nil
	n48.childSize--
	if n48.childSize >= childSizeMinNode48 {
		return n48
	}
	return n48.downgrade()
}

func (n48 *node48) anyChild() *artNode {
	for _, child := range n48.children {
		if child != nil {
			return child
		}
	}
	return nil
}

func (n48 *node48) upgrade() (n256 *node256) {
	n256 = &node256{baseNode: n48.baseNode}
	for i := range n48.childIndex {
		child := n48.findChild(uint8(i))
		if child == nil {
			continue
		}
		n256.addChild(uint8(i), *child)
	}
	return
}

func (n256 *node256) downgrade() (n48 *node48) {
	n48 = &node48{}
	for i := 0; i < len(n256.children) && n48.childSize < n256.childSize; i++ {
		child := n256.findChild(byte(i))
		if child == nil {
			continue
		}
		n48.addChild(byte(i), *child)
	}
	return
}

type node256 struct {
	baseNode
	children [childSizeMaxNode256]*artNode
}

func (n256 *node256) kind() kind {
	return kindNode256
}

func (n256 *node256) findChild(b byte) *artNode {
	return n256.children[b]
}

func (n256 *node256) addChild(b byte, child artNode) artNode {
	n256.children[b] = &child
	return n256
}

func (n256 *node256) removeChild(b byte) artNode {
	if n256.children[b] == nil {
		return n256
	}
	n256.children[b] = nil
	n256.childSize--
	if n256.childSize >= childSizeMinNode256 {
		return n256
	}
	return n256.downgrade()
}

func (n256 *node256) anyChild() *artNode {
	for _, child := range n256.children {
		if child != nil {
			return child
		}
	}
	return nil
}
