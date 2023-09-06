package main

import (
	"flag"
	"fmt"
)

// For now, this program only supports generating equillibrium vector addition problems
func main() {
	// Parse flags
	var nVec uint
	flag.UintVar(&nVec, "n", 2, "Number of vectors to generate (min 2)")
	flag.Parse()
	if nVec < 2 {
		nVec = 2
		fmt.Println("Invalid vector count, defaulting to 2")
	}

	// Generate vectors as whole integers with a minimum magnitude of 1
	set := make(VecEqSet, nVec)
	set.Generate()
	// Get component sums
	var iSum, jSum int
	for _, v := range set {
		iSum += v.X
		jSum += v.Y
	}
	fmt.Println("i component sum:", iSum)
	fmt.Println("j component sum:", jSum)
}
