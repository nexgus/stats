package machine

import (
	"crypto/sha256"
	"fmt"
	"os"
	"stats/pkg/disk"
	"strings"
)

func GetID() []byte {
	var deviceID []byte
	if device, err := disk.FindBootDevice(); err == nil {
		switch disk.GetBlockDeviceDriver(device) {
		case "sd":
			var model, wwid string
			if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/model", device)); err == nil {
				model = strings.Trim(string(data), " \t\n")
			}
			if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/wwid", device)); err == nil {
				wwid = strings.Trim(string(data), " \t\n")
			}
			deviceID = []byte(strings.Join([]string{model, wwid}, ""))
		case "mmcblk":
			var manfid, name, serial string
			if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/manfid", device)); err == nil {
				manfid = strings.Trim(string(data), " \t\n")
			}
			if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/name", device)); err == nil {
				name = strings.Trim(string(data), " \t\n")
			}
			if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/serial", device)); err == nil {
				serial = strings.Trim(string(data), " \t\n")
			}
			deviceID = []byte(strings.Join([]string{manfid, name, serial}, ""))
		}
	}
	machineID, _ := os.ReadFile("/etc/machine-id")

	machineID = append(machineID, deviceID...)
	hash := sha256.New()
	hash.Write(machineID)

	return hash.Sum(nil)
}
