package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// Generate vectors in the inclusive range (0,1)|(1,0) to (100,100)
const minComponent = 0
const maxComponent = 10

// Max allowed delta abs value before performing a 2-iter correction to 0---must be in range [2*minComponent, 2*maxComponent]
// NOTE: should be set below 2*maxComponent to produce better data
const maxAbsD = int(1.5 * maxComponent)

// A vector
type Vec struct {
	X int
	Y int

	// Vector angle measured CCW from +x axis in degrees
	// Angle is undefined if nil
	Angle *float64
}

// Calculate v.Angle from v.X and v.Y
func (v *Vec) calculateAngle() {
	if v.Y == 0 {
		return
	}
	v.Angle = new(float64)
	*v.Angle = math.Atan(float64(abs(v.X)/abs(v.Y))) * (180 / math.Pi)
	// Quadrant-correct theta
	if v.X < 0 { // Q2
		*v.Angle = 180 - *v.Angle
		if v.Y < 0 { // Q3
			*v.Angle = 270 - *v.Angle
		}
	} else if v.Y < 0 { // Q4
		*v.Angle = 360 - *v.Angle
	}
}

// VecEqSet - A set of vectors in equillibrium
type VecEqSet struct {
	Vecs []Vec

	params *Params

	// Generation state
	xState vecGenState
	yState vecGenState
}

type vecGenState struct {
	delta      int  // Deltas used to ensure equillibrium
	sign       int  // +/- sign, must be initialized to +-1 before genComponent is called
	correction uint // Count of remaining correction iterations for which signX/Y will stay constant respectively
}

// Generate a set of vectors in equillibrium
func (set *VecEqSet) Generate(params *Params) {
	// Initialize generation state
	set.xState = vecGenState{
		sign: 1}
	set.yState = set.xState // Copy

	for i := range set.Vecs {
		v := &set.Vecs[i]
		if i%2 == 0 {
			// Gen x component first (x == 0 is possible)
			set.genComponent(i, &v.X, &set.xState)
			if v.X == 0 {
				// Ensure empty vectors aren't generated
				for v.Y == 0 {
					set.genComponent(i, &v.Y, &set.yState)
				}
			} else {
				set.genComponent(i, &v.Y, &set.yState)
			}
		} else {
			// Gen y component first (y == 0 is possible)
			set.genComponent(i, &v.Y, &set.yState)
			if v.Y == 0 {
				// Ensure empty vectors aren't generated
				for v.X == 0 {
					set.genComponent(i, &v.X, &set.xState)
				}
			} else {
				set.genComponent(i, &v.X, &set.xState)
			}
		}
		if debug {
			if set.xState.correction > 0 || set.yState.correction > 0 { // Highlight vectors where component correction is in process
				fmt.Print("\x1b[34m")
			}
			log.Printf("Done with vector %d/%d:\n\tdX = %d, cX = %d\n\tdY = %d, cY = %d\n",
				i+1, len(set.Vecs), set.xState.delta, set.xState.correction, set.yState.delta, set.yState.correction)
			fmt.Print("\x1b[0m")
		}
		v.calculateAngle()
	}
}

// Generate a vector component with respect to delta, sign, and correction state
func (set *VecEqSet) genComponent(i int, component *int, state *vecGenState) {
	delta := &state.delta
	sign := &state.sign
	correction := &state.correction
	if i == len(set.Vecs)-1 { // Special single iter correction for the last vector of the set
		if *correction > 0 {
			*correction = 0
		}
		*component = -*delta
		*delta += *component
		return
	}
	// Generate component value
	min := minComponent
	max := maxComponent
	ad := abs(*delta)
	if *correction > 0 {
		if ad <= max {
			// Make it impossible to fully correct delta before the last correction iteration
			max = ad - (int(*correction) - 1)
		}
		if d := ad - max; d > 0 && (*correction) == 2 {
			// Bring abs(*delta) down to at least maxComponent to ensure the last iteration
			// corrects delta such that abs(delta) is contained within [min, maxComponent]
			min = d
		}
	} else if ad >= max {
		// Prevent ad from exceeding 2*max
		max = (2*max - ad)
	}

	if max == min {
		*component = (*sign) * min
	} else {
		for {
			*component = (*sign) * (min + rand.Intn(max-min))
			// Ensure the penultimate vector doesn't set delta to 0 (which can cause deadlocks)
			// or cause delta to exceed the maximum magnitude
			if i != len(set.Vecs)-2 || (*component != -*delta && abs(*delta+*component) <= max) {
				break
			}
		}
	}
	*delta += *component
	ad = abs(*delta)

	// Handle delta management and sign flipping behavior
	if *correction == 0 && (ad > maxAbsD || (ad > maxComponent && i == len(set.Vecs)-3)) {
		*correction += 2
		if *delta < 0 {
			*sign = 1
		} else {
			*sign = -1
		}
	} else if *correction > 0 {
		*correction--
	}
	// If delta correction is not currently needed, the component sign is flipped each iter to assist with data uniformity
	if *correction == 0 {
		*sign = -*sign
	}
}
