package scale

import (
	"strconv"
	"errors"
	"strings"
	"fmt"
)

var IntervalsAlias map[string]string

type Scale struct {
	Intervals []int
	Active []bool
	StrIntervals string
	Name string
}

const octave = 12
const seperator = ","

// TODO: implement this
func RelativeScales(intervals []int) []Scale {
	var ret []Scale
	for i := 0; i < len(intervals)-1; i++ {
		intervals = rot(intervals)
		ret = append(ret, NewScaleFromIntervals(intervals))
	}
	return ret
}

func sum(intervals []int) int {
	var total int
	for _, j := range intervals {
		total += j
	}
	return total
}

func appendInterval(length int, intervals []int, maxSemitones int) [][]int {
	var ret[][]int
	var newintervals = make([]int, len(intervals))
	var st int

	if length == 0 {
		if sum(intervals) != octave {
			return ret
		} else {
			return append(ret, intervals)
		}
	}
	for i := 1; i <= octave - (sum(intervals) + length) + 1; i++ {
		if i == 1 && len(intervals) >= maxSemitones-1 {
			for j := len(intervals) - 1; j >= len(intervals) - (maxSemitones-1); j-- {
				if intervals[j] == 1 {
					st++
				}
			}
		}
		if st < maxSemitones-1 {
			copy(newintervals, intervals)
			ret = append(ret, appendInterval(length-1, append(newintervals, i), maxSemitones)...)
		} else {
			st = 0
		}
	}
	return ret
}

func EnumIntervals(length int, maxSemitones int) [][]int {
	return appendInterval(length, []int{}, maxSemitones)
}

// convert string representation to slice of ints
func readIntervals(strIntervals *string) ([]int, error) {
	var intIntervals []int
	for _, strInterval := range strings.Split(*strIntervals, seperator) {
		val, _ := strconv.Atoi(string(strInterval))
		intIntervals = append(intIntervals, val)
	}
	return intIntervals, nil
}

// builds a slice of bool from binary string
func readActivity(strActivity *string, intervalLen int) ([]bool, error) {
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
	for len(activity) < intervalLen {
		activity = append(activity, false)
	}
	return activity, nil
}

func intervalsToString(intervals []int) string {
	var ret string
	for _, i := range intervals[:len(intervals)-1] {
		ret += strconv.Itoa(i)
		ret += ","
	}
	ret += strconv.Itoa(intervals[len(intervals)-1])
	return ret
}

// rotates a slice one step to the right
func rot(s []int) []int {
	return append(s[1:], s[0])
}

func CalcStrings(baseIntervals []int) {
	fmt.Println(intervalsToString(baseIntervals))
	for _, _ = range baseIntervals {
		baseIntervals = rot(baseIntervals)
		fmt.Println(intervalsToString(baseIntervals))
	}
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

func IdentifyIntervals(strIntervals string) string {
	ret := IntervalsAlias[strIntervals]
	if ret == "" {
		return "Common name unknown"
	} else {
		return ret
	}
}

func NewScale(strIntervals *string, strActivity *string, rot int) (Scale, error) {
	var s Scale
	var err error
	s.Intervals, _ = readIntervals(strIntervals)
	_, n := validIntervals(s.Intervals)
	s.Intervals, err = completeIntervals(s.Intervals, n)
	s.Intervals = applyMode(s.Intervals, rot - 1)
	s.StrIntervals = intervalsToString(s.Intervals)
	s.Name = IdentifyIntervals(s.StrIntervals)
	s.Active, _ = readActivity(strActivity, len(s.StrIntervals))
	return s, err
}

// assumes intervals are valid
func NewScaleFromIntervals(intervals []int) Scale {
	var s Scale
	s.Intervals = intervals
	s.StrIntervals = intervalsToString(intervals)
	s.Name = IdentifyIntervals(s.StrIntervals)
	for i := 0; i < octave; i++ {
		s.Active = append(s.Active, true)
	}
	return s
}
