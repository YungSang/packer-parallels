package common

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// Interface to help find the IP address of a running virtual machine.
type GuestIPFinder interface {
	GuestIP() (string, error)
}

// DHCPLeaseGuestLookup looks up the IP address of a guest using DHCP
// lease information from the VMware network devices.
type DHCPLeaseGuestLookup struct {
	// Driver that is being used (to find leases path)
	Driver Driver

	// MAC address of the guest.
	MACAddress string
}

func (f *DHCPLeaseGuestLookup) GuestIP() (string, error) {
	dhcpLeasesPath := f.Driver.DhcpLeasesPath()
	log.Printf("DHCP leases path: %s", dhcpLeasesPath)
	if dhcpLeasesPath == "" {
		return "", errors.New("no DHCP leases path found.")
	}

	var stdout bytes.Buffer

	cmd := exec.Command("grep", f.MACAddress, dhcpLeasesPath)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var ipAddress string

	ipLineRe := regexp.MustCompile(`^(.+?)=(.+?)$`)
	line := strings.TrimSpace(stdout.String())
	matches := ipLineRe.FindStringSubmatch(line)
	if matches != nil {
		ipAddress = matches[1]
	}

	if ipAddress == "" {
		return "", errors.New("IP not found for MAC in DHCP leases")
	}

	return ipAddress, nil
}
