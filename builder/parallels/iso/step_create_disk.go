package iso

import (
	"fmt"
	"strconv"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"

	prlcommon "github.com/yungsang/packer-parallels/builder/parallels/common"
)

// This step creates the virtual disk that will be used as the
// hard drive for the virtual machine.
type stepCreateDisk struct{}

func (s *stepCreateDisk) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*config)
	driver := state.Get("driver").(prlcommon.Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	command := []string{
		"set", vmName,
		"--device-set", "hdd0",
		"--size", strconv.FormatUint(uint64(config.DiskSize), 10),
		"--iface", config.HardDriveInterface,
	}

	ui.Say("Creating hard drive...")
	err := driver.PrlCtl(command...)
	if err != nil {
		err := fmt.Errorf("Error creating hard drive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepCreateDisk) Cleanup(state multistep.StateBag) {}
