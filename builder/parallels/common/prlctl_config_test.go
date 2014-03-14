package common

import (
	"reflect"
	"testing"
)

func TestPrlCtlConfigPrepare_PrlCtl(t *testing.T) {
	// Test with empty
	c := new(PrlCtlConfig)
	errs := c.Prepare(testConfigTemplate(t))
	if len(errs) > 0 {
		t.Fatalf("err: %#v", errs)
	}

	if !reflect.DeepEqual(c.PrlCtl, [][]string{}) {
		t.Fatalf("bad: %#v", c.PrlCtl)
	}

	// Test with a good one
	c = new(PrlCtlConfig)
	c.PrlCtl = [][]string{
		{"foo", "bar", "baz"},
	}
	errs = c.Prepare(testConfigTemplate(t))
	if len(errs) > 0 {
		t.Fatalf("err: %#v", errs)
	}

	expected := [][]string{
		[]string{"foo", "bar", "baz"},
	}

	if !reflect.DeepEqual(c.PrlCtl, expected) {
		t.Fatalf("bad: %#v", c.PrlCtl)
	}
}
