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

func imax(i, j int) int {
	if i > j {
		return i
	}
	return j
}
func imin(i, j int) int {
	if i < j {
		return i
	}
	return j
}

// Generate vectors in the inclusive range (0,1)|(1,0) to (100,100)
const minComponent = 0
const maxComponent = 50

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
type VecEqSet []Vec

// Generate a set of vectors in equillibrium
func (set VecEqSet) Generate() {
	var (
		// Deltas used to ensure equillibrium
		dX, dY int
		// +/- signs, must be set to -1 or 1
		signX int = 1
		signY int = 1
		// Count of remaining correction iterations for which signX and/or signY will not change respectively
		cX, cY uint
	)

	for i := range set {
		v := &set[i]
		if i%2 == 0 {
			// Gen x component first (x == 0 is possible)
			set.genComponent(i, &v.X, &dX, &signX, &cX)
			if v.X == 0 {
				// Ensure empty vectors aren't generated
				for v.Y == 0 {
					set.genComponent(i, &v.Y, &dY, &signY, &cY)
				}
			} else {
				set.genComponent(i, &v.Y, &dY, &signY, &cY)
			}
		} else {
			// Gen y component first (y == 0 is possible)
			set.genComponent(i, &v.Y, &dY, &signY, &cY)
			if v.Y == 0 {
				// Ensure empty vectors aren't generated
				for v.X == 0 {
					set.genComponent(i, &v.X, &dX, &signX, &cX)
				}
			} else {
				set.genComponent(i, &v.X, &dX, &signX, &cX)
			}
		}
		if debug {
			if cX+cY > 0 { // Highlight vectors where component correction is in process
				fmt.Print("\x1b[34m")
			}
			log.Printf("Done with vector %d/%d:\n\tdX = %d, cX = %d\n\tdY = %d, cY = %d\n", i+1, len(set), dX, cX, dY, cY)
			fmt.Print("\x1b[0m")
		}
		v.calculateAngle()
	}
}

// Generate a vector component with respect to delta, sign, and correction state
func (set VecEqSet) genComponent(i int, component *int, delta *int, sign *int, correction *uint) {
	if i == len(set)-1 { // Special single iter correction for the last vector of the set
		*component = -*delta
		*delta += *component
		return
	}
	// Generate component value
	min := minComponent
	max := maxComponent
	ad := abs(*delta)
	if *correction > 0 {
		switch *correction {
		case 3:
			fallthrough
		case 2:
			if d := ad - max; d > 0 && (*correction) == 2 {
				// Bring abs(*delta) down to at least maxComponent to ensure the final correction is legitimate
				min = d
			}
			if ad <= max {
				// Make it impossible to fully correct delta before the last correction iteration
				max = ad - (int(*correction) - 1)
			}
		case 1:
			*component = -*delta
			*delta += *component
			*correction--
			*sign = -*sign
			return
		}
	} else {
		// Prevent ad from exceeding 2*max
		if ad >= max {
			max = (2*max - ad)
		}
		// Must end with abs(*delta) <= 100
		/*if i == len(set)-2 {
			if ad > 0 {
				// Leave a delta for the final iteration
				max = (abs(*delta) - (min + 1))
			} else {
				*component = 0
				return
			}
			// Generating a penultimate component with a strength of 0 can cause a deadlock
			if min == 0 {
				min = 1
			}
		} else {
			// abs(*delta) can never exceed 2*max without introducing several tedious and error-prone modulos above
		}*/
	}

	if max == min {
		*component = (*sign) * min
	} else {
		if debug {
			//log.Printf("About to call Intn with value (max-min) = (%d-%d)\n", max, min)
		}
		*component = (*sign) * (min + rand.Intn(max-min))
	}
	*delta += *component

	// Handle delta management and sign flipping behavior
	if *correction == 0 && (abs(*delta) > maxAbsD || (abs(*delta) > maxComponent && i == len(set)-3)) {
		// Perform either 2 or 3 correction iters so that the last iter satisfies (len(set)-1) - i % 2 == 0
		*correction += 2 + uint(len(set)-(1+i))%2
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
