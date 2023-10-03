package fretboard

import (
	"io"
	"fmt"
	"strings"
	"github.com/raygarner/lute/guitarstring"
	"github.com/raygarner/lute/scale"
)

const chordMaxWidth = 4
const chordDesiredNotes = 4

type Fretboard struct {
	strings []guitarstring.GuitarString
	scale scale.Scale
	tonic int
	tuning string
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

func countNotes(chord []int) int {
	var total = 0
	for _, i := range chord {
		if i != 0 {
			total++
		}
	}
	return total
}

// degrees = number of degrees in scale
// len(chord) = number of strings, each element is the degree to be played on that string
// strings = the number of strings yet to be handled
// returns a list of chords
// TODO: remove duplicates
func EnumerateChords(degrees int, chord []int, strings int) [][]int {
	var ret [][]int
	var newchord = make([]int, len(chord))

	if strings == 0 {
		if countNotes(chord) != chordDesiredNotes {
			return ret
		} else {
			return append(ret, chord)
		}
	}
	for d := 0; d <= degrees; d++ {
		copy(newchord, chord)
		ret = append(ret, EnumerateChords(degrees, append(newchord, d), strings-1)...)
	}
	return ret
}

func (fb Fretboard) PrintChords(vertical bool, output io.Writer) {
	var chord []int
	var chords = EnumerateChords(len(fb.scale.Intervals), chord, len(fb.strings))
	var tmpfb Fretboard
	var playable bool
	var lowest, highest int
	for _, c := range chords {
		tmpfb, playable, lowest, highest = fb.applyChord(c)
		if playable {
			//fmt.Println(c)
			fmt.Fprintln(output, c)
			if highest < chordMaxWidth {
				lowest = 0
				highest = chordMaxWidth
			} else if lowest > guitarstring.NeckLength - chordMaxWidth {
				highest = guitarstring.NeckLength
				lowest = guitarstring.NeckLength - chordMaxWidth
			} else if highest - lowest < chordMaxWidth {
				highest = lowest + chordMaxWidth
			}
			if vertical {
				//tmpfb.Printv(lowest, highest)
				tmpfb.Printv(lowest, highest, output)
			} else {
				//tmpfb.Print(lowest, highest)
				tmpfb.Print(lowest, highest, output)
			}
			fmt.Fprintf(output, "\n\n")
		}
	}
}

// TODO: tidy this return
func (fb Fretboard) applyChord(chord []int) (Fretboard,bool,int,int) {
	var newfb Fretboard
	newfb = NewFretboard(fb.tuning, fb.scale, fb.tonic)
	var lowest = 999
	var highest = -1
	var playable bool
	var size = 0
	
	for i, _ := range newfb.strings {
		for f, _ := range newfb.strings[i].Frets {
			if newfb.strings[i].Frets[f] != chord[i] {
				newfb.strings[i].Frets[f] = 0
			} else {
				if chord[i] != 0 {
					if f < lowest {
						lowest = f
					}
					if f > highest {
						highest = f
					}
					size++
				}
			}
		}
	}
	if (highest - lowest) > chordMaxWidth-1 || size != chordDesiredNotes {
		playable = false
	} else {
		playable = true
	}
	return newfb, playable, lowest, highest
}

func printFrets(lowest int, highest int, output io.Writer) {
	//fmt.Printf("     ")
	fmt.Fprintf(output, "     ")
	for i := lowest+1; i <= highest; i++ {
		//fmt.Printf(" %2d ", i)
		fmt.Fprintf(output, " %2d ", i)
	}
	//fmt.Println()
	fmt.Fprintln(output, )
}

func (fb Fretboard) Print(lowest int, highest int, output io.Writer) {
	for _, gs := range fb.strings {
		gs.Print(lowest, highest, output)
	}
	//fmt.Println()
	fmt.Fprintln(output,)
	printFrets(lowest, highest, output)
}

func (fb Fretboard) PrintRow(fret int, output io.Writer) {
	//fmt.Printf("%2d  |", fret+1)
	fmt.Fprintf(output, "%2d  |", fret+1)
	for i := len(fb.strings)-1; i >= 0; i-- {
		fb.strings[i].PrintFret(fret, "", output)
	}
	//fmt.Println()
	fmt.Fprintln(output, )
}

func (fb Fretboard) Printv(lowest int, highest int, output io.Writer) {
	//fmt.Printf("    ")
	fmt.Fprintf(output, "    ")
	for j := len(fb.strings)-1; j >= 0; j-- {
		//fmt.Printf(" %s", fb.strings[j].Pitch)
		fmt.Fprintf(output, " %s", fb.strings[j].Pitch)
	}
	//fmt.Println()
	fmt.Fprintln(output,)
	//fmt.Printf("    ")
	fmt.Fprintf(output, "    ")
	//fmt.Println(strings.Repeat("=", 3*len(fb.strings)+1))
	fmt.Fprintln(output, strings.Repeat("=", 3*len(fb.strings)+1))
	if lowest < 0 {
		lowest = 0
	}
	if highest > guitarstring.NeckLength {
		highest = guitarstring.NeckLength
	}
	for i := lowest; i < highest; i++ {
		fb.PrintRow(i, output)
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
			offset += guitarstring.NeckLength
		}
		offsets = append(offsets, offset)
		names = append(names, tuning[i:i+2])
	}
	return offsets, names
}

func NewFretboard(tuning string, s scale.Scale, tonic int) Fretboard {
	var fb Fretboard
	fb.scale = s
	fb.tonic = tonic
	fb.tuning = tuning
	offsets, pitches := buildOffsets(tuning)
	for i, offset := range offsets {
		fb.strings = append(fb.strings,
			guitarstring.NewGuitarString((tonic+offset) % guitarstring.NeckLength, s,
			pitches[i]))
	}
	return fb
}
