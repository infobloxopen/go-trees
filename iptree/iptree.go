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

var (
	iPv4MaxMask = net.CIDRMask(iPv4Bits, iPv4Bits)
	iPv6MaxMask = net.CIDRMask(iPv6Bits, iPv6Bits)
)

var (
	masks32 = []uint32{
		0x00000000, 0x80000000, 0xc0000000, 0xe0000000,
		0xf0000000, 0xf8000000, 0xfc000000, 0xfe000000,
		0xff000000, 0xff800000, 0xffc00000, 0xffe00000,
		0xfff00000, 0xfff80000, 0xfffc0000, 0xfffe0000,
		0xffff0000, 0xffff8000, 0xffffc000, 0xffffe000,
		0xfffff000, 0xfffff800, 0xfffffc00, 0xfffffe00,
		0xffffff00, 0xffffff80, 0xffffffc0, 0xffffffe0,
		0xfffffff0, 0xfffffff8, 0xfffffffc, 0xfffffffe,
		0xffffffff}

	masks64 = []uint64{
		0x0000000000000000, 0x8000000000000000, 0xc000000000000000, 0xe000000000000000,
		0xf000000000000000, 0xf800000000000000, 0xfc00000000000000, 0xfe00000000000000,
		0xff00000000000000, 0xff80000000000000, 0xffc0000000000000, 0xffe0000000000000,
		0xfff0000000000000, 0xfff8000000000000, 0xfffc000000000000, 0xfffe000000000000,
		0xffff000000000000, 0xffff800000000000, 0xffffc00000000000, 0xffffe00000000000,
		0xfffff00000000000, 0xfffff80000000000, 0xfffffc0000000000, 0xfffffe0000000000,
		0xffffff0000000000, 0xffffff8000000000, 0xffffffc000000000, 0xffffffe000000000,
		0xfffffff000000000, 0xfffffff800000000, 0xfffffffc00000000, 0xfffffffe00000000,
		0xffffffff00000000, 0xffffffff80000000, 0xffffffffc0000000, 0xffffffffe0000000,
		0xfffffffff0000000, 0xfffffffff8000000, 0xfffffffffc000000, 0xfffffffffe000000,
		0xffffffffff000000, 0xffffffffff800000, 0xffffffffffc00000, 0xffffffffffe00000,
		0xfffffffffff00000, 0xfffffffffff80000, 0xfffffffffffc0000, 0xfffffffffffe0000,
		0xffffffffffff0000, 0xffffffffffff8000, 0xffffffffffffc000, 0xffffffffffffe000,
		0xfffffffffffff000, 0xfffffffffffff800, 0xfffffffffffffc00, 0xfffffffffffffe00,
		0xffffffffffffff00, 0xffffffffffffff80, 0xffffffffffffffc0, 0xffffffffffffffe0,
		0xfffffffffffffff0, 0xfffffffffffffff8, 0xfffffffffffffffc, 0xfffffffffffffffe,
		0xffffffffffffffff}
)

// Tree is a radix tree for IPv4 and IPv6 networks.
type Tree struct {
	root32 *numtree.Node32
	root64 *numtree.Node64
}

// Pair represents a key-value pair returned by Enumerate method.
type Pair struct {
	Key   *net.IPNet
	Value interface{}
}

type subTree64 *numtree.Node64

// NewTree creates empty tree.
func NewTree() *Tree {
	return &Tree{}
}

// InsertNet inserts value using given network as a key. The method returns new tree (old one remains unaffected).
func (t *Tree) InsertNet(n *net.IPNet, value interface{}) *Tree {
	if n == nil {
		return t
	}

	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		var (
			r32 *numtree.Node32
			r64 *numtree.Node64
		)

		if t != nil {
			r32 = t.root32
			r64 = t.root64
		}

		return &Tree{root32: r32.Insert(key, bits, value), root64: r64}
	}

	if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		var (
			r32 *numtree.Node32
			r64 *numtree.Node64
		)

		if t != nil {
			r32 = t.root32
			r64 = t.root64
		}

		if MSBits < numtree.Key64BitSize {
			return &Tree{root32: r32, root64: r64.Insert(MSKey, MSBits, value)}
		}

		var r *numtree.Node64
		if v, ok := r64.ExactMatch(MSKey, MSBits); ok {
			s, ok := v.(subTree64)
			if !ok {
				err := fmt.Errorf("invalid IPv6 tree: expected subTree64 value at 0x%016x, %d but got %T (%#v)",
					MSKey, MSBits, v, v)
				panic(err)
			}

			r = (*numtree.Node64)(s)
		}

		r = r.Insert(LSKey, LSBits, value)
		return &Tree{root32: r32, root64: r64.Insert(MSKey, MSBits, subTree64(r))}
	}

	return t
}

