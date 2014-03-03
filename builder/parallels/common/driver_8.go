package common

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type Parallels8Driver struct {
	// This is the path to the "prlctl" application.
	PrlCtlPath string
}

func (d *Parallels8Driver) Import(srcName, dstName string) error {
	args := []string{
		"clone", srcName,
		"--name", dstName,
	}

	return d.PrlCtl(args...)
}

func (d *Parallels8Driver) Delete(name string) error {
	return d.PrlCtl("delete", name)
}

func (d *Parallels8Driver) IsRunning(name string) (bool, error) {
	var stdout bytes.Buffer

	cmd := exec.Command(d.PrlCtlPath, "list", name, "--no-header", "--output", "status")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return false, err
	}

	status := strings.TrimSpace(stdout.String())

	if status == "running" {
		return true, nil
	}

	if status == "suspended" {
		return true, nil
	}

	if status == "paused" {
		return true, nil
	}

	return false, nil
}

func (d *Parallels8Driver) Stop(name string) error {
	if err := d.PrlCtl("stop", name); err != nil {
		return err
	}

	// We sleep here for a little bit to let the session "unlock"
	time.Sleep(2 * time.Second)

	return nil
}

func (d *Parallels8Driver) UseDefaults(name string) error {
	if err := d.PrlCtl("set", name, "--usedefanswers", "on"); err != nil {
		return err
	}

	return nil
}

func (d *Parallels8Driver) PrlCtl(args ...string) error {
	var stdout, stderr bytes.Buffer

	log.Printf("Executing prlctl: %#v", args)
	cmd := exec.Command(d.PrlCtlPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("prlctl error: %s", stderrString)
	}

	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)

	return err
}

func (d *Parallels8Driver) Verify() error {
	return nil
}

func (d *Parallels8Driver) Version() (string, error) {
	var stdout bytes.Buffer

	cmd := exec.Command(d.PrlCtlPath, "--version")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	versionOutput := strings.TrimSpace(stdout.String())
	log.Printf("prlctl --version output: %s", versionOutput)

	versionRe := regexp.MustCompile("[^.0-9]")
	matches := versionRe.Split(versionOutput, 2)
	if len(matches) == 0 || matches[0] == "" {
		return "", fmt.Errorf("No version found: %s", versionOutput)
	}

	log.Printf("prlctl version: %s", matches[0])
	return matches[0], nil
}

func (d *Parallels8Driver) SSHAddressFunc(config *SSHConfig) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		driver := state.Get("driver").(Driver)
		vmName := state.Get("vmName").(string)

		ui := state.Get("ui").(packer.Ui)
		ui.Say(fmt.Sprintf("Lookup up IP information... for [%s]", vmName))

		log.Println("Lookup up IP information...")

		var stdout bytes.Buffer

		cmd := exec.Command(d.PrlCtlPath, "list", "-i", vmName)
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			return "", err
		}

		var macAddress string

		macLineRe := regexp.MustCompile(`net0 (.+?) mac=(.+?) (.+?)$`)

		for _, line := range strings.Split(stdout.String(), "\n") {
			// Need to trim off CR character when running in windows
			line = strings.TrimSpace(line)

			matches := macLineRe.FindStringSubmatch(line)
			if matches != nil {
				macAddress = strings.ToLower(matches[2])
				break
			}
		}

		ui.Say(fmt.Sprintf("Mac Address = [%s]", macAddress))

		if macAddress == "" {
			return "", errors.New("MAC address not found")
		}

		ipLookup := &DHCPLeaseGuestLookup{
			Driver:     driver,
			MACAddress: macAddress,
		}

		ipAddress, err := ipLookup.GuestIP()
		if err != nil {
			log.Printf("IP lookup failed: %s", err)
			return "", fmt.Errorf("IP lookup failed: %s", err)
		}

		if ipAddress == "" {
			log.Println("IP is blank, no IP yet.")
			return "", errors.New("IP is blank")
		}

		log.Printf("Detected IP: %s", ipAddress)
		return fmt.Sprintf("%s:%d", ipAddress, config.SSHPort), nil
	}
}

func (d *Parallels8Driver) DhcpLeasesPath() string {
	return "/Library/Preferences/Parallels/parallels_dhcp_leases"
}
