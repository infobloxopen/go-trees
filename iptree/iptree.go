// Package iptree implements radix tree data structure for IPv4 networks.
package iptree

import (
	"net"

	"github.com/infobloxopen/go-trees/numtree"
)

const iPv4Bits = net.IPv4len * 8

// Tree is a radix tree for IPv4 networks.
type Tree struct {
	root *numtree.Node32
}

// NewTree creates empty tree.
func NewTree() *Tree {
	return &Tree{}
}

// InsertNet inserts value using given network as a key.
func (t *Tree) InsertNet(n *net.IPNet, value interface{}) *Tree {
	if n == nil {
		return t
	}

	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		return &Tree{
			root: t.root.Insert(key, bits, value)}
	}

	return t
}

func iPv4NetToUint32(n *net.IPNet) (uint32, int) {
	if len(n.IP) != net.IPv4len {
		return 0, -1
	}

	ones, bits := n.Mask.Size()
	if bits != iPv4Bits {
		return 0, -1
	}

	return (uint32(n.IP[0]) << 24) | (uint32(n.IP[1]) << 16) | (uint32(n.IP[2]) << 8) | uint32(n.IP[3]), ones
}
