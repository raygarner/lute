package main

import (
	"fmt"
	"flag"
	"strconv"
	"errors"
)

const frets = 12

var offsets_full = map[string]int {
	"en": 0,
	"f-": 0,
	"e+": 11,
	"fn": 11,
	"f+": 10,
	"g-": 10,
	"gn": 9,
	"g+": 8,
	"a-": 8,
	"an": 7,
	"a+": 6,
	"b-": 6,
	"bn": 5,
	"c-": 5,
	"cn": 4,
	"b+": 4,
	"c+": 3,
	"d-": 3,
	"dn": 2,
	"d+": 1,
	"e-": 1,
}

// prints fret numbers
func printFrets() {
	fmt.Printf("     ")
	for i := 1; i <= 12; i++ {
		fmt.Printf(" %2d ", i)
	}
	fmt.Println()
}

// prints one guitar string
func printGuitarString(guitarString [frets]int, pitch string) {
	fmt.Printf("%s ||", pitch)
	for _, degree := range guitarString {
		if degree != 0 {
			fmt.Printf(" %2d|", degree)
		} else {
			fmt.Printf("   |")
		}
	}
	fmt.Println()
}

// returns a guitar string with the scale data on it
func buildString(start int, intervals []int, active []bool) [frets]int {
	if (start == 0) {
		start = 11
	} else {
		start--
	}
	var currentFret = start
	var currentDegree = 0
	var guitarString [frets]int
	var degrees = len(intervals)
	if active[currentDegree % degrees] {
		guitarString[currentFret % frets] = (currentDegree % degrees) + 1
	}
	currentFret = (currentFret + intervals[currentDegree % degrees]) % frets
	currentDegree++
	for currentFret != start {
		if active[currentDegree % degrees] {
			guitarString[currentFret % frets] = (currentDegree % degrees) + 1
		}
		currentFret = (currentFret + intervals[currentDegree % degrees]) % frets
		currentDegree++
	}
	return guitarString
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

// return true if intervals wrap correctly
// also return the number of steps in the scale to form an octave
func validIntervals(intervals []int) (bool, int) {
	var total int
	var i int
	for total < 12 {
		total += intervals[i % len(intervals)]
		i++
	}
	return total == 12 && i % len(intervals) == 0, i
}

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

func buildOffsets(tuning string) ([]int, []string) {
	var offsets []int
	var names []string
	var lowest = tuning[len(tuning)-2:]
	var offset int

	fmt.Printf("%v\n", lowest)
	for i := 0; i < len(tuning); i += 2 {
		offset = offsets_full[tuning[i:i+2]] - offsets_full[lowest]
		if offset < 0 {
			offset += frets
		}
		offsets = append(offsets, offset)
		names = append(names, tuning[i:i+2])
	}
	return offsets, names
}

func main() {
	var strIntervals  = flag.String("i", "2212221", "the intervals of the scale in semitones")
	var mode = flag.Int("m", 1, "mode of the specified scale")
	var active = flag.String("a", "1111111", "which notes of the scale are active")
	var tonic = flag.Int("s", 8, "fret of the tonic note on the lowest string")
	var tuning = flag.String("t", "enbngndnanen", "the tuning of the instrument in descending order of pitch (works for any number of strings)")
	var intervals []int
	flag.Parse()
	intervals, _ = readIntervals(strIntervals)
	_, n := validIntervals(intervals)
	intervals, _ = completeIntervals(intervals, n)
	intervals = applyMode(intervals, *mode - 1)
	activity, _ := readActivity(active)
	offsets, string_names := buildOffsets(*tuning)
	for i, offset := range offsets {
		guitarString := buildString((*tonic+offset) % frets, intervals, activity)
		printGuitarString(guitarString, string_names[i])
	}

	fmt.Println()
	printFrets()
}