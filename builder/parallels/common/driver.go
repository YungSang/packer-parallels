package common

import (
	"log"
	"os/exec"

	"github.com/mitchellh/multistep"
)

// A driver is able to talk to Parallels and perform certain
// operations with it. Some of the operations on here may seem overly
// specific, but they were built specifically in mind to handle features
// of the Parallels builder for Packer, and to abstract differences in
// versions out of the builder steps, so sometimes the methods are
// extremely specific.
type Driver interface {
	// Create a SATA controller.
	//CreateSATAController(vm string, controller string) error

	// Delete a VM by name
	Delete(string) error

	// Import a VM
	Import(string, string) error

	// Checks if the VM with the given name is running.
	IsRunning(string) (bool, error)

	// Stop stops a running machine, forcefully.
	Stop(string) error

	// Use default settings from Parallels
	UseDefaults(string) error

	// PrlCtl executes the given prlctl command
	PrlCtl(...string) error

	// SSHAddress returns the SSH address for the VM that is being
	// managed by this driver.
	SSHAddressFunc(*SSHConfig) func(multistep.StateBag) (string, error)

	// Get the path to the DHCP leases file for the given device.
	DhcpLeasesPath() string

	// Verify checks to make sure that this driver should function
	// properly. If there is any indication the driver can't function,
	// this will return an error.
	Verify() error

	// Version reads the version of Parallels that is installed.
	Version() (string, error)
}

func NewDriver() (Driver, error) {
	var prlCtlPath string

	if prlCtlPath == "" {
		var err error
		prlCtlPath, err = exec.LookPath("prlctl")
		if err != nil {
			return nil, err
		}
	}

	log.Printf("prlctl path: %s", prlCtlPath)
	driver := &Parallels8Driver{prlCtlPath}
	if err := driver.Verify(); err != nil {
		return nil, err
	}

	return driver, nil
}
