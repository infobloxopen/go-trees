// Package domaintree implements radix tree data structure for domain names.
package domaintree

import (
	"strings"

	"github.com/infobloxopen/go-trees/strtree"
)

// Node is a radix tree for domain names.
type Node struct {
	branches *strtree.Tree

	hasValue bool
	value    interface{}
}

// Pair represents a key-value pair returned by Enumerate method.
type Pair struct {
	Key   string
	Value interface{}
}

// Insert puts value using given domain as a key. The method returns new tree (old one remains unaffected).
func (n *Node) Insert(d string, v interface{}) *Node {
	if n == nil {
		n = &Node{}
	} else {
		n = &Node{
			branches: n.branches,
			hasValue: n.hasValue,
			value:    n.value}
	}
	r := n

	labels := strings.Split(d, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		label := labels[i]

		item, ok := n.branches.Get(label)
		var next *Node
		if ok {
			next = item.(*Node)
			next = &Node{
				branches: next.branches,
				hasValue: next.hasValue,
				value:    next.value}
		} else {
			next = &Node{}
		}

		n.branches = n.branches.Insert(label, next)
		n = next
	}

	n.hasValue = true
	n.value = v

	return r
}

// Enumerate returns key-value pairs in given tree sorted by key first by top level domain label second by second level and so on.
func (n *Node) Enumerate() chan Pair {
	ch := make(chan Pair)

	go func() {
		defer close(ch)
		n.enumerate("", ch)
	}()

	return ch
}

// Get gets value for domain which is equal to domain in the tree or is a subdomain of existing domain.
func (n *Node) Get(d string) (interface{}, bool) {
	if n == nil {
		return nil, false
	}

	labels := strings.Split(d, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		label := labels[i]

		item, ok := n.branches.Get(label)
		if !ok {
			break
		}

		n = item.(*Node)
	}

	return n.value, n.hasValue
}

func (n *Node) enumerate(s string, ch chan Pair) {
	if n == nil {
		return
	}

	if n.hasValue {
		ch <- Pair{
			Key:   s,
			Value: n.value}
	}

	for item := range n.branches.Enumerate() {
		sub := item.Key
		if len(s) > 0 {
			sub += "." + s
		}
		node := item.Value.(*Node)

		node.enumerate(sub, ch)
	}
}
