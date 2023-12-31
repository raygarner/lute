package guitarstring

import (
	"io"
	"fmt"
	"github.com/raygarner/lute/scale"
)

const NeckLength = 12

type GuitarString struct {
	Frets []int
	Pitch string
}

func (gs GuitarString) PrintFret(f int, padding string, output io.Writer) {
	if gs.Frets[f] != 0 {
		//fmt.Printf(padding + "%2d|", gs.Frets[f])
		fmt.Fprintf(output, padding + "%2d|", gs.Frets[f])
	} else {
		//fmt.Printf(padding + "  |")
		fmt.Fprintf(output, padding + "  |")
	}
}

func (gs GuitarString) Print(lowest int, highest int, output io.Writer) {
	//fmt.Printf("%s ||", gs.Pitch)
	fmt.Fprintf(output, "%s ||", gs.Pitch)
	if lowest < 0 {
		lowest = 0
	}
	if highest > NeckLength {
		highest = NeckLength
	}
	for f := lowest; f < highest; f++ {
		gs.PrintFret(f, " ", output)
	}
	//fmt.Println()
	fmt.Fprintln(output, )
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
