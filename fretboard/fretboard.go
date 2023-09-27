package fretboard

import (
	"github.com/raygarner/lute/guitarstring"
	"github.com/raygarner/lute/scale"
	"fmt"
	"strings"
)

const frets = 12

type Fretboard struct {
	strings []guitarstring.GuitarString
}

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

func printFrets() {
	fmt.Printf("     ")
	for i := 1; i <= frets; i++ {
		fmt.Printf(" %2d ", i)
	}
	fmt.Println()
}

func (fb Fretboard) Print() {
	for _, gs := range fb.strings {
		gs.Print()
	}
	fmt.Println()
	printFrets()
}

func (fb Fretboard) PrintRow(fret int) {
	fmt.Printf("%2d  |", fret+1)
	for i := len(fb.strings)-1; i >= 0; i-- {
		fb.strings[i].PrintFret(fret, "")
	}
	fmt.Println()
}

func (fb Fretboard) Printv() {
	fmt.Printf("    ")
	for j := len(fb.strings)-1; j >= 0; j-- {
		fmt.Printf(" %s", fb.strings[j].Pitch)
	}
	fmt.Println()
	fmt.Printf("    ")
	fmt.Println(strings.Repeat("=", 3*len(fb.strings)+1))
	for i := 0; i < frets; i++ {
		fb.PrintRow(i)
	}
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

func NewFretboard(tuning string, s scale.Scale, tonic int) Fretboard {
	var fb Fretboard
	offsets, pitches := buildOffsets(tuning)
	for i, offset := range offsets {
		fb.strings = append(fb.strings,
			guitarstring.NewGuitarString((tonic+offset) % frets, s,
			pitches[i]))
	}
	return fb
}