type UpdateDescendantsCallback func(Pair) (interface{}, bool)

func updateDescendantsIPv4(n *numtree.Node32, isRoot bool, callback UpdateDescendantsCallback) {
	// we skip the first (root) node
	if !isRoot {
		mask := net.CIDRMask(int(n.Bits), iPv4Bits)
		key := &net.IPNet{IP: unpackUint32ToIP(n.Key).Mask(mask), Mask: mask}

		newValue, shouldUpdate := callback(Pair{Key: key, Value: n.Value})
		if shouldUpdate {
			n.Value = newValue
		}
	}

	c1, c2 := n.Children()
	if c1 != nil {
		updateDescendantsIPv4(c1, false, callback)
	}
	if c2 != nil {
		updateDescendantsIPv4(c2, false, callback)
	}
}

func updateDescendantsIPv6(n *numtree.Node64, callback UpdateDescendantsCallback, mkey uint64, mbits int) {
	// if big network, i.e. bits < 64
	if mbits < numtree.Key64BitSize {
		MSIP := append(unpackUint64ToIP(n.Key), make(net.IP, 8)...)

		// if the value contains another tree
		if s, ok := n.Value.(subTree64); ok {
			LSIP := unpackUint64ToIP(n.Key)
			mask := net.CIDRMask(numtree.Key64BitSize+int(s.Bits), iPv6Bits)
			key := &net.IPNet{IP: append(MSIP[0:8], LSIP...).Mask(mask), Mask: mask}

			if s.Value != nil {
				newValue, shouldUpdate := callback(Pair{Key: key, Value: s.Value})
				if shouldUpdate {
					s.Value = newValue
				}
			}

			c1, c2 := (*numtree.Node64)(s).Children()
			if c1 != nil {
				updateDescendantsIPv6(c1, callback, n.Key, int(n.Bits))
			}
			if c2 != nil {
				updateDescendantsIPv6(c2, callback, n.Key, int(n.Bits))
			}
			return
		}

		if n.Value != nil {
			mask := net.CIDRMask(int(n.Bits), iPv6Bits)
			key := &net.IPNet{IP: MSIP.Mask(mask), Mask: mask}
			newValue, shouldUpdate := callback(Pair{Key: key, Value: n.Value})
			if shouldUpdate {
				n.Value = newValue
			}
		}

		c1, c2 := n.Children()

		if c1 != nil {
			updateDescendantsIPv6(c1, callback, n.Key, int(n.Bits))
		}
		if c2 != nil {
			updateDescendantsIPv6(c2, callback, n.Key, int(n.Bits))
		}
		return
	}

	// for smaller networks, i.e. bits >= 64
	MSIP := append(unpackUint64ToIP(mkey), make(net.IP, 8)...)
	LSIP := unpackUint64ToIP(n.Key)
	mask := net.CIDRMask(numtree.Key64BitSize+int(n.Bits), iPv6Bits)
	key := &net.IPNet{IP: append(MSIP[0:8], LSIP...).Mask(mask), Mask: mask}

	if n.Value != nil {
		newValue, shouldUpdate := callback(Pair{Key: key, Value: n.Value})
		if shouldUpdate {
			n.Value = newValue
		}
	}

	c1, c2 := n.Children()
	if c1 != nil {
		updateDescendantsIPv6(c1, callback, mkey, mbits)
	}
	if c2 != nil {
		updateDescendantsIPv6(c2, callback, mkey, mbits)
	}
}

func (t *Tree) UpdateDescendants(n *net.IPNet, callback UpdateDescendantsCallback) {
	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		r := t.root32.FindNode(key, bits)
		if r == nil {
			return
		}

		updateDescendantsIPv4(r, true, callback)
		return
	}

	if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		r := t.root64.FindNode(MSKey, MSBits)
		if r == nil {
			return
		}

		if MSBits < numtree.Key64BitSize {
			c1, c2 := r.Children()
			if c1 != nil {
				updateDescendantsIPv6(c1, callback, MSKey, MSBits)
			}
			if c2 != nil {
				updateDescendantsIPv6(c2, callback, MSKey, MSBits)
			}
			return
		}

		s, ok := r.Value.(subTree64)
		if !ok {
			err := fmt.Errorf("invalid IPv6 tree: expected subTree64 value at 0x%016x, %d but got %T (%#v)",
				MSKey, MSBits, r, r)
			panic(err)
		}

		r2 := (*numtree.Node64)(s).FindNode(LSKey, LSBits)

		if r2 == nil {
			return
		}
		c1, c2 := r2.Children()
		if c1 != nil {
			updateDescendantsIPv6(c1, callback, MSKey, MSBits)
		}
		if c2 != nil {
			updateDescendantsIPv6(c2, callback, MSKey, MSBits)
		}
	}
}

