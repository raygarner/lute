package guitarstring

import (
	"fmt"
	"github.com/raygarner/lute/scale"
)

const neck_length = 12

type GuitarString struct {
	frets [neck_length]int
	Pitch string
}

func (gs GuitarString) PrintFret(f int, padding string) {
	if gs.frets[f] != 0 {
		fmt.Printf(padding + "%2d|", gs.frets[f])
	} else {
		fmt.Printf(padding + "  |")
	}
}

func (gs GuitarString) Print() {
	fmt.Printf("%s ||", gs.Pitch)
	for f := 0; f < neck_length; f++ {
		gs.PrintFret(f, " ")
	}
	fmt.Println()
}

func NewGuitarString(start int, s scale.Scale, pitch string) GuitarString {
	intervals := s.Intervals
	active := s.Active
	var gs GuitarString
	gs.Pitch = pitch
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
