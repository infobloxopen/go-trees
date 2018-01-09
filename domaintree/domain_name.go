package domaintree

import "github.com/infobloxopen/go-trees/dltree"

func split(s string) []dltree.DomainLabel {
	dn := make([]dltree.DomainLabel, getLabelsCount(s))
	if len(dn) > 0 {
		end := len(dn) - 1
		start := 0
		for i := range dn {
			label, p := dltree.MakeDomainLabel(s[start:])
			start += p + 1
			dn[end-i] = label
		}
	}

	return dn
}

func getLabelsCount(s string) int {
	labels := 0
	start := 0
	for {
		size, p := dltree.GetFirstLabelSize(s[start:])
		start += p + 1
		if start >= len(s) {
			if size > 0 {
				labels++
			}

			break
		}

		labels++
	}

	return labels
}
