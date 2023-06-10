package memory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Returns memory usage in byte.
// The first element is physical memory usage, including
//
//	[0] total: 	the total amount of physical RAM installed in the machine
//	[1] used:	this is calculated by total - free - buffers - cached
//	[2] free:	the amount of unused memory
//	[3] shared:	memory that is used by the tmpfs file system.
//	[4] buffers:	memory used for buffers
//	[5] cache:		memory used for cache
//	[6] available:	this is an estimation of the memory that is available to service memory requests from applications.
//
// The second element is virtual memory (SWAP) usage, incluing
//
//	[0] total
//	[1] used
//	[2] free
func GetUsage() ([][]int, error) {
	fp, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var p, v []int
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		key := fields[0]
		val := fields[1]

		switch key {
		case "MemTotal:": // p[0]
			fallthrough
		case "MemFree:": // p[1]
			fallthrough
		case "MemAvailable:": // p[2]
			fallthrough
		case "Buffers:": // p[3]
			fallthrough
		case "Cached:": // p[4]
			fallthrough
		case "Shmem:": // p[5]
			kb, _ := strconv.Atoi(val)
			p = append(p, kb*1024)
		case "SwapTotal:": // v[0]
			fallthrough
		case "SwapFree:": // v[1]
			kb, _ := strconv.Atoi(val)
			v = append(v, kb*1024)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return [][]int{
		{p[0], p[0] - p[1] - p[3] - p[4], p[1], p[5], p[3], p[4], p[2]},
		{v[0], v[0] - v[1], v[1]},
	}, nil
}
