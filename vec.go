package main

import (
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
const minComponent = 1
const maxComponent = 100

// Max allowed delta abs value before performing a 2-iter correction to 0---must be in range [2*minComponent, 2*maxComponent]
// NOTE: should be set below 2*maxComponent to produce better data
const maxAbsD = int(1.5 * maxComponent)

// A vector
type Vec struct {
	X int
	Y int

	// Vector angle measured CCW from +x axis in degrees
	Angle float64
}

// Calculate v.Angle from v.X and v.Y
func (v *Vec) calculateAngle() {
	v.Angle = math.Atan(float64(abs(v.X)/abs(v.Y))) * (180 / math.Pi)
	// Quadrant-correct theta
	if v.X < 0 { // Q2
		v.Angle = 180 - v.Angle
		if v.Y < 0 { // Q3
			v.Angle = 270 - v.Angle
		}
	} else if v.Y < 0 { // Q4
		v.Angle = 360 - v.Angle
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
		set.genComponent(i, &v.X, &dX, &signX, &cX)
		set.genComponent(i, &v.Y, &dY, &signY, &cY)
		v.calculateAngle()
	}
}

// Generate a vector component with respect to delta, sign, and correction state
func (set VecEqSet) genComponent(i int, component *int, delta *int, sign *int, correction *uint) {
	// Generate component value
	min := minComponent
	max := maxComponent
	if *correction > 0 {
		max -= int(*correction) - 1
		// Set min to ensure delta can be corrected to 0 in the remaining iters
		switch *correction {
		case 2:
			if d := abs(*delta) - maxComponent; d > 0 {
				min = d
			}
		case 1:
			min = abs(*delta)
		}
	}
	*component = (*sign) * (min + rand.Intn(max-min))
	*delta += *component

	// Handle delta management and sign flipping behavior
	if *correction == 0 && (*delta > maxAbsD || (*delta > maxComponent && i == len(set)-3)) {
		// Perform either 2 or 3 correction iters so that the last iter satisfies (len(set)-1) - i % 2 == 0
		*correction += 2 + uint(len(set)-(1+i))%2
		if *delta < 0 {
			*sign = -1
		} else {
			*sign = 1
		}
	} else if *correction > 0 {
		*correction--
	}
	// If delta correction is not currently needed, the component sign is flipped each iter to assist with data uniformity
	if *correction == 0 {
		*sign *= -1
	}
}
