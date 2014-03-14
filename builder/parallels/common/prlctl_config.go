package common

import (
	"fmt"

	"github.com/mitchellh/packer/packer"
)

type PrlCtlConfig struct {
	PrlCtl [][]string `mapstructure:"prlctl"`
}

func (c *PrlCtlConfig) Prepare(t *packer.ConfigTemplate) []error {
	if c.PrlCtl == nil {
		c.PrlCtl = make([][]string, 0)
	}

	errs := make([]error, 0)
	for i, args := range c.PrlCtl {
		for j, arg := range args {
			if err := t.Validate(arg); err != nil {
				errs = append(errs,
					fmt.Errorf("Error processing prlctl[%d][%d]: %s", i, j, err))
			}
		}
	}

	return errs
}
