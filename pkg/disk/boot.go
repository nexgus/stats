package disk

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

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
			splits := strings.Split(fields[0], "/")
			logicalDevice := splits[len(splits)-1]
			idx := len(logicalDevice) - 1
			letter := string(logicalDevice[idx])
			if _, err := strconv.Atoi(letter); err == nil {
				return logicalDevice[:idx], nil
			} else {
				return logicalDevice, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
