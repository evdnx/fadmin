package system

import (
	"github.com/evdnx/unixmint/bash"
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
	out, err := bash.Cmd("hostnamectl --json=short").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
