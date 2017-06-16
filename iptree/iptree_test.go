package iptree

import (
	"net"
	"testing"
)

func TestInsertNet(t *testing.T) {
	r := NewTree()

	newR := r.InsertNet(nil, "test")
	if newR != r {
		t.Errorf("Expected no changes inserting nil network but got:\n%s\n", newR.root.Dot())
	}

	newR = r.InsertNet(&net.IPNet{IP: nil, Mask: nil}, "test")
	if newR != r {
		t.Errorf("Expected no changes inserting invalid network but got:\n%s\n", newR.root.Dot())
	}

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	newR = r.InsertNet(n, "test")
	if newR == r {
		t.Errorf("Expected new root after insertion of new address but got previous")
	}
}

func TestIPv4NetToUint32(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	key, bits := iPv4NetToUint32(n)
	if key != 0xc0000200 || bits != 24 {
		t.Errorf("Expected 0xc0000200, 24 pair but got 0x%08x, %d", key, bits)
	}

	n = &net.IPNet{
		IP:   net.IP{0xc, 0x00, 0x02},
		Mask: net.IPMask{0xff, 0xff, 0xff, 0x00}}
	key, bits = iPv4NetToUint32(n)
	if bits >= 0 {
		t.Errorf("Expected negative number of bits but got %d", bits)
	}

	n = &net.IPNet{
		IP:   net.IP{0xc, 0x00, 0x02, 0x00},
		Mask: net.IPMask{0xff, 0x00, 0xff, 0x00}}
	key, bits = iPv4NetToUint32(n)
	if bits >= 0 {
		t.Errorf("Expected negative number of bits but got %d", bits)
	}
}
