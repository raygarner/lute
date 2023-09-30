package scale

import (
	"strconv"
	"errors"
	"strings"
	//"fmt"
)

// TODO: implement these
var ScaleAlias = map[string][]int {
	"Ionian": []int{2,2,1,2,2,2,1}, //ionian
	"Dorian": []int{2,2,1,2,2,2,1}, //dorian
	"Phrygian": []int{2,2,1,2,2,2,1}, //phrygian
	"Lydian": []int{2,2,1,2,2,2,1}, //lydian
	"Mixolydian": []int{2,2,1,2,2,2,1}, //mixolydian
	"Aeolian": []int{2,2,1,2,2,2,1}, //aeolian
	"Locrian": []int{2,2,1,2,2,2,1}, //locrian
	"Diminished": []int{2,2,1,2,2,2,1}, //diminished
	"Whole Tone": []int{2,2,1,2,2,2,1}, //whole tone
	"Chromatic": []int{2,2,1,2,2,2,1}, //chromatic
	"Harmonic Minor": []int{2,2,1,2,2,2,1}, //harmonic minor
	"Melodic Minor": []int{2,2,1,2,2,2,1}, //melodic minor
	"Overtone Dominant": []int{2,2,1,2,2,2,1}, //overtone dominant (4 of melodic minor)
}

type Scale struct {
	Intervals []int
	Active []bool
}

const octave = 12
const seperator = ","

func sum(intervals []int) int {
	var total int
	for _, j := range intervals {
		total += j
	}
	return total
}

func appendInterval(length int, intervals []int) [][]int {
	var ret[][]int
	var newintervals = make([]int, len(intervals))

	if length == 0 {
		if sum(intervals) != octave {
			return ret
		} else {
			return append(ret, intervals)
		}
	}
	for i := 1; i <= octave - (sum(intervals) + length) + 1; i++ {
		copy(newintervals, intervals)
		ret = append(ret, appendInterval(length-1, append(newintervals, i))...)
	}
	return ret
}

func EnumIntervals(length int) [][]int {
	return appendInterval(length, []int{})
}

// convert string representation to slice of ints
// TODO: read intervals seperated by commas
func readIntervals(strIntervals *string) ([]int, error) {
	var intIntervals []int
	for _, strInterval := range strings.Split(*strIntervals, seperator) {
		val, _ := strconv.Atoi(string(strInterval))
		intIntervals = append(intIntervals, val)
	}
	return intIntervals, nil
}

// builds a slice of bool from binary string
func readActivity(strActivity *string) ([]bool, error) {
	var activity []bool
	for _, strBit := range *strActivity {
		if string(strBit) == "1" {
			activity = append(activity, true)
		} else if string(strBit) == "0" {
			activity = append(activity, false)
		} else {
			return activity, errors.New("Invalid activity format")
		}
	}
	return activity, nil
}

// rotates a slice one step to the right
func rot(s []int) []int {
	return append(s[1:], s[0])
}

// gets the mode of a slice of intervals by repeatedly rotating it
func applyMode(intervals []int, mode int) []int {
	for i := 0; i < mode; i++ {
		intervals = rot(intervals)
	}
	return intervals
}

// return true if intervals wrap correctly
// also return the number of steps in the scale to form an octave
func validIntervals(intervals []int) (bool, int) {
	var total int
	var i int
	for total < octave {
		total += intervals[i % len(intervals)]
		i++
	}
	return total == octave && i % len(intervals) == 0, i
}

// length is the ammount of steps the scale takes to complete
// ie takes [2], 6 and returns [2,2,2,2,2,2] (whole tone scale)
func completeIntervals(intervals []int, length int) ([]int, error) {
	var i int
	for len(intervals) < length {
		intervals = append(intervals, intervals[i])
		i++
	}
	if sum(intervals) == octave {
		return intervals, nil
	} else {
		return intervals, errors.New("Intervals do not form an octave")
	}
}

func NewScale(strIntervals *string, strActivity *string, rot int) (Scale, error) {
	var s Scale
	var err error
	s.Intervals, _ = readIntervals(strIntervals)
	_, n := validIntervals(s.Intervals)
	s.Intervals, err = completeIntervals(s.Intervals, n)
	s.Intervals = applyMode(s.Intervals, rot - 1)
	s.Active, _ = readActivity(strActivity)
	return s, err
}

// assumes intervals are valid
func NewScaleFromIntervals(intervals []int) Scale {
	var s Scale
	s.Intervals = intervals
	for i := 0; i < octave; i++ {
		s.Active = append(s.Active, true)
	}
	return s
}
