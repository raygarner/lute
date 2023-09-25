package scale

import (
	"strconv"
	"errors"
)

type Scale struct {
	Intervals []int
	Active []bool
}

const octave = 12

// convert string representation to slice of ints
func readIntervals(strIntervals *string) ([]int, error) {
	var intIntervals []int
	for _, strInterval := range *strIntervals {
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
	return intervals, nil
}

func NewScale(strIntervals *string, strActivity *string, rot int) Scale {
	var s Scale
	s.Intervals, _ = readIntervals(strIntervals)
	_, n := validIntervals(s.Intervals)
	s.Intervals, _ = completeIntervals(s.Intervals, n)
	s.Intervals = applyMode(s.Intervals, rot - 1)
	s.Active, _ = readActivity(strActivity)
	return s
}