// InplaceInsertNet inserts (or replaces) value using given network as a key in current tree.
func (t *Tree) InplaceInsertNet(n *net.IPNet, value interface{}) {
	if n == nil {
		return
	}

	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		t.root32 = t.root32.InplaceInsert(key, bits, value)
	} else if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		if MSBits < numtree.Key64BitSize {
			t.root64 = t.root64.InplaceInsert(MSKey, MSBits, value)
		} else {
			if v, ok := t.root64.ExactMatch(MSKey, MSBits); ok {
				s, ok := v.(subTree64)
				if !ok {
					err := fmt.Errorf("invalid IPv6 tree: expected subTree64 value at 0x%016x, %d but got %T (%#v)",
						MSKey, MSBits, v, v)
					panic(err)
				}

				r := (*numtree.Node64)(s)
				newR := r.InplaceInsert(LSKey, LSBits, value)
				if newR != r {
					t.root64 = t.root64.InplaceInsert(MSKey, MSBits, subTree64(newR))
				}
			} else {
				var r *numtree.Node64
				r = r.InplaceInsert(LSKey, LSBits, value)
				t.root64 = t.root64.InplaceInsert(MSKey, MSBits, subTree64(r))
			}
		}
	}
}

// InplaceInsertNetCheckChildren does the same as InplaceInsertNet but additionally returns a boolean,
// which is true in case the new network being added already exists in the tree or if it becomes the parent node of an
// already existing node
func (t *Tree) InplaceInsertNetCheckChildren(n *net.IPNet, value interface{}) bool {
	var hasChildren bool
	if n == nil {
		return hasChildren
	}

	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		t.root32, hasChildren = t.root32.InplaceInsertWithHierarchyChange(key, bits, value)
		return hasChildren
	} else if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		if MSBits < numtree.Key64BitSize {
			t.root64, hasChildren = t.root64.InplaceInsertWithHierarchyChange(MSKey, MSBits, value)
			return hasChildren
		} else {
			if v, ok := t.root64.ExactMatch(MSKey, MSBits); ok {
				s, ok := v.(subTree64)
				if !ok {
					err := fmt.Errorf("invalid IPv6 tree: expected subTree64 value at 0x%016x, %d but got %T (%#v)",
						MSKey, MSBits, v, v)
					panic(err)
				}

				r := (*numtree.Node64)(s)
				var newR *numtree.Node64
				newR, hasChildren = r.InplaceInsertWithHierarchyChange(LSKey, LSBits, value)
				if newR != r {
					t.root64, _ = t.root64.InplaceInsertWithHierarchyChange(MSKey, MSBits, subTree64(newR))
				}
			} else {
				var r *numtree.Node64
				r, _ = r.InplaceInsertWithHierarchyChange(LSKey, LSBits, value)
				t.root64, hasChildren = t.root64.InplaceInsertWithHierarchyChange(MSKey, MSBits, subTree64(r))
			}
		}
	}
	return hasChildren
}

// InsertIP inserts value using given IP address as a key. The method returns new tree (old one remains unaffected).
func (t *Tree) InsertIP(ip net.IP, value interface{}) *Tree {
	return t.InsertNet(newIPNetFromIP(ip), value)
}

// InplaceInsertIP inserts (or replaces) value using given IP address as a key in current tree.
func (t *Tree) InplaceInsertIP(ip net.IP, value interface{}) {
	t.InplaceInsertNet(newIPNetFromIP(ip), value)
}

// Enumerate returns channel which is populated by key-value pairs of tree content.
func (t *Tree) Enumerate() chan Pair {
	ch := make(chan Pair)

	go func() {
		defer close(ch)

		if t == nil {
			return
		}

		t.enumerate(ch)
	}()

	return ch
}

