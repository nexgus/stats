package main

import (
	"fmt"
	"log"
	"time"

	"stats/pkg/cpu"
)

func main() {
	if info, err := cpu.GetInfo(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("CPU(s): %d\n", info[1])
		fmt.Printf("Model:  %s (%dS/%dC/%dT)\n",
			info[0], info[2], info[3], info[4])
	}

	if usages, err := cpu.GetUsage(100 * time.Millisecond); err != nil {
		log.Fatal(err)
	} else {
		for idx, usage := range usages {
			var key string
			if idx == 0 {
				key = "cpu"
			} else {
				key = fmt.Sprintf("cpu%d", idx-1)
			}
			fmt.Printf("%s: %.2f%%\n", key, usage*100)
		}
	}

	if loads, err := cpu.GetAvgLoad(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("1-min avg load:  %.2f\n", loads[0])
		fmt.Printf("5-min avg load:  %.2f\n", loads[1])
		fmt.Printf("15-min avg load: %.2f\n", loads[2])
	}
}
