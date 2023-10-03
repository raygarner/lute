package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/json"
	"github.com/raygarner/lute/scale"
	"github.com/raygarner/lute/fretboard"
	"github.com/raygarner/lute/guitarstring"
)

func main() {
	// Mode 1: map intervals
	var strIntervals = flag.String("i", "", "Intervals of the scale to map on the fretboard eg 2,2,1,2,2,2,1")
	var mode = flag.Int("m", 1, "Mode of the interval list passed to apply")
	var tonic = flag.Int("s", 0, "Fret of the tonic note on the lowest string ")
	var active = flag.String("a", "*", "Binary mask to apply to the scale toggling visibility of each note")
	var enumChords = flag.Bool("c", false, "Enumerate all playable 4 note chords from given scale")

	// Mode 2: enumerate scales
	var enumScales = flag.Int("e", 0, "Enumerate all possible 1 octave scales of n notes")

	// General options
	var database = flag.String("d", "./data/modes.json", "Path to JSON file containing aliases for interval permuations")
	var tuning = flag.String("t", "enbngndnanen", "Pitches of strings in descending order. Variable length.")
	var outputFile = flag.String("o", "", "Write output to specified file")
	var vertical = flag.Bool("v", false, "Print diagrams vertically instead of horizontally")

	flag.Parse()
	// if no mode is specified
	if *strIntervals == "" && *enumScales == 0 {
		fmt.Println("Please specify either -i or -e")
		fmt.Println("Use -h to display help")
		return
	}

	// general options
	data, err := ioutil.ReadFile(*database)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &scale.IntervalsAlias)
	if err != nil {
		fmt.Println(err)
	}
	_ = outputFile

	// -i mode
	if *strIntervals != "" {
		if *enumScales != 0 {
			fmt.Println("-i and -e cannot be used together")
			return
		}
		if *tonic == 0 {
			fmt.Println("Please specify the root using -s")
			return
		}
		if *active == "*" {
			*active = "111111111111"
		}
		s, err := scale.NewScale(strIntervals, active, *mode)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s: ", s.Name)
		fmt.Println(s.StrIntervals)
		fmt.Println()
		fb := fretboard.NewFretboard(*tuning, s, *tonic)
		if *vertical {
			fb.Printv(0, guitarstring.NeckLength)
		} else {
			fb.Print(0, guitarstring.NeckLength)
		}
		fmt.Println()
		fmt.Println()
		fmt.Println("Relative modes:")
		relatives := scale.RelativeScales(s.Intervals)
		for m, r := range relatives {
			fmt.Printf("%d: %s %s\n", m+2, r.StrIntervals, r.Name)
		}
		if *enumChords {
			fmt.Println()
			fmt.Println()
			fmt.Println("Enumerating all 4 note chords from given scale:")
			fmt.Println()
			fb.PrintChords(*vertical)
		}
		return
	}

	// -e mode
	if *enumScales != 0 {
		if *strIntervals != "" {
			fmt.Println("-i and -e cannot be used together")
			return
		}
		if *tonic == 0 {
			fmt.Println("Please specify the root using -s")
			return
		}
		if *enumChords {
			fmt.Println("-c cannot be used with -e")
			return
		}
		if *active != "*" {
			fmt.Println("-a cannot be used with -e")
			return
		}
		if *mode != 1 {
			fmt.Println("-m cannot be used with -e")
		}
		var newScale scale.Scale
		scales := scale.EnumIntervals(*enumScales)
		for _, s := range scales {
			fmt.Println()
			fmt.Println(s)
			newScale = scale.NewScaleFromIntervals(s)
			fb := fretboard.NewFretboard(*tuning, newScale, *tonic)
			if *vertical == false {
				fb.Print(0, guitarstring.NeckLength)
			} else {
				fb.Printv(0, guitarstring.NeckLength)
			}
		}
	}

}
