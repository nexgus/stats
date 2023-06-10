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
