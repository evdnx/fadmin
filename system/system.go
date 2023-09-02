package system

import (
	"runtime"

	"github.com/evdnx/unixmint/cmd"
)

type Host struct {
	Hostname                  string `json:"Hostname"`
	StaticHostname            string `json:"StaticHostname"`
	PrettyHostname            string `json:"PrettyHostname"`
	DefaultHostname           string `json:"DefaultHostname"`
	HostnameSource            string `json:"HostnameSource"`
	IconName                  string `json:"IconName"`
	Chassis                   string `json:"Chassis"`
	Deployment                string `json:"Deployment"`
	Location                  string `json:"Location"`
	KernelName                string `json:"KernelName"`
	KernelRelease             string `json:"KernelRelease"`
	KernelVersion             string `json:"KernelVersion"`
	OperatingSystemPrettyName string `json:"OperatingSystemPrettyName"`
	OperatingSystemCPEName    string `json:"OperatingSystemCPEName"`
	OperatingSystemHomeURL    string `json:"OperatingSystemHomeURL"`
	HardwareVendor            string `json:"HardwareVendor"`
	HardwareModel             string `json:"HardwareModel"`
	HardwareSerial            string `json:"HardwareSerial"`
	FirmwareVersion           string `json:"FirmwareVersion"`
	FirmwareVendor            string `json:"FirmwareVendor"`
	FirmwareDate              uint   `json:"FirmwareDate"`
	ProductUUID               string `json:"ProductUUID"`
}

func Info() (string, error) {
	switch runtime.GOOS {
	case "linux":
	case "freebsd":
	case "openbsd":
	case "netbsd":
	case "dragonfly":
	}

	out, err := cmd.Exec("hostnamectl --json=short").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
