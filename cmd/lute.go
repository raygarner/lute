package main

import (
	"flag"
	//"github.com/raygarner/lute/guitarstring"
	"github.com/raygarner/lute/scale"
	"github.com/raygarner/lute/fretboard"
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


func buildOffsets(tuning string) ([]int, []string) {
	var offsets []int
	var names []string
	var lowest = tuning[len(tuning)-2:]
	var offset int

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
	var active = flag.String("a", "111111111111", "which notes of the scale are active. One bit per degree, extra bits ignored")
	var tonic = flag.Int("s", 8, "fret of the tonic note on the lowest string")
	var tuning = flag.String("t", "enbngndnanen", "the tuning of the instrument in descending order of pitch (works for any number of strings)")
	flag.Parse()


	/*
	var fb []guitarstring.GuitarString
	//var fb fretboard.Fretboard
	var test fretboard.Fretboard
	_ = test
	s := scale.NewScale(strIntervals, active, *mode)
	offsets, string_names := buildOffsets(*tuning)
	for i, offset := range offsets {
		fb = append(fb,
			guitarstring.NewGuitarString((*tonic+offset) % frets, s,
			string_names[i]))
	}
	for _, gs := range fb {
		gs.Print()
	}
	*/

	s := scale.NewScale(strIntervals, active, *mode)
	fb := fretboard.NewFretboard(*tuning, s, *tonic)
	fb.Print()



}