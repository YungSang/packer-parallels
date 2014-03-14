package iso

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"

	prlcommon "github.com/yungsang/packer-parallels/builder/parallels/common"
)

const BuilderId = "yungsang.parallels"

type Builder struct {
	config config
	runner multistep.Runner
}

type config struct {
	common.PackerConfig `mapstructure:",squash"`

	prlcommon.OutputConfig   `mapstructure:",squash"`
	prlcommon.RunConfig      `mapstructure:",squash"`
	prlcommon.ShutdownConfig `mapstructure:",squash"`
	prlcommon.SSHConfig      `mapstructure:",squash"`
	prlcommon.PrlCtlConfig   `mapstructure:",squash"`

	DiskSize            uint     `mapstructure:"disk_size"`
	GuestOSType         string   `mapstructure:"guest_os_type"`
	GuestOSDistribution string   `mapstructure:"guest_os_distribution"`
	HardDriveInterface  string   `mapstructure:"hard_drive_interface"`
	ISOChecksum         string   `mapstructure:"iso_checksum"`
	ISOChecksumType     string   `mapstructure:"iso_checksum_type"`
	ISOUrls             []string `mapstructure:"iso_urls"`
	VMName              string   `mapstructure:"vm_name"`
	BootCommand         []string `mapstructure:"boot_command"`

	RawSingleISOUrl string `mapstructure:"iso_url"`

	tpl *packer.ConfigTemplate
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	md, err := common.DecodeConfig(&b.config, raws...)
	if err != nil {
		return nil, err
	}

	b.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, err
	}
	b.config.tpl.UserVars = b.config.PackerUserVars

	// Accumulate any errors and warnings
	errs := common.CheckUnusedConfig(md)
	errs = packer.MultiErrorAppend(
		errs, b.config.OutputConfig.Prepare(b.config.tpl, &b.config.PackerConfig)...)
	errs = packer.MultiErrorAppend(errs, b.config.RunConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.ShutdownConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.SSHConfig.Prepare(b.config.tpl)...)
	errs = packer.MultiErrorAppend(errs, b.config.PrlCtlConfig.Prepare(b.config.tpl)...)

	warnings := make([]string, 0)

	if b.config.DiskSize == 0 {
		b.config.DiskSize = 40000
	}

	if b.config.GuestOSType == "" && b.config.GuestOSDistribution == "" {
		b.config.GuestOSType = "other"
	}

	if b.config.HardDriveInterface == "" {
		b.config.HardDriveInterface = "ide"
	}

	if b.config.VMName == "" {
		b.config.VMName = fmt.Sprintf("packer-%s", b.config.PackerBuildName)
	}

	// Errors
	templates := map[string]*string{
		"guest_os_type":         &b.config.GuestOSType,
		"guest_os_distribution": &b.config.GuestOSDistribution,
		"hard_drive_interface":  &b.config.HardDriveInterface,
		"iso_checksum":          &b.config.ISOChecksum,
		"iso_checksum_type":     &b.config.ISOChecksumType,
		"iso_url":               &b.config.RawSingleISOUrl,
		"vm_name":               &b.config.VMName,
	}

	for n, ptr := range templates {
		var err error
		*ptr, err = b.config.tpl.Process(*ptr, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error processing %s: %s", n, err))
		}
	}

	for i, url := range b.config.ISOUrls {
		var err error
		b.config.ISOUrls[i], err = b.config.tpl.Process(url, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error processing iso_urls[%d]: %s", i, err))
		}
	}

	for i, command := range b.config.BootCommand {
		if err := b.config.tpl.Validate(command); err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("Error processing boot_command[%d]: %s", i, err))
		}
	}

	if b.config.HardDriveInterface != "ide" && b.config.HardDriveInterface != "sata" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("hard_drive_interface can only be ide or sata"))
	}

	if b.config.ISOChecksumType == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("The iso_checksum_type must be specified."))
	} else {
		b.config.ISOChecksumType = strings.ToLower(b.config.ISOChecksumType)
		if b.config.ISOChecksumType != "none" {
			if b.config.ISOChecksum == "" {
				errs = packer.MultiErrorAppend(
					errs, errors.New("Due to large file sizes, an iso_checksum is required"))
			} else {
				b.config.ISOChecksum = strings.ToLower(b.config.ISOChecksum)
			}

			if h := common.HashForType(b.config.ISOChecksumType); h == nil {
				errs = packer.MultiErrorAppend(
					errs,
					fmt.Errorf("Unsupported checksum type: %s", b.config.ISOChecksumType))
			}
		}
	}

	if b.config.RawSingleISOUrl == "" && len(b.config.ISOUrls) == 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("One of iso_url or iso_urls must be specified."))
	} else if b.config.RawSingleISOUrl != "" && len(b.config.ISOUrls) > 0 {
		errs = packer.MultiErrorAppend(
			errs, errors.New("Only one of iso_url or iso_urls may be specified."))
	} else if b.config.RawSingleISOUrl != "" {
		b.config.ISOUrls = []string{b.config.RawSingleISOUrl}
	}

	for i, url := range b.config.ISOUrls {
		b.config.ISOUrls[i], err = common.DownloadableURL(url)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Failed to parse iso_url %d: %s", i+1, err))
		}
	}

	// Warnings
	if b.config.ISOChecksumType == "none" {
		warnings = append(warnings,
			"A checksum type of 'none' was specified. Since ISO files are so big,\n"+
				"a checksum is highly recommended.")
	}

	if b.config.ShutdownCommand == "" {
		warnings = append(warnings,
			"A shutdown_command was not specified. Without a shutdown command, Packer\n"+
				"will forcibly halt the virtual machine, which may result in data loss.")
	}

	if errs != nil && len(errs.Errors) > 0 {
		return warnings, errs
	}

	return warnings, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Create the driver that we'll use to communicate with Parallels
	driver, err := prlcommon.NewDriver()
	if err != nil {
		return nil, fmt.Errorf("Failed creating Parallels driver: %s", err)
	}

	steps := []multistep.Step{
		&common.StepDownload{
			Checksum:     b.config.ISOChecksum,
			ChecksumType: b.config.ISOChecksumType,
			Description:  "ISO",
			ResultKey:    "iso_path",
			Url:          b.config.ISOUrls,
		},
		&prlcommon.StepOutputDir{
			Force: b.config.PackerForce,
			Path:  b.config.OutputDir,
		},
		new(stepCreateVM),
		new(stepCreateDisk),
		new(stepAttachISO),
		&prlcommon.StepPrlCtl{
			Commands: b.config.PrlCtl,
			Tpl:      b.config.tpl,
		},
		&prlcommon.StepRun{
			BootWait: b.config.BootWait,
		},
		&common.StepConnectSSH{
			SSHAddress:     driver.SSHAddressFunc(&b.config.SSHConfig),
			SSHConfig:      prlcommon.SSHConfigFunc(&b.config.SSHConfig),
			SSHWaitTimeout: b.config.SSHWaitTimeout,
		},
		new(common.StepProvision),
		&prlcommon.StepShutdown{
			Command: b.config.ShutdownCommand,
			Timeout: b.config.ShutdownTimeout,
		},
		&prlcommon.StepCleanFiles{
			OutputDir: b.config.OutputDir,
		},
		new(prlcommon.StepRemoveDevices),
	}

	// Setup the state bag
	state := new(multistep.BasicStateBag)
	state.Put("cache", cache)
	state.Put("config", &b.config)
	state.Put("driver", driver)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Run
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	ui.Say(fmt.Sprintf("OutputDir = %s", b.config.OutputDir))
	return prlcommon.NewArtifact(b.config.OutputDir)
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
