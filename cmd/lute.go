package main

import (
	"strconv"
	"fmt"
	"flag"
	"github.com/raygarner/lute/scale"
	"github.com/raygarner/lute/fretboard"
	"github.com/raygarner/lute/guitarstring"
)

func main() {
	var strIntervals  = flag.String("i", "2,2,1,2,2,2,1", "the intervals of the scale in semitones")
	var mode = flag.Int("m", 1, "mode of the specified scale")
	var active = flag.String("a", "111111111111", "which notes of the scale are active. One bit per degree, extra bits ignored")
	var tonic = flag.Int("s", 8, "fret of the tonic note on the lowest string")
	var tuning = flag.String("t", "enbngndnanen", "the tuning of the instrument in descending order of pitch (works for any number of strings)")
	var vertical = flag.Bool("v", false, "print fretboard vertically")
	var chords = flag.Bool("c", false, "enumerate all playable 4 note chords")
	var enum = flag.Int("e", 0, "enumerate all possible 1 octave scales of given length")
	flag.Parse()
	fmt.Println("Intervals: " + *strIntervals)
	fmt.Println("Mode: " + strconv.Itoa(*mode))
	fmt.Println("Active: " + *active)
	fmt.Println("Tonic: " + strconv.Itoa(*tonic))
	fmt.Println("Tuning: " + *tuning)
	fmt.Printf("Vertical: %v\n", *vertical)
	fmt.Printf("Chords: %v\n", *chords)
	fmt.Printf("Enumerate: %d\n", *enum)
	fmt.Println()
	fmt.Println()
	s, err := scale.NewScale(strIntervals, active, *mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s: ", s.Name)
	fmt.Println(s.StrIntervals)
	fmt.Println()
	fb := fretboard.NewFretboard(*tuning, s, *tonic)
	if *vertical == false {
		fb.Print(0, guitarstring.NeckLength)
	} else {
		fb.Printv(0, guitarstring.NeckLength)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("Relative modes:")
	relatives := scale.RelativeScales(s.Intervals)
	for m, r := range relatives {
		fmt.Printf("%d: %s %s\n", m+2, r.StrIntervals, r.Name)
	}
	if *chords {
		fmt.Println()
		fmt.Println()
		fmt.Println("Enumerating all 4 note chords from given scale:")
		fmt.Println()
		fb.PrintChords(*vertical)
	}
	var newScale scale.Scale
	if *enum > 0 {
		fmt.Printf("Enumerating all %d note scales:\n", *enum)
		scales := scale.EnumIntervals(*enum)
		for _, s := range scales {
			fmt.Println()
			fmt.Println(s)
			newScale = scale.NewScaleFromIntervals(s)
			fb = fretboard.NewFretboard(*tuning, newScale, *tonic)
			if *vertical == false {
				fb.Print(0, guitarstring.NeckLength)
			} else {
				fb.Printv(0, guitarstring.NeckLength)
			}
		}
	}
}