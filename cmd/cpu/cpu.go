package main

import (
	"fmt"
	"log"
	"time"

	"stats/pkg/cpu"
)

func main() {
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
