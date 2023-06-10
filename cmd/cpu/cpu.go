package main

import (
	"fmt"
	"log"
	"stats/pkg/cpu"
	"time"
)

func main() {
	if usages, err := cpu.Usage(100 * time.Millisecond); err != nil {
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
}
