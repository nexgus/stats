package disk

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func DeviceUniqueID(device string) string {
	switch getBlockDeviceDriver(device) {
	case "sd":
		var model, wwid string
		if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/model", device)); err == nil {
			model = strings.Trim(string(data), " \t\n")
		}
		if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/wwid", device)); err == nil {
			wwid = strings.Trim(string(data), " \t\n")
		}
		return strings.Join([]string{model, wwid}, "")
	default:
		return ""
	}
}

func getBlockDeviceDriver(device string) string {
	fp, err := os.Open(fmt.Sprintf("/sys/block/%s/device/uevent", device))
	if err != nil {
		return ""
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DRIVER=") {
			return strings.Split(line, "=")[1]
		}
	}

	return ""
}
