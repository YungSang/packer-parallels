package iso

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"

	prlcommon "github.com/yungsang/packer-parallels/builder/parallels/common"
)

// This step attaches the Parallels Tools as a inserted CD onto
// the virtual machine.
//
// Uses:
//   config *config
//   driver Driver
//   ui packer.Ui
//   vmName string
//
// Produces:
type stepInstallParallelsTools struct{}

func (stepInstallParallelsTools) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*config)
	driver := state.Get("driver").(prlcommon.Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	// If we're not attaching the guest additions then just return
	if config.ParallelsToolsMode != ParallelsToolsModeInstall {
		log.Println("Not installing Parallels Tools since it is disabled.")
		return multistep.ActionContinue
	}

	// Install the Parallels Tools into the VM
	ui.Say("Installing Parallels Tools into the VM...")
	command := []string{
		"installtools", vmName,
	}
	if err := driver.PrlCtl(command...); err != nil {
		err := fmt.Errorf("Error installing Parallels Tools: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (stepInstallParallelsTools) Cleanup(multistep.StateBag) {}
