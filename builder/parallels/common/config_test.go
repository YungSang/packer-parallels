package common

import (
	"testing"

	"github.com/mitchellh/packer/packer"
)

func testConfigTemplate(t *testing.T) *packer.ConfigTemplate {
	result, err := packer.NewConfigTemplate()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result
}
