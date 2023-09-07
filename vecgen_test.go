package main

import (
	"testing"
	"time"
)

func TestVecEqSet(t *testing.T) {
	const nVecMin = 100
	const nVecMax = 1_000_000
	// The number of goroutines which will be spawned for each vector size being tested
	const nprocPerSize = 4
	const testingDuration = 5 * time.Second

	var stopChannels []chan<- bool
	outChannels := make(map[int][]<-chan VecEqSet)
	for nVec := 100; nVec <= 1_000_000; nVec *= 10 {
		outChannels[nVec] = make([]<-chan VecEqSet, 0)
		for i := 0; i < nprocPerSize; i++ {
			// Create and append I/O channels
			stop := make(chan bool)
			out := make(chan VecEqSet)
			outChannels[nVec] = append(outChannels[nVec], out)
			stopChannels = append(stopChannels, stop)
			go func(nVec int, stop <-chan bool, out chan<- VecEqSet) {
				for {
					set := make(VecEqSet, nVec)
					set.Generate()
					select {}
				}

			}(nVec, stop, out)
		}

	}
}
