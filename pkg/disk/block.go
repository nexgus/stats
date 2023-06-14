package disk

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func DeviceUniqueID(device string) string {
	switch GetBlockDeviceDriver(device) {
	case "sd":
		var model, wwid string
		if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/model", device)); err == nil {
			model = strings.Trim(string(data), " \t\n")
		}
		if data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/wwid", device)); err == nil {
			wwid = strings.Trim(string(data), " \t\n")
		}
		return strings.Join([]string{model, wwid}, "")
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
		return strings.Join([]string{manfid, name, serial}, "")
	default:
		return ""
	}
}

func FindBootDevice() (string, error) {
	fp, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		// Examples:
		// tmpfs /run tmpfs rw,nosuid,nodev,noexec,relatime,size=3275448k,mode=755,inode64 0 0
		// /dev/sda2 / ext4 rw,relatime,errors=remount-ro 0 0
		// securityfs /sys/kernel/security securityfs rw,nosuid,nodev,noexec,relatime 0 0

		line := scanner.Text()
		fields := strings.Fields(line)
		if fields[1] == "/" {
			return NameLogicalToPhysical(fields[0]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func GetBlockDeviceDriver(device string) string {
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

func NameLogicalToPhysical(logical string) string {
	var physical string
	var idx int
	if strings.HasPrefix(logical, "sd") {
		for idx = len(logical) - 1; idx > 1; idx-- {
			char := string(logical[idx])
			if _, err := strconv.Atoi(char); err != nil {
				break
			}
		}
		physical = logical[:idx]
	} else if strings.HasPrefix(logical, "mmcblk") {
		idx = strings.Index(logical, "p")
		physical = logical[:idx]
	} else {
		physical = logical
	}

	return physical
}
