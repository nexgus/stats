package main

import (
	"fmt"
	"stats/pkg/machine"
)

func main() {
	mid := machine.GetID()
	fmt.Printf("%x\n", mid)
}
