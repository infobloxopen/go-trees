package domaintree

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/infobloxopen/go-trees/dltree"
)

func TestSplit(t *testing.T) {
	dn := ""
	labels := split(dn)
	if len(labels) != 0 {
		t.Errorf("Expected zero labels for empty domain name %q but got %d", dn, len(labels))
	}

	dn = "."
	labels = split(dn)
	if len(labels) != 0 {
		t.Errorf("Expected zero labels for root fqdn %q but got %d", dn, len(labels))
	}

	dn = "www\\.test.com"
	labels = split(dn)
	assertDomainName(labels, []string{
		"com",
		"www\\.test",
	}, dn, t)

	dn = "www.test.com."
	labels = split(dn)
	assertDomainName(labels, []string{
		"com",
		"test",
		"www",
	}, dn, t)
}

func TestGetLabelsCount(t *testing.T) {
	dn := ""
	c := getLabelsCount(dn)
	if c != 0 {
		t.Errorf("Expected zero labels for empty domain name %q but got %d", dn, c)
	}

	dn = "."
	c = getLabelsCount(dn)
	if c != 0 {
		t.Errorf("Expected zero labels for root fqdn %q but got %d", dn, c)
	}

	dn = "www\\.test.com"
	c = getLabelsCount(dn)
	if c != 2 {
		t.Errorf("Expected two labels for domain name %q but got %d", dn, c)
	}

	dn = "www.test.com."
	c = getLabelsCount(dn)
	if c != 3 {
		t.Errorf("Expected three labels for fqdn %q but got %d", dn, c)
	}
}

func assertDomainName(labels []dltree.DomainLabel, elabels []string, dn string, t *testing.T) {
	for i := range elabels {
		elabels[i] += "\n"
	}

	s := make([]string, len(labels))
	for i, label := range labels {
		s[i] = label.String() + "\n"
	}

	ctx := difflib.ContextDiff{
		A:        elabels,
		B:        s,
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("can't compare labels for domain name \"%s\": %s", dn, err))
	}

	if len(diff) > 0 {
		t.Errorf("Labels for domain name \"%s\" don't match:\n%s", dn, diff)
	}
}
