package cpu

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Return a slice with 3 float32 representing 1-min, 5-min, and 15-min average CPU/system loading.
// In fact, in Linux, this is more close to system loading.
func GetAvgLoad() ([]float32, error) {
	content, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	var loads []float32
	for idx, number := range strings.Fields(string(content)) {
		load, err := strconv.ParseFloat(number, 32)
		if err != nil {
			return nil, err
		}
		loads = append(loads, float32(load))

		if idx == 2 {
			break
		}
	}

	return loads, nil
}

// Returns the CPU information in a slice.
// The elements in the slice are:
//  0. Model string.
//  1. Total logical CPU count in integer.
//  2. Physical CPU count in integer.
//  3. Core count per physical CPU.
//  4. Thread count per core.
func GetInfo() ([]any, error) {
	fp, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var physical_id, model string
	var siblings int
	var cpus int    // cpu count
	var sockets int // physical cpu count
	var cores int   // core count per physical cpu
	var threads int // thread count per core
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ":")
		if len(splits) < 2 {
			continue
		}
		key := strings.Trim(splits[0], " \t\n")
		val := strings.Trim(strings.Join(splits[1:], ":"), " \t\n")

		switch key {
		case "processor":
			cpus, _ = strconv.Atoi(val)
		case "model name":
			if len(model) == 0 {
				model = val
			}
		case "physical id":
			if physical_id != val {
				sockets, _ = strconv.Atoi(val)
				physical_id = val
			}
		case "cpu cores":
			if cores == 0 {
				cores, _ = strconv.Atoi(val)
			}
		case "siblings":
			if siblings == 0 {
				siblings, _ = strconv.Atoi(val)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	cpus += 1
	sockets += 1
	threads = siblings / cores

	return []any{model, cpus, sockets, cores, threads}, nil
}

// Returns the CPU usage.
// The first one is aggregated. Starts from the second one, it is the usage of
// each single CPU.
func GetUsage(duration time.Duration) ([]float32, error) {
	cpus_0, err := ReadCpuTimes()
	if err != nil {
		return nil, fmt.Errorf("start: %v", err)
	}

	time.Sleep(duration)

	cpus_1, err := ReadCpuTimes()
	if err != nil {
		return nil, fmt.Errorf("stop: %v", err)
	}

	var cpus [][]float32
	for idx := 0; idx < len(cpus_0); idx++ {
		var diffs []float32
		cpu_times_0 := cpus_0[idx]
		cpu_times_1 := cpus_1[idx]
		for tidx := 0; tidx < len(cpu_times_0); tidx++ {
			diffs = append(diffs, cpu_times_1[tidx]-cpu_times_0[tidx])
		}
		cpus = append(cpus, diffs)
	}

	var usages []float32
	for _, diffs := range cpus {
		var total float32
		for _, diff := range diffs {
			total += diff
		}
		usages = append(usages, 1-diffs[3]/total)
	}

	return usages, nil
}

// Read /proc/stat and returns the values of each line starts with "cpu".
func ReadCpuTimes() ([][]float32, error) {
	fp, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var cpus [][]float32
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu") {
			break
		}

		var cpu_times []float32
		numbers := strings.Fields(line)[1:]
		for _, number := range numbers {
			cpu_time, err := strconv.ParseFloat(number, 32)
			if err != nil {
				return nil, err
			}
			cpu_times = append(cpu_times, float32(cpu_time))
		}
		cpus = append(cpus, cpu_times)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cpus, nil
}
