package guitarstring

import (
	"fmt"
	"github.com/raygarner/lute/scale"
)

const neck_length = 12

type GuitarString struct {
	frets [neck_length]int
	pitch string
}

func (gs GuitarString) Print() {
	fmt.Printf("%s ||", gs.pitch)
	for _, degree := range gs.frets {
		if degree != 0 {
			fmt.Printf(" %2d|", degree)
		} else {
			fmt.Printf("   |")
		}
	}
	fmt.Println()
}

func NewGuitarString(start int, s scale.Scale, pitch string) GuitarString {
	intervals := s.Intervals
	active := s.Active
	var gs GuitarString
	gs.pitch = pitch
	if (start == 0) {
		start = neck_length - 1
	} else {
		start--
	}
	var currentFret = start
	var currentDegree = 0
	var degrees = len(intervals)
	if active[currentDegree % degrees] {
		gs.frets[currentFret % neck_length] = (currentDegree % degrees) + 1
	}
	currentFret = (currentFret + intervals[currentDegree % degrees]) % neck_length
	currentDegree++
	for currentFret != start {
		if active[currentDegree % degrees] {
			gs.frets[currentFret % neck_length] = (currentDegree % degrees) + 1
		}
		currentFret = (currentFret + intervals[currentDegree % degrees]) % neck_length
		currentDegree++
	}
	return gs
}
