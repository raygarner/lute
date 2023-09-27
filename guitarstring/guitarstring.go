package guitarstring

import (
	"fmt"
	"github.com/raygarner/lute/scale"
)

const NeckLength = 12

type GuitarString struct {
	Frets []int
	Pitch string
}

func (gs GuitarString) PrintFret(f int, padding string) {
	if gs.Frets[f] != 0 {
		fmt.Printf(padding + "%2d|", gs.Frets[f])
	} else {
		fmt.Printf(padding + "  |")
	}
}

func (gs GuitarString) Print(lowest int, highest int) {
	fmt.Printf("%s ||", gs.Pitch)
	if lowest < 0 {
		lowest = 0
	}
	if highest > NeckLength {
		highest = NeckLength
	}
	for f := lowest; f < highest; f++ {
		gs.PrintFret(f, " ")
	}
	fmt.Println()
}

func NewGuitarString(start int, s scale.Scale, pitch string) GuitarString {
	intervals := s.Intervals
	active := s.Active
	var gs GuitarString
	gs.Frets = make([]int, NeckLength)
	gs.Pitch = pitch
	if (start == 0) {
		start = NeckLength - 1
	} else {
		start--
	}
	var currentFret = start
	var currentDegree = 0
	var degrees = len(intervals)
	if active[currentDegree % degrees] {
		gs.Frets[currentFret % NeckLength] = (currentDegree % degrees) + 1
	}
	currentFret = (currentFret + intervals[currentDegree % degrees]) % NeckLength
	currentDegree++
	for currentFret != start {
		if active[currentDegree % degrees] {
			gs.Frets[currentFret % NeckLength] = (currentDegree % degrees) + 1
		}
		currentFret = (currentFret + intervals[currentDegree % degrees]) % NeckLength
		currentDegree++
	}
	return gs
}
