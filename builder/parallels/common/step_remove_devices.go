package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// This step removes any devices (floppy disks, ISOs, etc.) from the
// machine that we may have added.
//
// Uses:
//   driver Driver
//   ui packer.Ui
//   vmName string
//
// Produces:
type StepRemoveDevices struct{}

func (s *StepRemoveDevices) Run(state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	if _, ok := state.GetOk("attachedIso"); ok {
		command := []string{
			"set", vmName,
			"--device-set", "cdrom0",
			"--disable", "--disconnect",
		}

		if err := driver.PrlCtl(command...); err != nil {
			err := fmt.Errorf("Error detaching ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepRemoveDevices) Cleanup(state multistep.StateBag) {
}
