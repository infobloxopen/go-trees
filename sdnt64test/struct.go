package main

type pair struct {
	k string
	v uint64
}

type mapper64 interface {
	Map(s string) uint64
}
