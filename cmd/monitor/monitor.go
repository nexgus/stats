package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"stats/pkg/cpu"
	"stats/pkg/queue"
)

func main() {
	wg := new(sync.WaitGroup)
	que := queue.New()
	count := 50

	wg.Add(1)
	go func(count int) {
		for counter := 0; counter < count; counter++ {
			if counter%10 == 0 && counter > 0 {
				fmt.Printf("WriteCpuUsage %d\n", counter)
			}
			t0 := time.Now()
			usage, err := cpu.GetUsage(50 * time.Millisecond)
			if err != nil {
				fmt.Printf("GetCpuUsage: %v", err)
				break
			} else {
				ts := float64(time.Now().UnixNano()) / 1000000000
				record := map[string]any{
					"timestamp": ts,
					"type":      "cpu usage",
					"data":      usage,
				}
				que.Push(record)
			}
			t := time.Since(t0)
			time.Sleep(100*time.Millisecond - t)
		}
		wg.Done()
	}(count)

	go func(count int) {
		for counter := 0; counter < count; counter++ {
			fmt.Printf("WriteCpuLoad\n")

			t0 := time.Now()
			load, err := cpu.GetAvgLoad()
			if err != nil {
				fmt.Printf("GetCpuLoad: %v", err)
				break
			} else {
				ts := float64(time.Now().UnixNano()) / 1000000000
				record := map[string]any{
					"timestamp": ts,
					"type":      "cpu load",
					"data":      load,
				}
				que.Push(record)
			}
			t := time.Since(t0)
			time.Sleep(time.Second - t)
		}
		wg.Done()
	}(50 / 10)

	wg.Add(1)
	go func(count int, filepath string) {
		for counter := 0; counter < count; {
			time.Sleep(time.Second)
			for que.NotEmpty() {
				records := que.GetAll().([]any)
				counter += len(records)
				fmt.Printf("Read: %d\n", counter)
				if fp, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664); err != nil {
					fmt.Printf("OpenHistory: %v\n", err)
				} else {
					for _, record := range records {
						line, _ := json.Marshal(record)
						fp.Write(line)
						fp.WriteString("\n")
					}
					fp.Close()
				}
			}
		}
		wg.Done()
	}(count, "./history.txt")

	wg.Wait()
}
