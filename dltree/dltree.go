// Package dltree implements red-black tree for key value pairs with domain label keys.
package dltree

import "github.com/infobloxopen/go-trees/domain"

// Tree is a red-black tree for key-value pairs where key is domain label.
type Tree struct {
	root *node
}

// Pair is a key-value pair representing tree node content.
type Pair struct {
	Key   string
	Value interface{}
}

// RawPair is a key-value pair representing tree node raw content.
type RawPair struct {
	Key   string
	Size  int
	Value interface{}
}

// NewTree creates empty tree.
func NewTree() *Tree {
	return new(Tree)
}

// Insert puts given key-value pair to the tree and returns pointer to new root.
func (t *Tree) Insert(key string, value interface{}) *Tree {
	var (
		n *node
	)

	if t != nil {
		n = t.root
	}

	dl, size, _ := domain.MakeLabel(key)
	return &Tree{root: n.insert(dl, size, value)}
}

// RawInsert puts given key-value pair to the tree and returns pointer to new root. Expects bindary domain label on input.
func (t *Tree) RawInsert(key string, size int, value interface{}) *Tree {
	var (
		n *node
	)

	if t != nil {
		n = t.root
	}

	return &Tree{root: n.insert(key, size, value)}
}

// InplaceInsert inserts or replaces given key-value pair in the tree. The method inserts data directly to current tree so make sure you have exclusive access to it.
func (t *Tree) InplaceInsert(key string, value interface{}) {
	dl, size, _ := domain.MakeLabel(key)
	t.root = t.root.inplaceInsert(dl, size, value)
}

// RawInplaceInsert inserts or replaces given key-value pair in the tree. The method inserts data directly to current tree so make sure you have exclusive access to it. Expects bindary domain label on input.
func (t *Tree) RawInplaceInsert(key string, size int, value interface{}) {
	t.root = t.root.inplaceInsert(key, size, value)
}

// Get returns value by given key.
func (t *Tree) Get(key string) (interface{}, bool) {
	if t == nil {
		return nil, false
	}

	dl, size, _ := domain.MakeLabel(key)
	return t.root.get(dl, size)
}

// RawGet returns value by given key. Expects bindary domain label on input.
func (t *Tree) RawGet(key string, size int) (interface{}, bool) {
	if t == nil {
		return nil, false
	}

	return t.root.get(key, size)
}

// Enumerate returns channel which is populated by key pair values in order of keys.
func (t *Tree) Enumerate() chan Pair {
	ch := make(chan Pair)

	go func() {
		defer close(ch)

		if t == nil {
			return
		}

		t.root.enumerate(ch)
	}()

	return ch
}

// RawEnumerate returns channel which is populated by key pair values in order of keys. Returns binary domain labels.
func (t *Tree) RawEnumerate() chan RawPair {
	ch := make(chan RawPair)

	go func() {
		defer close(ch)

		if t == nil {
			return
		}

		t.root.rawEnumerate(ch)
	}()

	return ch
}

// Delete removes node by given key. It returns copy of tree and true if node has been indeed deleted otherwise copy of tree and false.
func (t *Tree) Delete(key string) (*Tree, bool) {
	if t == nil {
		return nil, false
	}

	dl, size, _ := domain.MakeLabel(key)
	root, ok := t.root.del(dl, size)
	return &Tree{root: root}, ok
}

// RawDelete removes node by given key. It returns copy of tree and true if node has been indeed deleted otherwise copy of tree and false. Expects bindary domain label on input.
func (t *Tree) RawDelete(key string, size int) (*Tree, bool) {
	if t == nil {
		return nil, false
	}

	root, ok := t.root.del(key, size)
	return &Tree{root: root}, ok
}

// IsEmpty returns true if given tree has no nodes.
func (t *Tree) IsEmpty() bool {
	return t == nil || t.root == nil
}

// Dot dumps tree to Graphviz .dot format.
func (t *Tree) Dot() string {
	body := ""

	if t != nil {
		body = t.root.dot()
	}

	return "digraph d {\n" + body + "}\n"
}
