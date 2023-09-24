package guitarstring

import (
	"fmt"
)

type GuitarString struct {
	frets []int
	pitch string
	offset int
}

func (gs GuitarString) Print() {
	fmt.Printf("%s ||", gs.pitch)
	for _, degree := range gs.frets {
		if degree != 0 {
			fmt.Printf(" %2d|", degree)
		} else {
			fmt.Printf("   |")
		}
	}
	fmt.Println()
}

