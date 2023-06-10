package main

import (
	"fmt"
	"log"

	"stats/pkg/memory"
)

func main() {
	if usage, err := memory.GetUsage(); err != nil {
		log.Fatalf("GetUsage: %v", err)
	} else {
		fmt.Printf("Physical Memory\n")
		fmt.Printf("    Total:     %d\n", usage[0][0])
		fmt.Printf("    Used:      %d\n", usage[0][1])
		fmt.Printf("    Free:      %d\n", usage[0][2])
		fmt.Printf("    Shared:    %d\n", usage[0][3])
		fmt.Printf("    Buffers:   %d\n", usage[0][4])
		fmt.Printf("    Cache:     %d\n", usage[0][5])
		fmt.Printf("    Available: %d\n", usage[0][6])
		fmt.Printf("Virtual Memory\n")
		fmt.Printf("    Total:     %d\n", usage[1][0])
		fmt.Printf("    Used:      %d\n", usage[1][1])
		fmt.Printf("    Free:      %d\n", usage[1][2])
	}
}
