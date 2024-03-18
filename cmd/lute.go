package main

import (
	"fmt"
	"flag"
	"io"
	"encoding/json"
	"os"
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
	var buildChord = flag.String("b", "", "Steps to apply to build a chord (eg 3,3 for root inversion triad")
	var chordRoot = flag.Int("r", 0, "Degree of root note when using -b")
	var inversion = flag.Int("inv", 0, "Inversion of the chord inputting with -b")

	// Mode 2: enumerate scales
	var enumScales = flag.Int("e", 0, "Enumerate all possible 1 octave scales of n notes")
	var maxSemitones = flag.Int("f", 12, "Max number of consecutive semitones allowed in generated scales")

	// General options
	//var database = flag.String("d", "~/go/src/lute/data/modes.json", "Path to JSON file containing aliases for interval permuations")
	var database = flag.String("d", os.Getenv("HOME") + "/go/src/lute/data/modes.json", "Path to JSON file containing aliases for interval permuations")

	var tuning = flag.String("t", "enbngndnanen", "Pitches of strings in descending pitch order. Variable length.")
	var outputFile = flag.String("o", "", "Write output to specified file")
	var vertical = flag.Bool("v", false, "Print diagrams vertically instead of horizontally")

	var output io.Writer
	var err error

	flag.Parse()
	if *outputFile == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// if no mode is specified
	if *strIntervals == "" && *enumScales == 0 {
		fmt.Fprintln(output, "Please specify either -i or -e")
		fmt.Fprintln(output, "Use -h to display help")
		return
	}

	// general options
	data, err := os.ReadFile(*database)
	if err != nil {
		fmt.Fprintln(output, err)
	}
	err = json.Unmarshal(data, &scale.IntervalsAlias)
	if err != nil {
		fmt.Fprintln(output, err)
	}
	_ = outputFile

	// -i mode
	if *strIntervals != "" {
		if *enumScales != 0 {
			fmt.Fprintln(output, "-i and -e cannot be used together")
			return
		}
		if *tonic == 0 {
			fmt.Fprintln(output, "Please specify the root using -s")
			return
		}
		if *maxSemitones != 12 {
			fmt.Fprintln(output, "-f cannot be used with -s")
			return
		}
		if *active == "*" {
			*active = "111111111111"
		}
		s, err := scale.NewScale(strIntervals, active, *mode)
		if err != nil {
			fmt.Fprintln(output, err)
			return
		}
		fmt.Fprintf(output, "%s: ", s.Name)
		fmt.Fprintln(output, s.StrIntervals)
		fmt.Fprintln(output, )
		fb := fretboard.NewFretboard(*tuning, s, *tonic)
		if *vertical {
			fb.Printv(0, guitarstring.NeckLength, output)
		} else {
			fb.Print(0, guitarstring.NeckLength, output)
		}
		fmt.Fprintln(output, )
		fmt.Fprintln(output, )
		fmt.Fprintln(output, "Relative modes:")
		relatives := scale.RelativeScales(s.Intervals)
		for m, r := range relatives {
			fmt.Fprintf(output, "%d: %s %s\n", m+2, r.StrIntervals, r.Name)
		}
		if *enumChords {
			fmt.Fprintln(output, )
			fmt.Fprintln(output, )
			fmt.Fprintln(output, "Enumerating all 4 note chords from given scale:")
			fmt.Fprintln(output, )
			fb.PrintChords(*vertical, output)
		}
		if *buildChord != "" {
			if *chordRoot == 0 {
				fmt.Println(output, "Please specify chord root with -r")
				fmt.Fprintln(output, )
			} else {
				fmt.Fprintln(output, )
				fmt.Fprintln(output, )
				fmt.Fprintln(output, "Printing custom chord voicing")
				fmt.Fprintln(output, )
				fb.PrintChordVoicing(*chordRoot, buildChord, *inversion, output)
			}
		}
		return
	}

	// -e mode
	if *enumScales != 0 {
		if *strIntervals != "" {
			fmt.Fprintln(output, "-i and -e cannot be used together")
			return
		}
		if *tonic == 0 {
			fmt.Fprintln(output, "Please specify the root using -s")
			return
		}
		if *enumChords {
			fmt.Fprintln(output, "-c cannot be used with -e")
			return
		}
		if *active != "*" {
			fmt.Fprintln(output, "-a cannot be used with -e")
			return
		}
		if *mode != 1 {
			fmt.Fprintln(output, "-m cannot be used with -e")
		}
		var newScale scale.Scale
		scales := scale.EnumIntervals(*enumScales, *maxSemitones)
		for _, s := range scales {
			fmt.Fprintln(output, )
			fmt.Fprintln(output, s)
			newScale = scale.NewScaleFromIntervals(s)
			fb := fretboard.NewFretboard(*tuning, newScale, *tonic)
			if *vertical == false {
				fb.Print(0, guitarstring.NeckLength, output)
			} else {
				fb.Printv(0, guitarstring.NeckLength, output)
			}
		}
	}

}
