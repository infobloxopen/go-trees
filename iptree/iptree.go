// Package iptree implements radix tree data structure for IPv4 and IPv6 networks.
package iptree

import (
	"fmt"
	"net"

	"github.com/infobloxopen/go-trees/numtree"
)

const (
	iPv4Bits = net.IPv4len * 8
	iPv6Bits = net.IPv6len * 8
)

// Tree is a radix tree for IPv4 and IPv6 networks.
type Tree struct {
	root32 *numtree.Node32
	root64 *numtree.Node64
}

type subTree64 *numtree.Node64

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
			root32: t.root32.Insert(key, bits, value),
			root64: t.root64}
	}

	if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		if MSBits < numtree.Key64BitSize {
			return &Tree{
				root32: t.root32,
				root64: t.root64.Insert(MSKey, MSBits, value)}
		}

		var r *numtree.Node64
		if v, ok := t.root64.ExactMatch(MSKey, MSBits); ok {
			s, ok := v.(subTree64)
			if !ok {
				err := fmt.Errorf("invalid IPv6 tree: expected subTree64 value at 0x%016x, %d but got %T (%#v)",
					MSKey, MSBits, v, v)
				panic(err)
			}

			r = (*numtree.Node64)(s)
		}

		r = r.Insert(LSKey, LSBits, value)
		return &Tree{
			root32: t.root32,
			root64: t.root64.Insert(MSKey, MSBits, subTree64(r))}
	}

	return t
}

// GetByNet gets value for network which is equal to or contains given network.
func (t *Tree) GetByNet(n *net.IPNet) (interface{}, bool) {
	if n == nil {
		return nil, false
	}

	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		return t.root32.Match(key, bits)
	}

	if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		v, ok := t.root64.Match(MSKey, MSBits)
		if !ok || MSBits < numtree.Key64BitSize {
			return v, ok
		}

		s, ok := v.(subTree64)
		if !ok {
			return v, true
		}

		v, ok = (*numtree.Node64)(s).Match(LSKey, LSBits)
		if ok {
			return v, ok
		}

		return t.root64.Match(MSKey, numtree.Key64BitSize-1)
	}

	return nil, false
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

func iPv6NetToUint64Pair(n *net.IPNet) (uint64, int, uint64, int) {
	if len(n.IP) != net.IPv6len {
		return 0, -1, 0, -1
	}

	ones, bits := n.Mask.Size()
	if bits != iPv6Bits {
		return 0, -1, 0, -1
	}

	MSBits := numtree.Key64BitSize
	LSBits := 0
	if ones > numtree.Key64BitSize {
		LSBits = ones - numtree.Key64BitSize
	} else {
		MSBits = ones
	}

	return packIPToUint64(n.IP), MSBits, packIPToUint64(n.IP[8:]), LSBits
}

func packIPToUint64(x net.IP) uint64 {
	return (uint64(x[0]) << 56) | (uint64(x[1]) << 48) | (uint64(x[2]) << 40) | (uint64(x[3]) << 32) |
		(uint64(x[4]) << 24) | (uint64(x[5]) << 16) | (uint64(x[6]) << 8) | uint64(x[7])
}