// GetByNet gets value for network which is equal to or contains given network.
func (t *Tree) GetByNet(n *net.IPNet) (interface{}, bool) {
	if t == nil || n == nil {
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

// GetByIP gets value for network which is equal to or contains given IP address.
func (t *Tree) GetByIP(ip net.IP) (interface{}, bool) {
	return t.GetByNet(newIPNetFromIP(ip))
}

// DeleteByNet removes subtree which is contained by given network. The method returns new tree (old one remains unaffected) and flag indicating if deletion happens indeed.
func (t *Tree) DeleteByNet(n *net.IPNet) (*Tree, bool) {
	if t == nil || n == nil {
		return t, false
	}

	if key, bits := iPv4NetToUint32(n); bits >= 0 {
		r, ok := t.root32.Delete(key, bits)
		if ok {
			return &Tree{root32: r, root64: t.root64}, true
		}
	} else if MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n); MSBits >= 0 {
		r64 := t.root64
		if MSBits < numtree.Key64BitSize {
			r64, ok := r64.Delete(MSKey, MSBits)
			if ok {
				return &Tree{root32: t.root32, root64: r64}, true
			}
		} else if v, ok := r64.ExactMatch(MSKey, MSBits); ok {
			s, ok := v.(subTree64)
			if !ok {
				err := fmt.Errorf("invalid IPv6 tree: expected subTree64 value at 0x%016x, %d but got %T (%#v)",
					MSKey, MSBits, v, v)
				panic(err)
			}

			r, ok := (*numtree.Node64)(s).Delete(LSKey, LSBits)
			if ok {
				if r == nil {
					r64, _ = r64.Delete(MSKey, MSBits)
				} else {
					r64 = r64.Insert(MSKey, MSBits, subTree64(r))
				}

				return &Tree{root32: t.root32, root64: r64}, true
			}
		}
	}

	return t, false
}

// DeleteByIP removes node by given IP address. The method returns new tree (old one remains unaffected) and flag indicating if deletion happens indeed.
func (t *Tree) DeleteByIP(ip net.IP) (*Tree, bool) {
	return t.DeleteByNet(newIPNetFromIP(ip))
}

func (t *Tree) enumerate(ch chan Pair) {
	for n := range t.root32.Enumerate() {
		mask := net.CIDRMask(int(n.Bits), iPv4Bits)
		ch <- Pair{
			Key: &net.IPNet{
				IP:   unpackUint32ToIP(n.Key).Mask(mask),
				Mask: mask},
			Value: n.Value}
	}

	for n := range t.root64.Enumerate() {
		MSIP := append(unpackUint64ToIP(n.Key), make(net.IP, 8)...)
		if s, ok := n.Value.(subTree64); ok {
			for n := range (*numtree.Node64)(s).Enumerate() {
				LSIP := unpackUint64ToIP(n.Key)
				mask := net.CIDRMask(numtree.Key64BitSize+int(n.Bits), iPv6Bits)
				ch <- Pair{
					Key: &net.IPNet{
						IP:   append(MSIP[0:8], LSIP...).Mask(mask),
						Mask: mask},
					Value: n.Value}
			}
		} else {
			mask := net.CIDRMask(int(n.Bits), iPv6Bits)
			ch <- Pair{
				Key: &net.IPNet{
					IP:   MSIP.Mask(mask),
					Mask: mask},
				Value: n.Value}
		}
	}
}

func iPv4NetToUint32(n *net.IPNet) (uint32, int) {
	if len(n.IP) != net.IPv4len {
		return 0, -1
	}

	ones, bits := n.Mask.Size()
	if bits != iPv4Bits {
		return 0, -1
	}

	return packIPToUint32(n.IP), ones
}

func packIPToUint32(x net.IP) uint32 {
	return (uint32(x[0]) << 24) | (uint32(x[1]) << 16) | (uint32(x[2]) << 8) | uint32(x[3])
}

func unpackUint32ToIP(x uint32) net.IP {
	return net.IP{byte(x >> 24 & 0xff), byte(x >> 16 & 0xff), byte(x >> 8 & 0xff), byte(x & 0xff)}
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

func unpackUint64ToIP(x uint64) net.IP {
	return net.IP{
		byte(x >> 56 & 0xff),
		byte(x >> 48 & 0xff),
		byte(x >> 40 & 0xff),
		byte(x >> 32 & 0xff),
		byte(x >> 24 & 0xff),
		byte(x >> 16 & 0xff),
		byte(x >> 8 & 0xff),
		byte(x & 0xff)}
}

func newIPNetFromIP(ip net.IP) *net.IPNet {
	if ip4 := ip.To4(); ip4 != nil {
		return &net.IPNet{IP: ip4, Mask: iPv4MaxMask}
	}

	if ip6 := ip.To16(); ip6 != nil {
		return &net.IPNet{IP: ip6, Mask: iPv6MaxMask}
	}

	return nil
}
