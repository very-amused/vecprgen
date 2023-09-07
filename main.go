package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Flags
var nVec uint
var debug bool

// For now, this program only supports generating equillibrium vector addition problems
func main() {
	// Parse flags
	flag.UintVar(&nVec, "n", 2, "Number of vectors to generate (min 2)")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging output")
	flag.Parse()
	if nVec < 2 {
		nVec = 2
		fmt.Println("Invalid vector count, defaulting to 2")
	}
	// Parse params
	params := &Params{}
	params.SetDefault()
	var err error
	// Parse params.ini if present
	const iniFile = "params.ini"
	if _, err = os.Stat(iniFile); err == nil {
		if err = params.Parse(iniFile); err != nil {
			fmt.Printf("Error parsing %s: %s.\n", iniFile, err)
		}
	}

	// Generate vectors as whole integers with a minimum magnitude of 1
	set := VecEqSet{Vecs: make([]Vec, nVec)}
	set.Generate(params)
	// Get component sums
	var iSum, jSum int
	for _, v := range set.Vecs {
		//fmt.Printf("Vector %d = %di + %dj\n", i+1, v.X, v.Y)
		iSum += v.X
		jSum += v.Y
	}
	//fmt.Println()

	if debug {
		log.Println("i component sum:", iSum)
		log.Println("j component sum:", jSum)
		if iSum != 0 || jSum != 0 {
			os.Exit(1)
		}
	}
}
