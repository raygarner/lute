package scale

import (
	"strconv"
	"errors"
	"strings"
	"fmt"
)

var IntervalsAlias = map[string]string {
	"2,2,1,2,2,2,1": "Ionian",
	"2,1,2,2,2,1,2": "Dorian",
	"1,2,2,2,1,2,2": "Phrygian",
	"2,2,2,1,2,2,1": "Lydian",
	"2,2,1,2,2,1,2": "Mixolydian",
	"2,1,2,2,1,2,2": "Aeolian",
	"1,2,2,1,2,2,2": "Locrian",

	"2,1,2,2,1,3,1": "Harmonic Minor",
	"1,2,2,1,3,1,2": "Locrian n6",
	"2,2,1,3,1,2,1": "Augmented Major",
	"2,1,3,1,2,1,2": "Ukrainian Dorian",
	"1,3,1,2,1,2,2": "Phrygian Dominant",
	"3,1,2,1,2,2,1": "Lydian #2",
	"1,2,1,2,2,1,3": "Altered Diminished",

	"2,1,2,2,2,2,1": "Ascending Melodic Minor",
	"1,2,2,2,2,1,2": "Dorian b2",
	"2,2,2,2,1,2,1": "Lydian Augmented",
	"2,2,2,1,2,1,2": "Overtone Dominant",
	"2,2,1,2,1,2,2": "Aeolian Dominant",
	"2,1,2,1,2,2,2": "Half Diminished",
	"1,2,1,2,2,2,2": "Altered",

	"2,2,1,2,1,3,1": "Harmonic Major",
	"2,1,2,1,3,1,2": "Dorian b5",
	"1,2,1,3,1,2,2": "Phrygian b4",
	"2,1,3,1,2,2,1": "Lydian b3",
	"1,3,1,2,2,1,2": "Mixolydian b2",
	"3,1,2,2,1,2,1": "Lydian Augmented #2",
	"1,2,2,1,2,1,3": "Locrian bb7",

	"2,1,2,1,2,1,2,1": "Diminished",
	"1,2,1,2,1,2,1,2": "Dominant Diminished",
}

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

func IntervalsToString(intervals []int) string {
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
	fmt.Println(IntervalsToString(baseIntervals))
	for _, _ = range baseIntervals {
		baseIntervals = rot(baseIntervals)
		fmt.Println(IntervalsToString(baseIntervals))
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
	s.StrIntervals = IntervalsToString(s.Intervals)
	s.Name = IdentifyIntervals(s.StrIntervals)
	s.Active, _ = readActivity(strActivity)
	return s, err
}

// assumes intervals are valid
func NewScaleFromIntervals(intervals []int) Scale {
	var s Scale
	s.Intervals = intervals
	s.StrIntervals = IntervalsToString(intervals)
	s.Name = IdentifyIntervals(s.StrIntervals)
	for i := 0; i < octave; i++ {
		s.Active = append(s.Active, true)
	}
	return s
}
