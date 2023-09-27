package main

import (
	"strconv"
	"fmt"
	"flag"
	"github.com/raygarner/lute/scale"
	"github.com/raygarner/lute/fretboard"
)

func main() {
	var strIntervals  = flag.String("i", "2212221", "the intervals of the scale in semitones")
	var mode = flag.Int("m", 1, "mode of the specified scale")
	var active = flag.String("a", "111111111111", "which notes of the scale are active. One bit per degree, extra bits ignored")
	var tonic = flag.Int("s", 8, "fret of the tonic note on the lowest string")
	var tuning = flag.String("t", "enbngndnanen", "the tuning of the instrument in descending order of pitch (works for any number of strings)")
	var vertical = flag.Bool("v", false, "print fretboard vertically")
	flag.Parse()
	fmt.Println("Intervals: " + *strIntervals)
	fmt.Println("Mode: " + strconv.Itoa(*mode))
	fmt.Println("Active: " + *active)
	fmt.Println("Tonic: " + strconv.Itoa(*tonic))
	fmt.Println("Tuning: " + *tuning)
	fmt.Printf("Vertical: %v\n", *vertical)
	fmt.Println()
	s := scale.NewScale(strIntervals, active, *mode)
	fb := fretboard.NewFretboard(*tuning, s, *tonic)
	if *vertical == false {
		fb.Print()
	} else {
		fb.Printv()
	}
}